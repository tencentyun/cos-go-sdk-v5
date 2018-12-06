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
	u, _ := url.Parse("https://testhuanan-1253846586.cos.ap-guangzhou.myqcloud.com")
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

	lc := &cos.BucketPutLifecycleOptions{
		Rules: []cos.BucketLifecycleRule{
			{
				ID:     "1234",
				Filter: &cos.BucketLifecycleFilter{Prefix: "test"},
				Status: "Enabled",
				Transition: &cos.BucketLifecycleTransition{
					Days:         10,
					StorageClass: "Standard",
				},
			},
			{
				ID: "123422",
				// If used for all objecs set Prefix:""
				Filter: &cos.BucketLifecycleFilter{Prefix: "gg"},
				Status: "Disabled",
				Expiration: &cos.BucketLifecycleExpiration{
					Days: 10,
				},
			},
		},
	}
	_, err := c.Bucket.PutLifecycle(context.Background(), lc)
	if err != nil {
		panic(err)
	}
}
