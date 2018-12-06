package main

import (
	"context"
	"net/url"
	"os"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/hello.txt"
	opt := &cos.ObjectOptionsOptions{
		Origin: "http://www.qq.com",
		AccessControlRequestMethod: "PUT",
	}
	_, err := c.Object.Options(context.Background(), name, opt)
	if err != nil {
		panic(err)
	}
}
