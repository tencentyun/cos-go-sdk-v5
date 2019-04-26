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
	u, _ := url.Parse("https://alangz-1251668577.cos.ap-guangzhou.myqcloud.com")
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

	id := "test1"
	opt := &cos.BucketPutInventoryOptions{
		ID: id,
		// True or False
		IsEnabled:              "True",
		IncludedObjectVersions: "All",
		Filter: &cos.BucketInventoryFilter{
			Prefix: "test",
		},
		OptionalFields: &cos.BucketInventoryOptionalFields{
			BucketInventoryFields: []string{
				"Size", "LastModifiedDate",
			},
		},
		Schedule: &cos.BucketInventorySchedule{
			// Weekly or Daily
			Frequency: "Daily",
		},
		Destination: &cos.BucketInventoryDestination{
			BucketDestination: &cos.BucketInventoryDestinationContent{
				Bucket: "qcs::cos:ap-guangzhou::alangz-1251668577",
				Format: "CSV",
			},
		},
	}
	_, err := c.Bucket.PutBucketInventory(context.Background(), id, opt)
	if err != nil {
		panic(err)
	}
}
