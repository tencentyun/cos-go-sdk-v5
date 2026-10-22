package cos

import (
	"context"
	"errors"
	math_rand "math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// gatewayLBDefaultRefresh 是 DNS 快照的刷新间隔。
	gatewayLBDefaultRefresh = 5 * time.Second
	// gatewayLBDefaultBaseCooldown 是节点首次失败后的冷却时长：足够短，让偶发抖动的代价
	// 很小；连续失败时按 2^(streak-1) 翻倍，直到封顶 gatewayLBDefaultMaxCooldown。
	gatewayLBDefaultBaseCooldown = 5 * time.Second
	// gatewayLBDefaultMaxCooldown 是冷却时长的封顶值，替代原先单一固定的
	// gatewayLBDefaultCooldown（数值不变，仍是 60s）。
	gatewayLBDefaultMaxCooldown = 60 * time.Second
	// gatewayLBDefaultFailureResetWindow 是判定"上一次失败与这一次是否算连续"的最大间隔：
	// 超过这个间隔的失败视为全新序列，streak 从 1 重新起算。必须显著大于
	// gatewayLBDefaultMaxCooldown（这里留了 5 倍余量），否则一个"每次冷却刚到期就失败"的
	// 节点会被误判为不连续，永远升级不到封顶值，退避机制形同虚设。
	gatewayLBDefaultFailureResetWindow = 5 * time.Minute
	// gatewayLBDefaultDNSTimeout bounds each background DNS lookup.
	gatewayLBDefaultDNSTimeout = 2 * time.Second
)

// gatewayLBRegistry 按域名（经 canonicalHost 规整）缓存 *gatewayLBEntry，跨请求共享。
// 不做空闲淘汰：单进程访问的高性能桶域名数量有限。
var gatewayLBRegistry sync.Map

// gatewayLBEntry 是某个高性能桶短域名的 DNS 快照与节点健康状态。
type gatewayLBEntry struct {
	mu sync.Mutex
	// nodes 是当前 DNS 快照，整体替换、不做原地修改；至今未成功解析过时为 nil。
	nodes []net.IP
	// expireAt 是下次允许解析的时刻。零值落在过去，故首次访问必然触发一次解析。
	expireAt time.Time
	// resolving 标记是否已有一次解析在途，避免并发重复起 goroutine。
	resolving bool
	// cooldown 记录 ip.String() -> 冷却截止时间。
	cooldown map[string]time.Time
	// failStreak 记录 ip.String() -> 当前连续失败次数，用于计算指数退避的冷却时长。
	failStreak map[string]int
	// lastFailAt 记录 ip.String() -> 最近一次失败时间，用于判断这次失败与上一次是否
	// 仍在 FailureResetWindow 内、算作"连续"。
	lastFailAt map[string]time.Time
}

// entryFor 返回（或创建）某个域名对应的 gatewayLBEntry，host 须已经过 canonicalHost 规整。
func entryFor(host string) *gatewayLBEntry {
	if v, ok := gatewayLBRegistry.Load(host); ok {
		return v.(*gatewayLBEntry)
	}
	// LoadOrStore 保证并发下只有一个 entry 生效，多余的构造结果被丢弃。
	v, _ := gatewayLBRegistry.LoadOrStore(host, &gatewayLBEntry{})
	return v.(*gatewayLBEntry)
}

// resolve 返回当前 DNS 快照，惰性刷新、无常驻 goroutine，全程不阻塞调用方。返回的快照可能为
// nil（首次解析尚未完成，或持续解析失败），此时调用方回退逻辑域名——与门禁不满足时同一条路径。
//
// 到了刷新时刻且无解析在途时踢一次后台解析，结果由后续请求享用。不等待解析完成，是为了让共享
// 状态与任何单次请求的生命周期解耦：否则一个恰好被取消的请求会连带影响其他并发请求的选点。
// 代价是冷启动窗口（一次 DNS 的时间）内不做负载均衡。
func (e *gatewayLBEntry) resolve(host string) []net.IP {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.resolving && time.Now().After(e.expireAt) {
		e.resolving = true
		go e.refreshAsync(host)
	}
	return e.nodes
}

