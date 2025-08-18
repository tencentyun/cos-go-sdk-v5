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

func TestUniqueGrantID(t *testing.T) {
	ids := []string{"abc", "abc", "ab"}
	actual := uniqueGrantID(ids)
	if actual != "abc,ab" {
		t.Errorf("expect uniqueIDs to be abc,ab, got %v", actual)
	}
}
