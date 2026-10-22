package cos

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

// rapidGatewayTestHost 是测试用的高性能桶短域名。
const rapidGatewayTestHost = "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"

// stubRoundTripper 记录最后一次收到的请求并返回预置结果，用于观察请求被改写成什么样。
type stubRoundTripper struct {
	resp *http.Response
	err  error

	mu       sync.Mutex
	calls    int
	lastURL  string
	lastHost string
}

func (s *stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s.mu.Lock()
	s.calls++
	s.lastURL = "" // URL 可能为 nil：覆盖透传无 URL 请求的用例
	if req.URL != nil {
		s.lastURL = req.URL.Host
	}
	s.lastHost = req.Host
	s.mu.Unlock()
	if s.err != nil {
		return nil, s.err
	}
	resp := s.resp
	if resp == nil {
		resp = &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}
	}
	return resp, nil
}

func (s *stubRoundTripper) snapshot() (calls int, urlHost, host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls, s.lastURL, s.lastHost
}

// newGatewayLBTestRequest 构造请求，host 为空表示不设 req.Host（覆盖 Host 头回退的场景）。
func newGatewayLBTestRequest(t *testing.T, rawURL, host string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q) failed: %v", rawURL, err)
	}
	// http.NewRequest 会把 URL.Host（含端口）填入 req.Host，这里按用例需要覆盖或清空。
	req.Host = host
	return req
}

// resetGatewayLBEntry 清理注册表项，避免用例间互相影响（注册表是包级单例）。
func resetGatewayLBEntry(host string) {
	gatewayLBRegistry.Delete(canonicalHost(host))
}

// seedGatewayLBEntry 预置一份快照以跳过真实 DNS：expireAt 设为未来时刻，避免用例运行期间
// 触发后台解析。
func seedGatewayLBEntry(host string, ips ...net.IP) *gatewayLBEntry {
	e := entryFor(canonicalHost(host))
	e.mu.Lock()
	e.nodes = ips
	e.expireAt = time.Now().Add(time.Hour)
	e.resolving = false
	e.cooldown = nil
	e.mu.Unlock()
	return e
}

func TestCanonicalHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{"已是纯域名", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"},
		{"剥掉端口", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com:8080", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"},
		{"转为小写", "Bucket-X--1259654469.Cosrapid.AP-Guangzhou.Myqcloud.com", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"},
		{"去掉末尾根点", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com.", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"},
		{"端口与大小写同时处理", "Bucket-X--1259654469.Cosrapid.ap-guangzhou.myqcloud.com:80", "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"},
		{"IPv6 带端口", "[2001:db8::1]:80", "2001:db8::1"},
		{"空字符串", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := canonicalHost(tc.host); got != tc.want {
				t.Errorf("canonicalHost(%q) = %q, want %q", tc.host, got, tc.want)
			}
		})
	}
}

