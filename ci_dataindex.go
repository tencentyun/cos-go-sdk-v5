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
	Coscrc64         string            `json:"COSCRC64"`
	COSStorageClass  string            `json:"COSStorageClass"`
	CacheControl     string            `json:"CacheControl"`
	ContentType      string            `json:"ContentType"`
	CreateTime       string            `json:"CreateTime"`
	CustomID         string            `json:"CustomId"`
	CustomLabels     map[string]string `json:"CustomLabels"`
	DatasetName      string            `json:"DatasetName"`
	ETag             string            `json:"ETag"`
	FileModifiedTime string            `json:"FileModifiedTime"`
	Filename         string            `json:"Filename"`
	MediaType        string            `json:"MediaType"`
	ObjectACL        string            `json:"ObjectACL"`
	Size             int            `json:"Size"`
	URI              string            `json:"URI"`
	UpdateTime       string            `json:"UpdateTime"`
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
		RequestID string     `json:"RequestId"`
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

