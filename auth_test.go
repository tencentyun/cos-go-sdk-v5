package cos

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewAuthorization(t *testing.T) {
	expectAuthorization := `q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=host;x-cos-content-sha1;x-cos-stroage-class&q-url-param-list=&q-signature=ce4ac0ecbcdb30538b3fee0a97cc6389694ce53a`
	secretID := "QmFzZTY0IGlzIGEgZ2VuZXJp"
	secretKey := "AKIDZfbOA78asKUYBcXFrJD0a1ICvR98JM"
	host := "testbucket-125000000.cos.ap-guangzhou.myqcloud.com"
	uri := "http://testbucket-125000000.cos.ap-guangzhou.myqcloud.com/testfile2"
	startTime := time.Unix(int64(1480932292), 0)
	endTime := time.Unix(int64(1481012292), 0)

	req, _ := http.NewRequest("PUT", uri, nil)
	req.Header.Add("Host", host)
	req.Header.Add("x-cos-content-sha1", "db8ac1c259eb89d4a131b253bacfca5f319d54f2")
	req.Header.Add("x-cos-stroage-class", "nearline")

	authTime := &AuthTime{
		SignStartTime: startTime,
		SignEndTime:   endTime,
		KeyStartTime:  startTime,
		KeyEndTime:    endTime,
	}
	auth := newAuthorization(secretID, secretKey, req, authTime, true)

	if auth != expectAuthorization {
		t.Errorf("NewAuthorization returned \n%#v, want \n%#v", auth, expectAuthorization)
	}
}

func TestAuthorizationTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("AuthorizationTransport didn't add Authorization header")
		}
	})

	auth := &AuthorizationTransport{
		SecretID:  "test",
		SecretKey: "test",
	}
	auth.SetCredential("ak", "sk", "token")
	client.client.Transport = auth
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	client.doAPI(context.Background(), req, nil, true)
	client.GetCredential()
}

func TestCVMCredentialTransport(t *testing.T) {
	setup()
	defer teardown()
	uri := client.BaseURL.BucketURL.String()
	ak := "test_ak"
	sk := "test_sk"
	token := "test_token"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-cos-security-token") != token {
			t.Errorf("CVMCredentialTransport x-cos-security-token error, want:%v, return:%v\n", token, r.Header.Get("x-cos-security-token"))
		}
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("CVMCredentialTransport didn't add Authorization header")
		}
		field := strings.Split(auth, "&")
		if len(field) != 7 {
			t.Errorf("CVMCredentialTransport Authorization header format error: %v\n", auth)
		}
		st_et := strings.Split(strings.Split(field[2], "=")[1], ";")
		st, _ := strconv.ParseInt(st_et[0], 10, 64)
		et, _ := strconv.ParseInt(st_et[1], 10, 64)
		authTime := &AuthTime{
			SignStartTime: time.Unix(st, 0),
			SignEndTime:   time.Unix(et, 0),
			KeyStartTime:  time.Unix(st, 0),
			KeyEndTime:    time.Unix(et, 0),
		}
		host := strings.TrimLeft(uri, "http://")
		req, _ := http.NewRequest("GET", uri, nil)
		req.Header.Add("Host", host)
		expect := newAuthorization(ak, sk, req, authTime, true)
		if expect != auth {
			t.Errorf("CVMCredentialTransport Authorization error, want:%v, return:%v\n", expect, auth)
		}
	})

	// CVM http server
	cvm_mux := http.NewServeMux()
	cvm_server := httptest.NewServer(cvm_mux)
	defer cvm_server.Close()
	// 将默认 CVM Host 修改成测试IP:PORT
	defaultCVMMetaHost = strings.TrimLeft(cvm_server.URL, "http://")

	cvm_mux.HandleFunc("/"+defaultCVMCredURI, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "cvm_read_cos_only")
	})
	cvm_mux.HandleFunc("/"+defaultCVMCredURI+"/cvm_read_cos_only", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(`{
            "TmpSecretId": "%s",
            "TmpSecretKey": "%s",
            "ExpiredTime": %v,
            "Expiration": "now",
            "Token": "%s",
            "Code": "Success"
        }`, ak, sk, time.Now().Unix()+3600, token))
	})

	client.client.Transport = &CVMCredentialTransport{}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	client.doAPI(context.Background(), req, nil, true)

	req, _ = http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	client.doAPI(context.Background(), req, nil, true)
	client.GetCredential()
}

func TestDNSScatterTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("AuthorizationTransport didn't add Authorization header")
		}
	})

	client.client.Transport = &AuthorizationTransport{
		SecretID:  "test",
		SecretKey: "test",
		Transport: DNSScatterTransport,
	}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	client.doAPI(context.Background(), req, nil, true)
	client.GetCredential()
}

func TestCredentialTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("AuthorizationTransport didn't add Authorization header")
		}
	})

	client.client.Transport = &CredentialTransport{
		Credential: NewTokenCredential("test", "test", ""),
	}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	client.doAPI(context.Background(), req, nil, true)
	client.GetCredential()
}
