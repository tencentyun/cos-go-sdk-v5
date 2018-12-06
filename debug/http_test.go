package debug

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

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
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestDebugRequestTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Test-Response", "2333")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("test response body"))
	})

	w := bytes.NewBufferString("")
	client := http.Client{}

	client.Transport = &DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
		Writer:         w,
	}

	body := bytes.NewReader([]byte("test_request body"))
	req, _ := http.NewRequest("GET", server.URL, body)
	req.Header.Add("X-Test-Debug", "123")
	client.Do(req)

	b := make([]byte, 800)
	w.Read(b)
	info := string(b)
	if !strings.Contains(info, "GET / HTTP/1.1\r\n") ||
		!strings.Contains(info, "X-Test-Debug: 123\r\n") {
		t.Errorf("DebugRequestTransport debug info %#v don't contains request header", info)
	}
	if !strings.Contains(info, "\r\n\r\ntest_request body") {
		t.Errorf("DebugRequestTransport debug info  %#v don't contains request body", info)
	}

	if !strings.Contains(info, "HTTP/1.1 502 Bad Gateway\r\n") ||
		!strings.Contains(info, "X-Test-Response: 2333\r\n") {
		t.Errorf("DebugRequestTransport debug info  %#v don't contains response header", info)
	}

	if !strings.Contains(info, "\r\n\r\ntest response body") {
		t.Errorf("DebugRequestTransport debug info  %#v don't contains response body", info)
	}
}
