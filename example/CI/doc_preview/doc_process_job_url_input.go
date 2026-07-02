//go:build example_doc_process_job_url_input
// +build example_doc_process_job_url_input

// 示例：通过 URL 输入源创建文档转码任务
// 运行：go run -tags=example_doc_process_job_url_input ./example/CI/doc_preview
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	u, _ := url.Parse("https://test-1250000000.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-1250000000.ci.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	// InvokeCreateDocProcessJobsWithUrl:
	// 使用第三方 URL 作为文档转码任务的输入源（无需将源文件上传到 COS）。
	// Input.Url 与 Input.Object 二选一。
	opt := &cos.CreateDocProcessJobsOptions{
		Tag: "DocProcess",
		Input: &cos.DocProcessJobInput{
			Url: "https://example.com/foo.docx", // 第三方可访问的文档 URL
		},
		Operation: &cos.DocProcessJobOperation{
			Output: &cos.DocProcessJobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1250000000",
				Object: "doc/output-${Number}.png",
			},
			DocProcess: &cos.DocProcessJobDocProcess{
				TgtType:   "png",
				StartPage: 1,
				EndPage:   -1,
			},
		},
		QueueId: "p532fdead78444e649e1a4467c1cd19d3",
	}

	res, _, err := c.CI.CreateDocProcessJobs(context.Background(), opt)
	if err != nil {
		fmt.Printf("CreateDocProcessJobs error: %v\n", err)
		return
	}
	fmt.Printf("JobsDetail: %+v\n", res.JobsDetail)
}
