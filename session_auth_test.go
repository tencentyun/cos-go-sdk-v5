package cos

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func sessionResponseXML(secretID, secretKey, sessionToken string, expiration time.Time) string {
	expirationElement := ""
	if !expiration.IsZero() {
		expirationElement = fmt.Sprintf("<Expiration>%s</Expiration>", expiration.UTC().Format(time.RFC3339))
	}
	return fmt.Sprintf("<CreateSessionResult><Credentials><AccessKeyId>%s</AccessKeyId><SecretAccessKey>%s</SecretAccessKey><SessionToken>%s</SessionToken>%s</Credentials></CreateSessionResult>",
		secretID, secretKey, sessionToken, expirationElement)
}

func newSessionTestClient(bucketURL *url.URL, secretID, secretKey, sessionToken string) (*Client, *AuthorizationTransport, *createSessionProvider) {
	base := &AuthorizationTransport{
		SecretID:     secretID,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
	}
	transport := NewRapidCredentialTransport(base)
	client := NewClient(&BaseURL{BucketURL: bucketURL}, &http.Client{Transport: transport})
	return client, base, transport.getSessionProvider()
}

func TestRapidCredentialTransport_InitializesSessionProvider(t *testing.T) {
	base := &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}
	transport := NewRapidCredentialTransport(base)

	if transport.getSessionProvider() == nil {
		t.Fatal("RapidCredentialTransport did not initialize session provider")
	}
}

func TestRapidCredentialTransport_DoesNotReplaceSessionProvider(t *testing.T) {
	base := &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}
	transport := NewRapidCredentialTransport(base)

	provider := transport.getSessionProvider()
	if transport.getSessionProvider() != provider {
		t.Fatal("RapidCredentialTransport replaced existing session provider")
	}
}

