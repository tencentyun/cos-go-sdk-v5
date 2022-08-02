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

func createDocJob() {
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

	createJobOpt := &cos.CreateDocProcessJobsOptions{
		Tag: "DocProcess",
		Input: &cos.DocProcessJobInput{
			Object: "input/doc_preview.ppt",
		},
		Operation: &cos.DocProcessJobOperation{
			Output: &cos.DocProcessJobOutput{
				Region: "ap-chongqing",
				Object: "test-doc${Number}.png",
				Bucket: "test-1234567890",
			},
			DocProcess: &cos.DocProcessJobDocProcess{
				TgtType:     "png",
				StartPage:   1,
				EndPage:     2,
				ImageParams: "watermark/2/text/5paH5qGj6aKE6KeI/fontsize/20/gravity/NorthEast",
			},
		},
		QueueId: "p363fa5add7b94ca693f667a8d4807f54",
	}
	createJobRes, _, err := c.CI.CreateDocProcessJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// 5、DescribeDocProcessJob
	DescribeJobRes, _, err := c.CI.DescribeDocProcessJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func describeDocJob() {
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

	jobId := "d695b1bc6117c11eda300ad97400fb982"
	DescribeJobRes, _, err := c.CI.DescribeDocProcessJob(context.Background(), jobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)

}

func describeDocJobs() {
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

	DescribeJobsOpt := &cos.DescribeDocProcessJobsOptions{
		QueueId: "p363fa5add7b94ca693f667a8d4807f54",
		Tag:     "DocProcess",
	}
	DescribeJobsRes, _, err := c.CI.DescribeDocProcessJobs(context.Background(), DescribeJobsOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobsRes)
}

func main() {
	// createDocJob()
	// describeDocJob()
	// describeDocJobs()
}
