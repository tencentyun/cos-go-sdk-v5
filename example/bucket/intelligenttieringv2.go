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
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
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
	id := "test"
	opt := &cos.BucketPutIntelligentTieringOptions{
		Id:     id,
		Status: "Enabled",
		Tiering: []*cos.BucketIntelligentTieringTransition{
			{
				Days:       91,
				AccessTier: "ARCHIVE_ACCESS",
			},
		},
		Filter: &cos.BucketIntelligentTieringFilter{
			And: &cos.BucketIntelligentTieringFilterAnd{
				Prefix: "test",
				Tag: []*cos.BucketTaggingTag{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		},
	}
	_, err := c.Bucket.PutIntelligentTieringV2(context.Background(), opt)
	logStatus(err)

	res, _, err := c.Bucket.GetIntelligentTieringV2(context.Background(), id)
	logStatus(err)
	fmt.Printf("%+v\n", res)

	r, _, err := c.Bucket.ListIntelligentTiering(context.Background())
	logStatus(err)
	fmt.Printf("%+v\n", r)

	_, err = c.Bucket.DeleteIntelligentTiering(context.Background(), id)
	logStatus(err)
}
