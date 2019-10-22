package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/QcloudApi/qcloud_sign_golang"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

// Use Qcloud api github.com/QcloudApi/qcloud_sign_golang
// doc https://cloud.tencent.com/document/product/436/14048
type Credent struct {
	SessionToken string `json:"sessionToken"`
	TmpSecretID  string `json:"tmpSecretId"`
	TmpSecretKey string `json:"tmpSecretKey"`
}

type PolicyStatement struct {
	Action    []string                          `json:"action,omitempty"`
	Effect    string                            `json:"effect,omitempty"`
	Resource  []string                          `json:"resource,omitempty"`
	Condition map[string]map[string]interface{} `json:"condition,omitempty"`
}

type CAMPolicy struct {
	Statement []PolicyStatement   `json:"statement,omitempty"`
	Version   string              `json:"version,omitempty"`
	Principal map[string][]string `json:"principal,omitempty"`
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
	// 在环境变量中设置您的 SecretId 和 SecretKey
	secretID := os.Getenv("COS_SECRETID")
	secretKey := os.Getenv("COS_SECRETKEY")
	appid := "1259654469"       //替换成您的APPID
	bucket := "test-1259654469" //替换成您的bucket，格式：<bucketname-APPID>

	// 配置
	config := map[string]interface{}{"secretId": secretID, "secretKey": secretKey, "debug": false}

	policy := &CAMPolicy{
		Statement: []PolicyStatement{
			PolicyStatement{
				Action: []string{
					"name/cos:PostObject",
					"name/cos:PutObject",
				},
				Effect: "allow",
				Resource: []string{
                    //这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
					"qcs::cos:ap-guangzhou:uid/" + appid + ":" + bucket + "/exampleobject",
				},
			},
		},
		Version: "2.0",
	}
	bPolicy, err := json.Marshal(policy)
	if err != nil {
		fmt.Print("Error.", err)
		return
	}
	policyStr := string(bPolicy)
	// 请求参数
	params := map[string]interface{}{
		"Region": "gz",
		"Action": "GetFederationToken",
		"name":   "test",
		"policy": policyStr,
	}
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

	u, _ := url.Parse("https://" + bucket + ".cos.ap-guangzhou.myqcloud.com")
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

	name := "exampleobject"
	f := strings.NewReader("test")

	_, err = c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}

	name = "exampleobject"
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
