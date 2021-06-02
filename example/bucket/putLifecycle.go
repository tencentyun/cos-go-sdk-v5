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

	lc := &cos.BucketPutLifecycleOptions{
		Rules: []cos.BucketLifecycleRule{
			{
				ID:     "1234",
				Filter: &cos.BucketLifecycleFilter{Prefix: "test"},
				Status: "Enabled",
				Transition: []cos.BucketLifecycleTransition{
					{
						Days:         30,
						StorageClass: "STANDARD_IA",
					},
					{
						Days:         90,
						StorageClass: "ARCHIVE",
					},
				},
				Expiration: &cos.BucketLifecycleExpiration{
					Days: 360,
				},
				NoncurrentVersionExpiration: &cos.BucketLifecycleNoncurrentVersion{
					NoncurrentDays: 360,
				},
				NoncurrentVersionTransition: []cos.BucketLifecycleNoncurrentVersion{
					{
						NoncurrentDays: 90,
						StorageClass:   "ARCHIVE",
					},
					{
						NoncurrentDays: 180,
						StorageClass:   "DEEP_ARCHIVE",
					},
				},
				AbortIncompleteMultipartUpload: &cos.BucketLifecycleAbortIncompleteMultipartUpload{
					DaysAfterInitiation: 90,
				},
			},
		},
	}
	_, err := c.Bucket.PutLifecycle(context.Background(), lc)
	if err != nil {
		panic(err)
	}
}
