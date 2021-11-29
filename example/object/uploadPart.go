package main

import (
	"context"
	"fmt"
	"os"

	"net/url"
	"strings"

	"net/http"

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

func initUpload(c *cos.Client, name string) *cos.InitiateMultipartUploadResult {
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	log_status(err)
	fmt.Printf("%#v\n", v)
	return v
}

func main() {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID:  os.Getenv("COS_SECRETID"),
			// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv("COS_SECRETKEY"),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	optcom := &cos.CompleteMultipartUploadOptions{}
	name := "test/test_multi_upload.go"
	up := initUpload(c, name)
	uploadID := up.UploadID

	fd, err := os.Open("test")
	if err != nil {
		fmt.Printf("Open File Error: %v\n", err)
		return
	}
	defer fd.Close()
	stat, err := fd.Stat()
	if err != nil {
		fmt.Printf("Stat File Error: %v\n", err)
		return
	}
	opt := &cos.ObjectUploadPartOptions{
		Listener:      &cos.DefaultProgressListener{},
		ContentLength: stat.Size(),
	}
	resp, err := c.Object.UploadPart(
		context.Background(), name, uploadID, 1, fd, opt,
	)
	log_status(err)
	optcom.Parts = append(optcom.Parts, cos.Object{
		PartNumber: 1, ETag: resp.Header.Get("ETag"),
	})

	f := strings.NewReader("test heoo")
	resp, err = c.Object.UploadPart(
		context.Background(), name, uploadID, 2, f, nil,
	)
	log_status(err)
	optcom.Parts = append(optcom.Parts, cos.Object{
		PartNumber: 2, ETag: resp.Header.Get("ETag"),
	})

	_, _, err = c.Object.CompleteMultipartUpload(context.Background(), name, uploadID, optcom)
	log_status(err)
}
