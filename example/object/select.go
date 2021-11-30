package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"io/ioutil"
	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID:  os.Getenv("COS_SECRETID"),
			// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
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

	opt := &cos.ObjectSelectOptions{
		Expression:     "Select * from COSObject",
		ExpressionType: "SQL",
		InputSerialization: &cos.SelectInputSerialization{
			JSON: &cos.JSONInputSerialization{
				Type: "DOCUMENT",
			},
		},
		OutputSerialization: &cos.SelectOutputSerialization{
			JSON: &cos.JSONOutputSerialization{
				RecordDelimiter: "\n",
			},
		},
		RequestProgress: "TRUE",
	}
	res, err := c.Object.Select(context.Background(), "test.json", opt)
	if err != nil {
		panic(err)
	}
	defer res.Close()
	data, err := ioutil.ReadAll(res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("data: %v\n", string(data))
	resp, _ := res.(*cos.ObjectSelectResponse)
	fmt.Printf("data: %+v\n", resp.Frame)

	// Select to File
	_, err = c.Object.SelectToFile(context.Background(), "test.json", "./test.json", opt)
	if err != nil {
		panic(err)
	}
}
