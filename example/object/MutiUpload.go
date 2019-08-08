package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	u, _ := url.Parse("http://alanbj-1251668577.cos.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})

	opt := &cos.MultiUploadOptions{
		OptIni:   nil,
		PartSize: 1,
	}
	v, _, err := c.Object.MultiUpload(
		context.Background(), "test/gomulput1G", "./test1G", opt,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
