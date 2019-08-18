package main

import (
	"context"
	"github.com/agin719/cos-go-sdk-v5"
	"github.com/agin719/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"fmt"
)

func main() {
  u, _ := url.Parse("https://jojobucket-1259654469.cos.ap-guangzhou.myqcloud.com")
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

	opt := &cos.BucketDomainConfiguration{
    Status : "ENABLED",
    Name : "www.agin719.com",
    Type : "REST",
    ForcedReplacement: "CNAME",
	}

	resp, err := c.Bucket.PutDomain(context.Background(), opt)
	if err != nil {
		panic(err)
	}
  fmt.Println(resp.Header)
}

