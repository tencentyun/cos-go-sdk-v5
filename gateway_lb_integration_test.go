package cos

// 本文件是 Gateway DNS 负载均衡的集成测试：用假 DNS server 提供可控的解析结果，用假 Gateway
// 提供可控的响应，请求经由真实的 net/http 链路发出，因此覆盖到 lookupGatewayNodes /
// refreshAsync / RoundTrip 的真实行为，而不只是内存状态。
//
// These tests must not use t.Parallel because installFakeDNS replaces net.DefaultResolver.
// LB timings are immutable; tests advance per-entry state to simulate expiration.
//
// 节点 IP 统一取 192.0.2.0/24（RFC 5737 TEST-NET-1，文档与测试专用、不可路由）。这些地址只是
// 标签：实际 TCP 连接由 gatewayRouter 在 dial 层重定向到本地 httptest server，用例不产生任何
// 对外网络连接。

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	fakeDNSTypeA    = 1
	fakeDNSTypeAAAA = 28
)

// 三个用于承载假 Gateway 节点的地址。
const (
	nodeIPA = "192.0.2.1"
	nodeIPB = "192.0.2.2"
	nodeIPC = "192.0.2.3"
	// nodeIPv6A 用于验证 IPv6 后端节点（RFC 3849 文档专用地址段，不可路由）。
	nodeIPv6A = "2001:db8::1"
)

// nodeAddr 返回节点在改写后的连接目标形式（门禁保证端口恒为 80）。
func nodeAddr(ip string) string {
	return net.JoinHostPort(ip, "80")
}

// rapidTestHost 按用例名生成互不冲突的高性能桶短域名，避免包级注册表在用例间串味。统一带
// it 前缀，与单元测试文件里的域名彻底分开（注册表是包级单例，两个文件共享同一命名空间）。
func rapidTestHost(name string) string {
	return "it" + name + "-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
}

// ---------------------------------------------------------------------------
// 假 DNS server
// ---------------------------------------------------------------------------

// fakeDNSQuery 是 server 收到的一次查询。
type fakeDNSQuery struct {
	name  string
	qtype uint16
}

// fakeDNSServer 是一个最小 UDP DNS server，只处理单 question 的 A/AAAA 查询，答案可在运行期
// 替换——对齐 Python 集成测试里可变的 selected_ip，用来模拟扩缩容与 DNS 故障。
type fakeDNSServer struct {
	conn *net.UDPConn
	done chan struct{}

	mu      sync.Mutex
	answers map[string][]string // 域名（小写、无末尾点）-> A/AAAA 记录，按地址族自动分类
	failAll bool                // 为真时对所有查询立即返回 SERVFAIL
	silent  bool                // 为真时对所有查询不做任何应答，模拟丢包超时（区别于 failAll）
	queries []fakeDNSQuery
}

func startFakeDNS(t *testing.T) *fakeDNSServer {
	t.Helper()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatalf("net.ListenUDP() failed: %v", err)
	}
	s := &fakeDNSServer{
		conn:    conn,
		done:    make(chan struct{}),
		answers: make(map[string][]string),
	}
	go s.serve()
	t.Cleanup(func() {
		conn.Close()
		<-s.done
	})
	return s
}

// setAnswer 替换某域名的记录集合（IPv4 走 A、IPv6 走 AAAA，按地址族自动分类），传空表示该
// 域名存在但无记录。
func (s *fakeDNSServer) setAnswer(host string, ips ...string) {
	s.mu.Lock()
	s.answers[canonicalHost(host)] = ips
	s.failAll = false
	s.silent = false
	s.mu.Unlock()
}

// setFailure 让后续所有查询立即返回 SERVFAIL，模拟 DNS 服务器主动拒绝。
func (s *fakeDNSServer) setFailure() {
	s.mu.Lock()
	s.failAll = true
	s.silent = false
	s.mu.Unlock()
}

// setSilent 让后续所有查询都不返回任何应答，模拟 DNS 请求丢包——用于验证真正的超时路径，
// 与 setFailure 的"立即拒绝"是两种不同的故障形态。
func (s *fakeDNSServer) setSilent() {
	s.mu.Lock()
	s.failAll = false
	s.silent = true
	s.mu.Unlock()
}

// queriesFor 返回针对某域名的查询记录。
func (s *fakeDNSServer) queriesFor(host string) []fakeDNSQuery {
	want := canonicalHost(host)
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []fakeDNSQuery
	for _, q := range s.queries {
		if q.name == want {
			out = append(out, q)
		}
	}
	return out
}

func (s *fakeDNSServer) resetQueries() {
	s.mu.Lock()
	s.queries = nil
	s.mu.Unlock()
}

func (s *fakeDNSServer) serve() {
	defer close(s.done)
	buf := make([]byte, 512)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			return // conn 已关闭
		}
		resp, ok := s.buildResponse(buf[:n])
		if !ok {
			continue
		}
		s.conn.WriteToUDP(resp, addr)
	}
}

// buildResponse 按查询构造应答：复制原 header 与 question，再按记录追加 answer。
// silent 时不返回任何字节（模拟丢包超时）；否则按 qtype 只挑出对应地址族的记录作答，
// AAAA 查询在答案里没有 IPv6 地址时会正常返回零应答（即"域名存在但无 AAAA 记录"）。
func (s *fakeDNSServer) buildResponse(query []byte) ([]byte, bool) {
	name, qtype, end, ok := parseFakeDNSQuestion(query)
	if !ok {
		return nil, false
	}

	s.mu.Lock()
	s.queries = append(s.queries, fakeDNSQuery{name: name, qtype: qtype})
	silent := s.silent
	failAll := s.failAll
	ips := s.answers[name]
	s.mu.Unlock()

	if silent {
		return nil, false // 丢包：不回任何应答，调用方最终会等到自己的超时
	}

	resp := make([]byte, end, 512)
	copy(resp, query[:end])
	resp[2], resp[3] = 0x81, 0x80 // QR=1, RD=1, RA=1, RCODE=0
	if failAll {
		resp[3] = 0x82 // RCODE=2 SERVFAIL
	}
	resp[4], resp[5] = 0, 1 // QDCOUNT=1

	var answers int
	if !failAll {
		for _, ip := range ips {
			parsed := net.ParseIP(ip)
			if parsed == nil {
				continue
			}
			switch qtype {
			case fakeDNSTypeA:
				v4 := parsed.To4()
				if v4 == nil {
					continue // 该记录是纯 IPv6 地址，不出现在 A 应答里
				}
				resp = append(resp,
					0xC0, 0x0C, // NAME：指向 offset 12 的 QNAME
					0x00, 0x01, // TYPE=A
					0x00, 0x01, // CLASS=IN
					0x00, 0x00, 0x00, 0x1E, // TTL=30
					0x00, 0x04, // RDLENGTH=4
				)
				resp = append(resp, v4...)
				answers++
			case fakeDNSTypeAAAA:
				if parsed.To4() != nil {
					continue // 该记录是 IPv4 地址，不出现在 AAAA 应答里
				}
				resp = append(resp,
					0xC0, 0x0C, // NAME：指向 offset 12 的 QNAME
					0x00, 0x1C, // TYPE=AAAA
					0x00, 0x01, // CLASS=IN
					0x00, 0x00, 0x00, 0x1E, // TTL=30
					0x00, 0x10, // RDLENGTH=16
				)
				resp = append(resp, parsed.To16()...)
				answers++
			}
		}
	}
	resp[6], resp[7] = byte(answers>>8), byte(answers)
	resp[8], resp[9], resp[10], resp[11] = 0, 0, 0, 0 // NSCOUNT/ARCOUNT

	return resp, true
}