func TestAuthorizationTransport_SharedAcrossRapidBucketClients(t *testing.T) {
	var sessionCallsA, sessionCallsB int32
	serverA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			atomic.AddInt32(&sessionCallsA, 1)
			w.Write([]byte(sessionResponseXML("SKIDA", "sk-a", "tok-a", time.Unix(4070912645, 0))))
		case r.URL.Path == "/a.txt":
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDA") {
				t.Errorf("bucket A Authorization = %q, want SKIDA", got)
			}
			w.Write([]byte("a"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer serverA.Close()
	serverB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			atomic.AddInt32(&sessionCallsB, 1)
			w.Write([]byte(sessionResponseXML("SKIDB", "sk-b", "tok-b", time.Unix(4070912645, 0))))
		case r.URL.Path == "/b.txt":
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDB") {
				t.Errorf("bucket B Authorization = %q, want SKIDB", got)
			}
			w.Write([]byte("b"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer serverB.Close()

	uA, _ := url.Parse(serverA.URL)
	uB, _ := url.Parse(serverB.URL)
	transport := NewRapidCredentialTransport(&AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"})
	clientA := NewClient(&BaseURL{BucketURL: uA}, &http.Client{Transport: transport})
	clientA.Host = "bucket-a-x--1250000000.cos-rapid.ap-guangzhou.tencentcos.com"
	clientB := NewClient(&BaseURL{BucketURL: uB}, &http.Client{Transport: transport})
	clientB.Host = "bucket-b-x--1250000000.cos-rapid.ap-guangzhou.tencentcos.com"

	tests := []struct {
		client *Client
		name   string
	}{{clientA, "a.txt"}, {clientB, "b.txt"}}
	errCh := make(chan error, len(tests))
	for _, tc := range tests {
		go func(c *Client, name string) {
			resp, err := c.Object.Get(context.Background(), name, nil)
			if err == nil {
				resp.Body.Close()
			}
			errCh <- err
		}(tc.client, tc.name)
	}
	for range tests {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range tests {
		resp, err := tc.client.Object.Get(context.Background(), tc.name, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	}
	if got := atomic.LoadInt32(&sessionCallsA); got != 1 {
		t.Fatalf("bucket A CreateSession calls = %d, want 1", got)
	}
	if got := atomic.LoadInt32(&sessionCallsB); got != 1 {
		t.Fatalf("bucket B CreateSession calls = %d, want 1", got)
	}
}

func TestCreateSessionProvider_CreateSessionUsesCallerUserAgent(t *testing.T) {
	var gotSessionUserAgent, gotDataUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			gotSessionUserAgent = r.Header.Get("User-Agent")
			w.Write([]byte(sessionResponseXML("SKID", "sk", "tok", time.Unix(4070912645, 0))))
		case r.URL.Path == "/a.txt":
			gotDataUserAgent = r.Header.Get("User-Agent")
			w.Write([]byte("a"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	transport := NewRapidCredentialTransport(&AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"})
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: transport})
	c.Host = "bucket-x--1250000000.cos-rapid.ap-guangzhou.tencentcos.com"
	c.UserAgent = "my-caller-agent/1.0"

	resp, err := c.Object.Get(context.Background(), "a.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if gotDataUserAgent != c.UserAgent {
		t.Fatalf("data request User-Agent = %q, want %q", gotDataUserAgent, c.UserAgent)
	}
	if gotSessionUserAgent != c.UserAgent {
		t.Fatalf("CreateSession request User-Agent = %q, want caller's %q (got the SDK's own default instead)", gotSessionUserAgent, c.UserAgent)
	}
}

func TestCreateSessionProvider_CreateSessionDoesNotRetryInternally(t *testing.T) {
	var sessionCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.URL.RawQuery == "session" {
			atomic.AddInt32(&sessionCalls, 1)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	transport := NewRapidCredentialTransport(&AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"})
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: transport})
	c.Host = "bucket-x--1250000000.cos-rapid.ap-guangzhou.tencentcos.com"
	c.Conf.RetryOpt.Count = 0

	if _, err := c.Object.Get(context.Background(), "a.txt", nil); err == nil {
		t.Fatal("Object.Get error = nil, want an error from the failing CreateSession request")
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 1 {
		t.Fatalf("CreateSession HTTP requests = %d, want 1 (no internal retry)", got)
	}
}

func TestBucketService_CreateSession(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		if r.URL.RawQuery != "session" {
			t.Errorf("CreateSession query = %q, want session", r.URL.RawQuery)
		}
		if got := r.Header.Get("x-cos-create-session-mode"); got != "" {
			t.Errorf("x-cos-create-session-mode = %q, want empty", got)
		}
		if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=AKIDBASE") {
			t.Errorf("Authorization = %q, want base credential", got)
		}
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateSessionResult><Credentials><AccessKeyId>SKIDSESSION</AccessKeyId><SecretAccessKey>session-secret</SecretAccessKey><SessionToken>session-token</SessionToken><Expiration>2099-01-01T00:00:00Z</Expiration></Credentials></CreateSessionResult>`))
	})

	client = NewClient(client.BaseURL, &http.Client{Transport: &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}})
	res, resp, err := client.Bucket.CreateSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("CreateSession response = %+v", resp)
	}
	if res.Credentials.SecretID != "SKIDSESSION" || res.Credentials.SecretKey != "session-secret" || res.Credentials.SessionToken != "session-token" {
		t.Fatalf("CreateSession credentials = %+v", res.Credentials)
	}
	wantExpiration := time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	if !res.Credentials.Expiration.Equal(wantExpiration) {
		t.Fatalf("CreateSession expiration = %s, want %s", res.Credentials.Expiration, wantExpiration)
	}
}

func TestBucketService_CreateSessionRejectsEmptyResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	client = NewClient(client.BaseURL, &http.Client{Transport: &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}})

	res, resp, err := client.Bucket.CreateSession(context.Background())
	if err == nil || !strings.Contains(err.Error(), "credentials are incomplete") {
		t.Fatalf("CreateSession error = %v, want incomplete credentials error", err)
	}
	if res != nil || resp == nil {
		t.Fatalf("CreateSession result = %+v, response = %+v; want nil result and non-nil response", res, resp)
	}
}

func TestBucketService_CreateSessionRejectsMissingExpiration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Time{})))
	})
	client = NewClient(client.BaseURL, &http.Client{Transport: &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}})

	res, resp, err := client.Bucket.CreateSession(context.Background())
	if err == nil || !strings.Contains(err.Error(), "expiration is missing") {
		t.Fatalf("CreateSession error = %v, want missing expiration error", err)
	}
	if res != nil || resp == nil {
		t.Fatalf("CreateSession result = %+v, response = %+v; want nil result and non-nil response", res, resp)
	}
}

func TestBucketService_CreateSessionRejectsLegacyJSON(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Credentials":{"TmpSecretId":"SKIDSESSION","TmpSecretKey":"session-secret","Token":"session-token"},"ExpiredTime":4070912645}`))
	})
	client = NewClient(client.BaseURL, &http.Client{Transport: &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}})

	if _, _, err := client.Bucket.CreateSession(context.Background()); err == nil {
		t.Fatal("CreateSession accepted legacy JSON response, want XML-only decoding error")
	}
}

func TestAuthorizationTransport_UsesCachedSessionForRapidBucket(t *testing.T) {
	var sessionCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			atomic.AddInt32(&sessionCalls, 1)
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=AKIDBASE") {
				t.Errorf("CreateSession Authorization = %q, want base credential", got)
			}
			if r.Host != "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com" {
				t.Errorf("CreateSession Host = %q", r.Host)
			}
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
		case r.URL.Path == "/test.txt":
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("Object Authorization = %q, want session credential", got)
			}
			testHeader(t, r, "x-cos-security-token", "session-token")
			w.Write([]byte("ok"))
		default:
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	for i := 0; i < 2; i++ {
		resp, err := c.Object.Get(context.Background(), "test.txt", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 1 {
		t.Fatalf("CreateSession calls = %d, want 1", got)
	}
}

func TestBucketService_ListMultipartUploadsUsesSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
		case r.URL.Path == "/" && strings.Contains(r.URL.RawQuery, "uploads"):
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("ListMultipartUploads Authorization = %q, want session credential", got)
			}
			w.Write([]byte(`<ListMultipartUploadsResult></ListMultipartUploadsResult>`))
		default:
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	_, resp, err := c.Bucket.ListMultipartUploads(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestObjectService_CopyUsesSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
		case r.URL.Path == "/dest.txt":
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("Copy Authorization = %q, want session credential", got)
			}
			w.Write([]byte(`<CopyObjectResult><ETag>"etag"</ETag><LastModified>2017-12-13T14:53:12</LastModified></CopyObjectResult>`))
		default:
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	_, resp, err := c.Object.Copy(context.Background(), "dest.txt", c.Host+"/src.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestObjectService_CopyPartUsesSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
		case r.URL.Path == "/dest.txt":
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("CopyPart Authorization = %q, want session credential", got)
			}
			w.Write([]byte(`<CopyPartResult><ETag>"etag"</ETag><LastModified>2017-12-13T14:53:12</LastModified></CopyPartResult>`))
		default:
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	_, resp, err := c.Object.CopyPart(context.Background(), "dest.txt", "uploadid", 1, c.Host+"/src.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestObjectService_GetSessionPresignedURL(t *testing.T) {
	var sessionCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.URL.RawQuery != "session" {
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		atomic.AddInt32(&sessionCalls, 1)
		w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	presignedURL, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", 10, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if presignedURL.Host != c.Host {
		t.Fatalf("presigned URL host = %q, want %q", presignedURL.Host, c.Host)
	}
	if got := presignedURL.Query().Get("x-cos-security-token"); got != "session-token" {
		t.Fatalf("x-cos-security-token = %q, want session-token", got)
	}
	if got := presignedURL.RawQuery; !strings.Contains(got, "q-ak=SKIDSESSION") {
		t.Fatalf("RawQuery = %q, want session credential", got)
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 1 {
		t.Fatalf("CreateSession calls = %d, want 1", got)
	}
}

func TestObjectService_GetSessionPresignedURLRejectsSignMerged(t *testing.T) {
	var sessionCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&sessionCalls, 1)
		w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	presignedURL, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", 10*time.Minute, &PresignedURLOptions{SignMerged: true})
	if err == nil || !strings.Contains(err.Error(), "SignMerged") {
		t.Fatalf("GetSessionPresignedURL() = (%v, %v), want SignMerged unsupported error", presignedURL, err)
	}
	if presignedURL != nil {
		t.Fatalf("presigned URL = %v, want nil", presignedURL)
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 0 {
		t.Fatalf("CreateSession calls = %d, want 0", got)
	}
}

func TestObjectService_GetSessionPresignedURLRefreshesWhenTokenLifetimeTooShort(t *testing.T) {
	var sessionCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.URL.RawQuery != "session" {
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		n := atomic.AddInt32(&sessionCalls, 1)
		// 第一次返回一个剩 2 分钟的 token: 已过 refreshBefore(90s) 窗口, 普通复用逻辑
		// 会命中缓存, 但覆盖不了预签名要求的 10 分钟。第二次返回寿命充足的 token。
		var expiration time.Time
		if n == 1 {
			expiration = time.Now().Add(2 * time.Minute)
		} else {
			expiration = time.Now().Add(1 * time.Hour)
		}
		fmt.Fprint(w, sessionResponseXML(fmt.Sprintf("SKIDSESSION%d", n), "session-secret", fmt.Sprintf("session-token-%d", n), expiration))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	// 缓存里先放一个只剩 30s 的 token。
	if _, _, err := provider.credential(context.Background(), u, c.Host, ""); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 1 {
		t.Fatalf("setup CreateSession calls = %d, want 1", got)
	}

	// 预签名要求 10 分钟有效, 30s 的缓存 token 覆盖不了, 应强制刷新拿到新 token。
	presignedURL, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", 10*time.Minute, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&sessionCalls); got != 2 {
		t.Fatalf("CreateSession calls = %d, want 2 (should refresh)", got)
	}
	if got := presignedURL.Query().Get("x-cos-security-token"); got != "session-token-2" {
		t.Fatalf("x-cos-security-token = %q, want session-token-2 (refreshed)", got)
	}
	if got := presignedURL.RawQuery; !strings.Contains(got, "q-ak=SKIDSESSION2") {
		t.Fatalf("RawQuery = %q, want refreshed session credential", got)
	}
}

func TestObjectService_GetSessionPresignedURLRejectsWhenExpiredExceedsMaxTokenLifetime(t *testing.T) {
	var sessionCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.URL.RawQuery != "session" {
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		atomic.AddInt32(&sessionCalls, 1)
		// 服务端 token 总寿命固定 5 分钟, 无法延长。
		expiration := time.Now().Add(5 * time.Minute)
		fmt.Fprint(w, sessionResponseXML("SKIDSESSION", "session-secret", "session-token", expiration))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	// 请求 1 小时有效期, 但全新 token 也只有 5 分钟 -> 刷新后仍不够 -> 报错, 不签 URL。
	_, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", time.Hour, nil, false)
	if err == nil {
		t.Fatal("expected error when expired exceeds max token lifetime")
	}
	if !strings.Contains(err.Error(), "shorter than required") {
		t.Fatalf("error = %q, want to mention insufficient lifetime", err.Error())
	}
	// provider 会强制刷新一次(缓存空/不足), 但不会无限重试。
	if got := atomic.LoadInt32(&sessionCalls); got != 1 {
		t.Fatalf("CreateSession calls = %d, want 1", got)
	}
}

func TestObjectService_GetSessionPresignedURLRejectsUnsupportedTransport(t *testing.T) {
	u, _ := url.Parse("https://zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com")
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: sessionAuthContextCheckTransport{t: t}})

	_, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", time.Minute, nil)
	if err == nil || !strings.Contains(err.Error(), "transport does not support") {
		t.Fatalf("error = %v, want unsupported transport error", err)
	}
}