func TestIsAmbiguousGatewayStatus(t *testing.T) {
	tests := []struct {
		name string
		resp *http.Response
		want bool
	}{
		{"响应为空", nil, false},
		{"200 成功", &http.Response{StatusCode: http.StatusOK}, false},
		{"403 凭证事件不算节点故障", &http.Response{StatusCode: http.StatusForbidden}, false},
		{"404 业务应答不算节点故障", &http.Response{StatusCode: http.StatusNotFound}, false},
		{"500", &http.Response{StatusCode: http.StatusInternalServerError}, true},
		{"501 不在列", &http.Response{StatusCode: http.StatusNotImplemented}, false},
		{"502", &http.Response{StatusCode: http.StatusBadGateway}, true},
		{"503", &http.Response{StatusCode: http.StatusServiceUnavailable}, true},
		{"504", &http.Response{StatusCode: http.StatusGatewayTimeout}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isAmbiguousGatewayStatus(tc.resp); got != tc.want {
				t.Errorf("isAmbiguousGatewayStatus() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestShouldMarkGatewayNodeFailed(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		resp *http.Response
		err  error
		want bool
	}{
		{
			name: "成功响应不标记",
			ctx:  context.Background(),
			resp: &http.Response{StatusCode: http.StatusOK},
			want: false,
		},
		{
			name: "5xx 标记",
			ctx:  context.Background(),
			resp: &http.Response{StatusCode: http.StatusBadGateway},
			want: true,
		},
		{
			name: "普通网络错误标记",
			ctx:  context.Background(),
			err:  errors.New("dial tcp: connection refused"),
			want: true,
		},
		{
			name: "调用方取消不标记",
			ctx:  context.Background(),
			err:  context.Canceled,
			want: false,
		},
		{
			name: "包装后的取消错误不标记",
			ctx:  context.Background(),
			err:  fmt.Errorf("get %q: %w", "object", context.Canceled),
			want: false,
		},
		{
			name: "自定义 cause 但 ctx 已取消不标记",
			ctx:  canceledCtx,
			err:  errors.New("custom cancel cause"),
			want: false,
		},
		{
			name: "超时仍视为节点故障",
			ctx:  context.Background(),
			err:  context.DeadlineExceeded,
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := newGatewayLBTestRequest(t, "http://"+rapidGatewayTestHost+"/key", rapidGatewayTestHost)
			req = req.WithContext(tc.ctx)
			if got := shouldMarkGatewayNodeFailed(req, tc.resp, tc.err); got != tc.want {
				t.Errorf("shouldMarkGatewayNodeFailed() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUseGatewayLB(t *testing.T) {
	rapidHost := "bucket-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	ordinaryHost := "ordinary-1250000000.cos.ap-guangzhou.myqcloud.com"

	tests := []struct {
		name       string
		rawURL     string
		hostHeader string
		want       bool
	}{
		{
			name:       "Host 头为 rapid bucket，优先于 URL 生效",
			rawURL:     "http://" + ordinaryHost + "/key",
			hostHeader: rapidHost,
			want:       true,
		},
		{
			name:       "Host 头为普通桶，优先于 URL 生效（即使 URL 恰好是 rapid 命名）",
			rawURL:     "http://" + rapidHost + "/key",
			hostHeader: ordinaryHost,
			want:       false,
		},
		{
			name:       "Host 头带端口，剥掉端口后仍判定为 rapid bucket",
			rawURL:     "http://" + ordinaryHost + "/key",
			hostHeader: rapidHost + ":8443",
			want:       true,
		},
		{
			name:   "Host 头为空，回退 req.URL.Hostname()：rapid",
			rawURL: "http://" + rapidHost + "/key",
			want:   true,
		},
		{
			name:   "Host 头为空，回退 req.URL.Hostname()：普通桶",
			rawURL: "http://" + ordinaryHost + "/key",
			want:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := newGatewayLBTestRequest(t, tc.rawURL, tc.hostHeader)
			if got := useGatewayLB(req); got != tc.want {
				t.Errorf("useGatewayLB() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGatewayLBTransportEligible(t *testing.T) {
	for _, name := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
		t.Setenv(name, "")
	}

	tests := []struct {
		name      string
		rawURL    string
		transport http.RoundTripper
		want      bool
	}{
		{
			name:      "标准 HTTP 短域名满足门禁",
			rawURL:    "http://" + rapidGatewayTestHost + "/key",
			transport: &http.Transport{},
			want:      true,
		},
		{
			name:      "HTTPS 不介入",
			rawURL:    "https://" + rapidGatewayTestHost + "/key",
			transport: &http.Transport{},
			want:      false,
		},
		{
			name:      "显式端口不介入",
			rawURL:    "http://" + rapidGatewayTestHost + ":8080/key",
			transport: &http.Transport{},
			want:      false,
		},
		{
			name:      "显式默认端口同样不介入",
			rawURL:    "http://" + rapidGatewayTestHost + ":80/key",
			transport: &http.Transport{},
			want:      false,
		},
		{
			name:      "Host 为 IPv4 字面量不介入",
			rawURL:    "http://10.0.0.1/key",
			transport: &http.Transport{},
			want:      false,
		},
		{
			name:      "Host 为 IPv6 字面量不介入",
			rawURL:    "http://[2001:db8::1]/key",
			transport: &http.Transport{},
			want:      false,
		},
		{
			name:      "next 的 Proxy 为 nil 时视为无代理",
			rawURL:    "http://" + rapidGatewayTestHost + "/key",
			transport: &http.Transport{},
			want:      true,
		},
		{
			name:   "next 只要配置 Proxy 函数就不介入",
			rawURL: "http://" + rapidGatewayTestHost + "/key",
			transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lb := newRapidGatewayLBTransport(tc.transport)
			req := newGatewayLBTestRequest(t, tc.rawURL, "")
			if got := lb.eligible(req); got != tc.want {
				t.Errorf("eligible(%q) = %v, want %v", tc.rawURL, got, tc.want)
			}
		})
	}
}

func TestGatewayLBTransportProxyFunctionIsNotCalled(t *testing.T) {
	called := false
	transport := &http.Transport{Proxy: func(*http.Request) (*url.URL, error) {
		called = true
		return nil, errors.New("must not be called")
	}}
	lb := newRapidGatewayLBTransport(transport)
	req := newGatewayLBTestRequest(t, "http://"+rapidGatewayTestHost+"/key", "")
	if lb.eligible(req) {
		t.Fatal("eligible() = true with an explicit Proxy function, want false")
	}
	if called {
		t.Fatal("eligible() called the Proxy function, want conservative field-only gating")
	}
}

func TestHasHTTPProxyEnvironment(t *testing.T) {
	for _, name := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
		t.Setenv(name, "")
	}
	if hasHTTPProxyEnvironment() {
		t.Fatal("hasHTTPProxyEnvironment() = true with no proxy variables")
	}
	for _, name := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
		t.Run(name, func(t *testing.T) {
			for _, variable := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
				t.Setenv(variable, "")
			}
			t.Setenv(name, "http://127.0.0.1:3128")
			t.Setenv("NO_PROXY", rapidGatewayTestHost)
			if !hasHTTPProxyEnvironment() {
				t.Fatalf("hasHTTPProxyEnvironment() = false with %s set", name)
			}
		})
	}
}

func TestGatewayLBTransportDefaultTransportEnvironmentProxyGate(t *testing.T) {
	for _, name := range []string{"http_proxy", "HTTP_PROXY", "all_proxy", "ALL_PROXY"} {
		t.Setenv(name, "")
	}
	lb := newRapidGatewayLBTransport(nil)
	req := newGatewayLBTestRequest(t, "http://"+rapidGatewayTestHost+"/key", "")
	if !lb.eligible(req) {
		t.Fatal("eligible() = false with no proxy environment, want true")
	}

	t.Setenv("HTTP_PROXY", "http://127.0.0.1:3128")
	t.Setenv("NO_PROXY", rapidGatewayTestHost)
	if lb.eligible(req) {
		t.Fatal("eligible() = true with HTTP_PROXY set, want false even when NO_PROXY matches")
	}
}

func TestGatewayLBEntryPickFiltersCooldown(t *testing.T) {
	nodeA := net.ParseIP("10.0.0.1")
	nodeB := net.ParseIP("10.0.0.2")
	nodes := []net.IP{nodeA, nodeB}

	e := &gatewayLBEntry{}

	// 无冷却记录时两个节点都应能被选中。
	seen := make(map[string]bool)
	for i := 0; i < 200; i++ {
		ip, ok := e.pick(nodes)
		if !ok {
			t.Fatal("pick() returned ok=false while no node is cooling")
		}
		seen[ip.String()] = true
	}
	if !seen[nodeA.String()] || !seen[nodeB.String()] {
		t.Errorf("pick() only selected %v, want both nodes to be selectable", seen)
	}

	// 标记 A 冷却后，只应选出 B。
	e.markFailed(nodeA)
	for i := 0; i < 50; i++ {
		ip, ok := e.pick(nodes)
		if !ok {
			t.Fatal("pick() returned ok=false while node B is healthy")
		}
		if ip.String() != nodeB.String() {
			t.Fatalf("pick() = %s, want %s (node A is cooling)", ip, nodeB)
		}
	}
}

func TestGatewayLBEntryPickAllCooling(t *testing.T) {
	nodeA := net.ParseIP("10.0.0.1")
	nodeB := net.ParseIP("10.0.0.2")
	nodes := []net.IP{nodeA, nodeB}

	e := &gatewayLBEntry{}
	e.markFailed(nodeA)
	e.markFailed(nodeB)

	// 全部节点冷却中时不从冷却节点里硬选，交由调用方回退逻辑域名。
	if ip, ok := e.pick(nodes); ok {
		t.Errorf("pick() = (%s, true), want ok=false while all nodes are cooling", ip)
	}
}

func TestGatewayLBEntryCooldownExpires(t *testing.T) {
	node := net.ParseIP("10.0.0.1")
	nodes := []net.IP{node}

	e := &gatewayLBEntry{}
	// 直接写入一个已过期的冷却截止时间，避免用例等待真实的冷却时长。
	e.mu.Lock()
	e.cooldown = map[string]time.Time{node.String(): time.Now().Add(-time.Second)}
	e.mu.Unlock()

	ip, ok := e.pick(nodes)
	if !ok {
		t.Fatal("pick() returned ok=false, want the node to recover after cooldown expired")
	}
	if ip.String() != node.String() {
		t.Errorf("pick() = %s, want %s", ip, node)
	}
}

func TestGatewayLBEntryPruneCooldown(t *testing.T) {
	stale := net.ParseIP("10.0.0.9")
	kept := net.ParseIP("10.0.0.1")

	e := &gatewayLBEntry{}
	e.markFailed(stale)
	e.markFailed(kept)

	e.mu.Lock()
	e.pruneCooldown([]net.IP{kept})
	_, staleExists := e.cooldown[stale.String()]
	_, keptExists := e.cooldown[kept.String()]
	_, staleStreakExists := e.failStreak[stale.String()]
	_, keptStreakExists := e.failStreak[kept.String()]
	_, staleLastFailExists := e.lastFailAt[stale.String()]
	_, keptLastFailExists := e.lastFailAt[kept.String()]
	e.mu.Unlock()

	if staleExists {
		t.Errorf("pruneCooldown() kept cooldown record for %s, want it removed", stale)
	}
	if !keptExists {
		t.Errorf("pruneCooldown() removed cooldown record for %s, want it kept", kept)
	}
	if staleStreakExists || staleLastFailExists {
		t.Errorf("pruneCooldown() kept failStreak/lastFailAt record for %s, want it removed", stale)
	}
	if !keptStreakExists || !keptLastFailExists {
		t.Errorf("pruneCooldown() removed failStreak/lastFailAt record for %s, want it kept", kept)
	}
}

func TestGatewayLBEntryMarkFailedUsesBaseCooldownOnFirstFailure(t *testing.T) {
	node := net.ParseIP("10.0.0.1")
	e := &gatewayLBEntry{}
	before := time.Now()
	e.markFailed(node)

	e.mu.Lock()
	until := e.cooldown[node.String()]
	streak := e.failStreak[node.String()]
	e.mu.Unlock()

	if streak != 1 {
		t.Errorf("failStreak = %d, want 1 after the first failure", streak)
	}
	got := until.Sub(before)
	if got < gatewayLBDefaultBaseCooldown || got > gatewayLBDefaultBaseCooldown+time.Second {
		t.Errorf("first cooldown = %v, want approximately BaseCooldown (%v)", got, gatewayLBDefaultBaseCooldown)
	}
}

func TestGatewayLBEntryMarkFailedEscalatesExponentially(t *testing.T) {
	node := net.ParseIP("10.0.0.1")
	e := &gatewayLBEntry{}
	// 期望序列：5s, 10s, 20s, 40s, 60s(封顶), 60s(不再增长)。
	want := []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second, 40 * time.Second, 60 * time.Second, 60 * time.Second}
	for i, w := range want {
		before := time.Now()
		e.markFailed(node)
		e.mu.Lock()
		until := e.cooldown[node.String()]
		e.mu.Unlock()
		got := until.Sub(before)
		if got < w-time.Second || got > w+time.Second {
			t.Errorf("failure #%d cooldown = %v, want approximately %v", i+1, got, w)
		}
	}
}

func TestGatewayLBEntryConcurrentFailuresEscalatePerFailure(t *testing.T) {
	node := net.ParseIP("10.0.0.1")
	e := &gatewayLBEntry{}
	const failures = 16

	var wg sync.WaitGroup
	wg.Add(failures)
	for i := 0; i < failures; i++ {
		go func() {
			defer wg.Done()
			e.markFailed(node)
		}()
	}
	wg.Wait()

	key := node.String()
	e.mu.Lock()
	streak := e.failStreak[key]
	cooldown := e.cooldown[key].Sub(e.lastFailAt[key])
	e.mu.Unlock()
	if streak != failures {
		t.Errorf("failStreak = %d, want one increment per completed failure (%d)", streak, failures)
	}
	if cooldown != gatewayLBDefaultMaxCooldown {
		t.Errorf("cooldown = %v, want MaxCooldown %v", cooldown, gatewayLBDefaultMaxCooldown)
	}
}

func TestGatewayLBEntryMarkSucceededResetsStreak(t *testing.T) {
	node := net.ParseIP("10.0.0.1")
	e := &gatewayLBEntry{}
	e.markFailed(node) // streak=1, 5s
	e.markFailed(node) // streak=2, 10s
	e.markSucceeded(node)

	e.mu.Lock()
	_, streakExists := e.failStreak[node.String()]
	e.mu.Unlock()
	if streakExists {
		t.Fatal("markSucceeded() did not clear failStreak")
	}

	before := time.Now()
	e.markFailed(node) // 应重新从 base 起算，而不是延续之前的退避等级
	e.mu.Lock()
	until := e.cooldown[node.String()]
	e.mu.Unlock()
	got := until.Sub(before)
	if got < gatewayLBDefaultBaseCooldown || got > gatewayLBDefaultBaseCooldown+time.Second {
		t.Errorf("cooldown after markSucceeded = %v, want approximately BaseCooldown (%v)", got, gatewayLBDefaultBaseCooldown)
	}
}

func TestGatewayLBEntryMarkFailedResetsStreakAfterResetWindow(t *testing.T) {
	node := net.ParseIP("10.0.0.1")
	e := &gatewayLBEntry{}
	e.markFailed(node) // streak=1
	e.markFailed(node) // streak=2, 10s

	// 人为把 lastFailAt 设成超过 FailureResetWindow 之前，模拟"很久没再失败过"。
	e.mu.Lock()
	e.lastFailAt[node.String()] = time.Now().Add(-2 * gatewayLBDefaultFailureResetWindow)
	e.mu.Unlock()

	before := time.Now()
	e.markFailed(node) // 应视为全新序列，从 base 重新起算
	e.mu.Lock()
	until := e.cooldown[node.String()]
	streak := e.failStreak[node.String()]
	e.mu.Unlock()

	if streak != 1 {
		t.Errorf("failStreak = %d after exceeding FailureResetWindow, want 1 (new sequence)", streak)
	}
	got := until.Sub(before)
	if got < gatewayLBDefaultBaseCooldown || got > gatewayLBDefaultBaseCooldown+time.Second {
		t.Errorf("cooldown after reset window = %v, want approximately BaseCooldown (%v)", got, gatewayLBDefaultBaseCooldown)
	}
}

func TestEntryForNormalizedHostSharesEntry(t *testing.T) {
	host := "entry-share-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	first := entryFor(canonicalHost(host))
	second := entryFor(canonicalHost("Entry-Share-X--1259654469.Cosrapid.ap-guangzhou.myqcloud.com:80"))
	third := entryFor(canonicalHost(host + "."))

	if first != second || first != third {
		t.Error("entryFor() returned different entries for the same host written differently")
	}
}

func TestEntryForConcurrentCreatesSingleEntry(t *testing.T) {
	host := "entry-race-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	const goroutines = 32
	entries := make([]*gatewayLBEntry, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			entries[idx] = entryFor(canonicalHost(host))
		}(i)
	}
	wg.Wait()

	for i := 1; i < goroutines; i++ {
		if entries[i] != entries[0] {
			t.Fatalf("entryFor() returned %d distinct entries, want a single shared entry", i)
		}
	}
}

func TestGatewayLBEntryResolveReturnsUnexpiredSnapshot(t *testing.T) {
	host := "resolve-stale-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	node := net.ParseIP("10.0.0.1")
	e := seedGatewayLBEntry(host, node)

	// 快照未过期：直接返回内存中的快照，不触发任何解析。
	nodes := e.resolve(canonicalHost(host))
	if len(nodes) != 1 || nodes[0].String() != node.String() {
		t.Fatalf("resolve() = %v, want [%s]", nodes, node)
	}

	e.mu.Lock()
	resolving := e.resolving
	e.mu.Unlock()
	if resolving {
		t.Error("resolve() started a refresh while the snapshot is still fresh")
	}
}

func TestGatewayLBEntryResolveColdStartDoesNotBlock(t *testing.T) {
	e := &gatewayLBEntry{}
	start := time.Now()
	// 传空 host：标准库在进入解析器前就对空 host 短路返回错误，用例因此不产生网络活动。
	// 本用例验证的是 entry 层的调度行为，host 只是透传参数。
	nodes := e.resolve("")
	elapsed := time.Since(start)

	// 首次访问不等解析完成：立即返回空快照，由调用方回退逻辑域名。
	if len(nodes) != 0 {
		t.Errorf("resolve() = %v, want an empty snapshot before the first resolution completes", nodes)
	}
	if elapsed >= gatewayLBDefaultDNSTimeout {
		t.Errorf("resolve() took %v, want it to return without waiting for the DNS lookup", elapsed)
	}
	// resolving 由 resolve 在锁内同步置真；若后台解析已完成，则 expireAt 已被推进。
	// 两者任一成立即说明后台解析确实被触发了。
	e.mu.Lock()
	triggered := e.resolving || !e.expireAt.IsZero()
	e.mu.Unlock()
	if !triggered {
		t.Error("resolve() did not kick off a background resolution on first access")
	}
	waitForNoResolutionInFlight(t, e)
}

func TestGatewayLBEntryResolveRateLimitedAfterFailedResolution(t *testing.T) {
	e := &gatewayLBEntry{}
	// 复现一次解析失败后 refreshAsync 留下的状态：仍无快照，但刷新时刻已被推到未来。
	e.mu.Lock()
	e.expireAt = time.Now().Add(gatewayLBDefaultRefresh)
	e.mu.Unlock()

	// 刷新周期内的后续请求都拿到空快照且不再触发解析，否则域名持续解析失败时每个请求都会打一次
	// DNS。resolving 由 resolve 在锁内同步置真，故可同步断言。
	for i := 0; i < 10; i++ {
		if nodes := e.resolve(""); len(nodes) != 0 {
			t.Fatalf("resolve() = %v, want an empty snapshot", nodes)
		}
		e.mu.Lock()
		resolving := e.resolving
		e.mu.Unlock()
		if resolving {
			t.Fatal("resolve() started a resolution within the refresh interval, want it rate limited")
		}
	}
}

func TestGatewayLBEntryResolveDoesNotStackResolutions(t *testing.T) {
	e := &gatewayLBEntry{}
	// 预置"已有解析在途 + 已到刷新时刻"：此时任何调用都不应再起一次解析。
	e.mu.Lock()
	e.resolving = true
	e.mu.Unlock()

	const goroutines = 32
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			// 传空 host：一旦有调用误起解析，查询会立即失败并把 expireAt 推进。
			e.resolve("")
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond) // 留出时间让误起的解析暴露出来

	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.expireAt.IsZero() {
		t.Error("resolve() started another resolution while one was already in flight")
	}
}

func TestGatewayLBTransportRewritesConnectionTarget(t *testing.T) {
	host := "rewrite-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	node := net.ParseIP("10.0.0.1")
	seedGatewayLBEntry(host, node)

	stub := &stubRoundTripper{}
	lb := newRapidGatewayLBTransport(stub)
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)
	originalURL := req.URL

	resp, err := lb.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip() failed: %v", err)
	}
	resp.Body = http.NoBody

	_, gotURLHost, gotHost := stub.snapshot()
	if want := net.JoinHostPort(node.String(), "80"); gotURLHost != want {
		t.Errorf("req.URL.Host = %q, want %q", gotURLHost, want)
	}
	if gotHost != host {
		t.Errorf("req.Host = %q, want the logical host %q", gotHost, host)
	}
	// 改写不得影响调用方原始请求对象。
	if originalURL.Host != host {
		t.Errorf("caller's req.URL.Host was modified to %q, want %q", originalURL.Host, host)
	}
	if req.URL != originalURL {
		t.Error("RoundTrip() replaced the URL of the caller's request object")
	}
}

func TestGatewayLBTransportPinsHostHeaderWhenHostEmpty(t *testing.T) {
	host := "pin-host-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	seedGatewayLBEntry(host, net.ParseIP("10.0.0.1"))

	stub := &stubRoundTripper{}
	lb := newRapidGatewayLBTransport(stub)
	// 不设 req.Host：Host 为空时标准库会取 req.URL.Host，若不钉死就会把物理 IP 当 Host 头。
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", "")

	if _, err := lb.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip() failed: %v", err)
	}

	_, _, gotHost := stub.snapshot()
	if gotHost != host {
		t.Errorf("req.Host = %q, want the logical host %q", gotHost, host)
	}
}

