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
			SecretID:  os.Getenv("COS_SECRETID"),
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

	ctx := context.Background()

	// 由调用方创建 WorkerChannel 和 ResultChannel
	chjobs := make(chan *cos.Jobs, 100)
	chresults := make(chan *cos.Results, 10000)

	// 由调用方启动 worker 协程池，复用 cos.Worker 函数
	poolSize := 1
	for i := 0; i < poolSize; i++ {
		go cos.UploadWorker(ctx, c.Object, chjobs, chresults)
	}

	// 多线程上传对象，通过 WorkerChannel 和 ResultChannel 传入外部 channel
	opt := &cos.MultiUploadOptions{
		OptIni: &cos.InitiateMultipartUploadOptions{
			ACLHeaderOptions: nil,
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				XOptionHeader: &http.Header{
					"Content-Type": []string{"application/xml"},
				},
			},
		},
		WorkerChannel: chjobs,
		ResultChannel: chresults,
	}
	v, _, err := c.Object.Upload(ctx, "test", os.Args[1], opt)
	logStatus(err)
	fmt.Printf("Upload done, %v\n", v)

	// 上传完成后关闭 channel
	close(chjobs)
	close(chresults)
}
