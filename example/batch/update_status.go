package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	uin := "100010805041"
	appid := 1259654469
	jobid := "289b0ea1-5ac5-453d-8a61-7f452dd4a209"
	u, _ := url.Parse("https://" + uin + ".cos-control.ap-chengdu.myqcloud.com")
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

	opt := &cos.BatchUpdateStatusOptions{
		JobId:              jobid,
		RequestedJobStatus: "Ready", // 允许状态转换见 https://cloud.tencent.com/document/product/436/38604
		StatusUpdateReason: "to test",
	}
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}

	res, _, err := c.Batch.UpdateJobStatus(context.Background(), opt, headers)
	if err != nil {
		panic(err)
	}
	if res != nil {
		fmt.Printf("%+v", res)
	}
}
