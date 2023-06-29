package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	uin := "100010805041"
	appid := 1259654469
	jobid := "49e0dd01-27a6-41a6-97b2-dda3cca19223"
	u, _ := url.Parse("https://" + uin + ".cos-control.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BatchURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}

	_, err := c.Batch.DeleteJob(context.Background(), jobid, headers)
	if err != nil {
		panic(err)
	}
}
