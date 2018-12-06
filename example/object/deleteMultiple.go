package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"bytes"
	"io"

	"math/rand"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func genBigData(blockSize int) []byte {
	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func uploadMulti(c *cos.Client) []string {
	names := []string{}
	data := genBigData(1024 * 1024 * 1)
	ctx := context.Background()
	var r io.Reader
	var name string
	n := 3

	for n > 0 {
		name = fmt.Sprintf("test/test_multi_delete_%s", time.Now().Format(time.RFC3339))
		r = bytes.NewReader(data)

		c.Object.Put(ctx, name, r, nil)
		names = append(names, name)
		n--
	}
	return names
}

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
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
	ctx := context.Background()

	names := uploadMulti(c)
	names = append(names, []string{"a", "b", "c", "a+bc/xx&?+# "}...)
	obs := []cos.Object{}
	for _, v := range names {
		obs = append(obs, cos.Object{Key: v})
	}
	//sha1 := ""
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
		//XCosSha1: sha1,
		//Quiet: true,
	}

	c = cos.NewClient(b, &http.Client{
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

	v, _, err := c.Object.DeleteMulti(ctx, opt)
	if err != nil {
		panic(err)
	}

	for _, x := range v.DeletedObjects {
		fmt.Printf("deleted %s\n", x.Key)
	}
	for _, x := range v.Errors {
		fmt.Printf("error %s, %s, %s\n", x.Key, x.Code, x.Message)
	}
}