// parseFakeDNSQuestion 解析查询里的 QNAME/QTYPE，并返回 question 段的结束偏移。
// 查询报文里不会出现压缩指针，因此只处理普通 label。
func parseFakeDNSQuestion(msg []byte) (name string, qtype uint16, end int, ok bool) {
	if len(msg) < 12 {
		return "", 0, 0, false
	}
	var labels []string
	i := 12
	for {
		if i >= len(msg) {
			return "", 0, 0, false
		}
		length := int(msg[i])
		i++
		if length == 0 {
			break
		}
		if length > 63 || i+length > len(msg) {
			return "", 0, 0, false
		}
		labels = append(labels, string(msg[i:i+length]))
		i += length
	}
	if i+4 > len(msg) {
		return "", 0, 0, false
	}
	qtype = uint16(msg[i])<<8 | uint16(msg[i+1])
	return strings.ToLower(strings.Join(labels, ".")), qtype, i + 4, true
}

// installFakeDNS 把包级 net.DefaultResolver 指向假 server。PreferGo 必须为真，否则某些平台会
// 走 cgo resolver、绕开这里的 Dial。
func installFakeDNS(t *testing.T, s *fakeDNSServer) {
	t.Helper()
	original := net.DefaultResolver
	target := s.conn.LocalAddr().String()
	// 内层 dialer 必须显式给一个非 nil 的 Resolver（哪怕是零值），否则它在解析 target 这个
	// 字面量地址时仍会去读全局 net.DefaultResolver——这会跟测试结束时的 Cleanup 写回原值
	// 产生数据竞态（target 已经是字面量 "ip:port"，本不需要经过任何 DNS 解析）。
	innerDialer := &net.Dialer{Resolver: &net.Resolver{}}
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return innerDialer.DialContext(ctx, "udp", target)
		},
	}
	t.Cleanup(func() { net.DefaultResolver = original })
}

// ---------------------------------------------------------------------------
// 假 Gateway 节点
// ---------------------------------------------------------------------------

// gatewayRecord 是假 Gateway 收到的一次请求，字段选取对齐 Python 集成测试的 records。
type gatewayRecord struct {
	host       string // 服务端实际读到的 Host 头，用于验证 Host/签名不变量
	path       string
	rawQuery   string
	remoteAddr string // 客户端源地址，源端口相同即说明连接被复用
	body       []byte // 服务端实际收到的请求体，用于验证写请求 body 传输完整
}

// fakeGateway 是一个记录请求并可动态改变响应状态的假 Gateway 节点。
type fakeGateway struct {
	name   string
	server *httptest.Server

	mu      sync.Mutex
	status  int
	delay   time.Duration // 响应前的人工延迟，用于构造"请求卡住"场景（取消/超时测试）
	records []gatewayRecord
}

func startFakeGateway(t *testing.T, name string) *fakeGateway {
	t.Helper()
	g := &fakeGateway{name: name, status: http.StatusOK}
	g.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		g.mu.Lock()
		g.records = append(g.records, gatewayRecord{
			host:       r.Host,
			path:       r.URL.Path,
			rawQuery:   r.URL.RawQuery,
			remoteAddr: r.RemoteAddr,
			body:       body,
		})
		status := g.status
		delay := g.delay
		g.mu.Unlock()
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-r.Context().Done():
				return // 客户端已取消/超时断开，不再继续处理这次请求
			}
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(status)
		w.Write([]byte(g.name))
	}))
	t.Cleanup(g.server.Close)
	return g
}

// setStatus 切换该节点后续响应的状态码，用于模拟节点故障与恢复。
func (g *fakeGateway) setStatus(status int) {
	g.mu.Lock()
	g.status = status
	g.mu.Unlock()
}

// setDelay 让该节点后续的响应延迟指定时长，用于构造"请求卡住"场景，配合调用方的取消/超时
// 验证节点是否被正确地标记/不标记为故障。
func (g *fakeGateway) setDelay(d time.Duration) {
	g.mu.Lock()
	g.delay = d
	g.mu.Unlock()
}

func (g *fakeGateway) snapshot() []gatewayRecord {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]gatewayRecord(nil), g.records...)
}

func (g *fakeGateway) requestCount() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.records)
}

func (g *fakeGateway) reset() {
	g.mu.Lock()
	g.records = nil
	g.mu.Unlock()
}

// addr 返回该假 Gateway 的真实监听地址。
func (g *fakeGateway) addr(t *testing.T) string {
	t.Helper()
	u, err := url.Parse(g.server.URL)
	if err != nil {
		t.Fatalf("url.Parse(%q) failed: %v", g.server.URL, err)
	}
	return u.Host
}

// ---------------------------------------------------------------------------
// dial 层路由
// ---------------------------------------------------------------------------

// gatewayRouter 在 dial 层把"节点 IP:80"重定向到本地假 Gateway，并记录每次连接的目标地址。
// 门禁强制端口为 80，无法在快照里塞入 httptest 的随机端口，因此重定向必须发生在 dial 层——
// 这样也不需要改动门禁或 URL 改写逻辑。
type gatewayRouter struct {
	mu     sync.Mutex
	routes map[string]string
	dials  []string
	tr     *http.Transport
}

func newGatewayRouter() *gatewayRouter {
	return &gatewayRouter{routes: make(map[string]string)}
}

// route 把某个连接目标（形如 "192.0.2.1:80"）指向一个假 Gateway。
func (r *gatewayRouter) route(t *testing.T, addr string, g *fakeGateway) {
	t.Helper()
	r.mu.Lock()
	r.routes[addr] = g.addr(t)
	r.mu.Unlock()
}

// unroute 移除路由，使连接该地址时 dial 失败，模拟节点彻底不可达。
func (r *gatewayRouter) unroute(addr string) {
	r.mu.Lock()
	delete(r.routes, addr)
	r.mu.Unlock()
}

func (r *gatewayRouter) dialedAddrs() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.dials...)
}

func (r *gatewayRouter) resetDials() {
	r.mu.Lock()
	r.dials = nil
	r.mu.Unlock()
}

// transport 返回下一跳，同一个 router 始终复用同一个 Transport（对齐真实进程只有一个连接池的
// 情形，也让用例能通过 closeIdleConnections 精确控制连接池状态）。Proxy 保持 nil，使
// hasProxy() 走"Proxy == nil 即无代理"的确定分支，不受执行环境的代理环境变量影响。
func (r *gatewayRouter) transport() *http.Transport {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tr != nil {
		return r.tr
	}
	r.tr = &http.Transport{
		Proxy: nil,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			r.mu.Lock()
			r.dials = append(r.dials, addr)
			target, ok := r.routes[addr]
			r.mu.Unlock()
			if !ok {
				return nil, fmt.Errorf("fake router: no route for %s", addr)
			}
			return (&net.Dialer{}).DialContext(ctx, network, target)
		},
	}
	return r.tr
}

// closeIdleConnections 关闭连接池里的空闲连接，连带结束它们的读写 goroutine。goroutine 计数
// 类断言需要先排除 keep-alive 连接这个变量。
func (r *gatewayRouter) closeIdleConnections() {
	r.transport().CloseIdleConnections()
}

// ---------------------------------------------------------------------------
// 通用辅助
// ---------------------------------------------------------------------------

// forceSnapshotExpired advances only the entry's state; production timing constants remain immutable.
func forceSnapshotExpired(e *gatewayLBEntry) {
	e.mu.Lock()
	e.expireAt = time.Now().Add(-time.Second)
	e.mu.Unlock()
}

// forceCooldownExpired makes selected nodes eligible without waiting for the production cooldown.
func forceCooldownExpired(e *gatewayLBEntry, ips ...string) {
	e.mu.Lock()
	for _, ip := range ips {
		if _, ok := e.cooldown[ip]; ok {
			e.cooldown[ip] = time.Now().Add(-time.Second)
		}
	}
	e.mu.Unlock()
}

// newGatewayLBClient 构造一个经过负载均衡层的客户端。
func newGatewayLBClient(r *gatewayRouter) *http.Client {
	return &http.Client{Transport: newRapidGatewayLBTransport(r.transport())}
}

// sendGatewayRequest 发一次请求并读完 body（读完才会把连接放回池里，连接复用断言依赖这一点）。
func sendGatewayRequest(client *http.Client, host, path string) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, "http://"+host+path, nil)
	if err != nil {
		return 0, "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(body), nil
}