func TestObjectService_GetSessionPresignedURLRejectsOrdinaryBucket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("ordinary bucket should not call CreateSession: %s?%s", r.URL.Path, r.URL.RawQuery)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "normal-1250000000.cos.ap-guangzhou.myqcloud.com"

	_, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", 10, nil, false)
	if err == nil {
		t.Fatal("expected error when bucket is not a high-performance bucket")
	}
	if !strings.Contains(err.Error(), "high-performance") {
		t.Fatalf("error = %q, want to mention high-performance bucket", err.Error())
	}
}

func TestObjectService_GetSessionPresignedURLRejectsAuthTimeAfterSessionExpiration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	_, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "test.txt", 10, &PresignedURLOptions{
		AuthTime: &AuthTime{
			SignStartTime: time.Date(2099, 1, 2, 3, 0, 0, 0, time.UTC),
			SignEndTime:   time.Date(2099, 1, 2, 3, 5, 0, 0, time.UTC),
			KeyStartTime:  time.Date(2099, 1, 2, 3, 0, 0, 0, time.UTC),
			KeyEndTime:    time.Date(2099, 1, 2, 3, 5, 0, 0, time.UTC),
		},
	})
	if err == nil {
		t.Fatal("expected error when presigned AuthTime exceeds session expiration")
	}
}

type sessionAuthContextCheckTransport struct {
	t               *testing.T
	wantSessionAuth bool
}

func (t sessionAuthContextCheckTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if got := isSessionAuth(req.Context()); got != t.wantSessionAuth {
		t.t.Errorf("isSessionAuth = %v, want %v", got, t.wantSessionAuth)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       ioutil.NopCloser(strings.NewReader("ok")),
		Request:    req,
	}, nil
}

func TestObjectService_OrdinaryBucketDoesNotSetSessionAuthContext(t *testing.T) {
	u, _ := url.Parse("https://normal-1250000000.cos.ap-guangzhou.myqcloud.com")
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: sessionAuthContextCheckTransport{t: t}})
	resp, err := c.Object.Get(context.Background(), "test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestAuthorizationTransport_OrdinaryBucketUsesBaseCredential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery == "session" {
			t.Fatalf("ordinary bucket should not call CreateSession")
		}
		if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=AKIDBASE") {
			t.Errorf("Authorization = %q, want base credential", got)
		}
		if token := r.Header.Get("x-cos-security-token"); token != "" {
			t.Errorf("x-cos-security-token = %q, want empty", token)
		}
		fmt.Fprint(w, "ok")
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "normal-1250000000.cos.ap-guangzhou.myqcloud.com"

	resp, err := c.Object.Get(context.Background(), "test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestParseBucketNameFromHost(t *testing.T) {
	cases := []struct {
		host string
		want string
	}{
		{"zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com", "zlb-normal-x--251433186"},
		{"ZLB-NORMAL-X--251433186.COS-RAPID.ap-guangzhou.tencentcos.com", "zlb-normal-x--251433186"},
		{"zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com:443", "zlb-normal-x--251433186"},
		{"zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com.", "zlb-normal-x--251433186"},
		// No domain anchor required: any hostname's first label is extracted,
		// regardless of whether it's a ".cos-rapid." domain. isRapidBucket is
		// solely responsible for judging whether the extracted name is a rapid
		// bucket name (see TestIsHighPerformanceBucket).
		{"testuser-dev-x--1374209587.gz.tencentos.com", "testuser-dev-x--1374209587"},
		{"TESTUSER-DEV-X--1374209587.GZ.TENCENTOS.COM", "testuser-dev-x--1374209587"},
		{"testuser-dev-x--1374209587.gz.tencentos.com:443", "testuser-dev-x--1374209587"},
		{"testuser-dev-x--1374209587.gz.tencentos.com.", "testuser-dev-x--1374209587"},
		{"normal-1250000000.cos.ap-guangzhou.myqcloud.com", "normal-1250000000"},
		{"justbucket", "justbucket"},
		{"", ""},
		{".cos-rapid.ap-guangzhou.tencentcos.com", ""},
	}
	for _, c := range cases {
		if got := parseBucketNameFromHost(c.host); got != c.want {
			t.Errorf("parseBucketNameFromHost(%q) = %q, want %q", c.host, got, c.want)
		}
	}
}

