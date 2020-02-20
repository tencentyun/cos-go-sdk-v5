package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"io/ioutil"
	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

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
				ResponseBody:   false,
			},
		},
	})

	opt := &cos.ObjectSelectOptions{
		Expression:     "Select * from COSObject",
		ExpressionType: "SQL",
		InputSerialization: &cos.SelectInputSerialization{
			CSV: &cos.CSVInputSerialization{
				FileHeaderInfo: "IGNORE",
			},
		},
		OutputSerialization: &cos.SelectOutputSerialization{
			CSV: &cos.CSVOutputSerialization{
				RecordDelimiter: "\n",
			},
		},
		RequestProgress: "TRUE",
	}
	res, err := c.Object.Select(context.Background(), "test.csv", opt)
	if err != nil {
		panic(err)
	}
	defer res.Close()
	data, err := ioutil.ReadAll(res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("data: %v\n", string(data))
	resp, _ := res.(*cos.ObjectSelectResponse)
	fmt.Printf("data: %+v\n", resp.Frame)

	// Select To File
	_, err = c.Object.SelectToFile(context.Background(), "test.csv", "./test.csv", opt)
	if err != nil {
		panic(err)
	}
}
