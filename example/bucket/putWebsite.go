package main

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"os"
)

func main() {
  	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
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

	opt := &cos.BucketWebsiteConfiguration{
		Index: "index.html",
    	Error: "index_backup.html",
	    RedirectProtocol: "https",
    	Rules: []cos.WebsiteRoutingRule{
      		{
        		ConditionErrorCode: "404",
		        RedirectProtocol: "https",
    		    RedirectReplaceKey: "404.html",
      		},
			{
				ConditionPrefix : "docs/",
				RedirectProtocol : "https",
				RedirectReplaceKeyPrefix : "documents/",
			},
	
    	},
	}

	_, err := c.Bucket.PutWebsite(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}

