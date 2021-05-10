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
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
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

	opt := &cos.BucketPutPolicyOptions{
		Version: "2.0",
		Statement: []cos.BucketStatement{
			{
				Principal: map[string][]string{
					"qcs": []string{
						"qcs::cam::uin/100000000001:uin/100000000011", //替换成您想授予权限的账户uin
					},
				},
				Action: []string{
					"name/cos:GetObject",
				},
				Effect: "allow",
				Resource: []string{
					//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
					"qcs::cos:ap-guangzhou:uid/1259654469:test-1259654469/exampleobject",
				},
				Condition: map[string]map[string]interface{}{
					"ip_not_equal": map[string]interface{}{
						"qcs:ip": []string{
							"<ip>",
						},
					},
				},
			},
		},
	}

	_, err := c.Bucket.PutPolicy(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}
