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
	u, _ := url.Parse("https://alanbj-1251668577.cos.ap-beijing.myqcloud.com")
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

	opt := &cos.PutBucketReplicationOptions{
		// qcs::cam::uin/[UIN]:uin/[Subaccount]
		Role: "qcs::cam::uin/2779643970:uin/2779643970",
		Rule: []cos.BucketReplicationRule{
			{
				ID: "1",
				// Enabled or Disabled
				Status: "Enabled",
				Destination: &cos.ReplicationDestination{
					// qcs::cos:[Region]::[Bucketname-Appid]
					Bucket: "qcs::cos:ap-guangzhou::alangz-1251668577",
				},
			},
		},
	}
	_, err := c.Bucket.PutBucketReplication(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}
