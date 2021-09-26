package main

import (
	"context"
	"fmt"
	"github.com/agin719/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"os"
	"time"
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
	u, e := url.Parse("http://ap-guangzhou.migration.myqcloud.com")
	log_status(e)
	b := &cos.BaseURL{BucketURL: u, FetchURL: u}
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
	bucket := "test-1259654469"
	opt := &cos.PutFetchTaskOptions{
		Url: "http://" + bucket + ".cos.ap-guangzhou.myqcloud.com/exampleobject",
		Key: "exampleobject",
	}

	res, _, err := c.Object.PutFetchTask(context.Background(), bucket, opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)

	time.Sleep(time.Second * 3)

	rs, _, err := c.Object.GetFetchTask(context.Background(), bucket, res.Data.TaskId)
	log_status(err)
	fmt.Printf("res: %+v\n", rs)
}
