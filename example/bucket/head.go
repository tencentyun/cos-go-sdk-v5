package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	resp, err := c.Bucket.Head(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
}
