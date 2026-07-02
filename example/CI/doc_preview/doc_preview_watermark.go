//go:build example_doc_preview_watermark
// +build example_doc_preview_watermark

// Example: 同步给 PDF 加平铺水印（复用 DocPreview 接口，通过 DstType="watermark" 触发水印场景）
//
// 运行:
//
//	go run -tags=example_doc_preview_watermark example/CI/doc_preview/doc_preview_watermark.go
//
// 依赖环境变量: SECRETID / SECRETKEY / BUCKET / REGION
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	bucket := os.Getenv("BUCKET")
	region := os.Getenv("REGION")
	if bucket == "" || region == "" {
		fmt.Println("please set env BUCKET / REGION")
		return
	}
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
		},
	})

	resp, err := c.CI.DocPreview(context.Background(), "sample.pdf", &cos.DocPreviewOptions{
		SrcType:           "pdf",
		DstType:           "watermark", // 触发 PDF 加水印场景
		Type:              "Text",
		Text:              "confidential",
		Batch:             1,
		HorizontalSpacing: 200,
		VerticalSpacing:   150,
		Page:              1,
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status:", resp.StatusCode)
	out, _ := os.Create("watermark.jpg")
	defer out.Close()
	n, _ := io.Copy(out, resp.Body)
	fmt.Println("saved:", n, "bytes -> watermark.jpg")
}
