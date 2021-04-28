package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

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
	cu, _ := url.Parse("http://test-1259654469.pic.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	_, err := c.CI.PutGuetzli(context.Background())
	log_status(err)
	res, _, err := c.CI.GetGuetzli(context.Background())
	log_status(err)
	if res != nil && res.GuetzliStatus != "on" {
		fmt.Printf("Error Status: %v\n", res.GuetzliStatus)
	}
	time.Sleep(time.Second * 3)
	_, err = c.CI.DeleteGuetzli(context.Background())
	log_status(err)
	res, _, err = c.CI.GetGuetzli(context.Background())
	log_status(err)
	if res != nil && res.GuetzliStatus != "off" {
		fmt.Printf("Error Status: %v\n", res.GuetzliStatus)
	}

}
