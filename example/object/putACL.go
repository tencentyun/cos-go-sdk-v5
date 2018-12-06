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

	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"
	_, err := c.Object.PutACL(context.Background(), name, opt)
	if err != nil {
		panic(err)
	}

	// with body
	opt = &cos.ObjectPutACLOptions{
		Body: &cos.ACLXml{
			Owner: &cos.Owner{
				ID: "qcs::cam::uin/100000760461:uin/100000760461",
			},
			AccessControlList: []cos.ACLGrant{
				{
					Grantee: &cos.ACLGrantee{
						Type: "RootAccount",
						ID:   "qcs::cam::uin/100000760461:uin/100000760461",
					},

					Permission: "FULL_CONTROL",
				},
			},
		},
	}

	_, err = c.Object.PutACL(context.Background(), name, opt)
	if err != nil {
		panic(err)
	}
}
