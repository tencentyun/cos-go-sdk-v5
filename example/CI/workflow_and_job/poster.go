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

func PutPosterproductionTemplate() {
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
	PosterproductionTemplate := &cos.PosterproductionTemplateOptions{
		Input: &cos.PosterproductionInput{
			Object: "input/sample.psd",
		},
		Name: "test",
	}
	PutPosterproductionRes, _, err := c.CI.PutPosterproductionTemplate(context.Background(), PosterproductionTemplate)
	log_status(err)
	fmt.Printf("%+v\n", PutPosterproductionRes)
	fmt.Printf("%+v\n", &PutPosterproductionRes.Template)
}

func GetPosterproductionTemplate() {
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
	PutPosterproductionRes, _, err := c.CI.GetPosterproductionTemplate(context.Background(), "6444f12ae24d596cdbd774fb")
	log_status(err)
	fmt.Printf("%+v\n", PutPosterproductionRes)
	fmt.Printf("%+v\n", &PutPosterproductionRes.Template)
}

func GetPosterproductionTemplates() {
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
	opt := &cos.DescribePosterproductionTemplateOptions{
		PageNumber: 1,
		PageSize:   10,
	}
	PutPosterproductionRes, _, err := c.CI.GetPosterproductionTemplates(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", PutPosterproductionRes)
	fmt.Printf("%+v\n", &PutPosterproductionRes.TemplateList)
}

func InvokePosterProductionJob() {
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
	type autoInfo struct {
		TextMain string `xml:"text_main,omitempty"`
		TextSub  string `xml:"text_sub,omitempty"`
	}
	info := autoInfo{
		TextMain: "父亲节快乐",
		TextSub:  "献给最伟大的父亲!!!",
	}
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "PosterProduction",
		Operation: &cos.MediaProcessJobOperation{
			PosterProduction: &cos.PosterProduction{
				TemplateId: "6444f12ae24d596cdbd774fb",
				Info:       info,
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "poster/PosterProduction2.jpg",
			},
		},
		// todo 需要替换为自己的回调地址信息
		CallBack: "https://demo.org/callback",
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

func main() {
}