// sendGatewayRequestCtx 与 sendGatewayRequest 相同，但允许调用方提供自定义 ctx——用于验证
// 取消/超时场景下节点是否被正确地标记/不标记为故障。
func sendGatewayRequestCtx(client *http.Client, ctx context.Context, host, path string) (int, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+host+path, nil)
	if err != nil {
		return 0, "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(body), nil
}

// sendGatewayRequestBody 发一次带 body 的请求（PUT/POST/DELETE），用于验证 cloneRequest 与
// URL 改写不会干扰请求体的传输——这条性质此前的用例全部用 GET，从未被验证过。
func sendGatewayRequestBody(client *http.Client, host, path, method string, body []byte) (int, string, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, "http://"+host+path, reader)
	if err != nil {
		return 0, "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(respBody), nil
}

// ipStrings 把快照转成便于比较与打印的字符串集合。
func ipStrings(nodes []net.IP) []string {
	out := make([]string, 0, len(nodes))
	for _, ip := range nodes {
		out = append(out, ip.String())
	}
	return out
}

func sameStringSet(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := make(map[string]int, len(got))
	for _, v := range got {
		seen[v]++
	}
	for _, v := range want {
		if seen[v] == 0 {
			return false
		}
		seen[v]--
	}
	return true
}

// snapshotNodes returns a copy of the currently published snapshot.
func snapshotNodes(e *gatewayLBEntry) []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return ipStrings(e.nodes)
}

// waitForResolutionIdle waits until the current background lookup completes.
func waitForResolutionIdle(e *gatewayLBEntry, deadline time.Time) bool {
	for time.Now().Before(deadline) {
		e.mu.Lock()
		resolving := e.resolving
		e.mu.Unlock()
		if !resolving {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return false
}

// waitForSnapshot starts a fresh lookup after any previous lookup completes and
// verifies the snapshot published by that lookup rather than a stale snapshot.
func waitForSnapshot(t *testing.T, e *gatewayLBEntry, host string, want ...string) {
	t.Helper()
	lookupHost := canonicalHost(host)
	deadline := time.Now().Add(3 * time.Second)
	var got []string
	for time.Now().Before(deadline) {
		if !waitForResolutionIdle(e, deadline) {
			break
		}
		e.mu.Lock()
		e.nodes = nil
		e.expireAt = time.Now().Add(-time.Second)
		e.mu.Unlock()
		e.resolve(lookupHost)
		if !waitForResolutionIdle(e, deadline) {
			break
		}
		got = snapshotNodes(e)
		if sameStringSet(got, want) {
			return
		}
	}
	t.Fatalf("snapshot = %v, want %v within timeout", got, want)
}

// waitForNoResolutionInFlight waits for background lookup cleanup between tests.
func waitForNoResolutionInFlight(t *testing.T, e *gatewayLBEntry) {
	t.Helper()
	if !waitForResolutionIdle(e, time.Now().Add(3*time.Second)) {
		t.Fatal("a background resolution is still in flight")
	}
}

// cooldownKeys 返回当前处于冷却记录中的节点。
func cooldownKeys(e *gatewayLBEntry) []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]string, 0, len(e.cooldown))
	for k := range e.cooldown {
		out = append(out, k)
	}
	return out
}

// nodeFailureState returns the current cooldown deadline and failure streak.
func nodeFailureState(e *gatewayLBEntry, ip string) (time.Time, int, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	until, ok := e.cooldown[ip]
	return until, e.failStreak[ip], ok
}

// cooldownUntil returns the current cooldown deadline for an IP.
func cooldownUntil(e *gatewayLBEntry, ip string) time.Time {
	until, _, _ := nodeFailureState(e, ip)
	return until
}

// addCooldown preserves existing state while excluding selected nodes from pick.
func addCooldown(e *gatewayLBEntry, ips ...string) {
	e.mu.Lock()
	if e.cooldown == nil {
		e.cooldown = make(map[string]time.Time)
	}
	for _, ip := range ips {
		e.cooldown[ip] = time.Now().Add(gatewayLBDefaultMaxCooldown)
	}
	e.mu.Unlock()
}

// setCooldown replaces the candidate exclusions for tests that require exact control.
func setCooldown(e *gatewayLBEntry, ips ...string) {
	e.mu.Lock()
	e.cooldown = make(map[string]time.Time, len(ips))
	for _, ip := range ips {
		e.cooldown[ip] = time.Now().Add(gatewayLBDefaultMaxCooldown)
	}
	e.mu.Unlock()
}

// clearCooldownFor 只移除指定节点的冷却记录，保留其余记录——用于在观察生产代码自己设置的
// 冷却时，撤掉测试为"强制首选"而人为加的那一条。
func clearCooldownFor(e *gatewayLBEntry, ips ...string) {
	e.mu.Lock()
	for _, ip := range ips {
		delete(e.cooldown, ip)
	}
	e.mu.Unlock()
}

// aQueryCount returns observed A queries. Successful lookups normally issue one
// A query, while the Go resolver may retry SERVFAIL within a single lookup.
func aQueryCount(dns *fakeDNSServer, host string) int {
	var n int
	for _, q := range dns.queriesFor(host) {
		if q.qtype == fakeDNSTypeA {
			n++
		}
	}
	return n
}

// anyAQuery 报告上次 resetQueries 之后假 DNS 是否收到过任意 A 查询（不限定查询的域名）。
// 只统计逻辑域名会引入环境依赖：resolv.conf 的 search/ndots 配置会让 Go 解析器先查询带搜索
// 后缀的候选名，silent 模式下这些查询足以耗尽整个 ctx 超时预算，逻辑域名本身的查询可能根本
// 来不及发出。判定"解析确实发生过"应按"发过任意 A 查询"。
func anyAQuery(dns *fakeDNSServer) bool {
	dns.mu.Lock()
	defer dns.mu.Unlock()
	for _, q := range dns.queries {
		if q.qtype == fakeDNSTypeA {
			return true
		}
	}
	return false
}

// registerHost 让假 DNS 应答某域名，并清理该域名对应的注册表项。
func registerHost(t *testing.T, dns *fakeDNSServer, host string, ips ...string) *gatewayLBEntry {
	t.Helper()
	dns.setAnswer(host, ips...)
	resetGatewayLBEntry(host)
	entry := entryFor(canonicalHost(host))
	t.Cleanup(func() {
		waitForNoResolutionInFlight(t, entry)
		resetGatewayLBEntry(host)
	})
	return entry
}

// setupGatewayLB 准备一次集成测试所需的环境：假 DNS、dial 路由与经过负载均衡层的客户端。
func setupGatewayLB(t *testing.T, host string, ips ...string) (*fakeDNSServer, *gatewayRouter, *http.Client, *gatewayLBEntry) {
	t.Helper()
	dns := startFakeDNS(t)
	installFakeDNS(t, dns)
	entry := registerHost(t, dns, host, ips...)
	router := newGatewayRouter()
	return dns, router, newGatewayLBClient(router), entry
}

// ---------------------------------------------------------------------------
// A 组：负载均衡效果
// ---------------------------------------------------------------------------

