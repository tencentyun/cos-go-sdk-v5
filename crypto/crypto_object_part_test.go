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
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/crypto"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

func (s *CosTestSuite) TestMultiUpload_Normal() {
	name := "test/ObjectPut" + time.Now().Format(time.RFC3339)
	contentLength := int64(1024*1024*10 + 1)
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)

	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		PartSize: (contentLength / 16 / 3) * 16,
	}
	v, _, err := s.CClient.Object.InitiateMultipartUpload(context.Background(), name, nil, &cryptoCtx)
	assert.Nil(s.T(), err, "Init Failed")
	chunks, _, err := cos.SplitSizeIntoChunks(contentLength, cryptoCtx.PartSize)
	assert.Nil(s.T(), err, "Split Failed")
	optcom := &cos.CompleteMultipartUploadOptions{}
	for _, chunk := range chunks {
		opt := &cos.ObjectUploadPartOptions{
			ContentLength: chunk.Size,
		}
		f := bytes.NewReader(originData[chunk.OffSet : chunk.OffSet+chunk.Size])
		resp, err := s.CClient.Object.UploadPart(context.Background(), name, v.UploadID, chunk.Number, io.LimitReader(f, chunk.Size), opt, &cryptoCtx)
		assert.Nil(s.T(), err, "UploadPart failed")
		optcom.Parts = append(optcom.Parts, cos.Object{
			PartNumber: chunk.Number, ETag: resp.Header.Get("ETag"),
		})
	}
	_, _, err = s.CClient.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	assert.Nil(s.T(), err, "Complete Failed")

	resp, err := s.CClient.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestMultiUpload_DecryptWithKey() {
	name := "test/ObjectPut" + time.Now().Format(time.RFC3339)
	contentLength := int64(1024*1024*10 + 1)
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)
	f := bytes.NewReader(originData)

	// 分块上传
	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		PartSize: (contentLength / 16 / 3) * 16,
	}
	v, _, err := s.CClient.Object.InitiateMultipartUpload(context.Background(), name, nil, &cryptoCtx)
	assert.Nil(s.T(), err, "Init Failed")
	chunks, _, err := cos.SplitSizeIntoChunks(contentLength, cryptoCtx.PartSize)
	assert.Nil(s.T(), err, "Split Failed")
	optcom := &cos.CompleteMultipartUploadOptions{}
	for _, chunk := range chunks {
		opt := &cos.ObjectUploadPartOptions{
			ContentLength: chunk.Size,
			Listener:      &cos.DefaultProgressListener{},
		}
		resp, err := s.CClient.Object.UploadPart(context.Background(), name, v.UploadID, chunk.Number, io.LimitReader(f, chunk.Size), opt, &cryptoCtx)
		assert.Nil(s.T(), err, "UploadPart failed")
		optcom.Parts = append(optcom.Parts, cos.Object{
			PartNumber: chunk.Number, ETag: resp.Header.Get("ETag"),
		})
	}
	_, _, err = s.CClient.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	assert.Nil(s.T(), err, "Complete Failed")

	// 正常读取
	resp, err := s.Client.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	encryptedData, _ := ioutil.ReadAll(resp.Body)
	assert.NotEqual(s.T(), bytes.Compare(encryptedData, originData), 0, "encryptedData == originData")

	// 获取解密信息
	resp, err = s.CClient.Object.Head(context.Background(), name, nil)
	assert.Nil(s.T(), err, "HeadObject Failed")
	cipherKey := resp.Header.Get(coscrypto.COSClientSideEncryptionKey)
	cipherKeybs, err := base64.StdEncoding.DecodeString(cipherKey)
	assert.Nil(s.T(), err, "base64 Decode Failed")
	cipherIV := resp.Header.Get(coscrypto.COSClientSideEncryptionStart)
	cipherIVbs, err := base64.StdEncoding.DecodeString(cipherIV)
	assert.Nil(s.T(), err, "base64 Decode Failed")
	key, err := s.Master.Decrypt(cipherKeybs)
	assert.Nil(s.T(), err, "Master Decrypt Failed")
	iv, err := s.Master.Decrypt(cipherIVbs)
	assert.Nil(s.T(), err, "Master Decrypt Failed")

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

