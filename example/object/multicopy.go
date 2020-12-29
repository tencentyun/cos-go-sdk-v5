package main

import (
	"context"
	"net/url"
	"os"

	"net/http"

	"fmt"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func log_status(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
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

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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

	opt := &cos.MultiCopyOptions{
		OptCopy: &cos.ObjectCopyOptions{
			&cos.ObjectCopyHeaderOptions{
				XCosStorageClass: "Archive",
			},
			nil,
		},
		ThreadPoolSize: 10,
	}
	source := "exampleobject"
	soruceURL := fmt.Sprintf("%s/%s", u.Host, source)
	dest := fmt.Sprintf("destobject")
	res, _, err := c.Object.MultiCopy(context.Background(), dest, soruceURL, opt)
	log_status(err)
	fmt.Printf("res:%+v\n", res)
}
