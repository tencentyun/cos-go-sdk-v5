package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type WebsiteIndex struct {
	Suffix string 	`XML:",omitempty"`
}

type WebsiteRedirectProtocol struct {
	Protocol string	`XML:",omitempty"`
}

type WebsiteError struct {
	Key string	`xml:",omitempty"`
}

type WebsiteRuleCondition struct {
	ErrorCode string `xml:"HttpErrorCodeReturnedEquals,omitempty"`
	Prefix    string `xml:"KeyPrefixEquals,omitempty"`
}
type WebsiteRedirect struct {
	Protocol         string `xml:"Protocol,omitempty"`
	ReplaceKeyPrefix string `xml:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKey       string `xml:"ReplaceKeyWith,omitempty"`
}
type WebsiteRoutingRule struct {
	Condition *WebsiteRuleCondition	`xml:",omitempty"`
	Redirect  *WebsiteRedirect	`xml:",omitempty"`
}
type WebsiteRoutingRules struct {
	Rule []WebsiteRoutingRule `xml:"RoutingRule"`
}

type BucketWebsiteConfiguration struct {
	XMLName  xml.Name                 `xml:"WebsiteConfiguration"`
	Index    *WebsiteIndex            `xml:"IndexDocument"`
	Redirect *WebsiteRedirectProtocol `xml:"RedirectAllRequestsTo,omitempty"`
	Error    *WebsiteError            `xml:"ErrorDocument,omitempty"`
	Rules    *WebsiteRoutingRules     `xml:"RoutingRules,omitempty"`
}

type PPQ struct {
	id string `url:"id"`
	op string `url:"op"`
}

func (s *BucketService) PutWebsite(ctx context.Context, opt *BucketWebsiteConfiguration) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodPut,
		body:    opt,
		optQuery: PPQ{"aaa", "bbb"},
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