func TestGatewayLBTransportPreservesHostCase(t *testing.T) {
	host := "Case-Host-X--1259654469.Cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	seedGatewayLBEntry(host, net.ParseIP("10.0.0.1"))

	stub := &stubRoundTripper{}
	lb := newRapidGatewayLBTransport(stub)
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)

	if _, err := lb.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip() failed: %v", err)
	}

	// 必须回写原始大小写：上一层签名用的是原始 req.Host，写回小写形式会与之不一致。
	_, _, gotHost := stub.snapshot()
	if gotHost != host {
		t.Errorf("req.Host = %q, want the original host %q", gotHost, host)
	}
}

func TestGatewayLBTransportPassthrough(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		host   string
	}{
		{
			name:   "普通桶不介入",
			rawURL: "http://ordinary-1250000000.cos.ap-guangzhou.myqcloud.com/key",
			host:   "ordinary-1250000000.cos.ap-guangzhou.myqcloud.com",
		},
		{
			name:   "HTTPS 不介入",
			rawURL: "https://" + rapidGatewayTestHost + "/key",
			host:   rapidGatewayTestHost,
		},
		{
			name:   "显式端口不介入",
			rawURL: "http://" + rapidGatewayTestHost + ":8080/key",
			host:   rapidGatewayTestHost + ":8080",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetGatewayLBEntry(tc.host)
			defer resetGatewayLBEntry(tc.host)
			// 预置快照：若 LB 误介入，改写结果会与原始 Host 不同，从而被断言捕获。
			seedGatewayLBEntry(tc.host, net.ParseIP("10.0.0.1"))

			stub := &stubRoundTripper{}
			lb := newRapidGatewayLBTransport(stub)
			req := newGatewayLBTestRequest(t, tc.rawURL, tc.host)
			wantURLHost := req.URL.Host

			if _, err := lb.RoundTrip(req); err != nil {
				t.Fatalf("RoundTrip() failed: %v", err)
			}

			calls, gotURLHost, gotHost := stub.snapshot()
			if calls != 1 {
				t.Fatalf("next.RoundTrip called %d times, want 1", calls)
			}
			if gotURLHost != wantURLHost {
				t.Errorf("req.URL.Host = %q, want it untouched (%q)", gotURLHost, wantURLHost)
			}
			if gotHost != tc.host {
				t.Errorf("req.Host = %q, want %q", gotHost, tc.host)
			}
		})
	}
}

