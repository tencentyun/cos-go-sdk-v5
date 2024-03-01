package cos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Dataset struct {
	BindCount     int    `json:"BindCount"`
	CreateTime    string `json:"CreateTime"`
	DatasetName   string `json:"DatasetName"`
	Description   string `json:"Description"`
	FileCount     int    `json:"FileCount"`
	TemplateID    string `json:"TemplateId"`
	TotalFileSize int    `json:"TotalFileSize"`
	UpdateTime    string `json:"UpdateTime"`
}

type OptHeaders struct {
	XOptionHeader *http.Header `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type CreateDataSetOptions struct {
	DatasetName string      `json:"DatasetName" url:"-"`
	Description string      `json:"Description" url:"-"`
	TemplateId  string      `json:"TemplateId" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type CreateDataSetResult struct {
	Response struct {
		Dataset   Dataset `json:"Dataset"`
		RequestID string  `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) baseSend(ctx context.Context, opt interface{}, optionHeader *OptHeaders, uri string, method string) (*bytes.Buffer, *Response, error) {
	var buf bytes.Buffer
	var f *strings.Reader
	var sendOpt *sendOptions
	if optionHeader == nil {
		optionHeader = &OptHeaders{
			XOptionHeader: &http.Header{},
		}
	}
	optionHeader.XOptionHeader.Add("Content-Type", "application/json")
	optionHeader.XOptionHeader.Add("Accept", "application/json")
	if method == http.MethodGet {
		sendOpt = &sendOptions{
			baseURL:   s.client.BaseURL.CIURL,
			uri:       uri,
			method:    method,
			optHeader: optionHeader,
			optQuery:  opt,
			result:    &buf,
		}
	} else {
		if opt != nil {
			bs, err := json.Marshal(opt)
			if err != nil {
				return nil, nil, err
			}
			f = strings.NewReader(string(bs))
		}
		sendOpt = &sendOptions{
			baseURL:   s.client.BaseURL.CIURL,
			uri:       uri,
			method:    method,
			body:      f,
			optHeader: optionHeader,
			result:    &buf,
		}
	}
	resp, err := s.client.send(ctx, sendOpt)
	return &buf, resp, err
}

func (s *CIService) CreateDataSet(ctx context.Context, opt *CreateDataSetOptions) (*CreateDataSetResult, *Response, error) {
	var res CreateDataSetResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/dataset/create", http.MethodPost)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DescribeDatasetsOptions struct {
	MaxResults int64       `url:"maxresults,omitempty" url:"-"`
	Prefix     string      `url:"prefix,omitempty" url:"-"`
	Nexttoken  string      `url:"nexttoken,omitempty" url:"-"`
	OptHeaders *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}
type DescribeDatasetsResult struct {
	Response struct {
		Datasets  []Dataset `json:"Datasets"`
		NextToken string    `json:"NextToken"`
		RequestID string    `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DescribeDatasets(ctx context.Context, opt *DescribeDatasetsOptions) (*DescribeDatasetsResult, *Response, error) {
	var res DescribeDatasetsResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/datasets", http.MethodGet)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type UpdateDatasetOptions struct {
	DatasetName string      `json:"DatasetName" url:"-"`
	Description string      `json:"Description" url:"-"`
	TemplateId  string      `json:"TemplateId" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type UpdateDatasetResult struct {
	Response struct {
		Dataset   Dataset `json:"Dataset"`
		RequestID string  `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) UpdateDataset(ctx context.Context, opt *UpdateDatasetOptions) (*UpdateDatasetResult, *Response, error) {
	var res UpdateDatasetResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/dataset/update", http.MethodPost)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DeleteDatasetOptions struct {
	DatasetName string      `json:"DatasetName" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type DeleteDatasetResult struct {
	Response struct {
		Dataset   Dataset `json:"Dataset"`
		RequestID string  `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DeleteDataset(ctx context.Context, opt *DeleteDatasetOptions) (*DeleteDatasetResult, *Response, error) {
	var res DeleteDatasetResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/dataset", http.MethodDelete)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DescribeDatasetOptions struct {
	DatasetName string      `url:"datasetname,omitempty" url:"-"`
	Statistics  bool        `url:"statistics,omitempty" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type DescribeDatasetResult struct {
	Response struct {
		Dataset   Dataset `json:"Dataset"`
		RequestID string  `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DescribeDataset(ctx context.Context, opt *DescribeDatasetOptions) (*DescribeDatasetResult, *Response, error) {
	var res DescribeDatasetResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/dataset", http.MethodGet)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type File struct {
	URI          string             `json:"URI,omitempty"`
	CustomID     string             `json:"CustomId,omitempty"`
	CustomLabels *map[string]string `json:"CustomLabels,omitempty"`
	MediaType    string             `json:"MediaType,omitempty"`
	ContentType  string             `json:"contenttype,omitempty"`
}

type CreateFileMetaIndexOptions struct {
	DatasetName string      `json:"DatasetName,omitempty" url:"-"`
	File        *File       `json:"File,omitempty" url:"-"`
	Callback    string      `json:"Callback,omitempty" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type CreateFileMetaIndexResult struct {
	Response struct {
		EventID   string `json:"EventId"`
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) CreateFileMetaIndex(ctx context.Context, opt *CreateFileMetaIndexOptions) (*CreateFileMetaIndexResult, *Response, error) {
	var res CreateFileMetaIndexResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/filemeta/create", http.MethodPost)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type UpdateFileMetaIndexOptions struct {
	DatasetName string      `json:"DatasetName,omitempty" url:"-"`
	File        *File       `json:"File,omitempty" url:"-"`
	Callback    string      `json:"Callback,omitempty" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type UpdateFileMetaIndexResult struct {
	Response struct {
		EventID   string `json:"EventId"`
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) UpdateFileMetaIndex(ctx context.Context, opt *UpdateFileMetaIndexOptions) (*UpdateFileMetaIndexResult, *Response, error) {
	var res UpdateFileMetaIndexResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/filemeta/update", http.MethodPost)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DescribeFileMetaIndexOptions struct {
	DatasetName string      `json:"-" url:"datasetname,omitempty"`
	Uri         string      `json:"-" url:"uri,omitempty"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type FileInfo struct {
	DatasetName          string            `json:"DatasetName,omitempty"`
	OwnerID              string            `json:"OwnerID,omitempty"`
	ObjectId             string            `json:"ObjectId,omitempty"`
	CreateTime           string            `json:"CreateTime,omitempty"`
	UpdateTime           string            `json:"UpdateTime,omitempty"`
	URI                  string            `json:"URI,omitempty"`
	Filename             string            `json:"Filename,omitempty"`
	MediaType            string            `json:"MediaType,omitempty"`
	ContentType          string            `json:"ContentType,omitempty"`
	COSStorageClass      string            `json:"COSStorageClass,omitempty"`
	Coscrc64             string            `json:"COSCRC64,omitempty"`
	Size                 int               `json:"Size,omitempty"`
	CacheControl         string            `json:"CacheControl,omitempty"`
	ContentDisposition   string            `json:"ContentDisposition,omitempty"`
	ContentEncoding      string            `json:"ContentEncoding,omitempty"`
	ContentLanguage      string            `json:"ContentLanguage,omitempty"`
	ServerSideEncryption string            `json:"ServerSideEncryption,omitempty"`
	ETag                 string            `json:"ETag,omitempty"`
	FileModifiedTime     string            `json:"FileModifiedTime,omitempty"`
	CustomID             string            `json:"CustomId,omitempty"`
	CustomLabels         map[string]string `json:"CustomLabels,omitempty"`
	COSUserMeta          map[string]string `json:"COSUserMeta,omitempty"`
	ObjectACL            string            `json:"ObjectACL",omitempty`
	COSTagging           map[string]string `json:"COSTagging,omitempty"`
}

type DescribeFileMetaIndexResult struct {
	Response struct {
		Files     []FileInfo `json:"Files"`
		RequestID string     `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DescribeFileMetaIndex(ctx context.Context, opt *DescribeFileMetaIndexOptions) (*DescribeFileMetaIndexResult, *Response, error) {
	var res DescribeFileMetaIndexResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/filemeta", http.MethodGet)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DeleteFileMetaIndexOptions struct {
	DatasetName string      `url:"-" json:"DatasetName,omitempty"`
	Uri         string      `url:"-" json:"URI,omitempty"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type DeleteFileMetaIndexResult struct {
	Response struct {
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DeleteFileMetaIndex(ctx context.Context, opt *DeleteFileMetaIndexOptions) (*DeleteFileMetaIndexResult, *Response, error) {
	var res DeleteFileMetaIndexResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/filemeta", http.MethodDelete)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type Query struct {
	Operation  string   `json:"Operation,omitempty"`
	Field      string   `json:"Field,omitempty"`
	Value      string   `json:"Value,omitempty"`
	SubQueries []*Query `json:"SubQueries,omitempty"`
}

type Aggregation struct {
	Field     string `json:"Field,omitempty"`
	Operation string `json:"Operation,omitempty"`
}
type DatasetSimpleQueryOptions struct {
	DatasetName  string         `json:"DatasetName,omitempty" url:"-"`
	Query        *Query         `json:"Query,omitempty" url:"-"`
	Sort         string         `json:"Sort,omitempty" url:"-"`
	Order        string         `json:"Order,omitempty" url:"-"`
	MaxResults   string         `json:"MaxResults,omitempty" url:"-"`
	Aggregations []*Aggregation `json:"Aggregations,omitempty" url:"-"`
	NextToken    string         `json:"NextToken,omitempty" url:"-"`
	WithFields   []string       `json:"WithFields,omitempty"  url:"-"`
	OptHeaders   *OptHeaders    `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type Groups struct {
	Count int    `json:"Count"`
	Value string `json:"Value"`
}
type Aggregations struct {
	Field     string   `json:"Field"`
	Groups    []Groups `json:"Groups"`
	Operation string   `json:"Operation"`
	Value     float32  `json:"Value"`
}

type DatasetSimpleQueryResult struct {
	Response struct {
		Aggregations []Aggregations `json:"Aggregations"`
		Files        []FileInfo     `json:"Files"`
		NextToken    string         `json:"NextToken"`
		RequestID    string         `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DatasetSimpleQuery(ctx context.Context, opt *DatasetSimpleQueryOptions) (*DatasetSimpleQueryResult, *Response, error) {
	var res DatasetSimpleQueryResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/datasetquery/simple", http.MethodPost)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type CreateDatasetBindingOptions struct {
	DatasetName string      `json:"DatasetName,omitempty" url:"-"`
	URI         string      `json:"URI,omitempty" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type Binding struct {
	CreateTime  string `json:"CreateTime,omitempty" url:"-"`
	DatasetName string `json:"DatasetName,omitempty" url:"-"`
	Detail      string `json:"Detail,omitempty" url:"-"`
	State       string `json:"State,omitempty" url:"-"`
	URI         string `json:"URI,omitempty" url:"-"`
	UpdateTime  string `json:"UpdateTime,omitempty" url:"-"`
}
type CreateDatasetBindingResult struct {
	Response struct {
		Binding   Binding `json:"Binding,omitempty"`
		RequestID string  `json:"RequestId,omitempty"`
	} `json:"Response,omitempty"`
}

func (s *CIService) CreateDatasetBinding(ctx context.Context, opt *CreateDatasetBindingOptions) (*CreateDatasetBindingResult, *Response, error) {
	var res CreateDatasetBindingResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/datasetbinding/create", http.MethodPost)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DescribeDatasetBindingOptions struct {
	DatasetName string      `json:"-" url:"datasetname,omitempty"`
	URI         string      `json:"-" url:"uri,omitempty"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type DescribeDatasetBindingResult struct {
	Response struct {
		Binding   Binding `json:"Binding,omitempty"`
		RequestID string  `json:"RequestId,omitempty"`
	} `json:"Response,omitempty"`
}

func (s *CIService) DescribeDatasetBinding(ctx context.Context, opt *DescribeDatasetBindingOptions) (*DescribeDatasetBindingResult, *Response, error) {
	var res DescribeDatasetBindingResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/datasetbinding", http.MethodGet)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DescribeDatasetBindingsOptions struct {
	DatasetName string      `json:"-" url:"datasetname,omitempty"`
	MaxResults  int         `json:"-" url:"maxresults,omitempty"`
	NextToken   string      `json:"-" url:"nexttoken,omitempty"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type DescribeDatasetBindingsResult struct {
	Response struct {
		Bindings  []*Binding `json:"Bindings,omitempty"`
		NextToken string     `json:"NextToken,omitempty"`
		RequestID string     `json:"RequestId,omitempty"`
	} `json:"Response,omitempty"`
}

func (s *CIService) DescribeDatasetBindings(ctx context.Context, opt *DescribeDatasetBindingsOptions) (*DescribeDatasetBindingsResult, *Response, error) {
	var res DescribeDatasetBindingsResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/datasetbindings", http.MethodGet)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

type DeleteDatasetBindingOptions struct {
	DatasetName string      `json:"DatasetName,omitempty" url:"-"`
	URI         string      `json:"URI,omitempty" url:"-"`
	OptHeaders  *OptHeaders `header:"-,omitempty" url:"-" json:"-" xml:"-"`
}

type DeleteDatasetBindingResult struct {
	Response struct {
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

func (s *CIService) DeleteDatasetBinding(ctx context.Context, opt *DeleteDatasetBindingOptions) (*DeleteDatasetBindingResult, *Response, error) {
	var res DeleteDatasetBindingResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, opt.OptHeaders, "/datasetbinding", http.MethodDelete)
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}
