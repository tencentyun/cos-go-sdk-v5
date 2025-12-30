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

func logStatus(err error) {
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
	u, _ := url.Parse("https://cd-1259654469.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
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
					FollowQueryString: cos.Bool(true),
					HttpHeader: &cos.BucketOriginHttpHeader{
						NewHttpHeaders: []cos.OriginHttpHeader{
							{
								Key:   "Content-Type",
								Value: "csv",
							},
						},
						FollowHttpHeaders: []cos.OriginHttpHeader{
							{
								Key: "Content-Type",
							},
						},
					},
					FollowRedirection: cos.Bool(true),
				},
				OriginInfo: &cos.BucketOriginInfo{
					HostInfo: &cos.BucketOriginHostInfo{
						HostName: "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
						StandbyHostName: []*cos.BucketOriginStandbyHost{
							&cos.BucketOriginStandbyHost{
								Index:    1,
								HostName: "www.qq.com",
							},
							&cos.BucketOriginStandbyHost{
								Index:    2,
								HostName: "www.myqlcoud.com",
							},
						},
					},
				},
			},
		},
	}

	_, err := c.Bucket.PutOrigin(context.Background(), opt)
	logStatus(err)
	res, _, err := c.Bucket.GetOrigin(context.Background())
	logStatus(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.Rule)
	for _, rule := range res.Rule {
		fmt.Printf("%+v\n", rule.OriginInfo.HostInfo)
	}
	_, err = c.Bucket.DeleteOrigin(context.Background())
	logStatus(err)
}
