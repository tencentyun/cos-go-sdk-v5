package main

import (
	"bytes"
	"context"
	"encoding/base64"
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

func main() {
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
	obj := "pic/self.jpeg"
	opt := &cos.FaceEffectOptions{
		Type:      "face-segmentation",
		Whitening: 50,
	}
	res, _, err := c.CI.FaceEffect(context.Background(), obj, opt)
	log_status(err)

	if len(res.ResultImage) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.ResultImage)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("result_image.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}

	if len(res.ResultMask) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.ResultMask)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("result_mask.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
}