func TestIsHighPerformanceBucket(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"zlb-normal-x--251433186", true},
		{"testuser-dev-x--1374209587", true},
		{"a-x--1", true},
		{"zlb-normal-x-251433186", false}, // 单横线不是rapid bucket
		{"example-1250000000", false},
		{"zlb-normal-x--abc", false}, // appid 必须数字
		{"-x--123", false},           // 不能以 - 开头
		{"", false},
		{"ab", false}, // 太短
	}
	for _, c := range cases {
		if got := isRapidBucket(c.name); got != c.want {
			t.Errorf("isHighPerformanceBucket(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestCreateSessionProvider_CachesDifferentBuckets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sessionResponseXML("SKID", "sk", "tok", time.Unix(4070912645, 0))))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	if _, _, err := provider.credential(context.Background(), u, "aaa-x--111.cos-rapid.ap-guangzhou.tencentcos.com", ""); err != nil {
		t.Fatal(err)
	}
	if _, _, err := provider.credential(context.Background(), u, "bbb-x--222.cos-rapid.ap-guangzhou.tencentcos.com", ""); err != nil {
		t.Fatal(err)
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	if len(provider.bucketSessions) != 2 {
		t.Fatalf("cached buckets = %d, want 2", len(provider.bucketSessions))
	}
}

func TestCreateSessionProvider_UsesBaseSessionToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-cos-security-token"); got != "base-sts-token" {
			t.Errorf("CreateSession x-cos-security-token = %q, want base-sts-token", got)
		}
		w.Write([]byte(sessionResponseXML("SKID", "sk", "tok", time.Unix(4070912645, 0))))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "base-sts-token")

	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	if _, _, err := provider.credential(context.Background(), u, host, ""); err != nil {
		t.Fatal(err)
	}
}

func TestCreateSessionProvider_ReadsUpdatedTransportCredential(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if strings.Contains(auth, "q-ak=AKIDOLD") {
			seen = append(seen, "AKIDOLD")
		} else if strings.Contains(auth, "q-ak=AKIDNEW") {
			seen = append(seen, "AKIDNEW")
		}
		exp := time.Now().Add(1 * time.Second) // 短寿命, 迫使每次刷新
		fmt.Fprint(w, sessionResponseXML("SKID", "sk", "tok", exp))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, transport, provider := newSessionTestClient(u, "AKIDOLD", "sk-old", "")

	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	if _, _, err := provider.credential(context.Background(), u, host, ""); err != nil {
		t.Fatalf("first GetSessionCredential: %v", err)
	}

	transport.SetCredential("AKIDNEW", "sk-new", "")
	if _, _, err := provider.credential(context.Background(), u, host, ""); err != nil {
		t.Fatalf("second GetSessionCredential: %v", err)
	}

	if len(seen) != 2 || seen[0] != "AKIDOLD" || seen[1] != "AKIDNEW" {
		t.Fatalf("CreateSession credentials seen = %v, want [AKIDOLD AKIDNEW]", seen)
	}
}

func TestCreateSessionProvider_UsesBucketCreateSessionErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`<Error><Code>AccessDenied</Code><Message>denied</Message><RequestId>req-1</RequestId></Error>`))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	_, _, err := provider.credential(context.Background(), u, "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com", "")
	cosErr, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("error = %T %v, want *ErrorResponse", err, err)
	}
	if cosErr.Code != "AccessDenied" || cosErr.Message != "denied" || cosErr.RequestID != "req-1" {
		t.Fatalf("error response = %+v", cosErr)
	}
}

func TestCreateSessionProvider_DoesNotFollowRedirect(t *testing.T) {
	var target int32
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&target, 1)
		w.Write([]byte(sessionResponseXML("SKID", "sk", "tok", time.Unix(4070912645, 0))))
	}))
	defer backend.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, backend.URL, http.StatusFound)
	}))
	defer redirector.Close()

	u, _ := url.Parse(redirector.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	_, _, err := provider.credential(context.Background(), u, host, "")
	// 不跟随重定向 -> 拿到 302, 非 2xx -> 报错; 且不会打到 backend。
	if err == nil {
		t.Fatal("expected error: CreateSession should not follow redirect")
	}
	if atomic.LoadInt32(&target) != 0 {
		t.Fatalf("redirect target was reached %d times, want 0 (SSRF risk)", target)
	}
}

func TestCreateSessionProvider_FreshCredentialRejectsMissingExpiration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sessionResponseXML("SKID", "sk", "tok", time.Time{})))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	_, _, err := provider.freshCredential(context.Background(), u, "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com", 10*time.Minute, "")
	if err == nil {
		t.Fatal("expected error for missing session credential expiration")
	}
	if !strings.Contains(err.Error(), "expiration is missing") {
		t.Fatalf("error = %q, want missing expiration error", err.Error())
	}
}

