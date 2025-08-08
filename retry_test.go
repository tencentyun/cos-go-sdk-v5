package cos

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	normal_domain = "cos-sdk-err-retry-1253960454.cos.ap-chengdu.myqcloud.com"
	backup_domain = "cos-sdk-err-retry-1253960454.cos.ap-chengdu.tencentcos.cn"
	nocos_domain  = os.Getenv("ERR_HOST")
)

func retrysetup(host string) *Client {
	u, _ := url.Parse("http://" + host)
	cli := NewClient(&BaseURL{u, u, u, u, u, u}, &http.Client{
		Timeout: 6 * time.Second,
	})
	cli.Conf.RetryOpt.Count = 2
	return cli
}

func checkRetry(t *testing.T, resp *Response, domain string, retry bool, switched bool, e error) bool {
	if resp == nil {
		t.Errorf("checkRetry resp is nil, err: %v", e)
		return false
	}
	resp.Body.Close()
	// 不重试
	if !retry && resp.Request.Header.Get("X-Cos-Sdk-Retry") != "" {
		t.Errorf("X-Cos-Sdk-Retry is not empty, %v", resp.Request.Header.Get("X-Cos-Sdk-Retry"))
		return false
	}
	// 重试
	if retry && resp.Request.Header.Get("X-Cos-Sdk-Retry") != "true" {
		t.Errorf("X-Cos-Sdk-Retry is not true, %v", resp.Request.Header.Get("X-Cos-Sdk-Retry"))
		return false
	}
	// 不切换域名
	if !switched && resp.Request.Host != domain {
		t.Errorf("host is switch, return: %v, expect: %v", resp.Request.Host, domain)
		return false
	}
	// 切换域名
	if switched && resp.Request.Host != backup_domain {
		t.Errorf("host is not switch, return: %v, expect: %v", resp.Request.Host, backup_domain)
		return false
	}
	return true
}

func Test_Retry_GetObject_normaldomain_noswitch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := normal_domain
	cli := retrysetup(domain)
	var resp *Response
	var err error
	resp, err = cli.Object.Get(context.Background(), "200r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "200", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "206r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "500r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "500", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "shutdown", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "timeout", nil)
	checkRetry(t, resp, domain, true, false, err)
}

func Test_Retry_GetObject_nocosdomain_noswitch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := nocos_domain
	cli := retrysetup(domain)
	var resp *Response
	var err error
	resp, err = cli.Object.Get(context.Background(), "200r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "200", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "206r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "500r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "500", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "shutdown", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "timeout", nil)
	checkRetry(t, resp, domain, true, false, err)
}

func Test_Retry_GetObject_normaldomain_switch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := normal_domain
	cli := retrysetup(domain)
	cli.Conf.RetryOpt.AutoSwitchHost = true
	var resp *Response
	var err error
	resp, err = cli.Object.Get(context.Background(), "200r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "200", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "206r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "302r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "307r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "400r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "500r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "500", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "503r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "504r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "shutdown", nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Get(context.Background(), "timeout", nil)
	checkRetry(t, resp, domain, true, true, err)
}

func Test_Retry_GetObject_nocosdomain_switch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := nocos_domain
	cli := retrysetup(domain)
	cli.Conf.RetryOpt.AutoSwitchHost = true
	var resp *Response
	var err error
	resp, err = cli.Object.Get(context.Background(), "200r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "200", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "204", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "206r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "301", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "302", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "307", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "400", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "403", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404r", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "404", nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Get(context.Background(), "500r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "500", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "503", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504r", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "504", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "shutdown", nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Get(context.Background(), "timeout", nil)
	checkRetry(t, resp, domain, true, false, err)
}

func Test_Retry_PutObject_normaldomain_noswitch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := normal_domain
	cli := retrysetup(domain)
	var resp *Response
	var err error
	resp, err = cli.Object.Put(context.Background(), "200r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "200", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "500r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "500", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "shutdown", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "timeout", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
}

func Test_Retry_PutObject_nocosdomain_noswitch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := nocos_domain
	cli := retrysetup(domain)
	var resp *Response
	var err error

	resp, err = cli.Object.Put(context.Background(), "200r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "200", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "500r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "500", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "shutdown", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "timeout", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
}

func Test_Retry_PutObject_normaldomain_switch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := normal_domain
	cli := retrysetup(domain)
	cli.Conf.RetryOpt.AutoSwitchHost = true
	var resp *Response
	var err error
	resp, err = cli.Object.Put(context.Background(), "200r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "200", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "302r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "307r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "400r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "500r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "500", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "503r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "504r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "shutdown", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
	resp, err = cli.Object.Put(context.Background(), "timeout", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, true, err)
}

func Test_Retry_PutObject_nocosdomain_switch(t *testing.T) {
	if os.Getenv("ERR_HOST") == "" {
		t.Skip("ERR_HOST is empty, skip")
	}
	domain := nocos_domain
	cli := retrysetup(domain)
	cli.Conf.RetryOpt.AutoSwitchHost = true
	var resp *Response
	var err error
	resp, err = cli.Object.Put(context.Background(), "200r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "200", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "206", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "301", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "302", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "307", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "400", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "403", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "404", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, false, false, err)
	resp, err = cli.Object.Put(context.Background(), "500r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "500", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "503", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504r", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "504", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "shutdown", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
	resp, err = cli.Object.Put(context.Background(), "timeout", strings.NewReader(""), nil)
	checkRetry(t, resp, domain, true, false, err)
}
