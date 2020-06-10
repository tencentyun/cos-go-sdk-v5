package main

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	appid := "1259654469"
	bucket := "test-1259654469"
	c := sts.NewClient(
		os.Getenv("COS_SECRETID"),
		os.Getenv("COS_SECRETKEY"),
		nil,
	)
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          "ap-guangzhou",
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
						"name/cos:GetObject",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:ap-guangzhou:uid/" + appid + ":" + bucket + "/exampleobject",
					},
				},
			},
		},
	}
	res, err := c.GetCredential(opt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res.Credentials)

	//获取临时ak、sk、token
	tAk := res.Credentials.TmpSecretID
	tSk := res.Credentials.TmpSecretKey
	token := res.Credentials.SessionToken

	u, _ := url.Parse("https://" + bucket + ".cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 使用临时密钥
			SecretID:     tAk,
			SecretKey:    tSk,
			SessionToken: token,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "exampleobject"
	f := strings.NewReader("test")

	_, err = client.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}

	name = "exampleobject"
	f = strings.NewReader("test xxx")
	optc := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			//XCosACL: "public-read",
			XCosACL: "private",
		},
	}
	_, err = client.Object.Put(context.Background(), name, f, optc)
	if err != nil {
		panic(err)
	}

}
