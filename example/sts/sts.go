package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QcloudApi/qcloud_sign_golang"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"strings"
)

// Use Qcloud api github.com/QcloudApi/qcloud_sign_golang
// doc https://cloud.tencent.com/document/product/436/14048
type Credent struct {
	SessionToken string `json:"sessionToken"`
	TmpSecretID  string `json:"tmpSecretId"`
	TmpSecretKey string `json:"tmpSecretKey"`
}

// Data data in sts response body
type Data struct {
	Credentials Credent `json:"credentials"`
}

// Response sts response body
// In qcloud_sign_golang this response only return ak, sk and token
type Response struct {
	Dat Data `json:"data"`
}

func main() {
	// 替换实际的 SecretId 和 SecretKey
	secretID := "ak"
	secretKey := "sk"

	// 配置
	config := map[string]interface{}{"secretId": secretID, "secretKey": secretKey, "debug": false}

	// 请求参数
	params := map[string]interface{}{"Region": "gz", "Action": "GetFederationToken", "name": "alantong", "policy": "{\"statement\": [{\"action\": [\"name/cos:GetObject\",\"name/cos:PutObject\"],\"effect\": \"allow\",\"resource\":[\"qcs::cos:ap-guangzhou:uid/1253960454:prefix//1253960454/alangz/*\"]}],\"version\": \"2.0\"}"}

	// 发送请求
	retData, err := QcloudApi.SendRequest("sts", params, config)
	if err != nil {
		fmt.Print("Error.", err)
		return
	}
	r := &Response{}
	err = json.Unmarshal([]byte(retData), r)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取临时ak、sk、token
	tAk := r.Dat.Credentials.TmpSecretID
	tSk := r.Dat.Credentials.TmpSecretKey
	token := r.Dat.Credentials.SessionToken

	u, _ := url.Parse("https://alangz-1253960454.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
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

	name := "test/objectPut.go"
	f := strings.NewReader("test")

	_, err = c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}

	name = "test/put_option.go"
	f = strings.NewReader("test xxx")
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			//XCosACL: "public-read",
			XCosACL: "private",
		},
	}
	_, err = c.Object.Put(context.Background(), name, f, opt)
	if err != nil {
		panic(err)
	}
}
