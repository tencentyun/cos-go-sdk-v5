package main

import (
	"context"
	"fmt"
	"os"

	"net/url"
	"strings"

	"net/http"

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

func initUpload(c *cos.Client, name string) *cos.InitiateMultipartUploadResult {
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	log_status(err)
	fmt.Printf("%#v\n", v)
	return v
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
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	optcom := &cos.CompleteMultipartUploadOptions{}
	name := "test/test_multi_upload.go"
	up := initUpload(c, name)
	uploadID := up.UploadID

	fd, err := os.Open("test")
	if err != nil {
		fmt.Printf("Open File Error: %v\n", err)
		return
	}
	defer fd.Close()
	stat, err := fd.Stat()
	if err != nil {
		fmt.Printf("Stat File Error: %v\n", err)
		return
	}
	opt := &cos.ObjectUploadPartOptions{
		Listener:      &cos.DefaultProgressListener{},
		ContentLength: int(stat.Size()),
	}
	resp, err := c.Object.UploadPart(
		context.Background(), name, uploadID, 1, fd, opt,
	)
	optcom.Parts = append(optcom.Parts, cos.Object{
		PartNumber: 1, ETag: resp.Header.Get("ETag"),
	})
	log_status(err)

	f := strings.NewReader("test heoo")
	resp, err = c.Object.UploadPart(
		context.Background(), name, uploadID, 2, f, nil,
	)
	log_status(err)
	optcom.Parts = append(optcom.Parts, cos.Object{
		PartNumber: 2, ETag: resp.Header.Get("ETag"),
	})

	_, _, err = c.Object.CompleteMultipartUpload(context.Background(), name, uploadID, optcom)
	log_status(err)
}
