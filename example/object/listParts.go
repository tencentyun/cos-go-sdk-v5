package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
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

func uploadPart(c *cos.Client, name string, uploadID string, blockSize, n int) string {

	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		log_status(err)
	}
	s := fmt.Sprintf("%X", b)
	f := strings.NewReader(s)

	resp, err := c.Object.UploadPart(
		context.Background(), name, uploadID, n, f, nil,
	)
	log_status(err)
	fmt.Printf("%s\n", resp.Status)
	return resp.Header.Get("Etag")
}

func main() {
	u, _ := url.Parse("http://alangz-1251668577.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/test_list_parts.go"
	up := initUpload(c, name)
	uploadID := up.UploadID
	ctx := context.Background()
	blockSize := 1024 * 1024 * 3

	for i := 1; i < 5; i++ {
		uploadPart(c, name, uploadID, blockSize, i)
	}

	// opt := &cos.ObjectListPartsOptions{
	// 	MaxParts: "1",
	// }
	v, _, err := c.Object.ListParts(ctx, name, uploadID, nil)
	if err != nil {
		log_status(err)
		return
	}
	for _, p := range v.Parts {
		fmt.Printf("%d, %s, %d\n", p.PartNumber, p.ETag, p.Size)
	}
	fmt.Printf("%s\n", v.Initiator.ID)
	fmt.Printf("%s\n", v.Owner.ID)
}
