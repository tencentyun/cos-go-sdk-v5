package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:              "text/html",
			XCosServerSideEncryption: "AES256",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{},
	}
	name := "PutFromGoWithSSE-COS"
	content := "Put Object From Go With SSE-COS"
	f := strings.NewReader(content)
	_, err := c.Object.Put(context.Background(), name, f, opt)
	log_status(err)

	getopt := &cos.ObjectGetOptions{}
	var resp *cos.Response
	resp, err = c.Object.Get(context.Background(), name, getopt)
	log_status(err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyContent := string(bodyBytes)
	if bodyContent != content {
		log_status(errors.New("Content inconsistency"))
	}
}
