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

// 开启 Guetzli
func openGuetzli() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-guangzhou.myqcloud.com")
	cu, _ := url.Parse("http://test-1234567890.pic.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
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

	_, err := c.CI.PutGuetzli(context.Background())
	log_status(err)
}

// 查询 Guetzli
func getGuetzli() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-guangzhou.myqcloud.com")
	cu, _ := url.Parse("http://test-1234567890.pic.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
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

	res, _, err := c.CI.GetGuetzli(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 关闭 Guetzli
func closeGuetzli() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-guangzhou.myqcloud.com")
	cu, _ := url.Parse("http://test-1234567890.pic.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
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

	_, err := c.CI.DeleteGuetzli(context.Background())
	log_status(err)
}

// 下载时压缩
func compressWhenDownload(ctx context.Context, rawurl, obj, localpath, operation string, opt *cos.ObjectGetOptions, id ...string) {
	u, _ := url.Parse(rawurl)
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
	_, err := c.CI.GetToFile(ctx, obj, localpath, operation, opt, id...)
	log_status(err)
}

// 上传时压缩
func compressWhenUpload(ctx context.Context, rawurl, obj, localpath string, pic *cos.PicOperations) {
	u, _ := url.Parse(rawurl)
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
		nil,
		&cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	res, _, err := c.CI.PutFromFile(ctx, obj, localpath, opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.OriginalInfo)
	fmt.Printf("%+v\n", res.ProcessResults)
}

// 云上数据压缩
func compressWhenCloud(ctx context.Context, rawurl, obj string, pic *cos.PicOperations) {
	u, _ := url.Parse(rawurl)
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

	res, _, err := c.CI.ImageProcess(ctx, obj, pic)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.OriginalInfo)
	fmt.Printf("%+v\n", res.ProcessResults)
}

// webp 压缩
func webp() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时压缩
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.webp"
		operation := "imageMogr2/format/webp"
		compressWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时压缩
	{
		obj := "pic/webp/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "deer.webp",
					Rule:   "imageMogr2/format/webp",
				},
			},
		}
		compressWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据压缩
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "webp/deer1.webp",
					Rule:   "imageMogr2/format/webp",
				},
			},
		}
		compressWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// heif 压缩
func heif() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时压缩
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.heif"
		operation := "imageMogr2/format/heif"
		compressWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时压缩
	{
		obj := "pic/heif/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "deer.heif",
					Rule:   "imageMogr2/format/heif",
				},
			},
		}
		compressWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据压缩
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "heif/deer1.heif",
					Rule:   "imageMogr2/format/heif",
				},
			},
		}
		compressWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// tpg 压缩
func tpg() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时压缩
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.tpg"
		operation := "imageMogr2/format/tpg"
		compressWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时压缩
	{
		obj := "pic/tpg/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "deer.tpg",
					Rule:   "imageMogr2/format/tpg",
				},
			},
		}
		compressWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据压缩
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "tpg/deer1.tpg",
					Rule:   "imageMogr2/format/tpg",
				},
			},
		}
		compressWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// tpg 压缩
func avif() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时压缩
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.avif"
		operation := "imageMogr2/format/avif"
		compressWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时压缩
	{
		obj := "pic/avif/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "deer.avif",
					Rule:   "imageMogr2/format/avif",
				},
			},
		}
		compressWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据压缩
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "avif/deer1.avif",
					Rule:   "imageMogr2/format/avif",
				},
			},
		}
		compressWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

func main() {
	// openGuetzli()
	// closeGuetzli()
	// getGuetzli()
	// webp()
	// heif()
	// tpg()
	// avif()
}
