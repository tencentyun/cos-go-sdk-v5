package main

import (
	"context"
	"fmt"
	"net/url"

	"net/http"
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
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
	// Get 请求配置重试
	c.Conf.RetryOpt.Count = 3                   // 错误重试次数，默认重试3次
	c.Conf.RetryOpt.Interval = time.Millisecond // 错误重试间隔时间，默认0
	c.Conf.RetryOpt.StatusCode = []int{}        // 默认5xx都会重试，该参数配置其余需要重试的响应码

	name := "exampleobject"
	_, err := c.Object.Get(context.Background(), name, nil)
	log_status(err)
}