// 统一构造入口必须让 CreateSession 和数据面请求经过同一个 Gateway LB，同时分别使用基础凭证
// 与 session 凭证签名。该用例覆盖完整的真实 HTTP 链路，而不只检查 Transport 字段层级。
func TestNewRapidCredentialTransportIntegrationSessionAndDataUseGatewayLB(t *testing.T) {
	host := rapidTestHost("rapidtransport")
	dns := startFakeDNS(t)
	installFakeDNS(t, dns)
	entry := registerHost(t, dns, host, nodeIPA)
	router := newGatewayRouter()
	defer router.closeIdleConnections()

	var sessionCalls, dataCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != host {
			t.Errorf("Host = %q, want %q", r.Host, host)
		}
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			atomic.AddInt32(&sessionCalls, 1)
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=AKIDBASE") {
				t.Errorf("CreateSession Authorization = %q, want base credential", got)
			}
			fmt.Fprint(w, sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Now().Add(time.Hour)))
		case r.URL.Path == "/key":
			atomic.AddInt32(&dataCalls, 1)
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("data Authorization = %q, want session credential", got)
			}
			if got := r.Header.Get("x-cos-security-token"); got != "session-token" {
				t.Errorf("data session token = %q, want session-token", got)
			}
			fmt.Fprint(w, "data")
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	router.route(t, nodeAddr(nodeIPA), &fakeGateway{server: server})
	waitForSnapshot(t, entry, host, nodeIPA)

	bucketURL, err := url.Parse("http://" + host)
	if err != nil {
		t.Fatal(err)
	}
	base := &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret", Transport: router.transport()}
	client := NewClient(&BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: NewRapidCredentialTransport(base),
	})
	client.Conf.RetryOpt.Count = 0

	resp, err := client.Object.Get(context.Background(), "key", nil)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "data" {
		t.Fatalf("data response body = %q, want data", body)
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 1 {
		t.Fatalf("CreateSession requests = %d, want 1", got)
	}
	if got := atomic.LoadInt32(&dataCalls); got != 1 {
		t.Fatalf("data requests = %d, want 1", got)
	}
}

func TestGatewayLBIntegrationDistributesAcrossNodes(t *testing.T) {
	host := rapidTestHost("spread")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// 复用同一个 client（连接池已建立连接）连发多次：若选点发生在连接层，所有请求都会沿用
	// 首次建立的那条连接、只打到一个节点；两个节点都收到请求即证明选点是请求级的。
	const rounds = 200
	for i := 0; i < rounds; i++ {
		status, body, err := sendGatewayRequest(client, host, "/key")
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		if status != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i, status)
		}
		if body != "A" && body != "B" {
			t.Fatalf("request %d body = %q, want a gateway name", i, body)
		}
	}

	countA, countB := gwA.requestCount(), gwB.requestCount()
	if countA+countB != rounds {
		t.Fatalf("total requests = %d+%d, want %d", countA, countB, rounds)
	}
	if countA < rounds*3/10 || countA > rounds*7/10 {
		t.Errorf("distribution A=%d B=%d out of %d, want roughly even", countA, countB, rounds)
	}
}

func TestGatewayLBIntegrationForcedPreferredNode(t *testing.T) {
	host := rapidTestHost("forced")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// 逐个把每个节点强制为唯一候选，各自连发 3 次请求都应成功且落在该节点上：证明每个节点都
	// 单独可用，而不是靠统计分布侧面推断。
	cases := []struct {
		preferred string
		cooled    string
		want      string
		gw        *fakeGateway
		other     *fakeGateway
	}{
		{preferred: nodeIPA, cooled: nodeIPB, want: "A", gw: gwA, other: gwB},
		{preferred: nodeIPB, cooled: nodeIPA, want: "B", gw: gwB, other: gwA},
	}
	for _, tc := range cases {
		t.Run("prefer_"+tc.want, func(t *testing.T) {
			gwA.reset()
			gwB.reset()
			setCooldown(entry, tc.cooled)

			for i := 0; i < 3; i++ {
				status, body, err := sendGatewayRequest(client, host, "/key")
				if err != nil {
					t.Fatalf("request %d failed: %v", i, err)
				}
				if status != http.StatusOK || body != tc.want {
					t.Fatalf("request %d = (%d, %q), want (200, %q)", i, status, body, tc.want)
				}
			}
			if got := tc.gw.requestCount(); got != 3 {
				t.Errorf("preferred gateway got %d requests, want 3", got)
			}
			if got := tc.other.requestCount(); got != 0 {
				t.Errorf("cooling gateway got %d requests, want 0", got)
			}
		})
	}
}

func TestGatewayLBIntegrationKeepsLogicalHostHeader(t *testing.T) {
	host := rapidTestHost("hosthdr")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA)

	gw := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gw)
	waitForSnapshot(t, entry, host, nodeIPA)

	if _, body, err := sendGatewayRequest(client, host, "/dir/key?versionId=42"); err != nil || body != "A" {
		t.Fatalf("request = (%q, %v), want (\"A\", nil)", body, err)
	}

	// 只有真实服务端才能观测标准库最终写出的 Host 头，这是本条必须做成集成测试的原因。
	records := gw.snapshot()
	if len(records) != 1 {
		t.Fatalf("gateway got %d requests, want 1", len(records))
	}
	if records[0].host != host {
		t.Errorf("server-side Host = %q, want the logical host %q", records[0].host, host)
	}
	if records[0].path != "/dir/key" || records[0].rawQuery != "versionId=42" {
		t.Errorf("server-side target = %q?%q, want /dir/key?versionId=42",
			records[0].path, records[0].rawQuery)
	}
	// 连接目标确实是物理节点，说明上面的 Host 头不是"没改写"造成的假象。
	if dialed := router.dialedAddrs(); len(dialed) == 0 || dialed[0] != nodeAddr(nodeIPA) {
		t.Errorf("dialed = %v, want the first dial to target %s", dialed, nodeAddr(nodeIPA))
	}
}

func TestGatewayLBIntegrationSharesPoolAndKeepsHostsSeparate(t *testing.T) {
	hostA := rapidTestHost("poolone")
	hostB := rapidTestHost("pooltwo")

	dns, router, client, entryA := setupGatewayLB(t, hostA, nodeIPA)
	entryB := registerHost(t, dns, hostB, nodeIPA) // 两个桶域名解析到同一个 Gateway 节点

	gw := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gw)
	waitForSnapshot(t, entryA, hostA, nodeIPA)
	waitForSnapshot(t, entryB, hostB, nodeIPA)

	for _, host := range []string{hostA, hostB} {
		if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
			t.Fatalf("request to %s = (%q, %v), want (\"A\", nil)", host, body, err)
		}
	}

	records := gw.snapshot()
	if len(records) != 2 {
		t.Fatalf("gateway got %d requests, want 2", len(records))
	}
	// 按 (scheme, ip, port) 分池：两个不同桶域名改写到同一节点后共用一条连接，源端口相同。
	if records[0].remoteAddr != records[1].remoteAddr {
		t.Errorf("remote addrs = %q and %q, want the same connection to be reused",
			records[0].remoteAddr, records[1].remoteAddr)
	}
	// 共用连接不能让 Host 头串味：每个请求仍须携带自己的逻辑域名。
	if records[0].host != hostA || records[1].host != hostB {
		t.Errorf("Host headers = %q, %q, want %q, %q",
			records[0].host, records[1].host, hostA, hostB)
	}
}

func TestGatewayLBIntegrationIPv6NodeIsUsable(t *testing.T) {
	host := rapidTestHost("ipv6")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPv6A)

	gw := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPv6A), gw)
	waitForSnapshot(t, entry, host, nodeIPv6A)

	// 之前的 20 条用例全部用 IPv4 节点；快照里混入 IPv6 地址是 10.2 节记录的已知残余风险，
	// 但从未在集成层面验证过"选中一个 IPv6 节点时，改写、拨号、响应这条链路本身是否走得通"。
	status, body, err := sendGatewayRequest(client, host, "/key")
	if err != nil || status != http.StatusOK || body != "A" {
		t.Fatalf("request = (%d, %q, %v), want (200, \"A\", nil)", status, body, err)
	}
	// net.JoinHostPort 会给 IPv6 地址自动加方括号；确认改写后的连接目标是这个合法形式。
	if dialed := router.dialedAddrs(); len(dialed) != 1 || dialed[0] != nodeAddr(nodeIPv6A) {
		t.Errorf("dialed = %v, want %s", dialed, nodeAddr(nodeIPv6A))
	}
}

