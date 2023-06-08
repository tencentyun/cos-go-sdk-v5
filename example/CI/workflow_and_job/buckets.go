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

func getClient() *cos.Client {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	return c
}

// 查询开通媒体处理的存储桶
func describeMediaBucket() {
	c := getClient()
	opt := &cos.DescribeMediaProcessBucketsOptions{
		Regions: "ap-chongqing",
	}
	res, _, err := c.CI.DescribeMediaProcessBuckets(context.Background(), opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)
}

// 查询开通图片处理（异步）的存储桶
func describePicBucket() {
	c := getClient()
	opt := &cos.DescribePicProcessBucketsOptions{
		Regions: "ap-chongqing",
	}
	res, _, err := c.CI.DescribePicProcessBuckets(context.Background(), opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)
}

// 查询开通文档预览的存储桶
func describeDocBucket() {
	c := getClient()
	opt := &cos.DescribeDocProcessBucketsOptions{
		Regions: "ap-chongqing",
	}
	res, _, err := c.CI.DescribeDocProcessBuckets(context.Background(), opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)
}

// 查询开通AI 内容识别（异步）的存储桶
func describeAIBucket() {
	c := getClient()
	opt := &cos.DescribeAIProcessBucketsOptions{
		Regions: "ap-chongqing",
	}
	res, _, err := c.CI.DescribeAIProcessBuckets(context.Background(), opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)
}

// 查询开通智能语音的存储桶
func describeASRBucket() {
	c := getClient()
	opt := &cos.DescribeASRProcessBucketsOptions{
		Regions: "ap-chongqing",
	}
	res, _, err := c.CI.DescribeASRProcessBuckets(context.Background(), opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)
}

// 查询开通文件处理的存储桶
func describeFileBucket() {
	c := getClient()
	opt := &cos.DescribeFileProcessBucketsOptions{
		Regions: "ap-chongqing",
	}
	res, _, err := c.CI.DescribeFileProcessBuckets(context.Background(), opt)
	log_status(err)
	fmt.Printf("res: %+v\n", res)
}

func main() {
	// describeMediaBucket()
	// describePicBucket()
	// describeDocBucket()
	// describeAIBucket()
	// describeASRBucket()
	describeFileBucket()
}