func (s *CosTestSuite) TestMultiUpload_PutFromFile() {
	name := "test/ObjectPut" + time.Now().Format(time.RFC3339)
	filepath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filepath)
	assert.Nil(s.T(), err, "Create File Failed")
	defer os.Remove(filepath)

	contentLength := int64(1024*1024*10 + 1)
	originData := make([]byte, contentLength)
	_, err = rand.Read(originData)
	newfile.Write(originData)
	newfile.Close()

	m := md5.New()
	m.Write(originData)
	contentMD5 := m.Sum(nil)
	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		PartSize: (contentLength / 16 / 3) * 16,
	}
	v, _, err := s.CClient.Object.InitiateMultipartUpload(context.Background(), name, nil, &cryptoCtx)
	assert.Nil(s.T(), err, "Init Failed")
	_, chunks, _, err := cos.SplitFileIntoChunks(filepath, cryptoCtx.PartSize)
	assert.Nil(s.T(), err, "Split Failed")
	optcom := &cos.CompleteMultipartUploadOptions{}
	var wg sync.WaitGroup
	var mtx sync.Mutex
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chk cos.Chunk) {
			defer wg.Done()
			fd, err := os.Open(filepath)
			assert.Nil(s.T(), err, "Open File Failed")
			opt := &cos.ObjectUploadPartOptions{
				ContentLength: chk.Size,
			}
			fd.Seek(chk.OffSet, os.SEEK_SET)
			resp, err := s.CClient.Object.UploadPart(context.Background(), name, v.UploadID, chk.Number, io.LimitReader(fd, chk.Size), opt, &cryptoCtx)
			assert.Nil(s.T(), err, "UploadPart failed")
			mtx.Lock()
			optcom.Parts = append(optcom.Parts, cos.Object{
				PartNumber: chk.Number, ETag: resp.Header.Get("ETag"),
			})
			mtx.Unlock()
		}(chunk)
	}
	wg.Wait()
	sort.Sort(cos.ObjectList(optcom.Parts))
	_, _, err = s.CClient.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	assert.Nil(s.T(), err, "Complete Failed")

	downfile := "downfile" + time.Now().Format(time.RFC3339)
	_, err = s.CClient.Object.GetToFile(context.Background(), name, downfile, nil)
	assert.Nil(s.T(), err, "GetObject Failed")

	m = md5.New()
	fd, err := os.Open(downfile)
	assert.Nil(s.T(), err, "Open File Failed")
	defer os.Remove(downfile)
	defer fd.Close()
	io.Copy(m, fd)
	downContentMD5 := m.Sum(nil)
	assert.Equal(s.T(), bytes.Compare(contentMD5, downContentMD5), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestMultiUpload_GetWithRange() {
	name := "test/ObjectPut" + time.Now().Format(time.RFC3339)
	filepath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filepath)
	assert.Nil(s.T(), err, "Create File Failed")
	defer os.Remove(filepath)

	contentLength := int64(1024*1024*10 + 1)
	originData := make([]byte, contentLength)
	_, err = rand.Read(originData)
	newfile.Write(originData)
	newfile.Close()

	m := md5.New()
	m.Write(originData)
	contentMD5 := m.Sum(nil)
	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		PartSize: (contentLength / 16 / 3) * 16,
	}
	iniopt := &cos.InitiateMultipartUploadOptions{
		&cos.ACLHeaderOptions{
			XCosACL: "private",
		},
		&cos.ObjectPutHeaderOptions{
			ContentMD5:  base64.StdEncoding.EncodeToString(contentMD5),
			XCosMetaXXX: &http.Header{},
		},
	}
	iniopt.XCosMetaXXX.Add("x-cos-meta-isEncrypted", "true")

	v, _, err := s.CClient.Object.InitiateMultipartUpload(context.Background(), name, iniopt, &cryptoCtx)
	assert.Nil(s.T(), err, "Init Failed")
	_, chunks, _, err := cos.SplitFileIntoChunks(filepath, cryptoCtx.PartSize)
	assert.Nil(s.T(), err, "Split Failed")
	optcom := &cos.CompleteMultipartUploadOptions{}
	var wg sync.WaitGroup
	var mtx sync.Mutex
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chk cos.Chunk) {
			defer wg.Done()
			fd, err := os.Open(filepath)
			assert.Nil(s.T(), err, "Open File Failed")
			opt := &cos.ObjectUploadPartOptions{
				ContentLength: chk.Size,
			}
			fd.Seek(chk.OffSet, os.SEEK_SET)
			resp, err := s.CClient.Object.UploadPart(context.Background(), name, v.UploadID, chk.Number, io.LimitReader(fd, chk.Size), opt, &cryptoCtx)
			assert.Nil(s.T(), err, "UploadPart failed")
			mtx.Lock()
			optcom.Parts = append(optcom.Parts, cos.Object{
				PartNumber: chk.Number, ETag: resp.Header.Get("ETag"),
			})
			mtx.Unlock()
		}(chunk)
	}
	wg.Wait()
	sort.Sort(cos.ObjectList(optcom.Parts))
	_, _, err = s.CClient.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	assert.Nil(s.T(), err, "Complete Failed")

	// Range解密读取
	for i := 0; i < 10; i++ {
		math_rand.Seed(time.Now().UnixNano())
		rangeStart := math_rand.Int63n(contentLength)
		rangeEnd := rangeStart + math_rand.Int63n(contentLength-rangeStart)
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

	opt := &cos.ObjectGetOptions{
		Listener: &cos.DefaultProgressListener{},
	}
	resp, err := s.CClient.Object.Get(context.Background(), name, opt)
	assert.Nil(s.T(), err, "GetObject Failed")
	assert.Equal(s.T(), resp.Header.Get("x-cos-meta-isEncrypted"), "true", "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionCekAlg), coscrypto.AesCtrAlgorithm, "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionWrapAlg), coscrypto.CosKmsCryptoWrap, "meta data isn't consistent")
	assert.Equal(s.T(), resp.Header.Get(coscrypto.COSClientSideEncryptionUnencryptedContentMD5), base64.StdEncoding.EncodeToString(contentMD5), "meta data isn't consistent")
	defer resp.Body.Close()
	decryptedData, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")

	_, err = s.CClient.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestMultiUpload_GetWithNewClient() {
	name := "test/ObjectPut" + time.Now().Format(time.RFC3339)
	contentLength := int64(1024*1024*10 + 1)
	originData := make([]byte, contentLength)
	_, err := rand.Read(originData)

	cryptoCtx := coscrypto.CryptoContext{
		DataSize: contentLength,
		PartSize: (contentLength / 16 / 3) * 16,
	}
	v, _, err := s.CClient.Object.InitiateMultipartUpload(context.Background(), name, nil, &cryptoCtx)
	assert.Nil(s.T(), err, "Init Failed")
	chunks, _, err := cos.SplitSizeIntoChunks(contentLength, cryptoCtx.PartSize)
	assert.Nil(s.T(), err, "Split Failed")
	optcom := &cos.CompleteMultipartUploadOptions{}
	for _, chunk := range chunks {
		opt := &cos.ObjectUploadPartOptions{
			ContentLength: chunk.Size,
		}
		f := bytes.NewReader(originData[chunk.OffSet : chunk.OffSet+chunk.Size])
		resp, err := s.CClient.Object.UploadPart(context.Background(), name, v.UploadID, chunk.Number, io.LimitReader(f, chunk.Size), opt, &cryptoCtx)
		assert.Nil(s.T(), err, "UploadPart failed")
		optcom.Parts = append(optcom.Parts, cos.Object{
			PartNumber: chunk.Number, ETag: resp.Header.Get("ETag"),
		})
	}
	_, _, err = s.CClient.Object.CompleteMultipartUpload(context.Background(), name, v.UploadID, optcom)
	assert.Nil(s.T(), err, "Complete Failed")

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
		kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), kRegion)
		master, _ := coscrypto.CreateMasterKMS(kmsclient, "KMSID", material)
		client := coscrypto.NewCryptoClient(c, master)
		resp, err := client.Object.Get(context.Background(), name, nil)
		assert.Nil(s.T(), err, "Get Object Failed")
		defer resp.Body.Close()
		decryptedData, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
	}

	{
		// 使用相同的MatDesc客户端读取, 地域不一样，期待错误
		material := make(map[string]string)
		material["desc"] = "cos crypto suite test"
		diffRegion := "ap-shanghai"
		if diffRegion == kRegion {
			diffRegion = "ap-guangzhou"
		}
		kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), diffRegion)
		master, _ := coscrypto.CreateMasterKMS(kmsclient, "KMSID", material)
		client := coscrypto.NewCryptoClient(c, master)
		resp, err := client.Object.Get(context.Background(), name, nil)
		assert.Nil(s.T(), resp, "Get Object Failed")
		assert.NotNil(s.T(), err, "Get Object Failed")
	}

	{
		// 使用相同的MatDesc和KMSID客户端读取, 期待正确
		material := make(map[string]string)
		material["desc"] = "cos crypto suite test"
		kmsclient, _ := coscrypto.NewKMSClient(c.GetCredential(), kRegion)
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