func TestGatewayLBIntegrationWriteRequestBodyIsIntact(t *testing.T) {
	host := rapidTestHost("writebody")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA)

	gw := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gw)
	waitForSnapshot(t, entry, host, nodeIPA)

	// PUT/POST/DELETE 是常见的写路径；此前全部 20 条用例只发过 GET，从未验证过 cloneRequest +
	// URL 改写不会干扰请求体的传输。
	payload := []byte("integration test payload for the write path")
	status, _, err := sendGatewayRequestBody(client, host, "/key", http.MethodPut, payload)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("PUT status = %d, want 200", status)
	}

	records := gw.snapshot()
	if len(records) != 1 {
		t.Fatalf("gateway got %d requests, want 1", len(records))
	}
	if string(records[0].body) != string(payload) {
		t.Errorf("server-side body = %q, want %q", records[0].body, payload)
	}
	if records[0].host != host {
		t.Errorf("server-side Host = %q, want the logical host %q", records[0].host, host)
	}
	if dialed := router.dialedAddrs(); len(dialed) != 1 || dialed[0] != nodeAddr(nodeIPA) {
		t.Errorf("dialed = %v, want %s", dialed, nodeAddr(nodeIPA))
	}
}

// ---------------------------------------------------------------------------
// B 组：错误处理与故障切换
// ---------------------------------------------------------------------------

func TestGatewayLBIntegrationFailsOverOn5xx(t *testing.T) {
	host := rapidTestHost("failover5xx")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	gwA.setStatus(http.StatusBadGateway)
	setCooldown(entry, nodeIPB) // 强制首选故障节点 A

	status, body, err := sendGatewayRequest(client, host, "/key")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	// 5xx 必须原样上抛，由上层决定是否重试；本层不做内部重试。
	if status != http.StatusBadGateway || body != "A" {
		t.Fatalf("request = (%d, %q), want (502, \"A\")", status, body)
	}
	if gwB.requestCount() != 0 {
		t.Fatalf("gateway B got %d requests, want 0 (no internal retry)", gwB.requestCount())
	}
	if keys := cooldownKeys(entry); !containsString(keys, nodeIPA) {
		t.Fatalf("cooldown = %v, want node A to be marked after the 502", keys)
	}

	// 撤掉测试人为加的 B 冷却，只保留生产代码自己给 A 设的那条，观察后续流量走向。
	clearCooldownFor(entry, nodeIPB)
	gwA.reset()
	gwB.reset()
	for i := 0; i < 10; i++ {
		status, body, err := sendGatewayRequest(client, host, "/key")
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		if status != http.StatusOK || body != "B" {
			t.Fatalf("request %d = (%d, %q), want (200, \"B\")", i, status, body)
		}
	}
	if gwA.requestCount() != 0 {
		t.Errorf("cooling gateway A got %d requests, want 0", gwA.requestCount())
	}
}

func TestGatewayLBIntegrationFailsOverOnDialFailure(t *testing.T) {
	host := rapidTestHost("failoverdial")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// 移除路由使连接直接失败，模拟节点彻底不可达（比返回 5xx 更底层的故障）。
	router.unroute(nodeAddr(nodeIPA))
	setCooldown(entry, nodeIPB)

	if _, _, err := sendGatewayRequest(client, host, "/key"); err == nil {
		t.Fatal("request succeeded, want a dial error to be propagated")
	}
	if keys := cooldownKeys(entry); !containsString(keys, nodeIPA) {
		t.Fatalf("cooldown = %v, want node A to be marked after the dial failure", keys)
	}

	clearCooldownFor(entry, nodeIPB)
	for i := 0; i < 10; i++ {
		if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "B" {
			t.Fatalf("request %d = (%q, %v), want (\"B\", nil)", i, body, err)
		}
	}
	if gwA.requestCount() != 0 {
		t.Errorf("unreachable gateway A got %d requests, want 0", gwA.requestCount())
	}
}

func TestGatewayLBIntegrationRecoversAfterCooldownExpires(t *testing.T) {
	host := rapidTestHost("recover")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// A 故障并被冷却。
	gwA.setStatus(http.StatusServiceUnavailable)
	setCooldown(entry, nodeIPB)
	if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusServiceUnavailable {
		t.Fatalf("request = (%d, %v), want (503, nil)", status, err)
	}
	until, streak, ok := nodeFailureState(entry, nodeIPA)
	if !ok || !until.After(time.Now()) || streak != 1 {
		t.Fatalf("node A failure state = (deadline=%v, streak=%d, present=%v), want a future cooldown with streak 1", until, streak, ok)
	}
	clearCooldownFor(entry, nodeIPB)

	// Once A recovers and its cooldown expires, it must become eligible again.
	gwA.setStatus(http.StatusOK)
	forceCooldownExpired(entry, nodeIPA)
	addCooldown(entry, nodeIPB)
	if status, body, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusOK || body != "A" {
		t.Fatalf("request after A cooldown = (%d, %q, %v), want (200, \"A\", nil)", status, body, err)
	}
	clearCooldownFor(entry, nodeIPB)

	gwA.reset()
	gwB.reset()
	for i := 0; i < 100; i++ {
		if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusOK {
			t.Fatalf("request %d = (%d, %v), want (200, nil)", i, status, err)
		}
	}
	if gwA.requestCount() == 0 {
		t.Error("gateway A never selected after its cooldown expired")
	}
	if gwB.requestCount() == 0 {
		t.Error("gateway B never selected, want both nodes to be fair candidates again")
	}
}

// TestGatewayLBIntegrationBackoffEscalatesAcrossRequests 端到端验证指数退避：反复对同一节点
// 触发失败，冷却时长应逐次翻倍直到封顶；期间穿插一次成功，验证 streak 被清零、下次失败重新
// 从 BaseCooldown 起算。覆盖的是 RoundTrip -> markFailed/markSucceeded 这条真实调用链路，
// 不只是直接调用 gatewayLBEntry 方法的白盒单测。
func TestGatewayLBIntegrationBackoffEscalatesAcrossRequests(t *testing.T) {
	host := rapidTestHost("backoffescalate")

	dns := startFakeDNS(t)
	installFakeDNS(t, dns)
	entry := registerHost(t, dns, host, nodeIPA, nodeIPB)
	router := newGatewayRouter()
	client := newGatewayLBClient(router)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	gwA.setStatus(http.StatusServiceUnavailable)

	// Verify the production sequence while advancing A's deadline between rounds.
	want := []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second, 40 * time.Second, 60 * time.Second}
	for i, w := range want {
		addCooldown(entry, nodeIPB)
		before := time.Now()
		if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusServiceUnavailable {
			t.Fatalf("failure #%d request = (%d, %v), want (503, nil)", i+1, status, err)
		}
		got := cooldownUntil(entry, nodeIPA).Sub(before)
		if got < w-time.Second || got > w+time.Second {
			t.Errorf("failure #%d cooldown = %v, want approximately %v", i+1, got, w)
		}
		clearCooldownFor(entry, nodeIPB)
		forceCooldownExpired(entry, nodeIPA)
	}

	// A 恢复健康：一次成功应清空连续失败计数，下次失败重新从 BaseCooldown 起算。
	gwA.setStatus(http.StatusOK)
	addCooldown(entry, nodeIPB)
	if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusOK {
		t.Fatalf("recovery request = (%d, %v), want (200, nil)", status, err)
	}
	clearCooldownFor(entry, nodeIPB)

	gwA.setStatus(http.StatusServiceUnavailable)
	addCooldown(entry, nodeIPB)
	before := time.Now()
	if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusServiceUnavailable {
		t.Fatalf("post-recovery failure request = (%d, %v), want (503, nil)", status, err)
	}
	got := cooldownUntil(entry, nodeIPA).Sub(before)
	if got < gatewayLBDefaultBaseCooldown-time.Second || got > gatewayLBDefaultBaseCooldown+time.Second {
		t.Errorf("cooldown after recovery = %v, want approximately BaseCooldown (%v)", got, gatewayLBDefaultBaseCooldown)
	}
}

