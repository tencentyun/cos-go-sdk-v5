package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

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
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})

	// Case1 多线程上传对象
	opt := &cos.MultiUploadOptions{
		ThreadPoolSize: 5,
	}
	v, _, err := c.Object.Upload(
		context.Background(), "gomulput1G", "./test1G", opt,
	)
	log_status(err)
	fmt.Printf("Case1 done, %v\n", v)

	// Case2 多线程上传对象，查看上传进度
	opt.OptIni = &cos.InitiateMultipartUploadOptions{
		nil,
		&cos.ObjectPutHeaderOptions{
			Listener: &cos.DefaultProgressListener{},
		},
	}
	v, _, err = c.Object.Upload(
		context.Background(), "gomulput1G", "./test1G", opt,
	)
	log_status(err)
	fmt.Printf("Case2 done, %v\n", v)

}
