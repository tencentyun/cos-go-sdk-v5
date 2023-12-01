package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

func getClient() *cos.Client {
	u, _ := url.Parse("https://testpic-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://testpic-1253960454.ci.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
	return c
}

func CreatePlayKey() {
	c := getClient()
	Res, _, err := c.CI.CreatePlayKey(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", Res.PlayKeyList)
}

func GetPlayKey() {
	c := getClient()

	Res, _, err := c.CI.GetPlayKey(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", Res.PlayKeyList)

}

func main() {
	CreatePlayKey()
	GetPlayKey()
}
