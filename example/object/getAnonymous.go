package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"io/ioutil"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func upload(c *cos.Client, name string) {
	f := strings.NewReader("test")
	f = strings.NewReader("test xxx")
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	c.Object.Put(context.Background(), name, f, opt)
	return
}

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, nil)

	name := "test/anonymous_get.go"
	upload(c, name)

	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		panic(err)
		return
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
