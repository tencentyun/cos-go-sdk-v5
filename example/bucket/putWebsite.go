package main

import (
	"context"
	"github.com/agin719/cos-go-sdk-v5"
	"github.com/agin719/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
)

func main() {
	u, _ := url.Parse("https://bucket-1-1259654469.cos.ap-shenzhen-fsi.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "AKIDfcOzOmUkJfphOt6JJ6kCPQFsKfqrbIhu",
			SecretKey: "CCsLj86tUt6MUQAr44tBLNI3d3IxWvz1",
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	opt := &cos.BucketWebsiteConfiguration{
		Index: &cos.WebsiteIndex{Suffix: "index.html"},
	}

	_, err := c.Bucket.PutWebsite(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}