func TestGatewayLBTransportFallsBackWhenAllNodesCooling(t *testing.T) {
	host := "all-cooling-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	nodeA := net.ParseIP("10.0.0.1")
	nodeB := net.ParseIP("10.0.0.2")
	e := seedGatewayLBEntry(host, nodeA, nodeB)
	e.markFailed(nodeA)
	e.markFailed(nodeB)

	stub := &stubRoundTripper{}
	lb := newRapidGatewayLBTransport(stub)
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)

	if _, err := lb.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip() failed: %v", err)
	}

	// 全部节点冷却时应回退到未改写的逻辑域名，而不是从冷却节点里硬选一个。
	_, gotURLHost, _ := stub.snapshot()
	if gotURLHost != host {
		t.Errorf("req.URL.Host = %q, want the logical host %q", gotURLHost, host)
	}
}

func TestGatewayLBTransportMarksFailedNode(t *testing.T) {
	host := "mark-failed-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	node := net.ParseIP("10.0.0.1")
	e := seedGatewayLBEntry(host, node)

	stub := &stubRoundTripper{resp: &http.Response{StatusCode: http.StatusBadGateway, Header: make(http.Header)}}
	lb := newRapidGatewayLBTransport(stub)
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)

	resp, err := lb.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip() failed: %v", err)
	}
	// 5xx 响应必须原样上抛，由上层决定如何处理。
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("resp.StatusCode = %d, want %d", resp.StatusCode, http.StatusBadGateway)
	}

	e.mu.Lock()
	_, cooling := e.cooldown[node.String()]
	e.mu.Unlock()
	if !cooling {
		t.Errorf("node %s was not marked as cooling after a 502 response", node)
	}
}

