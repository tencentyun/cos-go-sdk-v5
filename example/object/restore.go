package main

import (
	"context"
	// "os"

	"net/url"

	"net/http"

	// "github.com/tencentyun/cos-go-sdk-v5"
	// "github.com/tencentyun/cos-go-sdk-v5/debug"
	"github.com/toranger/cos-go-sdk-v5"
	"github.com/toranger/cos-go-sdk-v5/debug"
)

func main() {
	ak := "AKID4ygEetn1tZ6UingT44tU5smniXNEthIo"
	sk := "gHnDWp7NuLlBVmxAoRyPl0PoLrQqBMQK"
	u, _ := url.Parse("https://alangz-1253960454.cos.ap-guangzhou.myqcloud.com")
	// u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// SecretID:  os.Getenv("COS_SECRETID"),
			// SecretKey: os.Getenv("COS_SECRETKEY"),
			SecretID:  ak,
			SecretKey: sk,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	opt := &cos.ObjectRestoreOptions{
		Days: 2,
		Tier: &cos.CASJobParameters{
			// Standard, Exepdited and Bulk
			Tier: "Expedited",
		},
	}
	name := "archivetest"
	_, err := c.Object.PutRestore(context.Background(), name, opt)
	if err != nil {
		panic(err)
	}
}
