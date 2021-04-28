package coscrypto

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	math_rand "math/rand"
	"time"
)

const (
	COSClientSideEncryptionKey                      string = "x-cos-meta-client-side-encryption-key"
	COSClientSideEncryptionStart                           = "x-cos-meta-client-side-encryption-start"
	COSClientSideEncryptionCekAlg                          = "x-cos-meta-client-side-encryption-cek-alg"
	COSClientSideEncryptionWrapAlg                         = "x-cos-meta-client-side-encryption-wrap-alg"
	COSClientSideEncryptionMatDesc                         = "x-cos-meta-client-side-encryption-matdesc"
	COSClientSideEncryptionUnencryptedContentLength        = "x-cos-meta-client-side-encryption-unencrypted-content-length"
	COSClientSideEncryptionUnencryptedContentMD5           = "x-cos-meta-client-side-encryption-unencrypted-content-md5"
	COSClientSideEncryptionDataSize                        = "x-cos-meta-client-side-encryption-data-size"
	COSClientSideEncryptionPartSize                        = "x-cos-meta-client-side-encryption-part-size"
	UserAgent                                              = "User-Agent"
)

const (
	CosKmsCryptoWrap   = "KMS/TencentCloud"
	AesCtrAlgorithm    = "AES/CTR/NoPadding"
	EncryptionUaSuffix = "COSEncryptionClient"
)

type MasterCipher interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	GetWrapAlgorithm() string
	GetMatDesc() string
}

type ContentCipherBuilder interface {
	ContentCipher() (ContentCipher, error)
	ContentCipherEnv(Envelope) (ContentCipher, error)
	GetMatDesc() string
}

type ContentCipher interface {
	EncryptContent(io.Reader) (io.ReadCloser, error)
	DecryptContent(io.Reader) (io.ReadCloser, error)
	Clone(cd CipherData) (ContentCipher, error)
	GetEncryptedLen(int64) int64
	GetCipherData() *CipherData
	GetAlignLen() int
}

type Envelope struct {
	IV                    string
	CipherKey             string
	MatDesc               string
	WrapAlg               string
	CEKAlg                string
	UnencryptedMD5        string
	UnencryptedContentLen string
}

func (el Envelope) IsValid() bool {
	return len(el.IV) > 0 &&
		len(el.CipherKey) > 0 &&
		len(el.WrapAlg) > 0 &&
		len(el.CEKAlg) > 0
}

func (el Envelope) String() string {
	return fmt.Sprintf("IV=%s&CipherKey=%s&WrapAlg=%s&CEKAlg=%s", el.IV, el.CipherKey, el.WrapAlg, el.CEKAlg)
}

type CipherData struct {
	IV            []byte
	Key           []byte
	MatDesc       string
	WrapAlgorithm string
	CEKAlgorithm  string
	EncryptedIV   []byte
	EncryptedKey  []byte
}

func (cd *CipherData) RandomKeyIv(keyLen int, ivLen int) error {
	math_rand.Seed(time.Now().UnixNano())

	// Key
	cd.Key = make([]byte, keyLen)
	if _, err := io.ReadFull(rand.Reader, cd.Key); err != nil {
		return err
	}

	// sizeof uint64
	if ivLen < 8 {
		return fmt.Errorf("ivLen:%d less than 8", ivLen)
	}

	// IV: | nonce: 8 bytes | Serial number: 8 bytes |
	cd.IV = make([]byte, ivLen)
	if _, err := io.ReadFull(rand.Reader, cd.IV[0:ivLen-8]); err != nil {
		return err
	}

	// only use 4 byte,in order not to overflow when SeekIV()
	randNumber := math_rand.Uint32()
	cd.SetIV(uint64(randNumber))
	return nil
}

func (cd *CipherData) SetIV(iv uint64) {
	ivLen := len(cd.IV)
	binary.BigEndian.PutUint64(cd.IV[ivLen-8:], iv)
}

func (cd *CipherData) GetIV() uint64 {
	ivLen := len(cd.IV)
	return binary.BigEndian.Uint64(cd.IV[ivLen-8:])
}

func (cd *CipherData) SeekIV(startPos uint64) {
	cd.SetIV(cd.GetIV() + startPos/uint64(len(cd.IV)))
}

func (cd *CipherData) Clone() CipherData {
	var cloneCd CipherData
	cloneCd = *cd

	cloneCd.Key = make([]byte, len(cd.Key))
	copy(cloneCd.Key, cd.Key)

	cloneCd.IV = make([]byte, len(cd.IV))
	copy(cloneCd.IV, cd.IV)

	cloneCd.EncryptedIV = make([]byte, len(cd.EncryptedIV))
	copy(cloneCd.EncryptedIV, cd.EncryptedIV)

	cloneCd.EncryptedKey = make([]byte, len(cd.EncryptedKey))
	copy(cloneCd.EncryptedKey, cd.EncryptedKey)

	return cloneCd
}
