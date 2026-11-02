package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"hash/crc64"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// generateRandomFile 生成一个指定大小的随机内容文件，用于示例演示上传。
// 返回文件路径（放在 os.TempDir() 下，调用方负责最终删除）。
func generateRandomFile(prefix string, size int64) (string, error) {
	f, err := os.CreateTemp("", prefix+"-*.bin")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.CopyN(f, rand.Reader, size); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// crc64OfFile 计算本地文件的 CRC64（ECMA 多项式），作为兜底比对依据。
func crc64OfFile(path string) (uint64, error) {
	fd, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	if _, err := io.Copy(h, fd); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

// completeWithHeadByUpload 演示高级上传接口 Upload 的 EnableCompleteWithHead 开关：
// Complete 失败并在重试过程中返回 404 NoSuchUpload 时，由 SDK 自动 HEAD 对象 + CRC64 比对
// 确认服务端是否已合并完成，避免把实际已上传成功的文件误判为失败。
//
// 开关开启后：
//   - 若 Conf.EnableCRC 且 !DisableChecksum：SDK 自动把本地 CRC64 传入做严格比对
//   - 否则：仅凭对象存在判定成功
func completeWithHeadByUpload(c *cos.Client, filepath, name string) {
	opt := &cos.MultiUploadOptions{
		PartSize:       32, // 32MB
		ThreadPoolSize: 3,
		// 开启 Complete 幂等兜底，默认 false
		EnableCompleteWithHead: true,
	}
	v, resp, err := c.Object.Upload(context.Background(), name, filepath, opt)
	if err != nil {
		fmt.Printf("[Upload] failed: %v\n", err)
		return
	}
	fmt.Printf("[Upload] succeeded, Key: %v, ETag: %v, StatusCode: %d\n", v.Key, v.ETag, resp.StatusCode)
}

// completeWithHeadByAPI 演示复合接口 CompleteMultipartUploadWithHead：
// 自行完成 Init / UploadPart 后调用；通过 Crc64CheckInHead 控制兜底时是否比对 CRC64。
//
// Crc64CheckInHead:
//   - nil：命中 NoSuchUpload 兜底时，HEAD 到对象即判定成功，不比对 CRC64
//   - 非 nil：比较本地 CRC64 与服务端 x-cos-hash-crc64ecma，一致才判定成功
func completeWithHeadByAPI(c *cos.Client, filepath, name string) {
	// 1. 本地 CRC64（可选，用于严格兜底比对）
	localcrc, err := crc64OfFile(filepath)
	if err != nil {
		fmt.Printf("[API] crc64 of file failed: %v\n", err)
		return
	}

	// 2. Init
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		fmt.Printf("[API] InitiateMultipartUpload failed: %v\n", err)
		return
	}
	uploadID := v.UploadID

	// 3. UploadPart（本示例仅上传一个分块，演示用）
	fd, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("[API] open file failed: %v\n", err)
		return
	}
	partResp, err := c.Object.UploadPart(context.Background(), name, uploadID, 1, fd, nil)
	fd.Close()
	if err != nil {
		fmt.Printf("[API] UploadPart failed: %v\n", err)
		return
	}

	// 4. CompleteMultipartUploadWithHead
	// 入参、出参与 CompleteMultipartUpload 完全一致，调用端可无缝替换
	optcom := &cos.CompleteMultipartUploadOptions{
		Parts:            []cos.Object{{PartNumber: 1, ETag: partResp.Header.Get("ETag")}},
		Crc64CheckInHead: &localcrc, // 若设为 nil 则兜底时仅校验对象存在
	}
	res, resp, err := c.Object.CompleteMultipartUploadWithHead(context.Background(), name, uploadID, optcom)
	if err != nil {
		// 兜底失败时错误仍是 CosError（%w 包装），上层可用 IsCOSError 识别
		if e, ok := cos.IsCOSError(err); ok {
			fmt.Printf("[API] CompleteMultipartUploadWithHead failed, cos code=%v, msg=%v, raw=%v\n",
				e.Code, e.Message, err)
			return
		}
		fmt.Printf("[API] CompleteMultipartUploadWithHead failed: %v\n", err)
		return
	}
	fmt.Printf("[API] CompleteMultipartUploadWithHead succeeded, Key: %v, ETag: %v, StatusCode: %d\n",
		res.Key, res.ETag, resp.StatusCode)
}

// dropInReplacement 演示无缝切换：不设 Crc64CheckInHead 时，CompleteMultipartUploadWithHead
// 的行为与 CompleteMultipartUpload 完全一致（签名、返回值、错误类型均不变）。
func dropInReplacement(c *cos.Client, name, uploadID string, parts []cos.Object) {
	optcom := &cos.CompleteMultipartUploadOptions{Parts: parts}
	// 把方法名从 CompleteMultipartUpload 换成 CompleteMultipartUploadWithHead 即可
	res, resp, err := c.Object.CompleteMultipartUploadWithHead(context.Background(), name, uploadID, optcom)
	if err != nil {
		fmt.Printf("[drop-in] failed: %v\n", err)
		return
	}
	fmt.Printf("[drop-in] succeeded, Key: %v, ETag: %v, StatusCode: %d\n", res.Key, res.ETag, resp.StatusCode)
}

func main() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
		},
	})

	// 生成一个 2.5MB 的随机内容文件作为演示用的上传源，避免依赖外部 ./test.file。
	// 大小 > PartSize(32MB) 时会走分块合并路径；此处小文件会走简单上传，
	// 若需触发真实分块逻辑可把 size 调整为 > 32*1024*1024。
	localFile, err := generateRandomFile("cos-cwh-example", 2*1024*1024+512*1024)
	if err != nil {
		fmt.Printf("generate random file failed: %v\n", err)
		return
	}
	defer os.Remove(localFile)
	fmt.Printf("random test file: %s\n", localFile)

	name := "exampleobject"

	// 方式一：高级上传接口开启 Complete 幂等兜底（推荐，无需手工算 CRC）
	completeWithHeadByUpload(c, localFile, name)

	// 方式二：直接使用复合接口（适合自行拆分 Init / UploadPart / Complete 的场景）
	completeWithHeadByAPI(c, localFile, "test/complete_with_head_api.txt")

	// 方式三：当作 CompleteMultipartUpload 的 drop-in replacement 使用
	_ = dropInReplacement
}