// refreshAsync 在后台解析一次并发布结果，是唯一的解析入口（首次解析与周期刷新共用）。用独立
// context：解析不属于任何一次请求，不应被恰好触发它的那次请求取消。
func (e *gatewayLBEntry) refreshAsync(host string) {
	ips := lookupGatewayNodes(context.Background(), host, gatewayLBDefaultDNSTimeout)
	e.mu.Lock()
	if len(ips) > 0 { // 解析失败或空结果时保留旧快照，避免一次抖动丢掉可用节点
		e.nodes = ips
		e.pruneCooldown(ips)
	}
	// 无论成败都推进：解析失败时同样限速，避免每个请求都触发一次解析。
	e.expireAt = time.Now().Add(gatewayLBDefaultRefresh)
	e.resolving = false
	e.mu.Unlock()
}

// lookupGatewayNodes 解析域名的 A/AAAA 记录，合并调用方 ctx 与固定超时，失败统一返回 nil。
func lookupGatewayNodes(ctx context.Context, host string, timeout time.Duration) []net.IP {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil
	}
	return normalizeGatewayNodes(addrs)
}

// normalizeGatewayNodes 把解析结果整理成快照：过滤空地址、按 ip.String() 去重，保持解析器
// 给出的顺序。去重是必要的——重复 IP 会在 pick 的候选里占多个位置，破坏均匀随机。
func normalizeGatewayNodes(addrs []net.IPAddr) []net.IP {
	ips := make([]net.IP, 0, len(addrs))
	seen := make(map[string]struct{}, len(addrs))
	for _, addr := range addrs {
		if len(addr.IP) == 0 {
			continue
		}
		// 用 ip.String() 而非字节做 key：同一 IPv4 有 4 字节和 16 字节两种表示，String()
		// 是两者的规范形式，也与 cooldown 的 key 保持一致。
		key := addr.IP.String()
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		ips = append(ips, addr.IP)
	}
	return ips
}

// pruneCooldown 删除已不在最新快照中的冷却/失败记录，避免节点 IP 频繁更换时这几张 map
// 无限增长。调用方须持有 e.mu。
func (e *gatewayLBEntry) pruneCooldown(nodes []net.IP) {
	if len(e.cooldown) == 0 && len(e.failStreak) == 0 && len(e.lastFailAt) == 0 {
		return
	}
	keep := make(map[string]struct{}, len(nodes))
	for _, ip := range nodes {
		keep[ip.String()] = struct{}{}
	}
	for k := range e.cooldown {
		if _, ok := keep[k]; !ok {
			delete(e.cooldown, k)
		}
	}
	for k := range e.failStreak {
		if _, ok := keep[k]; !ok {
			delete(e.failStreak, k)
		}
	}
	for k := range e.lastFailAt {
		if _, ok := keep[k]; !ok {
			delete(e.lastFailAt, k)
		}
	}
}

// pick 从快照中均匀随机选一个不在冷却期的节点。全部节点都在冷却时返回 ok=false，由调用方
// 回退到未改写的逻辑域名，而不是从已知失败的节点里硬选一个。
// 随机源沿用本包 init() 中已 seed 的全局 math_rand。
func (e *gatewayLBEntry) pick(nodes []net.IP) (net.IP, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	candidates := make([]net.IP, 0, len(nodes))
	for _, ip := range nodes {
		if until, cooling := e.cooldown[ip.String()]; !cooling || now.After(until) {
			candidates = append(candidates, ip)
		}
	}
	if len(candidates) == 0 {
		return nil, false
	}
	return candidates[math_rand.Intn(len(candidates))], true
}

// markFailed 记一次节点失败：首次失败冷却 gatewayLBDefaultBaseCooldown；若与上一次失败的
// 间隔未超过 gatewayLBDefaultFailureResetWindow，视为同一连续失败序列，冷却时长按
// base*2^(streak-1) 指数翻倍，封顶 gatewayLBDefaultMaxCooldown；否则视为新序列，从
// base 重新起算。冷却是选点提示而非永久拉黑，到期自动恢复。
func (e *gatewayLBEntry) markFailed(ip net.IP) {
	key := ip.String()
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	if e.failStreak == nil {
		e.failStreak = make(map[string]int)
	}
	if e.lastFailAt == nil {
		e.lastFailAt = make(map[string]time.Time)
	}
	if last, ok := e.lastFailAt[key]; !ok || now.Sub(last) > gatewayLBDefaultFailureResetWindow {
		e.failStreak[key] = 0
	}
	e.failStreak[key]++
	e.lastFailAt[key] = now

	cd := gatewayLBDefaultBaseCooldown
	for i := 1; i < e.failStreak[key] && cd < gatewayLBDefaultMaxCooldown; i++ {
		cd *= 2
		if cd <= 0 { // 溢出保护，正常参数下不会触发
			cd = gatewayLBDefaultMaxCooldown
			break
		}
	}
	if cd > gatewayLBDefaultMaxCooldown {
		cd = gatewayLBDefaultMaxCooldown
	}
	if e.cooldown == nil {
		e.cooldown = make(map[string]time.Time)
	}
	e.cooldown[key] = now.Add(cd)
}

