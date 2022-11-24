package coscrypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	KMSEndPoint = "kms.tencentcloudapi.com"
)

type MasterKMSCipher struct {
	Client  *kms.Client
	KmsId   string
	MatDesc string
}

type KMSClientOptions = func(*profile.HttpProfile)

func KMSEndpoint(endpoint string) KMSClientOptions {
	return func(pf *profile.HttpProfile) {
		pf.Endpoint = endpoint
	}
}

func NewKMSClient(cred *cos.Credential, region string, opt ...KMSClientOptions) (*kms.Client, error) {
	if cred == nil {
		fmt.Errorf("credential is nil")
	}
	credential := common.NewTokenCredential(
		cred.SecretID,
		cred.SecretKey,
		cred.SessionToken,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = KMSEndPoint

	for _, fn := range opt {
		fn(cpf.HttpProfile)
	}
	client, err := kms.NewClient(credential, region, cpf)
	return client, err
}

func CreateMasterKMS(client *kms.Client, kmsId string, desc map[string]string) (MasterCipher, error) {
	if kmsId == "" || client == nil {
		return nil, fmt.Errorf("KMS ID is empty or kms client is nil")
	}
	var kmsCipher MasterKMSCipher
	var jdesc string
	if len(desc) > 0 {
		bs, err := json.Marshal(desc)
		if err != nil {
			return nil, err
		}
		jdesc = string(bs)
	}
	kmsCipher.Client = client
	kmsCipher.KmsId = kmsId
	kmsCipher.MatDesc = jdesc
	return &kmsCipher, nil
}

func (kc *MasterKMSCipher) Encrypt(plaintext []byte) ([]byte, error) {
	request := kms.NewEncryptRequest()
	request.KeyId = common.StringPtr(kc.KmsId)
	request.EncryptionContext = common.StringPtr(kc.MatDesc)
	request.Plaintext = common.StringPtr(base64.StdEncoding.EncodeToString(plaintext))
	resp, err := kc.Client.Encrypt(request)
	if err != nil {
		return nil, err
	}
	return []byte(*resp.Response.CiphertextBlob), nil
}

func (kc *MasterKMSCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	request := kms.NewDecryptRequest()
	request.CiphertextBlob = common.StringPtr(string(ciphertext))
	request.EncryptionContext = common.StringPtr(kc.MatDesc)
	resp, err := kc.Client.Decrypt(request)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(*resp.Response.Plaintext)
}

func (kc *MasterKMSCipher) GetWrapAlgorithm() string {
	return CosKmsCryptoWrap
}

func (kc *MasterKMSCipher) GetMatDesc() string {
	return kc.MatDesc
}
