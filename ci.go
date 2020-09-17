package cos

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type PicOperations struct {
	IsPicInfo int                  `json:"is_pic_info,omitempty"`
	Rules     []PicOperationsRules `json:"rules,omitemtpy"`
}

type PicOperationsRules struct {
	Bucket string `json:"bucket,omitempty"`
	FileId string `json:"fileid"`
	Rule   string `json:"rule"`
}

func EncodePicOperations(pic *PicOperations) string {
	bs, err := json.Marshal(pic)
	if err != nil {
		return ""
	}
	return string(bs)
}

type CloudImageReuslt struct {
	XMLName       xml.Name          `xml:"UploadResult"`
	OriginalInfo  *PicOriginalInfo  `xml:"OriginalInfo,omitempty"`
	ProcessObject *PicProcessObject `xml:"ProcessResults>Object,omitempty"`
}

type PicOriginalInfo struct {
	Key       string        `xml:"Key,omitempty"`
	Location  string        `xml:"Location,omitempty"`
	ImageInfo *PicImageInfo `xml:"ImageInfo,omitempty"`
}

type PicImageInfo struct {
	Format  string `xml:"Format,omitempty"`
	Width   int    `xml:"Width,omitempty"`
	Height  int    `xml:"Height,omitempty"`
	Size    int    `xml:"Size,omitempty"`
	Quality int    `xml:"Quality,omitempty"`
}

type PicProcessObject struct {
	Key      string `xml:"Key,omitempty"`
	Location string `xml:"Location,omitempty"`
	Format   string `xml:"Format,omitempty"`
	Width    int    `xml:"Width,omitempty"`
	Height   int    `xml:"Height,omitempty"`
	Size     int    `xml:"Size,omitempty"`
	Quality  int    `xml:"Quality,omitempty"`
}

type CloudImageOptions struct {
	PicOperations string `header:"Pic-Operations" xml:"-" url:"-"`
}

func (s *ObjectService) PostCI(ctx context.Context, name string, opt *CloudImageOptions) (*CloudImageReuslt, *Response, error) {
	var res CloudImageReuslt
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name) + "?image_process",
		method:    http.MethodPost,
		optHeader: opt,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type CloudImageRecognitionInfo struct {
	Code    int    `xml:"Code,omitempty"`
	Msg     string `xml:"Msg,omitempty"`
	HitFlag int    `xml:"HitFlag,omitempty"`
	Score   int    `xml:"Score,omitempty"`
	Label   string `xml:"Label,omitempty"`
	Count   int    `xml:"Count,omitempty"`
}

type CloudImageRecognitionResult struct {
	PornInfo      *CloudImageRecognitionInfo `xml:"PornInfo,omitempty"`
	TerroristInfo *CloudImageRecognitionInfo `xml:"TerroristInfo,omitempty"`
	PoliticsInfo  *CloudImageRecognitionInfo `xml:"PoliticsInfo,omitempty"`
	AdsInfo       *CloudImageRecognitionInfo `xml:"AdsInfo,omitempty"`
}

func GetRecognitionResult(body io.ReadCloser) *CloudImageRecognitionResult {
	var res CloudImageRecognitionResult
	err := xml.NewDecoder(body).Decode(&res)
	if err != nil && err != io.EOF {
		return nil
	}
	return &res
}
