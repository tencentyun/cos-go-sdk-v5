package main

import (
	"context"
	"net/url"
	"os"
	"time"
	"net/http"

	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	u, _ := url.Parse("http://tencentyun02-1252448703.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_Key"),
			SecretKey: os.Getenv("COS_Secret"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	f,err:=os.Open("E:/cos-php-sdk.zip")
	if err!=nil {panic(err)}
	opt := &cos.MultiUploadOptions{
		OptIni: nil,
		PartSize:1,
	}
	v, _, err := c.Object.MultiUpload(
		context.Background(), "test", f, opt,
	)
	if err!=nil {panic(err)}
	fmt.Println(v)
}