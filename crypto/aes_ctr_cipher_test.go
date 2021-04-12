package coscrypto_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5/crypto"
	"io/ioutil"
	math_rand "math/rand"
)

type EmptyMasterCipher struct{}

func (mc EmptyMasterCipher) Encrypt(b []byte) ([]byte, error) {
	return b, nil
}
func (mc EmptyMasterCipher) Decrypt(b []byte) ([]byte, error) {
	return b, nil
}
func (mc EmptyMasterCipher) GetWrapAlgorithm() string {
	return "Test/EmptyWrapAlgo"
}
func (mc EmptyMasterCipher) GetMatDesc() string {
	return "Empty Desc"
}

func (s *CosTestSuite) TestCryptoObjectService_EncryptAndDecrypt() {
	var masterCipher EmptyMasterCipher
	builder := coscrypto.CreateAesCtrBuilder(masterCipher)

	contentCipher, err := builder.ContentCipher()
	assert.Nil(s.T(), err, "CryptoObject.CreateAesCtrBuilder Failed")

	dataSize := math_rand.Int63n(1024 * 1024 * 32)
	originData := make([]byte, dataSize)
	rand.Read(originData)
	// 加密
	r1 := bytes.NewReader(originData)
	reader1, err := contentCipher.EncryptContent(r1)
	assert.Nil(s.T(), err, "CryptoObject.contentCipher.Encrypt Failed")
	encryptedData, err := ioutil.ReadAll(reader1)
	assert.Nil(s.T(), err, "CryptoObject.Read Failed")

	// 解密
	r2 := bytes.NewReader(encryptedData)
	reader2, err := contentCipher.DecryptContent(r2)
	decryptedData, err := ioutil.ReadAll(reader2)
	assert.Nil(s.T(), err, "CryptoObject.Read Failed")
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
}

func (s *CosTestSuite) TestCryptoObjectService_Encrypt() {
	var masterCipher EmptyMasterCipher
	builder := coscrypto.CreateAesCtrBuilder(masterCipher)

	contentCipher, err := builder.ContentCipher()
	assert.Nil(s.T(), err, "CryptoObject.CreateAesCtrBuilder Failed")

	dataSize := math_rand.Int63n(1024 * 1024 * 32)
	originData := make([]byte, dataSize)
	rand.Read(originData)

	// 加密
	r := bytes.NewReader(originData)
	reader, err := contentCipher.EncryptContent(r)
	assert.Nil(s.T(), err, "CryptoObject.contentCipher.Encrypt Failed")
	encryptedData, err := ioutil.ReadAll(reader)
	assert.Nil(s.T(), err, "CryptoObject.Read Failed")

	// 直接解密
	cd := contentCipher.GetCipherData()
	block, err := aes.NewCipher(cd.Key)
	assert.Nil(s.T(), err, "CryptoObject.NewCipher Failed")
	decrypter := cipher.NewCTR(block, cd.IV)
	decryptedData := make([]byte, len(originData))
	decrypter.XORKeyStream(decryptedData, encryptedData)
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
}

func (s *CosTestSuite) TestCryptoObjectService_Decrypt() {
	var masterCipher EmptyMasterCipher
	builder := coscrypto.CreateAesCtrBuilder(masterCipher)

	contentCipher, err := builder.ContentCipher()
	assert.Nil(s.T(), err, "CryptoObject.CreateAesCtrBuilder Failed")
	dataSize := math_rand.Int63n(1024 * 1024 * 32)
	originData := make([]byte, dataSize)
	rand.Read(originData)

	// 直接加密
	cd := contentCipher.GetCipherData()
	block, err := aes.NewCipher(cd.Key)
	assert.Nil(s.T(), err, "CryptoObject.NewCipher Failed")
	encrypter := cipher.NewCTR(block, cd.IV)
	encryptedData := make([]byte, len(originData))
	encrypter.XORKeyStream(encryptedData, originData)

	// 解密
	r := bytes.NewReader(encryptedData)
	reader, err := contentCipher.DecryptContent(r)
	assert.Nil(s.T(), err, "CryptoObject.contentCipher.Encrypt Failed")
	decryptedData, err := ioutil.ReadAll(reader)
	assert.Nil(s.T(), err, "CryptoObject.Read Failed")
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "decryptData != originData")
}
