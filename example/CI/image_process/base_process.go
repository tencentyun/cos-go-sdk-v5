package main

import (
	"context"
	"encoding/base64"
	"fmt"
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

// 下载时处理
func processWhenDownload(ctx context.Context, rawurl, obj, localpath, operation string, opt *cos.ObjectGetOptions, id ...string) {
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

// 上传时处理
func processWhenUpload(ctx context.Context, rawurl, obj, localpath string, pic *cos.PicOperations) {
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
		ACLHeaderOptions: nil,
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
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

// 分块上传时处理
func processWhenPartUpload(ctx context.Context, rawurl, obj, localpath string, pic *cos.PicOperations) {
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
	opt := &cos.MultiUploadOptions{
		OptIni: &cos.InitiateMultipartUploadOptions{
			ACLHeaderOptions: nil,
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				XOptionHeader: &http.Header{},
			},
		},
	}

	opt.OptIni.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	res, _, err := c.Object.UploadWithPicOperations(context.Background(), obj, localpath, opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 云上数据处理
func processWhenCloud(ctx context.Context, rawurl, obj string, pic *cos.PicOperations) {
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

// 云上数据处理
func processWhenCloudWithHeader(ctx context.Context, rawurl, obj string, imageProcessHeader *cos.ImageProcessHeader) {
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
	res, _, err := c.CI.ImageProcessWithHeader(ctx, obj, imageProcessHeader)
	log_status(err)
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.OriginalInfo)
	fmt.Printf("%+v\n", res.ProcessResults)
}

// 添加盲水印
func blindWatermark() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./watermark3.jpg"
		operation := "watermark/3/type/2/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw=="
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/upload/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "watermark3.jpg",
					Rule:   "watermark/3/type/2/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/upload/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "watermark3.jpg",
					Rule:   "watermark/3/type/2/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 提取盲水印
func extractBlindWatermark() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 上传时处理
	{
		obj := "pic/upload/deer.jpg"
		filepath := "./watermark3.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "wm1.jpg",
					Rule:   "watermark/4/type/2/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/upload/watermark3.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "wm2.jpg",
					Rule:   "watermark/4/type/2/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 缩放
func thumbnail() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./thumbnail_50%.jpg"
		operation := "imageMogr2/thumbnail/!50p"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/thumbnail/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "thumbnail_50%.jpg",
					Rule:   "imageMogr2/thumbnail/!50p",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "thumbnail/thumbnail_70%.jpg",
					Rule:   "imageMogr2/thumbnail/!70p",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 裁剪
func tailor() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./cut_600x600.jpg"
		operation := "imageMogr2/cut/600x600x100x10"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/tailor/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "crop_500x.jpg",
					Rule:   "imageMogr2/crop/500x",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "tailor/iradius_300.jpg",
					Rule:   "imageMogr2/iradius/300",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 旋转
func rotate() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./rotate_45.jpg"
		operation := "imageMogr2/rotate/45"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/rotate/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "rotate_45.jpg",
					Rule:   "imageMogr2/rotate/45",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "rotate/horizontal.jpg",
					Rule:   "imageMogr2/flip/horizontal",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 格式变换
func format() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.webp"
		operation := "imageMogr2/format/webp"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/format/deer.jpg"
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
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "format/deer.png",
					Rule:   "imageMogr2/format/webp",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 质量变换
func quality() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./quality_60.jpg"
		operation := "imageMogr2/quality/60"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/quality/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "quality_60.jpg",
					Rule:   "imageMogr2/quality/60",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "quality/rquality_50.jpg",
					Rule:   "imageMogr2/rquality/50",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 高斯模糊
func blur() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./blur_8x5.jpg"
		operation := "imageMogr2/blur/8x5"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/blur/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "blur_8x5.jpg",
					Rule:   "imageMogr2/blur/8x5",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "blur/blur_10x3.jpg",
					Rule:   "imageMogr2/blur/10x3",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 亮度
func bright() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./bright_+50.jpg"
		operation := "imageMogr2/bright/50"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/bright/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "bright_+50.jpg",
					Rule:   "imageMogr2/bright/50",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "bright/bright_-30.jpg",
					Rule:   "imageMogr2/bright/-30",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 对比度
func contrast() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./contrast_+50.jpg"
		operation := "imageMogr2/contrast/50"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/contrast/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "contrast_+50.jpg",
					Rule:   "imageMogr2/contrast/50",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "contrast/contrast_-30.jpg",
					Rule:   "imageMogr2/contrast/-30",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 锐化
func sharpen() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./sharpen_70.jpg"
		operation := "imageMogr2/sharpen/70"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/sharpen/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "sharpen_70+50.jpg",
					Rule:   "imageMogr2/sharpen/70",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "sharpen/sharpen_120.jpg",
					Rule:   "imageMogr2/sharpen/120",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 灰度图
func grayscale() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./grayscale.jpg"
		operation := "imageMogr2/grayscale/1"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/grayscale/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "grayscale.jpg",
					Rule:   "imageMogr2/grayscale/1",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "grayscale/grayscale0.jpg",
					Rule:   "imageMogr2/grayscale/0",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 图片水印
func picWatermark() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./picwatermark.jpg"
		operation := "watermark/1/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==/gravity/northeast"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/picwatermark/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "picwatermark.jpg",
					Rule:   "watermark/1/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==/gravity/northeast",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "picwatermark/picwatermark1.jpg",
					Rule:   "watermark/1/image/aHR0cDovL2xpbGFuZy0xMjUzOTYwNDU0LmNvcy5hcC1jaG9uZ3FpbmcubXlxY2xvdWQuY29tL3BpYy9iZWFyLnBuZw==/gravity/northeast/batch/1/degree/45/dissolve/40/",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 文字水印
func textWatermark() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./textwatermark.jpg"
		operation := "watermark/2/text/Y2xvdWQudGVuY2VudC5jb20v/fill/IzNEM0QzRA/fontsize/20/dissolve/50/gravity/northeast/dx/20/dy/20/batch/1/degree/45"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/textwatermark/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "textwatermark.jpg",
					Rule:   "watermark/2/text/Y2xvdWQudGVuY2VudC5jb20v/fill/IzNEM0QzRA/fontsize/20/dissolve/50/gravity/northeast/dx/20/dy/20/batch/1/degree/45",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "textwatermark/textwatermark1.jpg",
					Rule:   "watermark/2/text/Y2xvdWQudGVuY2VudC5jb20v/fill/IzNEM0QzRA/fontsize/20/dissolve/50/gravity/northeast/dx/20/dy/20/batch/1/degree/45",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 获取图片基本信息
func getImageInfoBase() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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

	name := "pic/deer.jpg"
	operation := "imageInfo"
	_, err := c.CI.Get(context.Background(), name, operation, nil)
	log_status(err)
}

// 获取图片EXIF
func getImageInfoExif() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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

	name := "pic/deer.jpg"
	operation := "exif"
	_, err := c.CI.Get(context.Background(), name, operation, nil)
	log_status(err)
}

// 获取图片主色调
func getImageInfoImageAve() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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

	name := "pic/deer.jpg"
	operation := "imageAve"
	_, err := c.CI.Get(context.Background(), name, operation, nil)
	log_status(err)
}

// 去除元信息
func strip() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./strip.jpg"
		operation := "imageMogr2/strip"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/strip/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "strip.jpg",
					Rule:   "imageMogr2/strip",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "strip/strip1.jpg",
					Rule:   "imageMogr2/strip",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 快速缩略模板
func thumbnailTemplate() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./thumbnailTemplate.jpg"
		operation := "imageView2/1/w/400/h/600/q/85"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/thumbnailTemplate/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "400x600.jpg",
					Rule:   "imageView2/1/w/400/h/600/q/85",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "thumbnailTemplate/200x300.jpg",
					Rule:   "imageView2/1/w/200/h/300/q/85",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 限制图片大小
func limit() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./15k.png"
		operation := "imageMogr2/strip/format/png/size-limit/15k!"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/limit/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "15k.png",
					Rule:   "imageMogr2/strip/format/png/size-limit/15k!",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "limit/100k.png",
					Rule:   "imageMogr2/strip/format/png/size-limit/100k!",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 操作管道符
func pipe() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./pipe.jpg"
		operation := "imageMogr2/thumbnail/!50p|watermark/2/text/5pWw5o2u5LiH6LGh/fill/I0ZGRkZGRg==/fontsize/30/dx/20/dy/20"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/pipe/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "pipe.jpg",
					Rule:   "imageMogr2/thumbnail/!50p|watermark/2/text/5pWw5o2u5LiH6LGh/fill/I0ZGRkZGRg==/fontsize/30/dx/20/dy/20",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "pipe/pipe_80%.jpg",
					Rule:   "imageMogr2/thumbnail/!80p|watermark/2/text/5pWw5o2u5LiH6LGh/fill/I0ZGRkZGRg==/fontsize/30/dx/20/dy/20",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

func commonProcess() {
	rawurl := "https://test-1234567890.cos.ap-chongqing.myqcloud.com"
	// 下载时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.jpg"
		operation := "imageMogr2/xxx"
		processWhenDownload(context.Background(), rawurl, obj, filepath, operation, nil)
	}
	// 上传时处理
	{
		obj := "pic/pipe/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "pipe.jpg",
					Rule:   "imageMogr2/xxx",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 分块上传时处理
	{
		obj := "pic/deer.jpg"
		filepath := "./deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "pipe.jpg",
					Rule:   "imageMogr2/xxx",
				},
			},
		}
		processWhenPartUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "pic/deer.jpg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "pipe/pipe_80%.jpg",
					Rule:   "imageMogr2/xxx",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 添加文件盲水印3.0
func blindWatermark3() {
	rawurl := "https://test-1234567890.cos.ap-beijing.myqcloud.com"
	// 上传时处理
	{
		obj := "test_slim_blindWatermark3.png"
		filepath := "./test_slim.png"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "test_slim_blindWatermark3_upload.jpg",
					Rule:   "watermark/3/type/3/text/5L2g5aW95LiH6LGh/version/3.0",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "test_slim.png"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "test_slim_blindWatermark3_cloud.png",
					Rule:   "watermark/3/type/3/text/5L2g5aW95LiH6LGh/version/3.0",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 提取文字盲水印3.0
func extractBlindWatermark3() {
	rawurl := "https://test-1234567890.cos.ap-beijing.myqcloud.com"
	// 上传时处理
	{
		obj := "test_slim.png"
		filepath := "./test_slim.png"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "empty",
					Rule:   "watermark/4/type/3/version/3.0",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}
	// 云上数据处理
	{
		obj := "test_slim_blindWatermark3_cloud.png"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "empty",
					Rule:   "watermark/4/type/3/version/3.0",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}
}

// 文字水印和AIGC
func textWatermarkAndAIGC() {
	rawurl := "https://test-12500000.cos.ap-chongqing.myqcloud.com"

	aigcParam := make(map[string]string)
	aigcParam["Label"] = "label"
	aigcParam["ContentProducer"] = "contentProducer"
	aigcParam["ProduceID"] = "ProduceID"
	aigcParam["ReservedCode1"] = "ReservedCode1"
	aigcParam["ContentPropagator"] = "ContentPropagator"
	aigcParam["PropagateID"] = "PropagateID"
	aigcParam["ReservedCode2"] = "ReservedCode2"

	AIGCMetadata := "AIGCMetadata/"
	for k, v := range aigcParam {
		AIGCMetadata += k + "/" + base64.StdEncoding.EncodeToString([]byte(v)) + "/"
	}

	// 上传时处理
	{
		obj := "pic/textwatermark/deer.jpg"
		filepath := "./deer.jpeg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "textwatermark.jpg",
					Rule:   "imageMogr2/" + AIGCMetadata + "|watermark/2/text/6IW-6K6v5LqRwrfkuIfosaHkvJjlm74/fill/IzNEM0QzRA/fontsize/20/dissolve/50/gravity/northeast/dx/20/dy/20/batch/1/degree/45",
				},
			},
		}
		processWhenUpload(context.Background(), rawurl, obj, filepath, pic)
	}

	// 云上数据处理
	{
		obj := "pic/deer.jpeg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "textwatermark/textwatermark1.jpg",
					Rule:   "imageMogr2/" + AIGCMetadata + "|watermark/2/text/6IW-6K6v5LqRwrfkuIfosaHkvJjlm74/fill/IzNEM0QzRA/fontsize/20/dissolve/50/gravity/northeast/dx/20/dy/20/batch/1/degree/45",
				},
			},
		}
		processWhenCloud(context.Background(), rawurl, obj, pic)
	}

	// 云上数据处理
	{
		obj := "pic/deer.jpeg"
		pic := &cos.PicOperations{
			IsPicInfo: 1,
			Rules: []cos.PicOperationsRules{
				{
					FileId: "textwatermark/textwatermark1.jpg",
					Rule:   "imageMogr2/" + AIGCMetadata + "|watermark/2/text/6IW-6K6v5LqRwrfkuIfosaHkvJjlm74/fill/IzNEM0QzRA/fontsize/20/dissolve/50/gravity/northeast/dx/20/dy/20/batch/1/degree/45",
				},
			},
		}
		imageProcessHeader := &cos.ImageProcessHeader{
			PicOperations: cos.EncodePicOperations(pic),
			XOptionHeader: &http.Header{},
		}
		imageProcessHeader.XOptionHeader.Add("x-cos-meta-aigctag", "csig")
		processWhenCloudWithHeader(context.Background(), rawurl, obj, imageProcessHeader)
	}

}

func getAIGC() {
	u, _ := url.Parse("https://test-12500000.cos.ap-chongqing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv("SECRETID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv("SECRETKEY"),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
	name := "pic/textwatermark/textwatermark.jpg"
	opt := &cos.ObjectGetOptions{
		CiProcess: "ImageAIGCMetadata",
	}
	resp, err := c.Object.Get(context.Background(), name, opt)
	log_status(err)
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}

func main() {
	// commonProcess()
}
