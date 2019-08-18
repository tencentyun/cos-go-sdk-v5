package cos

import (
  "context"
  "net/http"
	"encoding/xml"
)

type BucketDomainConfiguration struct {
  XMLName   xml.Name  `xml:"DomainConfiguration"`
  Status    string    `xml:"DomainRule>Status.omitempty"`
  Name      string    `xml:"DomainRule>Name,omitempty"`
  Type      string    `xml:"DomainRule>Type,omitempty"`
  ForcedReplacement string `xml:"DomainRule>ForcedReplacement,omitempty"`
}

func (s *BucketService) PutDomain(ctx context.Context, opt *BucketDomainConfiguration) (*Response, error) {
  sendOpt := sendOptions {
    baseURL : s.client.BaseURL.BucketURL,
    uri : "/?domain",
    method : http.MethodPut,
    body : opt,
  }
  resp, err := s.client.send(ctx, &sendOpt)
  return resp, err
}

func (s *BucketService) GetDomain(ctx context.Context) (*BucketDomainConfiguration, *Response, error) {
  var res BucketDomainConfiguration
  sendOpt := sendOptions {
    baseURL : s.client.BaseURL.BucketURL,
    uri : "/?domain",
    method : http.MethodGet,
    result : &res,
  }
  resp, err := s.client.send(ctx, &sendOpt)
  return &res, resp, err
}
