package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/clbanning/mxj"
	"net/http"
	"strconv"
	"strings"
)

type BucketPutOriginOptions struct {
	XMLName xml.Name           `xml:"OriginConfiguration"`
	Rule    []BucketOriginRule `xml:"OriginRule"`
}

type BucketOriginRule struct {
	RulePriority    int                    `xml:"RulePriority,omitempty"`
	OriginType      string                 `xml:"OriginType,omitempty"`
	OriginCondition *BucketOriginCondition `xml:"OriginCondition,omitempty"`
	OriginParameter *BucketOriginParameter `xml:"OriginParameter,omitempty"`
	OriginInfo      *BucketOriginInfo      `xml:"OriginInfo,omitempty"`
	HttpStandbyCode *HTTPStandbyCode       `xml:"HTTPStandbyCode,omitempty"`
}

type BucketOriginCondition struct {
	HTTPStatusCode string `xml:"HTTPStatusCode,omitempty"`
	Prefix         string `xml:"Prefix,omitempty"`
	Suffix         string `xml:"Suffix,omitempty"`
}

type BucketOriginParameter struct {
	Protocol                       string                          `xml:"Protocol,omitempty"`
	TransparentErrorCode           *bool                           `xml:"TransparentErrorCode,omitempty"`
	FollowQueryString              *bool                           `xml:"FollowQueryString,omitempty"`
	HttpHeader                     *BucketOriginHttpHeader         `xml:"HttpHeader,omitempty"`
	FollowRedirection              *bool                           `xml:"FollowRedirection,omitempty"`
	FollowRedirectionConfiguration *FollowRedirectionConfiguration `xml:"FollowRedirectionConfiguration,omitempty"`
	HttpRedirectCode               string                          `xml:"HttpRedirectCode,omitempty"`
	CopyOriginData                 *bool                           `xml:"CopyOriginData,omitempty"`
}

type FollowRedirectionConfiguration struct {
	FollowOriginHeaders *bool `xml:"FollowOriginHeaders,omitempty"`
	FollowUrlAutoDecode *bool `xml:"FollowUrlAutoDecode,omitempty"`
}

type BucketOriginHttpHeader struct {
	FollowAllHeaders    *bool              `xml:"FollowAllHeaders,omitempty"`
	NewHttpHeaders      []OriginHttpHeader `xml:"NewHttpHeaders>Header,omitempty"`
	FollowHttpHeaders   []OriginHttpHeader `xml:"FollowHttpHeaders>Header,omitempty"`
	ForbidFollowHeaders []OriginHttpHeader `xml:"ForbidFollowHttpHeaders>Header,omitempty"`
}

type OriginHttpHeader struct {
	Key   string `xml:"Key,omitempty"`
	Value string `xml:"Value,omitempty"`
}

type BucketOriginInfo struct {
	HostInfo []*BucketOriginHostInfo `xml:"HostInfo,omitempty"`
	FileInfo *BucketOriginFileInfo   `xml:"FileInfo,omitempty"`
}

type BucketOriginHostInfo struct {
	HostName             string
	Weight               int64
	StandbyHostName_N    []string
	PrivateHost          *BucketOriginPrivateHost
	PrivateStandbyHost_N []*BucketOriginPrivateHost
}

type BucketOriginPrivateHost struct {
	Host               string                          `xml:"Host,omitempty"`
	CredentialProvider *BucketOriginCredentialProvider `xml:"CredentialProvider,omitempty"`
}
type BucketOriginCredentialProvider struct {
	AuthorizationAlgorithm string `xml:"AuthorizationAlgorithm,omitempty"`
	Region                 string `xml:"Region,omitempty"`
	SecretId               string `xml:"SecretId,omitempty"`
	SecretKey              string `xml:"SecretKey,omitempty"`
	EncryptedSecretKey     string `xml:"EncryptedSecretKey,omitempty"`
	Role                   string `xml:"Role,omitempty"`
}

type BucketOriginFileInfo struct {
	// 兼容旧版本
	PrefixDirective bool   `xml:"PrefixDirective,omitempty"`
	Prefix          string `xml:"Prefix,omitempty"`
	Suffix          string `xml:"Suffix,omitempty"`
	// 新版本
	PrefixConfiguration    *OriginPrefixConfiguration    `xml:"PrefixConfiguration,omitempty"`
	SuffixConfiguration    *OriginSuffixConfiguration    `xml:"SuffixConfiguration,omitempty"`
	FixedFileConfiguration *OriginFixedFileConfiguration `xml:"FixedFileConfiguration,omitempty"`
}

