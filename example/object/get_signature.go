package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func log_status(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Printf("ERROR: %v\n", err)
		// ERROR
	}
}

func main() {
	ak := os.Getenv("COS_SECRETID")
	sk := os.Getenv("COS_SECRETKEY")
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, nil)

	name := "中文测试"

	// 把相关header和query签入到签名中
	opt := &cos.PresignedURLOptions{
		Query:  &url.Values{},
		Header: &http.Header{},
	}
	opt.Query.Add("test_key", "中文测试")
	opt.Header.Add("x-cos-meta-test", "中文测试")

	// 获取签名
	auth := c.Object.GetSignature(context.Background(), http.MethodPut, name, ak, sk, time.Hour, opt)
	fmt.Printf("signature: %s\n", auth)

	cli := &http.Client{
		Transport: &debug.DebugRequestTransport{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
		},
	}
	// 请求需要和签名保持一致
	uristr := fmt.Sprintf("https://%s/%s?%s", u.Host, name, opt.Query.Encode())
	req, err := http.NewRequest(http.MethodPut, uristr, strings.NewReader("test"))
	req.Header.Add("x-cos-meta-test", "中文测试")
	req.Header.Add("Authorization", auth)

	resp, err := cli.Do(req)
	log_status(err)
	defer resp.Body.Close()
}
