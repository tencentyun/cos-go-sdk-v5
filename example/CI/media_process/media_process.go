package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
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

// PostSnapshot 获取媒体文件截图(ci域名)
// https://cloud.tencent.com/document/product/460/73407
func PostSnapshot() {
	c := getClient()
	PostSnapshotOpt := &cos.PostSnapshotOptions{
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Time:   "1",
		Width:  128,
		Height: 128,
		Format: "png",
		Output: &cos.JobOutput{
			Region: "ap-chongqing",
			Bucket: "test-1234567890",
			Object: "test.mp4.png",
		},
	}
	PostSnapshotRes, _, err := c.CI.PostSnapshot(context.Background(), PostSnapshotOpt)
	log_status(err)
	fmt.Printf("%+v\n", PostSnapshotRes)
}

// GetSnapshot 获取媒体文件截图(cos域名)
// https://cloud.tencent.com/document/product/460/49283
func GetSnapshot() {
	c := getClient()
	opt := &cos.GetSnapshotOptions{
		Time: 3,
	}
	resp, err := c.CI.GetSnapshot(context.Background(), "input/test.mp4", opt)
	log_status(err)
	defer resp.Body.Close()

	fd, err := os.OpenFile("test.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		log_status(err)
	}
	_, err = io.Copy(fd, resp.Body)
	fd.Close()
	log_status(err)
}

// GetMediaInfo 获取媒体文件信息(cos域名)
// https://cloud.tencent.com/document/product/460/49284
func GetMediaInfo() {
	c := getClient()
	res, _, err := c.CI.GetMediaInfo(context.Background(), "input/test.mp4", nil)
	log_status(err)
	fmt.Printf("res: %+v\n", res.MediaInfo)
}

// PostMediaInfo 获取媒体文件信息(ci域名)
func PostMediaInfo() {
	c := getClient()
	opt := &cos.GenerateMediaInfoOptions{
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
	}
	res, _, err := c.CI.GenerateMediaInfo(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res.MediaInfo)
}

// GetPrivateM3U8 获取私有 m3u8
// https://cloud.tencent.com/document/product/460/63738
func GetPrivateM3U8() {
	c := getClient()
	getPrivateM3U8Opt := &cos.GetPrivateM3U8Options{
		Expires: 3600,
	}
	res, err := c.CI.GetPrivateM3U8(context.Background(), "output/example.m3u8", getPrivateM3U8Opt)
	log_status(err)
	rspBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Printf("%s\n", rspBody)
}

// ModifyM3U8Token 在加密 M3U8 的 key uri 中增加 Token
// https://cloud.tencent.com/document/product/460/81153
func ModifyM3U8Token() {
	c := getClient()
	getPrivateM3U8Opt := &cos.ModifyM3U8TokenOptions{
		Token: "abc",
	}
	res, err := c.CI.ModifyM3U8Token(context.Background(), "output/example.m3u8", getPrivateM3U8Opt)
	log_status(err)
	rspBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Printf("%s\n", rspBody)
}

func main() {
}