func TestGatewayLBIntegrationFallsBackToLogicalHostWhenAllCooling(t *testing.T) {
	host := rapidTestHost("allcooling")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	// 逻辑域名单独指向第三个假 Gateway：请求落到它即证明回退路径真的被走通了，
	// 而不是仅仅"没有改写 URL"。
	fallback := startFakeGateway(t, "FALLBACK")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	router.route(t, net.JoinHostPort(host, "80"), fallback)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	setCooldown(entry, nodeIPA, nodeIPB)
	router.resetDials()

	status, body, err := sendGatewayRequest(client, host, "/key")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if status != http.StatusOK || body != "FALLBACK" {
		t.Fatalf("request = (%d, %q), want (200, \"FALLBACK\")", status, body)
	}
	dialed := router.dialedAddrs()
	if len(dialed) != 1 || dialed[0] != net.JoinHostPort(host, "80") {
		t.Errorf("dialed = %v, want a single dial to the logical host", dialed)
	}
	if gwA.requestCount() != 0 || gwB.requestCount() != 0 {
		t.Errorf("cooling nodes got requests (A=%d B=%d), want none",
			gwA.requestCount(), gwB.requestCount())
	}
}

func TestGatewayLBIntegrationDoesNotMarkOn4xx(t *testing.T) {
	host := rapidTestHost("clienterr")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), startFakeGateway(t, "B"))
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// 403 是凭证事件、404 是业务应答，换节点都无意义，因此不得计入节点故障。
	for _, status := range []int{http.StatusForbidden, http.StatusNotFound} {
		gwA.setStatus(status)
		setCooldown(entry, nodeIPB)
		got, _, err := sendGatewayRequest(client, host, "/key")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if got != status {
			t.Fatalf("status = %d, want %d", got, status)
		}
		if keys := cooldownKeys(entry); containsString(keys, nodeIPA) {
			t.Errorf("cooldown = %v after HTTP %d, want node A untouched", keys, status)
		}
	}
}

func TestGatewayLBIntegrationCancelDoesNotMarkNodeFailed(t *testing.T) {
	host := rapidTestHost("cancelnomark")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA)

	gw := startFakeGateway(t, "A")
	gw.setDelay(500 * time.Millisecond) // 故意卡住，留出时间让取消在请求进行中真实发生
	router.route(t, nodeAddr(nodeIPA), gw)
	waitForSnapshot(t, entry, host, nodeIPA)

	// 此前"取消不算节点故障"这条性质只用 stubRoundTripper{err: context.Canceled} 这种合成
	// 错误在单元测试里验证过；这里改用真实网络栈+真实中途取消，确认标准库返回的错误依然满足
	// errors.Is(err, context.Canceled)，且节点确实没有被误标记。
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(30*time.Millisecond, cancel)

	if _, _, err := sendGatewayRequestCtx(client, ctx, host, "/key"); !errors.Is(err, context.Canceled) {
		t.Fatalf("request error = %v, want context.Canceled", err)
	}
	if keys := cooldownKeys(entry); containsString(keys, nodeIPA) {
		t.Errorf("cooldown = %v, want the canceled request to leave the node untouched", keys)
	}
}

func TestGatewayLBIntegrationClientTimeoutMarksNodeFailed(t *testing.T) {
	host := rapidTestHost("timeoutmark")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)
	client.Timeout = 50 * time.Millisecond // 客户端整体超时，比下面 A 节点的延迟短得多

	gwA := startFakeGateway(t, "A")
	gwA.setDelay(300 * time.Millisecond) // 故意卡住，确保必然触发客户端超时而不是偶然
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// 跟上一条用例对照：context.DeadlineExceeded 与 context.Canceled 不同，设计上有意保留为
	// 节点故障（超时很可能正是节点慢/半死的症状），此前也只在单元测试里验证过。
	setCooldown(entry, nodeIPB) // 强制首选卡住的节点 A

	if _, _, err := sendGatewayRequest(client, host, "/key"); err == nil {
		t.Fatal("request succeeded, want a client-timeout error")
	}
	if keys := cooldownKeys(entry); !containsString(keys, nodeIPA) {
		t.Fatalf("cooldown = %v, want node A marked after the client timeout", keys)
	}
}

func containsString(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// C 组：扩缩容与 DNS 生命周期
// ---------------------------------------------------------------------------

func TestGatewayLBIntegrationScaleOut(t *testing.T) {
	host := rapidTestHost("scaleout")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA)

	// 扩容前：全部流量只能落在 A。
	for i := 0; i < 20; i++ {
		if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
			t.Fatalf("pre-scale request %d = (%q, %v), want (\"A\", nil)", i, body, err)
		}
	}

	dns.setAnswer(host, nodeIPA, nodeIPB) // 扩容：域名新增一个 Gateway 节点
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	gwA.reset()
	gwB.reset()
	for i := 0; i < 200; i++ {
		if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusOK {
			t.Fatalf("post-scale request %d = (%d, %v), want (200, nil)", i, status, err)
		}
	}
	if gwB.requestCount() == 0 {
		t.Error("new node B never received traffic after scale-out")
	}
	if gwA.requestCount() == 0 {
		t.Error("existing node A stopped receiving traffic after scale-out")
	}
}

func TestGatewayLBIntegrationScaleIn(t *testing.T) {
	host := rapidTestHost("scalein")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	gwB := startFakeGateway(t, "B")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), gwB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// 先让两个节点都被使用过，从而在连接池里各自留下连接。
	for i := 0; i < 50; i++ {
		if _, _, err := sendGatewayRequest(client, host, "/key"); err != nil {
			t.Fatalf("pre-scale request %d failed: %v", i, err)
		}
	}
	if gwA.requestCount() == 0 || gwB.requestCount() == 0 {
		t.Fatalf("both nodes should have served traffic before scale-in (A=%d B=%d)",
			gwA.requestCount(), gwB.requestCount())
	}

	dns.setAnswer(host, nodeIPA) // 缩容：B 下线
	waitForSnapshot(t, entry, host, nodeIPA)

	gwB.reset()
	router.resetDials()
	for i := 0; i < 100; i++ {
		if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
			t.Fatalf("post-scale request %d = (%q, %v), want (\"A\", nil)", i, body, err)
		}
	}
	if gwB.requestCount() != 0 {
		t.Errorf("removed node B still served %d requests", gwB.requestCount())
	}
	// 摘除的节点不再产生新连接。按设计（8.1 节）已建立的空闲连接不主动关闭，靠
	// IdleConnTimeout 自然回收，因此这里只断言"没有新请求打过去"，不断言旧连接已消失。
	if containsString(router.dialedAddrs(), nodeAddr(nodeIPB)) {
		t.Errorf("dialed %v, want no new connection to the removed node", router.dialedAddrs())
	}
}

func TestGatewayLBIntegrationScaleInPrunesCooldown(t *testing.T) {
	host := rapidTestHost("scaleinprune")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	router.route(t, nodeAddr(nodeIPA), startFakeGateway(t, "A"))
	router.route(t, nodeAddr(nodeIPB), startFakeGateway(t, "B"))
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	entry.markFailed(net.ParseIP(nodeIPB))
	if keys := cooldownKeys(entry); !containsString(keys, nodeIPB) {
		t.Fatalf("cooldown = %v, want node B marked", keys)
	}

	dns.setAnswer(host, nodeIPA) // 缩容掉正处于冷却中的节点
	waitForSnapshot(t, entry, host, nodeIPA)
	waitForNoResolutionInFlight(t, entry)

	if keys := cooldownKeys(entry); containsString(keys, nodeIPB) {
		t.Errorf("cooldown = %v, want the removed node's record pruned", keys)
	}
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
		t.Fatalf("request = (%q, %v), want (\"A\", nil)", body, err)
	}
}

