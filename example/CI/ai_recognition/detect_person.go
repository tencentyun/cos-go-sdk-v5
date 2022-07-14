package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

// 身份证识别 (云上处理)
func idCardOCRWhenCloud() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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
	obj := "pic/idcard_1.png"
	query := &cos.IdCardOCROptions{
		Config: &cos.IdCardOCROptionsConfig{
			CropPortrait:    true,
			CropIdCard:      true,
			CopyWarn:        true,
			BorderCheckWarn: true,
			ReshootWarn:     true,
			DetectPsWarn:    true,
			TempIdWarn:      true,
			InvalidDateWarn: true,
			Quality:         true,
			MultiCardDetect: true,
		},
	}
	res, _, err := c.CI.IdCardOCRWhenCloud(context.Background(), obj, query)
	log_status(err)
	fmt.Printf("%+v\n", res)
	if res.AdvancedInfo != nil && len(res.AdvancedInfo.IdCard) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.AdvancedInfo.IdCard)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("idcard_cut.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
	if res.AdvancedInfo != nil && len(res.AdvancedInfo.Portrait) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.AdvancedInfo.Portrait)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("idcard_portrait.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
}

// 身份证识别 (云上处理)
func idCardOCRWhenUpload() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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
	obj := "pic/idcard_0.png"
	localFile := "./idcard_0.png"
	query := &cos.IdCardOCROptions{
		Config: &cos.IdCardOCROptionsConfig{
			CropPortrait:    true,
			CropIdCard:      true,
			CopyWarn:        true,
			BorderCheckWarn: true,
			ReshootWarn:     true,
			DetectPsWarn:    true,
			TempIdWarn:      true,
			InvalidDateWarn: true,
			Quality:         true,
			MultiCardDetect: true,
		},
	}

	res, _, err := c.CI.IdCardOCRWhenUpload(context.Background(), obj, localFile, query, nil)
	log_status(err)
	fmt.Printf("%+v\n", res)
	if res.AdvancedInfo != nil && len(res.AdvancedInfo.IdCard) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.AdvancedInfo.IdCard)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("idcard_cut.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
	if res.AdvancedInfo != nil && len(res.AdvancedInfo.Portrait) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.AdvancedInfo.Portrait)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("idcard_portrait.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
}

// 获取数字验证码
func getLiveCode() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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
	res, _, err := c.CI.GetLiveCode(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 获取动作顺序
func getActionSequence() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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
	res, _, err := c.CI.GetActionSequence(context.Background())
	log_status(err)
	fmt.Printf("%+v\n", res)
}

// 活体人脸合身 (云上处理)
func livenessRecognitionWhenCloud() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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
	obj := "pic/self.mp4"
	query := &cos.LivenessRecognitionOptions{
		IdCard:       "111222xxxxxxxxxxxx",
		Name:         "张三",
		LivenessType: "SILENT",
		BestFrameNum: 2,
	}
	res, _, err := c.CI.LivenessRecognitionWhenCloud(context.Background(), obj, query)
	log_status(err)
	fmt.Printf("%f\n", res.Sim)
	if len(res.BestFrameBase64) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.BestFrameBase64)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("image0.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
	if len(res.BestFrameList) > 0 {
		for i, v := range res.BestFrameList {
			d, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fd, err := os.OpenFile("image"+strconv.Itoa(i+1)+".jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					r := bytes.NewReader(d)
					_, err = io.Copy(fd, r)
					if err != nil {
						fmt.Println(err.Error())
					}
					fd.Close()
				}
			}
		}
	}
}

// 活体人脸合身 (上传时处理)
func livenessRecognitionWhenUpload() {
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
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
	obj := "pic/self1.mp4"
	localFile := "./self.mp4"
	query := &cos.LivenessRecognitionOptions{
		IdCard:       "111222xxxxxxxxxxxx",
		Name:         "张三",
		LivenessType: "SILENT",
		BestFrameNum: 2,
	}

	res, _, err := c.CI.LivenessRecognitionWhenUpload(context.Background(), obj, localFile, query, nil)
	log_status(err)
	fmt.Printf("%f\n", res.Sim)
	if len(res.BestFrameBase64) > 0 {
		d, err := base64.StdEncoding.DecodeString(res.BestFrameBase64)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fd, err := os.OpenFile("image0.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				r := bytes.NewReader(d)
				_, err = io.Copy(fd, r)
				if err != nil {
					fmt.Println(err.Error())
				}
				fd.Close()
			}
		}
	}
	if len(res.BestFrameList) > 0 {
		for i, v := range res.BestFrameList {
			d, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fd, err := os.OpenFile("image"+strconv.Itoa(i+1)+".jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					r := bytes.NewReader(d)
					_, err = io.Copy(fd, r)
					if err != nil {
						fmt.Println(err.Error())
					}
					fd.Close()
				}
			}
		}
	}
}

func main() {
	// idCardOCRWhenCloud()
	// idCardOCRWhenUpload()
	// getLiveCode()
	// getActionSequence()
	// livenessRecognitionWhenCloud()
	livenessRecognitionWhenUpload()
}
