package main

import (
	"context"
	"fmt"
	"io"
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

// enhanceImageWhenDownload 图像增强, 下载时处理
func enhanceImageWhenDownload() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
	obj := "pic/heibai.jpeg"
	localPath := "test.jpeg"
	resp, err := c.CI.GetAIEnhanceImage(context.Background(), obj)
	log_status(err)
	if err == nil {
		fd, _ := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
		io.Copy(fd, resp.Body)
		fd.Close()
	}
}

// 图像增强, 上传时处理
func enhanceImageWhenUpload() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
	opt := &cos.ObjectPutOptions{
		nil,
		&cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "enhance.jpeg",
				Rule:   "ci-process=AIEnhanceImage",
			},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	name := "pic/enhance/heibai.jpeg"
	local_filename := "./heibai.jpeg"
	res, _, err := c.CI.PutFromFile(context.Background(), name, local_filename, opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.OriginalInfo)
	fmt.Printf("%+v\n", res.ProcessResults)
}

// 图像增强, 云上处理
func enhanceImageWhenCloud() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "enhance.jpeg",
				Rule:   "ci-process=AIEnhanceImage",
			},
		},
	}

	key := "pic/heibai.jpeg"
	res, _, err := c.CI.ImageProcess(context.Background(), key, pic)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func main() {
	// enhanceImageWhenDownload()
	// enhanceImageWhenUpload()
	enhanceImageWhenCloud()
}
