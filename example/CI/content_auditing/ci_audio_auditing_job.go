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
	bu, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	cu, _ := url.Parse("https://test-1259654469.ci.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: bu, CIURL: cu}
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
	opt := &cos.PutAudioAuditingJobOptions{
		InputObject: "test.mp3",
		Conf:        &cos.AudioAuditingJobConf{},
	}

	res, _, err := c.CI.PutAudioAuditingJob(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)

	time.Sleep(3 * time.Second)
	res2, _, err := c.CI.GetAudioAuditingJob(context.Background(), res.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", res2)
}
