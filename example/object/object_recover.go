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
		fmt.Println("WARN: Resource is not existed: %v", err)
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

var (
	srcBucket       = "test-1259654469"
	dstBucket       = "test2-1259654469"
	srcBucketRegion = "ap-guangzhou"
	dstBucketRegion = "ap-guangzhou"

	srcCosClient *cos.Client
	dstCosClient *cos.Client
)

func newClient(bucket, region string) *cos.Client {
	u, _ := url.Parse(fmt.Sprintf("https://%v.cos.%v.myqcloud.com", bucket, region))
	b := &cos.BaseURL{
		BucketURL: u,
	}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
}

func recoverObj(key, versionId string) {
	sourceURL := fmt.Sprintf("%v.cos.%v.myqcloud.com/%v?versionId=%v", srcBucket, srcBucketRegion, key, versionId)
	_, _, err := dstCosClient.Object.MultiCopy(context.Background(), key, sourceURL, nil)
	if err != nil {
		log_status(err)
	}
}

func main() {
	srcCosClient = newClient(srcBucket, srcBucketRegion)
	dstCosClient = newClient(dstBucket, dstBucketRegion)

	keyMarker := ""
	versionIdMarker := ""
	isTruncated := true
	opt := &cos.BucketGetObjectVersionsOptions{
		EncodingType: "url",
	}
	var preKey string
	for isTruncated {
		opt.KeyMarker = keyMarker
		opt.VersionIdMarker = versionIdMarker
		v, _, err := srcCosClient.Bucket.GetObjectVersions(context.Background(), opt)
		if err != nil {
			log_status(err)
			break
		}
		for _, vc := range v.Version {
			// 最新非deletemarker对象
			if preKey != vc.Key {
				preKey = vc.Key
				key, _ := cos.DecodeURIComponent(vc.Key)
				fmt.Printf("key: %v, versionId: %v, lastest: %v\n", key, vc.VersionId, vc.IsLatest)
				recoverObj(key, vc.VersionId)
			}
		}
		keyMarker = v.NextKeyMarker
		versionIdMarker = v.NextVersionIdMarker
		isTruncated = v.IsTruncated
	}
}
