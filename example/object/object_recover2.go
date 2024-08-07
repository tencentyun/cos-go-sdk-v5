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

func logStatus(err error) {
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

// 恢复多版本桶srcBucket的数据
var (
	srcBucket       = "test-1259654469"
	srcBucketRegion = "ap-guangzhou"

	srcCosClient *cos.Client

	copyObjs = map[string]struct{}{}
)

func newClient(bucket, region string) *cos.Client {
	u, _ := url.Parse(fmt.Sprintf("https://%v.cos.%v.myqcloud.com", bucket, region))
	b := &cos.BaseURL{
		BucketURL: u,
	}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
}

// 恢复数据
func recoverObj(key, versionId string) {
	sourceURL := fmt.Sprintf("%v.cos.%v.myqcloud.com/%v?versionId=%v", srcBucket, srcBucketRegion, key, versionId)
	_, _, err := srcCosClient.Object.MultiCopy(context.Background(), key, sourceURL, nil)
	if err != nil {
		logStatus(err)
	}
}

func main() {
	srcCosClient = newClient(srcBucket, srcBucketRegion)

	keyMarker := ""
	versionIdMarker := ""
	isTruncated := true
	opt := &cos.BucketGetObjectVersionsOptions{
		EncodingType: "url",
	}
	for isTruncated {
		opt.KeyMarker = keyMarker
		opt.VersionIdMarker = versionIdMarker
		v, _, err := srcCosClient.Bucket.GetObjectVersions(context.Background(), opt)
		if err != nil {
			logStatus(err)
			break
		}
		for _, vc := range v.DeleteMarker {
			if vc.IsLatest {
				// 对象被删除，需要恢复
				copyObjs[vc.Key] = struct{}{}
			}
		}
		for _, vc := range v.Version {
			// 按最新恢复
			if _, ok := copyObjs[vc.Key]; ok {
				delete(copyObjs, vc.Key)
				key, _ := cos.DecodeURIComponent(vc.Key)
				fmt.Printf("key: %v, versionId: %v\n", key, vc.VersionId)
				recoverObj(key, vc.VersionId)
			}
		}
		keyMarker = v.NextKeyMarker
		versionIdMarker = v.NextVersionIdMarker
		isTruncated = v.IsTruncated
	}
}
