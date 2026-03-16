package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the COS client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a cos.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	u, _ := url.Parse(server.URL)
	client = NewClient(&BaseURL{u, u, u, u, u, u}, nil)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

// Helper function to test that a value is marshalled to XML as expected.
func testXMLMarshal(t *testing.T, v interface{}, want string) {
	j, err := xml.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %v", v)
	}

	w := new(bytes.Buffer)
	err = xml.NewEncoder(w).Encode([]byte(want))
	if err != nil {
		t.Errorf("String is not valid json: %s", want)
	}

	if w.String() != string(j) {
		t.Errorf("xml.Marshal(%q) returned %s, want %s", v, j, w)
	}

	// now go the other direction and make sure things unmarshal as expected
	u := reflect.ValueOf(v).Interface()
	if err := xml.Unmarshal([]byte(want), u); err != nil {
		t.Errorf("Unable to unmarshal XML for %v", want)
	}

	if !reflect.DeepEqual(v, u) {
		t.Errorf("xml.Unmarshal(%q) returned %s, want %s", want, u, v)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil, nil)

	if got, want := c.BaseURL.ServiceURL.String(), defaultServiceBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, UserAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}
}

func TestNewBucketURL_secure_false(t *testing.T) {
	u, _ := NewBucketURL("bname-idx", "ap-guangzhou", false)
	got := u.String()
	want := "http://bname-idx.cos.ap-guangzhou.myqcloud.com"
	if got != want {
		t.Errorf("NewBucketURL is %v, want %v", got, want)
	}
	_, err := NewBucketURL("", "ap-guangzhou", false)
	if err == nil {
		t.Errorf("NewBucketURL should return error")
	}
	_, err = NewBucketURL("bname-idx", "", false)
	if err == nil {
		t.Errorf("NewBucketURL should return error")
	}
}

func TestNewBucketURL_secure_true(t *testing.T) {
	u, _ := NewBucketURL("bname-idx", "ap-guangzhou", true)
	got := u.String()
	want := "https://bname-idx.cos.ap-guangzhou.myqcloud.com"
	if got != want {
		t.Errorf("NewBucketURL is %v, want %v", got, want)
	}
}

func TestClient_doAPI(t *testing.T) {
	setup()
	defer teardown()

}

func TestNewAuthTime(t *testing.T) {
	a := NewAuthTime(time.Hour)
	if a.SignStartTime != a.KeyStartTime ||
		a.SignEndTime != a.SignEndTime ||
		a.SignStartTime.Add(time.Hour) != a.SignEndTime {
		t.Errorf("NewAuthTime request got %+v is not valid", a)
	}
}

func Test_addHeaderOptions(t *testing.T) {
	val := &XOptionalValue{
		&http.Header{},
	}
	val.Header.Add("key", "value")
	ctx := context.WithValue(context.Background(), XOptionalKey, val)
	res, err := addHeaderOptions(ctx, http.Header{}, nil)
	if err != nil {
		t.Errorf("addHeaderOptions return failed: %v", err)
	}
	if res.Get("key") != "value" {
		t.Errorf("addHeaderOptions failed")
	}
}

func Test_SwitchHost(t *testing.T) {
	u, _ := url.Parse("https://example-125000000.cos.ap-chengdu.myqcloud.com/123")
	res := toSwitchHost(u)
	want := "https://example-125000000.cos.ap-chengdu.tencentcos.cn/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}

	u, _ = url.Parse("https://example-125000000.cos.ap-chengdu.tencentcos.cn/123")
	res = toSwitchHost(u)
	want = "https://example-125000000.cos.ap-chengdu.tencentcos.cn/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}

	u, _ = url.Parse("https://service.cos.myqcloud.com/123")
	res = toSwitchHost(u)
	want = "https://service.cos.myqcloud.com/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}

	u, _ = url.Parse("https://example-125000000.file.myqcloud.com/123")
	res = toSwitchHost(u)
	want = "https://example-125000000.file.myqcloud.com/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}

	u, _ = url.Parse("http://example-125000000.cos.ap-chengdu.myqcloud.com:80/123")
	res = toSwitchHost(u)
	want = "http://example-125000000.cos.ap-chengdu.tencentcos.cn:80/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}

	u, _ = url.Parse("https://example-125000000.cos-website.ap-chengdu.myqcloud.com:443/123")
	res = toSwitchHost(u)
	want = "https://example-125000000.cos-website.ap-chengdu.myqcloud.com:443/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}

	u, _ = url.Parse("https://example-125000000.cos.accelerate.myqcloud.com:443/123")
	res = toSwitchHost(u)
	want = "https://example-125000000.cos.accelerate.myqcloud.com:443/123"
	if res.String() != want {
		t.Errorf("toSwitchHost failed, expect: %v, res: %v", want, res.String())
	}
}

