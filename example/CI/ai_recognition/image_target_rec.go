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

// 图⽚⼈⻋通⽤识别功能(上传时)
func imageTargetRecWhenUpload() {
	u, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
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

	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: nil,
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "/pic/out.jpeg",
				Rule:   "ci-process=ImageTargetRec&detect-type=car,face,plate,body&cover=1",
			},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	name := "pic/cup.png"
	local_filename := "./test_output_02_000.jpg"
	res, _, err := c.CI.PutFromFile(context.Background(), name, local_filename, opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.OriginalInfo)
	fmt.Printf("%+v\n", res.ProcessResults)
	fmt.Printf("%+v\n", res.ImgTargetRecResult)
}

// 图⽚⼈⻋通⽤识别功能(云上数据处理)
func imageTargetRecWhenCloud() {
	u, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
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
				ResponseBody:   true,
			},
		},
	})

	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "/pic/cup_out.jpeg",
				Rule:   "ci-process=ImageTargetRec&detect-type=car,face,plate,body&cover=1",
			},
		},
	}

	key := "pic/cup.png"
	res, _, err := c.CI.ImageProcess(context.Background(), key, pic)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.ImgTargetRecResult)
}

func main() {
	// imageTargetRecWhenUpload()
	imageTargetRecWhenCloud()
}
