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

	v, _, err := c.Bucket.GetLifecycle(context.Background())
	if err != nil {
		panic(err)
	}
	for _, r := range v.Rules {
		fmt.Printf("%+v\n", r.ID)
		fmt.Printf("%+v\n", r.Filter)
		fmt.Printf("%+v\n", r.Status)
		fmt.Printf("%+v\n", r.Transition)
		fmt.Printf("%+v\n", r.Expiration)
		fmt.Printf("%+v\n", r.NoncurrentVersionExpiration)
		fmt.Printf("%+v\n", r.NoncurrentVersionTransition)
		fmt.Printf("%+v\n", r.AbortIncompleteMultipartUpload)
	}
}
