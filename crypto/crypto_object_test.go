package coscrypto_test

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/crypto"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

const (
	kAppid  = 1259654469
	kBucket = "cosgosdktest-1259654469"
	kRegion = "ap-guangzhou"
)

type CosTestSuite struct {
	suite.Suite
	Client  *cos.Client
	CClient *coscrypto.CryptoClient
	Master  coscrypto.MasterCipher
}

func (s *CosTestSuite) SetupSuite() {
	u, _ := url.Parse("https://" + kBucket + ".cos." + kRegion + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	s.Client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})
	material := make(map[string]string)
	material["desc"] = "cos crypto suite test"
	kmsclient, _ := coscrypto.NewKMSClient(s.Client.GetCredential(), kRegion)
	s.Master, _ = coscrypto.CreateMasterKMS(kmsclient, os.Getenv("KMSID"), material)
	s.CClient = coscrypto.NewCryptoClient(s.Client, s.Master)
	opt := &cos.BucketPutOptions{
		XCosACL: "public-read",
	}
	r, err := s.Client.Bucket.Put(context.Background(), opt)
	if err != nil && r != nil && r.StatusCode == 409 {
		fmt.Println("BucketAlreadyOwnedByYou")
	} else if err != nil {
		assert.Nil(s.T(), err, "PutBucket Failed")
	}
}

