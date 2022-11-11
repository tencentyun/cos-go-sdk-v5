package coscrypto

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"os"
	"strconv"
)

type CryptoObjectService struct {
	*cos.ObjectService
	cryptoClient *CryptoClient
}

type CryptoClient struct {
	*cos.Client
	Object               *CryptoObjectService
	ContentCipherBuilder ContentCipherBuilder

	userAgent string
}

func NewCryptoClient(client *cos.Client, masterCipher MasterCipher) *CryptoClient {
	cc := &CryptoClient{
		Client: client,
		Object: &CryptoObjectService{
			client.Object,
			nil,
		},
		ContentCipherBuilder: CreateAesCtrBuilder(masterCipher),
	}
	cc.userAgent = cc.Client.UserAgent + "/" + EncryptionUaSuffix
	cc.Object.cryptoClient = cc

	return cc
}

func (s *CryptoObjectService) Put(ctx context.Context, name string, r io.Reader, opt *cos.ObjectPutOptions) (*cos.Response, error) {
	cc, err := s.cryptoClient.ContentCipherBuilder.ContentCipher()
	if err != nil {
		return nil, err
	}
	reader, err := cc.EncryptContent(r)
	if err != nil {
		return nil, err
	}
	opt = cos.CloneObjectPutOptions(opt)
	totalBytes, err := cos.GetReaderLen(r)
	if err != nil && opt != nil && opt.Listener != nil && opt.ContentLength == 0 {
		return nil, err
	}
	if err == nil {
		if opt != nil && opt.ContentLength == 0 {
			// 如果未设置Listener, 非bytes.Buffer/bytes.Reader/strings.Reader/os.File由用户指定Contength
			if opt.Listener != nil || cos.IsLenReader(r) {
				opt.ContentLength = totalBytes
			}
		}
	}
	if opt.XOptionHeader == nil {
		opt.XOptionHeader = &http.Header{}
	}
	if opt.ContentMD5 != "" {
		opt.XOptionHeader.Add(COSClientSideEncryptionUnencryptedContentMD5, opt.ContentMD5)
		opt.ContentMD5 = ""
	}
	if opt.ContentLength != 0 {
		opt.XOptionHeader.Add(COSClientSideEncryptionUnencryptedContentLength, strconv.FormatInt(opt.ContentLength, 10))
		opt.ContentLength = cc.GetEncryptedLen(opt.ContentLength)
	}
	opt.XOptionHeader.Add(UserAgent, s.cryptoClient.userAgent)
	addCryptoHeaders(opt.XOptionHeader, cc.GetCipherData())

	return s.ObjectService.Put(ctx, name, reader, opt)
}

func (s *CryptoObjectService) PutFromFile(ctx context.Context, name, filePath string, opt *cos.ObjectPutOptions) (resp *cos.Response, err error) {
	nr := 0
	for nr < 3 {
		fd, e := os.Open(filePath)
		if e != nil {
			err = e
			return
		}
		resp, err = s.Put(ctx, name, fd, opt)
		if err != nil {
			nr++
			fd.Close()
			continue
		}
		fd.Close()
		break
	}
	return
}

func (s *CryptoObjectService) Get(ctx context.Context, name string, opt *cos.ObjectGetOptions, id ...string) (*cos.Response, error) {
	meta, err := s.ObjectService.Head(ctx, name, nil, id...)
	if err != nil {
		return meta, err
	}
	_isEncrypted := isEncrypted(&meta.Header)
	if !_isEncrypted {
		return s.ObjectService.Get(ctx, name, opt, id...)
	}

	envelope, err := getEnvelopeFromHeader(&meta.Header)
	if err != nil {
		return nil, err
	}
	if !envelope.IsValid() {
		return nil, fmt.Errorf("get envelope from header failed, object:%v", name)
	}
	encryptMatDesc := s.cryptoClient.ContentCipherBuilder.GetMatDesc()
	if envelope.MatDesc != encryptMatDesc {
		return nil, fmt.Errorf("provided master cipher error, want:%v, return:%v, object:%v", encryptMatDesc, envelope.MatDesc, name)
	}

	cc, err := s.cryptoClient.ContentCipherBuilder.ContentCipherEnv(envelope)
	if err != nil {
		return nil, fmt.Errorf("get content cipher from envelope failed: %v, object:%v", err, name)
	}

	opt = cos.CloneObjectGetOptions(opt)
	if opt.XOptionHeader == nil {
		opt.XOptionHeader = &http.Header{}
	}
	optRange, err := cos.GetRangeOptions(opt)
	if err != nil {
		return nil, err
	}
	discardAlignLen := int64(0)
	// Range请求
	if optRange != nil && optRange.HasStart {
		// 加密block对齐
		adjustStart := adjustRangeStart(optRange.Start, int64(cc.GetAlignLen()))
		discardAlignLen = optRange.Start - adjustStart
		if discardAlignLen > 0 {
			optRange.Start = adjustStart
			opt.Range = cos.FormatRangeOptions(optRange)
		}

		cd := cc.GetCipherData().Clone()
		cd.SeekIV(uint64(adjustStart))
		cc, err = cc.Clone(cd)
		if err != nil {
			return nil, fmt.Errorf("ContentCipher Clone failed:%v, bject:%v", err, name)
		}
	}
	opt.XOptionHeader.Add(UserAgent, s.cryptoClient.userAgent)
	resp, err := s.ObjectService.Get(ctx, name, opt, id...)
	if err != nil {
		return resp, err
	}
	resp.Body, err = cc.DecryptContent(resp.Body)
	if err != nil {
		return resp, err
	}
	// 抛弃多读取的数据
	if discardAlignLen > 0 {
		resp.Body = &cos.DiscardReadCloser{
			RC:      resp.Body,
			Discard: int(discardAlignLen),
		}
	}
	return resp, err
}

