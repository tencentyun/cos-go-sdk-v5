package main

import (
	"context"
	"fmt"

	"github.com/tencentyun/cos-go-sdk-v5"
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

// InvokeDNAJob 提交一个DNA任务
// https://cloud.tencent.com/document/product/460/96115
func InvokeDNAJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "DNA",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			DnaConfig: &cos.DnaConfig{
				RuleType: "GetFingerPrint",
				DnaDbId:  "xxx",
			},
			UserData: "This is my DNA job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// GetDnaDb 查询 DNA 库列表
// https://cloud.tencent.com/document/product/460/96117
func GetDnaDb() {
	c := getClient()
	opt := &cos.GetDnaDbOptions{
		PageNumber: "2",
		PageSize:   "10",
	}
	res, _, err := c.CI.GetDnaDb(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// GetDnaDbFiles 获取 DNA 库中文件列表
// https://cloud.tencent.com/document/product/460/96116
func GetDnaDbFiles() {
	c := getClient()
	opt := &cos.GetDnaDbFilesOptions{
		PageNumber: "2",
		PageSize:   "10",
	}
	res, _, err := c.CI.GetDnaDbFiles(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func main() {
}
