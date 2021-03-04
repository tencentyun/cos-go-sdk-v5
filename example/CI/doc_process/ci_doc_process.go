package main

import (
	"context"
	"fmt"
	"io"
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
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	cu, _ := url.Parse("https://test-1259654469.ci.ap-guangzhou.myqcloud.com")
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

	// 1、UpdateDocProcessQueue
	updateQueueOpt := &cos.UpdateDocProcessQueueOptions{
		Name:    "queue-doc-process-1",
		QueueID: "p111a8dd208104ce3b11c78398f658ca8",
		State:   "Active",
		NotifyConfig: &cos.DocProcessQueueNotifyConfig{
			State: "Off",
		},
	}
	updateQueueRes, _, err := c.CI.UpdateDocProcessQueue(context.Background(), updateQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", updateQueueRes)

	// 2、DescribeDocProcessQueues
	DescribeQueueOpt := &cos.DescribeDocProcessQueuesOptions{
		QueueIds:   "p111a8dd208104ce3b11c78398f658ca8,p4318f85d2aa14c43b1dba6f9b78be9b3,aacb2bb066e9c4478834d4196e76c49d3",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeDocProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)

	// 3、DescribeDocProcessBuckets
	BucketsOpt := &cos.DescribeDocProcessBucketsOptions{
		Regions: "All",
	}
	BucketsRes, _, err := c.CI.DescribeDocProcessBuckets(context.Background(), BucketsOpt)
	log_status(err)
	fmt.Printf("%+v\n", BucketsRes)

	// 4、CreateDocProcessJobs
	createJobOpt := &cos.CreateDocProcessJobsOptions{
		Tag: "DocProcess",
		Input: &cos.DocProcessJobInput{
			Object: "form.pdf",
		},
		Operation: &cos.DocProcessJobOperation{
			Output: &cos.DocProcessJobOutput{
				Region: "ap-guangzhou",
				Object: "test-doc${Number}",
				Bucket: "test-1259654469",
			},
			DocProcess: &cos.DocProcessJobDocProcess{
				TgtType:     "png",
				StartPage:   1,
				EndPage:     -1,
				ImageParams: "watermark/1/image/aHR0cDovL3Rlc3QwMDUtMTI1MTcwNDcwOC5jb3MuYXAtY2hvbmdxaW5nLm15cWNsb3VkLmNvbS8xLmpwZw==/gravity/southeast",
			},
		},
		QueueId: "p111a8dd208104ce3b11c78398f658ca8",
	}
	createJobRes, _, err := c.CI.CreateDocProcessJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// 5、DescribeDocProcessJob
	DescribeJobRes, _, err := c.CI.DescribeDocProcessJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)

	// 6、DescribeDocProcessJobs
	DescribeJobsOpt := &cos.DescribeDocProcessJobsOptions{
		QueueId: "p111a8dd208104ce3b11c78398f658ca8",
		Tag:     "DocProcess",
	}
	DescribeJobsRes, _, err := c.CI.DescribeDocProcessJobs(context.Background(), DescribeJobsOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobsRes)

	// 7、doc-preview
	opt := &cos.DocPreviewOptions{
		Page: 1,
	}
	resp, err := c.CI.DocPreview(context.Background(), "form.pdf", opt)
	log_status(err)
	fd, _ := os.OpenFile("form.pdf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	io.Copy(fd, resp.Body)
	fd.Close()

}
