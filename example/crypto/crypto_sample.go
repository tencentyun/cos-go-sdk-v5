package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/crypto"
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
	os.Exit(1)
}

func cos_max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func simple_put_object() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
	// Case1 上传对象
	name := "test/example2"

	fmt.Println("============== simple_put_object ======================")
	// 该标识信息唯一确认一个主加密密钥, 解密时，需要传入相同的标识信息
	// KMS加密时，该信息设置成EncryptionContext，最大支持1024字符，如果Encrypt指定了该参数，则在Decrypt 时需要提供同样的参数
	materialDesc := make(map[string]string)
	//materialDesc["desc"] = "material information of your master encrypt key"

	// 创建KMS客户端, 可通过 coscrypto.KMSEndpoint 指定KMS域名
	// kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), "ap-guangzhou", coscrypto.KMSEndpoint("kms.internal.tencentcloudapi.com"))
	kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), "ap-guangzhou")
	// 创建KMS主加密密钥，标识信息和主密钥一一对应
	kmsID := os.Getenv("KMSID")
	masterCipher, _ := coscrypto.CreateMasterKMS(kmsclient, kmsID, materialDesc)
	// 创建加密客户端
	client := coscrypto.NewCryptoClient(c, masterCipher)

	contentLength := 1024*1024*10 + 1
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)
	// 加密上传
	_, err = client.Object.Put(context.Background(), name, f, nil)
	log_status(err)

	math_rand.Seed(time.Now().UnixNano())
	rangeStart := math_rand.Intn(contentLength)
	rangeEnd := rangeStart + math_rand.Intn(contentLength-rangeStart)
	opt := &cos.ObjectGetOptions{
		Range: fmt.Sprintf("bytes=%v-%v", rangeStart, rangeEnd),
	}
	// 解密下载
	resp, err := client.Object.Get(context.Background(), name, opt)
	log_status(err)
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	if bytes.Compare(decryptedData, originData[rangeStart:rangeEnd+1]) != 0 {
		fmt.Println("Error: encryptedData != originData")
	}
}

func simple_put_object_from_file() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
	// Case1 上传对象
	name := "test/example1"

	fmt.Println("============== simple_put_object_from_file ======================")
	// 该标识信息唯一确认一个主加密密钥, 解密时，需要传入相同的标识信息
	// KMS加密时，该信息设置成EncryptionContext，最大支持1024字符，如果Encrypt指定了该参数，则在Decrypt 时需要提供同样的参数
	materialDesc := make(map[string]string)
	//materialDesc["desc"] = "material information of your master encrypt key"

	// 创建KMS客户端
	kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), "ap-guangzhou")
	// 创建KMS主加密密钥，标识信息和主密钥一一对应
	kmsID := os.Getenv("KMSID")
	masterCipher, _ := coscrypto.CreateMasterKMS(kmsclient, kmsID, materialDesc)
	// 创建加密客户端
	client := coscrypto.NewCryptoClient(c, masterCipher)

	filepath := "test"
	fd, err := os.Open(filepath)
	log_status(err)
	defer fd.Close()
	m := md5.New()
	io.Copy(m, fd)
	originDataMD5 := m.Sum(nil)

	// 加密上传
	_, err = client.Object.PutFromFile(context.Background(), name, filepath, nil)
	log_status(err)

	// 解密下载
	_, err = client.Object.GetToFile(context.Background(), name, "./test.download", nil)
	log_status(err)

	fd, err = os.Open("./test.download")
	log_status(err)
	defer fd.Close()
	m = md5.New()
	io.Copy(m, fd)
	decryptedDataMD5 := m.Sum(nil)

	if bytes.Compare(decryptedDataMD5, originDataMD5) != 0 {
		fmt.Println("Error: encryptedData != originData")
	}
}