func TestGatewayLBTransportDoesNotMarkOnCanceledContext(t *testing.T) {
	host := "no-mark-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	node := net.ParseIP("10.0.0.1")
	e := seedGatewayLBEntry(host, node)
	e.markFailed(node)
	e.markFailed(node)

	key := node.String()
	e.mu.Lock()
	e.cooldown[key] = time.Now().Add(-time.Second)
	wantStreak := e.failStreak[key]
	wantLastFailAt := e.lastFailAt[key]
	wantCooldown := e.cooldown[key]
	e.mu.Unlock()

	stub := &stubRoundTripper{err: context.Canceled}
	lb := newRapidGatewayLBTransport(stub)
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)

	if _, err := lb.RoundTrip(req); !errors.Is(err, context.Canceled) {
		t.Fatalf("RoundTrip() error = %v, want context.Canceled", err)
	}

	e.mu.Lock()
	gotStreak := e.failStreak[key]
	gotLastFailAt := e.lastFailAt[key]
	gotCooldown := e.cooldown[key]
	e.mu.Unlock()
	if gotStreak != wantStreak {
		t.Errorf("failStreak after canceled request = %d, want unchanged value %d", gotStreak, wantStreak)
	}
	if !gotLastFailAt.Equal(wantLastFailAt) {
		t.Errorf("lastFailAt after canceled request = %v, want unchanged value %v", gotLastFailAt, wantLastFailAt)
	}
	if !gotCooldown.Equal(wantCooldown) {
		t.Errorf("cooldown after canceled request = %v, want unchanged value %v", gotCooldown, wantCooldown)
	}
}