func (s *CryptoObjectService) GetToFile(ctx context.Context, name, localpath string, opt *cos.ObjectGetOptions, id ...string) (*cos.Response, error) {
	resp, err := s.Get(ctx, name, opt, id...)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	// If file exist, overwrite it
	fd, err := os.OpenFile(localpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return resp, err
	}

	_, err = io.Copy(fd, resp.Body)
	fd.Close()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *CryptoObjectService) MultiUpload(ctx context.Context, name string, filepath string, opt *cos.MultiUploadOptions) (*cos.CompleteMultipartUploadResult, *cos.Response, error) {
	return s.Upload(ctx, name, filepath, opt)
}

func (s *CryptoObjectService) Upload(ctx context.Context, name string, filepath string, opt *cos.MultiUploadOptions) (*cos.CompleteMultipartUploadResult, *cos.Response, error) {
	return nil, nil, fmt.Errorf("CryptoObjectService doesn't support Upload Now")
}

func (s *CryptoObjectService) Download(ctx context.Context, name string, filepath string, opt *cos.MultiDownloadOptions) (*cos.Response, error) {
	return nil, fmt.Errorf("CryptoObjectService doesn't support Download Now")
}

func adjustRangeStart(start int64, alignLen int64) int64 {
	return (start / alignLen) * alignLen
}

func addCryptoHeaders(header *http.Header, cd *CipherData) {
	if cd.MatDesc != "" {
		header.Add(COSClientSideEncryptionMatDesc, cd.MatDesc)
	}
	// encrypted key
	strEncryptedKey := base64.StdEncoding.EncodeToString(cd.EncryptedKey)
	header.Add(COSClientSideEncryptionKey, strEncryptedKey)

	// encrypted iv
	strEncryptedIV := base64.StdEncoding.EncodeToString(cd.EncryptedIV)
	header.Add(COSClientSideEncryptionStart, strEncryptedIV)

	header.Add(COSClientSideEncryptionWrapAlg, cd.WrapAlgorithm)
	header.Add(COSClientSideEncryptionCekAlg, cd.CEKAlgorithm)
}

func getEnvelopeFromHeader(header *http.Header) (Envelope, error) {
	var envelope Envelope

	envelope.CipherKey = header.Get(COSClientSideEncryptionKey)
	decodedKey, err := base64.StdEncoding.DecodeString(envelope.CipherKey)
	if err != nil {
		return envelope, err
	}
	envelope.CipherKey = string(decodedKey)

	envelope.IV = header.Get(COSClientSideEncryptionStart)
	decodedIV, err := base64.StdEncoding.DecodeString(envelope.IV)
	if err != nil {
		return envelope, err
	}
	envelope.IV = string(decodedIV)

	envelope.MatDesc = header.Get(COSClientSideEncryptionMatDesc)
	envelope.WrapAlg = header.Get(COSClientSideEncryptionWrapAlg)
	envelope.CEKAlg = header.Get(COSClientSideEncryptionCekAlg)
	return envelope, nil
}

func isEncrypted(header *http.Header) bool {
	encryptedKey := header.Get(COSClientSideEncryptionKey)
	if len(encryptedKey) > 0 {
		return true
	}
	return false
}
