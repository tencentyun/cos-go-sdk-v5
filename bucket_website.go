package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type WebsiteRoutingRule struct {
  ConditionErrorCode string `xml:"Condition>HttpErrorCodeReturnedEquals,omitempty"`
	ConditionPrefix    string `xml:"Condition>KeyPrefixEquals,omitempty"`

  RedirectProtocol         string `xml:"Redirect>Protocol,omitempty"`
	RedirectReplaceKey       string `xml:"Redirect>ReplaceKeyWith,omitempty"`
	RedirectReplaceKeyPrefix string `xml:"Redirect>ReplaceKeyPrefixWith,omitempty"`
}

type BucketWebsiteConfiguration struct {
	XMLName           xml.Name     `xml:"WebsiteConfiguration"`
	Index             string       `xml:"IndexDocument>Suffix"`
	RedirectProtocol  string       `xml:"RedirectAllRequestsTo>Protocol,omitempty"`
	Error             string       `xml:"ErrorDocument>Key,omitempty"`
	Rules             []WebsiteRoutingRule     `xml:"RoutingRules>RoutingRule,omitempty"`
}

func (s *BucketService) PutWebsite(ctx context.Context, opt *BucketWebsiteConfiguration) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

func (s *BucketService) GetWebsite(ctx context.Context) (*BucketWebsiteConfiguration, *Response, error) {
	var res BucketWebsiteConfiguration
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodGet,
		result:  &res,
	}
  resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

func (s *BucketService) DelWebsite(ctx context.Context) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