type OriginPrefixConfiguration struct {
	Prefix             string `xml:"Prefix,omitempty"`
	ReplacedWithPrefix string `xml:"ReplacedWithPrefix,omitempty"`
}

type OriginSuffixConfiguration struct {
	Suffix             string `xml:"Suffix,omitempty"`
	ReplacedWithSuffix string `xml:"ReplacedWithSuffix,omitempty"`
}

type OriginFixedFileConfiguration struct {
	FixedFilePath string `xml:"FixedFilePath,omitempty"`
}

type HTTPStandbyCode struct {
	StatusCode []string `xml:"StatusCode,omitempty"`
}

type BucketGetOriginResult BucketPutOriginOptions

func (this *BucketOriginHostInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if this == nil {
		return nil
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	if this.HostName != "" {
		err = e.EncodeElement(this.HostName, xml.StartElement{Name: xml.Name{Local: "HostName"}})
		if err != nil {
			return err
		}
	}
	if this.Weight != 0 {
		err = e.EncodeElement(this.Weight, xml.StartElement{Name: xml.Name{Local: "Weight"}})
		if err != nil {
			return err
		}
	}
	if this.PrivateHost != nil {
		err = e.EncodeElement(this.PrivateHost, xml.StartElement{Name: xml.Name{Local: "PrivateHost"}})
		if err != nil {
			return err
		}
	}
	for index, standByHostName := range this.StandbyHostName_N {
		err = e.EncodeElement(standByHostName, xml.StartElement{Name: xml.Name{Local: fmt.Sprintf("StandbyHostName_%v", index+1)}})
		if err != nil {
			return err
		}
	}
	for index, privateStandbyHost := range this.PrivateStandbyHost_N {
		err = e.EncodeElement(privateStandbyHost, xml.StartElement{Name: xml.Name{Local: fmt.Sprintf("PrivateStandbyHost_%v", index+1)}})
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (this *BucketOriginHostInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var val struct {
		XMLName xml.Name
		Inner   []byte `xml:",innerxml"`
	}
	err := d.DecodeElement(&val, &start)
	if err != nil {
		return err
	}
	str := "<HostInfo>" + string(val.Inner) + "</HostInfo>"
	myMxjMap, err := mxj.NewMapXml([]byte(str))
	if err != nil {
		return err
	}
	myMap, ok := myMxjMap["HostInfo"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("XML HostInfo Parse failed")
	}

	privateStandbyHostMap := make(map[string]*BucketOriginPrivateHost)
	var total int
	for key, value := range myMap {
		if key == "HostName" {
			this.HostName = value.(string)
		}
		if key == "Weight" {
			v := value.(string)
			this.Weight, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
		}
		if strings.HasPrefix(key, "StandbyHostName_") {
			total++
		}
		if key == "PrivateHost" {
			if _, ok := value.(map[string]interface{}); ok {
				this.PrivateHost = &BucketOriginPrivateHost{}
				err = mxj.Map(value.(map[string]interface{})).Struct(this.PrivateHost)
				if err != nil {
					return err
				}
			}
		}
		if strings.HasPrefix(key, "PrivateStandbyHost_") {
			if _, ok := value.(map[string]interface{}); ok {
				var privateStandbyHost_N BucketOriginPrivateHost
				err = mxj.Map(value.(map[string]interface{})).Struct(&privateStandbyHost_N)
				if err != nil {
					return err
				}
				privateStandbyHostMap[key] = &privateStandbyHost_N
			}
		}
	}
	// 按顺序执行
	for i := 1; i <= total; i++ {
		key := fmt.Sprintf("StandbyHostName_%v", i)
		this.StandbyHostName_N = append(this.StandbyHostName_N, myMap[key].(string))
	}
	for i := 1; i <= len(privateStandbyHostMap); i++ {
		key := fmt.Sprintf("PrivateStandbyHost_%v", i)
		this.PrivateStandbyHost_N = append(this.PrivateStandbyHost_N, privateStandbyHostMap[key])
	}

	return nil
}

func (s *BucketService) PutOrigin(ctx context.Context, opt *BucketPutOriginOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?origin",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetOrigin(ctx context.Context) (*BucketGetOriginResult, *Response, error) {
	var res BucketGetOriginResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?origin",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return &res, resp, err
}

func (s *BucketService) DeleteOrigin(ctx context.Context) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?origin",
		method:  http.MethodDelete,
	}
	resp, err := s.client.doRetry(ctx, sendOpt)
	return resp, err
}
