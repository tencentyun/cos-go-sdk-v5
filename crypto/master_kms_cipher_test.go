package coscrypto_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/crypto"
	"os"
)

func (s *CosTestSuite) TestMasterKmsCipher_TestKmsClient() {
	kmsclient, _ := coscrypto.NewKMSClient(&cos.Credential{
		SecretID:  os.Getenv("COS_SECRETID"),
		SecretKey: os.Getenv("COS_SECRETKEY"),
	}, kRegion)

	originData := make([]byte, 1024)
	_, err := rand.Read(originData)

	ctx := make(map[string]string)
	ctx["desc"] = string(originData[:10])
	bs, _ := json.Marshal(ctx)
	ctxJson := string(bs)
	enReq := kms.NewEncryptRequest()
	enReq.KeyId = common.StringPtr(os.Getenv("KMSID"))
	enReq.EncryptionContext = common.StringPtr(ctxJson)
	enReq.Plaintext = common.StringPtr(base64.StdEncoding.EncodeToString(originData))
	enResp, err := kmsclient.Encrypt(enReq)
	assert.Nil(s.T(), err, "Encrypt Failed")
	encryptedData := []byte(*enResp.Response.CiphertextBlob)

	deReq := kms.NewDecryptRequest()
	deReq.CiphertextBlob = common.StringPtr(string(encryptedData))
	deReq.EncryptionContext = common.StringPtr(ctxJson)
	deResp, err := kmsclient.Decrypt(deReq)
	assert.Nil(s.T(), err, "Decrypt Failed")
	decryptedData, err := base64.StdEncoding.DecodeString(*deResp.Response.Plaintext)
	assert.Nil(s.T(), err, "base64 Decode Failed")
	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "originData != decryptedData")
}

func (s *CosTestSuite) TestMasterKmsCipher_TestNormal() {
	kmsclient, _ := coscrypto.NewKMSClient(&cos.Credential{
		SecretID:  os.Getenv("COS_SECRETID"),
		SecretKey: os.Getenv("COS_SECRETKEY"),
	}, kRegion)

	desc := make(map[string]string)
	desc["test"] = "TestMasterKmsCipher_TestNormal"
	master, err := coscrypto.CreateMasterKMS(kmsclient, os.Getenv("KMSID"), desc)
	assert.Nil(s.T(), err, "CreateMasterKMS Failed")

	originData := make([]byte, 1024)
	_, err = rand.Read(originData)

	encryptedData, err := master.Encrypt(originData)
	assert.Nil(s.T(), err, "Encrypt Failed")

	decryptedData, err := master.Decrypt(encryptedData)
	assert.Nil(s.T(), err, "Decrypt Failed")

	assert.Equal(s.T(), bytes.Compare(originData, decryptedData), 0, "originData != decryptedData")
}

func (s *CosTestSuite) TestMasterKmsCipher_TestError() {
	kmsclient, _ := coscrypto.NewKMSClient(&cos.Credential{
		SecretID:  os.Getenv("COS_SECRETID"),
		SecretKey: os.Getenv("COS_SECRETKEY"),
	}, kRegion)

	desc := make(map[string]string)
	desc["test"] = "TestMasterKmsCipher_TestNormal"
	master, err := coscrypto.CreateMasterKMS(kmsclient, "ErrorKMSID", desc)
	assert.Nil(s.T(), err, "CreateMasterKMS Failed")

	originData := make([]byte, 1024)
	_, err = rand.Read(originData)

	encryptedData, err := master.Encrypt(originData)
	assert.NotNil(s.T(), err, "Encrypt Failed")

	decryptedData, err := master.Decrypt(encryptedData)
	assert.NotNil(s.T(), err, "Decrypt Failed")

	assert.NotEqual(s.T(), bytes.Compare(originData, decryptedData), 0, "originData != decryptedData")
}
