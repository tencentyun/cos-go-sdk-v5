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

func logStatus(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
	} else {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func main() {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台"存储桶概览"查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://cd-1259654469.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})

	ctx := context.Background()

	// 由调用方创建 WorkerChannel 和 ResultChannel
	chjobs := make(chan *cos.Jobs, 100)
	chresults := make(chan *cos.Results, 10000)

	// 由调用方启动 downloadWorker 协程池，复用 cos.DownloadWorker 函数
	poolSize := 5
	for i := 0; i < poolSize; i++ {
		go cos.DownloadWorker(ctx, c.Object, chjobs, chresults)
	}

	// 多线程下载对象，通过 WorkerChannel 和 ResultChannel 传入外部 channel
	opt := &cos.MultiDownloadOptions{
		WorkerChannel: chjobs,
		ResultChannel: chresults,
	}
	resp, err := c.Object.Download(ctx, "test", "./test1G", opt)
	logStatus(err)
	fmt.Printf("done, %v\n", resp.Header)

	// 下载完成后关闭 channel
	close(chjobs)
	close(chresults)
}
