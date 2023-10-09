package main

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv("COS_SECRETID"),
			// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey:    os.Getenv("COS_SECRETKEY"),
			SessionToken: "<token>", // 请替换成您的临时密钥
		},
	})

	name := "exampleobject"
	ctx := context.Background()

	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL2(ctx, http.MethodPut, name, time.Hour, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("url: %v\n", presignedURL.String())
	// 通过预签名方式上传对象
	data := "test upload with presignedURL"
	f := strings.NewReader(data)
	req, err := http.NewRequest(http.MethodPut, presignedURL.String(), f)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
