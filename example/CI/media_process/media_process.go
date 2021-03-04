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

func main() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	fmt.Printf("%+v\n", os.Getenv("COS_SECRETID"))
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Transcode",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/go_117374C.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Transcode: &cos.Transcode{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.264",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
				TimeInterval: &cos.TimeInterval{
					Start:    "10",
					Duration: "",
				},
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJobs(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}