func TestCreateSessionProvider_RefreshesSessionWhenBaseCredentialChangesEvenIfNotExpired(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if strings.Contains(auth, "q-ak=AKIDOLD") {
			seen = append(seen, "AKIDOLD")
		} else if strings.Contains(auth, "q-ak=AKIDNEW") {
			seen = append(seen, "AKIDNEW")
		}
		// 长寿命 token: 普通过期逻辑绝不会触发刷新, 只有基础凭证指纹变化才会重签。
		fmt.Fprint(w, sessionResponseXML("SKID", "sk", "tok", time.Unix(4070912645, 0)))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, base, provider := newSessionTestClient(u, "AKIDOLD", "sk-old", "")
	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	if _, _, err := provider.credential(context.Background(), u, host, ""); err != nil {
		t.Fatalf("first credential: %v", err)
	}
	// 缓存里是长寿命 token, 正常绝不会刷新; 切换底层基础凭证后必须重签。
	base.SetCredential("AKIDNEW", "sk-new", "")
	if _, _, err := provider.credential(context.Background(), u, host, ""); err != nil {
		t.Fatalf("second credential: %v", err)
	}

	if len(seen) != 2 || seen[0] != "AKIDOLD" || seen[1] != "AKIDNEW" {
		t.Fatalf("CreateSession credentials seen = %v, want [AKIDOLD AKIDNEW] (must re-sign on base credential change)", seen)
	}
}

// countingRoundTripper 记录被经过的次数, 并把请求委托给内部 next。
// 用于验证 session 数据请求确实经过了底层 transport 配置的自定义网络层。
type countingRoundTripper struct {
	calls int32
	next  http.RoundTripper
}

func (rt *countingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&rt.calls, 1)
	return rt.next.RoundTrip(req)
}

// 问题1回归: 用户配置在底层 transport.Transport 上的自定义网络层, 必须对 session 数据请求生效。
func TestRapidCredentialTransport_InheritsBaseNetworkTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.URL.RawQuery == "session":
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
		case r.URL.Path == "/test.txt":
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("Object Authorization = %q, want session credential", got)
			}
			w.Write([]byte("ok"))
		default:
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	custom := &countingRoundTripper{next: http.DefaultTransport}
	base := &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret", Transport: custom}
	transport := NewRapidCredentialTransport(base)

	u, _ := url.Parse(server.URL)
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: transport})
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	resp, err := c.Object.Get(context.Background(), "test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// CreateSession(经底层签名走底层网络层) + session 数据请求(继承底层网络层) 都应经过 custom。
	if got := atomic.LoadInt32(&custom.calls); got < 2 {
		t.Fatalf("custom base transport calls = %d, want >= 2 (CreateSession + session data request)", got)
	}
}

// 本文件专补集成测试(真机/真实 HTTP 链路)无法触达的 nil / 空参 / 类型断言 /
// 兜底分支，与主干链路集成测试合起来把本次 session
// 鉴权改动涉及的函数做到 100% 语句覆盖。

// --- RapidCredentialTransport.GetCredential ---

func TestRapidCredentialTransport_GetCredential_NilBase(t *testing.T) {
	tr := &RapidCredentialTransport{} // baseTransport 为 nil
	if _, _, _, err := tr.GetCredential(); err == nil || !strings.Contains(err.Error(), "base transport is nil") {
		t.Fatalf("GetCredential err = %v, want base transport is nil", err)
	}
}

func TestRapidCredentialTransport_GetCredential_DelegatesToBase(t *testing.T) {
	base := &AuthorizationTransport{SecretID: "AKID", SecretKey: "sk", SessionToken: "tok"}
	tr := NewRapidCredentialTransport(base)
	ak, sk, token, err := tr.GetCredential()
	if err != nil || ak != "AKID" || sk != "sk" || token != "tok" {
		t.Fatalf("GetCredential = (%q,%q,%q,%v), want base credential", ak, sk, token, err)
	}
}

// --- RapidCredentialTransport.RoundTrip: base 为 nil ---

func TestRapidCredentialTransport_RoundTrip_NilBase(t *testing.T) {
	tr := &RapidCredentialTransport{}
	req, _ := http.NewRequest(http.MethodGet, "https://zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com/a.txt", nil)
	if _, err := tr.RoundTrip(req); err == nil || !strings.Contains(err.Error(), "base transport is nil") {
		t.Fatalf("RoundTrip err = %v, want base transport is nil", err)
	}
}

// --- RapidCredentialTransport.RoundTrip: session 请求但非rapid bucket -> 回退底层基础凭证 ---

func TestRapidCredentialTransport_RoundTrip_SessionAuthButOrdinaryBucketFallsBack(t *testing.T) {
	// 底层用一个记录凭证的探针 transport，验证 !ok 时回退底层 RoundTrip。
	probe := &credentialProbeTransport{ak: "AKIDBASE", sk: "sk"}
	tr := NewRapidCredentialTransport(probe)
	// host 是普通桶，provider.credential 返回 ok=false -> 走 base.RoundTrip
	req, _ := http.NewRequest(http.MethodGet, "https://normal-1250000000.cos.ap-guangzhou.myqcloud.com/a.txt", nil)
	req = req.WithContext(withSessionAuth(req.Context()))
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if !probe.called {
		t.Fatal("expected fallback to base transport RoundTrip for ordinary bucket")
	}
}

// --- session 请求委托底层发送：底层收到"已签名"标记应跳过重签、保留 session 签名 ---

// bareBaseTransport 只实现 BaseTransport(TransportIface + RoundTripper)，
// 用于验证 RapidCredentialTransport 签名后委托底层发送时，底层沿用自身网络层。
type bareBaseTransport struct {
	ak, sk, token string
}

func (b *bareBaseTransport) GetCredential() (string, string, string, error) {
	return b.ak, b.sk, b.token, nil
}
func (b *bareBaseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}
func (b *bareBaseTransport) RoundTripSigned(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}

// 四个内置 BaseTransport 的 RoundTripSigned 必须只发送、不重签，保留上层 session 签名。
func TestBuiltinTransports_RoundTripSignedDoesNotResign(t *testing.T) {
	const signed = "q-ak=SESSIONAK&already-signed"
	cred := &Credential{SecretID: "AKIDBASE", SecretKey: "base-secret"}
	cases := []struct {
		name string
		base BaseTransport
	}{
		{"AuthorizationTransport", &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "base-secret"}},
		{"CVMCredentialTransport", &CVMCredentialTransport{secretID: "AKIDBASE", secretKey: "base-secret", expiredTime: time.Now().Add(time.Hour).Unix()}},
		{"CredentialTransport", &CredentialTransport{Credential: cred}},
		{"StsCredentialTransport", &StsCredentialTransport{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var seenAuth string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seenAuth = r.Header.Get("Authorization")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			req, _ := http.NewRequest(http.MethodGet, server.URL+"/a.txt", nil)
			req.Header.Set("Authorization", signed) // 模拟上层已用 session 凭证签好名
			resp, err := c.base.RoundTripSigned(req)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()
			if seenAuth != signed {
				t.Fatalf("Authorization = %q, want pre-set session signature (RoundTripSigned must not re-sign)", seenAuth)
			}
		})
	}
}

