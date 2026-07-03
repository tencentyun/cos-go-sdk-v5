//go:build example_doc_preview_ofd
// +build example_doc_preview_ofd

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// OFD → JPG 预览示例
// 运行: go run -tags=example_doc_preview_ofd example/CI/doc_preview/doc_preview_ofd.go
func main() {
	u, _ := url.Parse("https://xxx-125xxxxxxxx.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
		},
	})
	resp, err := c.CI.DocPreview(context.Background(), "sample.ofd", &cos.DocPreviewOptions{
		SrcType: "ofd",
		DstType: "jpg",
		Page:    1,
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile("preview.jpg", data, 0644)
	fmt.Println("done")
}
