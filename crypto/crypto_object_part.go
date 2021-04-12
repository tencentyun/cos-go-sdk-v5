package coscrypto

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"strconv"
)

type CryptoContext struct {
	DataSize      int64
	PartSize      int64
	ContentCipher ContentCipher
}

func partSizeIsValid(partSize int64, alignLen int64) bool {
	if partSize%alignLen == 0 {
		return true
	}
	return false
}

func (s *CryptoObjectService) InitiateMultipartUpload(ctx context.Context, name string, opt *cos.InitiateMultipartUploadOptions, cryptoCtx *CryptoContext) (*cos.InitiateMultipartUploadResult, *cos.Response, error) {
	contentCipher, err := s.cryptoClient.ContentCipherBuilder.ContentCipher()
	if err != nil {
		return nil, nil, err
	}
	if !partSizeIsValid(cryptoCtx.PartSize, int64(contentCipher.GetAlignLen())) {
		return nil, nil, fmt.Errorf("PartSize is invalid, it should be %v aligned", contentCipher.GetAlignLen())
	}
	// 添加自定义头部
	cryptoCtx.ContentCipher = contentCipher
	opt = cos.CloneInitiateMultipartUploadOptions(opt)
	if opt.XOptionHeader == nil {
		opt.XOptionHeader = &http.Header{}
	}
	if opt.ContentMD5 != "" {
		opt.XOptionHeader.Add(COSClientSideEncryptionUnencryptedContentMD5, opt.ContentMD5)
		opt.ContentMD5 = ""
	}
	if cryptoCtx.DataSize > 0 {
		opt.XOptionHeader.Add(COSClientSideEncryptionDataSize, strconv.FormatInt(cryptoCtx.DataSize, 10))
	}
	opt.XOptionHeader.Add(COSClientSideEncryptionPartSize, strconv.FormatInt(cryptoCtx.PartSize, 10))
	addCryptoHeaders(opt.XOptionHeader, contentCipher.GetCipherData())

	return s.ObjectService.InitiateMultipartUpload(ctx, name, opt)
}

func (s *CryptoObjectService) UploadPart(ctx context.Context, name, uploadID string, partNumber int, r io.Reader, opt *cos.ObjectUploadPartOptions, cryptoCtx *CryptoContext) (*cos.Response, error) {
	if cryptoCtx.PartSize == 0 {
		return nil, fmt.Errorf("CryptoContext's PartSize is zero")
	}
	opt = cos.CloneObjectUploadPartOptions(opt)
	if cryptoCtx.ContentCipher == nil {
		return nil, fmt.Errorf("ContentCipher is nil, Please call the InitiateMultipartUpload")
	}
	totalBytes, err := cos.GetReaderLen(r)
	if err == nil {
		// 与 go http 保持一致, 非bytes.Buffer/bytes.Reader/strings.Reader需用户指定ContentLength
		if opt != nil && opt.ContentLength == 0 && cos.IsLenReader(r) {
			opt.ContentLength = totalBytes
		}
	}
	cd := cryptoCtx.ContentCipher.GetCipherData().Clone()
	cd.SeekIV(uint64(partNumber-1) * uint64(cryptoCtx.PartSize))
	cc, err := cryptoCtx.ContentCipher.Clone(cd)
	opt.ContentLength = cc.GetEncryptedLen(opt.ContentLength)
	if err != nil {
		return nil, err
	}
	reader, err := cc.EncryptContent(r)
	if err != nil {
		return nil, err
	}
	return s.ObjectService.UploadPart(ctx, name, uploadID, partNumber, reader, opt)
}
