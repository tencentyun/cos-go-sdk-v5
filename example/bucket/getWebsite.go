package main

import (
	"context"
	"github.com/agin719/cos-go-sdk-v5"
	"github.com/agin719/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
  "encoding/xml"
  "os"
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

	res, _, err := c.Bucket.GetWebsite(context.Background())
	if err != nil {
		panic(err)
	}
  body, _ := xml.Marshal(res)
  os.Stdout.Write(body)
}

