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

func getDocQueue() {
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

	DescribeQueueOpt := &cos.DescribeDocProcessQueuesOptions{
		//QueueIds:   "p111a8dd208104ce3b11c78398f658ca8,p4318f85d2aa14c43b1dba6f9b78be9b3,aacb2bb066e9c4478834d4196e76c49d3",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeDocProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

func updateDocQueue() {
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

	updateQueueOpt := &cos.UpdateDocProcessQueueOptions{
		Name:    "my doc queue",
		QueueID: "p363fa5add7b94ca693f667a8d4807f54",
		State:   "Active",
		NotifyConfig: &cos.DocProcessQueueNotifyConfig{
			State: "Off",
		},
	}
	updateQueueRes, _, err := c.CI.UpdateDocProcessQueue(context.Background(), updateQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", updateQueueRes)
}

func main() {
	// updateDocQueue()
	getDocQueue()
}
