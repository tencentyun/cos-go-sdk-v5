package main

import (
	"context"
	"fmt"
	"hash/crc64"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func logStatus(err error) {
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

// calcFileCRC64 计算本地文件的 CRC64-ECMA 值
func calcFileCRC64(filePath string) (uint64, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	tb := crc64.MakeTable(crc64.ECMA)
	h := crc64.New(tb)
	if _, err := io.Copy(h, fd); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

// verifyCRC64 对比本地文件 CRC64 与响应头 x-cos-hash-crc64ecma 是否一致
func verifyCRC64(localFile string, resp *cos.Response) {
	localCRC, err := calcFileCRC64(localFile)
	if err != nil {
		fmt.Printf("  计算本地文件 CRC64 失败: %v\n", err)
		return
	}
	serverCRCStr := resp.Header.Get("x-cos-hash-crc64ecma")
	if serverCRCStr == "" {
		fmt.Printf("  响应头无 x-cos-hash-crc64ecma，跳过比对\n")
		return
	}
	serverCRC, _ := strconv.ParseUint(serverCRCStr, 10, 64)
	if localCRC == serverCRC {
		fmt.Printf("  CRC64 验证通过: 本地文件CRC64=%v, 服务端CRC64=%v\n", localCRC, serverCRC)
	} else {
		fmt.Printf("  CRC64 验证失败! 本地文件CRC64=%v, 服务端CRC64=%v\n", localCRC, serverCRC)
	}
}

func main() {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台"存储桶概览"查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	name := "test/example2"

	// Case1 使用 GetToFile 下载对象，SDK 内部自动进行 CRC64 校验
	// 如果服务端返回的 x-cos-hash-crc64ecma 头部与本地计算的 CRC64 不一致，会返回 "verification failed" 错误
	fmt.Println("========== Case1: GetToFile CRC64 校验 ==========")
	localFile1 := "gettofile_crc_test"
	resp, err := c.Object.GetToFile(context.Background(), name, localFile1, nil)
	if err != nil {
		logStatus(err)
	} else {
		fmt.Printf("GetToFile 成功, x-cos-hash-crc64ecma: %v\n",
			resp.Header.Get("x-cos-hash-crc64ecma"))
		verifyCRC64(localFile1, resp)
	}
	os.Remove(localFile1)

	// Case2 使用 GetToFile 带 Listener 下载对象，同时查看进度和进行 CRC64 校验
	fmt.Println("\n========== Case2: GetToFile 带进度 Listener CRC64 校验 ==========")
	localFile2 := "gettofile_listener_crc_test"
	opt := &cos.ObjectGetOptions{
		Listener: &cos.DefaultProgressListener{},
	}
	resp, err = c.Object.GetToFile(context.Background(), name, localFile2, opt)
	if err != nil {
		logStatus(err)
	} else {
		fmt.Printf("GetToFile 带 Listener 成功, x-cos-hash-crc64ecma: %v\n",
			resp.Header.Get("x-cos-hash-crc64ecma"))
		verifyCRC64(localFile2, resp)
	}
	os.Remove(localFile2)

	// Case3 使用 Download 多协程分块下载对象，SDK 内部使用 CRC64Combine 合并各分块 CRC 并与服务端返回的总 CRC 进行校验
	fmt.Println("\n========== Case3: Download 多分块 CRC64 合并校验 ==========")
	localFile3 := "download_crc_test"
	downOpt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
		PartSize:       1, // 1MB 分块大小
	}
	resp, err = c.Object.Download(context.Background(), name, localFile3, downOpt)
	if err != nil {
		logStatus(err)
	} else {
		fmt.Printf("Download 成功, x-cos-hash-crc64ecma: %v\n",
			resp.Header.Get("x-cos-hash-crc64ecma"))
		verifyCRC64(localFile3, resp)
	}
	os.Remove(localFile3)

	// Case4 使用 Download 带 CheckPoint 断点续载 + CRC64 校验
	fmt.Println("\n========== Case4: Download 带 CheckPoint CRC64 校验 ==========")
	localFile4 := "download_checkpoint_crc_test"
	downOpt2 := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
		PartSize:       1,
		CheckPoint:     true,
	}
	resp, err = c.Object.Download(context.Background(), name, localFile4, downOpt2)
	if err != nil {
		logStatus(err)
	} else {
		fmt.Printf("Download 带 CheckPoint 成功, x-cos-hash-crc64ecma: %v\n",
			resp.Header.Get("x-cos-hash-crc64ecma"))
		verifyCRC64(localFile4, resp)
	}
	os.Remove(localFile4)

	// Case5 禁用 CRC 校验（DisableChecksum），下载时跳过 CRC64 校验
	fmt.Println("\n========== Case5: Download 禁用 CRC 校验 ==========")
	localFile5 := "download_nocrc_test"
	downOpt3 := &cos.MultiDownloadOptions{
		ThreadPoolSize:  5,
		PartSize:        1,
		DisableChecksum: true,
	}
	resp, err = c.Object.Download(context.Background(), name, localFile5, downOpt3)
	if err != nil {
		logStatus(err)
	} else {
		fmt.Printf("Download 禁用 CRC 校验成功, x-cos-hash-crc64ecma: %v\n",
			resp.Header.Get("x-cos-hash-crc64ecma"))
		verifyCRC64(localFile5, resp)
	}
	os.Remove(localFile5)
}
