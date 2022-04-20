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

// DescribeTemplate 搜索模板
func DescribeTemplate() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	opt := &cos.DescribeMediaTemplateOptions{
		Tag:        "Transcode",
		PageNumber: 1,
		PageSize:   5,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// DeleteTemplate 删除模板
func DeleteTemplate() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	DescribeTemplateRes, _, err := c.CI.DeleteMediaTemplate(context.Background(), "t11c1b0a3fb304463096e828a40a013579")
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateTranscodeTemplate TODO
func CreateTranscodeTemplate() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// CreateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
		Tag:  "Transcode",
		Name: "transtpl-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
		},
		Audio: &cos.Audio{
			Codec: "aac",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaTranscodeTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateTranscodeTemplate TODO
func UpdateTranscodeTemplate() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// CreateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
		Tag:  "Transcode",
		Name: "transtpl-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
			Crf:   "26",
		},
		Audio: &cos.Audio{
			Codec: "aac",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	templateId := "t1d31d58d8a4204d7396087f56a448abd5"
	createTplRes, _, err := c.CI.UpdateMediaTranscodeTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)

}

func main() {
	// DescribeTemplate()
	// DeleteTemplate()
	// CreateTranscodeTemplate()
	// UpdateTranscodeTemplate()
}
