package coscrypto

import (
	"io"
)

const (
	aesKeySize = 32
	ivSize     = 16
)

type aesCtrCipherBuilder struct {
	MasterCipher MasterCipher
}

type aesCtrCipher struct {
	CipherData CipherData
	Cipher     Cipher
}

func CreateAesCtrBuilder(cipher MasterCipher) ContentCipherBuilder {
	return aesCtrCipherBuilder{MasterCipher: cipher}
}

func (builder aesCtrCipherBuilder) createCipherData() (CipherData, error) {
	var cd CipherData
	var err error
	err = cd.RandomKeyIv(aesKeySize, ivSize)
	if err != nil {
		return cd, err
	}

	cd.WrapAlgorithm = builder.MasterCipher.GetWrapAlgorithm()
	cd.CEKAlgorithm = AesCtrAlgorithm
	cd.MatDesc = builder.MasterCipher.GetMatDesc()

	// EncryptedKey
	cd.EncryptedKey, err = builder.MasterCipher.Encrypt(cd.Key)
	if err != nil {
		return cd, err
	}

	// EncryptedIV
	cd.EncryptedIV, err = builder.MasterCipher.Encrypt(cd.IV)
	if err != nil {
		return cd, err
	}

	return cd, nil
}

func (builder aesCtrCipherBuilder) contentCipherCD(cd CipherData) (ContentCipher, error) {
	cipher, err := newAesCtr(cd)
	if err != nil {
		return nil, err
	}

	return &aesCtrCipher{
		CipherData: cd,
		Cipher:     cipher,
	}, nil
}

func (builder aesCtrCipherBuilder) ContentCipher() (ContentCipher, error) {
	cd, err := builder.createCipherData()
	if err != nil {
		return nil, err
	}
	return builder.contentCipherCD(cd)
}

func (builder aesCtrCipherBuilder) ContentCipherEnv(envelope Envelope) (ContentCipher, error) {
	var cd CipherData
	cd.EncryptedKey = make([]byte, len(envelope.CipherKey))
	copy(cd.EncryptedKey, []byte(envelope.CipherKey))

	plainKey, err := builder.MasterCipher.Decrypt([]byte(envelope.CipherKey))
	if err != nil {
		return nil, err
	}
	cd.Key = make([]byte, len(plainKey))
	copy(cd.Key, plainKey)

	cd.EncryptedIV = make([]byte, len(envelope.IV))
	copy(cd.EncryptedIV, []byte(envelope.IV))

	plainIV, err := builder.MasterCipher.Decrypt([]byte(envelope.IV))
	if err != nil {
		return nil, err
	}

	cd.IV = make([]byte, len(plainIV))
	copy(cd.IV, plainIV)

	cd.MatDesc = envelope.MatDesc
	cd.WrapAlgorithm = envelope.WrapAlg
	cd.CEKAlgorithm = envelope.CEKAlg

	return builder.contentCipherCD(cd)
}

func (builder aesCtrCipherBuilder) GetMatDesc() string {
	return builder.MasterCipher.GetMatDesc()
}

func (cc *aesCtrCipher) EncryptContent(src io.Reader) (io.ReadCloser, error) {
	reader := cc.Cipher.Encrypt(src)
	return &CryptoEncrypter{Body: src, Encrypter: reader}, nil
}

func (cc *aesCtrCipher) DecryptContent(src io.Reader) (io.ReadCloser, error) {
	reader := cc.Cipher.Decrypt(src)
	return &CryptoDecrypter{Body: src, Decrypter: reader}, nil
}

func (cc *aesCtrCipher) GetCipherData() *CipherData {
	return &(cc.CipherData)
}

func (cc *aesCtrCipher) GetEncryptedLen(plainTextLen int64) int64 {
	return plainTextLen
}

func (cc *aesCtrCipher) GetAlignLen() int {
	return len(cc.CipherData.IV)
}

func (cc *aesCtrCipher) Clone(cd CipherData) (ContentCipher, error) {
	cipher, err := newAesCtr(cd)
	if err != nil {
		return nil, err
	}

	return &aesCtrCipher{
		CipherData: cd,
		Cipher:     cipher,
	}, nil
}
