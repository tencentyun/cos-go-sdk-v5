package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const defaultSessionRefreshBefore = 90 * time.Second

// sessionTokenExpiredCode 是服务端对 session 凭证过期返回的错误码。数据面请求用缓存
// 凭证签名，因网络延迟到达服务端时凭证可能刚好过期，服务端据此返回该码，客户端识别后刷新
// 凭证重试一次。
const sessionTokenExpiredCode = "ExpiredToken"

// rapidBucketRE 匹配 rapid bucket 桶名格式：<name>-x--<appid>。
var rapidBucketRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*-x--[0-9]+$`)

type sessionAuthContextKey struct{}

func withSessionAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, sessionAuthContextKey{}, true)
}

func isSessionAuth(ctx context.Context) bool {
	v, _ := ctx.Value(sessionAuthContextKey{}).(bool)
	return v
}

func (c *Client) useSessionAuth() bool {
	host := c.Host
	if host == "" && c.BaseURL != nil && c.BaseURL.BucketURL != nil {
		host = c.BaseURL.BucketURL.Host
	}
	bucketName := parseBucketNameFromHost(host)
	return bucketName != "" && isRapidBucket(bucketName)
}

type CreateSessionResult struct {
	XMLName     xml.Name          `xml:"CreateSessionResult"`
	Credentials SessionCredential `xml:"Credentials"`
}

// SessionCredential 是 CreateSession 返回的rapid bucket session 凭证。
type SessionCredential struct {
	SecretID     string    `xml:"AccessKeyId"`
	SecretKey    string    `xml:"SecretAccessKey"`
	SessionToken string    `xml:"SessionToken"`
	Expiration   time.Time `xml:"Expiration"`
}

// createSessionProvider 按桶签发并缓存 session 凭证，由 RapidCredentialTransport 持有，可供多个单桶 Client 共享。
type createSessionProvider struct {
	owner         *RapidCredentialTransport
	refreshBefore time.Duration

	mu             sync.Mutex
	bucketSessions map[string]*bucketSession
}

type bucketSession struct {
	mu   sync.Mutex
	cred *SessionCredential
	// baseSecretID 记录签发该 session 时的基础 SecretID，变化即视为凭证已切换，需重签。
	baseSecretID string
}

func newCreateSessionProvider(owner *RapidCredentialTransport) *createSessionProvider {
	return &createSessionProvider{
		owner:          owner,
		refreshBefore:  defaultSessionRefreshBefore,
		bucketSessions: make(map[string]*bucketSession),
	}
}

// baseTransport 返回 owner 当前持有的底层 transport。
func (p *createSessionProvider) baseTransport() BaseTransport {
	if p.owner == nil {
		return nil
	}
	return p.owner.baseTransport
}

// baseSecretID 返回底层 transport 当前的基础 SecretID，用于 session 缓存的身份指纹比对。
func (p *createSessionProvider) baseSecretID() string {
	bt := p.baseTransport()
	if bt == nil {
		return ""
	}
	ak, _, _, _ := bt.GetCredential()
	return ak
}

// credential 返回指定请求桶的 session 凭证（获取或复用缓存）。userAgent 用于 createSession 内部临时 Client 的 User-Agent，传空则用 SDK 默认值。
func (p *createSessionProvider) credential(ctx context.Context, requestURL *url.URL, host string, userAgent string) (*SessionCredential, bool, error) {
	return p.getSessionCredential(ctx, requestURL, host, p.refreshBefore, userAgent)
}

// freshCredential 返回剩余寿命不低于 minRemaining 的 session 凭证，不满足时强制刷新一次，仍不足则报错。用于预签名。
func (p *createSessionProvider) freshCredential(ctx context.Context, requestURL *url.URL, host string, minRemaining time.Duration, userAgent string) (*SessionCredential, bool, error) {
	cred, ok, err := p.getSessionCredential(ctx, requestURL, host, minRemaining, userAgent)
	if err != nil || !ok {
		return cred, ok, err
	}
	// 缺失 Expiration 时无法确认覆盖 minRemaining，按错误处理。
	if cred.Expiration.IsZero() {
		return nil, false, fmt.Errorf("session credential expiration is missing")
	}
	if time.Now().Add(minRemaining).After(cred.Expiration) {
		return nil, false, fmt.Errorf("session credential lifetime (until %s) is shorter than required %s", cred.Expiration.Format(time.RFC3339), minRemaining)
	}
	return cred, true, nil
}

// getBucketSession 返回指定桶的 session 缓存条目，不存在时创建一个空条目。
func (p *createSessionProvider) getBucketSession(bucketName string) *bucketSession {
	p.mu.Lock()
	defer p.mu.Unlock()
	bs := p.bucketSessions[bucketName]
	if bs == nil {
		bs = &bucketSession{}
		p.bucketSessions[bucketName] = bs
	}
	return bs
}

// getSessionCredential 在缓存凭证剩余寿命大于 reuseWindow 时复用，否则刷新。
func (p *createSessionProvider) getSessionCredential(ctx context.Context, requestURL *url.URL, host string, reuseWindow time.Duration, userAgent string) (*SessionCredential, bool, error) {
	bucketName := parseBucketNameFromHost(host)
	if bucketName == "" || !isRapidBucket(bucketName) {
		return nil, false, nil
	}
	bs := p.getBucketSession(bucketName)
	bs.mu.Lock()
	defer bs.mu.Unlock()

	// 命中缓存需同时满足：剩余寿命大于 reuseWindow，且基础 SecretID 未变。
	curSecretID := p.baseSecretID()
	if bs.cred != nil && bs.baseSecretID == curSecretID && time.Now().Add(reuseWindow).Before(bs.cred.Expiration) {
		return bs.cred, true, nil
	}

	cred, err := p.createSession(ctx, requestURL, host, userAgent)
	if err != nil {
		// 刷新失败时，若旧凭证仍在有效期内且基础 SecretID 未变则继续复用，容忍瞬时失败或限频。
		if bs.cred != nil && bs.baseSecretID == curSecretID && time.Now().Before(bs.cred.Expiration) {
			return bs.cred, true, nil
		}
		return nil, false, err
	}
	bs.cred = cred
	bs.baseSecretID = curSecretID
	return cred, true, nil
}

func (p *createSessionProvider) createSession(ctx context.Context, requestURL *url.URL, host string, userAgent string) (*SessionCredential, error) {
	if requestURL == nil {
		return nil, fmt.Errorf("request URL is nil")
	}
	bt := p.baseTransport()
	if bt == nil {
		return nil, fmt.Errorf("base transport is nil")
	}
	bucketURL := &url.URL{Scheme: requestURL.Scheme, Host: requestURL.Host}
	client := NewClient(
		&BaseURL{BucketURL: bucketURL},
		&http.Client{Transport: bt},
	)
	client.Host = host
	if userAgent != "" {
		client.UserAgent = userAgent
	}
	// not retry
	client.Conf.RetryOpt.Count = 0
	result, _, err := client.Bucket.CreateSession(ctx)
	if err != nil {
		return nil, err
	}
	return &result.Credentials, nil
}

func parseBucketNameFromHost(host string) string {
	if host == "" {
		return ""
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	if i := strings.IndexByte(host, '.'); i >= 0 {
		return host[:i]
	}
	return host
}

func isRapidBucket(bucketName string) bool {
	return rapidBucketRE.MatchString(bucketName)
}

// BaseTransport 是 RapidCredentialTransport 的底层 transport 约束：提供基础凭证、发送请求，并发送已由上层签名的请求。
// SDK 内置的 AuthorizationTransport、CVMCredentialTransport、StsCredentialTransport、CredentialTransport 均满足该接口。
type BaseTransport interface {
	TransportIface
	http.RoundTripper

	// RoundTripSigned 发送 Authorization 已由调用方签好名的请求，仅进入网络发送层，不得获取凭证或修改签名头。
	RoundTripSigned(req *http.Request) (*http.Response, error)
}

// RapidCredentialTransport 是rapid bucket专用 Transport：组合一个底层 transport 提供基础凭证，CreateSession 用基础凭证签名，数据面请求用签发的 session 凭证签名。基础凭证切换后旧 session 因指纹失配自动重签。
//
// 用法：经 NewRapidCredentialTransport(base) 构造，再作为 http.Client.Transport 传入 NewClient。更新凭证请调用底层 transport 自身的方法。
type RapidCredentialTransport struct {
	// baseTransport 提供基础凭证，构造时注入且此后不可变。
	baseTransport BaseTransport

	provider *createSessionProvider
}

// NewRapidCredentialTransport 构造rapid bucket专用 Transport，注入提供基础凭证的底层 transport。base 允许为 nil，此时后续请求返回 "base transport is nil" 错误。
func NewRapidCredentialTransport(base BaseTransport) *RapidCredentialTransport {
	t := &RapidCredentialTransport{baseTransport: base}
	t.provider = newCreateSessionProvider(t)
	return t
}

// getSessionProvider 返回按桶缓存 session 凭证的 provider。
func (t *RapidCredentialTransport) getSessionProvider() *createSessionProvider {
	return t.provider
}

// GetCredential 透传底层基础凭证，满足 TransportIface。
func (t *RapidCredentialTransport) GetCredential() (string, string, string, error) {
	if t.baseTransport == nil {
		return "", "", "", fmt.Errorf("base transport is nil")
	}
	return t.baseTransport.GetCredential()
}

// RoundTrip 实现 http.RoundTripper：非 session 请求或目标非rapid bucket委托底层以基础凭证签名发送。
// 否则用 session 凭证签名后交由底层 RoundTripSigned 发送；若服务端返回 session 凭证过期
// （ExpiredToken），刷新凭证后重发一次，以容忍缓存凭证因网络延迟到达时刚好过期的场景。
func (t *RapidCredentialTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.baseTransport == nil {
		return nil, fmt.Errorf("base transport is nil")
	}

	// 目标 host 优先取 req.Host，为空则回退 req.URL.Hostname()。
	host := req.Host
	if host == "" {
		host = req.URL.Hostname()
	}
	// 非 session 请求（含 CreateSession 自身）或目标非rapid bucket，委托底层以基础凭证签名。
	if !isSessionAuth(req.Context()) || !isRapidBucket(parseBucketNameFromHost(host)) {
		return t.baseTransport.RoundTrip(req)
	}

	req = cloneRequest(req) // per RoundTrip contract

	cred, ok, err := t.provider.credential(req.Context(), req.URL, host, req.Header.Get("User-Agent"))
	if err != nil {
		return nil, err
	}
	// 前置已判定为rapid bucket，provider 必得相同结论，ok=false 属内部不一致，显式报错。
	if !ok {
		return nil, fmt.Errorf("cos: no session credential resolved for rapid bucket host %q", host)
	}

	// 用 session 凭证签名后交由底层 RoundTripSigned 发送，仅发送不再签名。
	authTime := NewAuthTime(defaultAuthExpire)
	AddAuthorizationHeader(cred.SecretID, cred.SecretKey, cred.SessionToken, req, authTime)
	resp, err := t.baseTransport.RoundTripSigned(req)
	if err != nil {
		return nil, err
	}

	// session 凭证过期（缓存凭证因网络延迟到达服务端时刚好过期）：把过期响应改写成 503 交给
	// 上层重试——ExpiredToken 本身是 400，上层通常不重试 4xx。
	if isExpiredTokenResponse(resp) {
		rewriteToServiceUnavailable(resp)
	}
	return resp, nil
}

// isExpiredTokenResponse 判断响应是否为 session 凭证过期（错误码 ExpiredToken）。读取并复原
// 响应体，确保未命中时上层仍可正常读取。
func isExpiredTokenResponse(resp *http.Response) bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	// 过期是非 2xx 错误响应，2xx 直接跳过不动 body。
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return false
	}
	data, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	// 无论解析成功与否都复原 body，避免上层丢 body。
	resp.Body = io.NopCloser(bytes.NewReader(data))
	if err != nil {
		return false
	}
	var e ErrorResponse
	if xml.Unmarshal(data, &e) != nil {
		return false
	}
	return e.Code == sessionTokenExpiredCode
}

// rewriteToServiceUnavailable 把响应状态码改写成 503，使上层按可重试处理；保留原始响应体不变，
// 上层仍可从 body 读到 ExpiredToken 原始错误码定位根因。
func rewriteToServiceUnavailable(resp *http.Response) {
	resp.StatusCode = http.StatusServiceUnavailable
	resp.Status = fmt.Sprintf("%d %s", http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
}
