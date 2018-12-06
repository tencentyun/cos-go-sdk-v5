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
	b := &cos.BaseURL{
		BucketURL: u,
	}
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

	opt := &cos.BucketPutCORSOptions{
		Rules: []cos.BucketCORSRule{
			{
				AllowedOrigins: []string{"http://www.qq.com"},
				AllowedMethods: []string{"PUT", "GET"},
				AllowedHeaders: []string{"x-cos-meta-test", "x-cos-xx"},
				MaxAgeSeconds:  500,
				ExposeHeaders:  []string{"x-cos-meta-test1"},
			},
			{
				ID:             "1234",
				AllowedOrigins: []string{"http://www.baidu.com", "twitter.com"},
				AllowedMethods: []string{"PUT", "GET"},
				MaxAgeSeconds:  500,
			},
		},
	}
	_, err := c.Bucket.PutCORS(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}
