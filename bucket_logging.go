package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketLoggingEnabled main struct of logging
type BucketLoggingEnabled struct {
	TargetBucket string `xml:"TargetBucket"`
	TargetPrefix string `xml:"TargetPrefix"`
}

// BucketPutLoggingOptions is the options of PutBucketLogging
type BucketPutLoggingOptions struct {
	XMLName        xml.Name              `xml:"BucketLoggingStatus"`
	LoggingEnabled *BucketLoggingEnabled `xml:"LoggingEnabled"`
}

// BucketGetLoggingResult is the result of GetBucketLogging
type BucketGetLoggingResult struct {
	XMLName        xml.Name              `xml:"BucketLoggingStatus"`
	LoggingEnabled *BucketLoggingEnabled `xml:"LoggingEnabled"`
}

// PutBucketLogging https://cloud.tencent.com/document/product/436/17054
func (s *BucketService) PutBucketLogging(ctx context.Context, opt *BucketPutLoggingOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?logging",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// GetBucketLogging https://cloud.tencent.com/document/product/436/17053
func (s *BucketService) GetBucketLogging(ctx context.Context) (*BucketGetLoggingResult, *Response, error) {
	var res BucketGetLoggingResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?logging",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err

}
