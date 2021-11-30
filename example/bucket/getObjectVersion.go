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

func log_status(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Printf("ERROR: %v\n", err)
		// ERROR
	}
}

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

	keyMarker := ""
	versionIdMarker := ""
	isTruncated := true
	opt := &cos.BucketGetObjectVersionsOptions{}
	for isTruncated {
		opt.KeyMarker = keyMarker
		opt.VersionIdMarker = versionIdMarker
		v, _, err := c.Bucket.GetObjectVersions(context.Background(), opt)
		if err != nil {
			log_status(err)
			break
		}
		for _, vc := range v.Version {
			fmt.Printf("Version: %v, %v, %v, %v\n", vc.Key, vc.Size, vc.VersionId, vc.IsLatest)
		}
		for _, dc := range v.DeleteMarker {
			fmt.Printf("DeleteMarker: %v, %v, %v\n", dc.Key, dc.VersionId, dc.IsLatest)
		}
		keyMarker = v.NextKeyMarker
		versionIdMarker = v.NextVersionIdMarker
		isTruncated = v.IsTruncated
	}
}
