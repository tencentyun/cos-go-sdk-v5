package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

// DescribeInventoryTriggerJob TODO
func DescribeInventoryTriggerJob() {
	u, _ := url.Parse("https://test-123456789.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-123456789.ci.ap-chongqing.myqcloud.com")
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
	jobId := "babc6cfc8cd2111ecb09a52540038936c"
	DescribeWorkflowRes, _, err := c.CI.DescribeInventoryTriggerJob(context.Background(), jobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

// DescribeInventoryTriggerJobs TODO
func DescribeInventoryTriggerJobs() {
	u, _ := url.Parse("https://test-123456789.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-123456789.ci.ap-chongqing.myqcloud.com")
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
	opt := &cos.DescribeInventoryTriggerJobsOptions{
		States: "All",
	}
	DescribeWorkflowRes, _, err := c.CI.DescribeInventoryTriggerJobs(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

// CreateInventoryTriggerJob TODO
func CreateInventoryTriggerJob() {
	u, _ := url.Parse("https://test-123456789.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-123456789.ci.ap-chongqing.myqcloud.com")
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
	rand.Seed(time.Now().UnixNano())
	opt := &cos.CreateInventoryTriggerJobOptions{
		Name: "trigger-" + strconv.Itoa(rand.Intn(100)),
		Input: &cos.InventoryTriggerJobInput{
			Manifest: "https://test-123456789.cos.ap-chongqing.myqcloud.com/cos_bucket_inventory/123456789/test/menu_instant_20220506171340/20220506/manifest.json",
		},
		Operation: &cos.InventoryTriggerJobOperation{
			WorkflowIds: "web6ac56c1ef54dbfa44d7f4103203be9",
			TimeInterval: cos.InventoryTriggerJobOperationTimeInterval{
				Start: "2002-02-16T10:45:12+0800",
				End:   "2022-05-16T10:45:12+0800",
			},
		},
	}
	DescribeInventoryTriggerJobRes, _, err := c.CI.CreateInventoryTriggerJob(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeInventoryTriggerJobRes)
}

// CreateInventoryTriggerJobByParam TODO
func CreateInventoryTriggerJobByParam() {
	u, _ := url.Parse("https://test-123456789.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-123456789.ci.ap-chongqing.myqcloud.com")
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
	rand.Seed(time.Now().UnixNano())
	opt := &cos.CreateInventoryTriggerJobOptions{
		Name: "trigger-" + strconv.Itoa(rand.Intn(100)),
		Input: &cos.InventoryTriggerJobInput{
			Prefix: "input/",
		},
		Type: "Job",
		Operation: &cos.InventoryTriggerJobOperation{
			QueueId: "pa27b2bd96bef43b6baba820175485532",
			TimeInterval: cos.InventoryTriggerJobOperationTimeInterval{
				Start: "2002-02-16T10:45:12+0800",
				End:   "2023-05-16T10:45:12+0800",
			},
			Tag: "Transcode",
			JobParam: &cos.InventoryTriggerJobOperationJobParam{
				TemplateId: "t00daf332ba39049f8bfb899c1ed0134b0",
			},
			JobLevel: 1,
			UserData: "This is my CreateInventoryTriggerJob",
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-123456789",
				Object: "output/${InputName}_${InventoryTriggerJobId}.${ext}",
			},
		},
	}
	DescribeInventoryTriggerJobRes, _, err := c.CI.CreateInventoryTriggerJob(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeInventoryTriggerJobRes)
}

// CancelInventoryTriggerJobs TODO
func CancelInventoryTriggerJobs() {
	u, _ := url.Parse("https://test-123456789.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-123456789.ci.ap-chongqing.myqcloud.com")
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
	jobId := "b56a3bbc0cd3011ecb09a52540038936c"
	_, err := c.CI.CancelInventoryTriggerJob(context.Background(), jobId)
	log_status(err)
}

func main() {
}
