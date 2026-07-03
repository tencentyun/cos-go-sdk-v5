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

// 示例1: Description 模式 + general 模板 + 多图（ObjectKey）
func describeWithLabels() {
	c := newClient()
	opt := &cos.CreateAIImageAnalysisOptions{
		Input: &cos.AIImageAnalysisInput{
			Message: &cos.AIImageAnalysisMessage{
				Content: &cos.AIImageAnalysisContent{
					Part: []cos.AIImageAnalysisPart{
						{Type: "Image", ObjectKey: "test/img1.jpg"},
						{Type: "Image", ObjectKey: "test/img2.jpg"},
					},
				},
			},
		},
		Conf: &cos.AIImageAnalysisConf{
			Type:         "Description",
			TemplateName: "general",
		},
	}
	res, _, err := c.CI.CreateAIImageAnalysis(context.Background(), opt)
	log_status(err)
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}

// 示例2: Description 模式 + ecommerce 模板（商品场景）
func describeEcommerce() {
	c := newClient()
	opt := &cos.CreateAIImageAnalysisOptions{
		Input: &cos.AIImageAnalysisInput{
			Message: &cos.AIImageAnalysisMessage{
				Content: &cos.AIImageAnalysisContent{
					Part: []cos.AIImageAnalysisPart{
						{Type: "Image", ObjectKey: "shop/a.jpg"},
						{Type: "Image", ObjectKey: "shop/b.jpg"},
					},
				},
			},
		},
		Conf: &cos.AIImageAnalysisConf{
			Type:         "Description",
			TemplateName: "ecommerce",
		},
	}
	res, _, err := c.CI.CreateAIImageAnalysis(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 示例3: Custom 模式 + 图文混排（Url + Prompt）
func describeCustom() {
	c := newClient()
	opt := &cos.CreateAIImageAnalysisOptions{
		Input: &cos.AIImageAnalysisInput{
			Message: &cos.AIImageAnalysisMessage{
				Content: &cos.AIImageAnalysisContent{
					Part: []cos.AIImageAnalysisPart{
						{Type: "Image", Url: "https://example.com/a.jpg"},
						{Type: "Image", ObjectKey: "test/b.jpg"},
						{Type: "Text", Text: "请用一句话总结这两张图的共同点"},
					},
				},
			},
		},
		Conf: &cos.AIImageAnalysisConf{Type: "Custom"},
	}
	res, _, err := c.CI.CreateAIImageAnalysis(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 示例4: 指定 AiModel
func describeWithCustomModel() {
	c := newClient()
	opt := &cos.CreateAIImageAnalysisOptions{
		Input: &cos.AIImageAnalysisInput{
			Message: &cos.AIImageAnalysisMessage{
				Content: &cos.AIImageAnalysisContent{
					Part: []cos.AIImageAnalysisPart{
						{Type: "Image", ObjectKey: "test/c.jpg"},
					},
				},
			},
		},
		Conf: &cos.AIImageAnalysisConf{
			Type:         "Description",
			TemplateName: "general",
			AiModel:      "qwen3.5-4b",
		},
	}
	res, _, err := c.CI.CreateAIImageAnalysis(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func main() {
	describeWithLabels()
	// describeEcommerce()
	// describeCustom()
	// describeWithCustomModel()
}
