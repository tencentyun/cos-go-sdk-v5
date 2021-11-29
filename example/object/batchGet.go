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

func download(wg *sync.WaitGroup, c *cos.Client, keysCh <-chan []string) {
	defer wg.Done()
	for keys := range keysCh {
		key := keys[0]
		filename := keys[1]
		_, err := c.Object.GetToFile(context.Background(), key, filename, nil)
		if err != nil {
			log_status(err)
		}
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