func TestGatewayLBTransportDoesNotRetryInternally(t *testing.T) {
	host := "no-retry-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	seedGatewayLBEntry(host, net.ParseIP("10.0.0.1"), net.ParseIP("10.0.0.2"))

	stub := &stubRoundTripper{err: errors.New("dial tcp: connection refused")}
	lb := newRapidGatewayLBTransport(stub)
	req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)

	if _, err := lb.RoundTrip(req); err == nil {
		t.Fatal("RoundTrip() returned nil error, want the transport error to be propagated")
	}

	// 即使还有其他健康节点，也只发送一次：重试完全交给上层。
	if calls, _, _ := stub.snapshot(); calls != 1 {
		t.Errorf("next.RoundTrip called %d times, want exactly 1 (no internal retry)", calls)
	}
}

func TestGatewayLBTransportDistributesAcrossNodes(t *testing.T) {
	host := "spread-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com"
	resetGatewayLBEntry(host)
	defer resetGatewayLBEntry(host)

	nodeA := net.ParseIP("10.0.0.1")
	nodeB := net.ParseIP("10.0.0.2")
	seedGatewayLBEntry(host, nodeA, nodeB)

	recorder := &stubRoundTripper{}
	lb := newRapidGatewayLBTransport(recorder)
	targets := make(map[string]int)

	// 复用同一个 Transport 连发多次：选点是请求级的，不依赖是否新建连接。
	for i := 0; i < 200; i++ {
		req := newGatewayLBTestRequest(t, "http://"+host+"/key", host)
		if _, err := lb.RoundTrip(req); err != nil {
			t.Fatalf("RoundTrip() failed: %v", err)
		}
		_, gotURLHost, _ := recorder.snapshot()
		targets[gotURLHost]++
	}

	wantA := net.JoinHostPort(nodeA.String(), "80")
	wantB := net.JoinHostPort(nodeB.String(), "80")
	if targets[wantA] == 0 || targets[wantB] == 0 {
		t.Errorf("request distribution = %v, want both %s and %s to be selected", targets, wantA, wantB)
	}
}

