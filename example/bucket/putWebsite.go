package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

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

	opt := &cos.BucketPutWebsiteOptions{
		Index:            "index.html",
		Error:            &cos.ErrorDocument{"index_backup.html"},
		RedirectProtocol: &cos.RedirectRequestsProtocol{"https"},
		RoutingRules: &cos.WebsiteRoutingRules{
			[]cos.WebsiteRoutingRule{
				{
					ConditionErrorCode: "404",
					RedirectProtocol:   "https",
					RedirectReplaceKey: "404.html",
				},
				{
					ConditionPrefix:          "docs/",
					RedirectProtocol:         "https",
					RedirectReplaceKeyPrefix: "documents/",
				},
			},
		},
	}

	_, err := c.Bucket.PutWebsite(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}