func TestGatewayLBIntegrationRefreshFailureKeepsSnapshot(t *testing.T) {
	host := rapidTestHost("keepsnapshot")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	router.route(t, nodeAddr(nodeIPA), startFakeGateway(t, "A"))
	router.route(t, nodeAddr(nodeIPB), startFakeGateway(t, "B"))
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)

	// DNS 转为故障。先确认刷新确实被尝试过（查询数增长），再断言快照未被空结果覆盖，
	// 否则"快照没变"可能只是因为压根没刷新，断言就失去意义。
	dns.resetQueries()
	dns.setFailure()
	forceSnapshotExpired(entry)
	deadline := time.Now().Add(3 * time.Second)
	for aQueryCount(dns, host) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("no refresh attempt observed after injecting the DNS failure")
		}
		entry.resolve(canonicalHost(host))
		time.Sleep(2 * time.Millisecond)
	}
	waitForNoResolutionInFlight(t, entry)

	if got := ipStrings(entry.resolve(canonicalHost(host))); !sameStringSet(got, []string{nodeIPA, nodeIPB}) {
		t.Fatalf("snapshot = %v, want the last successful snapshot to be preserved", got)
	}
	// 快照仍在，请求继续走负载均衡而不是回退。
	router.resetDials()
	if status, _, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusOK {
		t.Fatalf("request = (%d, %v), want (200, nil)", status, err)
	}
	dialed := router.dialedAddrs()
	if len(dialed) == 0 || (dialed[0] != nodeAddr(nodeIPA) && dialed[0] != nodeAddr(nodeIPB)) {
		t.Errorf("dialed = %v, want traffic to keep using the preserved snapshot", dialed)
	}
}

func TestGatewayLBIntegrationColdFailureIsRateLimitedThenSelfHeals(t *testing.T) {
	host := rapidTestHost("coldheal")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA)

	gwA := startFakeGateway(t, "A")
	fallback := startFakeGateway(t, "FALLBACK")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, net.JoinHostPort(host, "80"), fallback)

	// The first request must trigger one failed lookup and fall back to the logical host.
	dns.setFailure()
	const rounds = 60
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "FALLBACK" {
		t.Fatalf("first request = (%q, %v), want (\"FALLBACK\", nil)", body, err)
	}
	waitForNoResolutionInFlight(t, entry)
	queriesAfterFailure := aQueryCount(dns, host)
	if queriesAfterFailure == 0 {
		t.Fatal("the first failed lookup produced no DNS query")
	}
	entry.mu.Lock()
	nextRefresh := entry.expireAt
	entry.mu.Unlock()
	if !nextRefresh.After(time.Now()) {
		t.Fatalf("next refresh = %v, want a future deadline after DNS failure", nextRefresh)
	}

	// Requests inside the failed-refresh window must not start another lookup.
	for i := 1; i < rounds; i++ {
		if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "FALLBACK" {
			t.Fatalf("request %d = (%q, %v), want (\"FALLBACK\", nil)", i, body, err)
		}
	}
	waitForNoResolutionInFlight(t, entry)
	if got := aQueryCount(dns, host); got != queriesAfterFailure {
		t.Fatalf("DNS queries after %d rate-limited requests = %d, want %d", rounds, got, queriesAfterFailure)
	}

	// Expiring the snapshot simulates the next scheduled refresh after DNS recovers.
	dns.setAnswer(host, nodeIPA)
	waitForSnapshot(t, entry, host, nodeIPA)
	router.resetDials()
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
		t.Fatalf("request after recovery = (%q, %v), want (\"A\", nil)", body, err)
	}
	if dialed := router.dialedAddrs(); len(dialed) == 0 || dialed[0] != nodeAddr(nodeIPA) {
		t.Errorf("dialed = %v, want traffic to move onto the node once DNS recovered", dialed)
	}
}

func TestGatewayLBIntegrationDNSTimeoutFallsBackThenRecovers(t *testing.T) {
	host := rapidTestHost("dnstimeout")

	// 不用 setupGatewayLB：需要一开始就让 DNS 静默丢包，而不是先给一个正常答案。
	dns := startFakeDNS(t)
	dns.setSilent() // 模拟 DNS 服务器丢包不回应——与 setFailure 的"立即 SERVFAIL"是两种不同的
	// 故障形态，此前的用例只测过后者。
	installFakeDNS(t, dns)
	resetGatewayLBEntry(host)
	entry := entryFor(canonicalHost(host))
	t.Cleanup(func() {
		waitForNoResolutionInFlight(t, entry)
		resetGatewayLBEntry(host)
	})

	router := newGatewayRouter()
	client := newGatewayLBClient(router)
	fallback := startFakeGateway(t, "FALLBACK")
	router.route(t, net.JoinHostPort(host, "80"), fallback)

	// 冷启动不得阻塞在 DNS 上：结构性断言是"请求返回时后台解析仍在途"——若请求真的等了
	// 解析，返回时 resolving 必然已经为 false。不测经过时长：墙钟断言在慢机器上会误报。
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "FALLBACK" {
		t.Fatalf("request = (%q, %v), want (\"FALLBACK\", nil)", body, err)
	}
	entry.mu.Lock()
	resolving := entry.resolving
	entry.mu.Unlock()
	if !resolving {
		t.Fatal("background DNS lookup was not in flight when the fallback request returned")
	}

	// 确认后台确实发出过查询（而不是压根没来得及触发解析），否则上面的回退可能只是巧合，
	// 断言会失去意义。
	deadline := time.Now().Add(gatewayLBDefaultDNSTimeout + time.Second)
	for !anyAQuery(dns) {
		if time.Now().After(deadline) {
			t.Fatal("no DNS query was observed while the DNS server stayed silent")
		}
		time.Sleep(5 * time.Millisecond)
	}
	// 解析最终结束（resolving 清零）即证明超时被接线生效、没有永久挂起；精确的时长测量是
	// 墙钟断言，在慢机器上会误报，故不做。
	waitForNoResolutionInFlight(t, entry)

	// After DNS recovers, expiring the snapshot triggers a real lookup and recovery.
	gw := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gw)
	dns.setAnswer(host, nodeIPA)
	waitForSnapshot(t, entry, host, nodeIPA)
	router.resetDials()
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
		t.Fatalf("request after recovery = (%q, %v), want (\"A\", nil)", body, err)
	}
}

func TestGatewayLBIntegrationDeduplicatesResolvedNodes(t *testing.T) {
	host := rapidTestHost("dedup")
	_, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPA, nodeIPA)

	gwA := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gwA)

	// DNS 返回三条相同的 A 记录，快照必须只保留一个候选，否则该节点会在 pick 里占多个位置。
	waitForSnapshot(t, entry, host, nodeIPA)
	if nodes := entry.resolve(canonicalHost(host)); len(nodes) != 1 {
		t.Fatalf("snapshot = %v, want a single node after dedup", ipStrings(nodes))
	}
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
		t.Fatalf("request = (%q, %v), want (\"A\", nil)", body, err)
	}
}

// ---------------------------------------------------------------------------
// D 组：并发与资源
// ---------------------------------------------------------------------------

func TestGatewayLBIntegrationConcurrentColdStartResolvesOnce(t *testing.T) {
	host := rapidTestHost("singleflight")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA)

	router.route(t, nodeAddr(nodeIPA), startFakeGateway(t, "A"))
	router.route(t, net.JoinHostPort(host, "80"), startFakeGateway(t, "FALLBACK"))

	// Concurrent cold-start requests must share the single in-flight lookup.
	const goroutines = 100
	var wg sync.WaitGroup
	errCh := make(chan error, goroutines)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			status, body, err := sendGatewayRequestCtx(client, ctx, host, "/key")
			if err != nil || status != http.StatusOK || (body != "A" && body != "FALLBACK") {
				errCh <- fmt.Errorf("request = (%d, %q, %v)", status, body, err)
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Errorf("concurrent cold-start %v, want a successful gateway or fallback response", err)
	}
	waitForNoResolutionInFlight(t, entry)

	if got := aQueryCount(dns, host); got != 1 {
		t.Errorf("observed %d resolutions for %d concurrent cold-start requests, want exactly 1",
			got, goroutines)
	}
}

