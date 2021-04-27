package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
	// 创建文件夹
	name := "example/"
	_, err := c.Object.Put(context.Background(), name, strings.NewReader(""), nil)
	log_status(err)

	// 查看文件夹是否存在
	_, err = c.Object.Head(context.Background(), name, nil)
	log_status(err)

	// 删除文件夹
	_, err = c.Object.Delete(context.Background(), name)
	log_status(err)

	// 上传到虚拟目录
	dir := "exampledir/"
	filename := "exampleobject"
	key := dir + filename
	f := strings.NewReader("test file")
	_, err = c.Object.Put(context.Background(), key, f, nil)
	log_status(err)

	// 删除文件夹内所有文件
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:  dir,
		MaxKeys: 1000,
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			log_status(err)
			break
		}
		for _, content := range v.Contents {
			_, err = c.Object.Delete(context.Background(), content.Key)
			if err != nil {
				log_status(err)
			}
		}
		isTruncated = v.IsTruncated
		marker = v.NextMarker
	}
}
