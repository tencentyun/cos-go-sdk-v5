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
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
	opt := &cos.BucketPutObjectLockOptions{
		ObjectLockEnabled: "Enabled",
		Rule: &cos.ObjectLockRule{
			Days: 1,
		},
	}
	_, err := c.Bucket.PutObjectLockConfiguration(context.Background(), opt)
	logStatus(err)

	res, _, err := c.Bucket.GetObjectLockConfiguration(context.Background())
	logStatus(err)
	fmt.Printf("%+v\n", res)

	ropt := &cos.ObjectPutRetentionOptions{
		RetainUntilDate: "2022-12-10T08:34:48.000Z",
		Mode: "COMPLIANCE",
	}
	_, err = c.Object.PutRetention(context.Background(), "test", ropt)
	logStatus(err)

	r, _, err := c.Object.GetRetention(context.Background(), "test", nil)
	logStatus(err)
	fmt.Printf("%+v\n", r)
}
