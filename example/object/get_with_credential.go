package main

import (
	"context"
	"fmt"
	"net/url"

	"net/http"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func log_status(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Printf("ERROR: %v\n", err)
		// ERROR
	}
}

type Credential struct {
}

// 需实现 CredentialIface 三个方法
func (c *Credential) GetSecretId() string {
	return os.Getenv("COS_SECRETID")
}

func (c *Credential) GetSecretKey() string {
	return os.Getenv("COS_SECRETKEY")
}

func (c *Credential) GetToken() string {
	return ""
}

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		// 使用 CredentialsTransport
		Transport: &cos.CredentialTransport{
			// 通过 CredentialIface 获取密钥, 需实现 GetSecretKey，GetSecretId，GetToken 方法。
			Credential: &Credential{},
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	name := "exampleobject"
	_, err := c.Object.Get(context.Background(), name, nil)
	log_status(err)
}