func multi_put_object() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
	// Case1 上传对象
	name := "test/example1"

	fmt.Println("============== multi_put_object ======================")
	// 该标识信息唯一确认一个主加密密钥, 解密时，需要传入相同的标识信息
	// KMS加密时，该信息设置成EncryptionContext，最大支持1024字符，如果Encrypt指定了该参数，则在Decrypt 时需要提供同样的参数
	materialDesc := make(map[string]string)
	//materialDesc["desc"] = "material information of your master encrypt key"

	// 创建KMS客户端
	kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), "ap-guangzhou")
	// 创建KMS主加密密钥，标识信息和主密钥一一对应
	kmsID := os.Getenv("KMSID")
	masterCipher, _ := coscrypto.CreateMasterKMS(kmsclient, kmsID, materialDesc)
	// 创建加密客户端
	client := coscrypto.NewCryptoClient(c, masterCipher)

	contentLength := int64(1024*1024*10 + 1)
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)
	log_status(err)

	// 分块上传
	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		// 每个分块需要16字节对齐
		PartSize: cos_max(1024*1024, (contentLength/16/3)*16),
	}
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), name, nil, &cryptoCtx)
	log_status(err)
	// 切分数据
	chunks, _, err := cos.SplitSizeIntoChunks(contentLength, cryptoCtx.PartSize)
	log_status(err)
	optcom := &cos.CompleteMultipartUploadOptions{}
	for _, chunk := range chunks {
		opt := &cos.ObjectUploadPartOptions{
			ContentLength: chunk.Size,
		}
		f := bytes.NewReader(originData[chunk.OffSet : chunk.OffSet+chunk.Size])
		resp, err := client.Object.UploadPart(context.Background(), name, v.UploadID, chunk.Number, f, opt, &cryptoCtx)
		log_status(err)
		optcom.Parts = append(optcom.Parts, cos.Object{
			PartNumber: chunk.Number, ETag: resp.Header.Get("ETag"),
		})
	}
	_, _, err = client.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	log_status(err)

	resp, err := client.Object.Get(context.Background(), name, nil)
	log_status(err)
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	if bytes.Compare(decryptedData, originData) != 0 {
		fmt.Println("Error: encryptedData != originData")
	}
}

func multi_put_object_from_file() {
	u, _ := url.Parse("https://test-1259654469.cos.ap-guangzhou.myqcloud.com")
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
	// Case1 上传对象
	name := "test/example1"

	fmt.Println("============== multi_put_object_from_file ======================")
	// 该标识信息唯一确认一个主加密密钥, 解密时，需要传入相同的标识信息
	// KMS加密时，该信息设置成EncryptionContext，最大支持1024字符，如果Encrypt指定了该参数，则在Decrypt 时需要提供同样的参数
	materialDesc := make(map[string]string)
	//materialDesc["desc"] = "material information of your master encrypt key"

	// 创建KMS客户端
	kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), "ap-guangzhou")
	// 创建KMS主加密密钥，标识信息和主密钥一一对应
	kmsID := os.Getenv("KMSID")
	masterCipher, _ := coscrypto.CreateMasterKMS(kmsclient, kmsID, materialDesc)
	// 创建加密客户端
	client := coscrypto.NewCryptoClient(c, masterCipher)

	filepath := "test"
	stat, err := os.Stat(filepath)
	log_status(err)
	contentLength := stat.Size()

	// 分块上传
	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		// 每个分块需要16字节对齐
		PartSize: cos_max(1024*1024, (contentLength/16/3)*16),
	}
	// 切分数据
	_, chunks, _, err := cos.SplitFileIntoChunks(filepath, cryptoCtx.PartSize)
	log_status(err)

	// init mulitupload
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), name, nil, &cryptoCtx)
	log_status(err)

	// part upload
	optcom := &cos.CompleteMultipartUploadOptions{}
	for _, chunk := range chunks {
		fd, err := os.Open(filepath)
		log_status(err)
		opt := &cos.ObjectUploadPartOptions{
			ContentLength: chunk.Size,
		}
		fd.Seek(chunk.OffSet, os.SEEK_SET)
		resp, err := client.Object.UploadPart(context.Background(), name, v.UploadID, chunk.Number, cos.LimitReadCloser(fd, chunk.Size), opt, &cryptoCtx)
		log_status(err)
		optcom.Parts = append(optcom.Parts, cos.Object{
			PartNumber: chunk.Number, ETag: resp.Header.Get("ETag"),
		})
	}
	// complete upload
	_, _, err = client.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	log_status(err)

	_, err = client.Object.GetToFile(context.Background(), name, "test.download", nil)
	log_status(err)
}

func main() {
	simple_put_object()
	simple_put_object_from_file()
	multi_put_object()
	multi_put_object_from_file()
}