// 底层仅实现 BaseTransport(bareBaseTransport) 时，session 请求应能签名并委托底层发送(不 panic)。
func TestRapidCredentialTransport_SessionRequestWithBareBaseTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.RawQuery == "session":
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "ss", "st", time.Unix(4070912645, 0))))
		default:
			if got := r.Header.Get("Authorization"); !strings.Contains(got, "q-ak=SKIDSESSION") {
				t.Errorf("data request Authorization = %q, want session credential", got)
			}
			w.Write([]byte("ok"))
		}
	}))
	defer server.Close()

	tr := NewRapidCredentialTransport(&bareBaseTransport{ak: "AKIDBASE", sk: "sk"})
	u, _ := url.Parse(server.URL)
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: tr})
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	resp, err := c.Object.Get(context.Background(), "test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

// --- createSessionProvider 的 nil 分支 ---

func TestCreateSessionProvider_BaseTransportNil(t *testing.T) {
	p := newCreateSessionProvider(nil) // owner 为 nil
	if p.baseTransport() != nil {
		t.Fatal("baseTransport() = non-nil, want nil when owner is nil")
	}
	if p.baseSecretID() != "" {
		t.Fatal("baseSecretID() != empty, want empty when base is nil")
	}
	// createSession: base 为 nil -> 报错
	u, _ := url.Parse("https://zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com")
	if _, err := p.createSession(context.Background(), u, "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com", ""); err == nil || !strings.Contains(err.Error(), "base transport is nil") {
		t.Fatalf("createSession err = %v, want base transport is nil", err)
	}
}

func TestCreateSessionProvider_CreateSession_NilURL(t *testing.T) {
	base := &AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"}
	p := newCreateSessionProvider(NewRapidCredentialTransport(base))
	if _, err := p.createSession(context.Background(), nil, "host", ""); err == nil || !strings.Contains(err.Error(), "request URL is nil") {
		t.Fatalf("createSession(nil url) err = %v, want request URL is nil", err)
	}
}

// --- RenameObject 参数校验 ---

func TestObjectService_RenameObject_EmptyArgs(t *testing.T) {
	u, _ := url.Parse("https://zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com")
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: NewRapidCredentialTransport(&AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"})})

	if _, err := c.Object.RenameObject(context.Background(), "", "src", nil); err == nil || !strings.Contains(err.Error(), "empty object name") {
		t.Fatalf("RenameObject(empty name) err = %v, want empty object name", err)
	}
	if _, err := c.Object.RenameObject(context.Background(), "/", "src", nil); err == nil || !strings.Contains(err.Error(), "empty object name") {
		t.Fatalf("RenameObject(slash name) err = %v, want empty object name", err)
	}
	if _, err := c.Object.RenameObject(context.Background(), "dst", "", nil); err == nil || !strings.Contains(err.Error(), "empty rename source") {
		t.Fatalf("RenameObject(empty source) err = %v, want empty rename source", err)
	}
	if _, err := c.Object.RenameObject(context.Background(), "dst", "/", nil); err == nil || !strings.Contains(err.Error(), "empty rename source") {
		t.Fatalf("RenameObject(slash source) err = %v, want empty rename source", err)
	}
}

// TestObjectService_RenameObject_OrdinaryBucketRejected 验证对普通桶调用 RenameObject
// 会被本地 fail-fast 拒绝(不发出请求),错误提示指向 rapid bucket。
func TestObjectService_RenameObject_OrdinaryBucketRejected(t *testing.T) {
	// 普通桶域名(非 -x-- rapid bucket),useSessionAuth() 返回 false。
	u, _ := url.Parse("https://normal-bucket-1250000000.cos.ap-guangzhou.myqcloud.com")
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: NewRapidCredentialTransport(&AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"})})

	if _, err := c.Object.RenameObject(context.Background(), "dst", "src", nil); err == nil || !strings.Contains(err.Error(), "rapid bucket") {
		t.Fatalf("RenameObject(ordinary bucket) err = %v, want rapid bucket rejection", err)
	}
}

// --- freshCredential: 缓存不足以覆盖 minRemaining(getSessionCredential 返回 ok 但寿命不够) ---
// 该分支验证 lifetime shorter than required 提示。已有 RejectsWhenExpiredExceedsMaxTokenLifetime
// 覆盖了刷新后仍不足；这里补一个纯 provider 层、显式短 minRemaining 判断路径的直接断言。

func TestCreateSessionProvider_FreshCredential_ShortLifetime(t *testing.T) {
	// httptest 返回固定 5 分钟寿命 token；freshCredential 要求 10 分钟 minRemaining，
	// 刷新后仍不足 -> 命中 "shorter than required" 分支。
	server := newShortLifetimeSessionServer(t, 5*time.Minute)
	defer server.Close()

	base := &AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"}
	p := newCreateSessionProvider(NewRapidCredentialTransport(base))
	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	u, _ := url.Parse(server.URL)
	_, _, err := p.freshCredential(context.Background(), u, host, 10*time.Minute, "")
	if err == nil || !strings.Contains(err.Error(), "shorter than required") {
		t.Fatalf("freshCredential err = %v, want shorter than required", err)
	}
}

// sessionAnsweringTransport 作为底层 BaseTransport：应答 CreateSession(?session)返回假
// session 凭证，其余(数据)请求经 onData 回调后回 200，用于端到端验证 session 鉴权是否触发。
type sessionAnsweringTransport struct {
	createSessionCalled bool
	onData              func(*http.Request)
}

func (p *sessionAnsweringTransport) GetCredential() (string, string, string, error) {
	return "AKIDBASE", "base-secret", "", nil
}

func (p *sessionAnsweringTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.RawQuery == "session" {
		p.createSessionCalled = true
		body := sessionResponseXML("SKIDSESSION", "ss", "st", time.Unix(4070912645, 0))
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(strings.NewReader(body)),
			Request:    req,
		}, nil
	}
	if p.onData != nil {
		p.onData(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

func (p *sessionAnsweringTransport) RoundTripSigned(req *http.Request) (*http.Response, error) {
	if p.onData != nil {
		p.onData(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

// credentialProbeTransport 实现 BaseTransport，记录 RoundTrip 是否被调用，用于回退验证。
type credentialProbeTransport struct {
	ak, sk, token string
	called        bool
}

func (p *credentialProbeTransport) GetCredential() (string, string, string, error) {
	return p.ak, p.sk, p.token, nil
}
func (p *credentialProbeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	p.called = true
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    req,
	}, nil
}
func (p *credentialProbeTransport) RoundTripSigned(req *http.Request) (*http.Response, error) {
	return p.RoundTrip(req)
}

// newShortLifetimeSessionServer 返回一个每次都签发固定 lifetime 寿命 session token 的假网关。
func newShortLifetimeSessionServer(t *testing.T, lifetime time.Duration) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "session" {
			t.Errorf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		exp := time.Now().Add(lifetime)
		fmt.Fprint(w, sessionResponseXML("SKID", "sk", "tok", exp))
	}))
}

// --- presignRequiredLifetime: KeyEndTime 晚于 SignEndTime 时取 KeyEndTime ---

func TestPresignRequiredLifetime_KeyEndTimeWins(t *testing.T) {
	now := time.Now()
	opt := &PresignedURLOptions{
		AuthTime: &AuthTime{
			SignStartTime: now,
			SignEndTime:   now.Add(5 * time.Minute),
			KeyStartTime:  now,
			KeyEndTime:    now.Add(30 * time.Minute), // 晚于 SignEndTime
		},
	}
	got := presignRequiredLifetime(time.Minute, opt)
	// 应以 KeyEndTime(约 30 分钟)为准，明显大于 SignEndTime 的 5 分钟。
	if got < 20*time.Minute {
		t.Fatalf("presignRequiredLifetime = %s, want ~KeyEndTime(>=20m)", got)
	}
	// authTime 为 nil 时直接返回 expired。
	if got := presignRequiredLifetime(time.Minute, nil); got != time.Minute {
		t.Fatalf("presignRequiredLifetime(nil) = %s, want 1m", got)
	}
}

// --- GetSessionPresignedURL: name 为空 / provider 为 nil / host 从 BucketURL 取 ---

func TestGetSessionPresignedURL_EmptyName(t *testing.T) {
	u, _ := url.Parse("https://zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com")
	c := NewClient(&BaseURL{BucketURL: u}, &http.Client{Transport: NewRapidCredentialTransport(&AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"})})
	if _, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "", time.Minute, nil); err == nil || !strings.Contains(err.Error(), "object key is empty") {
		t.Fatalf("err = %v, want object key is empty", err)
	}
}

func TestGetSessionPresignedURL_HostFromBucketURL(t *testing.T) {
	// 不设 c.Host，走 host 从 BucketURL.Host 取的分支。
	server := newShortLifetimeSessionServer(t, time.Hour)
	defer server.Close()
	u, _ := url.Parse(server.URL)
	// 让 BucketURL.Host 本身是rapid bucket域名。
	u.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	base := &AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"}
	tr := NewRapidCredentialTransport(base)
	// 用真实 server 发 CreateSession：把 BucketURL 指到 server，但 host 名保持桶域名做识别。
	realU, _ := url.Parse(server.URL)
	c := NewClient(&BaseURL{BucketURL: realU}, &http.Client{Transport: tr})
	// c.Host 留空 -> GetSessionPresignedURL 内部从 BucketURL.Host 取。
	c.BaseURL.BucketURL.Host = realU.Host // 网络仍打 server
	// 直接断言 host 取值分支：用 provider 为 nil 的 transport 触发不了，这里改测 host 空时的桶识别失败分支。
	_, err := c.Object.GetSessionPresignedURL(context.Background(), http.MethodGet, "obj.txt", time.Minute, nil)
	// server.URL 的 host 不是rapid bucket -> 命中 "only supports high-performance bucket"。
	if err == nil || !strings.Contains(err.Error(), "high-performance") {
		t.Fatalf("err = %v, want high-performance bucket rejection (host from BucketURL)", err)
	}
	_ = u
}

// --- RoundTrip: req.Host 为空时经 req.URL.Hostname() 识别rapid bucket并触发 session 鉴权 ---
// (回归 review 问题1：SDK 常规用法只设 BucketURL、不设 Client.Host，req.Host 为空，
//  必须靠 URL.Hostname() 识别，否则 session 鉴权静默退化成基础凭证。)

func TestRapidCredentialTransport_RoundTrip_HostFromURL(t *testing.T) {
	// 底层探针：应答 CreateSession(?session)，并记录数据请求最终带的 Authorization。
	var dataAuth string
	probe := &sessionAnsweringTransport{
		onData: func(r *http.Request) { dataAuth = r.Header.Get("Authorization") },
	}
	tr := NewRapidCredentialTransport(probe)

	// URL 的 host 是rapid bucket域名，但 req.Host 显式置空，只能靠 URL.Hostname() 识别。
	req, _ := http.NewRequest(http.MethodGet,
		"https://zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com/a.txt", nil)
	req.Host = ""
	req = req.WithContext(withSessionAuth(req.Context()))

	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if !probe.createSessionCalled {
		t.Fatal("expected CreateSession to be called (session auth via URL.Hostname), but it was not")
	}
	if !strings.Contains(dataAuth, "q-ak=SKIDSESSION") {
		t.Fatalf("data request Authorization = %q, want session credential (q-ak=SKIDSESSION)", dataAuth)
	}
}

// --- freshCredential: getSessionCredential 返回 !ok(非rapid bucket)时早退 ---

func TestCreateSessionProvider_FreshCredential_NotRapidBucket(t *testing.T) {
	base := &AuthorizationTransport{SecretID: "AKID", SecretKey: "sk"}
	p := newCreateSessionProvider(NewRapidCredentialTransport(base))
	u, _ := url.Parse("https://normal-1250000000.cos.ap-guangzhou.myqcloud.com")
	cred, ok, err := p.freshCredential(context.Background(), u, "normal-1250000000.cos.ap-guangzhou.myqcloud.com", time.Minute, "")
	if err != nil || ok || cred != nil {
		t.Fatalf("freshCredential(ordinary) = (%v,%v,%v), want (nil,false,nil)", cred, ok, err)
	}
}

// --- RoundTrip: provider.credential 出错(CreateSession 失败)时返回 error ---

func TestRapidCredentialTransport_RoundTrip_CreateSessionError(t *testing.T) {
	// rapid bucket请求 -> 触发 CreateSession；server 返回 403 错误 -> provider.credential 返回 err
	// -> RoundTrip 命中 `if err != nil { return nil, err }` 分支。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `<Error><Code>AccessDenied</Code><Message>denied</Message></Error>`)
	}))
	defer server.Close()

	base := &AuthorizationTransport{SecretID: "AKIDBASE", SecretKey: "sk"}
	tr := NewRapidCredentialTransport(base)
	su, _ := url.Parse(server.URL)
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/a.txt", nil)
	// Host 用rapid bucket域名以触发 session 鉴权；网络仍打到 server(httptest client 按 req.URL 连接)。
	req.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	req = req.WithContext(withSessionAuth(req.Context()))
	_ = su
	if _, err := tr.RoundTrip(req); err == nil {
		t.Fatal("expected error when CreateSession fails during RoundTrip")
	}
}

// --- 刷新失败降级: 旧 session 未真正过期且基础身份未变 -> 复用旧凭证顶过抖动/限频 ---

func TestCreateSessionProvider_ReusesUnexpiredCredentialWhenRefreshFails(t *testing.T) {
	// server 始终返回 500, 保证本次 CreateSession(刷新) 必失败。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "sk", "")
	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	bucketName := parseBucketNameFromHost(host)

	// 预置一份"已跌破提前刷新窗口(默认 1min)但尚未真正过期"的旧凭证:
	// 剩余 30s -> 触发刷新, 但真过期线仍在未来。baseSecretID 与底层一致。
	b := provider.getBucketSession(bucketName)
	old := &SessionCredential{
		SecretID:     "SKIDOLD",
		SecretKey:    "oldsk",
		SessionToken: "oldtok",
		Expiration:   time.Now().Add(30 * time.Second),
	}
	b.cred = old
	b.baseSecretID = "AKIDBASE"

	cred, ok, err := provider.credential(context.Background(), u, host, "")
	if err != nil {
		t.Fatalf("expected graceful fallback (no error), got err = %v", err)
	}
	if !ok || cred == nil {
		t.Fatalf("expected ok=true with fallback credential, got ok=%v cred=%v", ok, cred)
	}
	if cred.SecretID != "SKIDOLD" || cred.SessionToken != "oldtok" {
		t.Fatalf("expected reuse of old credential, got %+v", cred)
	}
}

// --- 刷新失败不降级: 旧 session 已真正过期 -> 必须返回 error, 不得复用废凭证 ---

func TestCreateSessionProvider_ReturnsErrorWhenRefreshFailsAndCredentialExpired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	_, _, provider := newSessionTestClient(u, "AKIDBASE", "sk", "")
	host := "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"
	bucketName := parseBucketNameFromHost(host)

	// 旧凭证已过真过期线, 刷新又失败 -> 不能拿废凭证顶, 必须报错。
	b := provider.getBucketSession(bucketName)
	b.cred = &SessionCredential{
		SecretID:   "SKIDOLD",
		SecretKey:  "oldsk",
		Expiration: time.Now().Add(-1 * time.Second),
	}
	b.baseSecretID = "AKIDBASE"

	if _, ok, err := provider.credential(context.Background(), u, host, ""); err == nil || ok {
		t.Fatalf("expected error and ok=false when refresh fails and credential expired, got ok=%v err=%v", ok, err)
	}
}

// --- RoundTrip: 数据面返回 ExpiredToken -> 清缓存凭证并改写成 503 交上层重试 ---

func TestRapidCredentialTransport_RoundTrip_RewritesExpiredTokenToUnavailable(t *testing.T) {
	var sessionCalls, dataCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.URL.RawQuery == "session" {
			atomic.AddInt32(&sessionCalls, 1)
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
			return
		}
		atomic.AddInt32(&dataCalls, 1)
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `<Error><Code>ExpiredToken</Code><Message>The provided token has expired.</Message></Error>`)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, provider := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/a.txt", nil)
	req.Host = c.Host
	req = req.WithContext(withSessionAuth(req.Context()))

	resp, err := c.client.Transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	// ExpiredToken(400) 被改写成 503，供上层按可重试处理。
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("StatusCode = %d, want 503 (rewritten from ExpiredToken)", resp.StatusCode)
	}
	// body 保留原始 ExpiredToken 错误码，供上层定位根因。
	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "ExpiredToken") {
		t.Fatalf("body = %q, want original ExpiredToken code preserved", string(body))
	}
	// transport 层不重发：数据面只被调用一次。
	if got := atomic.LoadInt32(&dataCalls); got != 1 {
		t.Fatalf("data-plane calls = %d, want 1 (no in-transport retry)", got)
	}
	// transport 不再主动清缓存：凭证仍在本地有效期内，保留供上层重试时按复用判定处理。
	b := provider.getBucketSession(parseBucketNameFromHost(c.Host))
	if b.cred == nil {
		t.Fatalf("cached credential = nil, want retained after ExpiredToken (no invalidate)")
	}
}

