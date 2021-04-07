package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func upload(wg *sync.WaitGroup, c *cos.Client, keysCh <-chan string) {
	defer wg.Done()
	for key := range keysCh {
		// 下载文件到当前目录
		_, filename := filepath.Split(key)
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
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
	keysCh := make(chan string, 2)
	keys := []string{"test/test1", "test/test2", "test/test3"}
	var wg sync.WaitGroup
	threadpool := 2
	for i := 0; i < threadpool; i++ {
		wg.Add(1)
		go upload(&wg, c, keysCh)
	}
	for _, key := range keys {
		keysCh <- key
	}
	close(keysCh)
	wg.Wait()
}