// markSucceeded 清空该节点的连续失败记录：代表它这次已经证实健康，下次再失败会从
// gatewayLBDefaultBaseCooldown 重新起算，而不是延续之前的退避等级。不会提前解除本次已经
// 生效的冷却——正常路径下 pick() 不会选中仍在冷却期的节点，故不会对正在冷却的节点调用它。
func (e *gatewayLBEntry) markSucceeded(ip net.IP) {
	key := ip.String()
	e.mu.Lock()
	delete(e.failStreak, key)
	delete(e.lastFailAt, key)
	e.mu.Unlock()
}

// useGatewayLB 判断该请求的目标桶是否应经内部 Gateway 负载均衡 Transport 转发，供
// AuthorizationTransport/CVMCredentialTransport/CredentialTransport/StsCredentialTransport
// 的 transport(req) 在构造 rapidGatewayLBTransport 前调用。判断口径与 useSessionAuth
// 一致：优先取 Host 头（签名/逻辑 Host），为空才回退 req.URL.Hostname()，避免 BaseURL 与
// Host 头不一致（如自定义域名 + rapid bucket Host 头）时误判。
func useGatewayLB(req *http.Request) bool {
	host := req.Host
	if host == "" && req.URL != nil {
		host = req.URL.Hostname()
	}
	bucketName := parseBucketNameFromHost(host)
	return bucketName != "" && isRapidBucket(bucketName)
}

// rapidGatewayLBTransport 是高性能桶 Gateway 的内部负载均衡 Transport：每个请求都重新选一个
// 物理节点，只改写连接目标，保持 HTTP Host 与签名 Host 为逻辑短域名。
//
// 不做内部重试：失败时只标记节点冷却并把结果原样上抛，重试由上层决定。因为选点是请求级、
// 冷却状态跨请求共享，上层的任意一次重试都会重新选点并自然避开刚失败的节点。
type rapidGatewayLBTransport struct {
	next http.RoundTripper
}

// newRapidGatewayLBTransport 构造内部负载均衡 Transport。next 为 nil 时回落 http.DefaultTransport。
func newRapidGatewayLBTransport(next http.RoundTripper) *rapidGatewayLBTransport {
	return &rapidGatewayLBTransport{next: next}
}

// transport 返回下一跳，未设置时回落 http.DefaultTransport，与 AuthorizationTransport 一致。
func (t *rapidGatewayLBTransport) transport() http.RoundTripper {
	if t.next != nil {
		return t.next
	}
	return http.DefaultTransport
}

