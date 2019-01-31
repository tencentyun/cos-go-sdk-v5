package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"io"
	"io/ioutil"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
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

	// Case1 Download object into ReadCloser(). the body needs to be closed
	name := "test/hello.txt"
	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))

	// Case2 Download object to local file. the body needs to be closed
	fd, err := os.OpenFile("hello.txt", os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	resp, err = c.Object.Get(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	io.Copy(fd, resp.Body)
	resp.Body.Close()

	// Case3 Download object to local file path
	err = c.Object.GetToFile(context.Background(), name, "hello_1.txt", nil)
	if err != nil {
		panic(err)
	}

	// Case4 Download object with range header, can used to concurrent download
	opt := &cos.ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}
	resp, err = c.Object.Get(context.Background(), name, opt)
	if err != nil {
		panic(err)
	}
	bs, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
