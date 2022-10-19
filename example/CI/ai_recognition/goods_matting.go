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

// 商品抠图(上传时)
func goodsMattingWhenUpload() {
	u, _ := url.Parse("https://test-1253960454.cos.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
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
				FileId: "pic/out.jpeg",
				Rule:   "ci-process=GoodsMatting&center-layout=1&padding-layout=500x500",
			},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	name := "pic/test/cup.jpeg"
	local_filename := "./cup.jpeg"
	res, _, err := c.CI.PutFromFile(context.Background(), name, local_filename, opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.OriginalInfo)
	fmt.Printf("%+v\n", res.ProcessResults)
}

// 商品抠图(下载时)
func goodsMattingWhenDownload() {
	u, _ := url.Parse("https://test-1253960454.pic.ap-chongqing.myqcloud.com")
	// u, _ := url.Parse("https://test-1253960454.cos.ap-chongqing.myqcloud.com")
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

	key := "pic/cup.jpeg"
	localPath := "cup.jpeg"
	opt := &cos.GoodsMattingptions{
		CenterLayout:  "1",
		PaddingLayout: "500x500",
	}
	resp, err := c.CI.GoodsMattingWithOpt(context.Background(), key, opt)
	log_status(err)
	fd, _ := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	io.Copy(fd, resp.Body)
	fd.Close()
}

// 商品抠图(云上数据处理)
func goodsMattingWhenCloud() {
	u, _ := url.Parse("https://test-1253960454.cos.ap-chongqing.myqcloud.com")
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

	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "pic/cup_out.jpeg",
				Rule:   "ci-process=GoodsMatting&center-layout=1&padding-layout=500x500",
			},
		},
	}

	key := "pic/cup.jpeg"
	res, _, err := c.CI.ImageProcess(context.Background(), key, pic)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func main() {
	goodsMattingWhenDownload()
	// goodsMattingWhenUpload()
	// goodsMattingWhenCloud()
}