// RoundTrip 实现 http.RoundTripper：非高性能桶短域名、不满足传输门禁、尚无可用快照、全部节点
// 冷却中这四种情况都原样透传，不改写、不报错；其余情况改写连接目标，并在失败时标记节点冷却、
// 成功时清空该节点的连续失败记录。
func (t *rapidGatewayLBTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		// 与其在本层解引用 nil 而 panic，不如原样透传，由标准库按既有方式报错。
		return t.transport().RoundTrip(req)
	}
	// rawHost 保留调用方原始写法（可能含端口、大小写），用于回写 Host 头；lookupHost 是规整
	// 后的纯域名，用于 DNS 查询与缓存 key。必须分开：解析接口不接受 "host:port"。
	rawHost := req.Host
	if rawHost == "" {
		rawHost = req.URL.Host // 与标准库写 Host 头时的回退方式一致
	}
	lookupHost := canonicalHost(rawHost)
	if !isRapidBucket(parseBucketNameFromHost(lookupHost)) {
		return t.transport().RoundTrip(req)
	}
	if !t.eligible(req) {
		return t.transport().RoundTrip(req)
	}

	entry := entryFor(lookupHost)
	nodes := entry.resolve(lookupHost)
	if len(nodes) == 0 {
		return t.transport().RoundTrip(req)
	}
	node, ok := entry.pick(nodes)
	if !ok {
		return t.transport().RoundTrip(req)
	}

	// cloneRequest 只深拷贝 Header，URL 仍与原请求共享指针，故需再深拷贝一次 URL，
	// 避免改写作用到调用方的请求对象上。
	req = cloneRequest(req) // per RoundTrip contract
	// 显式钉死 Host 头：Host 为空时标准库会取 req.URL.Host，而它下一步就要被改成物理 IP。
	// 必须用 rawHost，写回规整后的变体会与上一层签名所用的值不一致。
	req.Host = rawHost
	// eligible 已保证 scheme 为 http 且无显式端口，故端口恒为 80。
	urlCopy := *req.URL
	urlCopy.Host = net.JoinHostPort(node.String(), "80")
	req.URL = &urlCopy

	resp, err := t.transport().RoundTrip(req)
	if shouldMarkGatewayNodeFailed(req, resp, err) {
		entry.markFailed(node)
	} else if err == nil && resp != nil {
		entry.markSucceeded(node)
	}
	return resp, err
}

// eligible 判定本次请求是否满足传输门禁，四条与 Python SDK 的传输类门禁一一对应。
func (t *rapidGatewayLBTransport) eligible(req *http.Request) bool {
	if req.URL.Scheme != "http" {
		// 改写为裸 IP 会同时改变 TLS SNI 与证书校验用的 ServerName，导致 HTTPS 必然失败。
		return false
	}
	if req.URL.Port() != "" {
		// 生产 Gateway 短域名不带显式端口，出现端口即非标准场景，不介入。
		return false
	}
	if net.ParseIP(req.URL.Hostname()) != nil {
		// Host 已是 IP：调用方可能刻意直连某个实例，不应用自己解析的 IP 集合覆盖其选择。
		return false
	}
	if t.hasProxy() {
		// Rewriting the target conflicts with proxy routing.
		return false
	}
	return true
}

// hasProxy conservatively disables LB whenever proxy routing may be configured.
// It intentionally ignores NO_PROXY because rewriting a logical host to a physical
// IP can change the bypass decision made by the underlying transport.
func (t *rapidGatewayLBTransport) hasProxy() bool {
	if t.next == nil {
		return hasHTTPProxyEnvironment()
	}
	if ht, ok := t.next.(*http.Transport); ok {
		return ht.Proxy != nil
	}
	return hasHTTPProxyEnvironment()
}

func hasHTTPProxyEnvironment() bool {
	for _, name := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
		if os.Getenv(name) != "" {
			return true
		}
	}
	return false
}

// canonicalHost 把 Host 字符串规整为纯域名：剥掉端口、去掉末尾根点、转为小写。
func canonicalHost(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	return strings.ToLower(strings.TrimSuffix(host, "."))
}

// shouldMarkGatewayNodeFailed 判定本次失败是否应记作该物理节点的问题。
//
// 排除调用方主动取消：取消与节点健康无关，误记账会让健康节点被冷却并连累并发请求的选点。
// DeadlineExceeded 有意保留为节点故障：超时虽有歧义，但很可能正是节点慢的症状。
func shouldMarkGatewayNodeFailed(req *http.Request, resp *http.Response, err error) bool {
	if err == nil {
		return isAmbiguousGatewayStatus(resp)
	}
	// 取消时 Transport 返回的就是 context.Cause(ctx)、未做包装，故 errors.Is 可靠；再查一次
	// ctx.Err() 以覆盖 WithCancelCause 传入自定义 cause 的情况。
	if errors.Is(err, context.Canceled) || req.Context().Err() == context.Canceled {
		return false
	}
	return true
}

// isAmbiguousGatewayStatus 判定响应是否属于节点级故障信号。现役 Gateway 以 500 表示部分后端
// 传输失败，故与 502/503/504 同类；普通 4xx 是业务应答、403 是凭证事件，换节点无意义。
func isAmbiguousGatewayStatus(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

var _ http.RoundTripper = (*rapidGatewayLBTransport)(nil)
