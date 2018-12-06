package main

import (
	"context"
	"net/url"
	"os"

	"net/http"

	"fmt"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_Key"),
			SecretKey: os.Getenv("COS_Secret"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/uploadFile.go"
	f, err := os.Open(os.Args[0])
	if err != nil {
		panic(err)
	}
	s, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Size())
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: int(s.Size()),
		},
	}
	//opt.ContentLength = int(s.Size())

	_, err = c.Object.Put(context.Background(), name, f, opt)
	if err != nil {
		panic(err)
	}
}
