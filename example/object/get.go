package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

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

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	// Case1 通过resp.Body下载对象，Body需要关闭
	name := "test/example"
	resp, err := c.Object.Get(context.Background(), name, nil)
	log_status(err)

	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))

	// Case2 下载对象到文件. Body需要关闭
	fd, err := os.OpenFile("test", os.O_WRONLY|os.O_CREATE, 0660)
	log_status(err)

	defer fd.Close()
	resp, err = c.Object.Get(context.Background(), name, nil)
	log_status(err)

	io.Copy(fd, resp.Body)
	resp.Body.Close()

	// Case3 下载对象到文件
	_, err = c.Object.GetToFile(context.Background(), name, "test", nil)
	log_status(err)

	// Case4 range下载对象，可以根据range实现并发下载
	opt := &cos.ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}
	resp, err = c.Object.Get(context.Background(), name, opt)
	log_status(err)
	bs, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))

	// Case5 下载对象到文件，查看下载进度
	opt = &cos.ObjectGetOptions{
		Listener: &cos.DefaultProgressListener{},
	}
	_, err = c.Object.GetToFile(context.Background(), name, "test", opt)
	log_status(err)
}