func (s *CosTestSuite) TestPutGetDeleteObject_DecryptWithKey_10MB() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	originData := make([]byte, 1024*1024*10+1)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	// 加密存储
	_, err = s.CClient.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	// 获取解密信息
	resp, err := s.CClient.Object.Head(context.Background(), name, nil)
	assert.Nil(s.T(), err, "HeadObject Failed")
	cipherKey := resp.Header.Get(coscrypto.COSClientSideEncryptionKey)
	cipherIV := resp.Header.Get(coscrypto.COSClientSideEncryptionStart)
	key, err := s.Master.Decrypt([]byte(cipherKey))
	assert.Nil(s.T(), err, "Master Decrypt Failed")
	iv, err := s.Master.Decrypt([]byte(cipherIV))
	assert.Nil(s.T(), err, "Master Decrypt Failed")

	// 正常读取
	resp, err = s.Client.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	encryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.NotEqual(s.T(), bytes.Compare(encryptedData, originData), 0, "encryptedData == originData")

	// 手动解密
	block, err := aes.NewCipher(key)
	assert.Nil(s.T(), err, "NewCipher Failed")
	decrypter := cipher.NewCTR(block, iv)
	decryptedData := make([]byte, len(originData))
	decrypter.XORKeyStream(decryptedData, encryptedData)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_Normal_10MB() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	originData := make([]byte, 1024*1024*10+1)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	// 加密存储
	_, err = s.CClient.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	// 解密读取
	resp, err := s.CClient.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_ZeroFile() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	// 加密存储
	_, err := s.CClient.Object.Put(context.Background(), name, bytes.NewReader([]byte("")), nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	// 解密读取
	resp, err := s.CClient.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), bytes.Compare([]byte(""), decryptedData), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_WithMetaData() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	originData := make([]byte, 1024*1024*10+1)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	m := md5.New()
	m.Write(originData)
	contentMD5 := m.Sum(nil)
	opt := &cos.ObjectPutOptions{
		&cos.ACLHeaderOptions{
			XCosACL: "private",
		},
		&cos.ObjectPutHeaderOptions{
			ContentLength: 1024*1024*10 + 1,
			ContentMD5:    base64.StdEncoding.EncodeToString(contentMD5),
			XCosMetaXXX:   &http.Header{},
		},
	}
	opt.XCosMetaXXX.Add("x-cos-meta-isEncrypted", "true")
	// 加密存储
	_, err = s.CClient.Object.Put(context.Background(), name, f, opt)
	assert.Nil(s.T(), err, "PutObject Failed")

	// 解密读取
	resp, err := s.CClient.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
	assert.Equal(s.T(), resp.Header.Get("x-cos-meta-isEncrypted"), "true", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionCekAlg), "AES/CTR/NoPadding", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionWrapAlg), "COS/KMS/Crypto", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionUnencryptedContentMD5), base64.StdEncoding.EncodeToString(contentMD5), "meta data isn't consistent")
	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_ByFile() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	filepath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filepath)
	assert.Nil(s.T(), err, "Create File Failed")
	defer os.Remove(filepath)

	originData := make([]byte, 1024*1024*10+1)
	_, err = rand.Read(originData)
	newfile.Write(originData)
	newfile.Close()

	m := md5.New()
	m.Write(originData)
	contentMD5 := m.Sum(nil)
	opt := &cos.ObjectPutOptions{
		&cos.ACLHeaderOptions{
			XCosACL: "private",
		},
		&cos.ObjectPutHeaderOptions{
			ContentLength: 1024*1024*10 + 1,
			ContentMD5:    base64.StdEncoding.EncodeToString(contentMD5),
			XCosMetaXXX:   &http.Header{},
		},
	}
	opt.XCosMetaXXX.Add("x-cos-meta-isEncrypted", "true")
	// 加密存储
	_, err = s.CClient.Object.PutFromFile(context.Background(), name, filepath, opt)
	assert.Nil(s.T(), err, "PutFromFile Failed")

	// 解密读取
	downfile := "downfile" + time.Now().Format(time.RFC3339)
	resp, err := s.CClient.Object.GetToFile(context.Background(), name, downfile, nil)
	assert.Nil(s.T(), err, "GetToFile Failed")
	assert.Equal(s.T(), resp.Header.Get("x-cos-meta-isEncrypted"), "true", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionCekAlg), "AES/CTR/NoPadding", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionWrapAlg), "COS/KMS/Crypto", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionUnencryptedContentMD5), base64.StdEncoding.EncodeToString(contentMD5), "meta data isn't consistent")

	fd, err := os.Open(downfile)
	assert.Nil(s.T(), err, "Open File Failed")
	defer os.Remove(downfile)
	defer fd.Close()
	m = md5.New()
	io.Copy(m, fd)
	downContentMD5 := m.Sum(nil)
	assert.Equal(s.T(), bytes.Compare(contentMD5, downContentMD5), 0, "decryptData != originData")
	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_DecryptWithNewClient_10MB() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	originData := make([]byte, 1024*1024*10+1)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	// 加密存储
	_, err = s.CClient.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	u, _ := url.Parse("https://" + kBucket + ".cos." + kRegion + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})
	{
		// 使用不同的MatDesc客户端读取, 期待错误
		material := make(map[string]string)
		material["desc"] = "cos crypto suite test 2"
		kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), kRegion)
		master, _ := coscrypto.CreateMasterKMS(kmsclient, os.Getenv("KMSID"), material)
		client := coscrypto.NewCryptoClient(c, master)
		resp, err := client.Object.Get(context.Background(), name, nil)
		assert.Nil(s.T(), resp, "Get Object Failed")
		assert.NotNil(s.T(), err, "Get Object Failed")
	}

	{
		// 使用相同的MatDesc客户端读取, 但KMSID不一样，期待正确，kms解密是不需要KMSID
		material := make(map[string]string)
		material["desc"] = "cos crypto suite test"
		kmsclient, _ := coscrypto.NewKMSClient(s.Client.GetCredential(), kRegion)
		master, _ := coscrypto.CreateMasterKMS(kmsclient, "KMSID", material)
		client := coscrypto.NewCryptoClient(c, master)
		resp, err := client.Object.Get(context.Background(), name, nil)
		assert.Nil(s.T(), err, "Get Object Failed")
		defer resp.Body.Close()
		decryptedData, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
	}

	{
		// 使用相同的MatDesc和KMSID客户端读取, 期待正确
		material := make(map[string]string)
		material["desc"] = "cos crypto suite test"
		kmsclient, _ := coscrypto.NewKMSClient(s.Client.GetCredential(), kRegion)
		master, _ := coscrypto.CreateMasterKMS(kmsclient, os.Getenv("KMSID"), material)
		client := coscrypto.NewCryptoClient(c, master)
		resp, err := client.Object.Get(context.Background(), name, nil)
		assert.Nil(s.T(), err, "Get Object Failed")
		defer resp.Body.Close()
		decryptedData, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
	}

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_RangeGet() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	contentLength := 1024*1024*10 + 1
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	// 加密存储
	_, err = s.CClient.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	// Range解密读取
	for i := 0; i < 10; i++ {
		math_rand.Seed(time.Now().UnixNano())
		rangeStart := math_rand.Intn(contentLength)
		rangeEnd := rangeStart + math_rand.Intn(contentLength-rangeStart)
		if rangeEnd == rangeStart || rangeStart >= contentLength-1 {
			continue
		}
		opt := &cos.ObjectGetOptions{
			Range: fmt.Sprintf("bytes=%v-%v", rangeStart, rangeEnd),
		}
		resp, err := s.CClient.Object.Get(context.Background(), name, opt)
		assert.Nil(s.T(), err, "GetObject Failed")
		defer resp.Body.Close()
		decryptedData, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(s.T(), bytes.Compare(originData[rangeStart:rangeEnd+1], decryptedData), 0, "decryptData != originData")
	}

	// 解密读取
	resp, err := s.CClient.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObject_WithListenerAndRange() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	contentLength := 1024*1024*10 + 1
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	// 加密存储
	popt := &cos.ObjectPutOptions{
		nil,
		&cos.ObjectPutHeaderOptions{
			Listener: &cos.DefaultProgressListener{},
		},
	}
	_, err = s.CClient.Object.Put(context.Background(), name, f, popt)
	assert.Nil(s.T(), err, "PutObject Failed")

	// Range解密读取
	for i := 0; i < 10; i++ {
		math_rand.Seed(time.Now().UnixNano())
		rangeStart := math_rand.Intn(contentLength)
		rangeEnd := rangeStart + math_rand.Intn(contentLength-rangeStart)
		if rangeEnd == rangeStart || rangeStart >= contentLength-1 {
			continue
		}
		opt := &cos.ObjectGetOptions{
			Range:    fmt.Sprintf("bytes=%v-%v", rangeStart, rangeEnd),
			Listener: &cos.DefaultProgressListener{},
		}
		resp, err := s.CClient.Object.Get(context.Background(), name, opt)
		assert.Nil(s.T(), err, "GetObject Failed")
		defer resp.Body.Close()
		decryptedData, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(s.T(), bytes.Compare(originData[rangeStart:rangeEnd+1], decryptedData), 0, "decryptData != originData")
	}
	// 解密读取
	opt := &cos.ObjectGetOptions{
		Listener: &cos.DefaultProgressListener{},
	}
	resp, err := s.CClient.Object.Get(context.Background(), name, opt)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func TestCosTestSuite(t *testing.T) {
	suite.Run(t, new(CosTestSuite))
}

func (s *CosTestSuite) TearDownSuite() {
}
