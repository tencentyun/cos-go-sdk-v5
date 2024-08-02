package main

import (
	"context"
	"fmt"
	"os"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func logStatus(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("Code: %v\n", e.Code)
		fmt.Printf("Message: %v\n", e.Message)
		fmt.Printf("Resource: %v\n", e.Resource)
		fmt.Printf("RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Println(err)
		// ERROR
	}
}

func main() {
	c := cos.NewClient(nil, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	s, _, err := c.Service.Get(context.Background())
	if err != nil {
		logStatus(err)
	}

	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
}
