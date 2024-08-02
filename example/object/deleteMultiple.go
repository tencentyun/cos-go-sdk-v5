package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"bytes"
	"io"

	"math/rand"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func logStatus(err error) {
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

func genBigData(blockSize int) []byte {
	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		logStatus(err)
	}
	return b
}

func uploadMulti(c *cos.Client) []string {
	names := []string{}
	data := genBigData(1024 * 1024 * 1)
	ctx := context.Background()
	var r io.Reader
	var name string
	n := 3

	for n > 0 {
		name = fmt.Sprintf("test/test_multi_delete_%s", time.Now().Format(time.RFC3339))
		r = bytes.NewReader(data)

		c.Object.Put(ctx, name, r, nil)
		names = append(names, name)
		n--
	}
	return names
}

func main() {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv("SECRETID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv("SECRETKEY"),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	ctx := context.Background()

	names := uploadMulti(c)
	names = append(names, []string{"a", "b", "c", "a+bc/xx&?+# "}...)
	obs := []cos.Object{}
	for _, v := range names {
		obs = append(obs, cos.Object{Key: v})
	}
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
		//Quiet: true,
	}

	v, _, err := c.Object.DeleteMulti(ctx, opt)
	logStatus(err)

	for _, x := range v.DeletedObjects {
		fmt.Printf("deleted %s\n", x.Key)
	}
	for _, x := range v.Errors {
		fmt.Printf("error %s, %s, %s\n", x.Key, x.Code, x.Message)
	}
}
