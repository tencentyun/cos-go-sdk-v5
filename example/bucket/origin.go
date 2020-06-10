package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func log_status(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("Code: %v\n", e.Code)
		fmt.Printf("Message: %v\n", e.Message)
		fmt.Printf("Resource: %v\n", e.Resource)
		fmt.Printf("RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Println(err)
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

	opt := &cos.BucketPutOriginOptions{
		Rule: []cos.BucketOriginRule{
			{
				OriginType: "Proxy",
				OriginCondition: &cos.BucketOriginCondition{
					HTTPStatusCode: "404",
					Prefix:         "",
				},
				OriginParameter: &cos.BucketOriginParameter{
					Protocol:          "FOLLOW",
					FollowQueryString: true,
					HttpHeader: &cos.BucketOriginHttpHeader{
						NewHttpHeaders: []cos.OriginHttpHeader{
							{
								Key:   "x-cos-ContentType",
								Value: "csv",
							},
						},
						FollowHttpHeaders: []cos.OriginHttpHeader{
							{
								Key: "Content-Type",
							},
						},
					},
					FollowRedirection: true,
				},
				OriginInfo: &cos.BucketOriginInfo{
					HostInfo: "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
				},
			},
		},
	}

	_, err := c.Bucket.PutOrigin(context.Background(), opt)
	log_status(err)
	res, _, err := c.Bucket.GetOrigin(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.Rule)
	_, err = c.Bucket.DeleteOrigin(context.Background())
	log_status(err)
}
