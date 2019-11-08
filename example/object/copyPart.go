package main

import (
	"context"
	"fmt"
	"os"

	"net/url"

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

	sourceUrl := "test-1253846586.cos.ap-guangzhou.myqcloud.com/source/copy_multi_upload.go"
	name := "test/test_multi_upload.go"
	up := initUpload(c, name)
	uploadID := up.UploadID

	opt := &cos.ObjectCopyPartOptions{}
	res, _, err := c.Object.CopyPart(
		context.Background(), name, uploadID, 1, sourceUrl, opt)
	if err != nil {
		panic(err)
	}
	fmt.Println("ETag:", res.ETag)

	completeOpt := &cos.CompleteMultipartUploadOptions{}
	completeOpt.Parts = append(completeOpt.Parts, cos.Object{
		PartNumber: 1,
		ETag:       res.ETag,
	})
	v, resp, err := c.Object.CompleteMultipartUpload(
		context.Background(), name, uploadID, completeOpt,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", resp.Status)
	fmt.Printf("%#v\n", v)
	fmt.Printf("%s\n", v.Location)
}
