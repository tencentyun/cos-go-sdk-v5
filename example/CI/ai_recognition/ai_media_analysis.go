package main

import (
	"context"
	"encoding/json"
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

// newClient 同时配置 BucketURL 和 CIURL（本接口走 CI 域名）。
func newClient() *cos.Client {
	u, _ := url.Parse("https://test-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-1253960454.ci.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
}

// 示例1: Description 模式 + FTM 模板（影视传媒）
func descriptionFTM() {
	c := newClient()
	opt := &cos.CreateAIMediaAnalysisOptions{
		Input: &cos.AIMediaAnalysisInput{
			Message: &cos.AIMediaAnalysisMessage{
				Content: &cos.AIMediaAnalysisContent{
					Part: []cos.AIMediaAnalysisPart{
						{Type: "Video", ObjectKey: "media/drama_clip.mp4", Fps: 1},
					},
				},
			},
		},
		Conf: &cos.AIMediaAnalysisConf{
			Type:         "Description",
			TemplateName: "FTM",
		},
	}
	res, _, err := c.CI.CreateAIMediaAnalysis(context.Background(), opt)
	log_status(err)
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}

// 示例2: Description 模式 + 默认模板 + Url 方式
func descriptionDefault() {
	c := newClient()
	opt := &cos.CreateAIMediaAnalysisOptions{
		Input: &cos.AIMediaAnalysisInput{
			Message: &cos.AIMediaAnalysisMessage{
				Content: &cos.AIMediaAnalysisContent{
					Part: []cos.AIMediaAnalysisPart{
						{Type: "Video", Url: "https://example.com/video.mp4", Fps: 1},
					},
				},
			},
		},
		Conf: &cos.AIMediaAnalysisConf{
			Type: "Description",
		},
	}
	res, _, err := c.CI.CreateAIMediaAnalysis(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 示例3: Custom 模式 + 自定义 Prompt
func customAnalysis() {
	c := newClient()
	opt := &cos.CreateAIMediaAnalysisOptions{
		Input: &cos.AIMediaAnalysisInput{
			Message: &cos.AIMediaAnalysisMessage{
				Content: &cos.AIMediaAnalysisContent{
					Part: []cos.AIMediaAnalysisPart{
						{Type: "Video", ObjectKey: "warehouse/monitor.mp4", Fps: 0.5},
						{Type: "Text", Text: "请统计视频中出现的人物数量，描述他们的着装与主要活动，并指出是否存在异常行为。"},
					},
				},
			},
		},
		Conf: &cos.AIMediaAnalysisConf{Type: "Custom"},
	}
	res, _, err := c.CI.CreateAIMediaAnalysis(context.Background(), opt)
	log_status(err)
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}

func main() {
	descriptionFTM()
	// descriptionDefault()
	// customAnalysis()
}
