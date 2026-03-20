package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func logStatus(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
	} else {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func main() {
	ak := os.Getenv("SECRETID")
	sk := os.Getenv("SECRETKEY")

	// 存储桶名称，由bucketname-appid 组成
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
			Expire:    time.Hour,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/multipart_presigned_example"
	ctx := context.Background()

	// ===================== Step 1: InitiateMultipartUpload 预签名 =====================
	// InitiateMultipartUpload 需要 POST 方法，query 参数带 uploads
	initOpt := &cos.PresignedURLOptions{
		Query: &url.Values{},
	}
	initOpt.Query.Set("uploads", "")

	initPresignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodPost, name, ak, sk, time.Hour, initOpt)
	if err != nil {
		logStatus(err)
		return
	}
	fmt.Printf("InitMultipartUpload PresignedURL: %s\n\n", initPresignedURL.String())

	// 使用预签名 URL 发起 InitiateMultipartUpload 请求
	initReq, _ := http.NewRequest(http.MethodPost, initPresignedURL.String(), nil)
	initResp, err := http.DefaultClient.Do(initReq)
	if err != nil {
		fmt.Printf("InitiateMultipartUpload request failed: %v\n", err)
		return
	}
	defer initResp.Body.Close()

	initBody, _ := ioutil.ReadAll(initResp.Body)
	if initResp.StatusCode != 200 {
		fmt.Printf("InitiateMultipartUpload failed, status: %s, body: %s\n", initResp.Status, string(initBody))
		return
	}

	var initResult cos.InitiateMultipartUploadResult
	if err := xml.Unmarshal(initBody, &initResult); err != nil {
		fmt.Printf("Parse InitiateMultipartUpload response failed: %v\n", err)
		return
	}
	uploadID := initResult.UploadID
	fmt.Printf("InitiateMultipartUpload 成功, UploadID: %s\n\n", uploadID)

	// ===================== Step 2: UploadPart 预签名 =====================
	// 模拟 3 个分块，每块 1MB 随机数据
	partCount := 3
	partSize := 1 * 1024 * 1024 // 1MB
	parts := make([]cos.Object, 0, partCount)

	for i := 1; i <= partCount; i++ {
		// 生成预签名 URL，query 参数带 partNumber 和 uploadId
		partOpt := &cos.PresignedURLOptions{
			Query: &url.Values{},
		}
		partOpt.Query.Set("partNumber", fmt.Sprintf("%d", i))
		partOpt.Query.Set("uploadId", uploadID)

		partPresignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodPut, name, ak, sk, time.Hour, partOpt)
		if err != nil {
			logStatus(err)
			return
		}
		fmt.Printf("UploadPart %d PresignedURL: %s\n\n", i, partPresignedURL.String())

		// 生成随机数据
		data := make([]byte, partSize)
		rand.Read(data)

		// 使用预签名 URL 上传分块
		partReq, _ := http.NewRequest(http.MethodPut, partPresignedURL.String(), bytes.NewReader(data))
		partReq.ContentLength = int64(len(data))
		partResp, err := http.DefaultClient.Do(partReq)
		if err != nil {
			fmt.Printf("UploadPart %d request failed: %v\n", i, err)
			return
		}
		io.Copy(ioutil.Discard, partResp.Body)
		partResp.Body.Close()

		if partResp.StatusCode != 200 {
			fmt.Printf("UploadPart %d failed, status: %s\n", i, partResp.Status)
			return
		}

		etag := partResp.Header.Get("ETag")
		fmt.Printf("UploadPart %d 成功, ETag: %s\n\n", i, etag)
		parts = append(parts, cos.Object{
			PartNumber: i,
			ETag:       etag,
		})
	}

	// ===================== Step 3: CompleteMultipartUpload 预签名 =====================
	completeOpt := &cos.PresignedURLOptions{
		Query: &url.Values{},
	}
	completeOpt.Query.Set("uploadId", uploadID)

	completePresignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodPost, name, ak, sk, time.Hour, completeOpt)
	if err != nil {
		logStatus(err)
		return
	}
	fmt.Printf("CompleteMultipartUpload PresignedURL: %s\n\n", completePresignedURL.String())

	// 构造 CompleteMultipartUpload 的 XML Body
	completeBody := &cos.CompleteMultipartUploadOptions{
		Parts: parts,
	}
	xmlData, err := xml.Marshal(completeBody)
	if err != nil {
		fmt.Printf("Marshal CompleteMultipartUpload body failed: %v\n", err)
		return
	}

	// 使用预签名 URL 完成分块上传
	completeReq, _ := http.NewRequest(http.MethodPost, completePresignedURL.String(), bytes.NewReader(xmlData))
	completeReq.Header.Set("Content-Type", "application/xml")
	completeResp, err := http.DefaultClient.Do(completeReq)
	if err != nil {
		fmt.Printf("CompleteMultipartUpload request failed: %v\n", err)
		return
	}
	defer completeResp.Body.Close()

	completeRespBody, _ := ioutil.ReadAll(completeResp.Body)
	if completeResp.StatusCode != 200 {
		fmt.Printf("CompleteMultipartUpload failed, status: %s, body: %s\n", completeResp.Status, string(completeRespBody))
		return
	}

	var completeResult cos.CompleteMultipartUploadResult
	if err := xml.Unmarshal(completeRespBody, &completeResult); err != nil {
		fmt.Printf("Parse CompleteMultipartUpload response failed: %v\n", err)
		return
	}

	fmt.Printf("分块上传完成!\n")
	fmt.Printf("  Location: %s\n", completeResult.Location)
	fmt.Printf("  Bucket:   %s\n", completeResult.Bucket)
	fmt.Printf("  Key:      %s\n", completeResult.Key)
	fmt.Printf("  ETag:     %s\n", completeResult.ETag)
}
