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
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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

	ok, err := c.Bucket.IsExist(context.Background())
	if err == nil && ok {
		fmt.Printf("bucket exists\n")
	} else if err != nil {
		fmt.Printf("head bucket failed: %v\n", err)
	} else {
		fmt.Printf("bucket does not exist\n")
	}
}
