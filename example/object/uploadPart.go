package main

import (
	"context"
	"fmt"
	"os"

	"net/url"
	"strings"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func initUpload(c *cos.Client, name string) *cos.InitiateMultipartUploadResult {
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", v)
	return v
}

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

	name := "test/test_multi_upload.go"
	up := initUpload(c, name)
	uploadID := up.UploadID

	f := strings.NewReader("test heoo")
	_, err := c.Object.UploadPart(
		context.Background(), name, uploadID, 1, f, nil,
	)
	if err != nil {
		panic(err)
	}
}
