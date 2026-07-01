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

// mustGetEnv 读取必填环境变量；缺失时打印友好提示并以非零码退出。
// demo 与 e2e 用例均通过此函数读取凭证，禁止硬编码。
func mustGetEnv(key, hint string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	fmt.Fprintf(os.Stderr, "\n[ENV MISSING] 缺少必需环境变量: %s\n", key)
	if hint != "" {
		fmt.Fprintf(os.Stderr, "              说明: %s\n", hint)
	}
	fmt.Fprintln(os.Stderr, "\n请在 shell 中导出后重试，例如：")
	fmt.Fprintln(os.Stderr, "  export COS_SECRETID=<你的 SecretId>")
	fmt.Fprintln(os.Stderr, "  export COS_SECRETKEY=<你的 SecretKey>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "如何获取密钥: https://console.cloud.tencent.com/cam/capi")
	os.Exit(1)
	return ""
}

// 该文件演示 HDR 图片处理服务开关接口（Bucket 级别）的标准用法。
// 参考文档: https://cloud.tencent.com/document/product/460/118210
//
// 三个接口：
//   - PutHDRImageProcessing    开通 HDR 图片处理（PUT /?hdr-image-processing）
//   - GetHDRImageProcessing    查询 HDR 图片处理状态（GET /?hdr-image-processing）
//   - DeleteHDRImageProcessing 关闭 HDR 图片处理（DELETE /?hdr-image-processing）
//
// 运行方式:
//   export COS_SECRETID=<你的 SecretId>
//   export COS_SECRETKEY=<你的 SecretKey>
//   go run hdr_image_processing.go

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

// newCIClient 构造一个指向 CI 域名的客户端，演示中统一使用该方法以便复用。
func newCIClient() *cos.Client {
	u, _ := url.Parse("https://test-1234567890.pic.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{CIURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  mustGetEnv("COS_SECRETID", "腾讯云访问密钥 SecretId"),
			SecretKey: mustGetEnv("COS_SECRETKEY", "腾讯云访问密钥 SecretKey"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
}

// putHDRImageProcessing 开通 HDR 图片处理服务。
// HdrMode 可选值：
//   - "API"      : 仅通过 API 调用支持对 HDR 图片进行处理
//   - "Auto"     : 无需携带 HDR 参数，使用万象基础处理参数时自动支持 HDR
//   - "Auto,API" : 同时启用以上两种模式（推荐）
func putHDRImageProcessing() {
	c := newCIClient()
	opt := &cos.HDRImageProcessingOptions{
		HdrMode: "Auto,API",
	}
	_, err := c.CI.PutHDRImageProcessing(context.Background(), opt)
	log_status(err)
}

// getHDRImageProcessing 查询当前 HDR 图片处理服务状态。
// 返回 Status 为 "on" / "off"，HdrMode 与开通时保持一致。
func getHDRImageProcessing() {
	c := newCIClient()
	res, _, err := c.CI.GetHDRImageProcessing(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// deleteHDRImageProcessing 关闭 HDR 图片处理服务。
func deleteHDRImageProcessing() {
	c := newCIClient()
	_, err := c.CI.DeleteHDRImageProcessing(context.Background())
	log_status(err)
}

func main() {
	putHDRImageProcessing()
	getHDRImageProcessing()
	// deleteHDRImageProcessing()
	// getHDRImageProcessing()
}