func TestGatewayLBIntegrationDoesNotLeakGoroutines(t *testing.T) {
	host := rapidTestHost("leak")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	gwA := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gwA)
	router.route(t, nodeAddr(nodeIPB), startFakeGateway(t, "B"))
	router.route(t, net.JoinHostPort(host, "80"), startFakeGateway(t, "FALLBACK"))
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	// Exercise initial resolution, normal traffic, node failure, DNS failure and recovery.
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)
	for i := 0; i < 50; i++ {
		sendGatewayRequest(client, host, "/key")
	}
	gwA.setStatus(http.StatusBadGateway)
	for i := 0; i < 20; i++ {
		sendGatewayRequest(client, host, "/key")
	}
	gwA.setStatus(http.StatusOK)
	dns.resetQueries()
	dns.setFailure()
	forceSnapshotExpired(entry)
	for i := 0; i < 20; i++ {
		sendGatewayRequest(client, host, "/key")
	}
	waitForNoResolutionInFlight(t, entry)
	if !anyAQuery(dns) {
		t.Fatal("goroutine leak scenario did not execute the DNS failure refresh")
	}

	// Exercise the context-timeout path as well as immediate SERVFAIL.
	dns.resetQueries()
	dns.setSilent()
	forceSnapshotExpired(entry)
	entry.resolve(canonicalHost(host))
	waitForNoResolutionInFlight(t, entry)
	if !anyAQuery(dns) {
		t.Fatal("goroutine leak scenario did not execute the silent DNS timeout")
	}

	dns.setAnswer(host, nodeIPA)
	waitForSnapshot(t, entry, host, nodeIPA)
	dns.setAnswer(host, nodeIPA, nodeIPB)
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)
	waitForNoResolutionInFlight(t, entry)
	router.closeIdleConnections()

	// All lookup and connection goroutines created after the baseline must settle.
	deadline := time.Now().Add(5 * time.Second)
	for {
		current := runtime.NumGoroutine()
		if current <= baseline {
			return
		}
		if time.Now().After(deadline) {
			t.Errorf("goroutines = %d, want it to settle back to the pre-lookup baseline %d", current, baseline)
			return
		}
		runtime.Gosched()
		time.Sleep(20 * time.Millisecond)
	}
}

func TestGatewayLBIntegrationConcurrentRequestsDuringRefresh(t *testing.T) {
	host := rapidTestHost("racyrefresh")
	dns, router, client, entry := setupGatewayLB(t, host, nodeIPA, nodeIPB)

	router.route(t, nodeAddr(nodeIPA), startFakeGateway(t, "A"))
	router.route(t, nodeAddr(nodeIPB), startFakeGateway(t, "B"))
	gwC := startFakeGateway(t, "C")
	router.route(t, nodeAddr(nodeIPC), gwC)
	router.route(t, net.JoinHostPort(host, "80"), startFakeGateway(t, "FALLBACK"))
	waitForSnapshot(t, entry, host, nodeIPA, nodeIPB)
	dns.resetQueries()

	// 请求、扩缩容刷新、失败标记三者并发进行，主要目的是在 -race 下暴露竞态。
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var totalCount, successCount int64

	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					atomic.AddInt64(&totalCount, 1)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					status, _, err := sendGatewayRequestCtx(client, ctx, host, "/key")
					cancel()
					if err == nil && status == http.StatusOK {
						atomic.AddInt64(&successCount, 1)
					}
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		answers := [][]string{
			{nodeIPA, nodeIPB},
			{nodeIPA},
			{nodeIPA, nodeIPB, nodeIPC},
			{nodeIPB, nodeIPC},
		}
		for i := 0; ; i++ {
			select {
			case <-stop:
				return
			default:
				dns.setAnswer(host, answers[i%len(answers)]...)
				forceSnapshotExpired(entry)
				time.Sleep(3 * time.Millisecond)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				entry.markFailed(net.ParseIP(nodeIPA))
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	time.Sleep(300 * time.Millisecond)
	close(stop)
	wg.Wait()
	waitForNoResolutionInFlight(t, entry)

	total, success := atomic.LoadInt64(&totalCount), atomic.LoadInt64(&successCount)
	if total == 0 {
		t.Fatal("no requests were attempted")
	}
	if success != total {
		t.Errorf("successful requests during concurrent scale/failure churn = %d/%d, want all requests to succeed", success, total)
	}
	if queries := aQueryCount(dns, host); queries == 0 {
		t.Fatal("no DNS refresh occurred during concurrent churn")
	}

	// Publish a deterministic final snapshot and prove requests use the new node.
	dns.setAnswer(host, nodeIPC)
	waitForSnapshot(t, entry, host, nodeIPC)
	gwC.reset()
	if status, body, err := sendGatewayRequest(client, host, "/key"); err != nil || status != http.StatusOK || body != "C" {
		t.Fatalf("request after concurrent refresh = (%d, %q, %v), want (200, \"C\", nil)", status, body, err)
	}
	if gwC.requestCount() != 1 {
		t.Fatalf("gateway C requests after final refresh = %d, want 1", gwC.requestCount())
	}
}

// ---------------------------------------------------------------------------
// E 组：门禁在真实链路上的表现
// ---------------------------------------------------------------------------

// TestGatewayLBIntegrationCustomTransportWithEnvProxyIsSkipped verifies the
// conservative gate for an opaque RoundTripper: any HTTP/ALL proxy environment
// variable disables LB, even when NO_PROXY contains the logical host.
func TestGatewayLBIntegrationCustomTransportWithEnvProxyIsSkipped(t *testing.T) {
	for _, name := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
		t.Setenv(name, "")
	}
	t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
	t.Setenv("NO_PROXY", ".myqcloud.com")
	host := rapidTestHost("envproxy")
	_, router, _, entry := setupGatewayLB(t, host, nodeIPA)

	gw := startFakeGateway(t, "A")
	router.route(t, nodeAddr(nodeIPA), gw)
	router.route(t, net.JoinHostPort(host, "80"), gw)
	waitForSnapshot(t, entry, host, nodeIPA)

	// 包一层自定义 RoundTripper：它内部仍转交给带路由的 transport，所以一旦 LB 误判为"无代理"
	// 而介入，dial 目标就会变成节点 IP，从而被下面的断言捕获。
	inner := router.transport()
	custom := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return inner.RoundTrip(req)
	})
	client := &http.Client{Transport: newRapidGatewayLBTransport(custom)}

	router.resetDials()
	if _, body, err := sendGatewayRequest(client, host, "/key"); err != nil || body != "A" {
		t.Fatalf("request = (%q, %v), want (\"A\", nil)", body, err)
	}
	dialed := router.dialedAddrs()
	if len(dialed) != 1 || dialed[0] != net.JoinHostPort(host, "80") {
		t.Errorf("dialed = %v, want the logical host (LB must stay out of the way)", dialed)
	}
}

func TestGatewayLBIntegrationHTTPSIsNotRewritten(t *testing.T) {
	host := rapidTestHost("httpsgate")
	_, router, _, entry := setupGatewayLB(t, host, nodeIPA)

	router.route(t, nodeAddr(nodeIPA), startFakeGateway(t, "A"))
	waitForSnapshot(t, entry, host, nodeIPA)

	client := newGatewayLBClient(router)
	router.resetDials()

	// 只断言连接目标未被改写。这里不追求 TLS 握手成功：httptest 的自签证书是签给
	// example.com/127.0.0.1 的，而 SNI 会用逻辑域名，握手必然失败——而握手之前的 dial
	// 目标已经足以证明门禁生效。
	req, err := http.NewRequest(http.MethodGet, "https://"+host+"/key", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}
	if resp, err := client.Do(req); err == nil {
		resp.Body.Close()
	}

	dialed := router.dialedAddrs()
	if len(dialed) != 1 || dialed[0] != net.JoinHostPort(host, "443") {
		t.Errorf("dialed = %v, want a single dial to the logical host on 443", dialed)
	}
}

// roundTripperFunc 让普通函数满足 http.RoundTripper，用于构造"非 *http.Transport 的下一跳"。
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
