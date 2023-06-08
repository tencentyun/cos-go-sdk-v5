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

// GetMediaQueue 获取媒体处理队列
func GetMediaQueue() {
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
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		PageNumber: 1,
		PageSize:   2,
		Category:   "CateAll",
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

// GetPicQueue 获取图片处理队列
func GetPicQueue() {
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
	DescribeQueueOpt := &cos.DescribePicProcessQueuesOptions{
		PageNumber: 1,
		PageSize:   2,
		Category:   "CateAll",
	}
	DescribeQueueRes, _, err := c.CI.DescribePicProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

// GetAIQueue 获取AI 内容识别队列
func GetAIQueue() {
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
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		PageNumber: 1,
		PageSize:   2,
		Category:   "CateAll",
	}
	DescribeQueueRes, _, err := c.CI.DescribeAIProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

// GetASRQueue 获取语音识别队列
func GetASRQueue() {
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
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		PageNumber: 1,
		PageSize:   2,
		Category:   "CateAll",
	}
	DescribeQueueRes, _, err := c.CI.DescribeASRProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

func UpdateMediaQueue() {
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
	DescribeQueueOpt := &cos.UpdateMediaProcessQueueOptions{
		Name:    "queue-transcode",
		QueueID: "pa27b2bd96bef43b6baba820175485532",
		State:   "Active",
		NotifyConfig: &cos.MediaProcessQueueNotifyConfig{
			State:        "On",
			Url:          "http://www.callback.com",
			Event:        "TaskFinish",
			Type:         "Url",
			ResultFormat: "JSON",
		},
	}
	DescribeQueueRes, _, err := c.CI.UpdateMediaProcessQueue(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

func UpdatePicQueue() {
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
	DescribeQueueOpt := &cos.UpdateMediaProcessQueueOptions{
		Name:    "queue-pic",
		QueueID: "pc0393837f562409586a051979cad0d72",
		State:   "Active",
		NotifyConfig: &cos.MediaProcessQueueNotifyConfig{
			State:        "On",
			Url:          "http://www.callback.com",
			Event:        "TaskFinish",
			Type:         "Url",
			ResultFormat: "JSON",
		},
	}
	DescribeQueueRes, _, err := c.CI.UpdateMediaProcessQueue(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

func UpdateAIQueue() {
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
	DescribeQueueOpt := &cos.UpdateMediaProcessQueueOptions{
		Name:    "queue-ai",
		QueueID: "pa7b0400f4e0041ac849ab12104dedce9",
		State:   "Active",
		NotifyConfig: &cos.MediaProcessQueueNotifyConfig{
			State:        "On",
			Url:          "http://www.callback.com",
			Event:        "TaskFinish",
			Type:         "Url",
			ResultFormat: "JSON",
		},
	}
	DescribeQueueRes, _, err := c.CI.UpdateMediaProcessQueue(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

func UpdateASRQueue() {
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
	DescribeQueueOpt := &cos.UpdateMediaProcessQueueOptions{
		Name:    "queue-asr",
		QueueID: "pe91d0af11fc14337987ff0c34f8b0886",
		State:   "Active",
		NotifyConfig: &cos.MediaProcessQueueNotifyConfig{
			State:        "On",
			Url:          "http://www.callback.com",
			Event:        "TaskFinish",
			Type:         "Url",
			ResultFormat: "JSON",
		},
	}
	DescribeQueueRes, _, err := c.CI.UpdateMediaProcessQueue(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
}

func main() {
	// UpdateMediaQueue()
	// UpdatePicQueue()
	// UpdateAIQueue()
	// UpdateASRQueue()
	GetMediaQueue()
	GetPicQueue()
	GetAIQueue()
	GetASRQueue()
}
