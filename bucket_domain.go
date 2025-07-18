package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type BucketDomainRule struct {
	Status            string `xml:"Status,omitempty"`
	Name              string `xml:"Name,omitempty"`
	Type              string `xml:"Type,omitempty"`
	ForcedReplacement string `xml:"ForcedReplacement,omitempty"`
}

type BucketPutDomainOptions struct {
	XMLName xml.Name           `xml:"DomainConfiguration"`
	Rules   []BucketDomainRule `xml:"DomainRule,omitempty"`
}
type BucketGetDomainResult BucketPutDomainOptions

func (s *BucketService) PutDomain(ctx context.Context, opt *BucketPutDomainOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?domain",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return resp, err
}

type BucketGetDomainOptions struct {
	XOptionHeader *http.Header `header:"-,omitempty" url:"-" xml:"-"`
}

func (s *BucketService) GetDomain(ctx context.Context, opt ...*BucketGetDomainOptions) (*BucketGetDomainResult, *Response, error) {
	var gopt *BucketGetDomainOptions
	if len(opt) > 0 {
		gopt = opt[0]
	}
	var res BucketGetDomainResult
	sendOpt := &sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/?domain",
		method:    http.MethodGet,
		result:    &res,
		optHeader: gopt,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return &res, resp, err
}

type BucketDeleteDomainOptions struct {
	XOptionHeader *http.Header `header:"-,omitempty" url:"-" xml:"-"`
}

func (s *BucketService) DeleteDomain(ctx context.Context, opt ...*BucketDeleteDomainOptions) (*Response, error) {
	var dopt *BucketDeleteDomainOptions
	if len(opt) > 0 {
		dopt = opt[0]
	}
	sendOpt := &sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/?domain",
		method:    http.MethodDelete,
		optHeader: dopt,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return resp, err
}

type BucketPutDomainCertificateOptions struct {
	XMLName         xml.Name                     `xml:"DomainCertificate"`
	CertificateInfo *BucketDomainCertificateInfo `xml:"CertificateInfo"`
	DomainList      []string                     `xml:"DomainList>DomainName"`
}

type BucketDomainCertificateInfo struct {
	CertType   string                  `xml:"CertType,omitempty"`
	CustomCert *BucketDomainCustomCert `xml:"CustomCert,omitempty"`
}
type BucketDomainCustomCert struct {
	CertId     string `xml:"CertId,omitempty"`
	Cert       string `xml:"Cert,omitempty"`
	PrivateKey string `xml:"PrivateKey,omitempty"`
}

func (s *BucketService) PutDomainCertificate(ctx context.Context, opt *BucketPutDomainCertificateOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?domaincertificate",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return resp, err
}

type BucketGetDomainCertificateResult struct {
	XMLName         xml.Name                        `xml:"DomainCertificate"`
	Status          string                          `xml:"Status,omitempty"`
	CertificateInfo *GetBucketDomainCertificateInfo `xml:"CertificateInfo"`
}

type GetBucketDomainCertificateInfo struct {
	CertId         string `xml:"CertId"`
	ValidityBegin  int64  `xml:"ValidityBegin"`
	ValidityExpire int64  `xml:"ValidityExpire"`
}

type BucketGetDomainCertificateOptions struct {
	DomainName string `url:"domainname"`
}

type BucketDeleteDomainCertificateOptions BucketGetDomainCertificateOptions

func (s *BucketService) GetDomainCertificate(ctx context.Context, opt *BucketGetDomainCertificateOptions) (*BucketGetDomainCertificateResult, *Response, error) {
	var res BucketGetDomainCertificateResult
	sendOpt := &sendOptions{
		baseURL:  s.client.BaseURL.BucketURL,
		uri:      "/?domaincertificate",
		method:   http.MethodGet,
		optQuery: opt,
		result:   &res,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return &res, resp, err
}

func (s *BucketService) DeleteDomainCertificate(ctx context.Context, opt *BucketDeleteDomainCertificateOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL:  s.client.BaseURL.BucketURL,
		uri:      "/?domaincertificate",
		method:   http.MethodDelete,
		optQuery: opt,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return resp, err
}