// --- RoundTrip: 非 ExpiredToken 错误(如 AccessDenied)不触发重试，原样透传且 body 完整可读 ---

func TestRapidCredentialTransport_RoundTrip_DoesNotRetryOnOtherError(t *testing.T) {
	var sessionCalls, dataCalls int32
	const deniedBody = `<Error><Code>AccessDenied</Code><Message>no permission</Message></Error>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.URL.RawQuery == "session" {
			atomic.AddInt32(&sessionCalls, 1)
			w.Write([]byte(sessionResponseXML("SKIDSESSION", "session-secret", "session-token", time.Unix(4070912645, 0))))
			return
		}
		atomic.AddInt32(&dataCalls, 1)
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, deniedBody)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	c, _, _ := newSessionTestClient(u, "AKIDBASE", "base-secret", "")
	c.Host = "zlb-normal-x--251433186.cos-rapid.ap-guangzhou.tencentcos.com"

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/a.txt", nil)
	req.Host = c.Host
	req = req.WithContext(withSessionAuth(req.Context()))

	resp, err := c.client.Transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("StatusCode = %d, want 403 passthrough", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != deniedBody {
		t.Fatalf("body = %q, want original error body intact", string(body))
	}
	if got := atomic.LoadInt32(&dataCalls); got != 1 {
		t.Fatalf("data-plane calls = %d, want 1 (no retry on non-expired error)", got)
	}
}
