package main

import (
	"context"
	"net/url"
	"os"
	"time"

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
	startTime := time.Now()

	tg := &cos.BucketPutTaggingOptions{
		TagSet: []cos.BucketTaggingTag{
			{
				Key:   "test_k2",
				Value: "test_v2",
			},
			{
				Key:   "test_k3",
				Value: "test_v3",
			},
			{
				Key:   startTime.Format("02_Jan_06_15_04_MST"),
				Value: "test_time",
			},
		},
	}
	_, err := c.Bucket.PutTagging(context.Background(), tg)
	if err != nil {
		panic(err)
	}
}
