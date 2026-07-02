//go:build example_pm3u8_user_data
// +build example_pm3u8_user_data

package main

// 演示：使用 EncodeUserData helper + GetPrivateM3U8 生成一个带 user-data 的
// 私有 m3u8 预签名 URL，返回给前端 / 播放器直接播放。
//
// 相比 GetPrivateM3U8 直接下载 m3u8 内容，"输出签名 URL" 是更常见的用法：
// 后端只负责签发 URL 并返回给业务方，播放器 / CDN 会自己去 GET，
// 服务端会把 user-data 透传到 m3u8 内每个 ts 分片 URL 中，方便鉴权 / 审计追踪。
//
// 运行方式（单文件独立构建，避免与目录内其他 main 冲突）：
//   go run -tags example_pm3u8_user_data ./example/CI/media_process/pm3u8_user_data.go
//
// 前置环境变量：
//   COS_SECRETID / COS_SECRETKEY  —— 腾讯云 API 密钥
//   COS_BUCKET_URL                —— 形如 https://<bucket>-<appid>.cos.<region>.myqcloud.com
//   COS_OBJECT_KEY                —— 已完成边转边播的 m3u8 对象 key

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	bucketURL, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	b := &cos.BaseURL{BucketURL: bucketURL}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	key := os.Getenv("COS_OBJECT_KEY")

	// 1) 使用 helper 生成 user-data
	userData, err := cos.EncodeUserData(map[string]string{
		"uid":     "u_10001",
		"traceid": "trace-abc",
	})
	if err != nil {
		fmt.Printf("EncodeUserData err: %v\n", err)
		return
	}
	fmt.Printf("user-data = %s\n", userData)

	// 2) 拼接 ci-process=pm3u8 相关的 query 参数
	q := &url.Values{}
	q.Set("ci-process", "pm3u8")
	q.Set("expires", "3600")
	q.Set("user-data", userData)

	// 3) 生成 GET 请求的预签名 URL —— 这才是最常见的用法：
	//    把签名 URL 返回给客户端 / 播放器，由它们自己拉 m3u8 与 ts 分片
	presigned, err := c.Object.GetPresignedURL2(
		context.Background(),
		http.MethodGet,
		key,
		15*time.Minute, // URL 有效期
		&cos.PresignedURLOptions{
			Query: q,
		},
	)
	if err != nil {
		fmt.Printf("GetPresignedURL2 err: %v\n", err)
		return
	}
	fmt.Printf("presigned pm3u8 url = %s\n", presigned.String())

	// （可选）如果确实需要在服务端直接拿到 m3u8 内容，
	// 可以继续走 GetPrivateM3U8：
	//   resp, _ := c.CI.GetPrivateM3U8(ctx, key,
	//       &cos.GetPrivateM3U8Options{Expires: 3600, UserData: userData})
}
