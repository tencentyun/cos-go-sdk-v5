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

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	secretID := os.Getenv("COS_SECRETID")
	secretKey := os.Getenv("COS_SECRETKEY")
	if secretID == "" || secretKey == "" {
		panic("COS_SECRETID or COS_SECRETKEY is invalid")
	}

	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	// Case1 Download object into ReadCloser(). the body needs to be closed
	name := "test1.txt"
	resp, err := c.Object.Get(context.Background(), name, nil)
	panicError(err)

	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))

	// Case2 Download object to local file. the body needs to be closed
	fd, err := os.OpenFile("hello.txt", os.O_WRONLY|os.O_CREATE, 0660)
	panicError(err)

	defer fd.Close()
	resp, err = c.Object.Get(context.Background(), name, nil)
	panicError(err)
	io.Copy(fd, resp.Body)
	resp.Body.Close()

	// Case3 Download object to local file path
	_, err = c.Object.GetToFile(context.Background(), name, "hello_1.txt", nil)
	panicError(err)

	// Case4 Download object with range header, can used to concurrent download
	opt := &cos.ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}
	resp, err = c.Object.Get(context.Background(), name, opt)
	panicError(err)
	bs, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