func Test_CheckRetrieable(t *testing.T) {
	setup()
	defer teardown()

	u, _ := url.Parse("https://example-125000000.cos.ap-chengdu.myqcloud.com/123")
	wanted := "https://example-125000000.cos.ap-chengdu.tencentcos.cn/123"
	client.Conf.RetryOpt.AutoSwitchHost = true
	res, retry := client.CheckRetrieable(u, nil, errors.New("err"), true)
	if retry != true || res.String() != wanted {
		t.Errorf("CheckRetrieable failed, switch: %v, retry: %v", res.String(), retry)
	}
}

func Test_BaseURL(t *testing.T) {
	u, _ := url.Parse("https://example-125000000.cos.ap-chengdu.myqcloud.com")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos-website.ap-chengdu.myqcloud.com")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos-internal.ap-chengdu.tencentcos.cn")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos.ap-chengdu.tencentcos.cn")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos.accelerate.myqcloud.com")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos-internal.accelerate.tencentcos.cn")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://example-125000000.cos.ap-chengdu.myqcloud.com:8080")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://example-125000000.cos-internal.ap-chengdu.tencentcos.cn:80")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://cluster-1.cos-2.ap-guangzhou.myqcloud.com")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://test-1250000.global.tencentcos.cn")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://1.cos-c-internal.ap-singapore.tencentcos.cn")
	if !(&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}

	u, _ = url.Parse("https://example-125000000.cos.ap-chengdu@123.com/.myqcloud.com")
	if (&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos.ap-chengdu@123.com/.myqcloud.com:443")
	if (&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("https://example-125000000.cos.ap-chengdu@123.com/.myqcloud.com")
	if (&BaseURL{BucketURL: u, ServiceURL: u, BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}

	u, _ = url.Parse("https://service.cos.myqcloud.com")
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://service.cos-internal.tencentcos.cn")
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://service.cos.tencentcos.cn:80")
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://cos.ap-guangzhou.myqcloud.com")
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://cos.ap-guangzhou.myqcloud.com")
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://cos.ap-guangzhou.tencentcos.cn:8080")
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://service.cos@qq.com/.myqcloud.com")
	if (&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://service.cos@qq.com/.myqcloud.com")
	if (&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}

	u, _ = url.Parse("http://123.cos-control.ap-guangzhou.myqcloud.com")
	if !(&BaseURL{BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://123.cos-control.ap-guangzhou.tencentcos.cn")
	if !(&BaseURL{BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
	u, _ = url.Parse("http://123.cos-control.ap-guangzhou.myqcloud.com:8080")
	if !(&BaseURL{BatchURL: u}).Check() {
		t.Errorf("BaseURL check failed: %v", u)
	}
}

// ==================== hostSuffix 正则测试 ====================

func Test_HostSuffix(t *testing.T) {
	// 应匹配的 COS 标准域名后缀
	matchCases := []string{
		"https://bucket-123.cos.ap-guangzhou.myqcloud.com",
		"http://bucket-123.cos-internal.ap-beijing.myqcloud.com",
		"https://bucket-123.cos-website.ap-shanghai.myqcloud.com",
		"https://bucket-123.ci.ap-chengdu.myqcloud.com",
		"https://bucket-123.cos.ap-guangzhou.tencentcos.cn",
		"http://bucket-123.cos-internal.ap-beijing.tencentcos.cn",
		"https://bucket-123.cos-website.ap-shanghai.tencentcos.cn",
		"https://bucket-123.ci.ap-chengdu.tencentcos.cn",
		"https://file.myqcloud.com",
		"https://file.tencentcos.cn",
		"https://bucket-123.cos.accelerate.myqcloud.com",
		"https://cos.ap-guangzhou.myqcloud.com",
		"https://bucket-123.cos.ap-guangzhou.myqcloud.com/key/path",
	}
	for _, c := range matchCases {
		if !hostSuffix.MatchString(c) {
			t.Errorf("hostSuffix should match: %v", c)
		}
	}

	// 不应匹配的非 COS 域名
	noMatchCases := []string{
		"https://example.com",
		"https://bucket-123.s3.amazonaws.com",
		"https://myqcloud.com",
		"https://tencentcos.cn",
		"https://bucket.other.mysite.com",
	}
	for _, c := range noMatchCases {
		if hostSuffix.MatchString(c) {
			t.Errorf("hostSuffix should not match: %v", c)
		}
	}
}

// ==================== hostPrefix 正则测试 ====================

func Test_HostPrefix(t *testing.T) {
	// 合法的带 bucket-appid 前缀的域名
	matchCases := []string{
		"https://bucket-123456.cos.ap-guangzhou.myqcloud.com",
		"http://test-appid-999.cos-internal.ap-beijing.myqcloud.com",
		"https://my-bucket-123.cos-website.ap-shanghai.tencentcos.cn",
		"https://abc-0.ci.ap-chengdu.myqcloud.com",
		// 无 bucket 前缀但 hostPrefix 允许（{0,1} 表示可选）
		"https://cos.ap-guangzhou.myqcloud.com",
		"http://cos-internal.ap-beijing.tencentcos.cn",
		"https://file.myqcloud.com",
		"https://file.tencentcos.cn",
	}
	for _, c := range matchCases {
		if !hostPrefix.MatchString(c) {
			t.Errorf("hostPrefix should match: %v", c)
		}
	}

	// 不合法的前缀（bucket 名不以数字结尾的 appid 格式，或 SSRF 类攻击 URL）
	noMatchCases := []string{
		"https://bucket.cos.ap-guangzhou@attack.com/.myqcloud.com",
		"https://example.com",
	}
	for _, c := range noMatchCases {
		if hostPrefix.MatchString(c) {
			t.Errorf("hostPrefix should not match: %v", c)
		}
	}
}

// ==================== domainSuffix 正则测试 ====================

func Test_DomainSuffix(t *testing.T) {
	// 应匹配：以 .myqcloud.com 或 .tencentcos.cn 结尾（可带端口）
	matchCases := []string{
		"cos.ap-guangzhou.myqcloud.com",
		"bucket-123.cos.ap-guangzhou.myqcloud.com",
		"bucket-123.cos.ap-guangzhou.myqcloud.com:8080",
		"bucket-123.cos.ap-guangzhou.tencentcos.cn",
		"bucket-123.cos.ap-guangzhou.tencentcos.cn:443",
		"service.cos.myqcloud.com",
		"service.cos-internal.tencentcos.cn:80",
		"https://bucket-123.cos.ap-guangzhou-3.myqcloud.com",
		"http://bucket-123.cos.ap-guangzhou-2.tencentcos.cn:8080",
	}
	for _, c := range matchCases {
		if !domainSuffix.MatchString(c) {
			t.Errorf("domainSuffix should match: %v", c)
		}
	}

	// 不应匹配
	noMatchCases := []string{
		"myqcloud.com",         // 缺少子域名前缀的点
		"tencentcos.cn",        // 缺少子域名前缀的点
		"example.com",          // 非腾讯云域名
		"cos.amazonaws.com",    // 非腾讯云域名
		"bucket.myqcloud.org",  // 后缀不对
		"bucket.tencentcos.com", // 后缀不对
	}
	for _, c := range noMatchCases {
		if domainSuffix.MatchString(c) {
			t.Errorf("domainSuffix should not match: %v", c)
		}
	}
}

// ==================== bucketDomainChecker 正则测试 ====================

func Test_BucketDomainChecker(t *testing.T) {
	// 合法的 bucket 域名格式
	matchCases := []string{
		"https://bucket-123.cos.ap-guangzhou.myqcloud.com",
		"http://bucket-123.cos-internal.ap-beijing.tencentcos.cn",
		"https://bucket-123.cos.accelerate.myqcloud.com",
		"http://service.cos.myqcloud.com",
		"https://cos.ap-guangzhou.myqcloud.com",
		"http://bucket-123.cos.ap-guangzhou.myqcloud.com:8080",
		"https://bucket-123.cos.ap-guangzhou.tencentcos.cn:443",
		"https://cluster-1.cos-2.ap-guangzhou.myqcloud.com",
		"https://test-1250000.global.tencentcos.cn",
		"http://1.cos-c-internal.ap-singapore.tencentcos.cn",
		// 无 scheme
		"bucket-123.cos.ap-guangzhou.myqcloud.com",
	}
	for _, c := range matchCases {
		if !bucketDomainChecker.MatchString(c) {
			t.Errorf("bucketDomainChecker should match: %v", c)
		}
	}

	// 不合法的域名格式（SSRF 攻击类）
	noMatchCases := []string{
		"https://cos@qq.com/.myqcloud.com",
		"https://bucket.cos.ap-chengdu@123.com/.myqcloud.com",
		"https://bucket.cos.ap-chengdu@123.com/.myqcloud.com:443",
		// 大写字母
		"https://Bucket-123.cos.ap-guangzhou.myqcloud.com",
		// 特殊字符
		"https://bucket_123.cos.ap-guangzhou.myqcloud.com",
	}
	for _, c := range noMatchCases {
		if bucketDomainChecker.MatchString(c) {
			t.Errorf("bucketDomainChecker should not match: %v", c)
		}
	}
}

// ==================== checkURL 函数测试 ====================

func Test_CheckURL(t *testing.T) {
	// nil URL 应返回 false
	if checkURL(nil) {
		t.Errorf("checkURL(nil) should return false")
	}

	// 合法 URL：COS 标准域名带 bucket-appid 前缀
	validCases := []string{
		"https://bucket-123.cos.ap-guangzhou.myqcloud.com",
		"http://bucket-123.cos-internal.ap-beijing.myqcloud.com",
		"https://bucket-123.cos-website.ap-shanghai.tencentcos.cn",
		"https://bucket-123.ci.ap-chengdu.myqcloud.com",
		"https://bucket-123.cos.accelerate.myqcloud.com",
		"https://file.myqcloud.com",
		"https://file.tencentcos.cn",
		// 非 COS 域名（自定义域名），checkURL 应返回 true
		"https://cdn.example.com",
		"https://my-custom-domain.com/path",
		// 无 bucket 前缀的 COS 服务域名也是合法的
		"https://cos.ap-guangzhou.myqcloud.com",
	}
	for _, c := range validCases {
		u, _ := url.Parse(c)
		if !checkURL(u) {
			t.Errorf("checkURL should return true for: %v", c)
		}
	}

	// SSRF 类攻击 URL: url.Parse 后 hostname 为 attack.com（@ 前为 userinfo），
	// 不匹配 COS 域名后缀，checkURL 不拦截（此类场景由 BaseURL.Check 拦截）
	ssrfCases := []string{
		"https://bucket.cos.ap-guangzhou@attack.com/.myqcloud.com",
	}
	for _, c := range ssrfCases {
		u, _ := url.Parse(c)
		if !checkURL(u) {
			t.Errorf("checkURL should return true for non-COS hostname SSRF URL (handled by BaseURL.Check): %v", c)
		}
	}

	// 缺少 Scheme
	u, _ := url.Parse("bucket-123.cos.ap-guangzhou.myqcloud.com")
	if checkURL(u) {
		t.Errorf("checkURL should return false when scheme is empty")
	}
}

// ==================== BaseURL.Check + innerCheck 综合测试 ====================

func Test_BaseURL_InnerCheck(t *testing.T) {
	// nil URL 字段应通过（innerCheck 对 nil 返回 true）
	if !(&BaseURL{}).Check() {
		t.Errorf("BaseURL with all nil should pass Check")
	}

	// 只设置部分字段为 nil，其他为合法
	u, _ := url.Parse("https://bucket-123.cos.ap-guangzhou.myqcloud.com")
	if !(&BaseURL{BucketURL: u}).Check() {
		t.Errorf("BaseURL with only BucketURL should pass Check")
	}
	if !(&BaseURL{ServiceURL: u}).Check() {
		t.Errorf("BaseURL with only ServiceURL should pass Check")
	}
	if !(&BaseURL{BatchURL: u}).Check() {
		t.Errorf("BaseURL with only BatchURL should pass Check")
	}

	// 无 scheme 应失败
	noScheme, _ := url.Parse("bucket-123.cos.ap-guangzhou.myqcloud.com")
	if (&BaseURL{BucketURL: noScheme}).Check() {
		t.Errorf("BaseURL without scheme should fail Check")
	}

	// 非腾讯云域名（不匹配 domainSuffix），不校验格式，应通过
	custom, _ := url.Parse("https://my-custom-cdn.example.com")
	if !(&BaseURL{BucketURL: custom}).Check() {
		t.Errorf("BaseURL with custom domain should pass Check")
	}

	// 腾讯云域名但格式非法（SSRF 类攻击）
	ssrf, _ := url.Parse("https://cos@qq.com/.myqcloud.com")
	if (&BaseURL{BucketURL: ssrf}).Check() {
		t.Errorf("BaseURL with SSRF URL in BucketURL should fail Check")
	}
	if (&BaseURL{ServiceURL: ssrf}).Check() {
		t.Errorf("BaseURL with SSRF URL in ServiceURL should fail Check")
	}
	if (&BaseURL{BatchURL: ssrf}).Check() {
		t.Errorf("BaseURL with SSRF URL in BatchURL should fail Check")
	}

	// 带端口的合法 URL
	withPort, _ := url.Parse("http://bucket-123.cos.ap-guangzhou.myqcloud.com:8080")
	if !(&BaseURL{BucketURL: withPort}).Check() {
		t.Errorf("BaseURL with port should pass Check: %v", withPort)
	}

	// 带端口的 tencentcos.cn 域名
	withPortTcn, _ := url.Parse("https://bucket-123.cos.ap-guangzhou.tencentcos.cn:443")
	if !(&BaseURL{BucketURL: withPortTcn}).Check() {
		t.Errorf("BaseURL with tencentcos.cn port should pass Check: %v", withPortTcn)
	}

	// 三个字段其中一个非法，整体应失败
	valid, _ := url.Parse("https://bucket-123.cos.ap-guangzhou.myqcloud.com")
	if (&BaseURL{BucketURL: valid, ServiceURL: ssrf, BatchURL: valid}).Check() {
		t.Errorf("BaseURL should fail if any URL is invalid")
	}
	if (&BaseURL{BucketURL: ssrf, ServiceURL: valid, BatchURL: valid}).Check() {
		t.Errorf("BaseURL should fail if BucketURL is invalid")
	}
	if (&BaseURL{BucketURL: valid, ServiceURL: valid, BatchURL: ssrf}).Check() {
		t.Errorf("BaseURL should fail if BatchURL is invalid")
	}

	// 带路径后缀的合法 URL（尾部 / 会被 TrimRight 去除）
	withSlash, _ := url.Parse("https://bucket-123.cos.ap-guangzhou.myqcloud.com/")
	if !(&BaseURL{BucketURL: withSlash}).Check() {
		t.Errorf("BaseURL with trailing slash should pass Check: %v", withSlash)
	}
}
