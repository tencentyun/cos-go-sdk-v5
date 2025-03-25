package main

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv("SECRETID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey:    os.Getenv("SECRETKEY"),
			SessionToken: "<token>", // 请替换成您的临时密钥
		},
	})

	name := "exampleobject"
	ctx := context.Background()

	// 获取预签名
	// http Method需要和实际http请求一致，如PUT请求设置成http.MethodPut，GET请求设置成http.MethodGet
	presignedURL, err := c.Object.GetPresignedURL2(ctx, http.MethodPut, name, time.Hour, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("url: %v\n", presignedURL.String())
}
