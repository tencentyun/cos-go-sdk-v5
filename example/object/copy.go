package main

import (
	"context"
	"net/url"
	"os"
	"strings"

	"net/http"

	"fmt"
	"io/ioutil"
	"time"

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
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	source := "test/objectMove1.go"
	expected := "test"
	f := strings.NewReader(expected)

	_, err := c.Object.Put(context.Background(), source, f, nil)
	log_status(err)

	soruceURL := fmt.Sprintf("%s/%s", u.Host, source)
	dest := fmt.Sprintf("test/objectMove_%d.go", time.Now().Nanosecond())
	//opt := &cos.ObjectCopyOptions{}
	res, _, err := c.Object.Copy(context.Background(), dest, soruceURL, nil)
	log_status(err)
	fmt.Printf("%+v\n\n", res)

	resp, err := c.Object.Get(context.Background(), dest, nil)
	log_status(err)

    bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	result := string(bs)
	if result != expected {
		panic(fmt.Sprintf("%s != %s", result, expected))
	}
}
