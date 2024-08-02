package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

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

func download(wg *sync.WaitGroup, c *cos.Client, keysCh <-chan []string) {
	defer wg.Done()
	for keys := range keysCh {
		key := keys[0]
		filename := keys[1]
		_, err := c.Object.GetToFile(context.Background(), key, filename, nil)
		if err != nil {
			logStatus(err)
		}
	}
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
				RequestHeader: false,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	// 多线程执行
	keysCh := make(chan []string, 3)
	var wg sync.WaitGroup
	threadpool := 3
	for i := 0; i < threadpool; i++ {
		wg.Add(1)
		go download(&wg, c, keysCh)
	}
	isTruncated := true
	prefix := "dir" // 下载 dir 目录下所有文件
	marker := ""
	localDir := "./local/"
	for isTruncated {
		opt := &cos.BucketGetOptions{
			Prefix:       prefix,
			Marker:       marker,
			EncodingType: "url", // url编码
		}
		// 列出目录
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, c := range v.Contents {
			key, _ := cos.DecodeURIComponent(c.Key) //EncodingType: "url"，先对 key 进行 url decode
			localfile := localDir + key
			if _, err := os.Stat(path.Dir(localfile)); err != nil && os.IsNotExist(err) {
				os.MkdirAll(path.Dir(localfile), os.ModePerm)
			}
			// 目录不需要下载
			if strings.HasSuffix(localfile, "/") {
				continue
			}
			keysCh <- []string{key, localfile}
		}
		marker = v.NextMarker
		isTruncated = v.IsTruncated
	}
	close(keysCh)
	wg.Wait()
}