func TestGatewayLBTransportNilNextUsesDefaultTransport(t *testing.T) {
	lb := newRapidGatewayLBTransport(nil)
	if lb.transport() != http.DefaultTransport {
		t.Error("transport() did not fall back to http.DefaultTransport when next is nil")
	}
}

func TestGatewayLBTransportNilURLPassesThrough(t *testing.T) {
	stub := &stubRoundTripper{}
	lb := newRapidGatewayLBTransport(stub)

	// 即使作为内部 RoundTripper，也必须透传手工构造的无 URL 请求而不是 panic。
	if _, err := lb.RoundTrip(&http.Request{Method: http.MethodGet}); err != nil {
		t.Fatalf("RoundTrip() failed: %v", err)
	}
	if calls, _, _ := stub.snapshot(); calls != 1 {
		t.Errorf("next.RoundTrip called %d times, want 1", calls)
	}
}

func TestNormalizeGatewayNodes(t *testing.T) {
	tests := []struct {
		name  string
		addrs []net.IPAddr
		want  []string
	}{
		{
			name:  "空输入",
			addrs: nil,
			want:  nil,
		},
		{
			name: "去掉完全重复的地址",
			addrs: []net.IPAddr{
				{IP: net.ParseIP("10.0.0.1")},
				{IP: net.ParseIP("10.0.0.1")},
				{IP: net.ParseIP("10.0.0.2")},
			},
			want: []string{"10.0.0.1", "10.0.0.2"},
		},
		{
			name: "同一 IPv4 的 4 字节与 16 字节表示视为重复",
			addrs: []net.IPAddr{
				{IP: net.ParseIP("10.0.0.1")},
				{IP: net.ParseIP("10.0.0.1").To4()},
			},
			want: []string{"10.0.0.1"},
		},
		{
			name: "过滤空地址",
			addrs: []net.IPAddr{
				{IP: nil},
				{IP: net.IP{}},
				{IP: net.ParseIP("10.0.0.1")},
			},
			want: []string{"10.0.0.1"},
		},
		{
			name: "保持解析器给出的顺序，IPv4/IPv6 不做区分",
			addrs: []net.IPAddr{
				{IP: net.ParseIP("2001:db8::1")},
				{IP: net.ParseIP("10.0.0.1")},
			},
			want: []string{"2001:db8::1", "10.0.0.1"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeGatewayNodes(tc.addrs)
			if len(got) != len(tc.want) {
				t.Fatalf("normalizeGatewayNodes() = %v, want %v", got, tc.want)
			}
			for i, ip := range got {
				if ip.String() != tc.want[i] {
					t.Errorf("node[%d] = %s, want %s", i, ip, tc.want[i])
				}
			}
		})
	}
}

func TestGatewayLBEntryPickIsUniformWithoutDuplicates(t *testing.T) {
	nodeA := net.ParseIP("10.0.0.1")
	nodeB := net.ParseIP("10.0.0.2")

	e := &gatewayLBEntry{}
	counts := make(map[string]int)
	const rounds = 4000
	for i := 0; i < rounds; i++ {
		ip, ok := e.pick([]net.IP{nodeA, nodeB})
		if !ok {
			t.Fatal("pick() returned ok=false while no node is cooling")
		}
		counts[ip.String()]++
	}

	// 无重复项时两个节点应各占约一半。阈值取 ±10%，足以捕获"某节点占两个候选位"造成的
	// 2:1 倾斜，同时不会因随机波动而偶发失败。
	for _, ip := range []net.IP{nodeA, nodeB} {
		got := counts[ip.String()]
		if got < rounds*4/10 || got > rounds*6/10 {
			t.Errorf("node %s selected %d/%d times, want roughly half", ip, got, rounds)
		}
	}
}

func TestLookupGatewayNodesHonorsContextTimeout(t *testing.T) {
	// 用一个已过期的 ctx 直接验证超时被接线到 resolver：若实现忽略 ctx，这个调用会真的发起
	// DNS 解析并对不可路由域名等到自身超时，测试将长时间挂住后失败。已过期的 ctx 使 resolver
	// 立即短路返回，因此本用例零网络活动、零等待，且与机器快慢无关。
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	if ips := lookupGatewayNodes(ctx, "timeout-probe-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com", gatewayLBDefaultDNSTimeout); ips != nil {
		t.Fatalf("lookupGatewayNodes() = %v, want nil for an expired context", ips)
	}
}
