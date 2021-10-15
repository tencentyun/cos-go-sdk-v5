package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

type JobInput struct {
	Object string `xml:"Object,omitempty"`
}

type JobOutput struct {
	Region       string `xml:"Region,omitempty"`
	Bucket       string `xml:"Bucket,omitempty"`
	Object       string `xml:"Object,omitempty"`
	SpriteObject string `xml:"SpriteObject,omitempty"`
}

type Container struct {
	Format string `xml:"Format"`
}

type Video struct {
	Codec         string `xml:"Codec"`
	Width         string `xml:"Width"`
	Height        string `xml:"Height"`
	Fps           string `xml:"Fps"`
	Remove        string `xml:"Remove,omitempty"`
	Profile       string `xml:"Profile"`
	Bitrate       string `xml:"Bitrate"`
	Crf           string `xml:"Crf"`
	Gop           string `xml:"Gop"`
	Preset        string `xml:"Preset"`
	Bufsize       string `xml:"Bufsize"`
	Maxrate       string `xml:"Maxrate"`
	HlsTsTime     string `xml:"HlsTsTime"`
	Pixfmt        string `xml:"Pixfmt"`
	LongShortMode string `xml:"LongShortMode"`
}

type TimeInterval struct {
	Start    string `xml:"Start"`
	Duration string `xml:"Duration"`
}

type Audio struct {
	Codec      string `xml:"Codec"`
	Samplerate string `xml:"Samplerate"`
	Bitrate    string `xml:"Bitrate"`
	Channels   string `xml:"Channels"`
	Remove     string `xml:"Remove,omitempty"`
}

type TransConfig struct {
	AdjDarMethod          string `xml:"AdjDarMethod"`
	IsCheckReso           string `xml:"IsCheckReso"`
	ResoAdjMethod         string `xml:"ResoAdjMethod"`
	IsCheckVideoBitrate   string `xml:"IsCheckVideoBitrate"`
	VideoBitrateAdjMethod string `xml:"VideoBitrateAdjMethod"`
	IsCheckAudioBitrate   string `xml:"IsCheckAudioBitrate"`
	AudioBitrateAdjMethod string `xml:"AudioBitrateAdjMethod"`
}

type Transcode struct {
	Container    *Container    `xml:"Container,omitempty"`
	Video        *Video        `xml:"Video,omitempty"`
	TimeInterval *TimeInterval `xml:"TimeInterval,omitempty"`
	Audio        *Audio        `xml:"Audio,omitempty"`
	TransConfig  *TransConfig  `xml:"TransConfig,omitempty"`
}

type Image struct {
	Url          string `xml:"Url,omitempty"`
	Mode         string `xml:"Mode,omitempty"`
	Width        string `xml:"Width,omitempty"`
	Height       string `xml:"Height,omitempty"`
	Transparency string `xml:"Transparency,omitempty"`
	Background   string `xml:"Background,omitempty"`
}

type Text struct {
	FontSize     string `xml:"FontSize,omitempty"`
	FontType     string `xml:"FontType,omitempty"`
	FontColor    string `xml:"FontColor,omitempty"`
	Transparency string `xml:"Transparency,omitempty"`
	Text         string `xml:"Text,omitempty"`
}

type Watermark struct {
	Type      string `xml:"Type,omitempty"`
	Pos       string `xml:"Pos,omitempty"`
	LocMode   string `xml:"LocMode,omitempty"`
	Dx        string `xml:"Dx,omitempty"`
	Dy        string `xml:"Dy,omitempty"`
	StartTime string `xml:"StartTime,omitempty"`
	EndTime   string `xml:"EndTime,omitempty"`
	Image     *Image `xml:"Image,omitempty"`
	Text      *Text  `xml:"Text,omitempty"`
}
type ConcatFragment struct {
	Url       string `xml:"Url,omitempty"`
	StartTime string `xml:"StartTime,omitempty"`
	EndTime   string `xml:"EndTime,omitempty"`
}
type ConcatTemplate struct {
	ConcatFragment []ConcatFragment `xml:"ConcatFragment,omitempty"`
	Audio          *Audio           `xml:"Audio,omitempty"`
	Video          *Video           `xml:"Video,omitempty"`
	Container      *Container       `xml:"Container,omitempty"`
	Index          string           `xml:"Index,omitempty"`
}

type SpriteSnapshotConfig struct {
	CellHeight string `xml:"CellHeight,omitempty"`
	CellWidth  string `xml:"CellWidth,omitempty"`
	Color      string `xml:"Color,omitempty"`
	Columns    string `xml:"Columns,omitempty"`
	Lines      string `xml:"Lines,omitempty"`
	Margin     string `xml:"Margin,omitempty"`
	Padding    string `xml:"Padding,omitempty"`
}
type Snapshot struct {
	Mode                 string                `xml:"Mode,omitempty"`
	Start                string                `xml:"Start,omitempty"`
	TimeInterval         string                `xml:"TimeInterval,omitempty"`
	Count                string                `xml:"Count,omitempty"`
	Width                string                `xml:"Width,omitempty"`
	Height               string                `xml:"Height,omitempty"`
	SnapshotOutMode      string                `xml:"SnapshotOutMode,omitempty"`
	SpriteSnapshotConfig *SpriteSnapshotConfig `xml:"SpriteSnapshotConfig,omitempty"`
}

// 有意和转码区分，两种任务关注的参数不一样避免干扰
type AnimationVideo struct {
	Codec                      string `xml:"Codec"`
	Width                      string `xml:"Width"`
	Height                     string `xml:"Height"`
	Fps                        string `xml:"Fps"`
	AnimateOnlyKeepKeyFrame    string `xml:"AnimateOnlyKeepKeyFrame"`
	AnimateTimeIntervalOfFrame string `xml:"AnimateTimeIntervalOfFrame"`
	AnimateFramesPerSecond     string `xml:"AnimateFramesPerSecond"`
	Quality                    string `xml:"Quality"`
}

type Animation struct {
	Container    *Container      `xml:"Container,omitempty"`
	Video        *AnimationVideo `xml:"Video,omitempty"`
	TimeInterval *TimeInterval   `xml:"TimeInterval,omitempty"`
}

type Segment struct {
	Format   string `xml:"Format,omitempty"`
	Duration string `xml:"Duration,omitempty"`
}

type MediaProcessJobOperation struct {
	Tag                 string          `xml:"Tag,omitempty"`
	Output              *JobOutput      `xml:"Output,omitempty"`
	Transcode           *Transcode      `xml:"Transcode,omitempty"`
	Watermark           *Watermark      `xml:"Watermark,omitempty"`
	TemplateId          string          `xml:"TemplateId,omitempty"`
	WatermarkTemplateId []string        `xml:"WatermarkTemplateId,omitempty"`
	ConcatTemplate      *ConcatTemplate `xml:"ConcatTemplate,omitempty"`
	Snapshot            *Snapshot       `xml:"Snapshot,omitempty"`
	Animation           *Animation      `xml:"Animation,omitempty"`
	Segment             *Segment        `xml:"Segment,omitempty"`
}

type CreateMediaJobsOptions struct {
	XMLName   xml.Name                  `xml:"Request"`
	Tag       string                    `xml:"Tag,omitempty"`
	Input     *JobInput                 `xml:"Input,omitempty"`
	Operation *MediaProcessJobOperation `xml:"Operation,omitempty"`
	QueueId   string                    `xml:"QueueId,omitempty"`
	CallBack  string                    `xml:"CallBack,omitempty"`
}

type MediaProcessJobDetail struct {
	Code         string                    `xml:"Code,omitempty"`
	Message      string                    `xml:"Message,omitempty"`
	JobId        string                    `xml:"JobId,omitempty"`
	Tag          string                    `xml:"Tag,omitempty"`
	State        string                    `xml:"State,omitempty"`
	CreationTime string                    `xml:"CreationTime,omitempty"`
	QueueId      string                    `xml:"QueueId,omitempty"`
	Input        *JobInput                 `xml:"Input,omitempty"`
	Operation    *MediaProcessJobOperation `xml:"Operation,omitempty"`
}

type CreateMediaJobsResult struct {
	XMLName    xml.Name               `xml:"Response"`
	JobsDetail *MediaProcessJobDetail `xml:"JobsDetail,omitempty"`
}

type CreateMultiMediaJobsOptions struct {
	XMLName   xml.Name                   `xml:"Request"`
	Tag       string                     `xml:"Tag,omitempty"`
	Input     *JobInput                  `xml:"Input,omitempty"`
	Operation []MediaProcessJobOperation `xml:"Operation,omitempty"`
	QueueId   string                     `xml:"QueueId,omitempty"`
	CallBack  string                     `xml:"CallBack,omitempty"`
}

type CreateMultiMediaJobsResult struct {
	XMLName    xml.Name                `xml:"Response"`
	JobsDetail []MediaProcessJobDetail `xml:"JobsDetail,omitempty"`
}

type MediaProcessJobsNotifyBody struct {
	XMLName    xml.Name `xml:"Response"`
	EventName  string   `xml:"EventName"`
	JobsDetail struct {
		Code         string `xml:"Code"`
		CreationTime string `xml:"CreationTime"`
		EndTime      string `xml:"EndTime"`
		Input        struct {
			BucketId string `xml:"BucketId"`
			Object   string `xml:"Object"`
			Region   string `xml:"Region"`
		} `xml:"Input"`
		JobId     string `xml:"JobId"`
		Message   string `xml:"Message"`
		Operation struct {
			MediaInfo struct {
				Format struct {
					Text           string `xml:",chardata"`
					Bitrate        string `xml:"Bitrate"`
					Duration       string `xml:"Duration"`
					FormatLongName string `xml:"FormatLongName"`
					FormatName     string `xml:"FormatName"`
					NumProgram     string `xml:"NumProgram"`
					NumStream      string `xml:"NumStream"`
					Size           string `xml:"Size"`
					StartTime      string `xml:"StartTime"`
				} `xml:"Format"`
				Stream struct {
					Audio []struct {
						Bitrate        string `xml:"Bitrate"`
						Channel        string `xml:"Channel"`
						ChannelLayout  string `xml:"ChannelLayout"`
						CodecLongName  string `xml:"CodecLongName"`
						CodecName      string `xml:"CodecName"`
						CodecTag       string `xml:"CodecTag"`
						CodecTagString string `xml:"CodecTagString"`
						CodecTimeBase  string `xml:"CodecTimeBase"`
						Duration       string `xml:"Duration"`
						Index          string `xml:"Index"`
						Language       string `xml:"Language"`
						SampleFmt      string `xml:"SampleFmt"`
						SampleRate     string `xml:"SampleRate"`
						StartTime      string `xml:"StartTime"`
						Timebase       string `xml:"Timebase"`
					} `xml:"Audio"`
					Subtitle string `xml:"Subtitle"`
					Video    []struct {
						AvgFps         string `xml:"AvgFps"`
						Bitrate        string `xml:"Bitrate"`
						CodecLongName  string `xml:"CodecLongName"`
						CodecName      string `xml:"CodecName"`
						CodecTag       string `xml:"CodecTag"`
						CodecTagString string `xml:"CodecTagString"`
						CodecTimeBase  string `xml:"CodecTimeBase"`
						Dar            string `xml:"Dar"`
						Duration       string `xml:"Duration"`
						Fps            string `xml:"Fps"`
						HasBFrame      string `xml:"HasBFrame"`
						Height         string `xml:"Height"`
						Index          string `xml:"Index"`
						Language       string `xml:"Language"`
						Level          string `xml:"Level"`
						NumFrames      string `xml:"NumFrames"`
						PixFormat      string `xml:"PixFormat"`
						Profile        string `xml:"Profile"`
						RefFrames      string `xml:"RefFrames"`
						Rotation       string `xml:"Rotation"`
						Sar            string `xml:"Sar"`
						StartTime      string `xml:"StartTime"`
						Timebase       string `xml:"Timebase"`
						Width          string `xml:"Width"`
					} `xml:"Video"`
				} `xml:"Stream"`
			} `xml:"MediaInfo"`
			MediaResult struct {
				OutputFile struct {
					Bucket       string   `xml:"Bucket"`
					ObjectName   []string `xml:"ObjectName"`
					ObjectPrefix string   `xml:"ObjectPrefix"`
					Region       string   `xml:"Region"`
				} `xml:"OutputFile"`
			} `xml:"MediaResult"`
			Output struct {
				Bucket string `xml:"Bucket"`
				Object string `xml:"Object"`
				Region string `xml:"Region"`
			} `xml:"Output"`
			TemplateId   string `xml:"TemplateId"`
			TemplateName string `xml:"TemplateName"`
		} `xml:"Operation"`
		QueueId   string `xml:"QueueId"`
		StartTime string `xml:"StartTime"`
		State     string `xml:"State"`
		Tag       string `xml:"Tag"`
	} `xml:"JobsDetail"`
}

type WorkflowExecutionNotifyBody struct {
	XMLName           xml.Name `xml:"Response"`
	EventName         string   `xml:"EventName"`
	WorkflowExecution struct {
		RunId      string `xml:"RunId"`
		BucketId   string `xml:"BucketId"`
		Object     string `xml:"Object"`
		CosHeaders []struct {
			Key   string `xml:"Key"`
			Value string `xml:"Value"`
		} `xml:"CosHeaders"`
		WorkflowId   string `xml:"WorkflowId"`
		WorkflowName string `xml:"WorkflowName"`
		CreateTime   string `xml:"CreateTime"`
		State        string `xml:"State"`
		Tasks        []struct {
			Type                  string `xml:"Type"`
			CreateTime            string `xml:"CreateTime"`
			EndTime               string `xml:"EndTime"`
			State                 string `xml:"State"`
			JobId                 string `xml:"JobId"`
			Name                  string `xml:"Name"`
			TemplateId            string `xml:"TemplateId"`
			TemplateName          string `xml:"TemplateName"`
			TranscodeTemplateId   string `xml:"TranscodeTemplateId,omitempty"`
			TranscodeTemplateName string `xml:"TranscodeTemplateName,omitempty"`
			HdrMode               string `xml:"HdrMode,omitempty"`
		} `xml:"Tasks"`
	} `xml:"WorkflowExecution"`
}

func (s *CIService) CreateMultiMediaJobs(ctx context.Context, opt *CreateMultiMediaJobsOptions) (*CreateMultiMediaJobsResult, *Response, error) {
	var res CreateMultiMediaJobsResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.CIURL,
		uri:     "/jobs",
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

func (s *CIService) CreateMediaJobs(ctx context.Context, opt *CreateMediaJobsOptions) (*CreateMediaJobsResult, *Response, error) {
	var res CreateMediaJobsResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.CIURL,
		uri:     "/jobs",
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type DescribeMediaProcessJobResult struct {
	XMLName        xml.Name               `xml:"Response"`
	JobsDetail     *MediaProcessJobDetail `xml:"JobsDetail,omitempty"`
	NonExistJobIds string                 `xml:"NonExistJobIds,omitempty"`
}

func (s *CIService) DescribeMediaJob(ctx context.Context, jobid string) (*DescribeMediaProcessJobResult, *Response, error) {
	var res DescribeMediaProcessJobResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.CIURL,
		uri:     "/jobs/" + jobid,
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type DescribeMediaJobsOptions struct {
	QueueId           string `url:"queueId,omitempty"`
	Tag               string `url:"tag,omitempty"`
	OrderByTime       string `url:"orderByTime,omitempty"`
	NextToken         string `url:"nextToken,omitempty"`
	Size              int    `url:"size,omitempty"`
	States            string `url:"states,omitempty"`
	StartCreationTime string `url:"startCreationTime,omitempty"`
	EndCreationTime   string `url:"endCreationTime,omitempty"`
}

type DescribeMediaJobsResult struct {
	XMLName    xml.Name                `xml:"Response"`
	JobsDetail []MediaProcessJobDetail `xml:"JobsDetail,omitempty"`
	NextToken  string                  `xml:"NextToken,omitempty"`
}

func (s *CIService) DescribeMediaJobs(ctx context.Context, opt *DescribeMediaJobsOptions) (*DescribeMediaJobsResult, *Response, error) {
	var res DescribeMediaJobsResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.CIURL,
		uri:      "/jobs",
		optQuery: opt,
		method:   http.MethodGet,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type DescribeMediaProcessQueuesOptions struct {
	QueueIds   string `url:"queueIds,omitempty"`
	State      string `url:"state,omitempty"`
	PageNumber int    `url:"pageNumber,omitempty"`
	PageSize   int    `url:"pageSize,omitempty"`
}

type DescribeMediaProcessQueuesResult struct {
	XMLName      xml.Name            `xml:"Response"`
	RequestId    string              `xml:"RequestId,omitempty"`
	TotalCount   int                 `xml:"TotalCount,omitempty"`
	PageNumber   int                 `xml:"PageNumber,omitempty"`
	PageSize     int                 `xml:"PageSize,omitempty"`
	QueueList    []MediaProcessQueue `xml:"QueueList,omitempty"`
	NonExistPIDs []string            `xml:"NonExistPIDs,omitempty"`
}

type MediaProcessQueue struct {
	QueueId       string                         `xml:"QueueId,omitempty"`
	Name          string                         `xml:"Name,omitempty"`
	State         string                         `xml:"State,omitempty"`
	MaxSize       int                            `xml:"MaxSize,omitempty"`
	MaxConcurrent int                            `xml:"MaxConcurrent,omitempty"`
	UpdateTime    string                         `xml:"UpdateTime,omitempty"`
	CreateTime    string                         `xml:"CreateTime,omitempty"`
	NotifyConfig  *MediaProcessQueueNotifyConfig `xml:"NotifyConfig,omitempty"`
}

type MediaProcessQueueNotifyConfig struct {
	Url   string `xml:"Url,omitempty"`
	State string `xml:"State,omitempty"`
	Type  string `xml:"Type,omitempty"`
	Event string `xml:"Event,omitempty"`
}

func (s *CIService) DescribeMediaProcessQueues(ctx context.Context, opt *DescribeMediaProcessQueuesOptions) (*DescribeMediaProcessQueuesResult, *Response, error) {
	var res DescribeMediaProcessQueuesResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.CIURL,
		uri:      "/queue",
		optQuery: opt,
		method:   http.MethodGet,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type UpdateMediaProcessQueueOptions struct {
	XMLName      xml.Name                       `xml:"Request"`
	Name         string                         `xml:"Name,omitempty"`
	QueueID      string                         `xml:"QueueID,omitempty"`
	State        string                         `xml:"State,omitempty"`
	NotifyConfig *MediaProcessQueueNotifyConfig `xml:"NotifyConfig,omitempty"`
}

type UpdateMediaProcessQueueResult struct {
	XMLName   xml.Name           `xml:"Response"`
	RequestId string             `xml:"RequestId"`
	Queue     *MediaProcessQueue `xml:"Queue"`
}

func (s *CIService) UpdateMediaProcessQueue(ctx context.Context, opt *UpdateMediaProcessQueueOptions) (*UpdateMediaProcessQueueResult, *Response, error) {
	var res UpdateMediaProcessQueueResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.CIURL,
		uri:     "/queue/" + opt.QueueID,
		body:    opt,
		method:  http.MethodPut,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type DescribeMediaProcessBucketsOptions struct {
	Regions     string `url:"regions,omitempty"`
	BucketNames string `url:"bucketNames,omitempty"`
	BucketName  string `url:"bucketName,omitempty"`
	PageNumber  int    `url:"pageNumber,omitempty"`
	PageSize    int    `url:"pageSize,omitempty"`
}

type DescribeMediaProcessBucketsResult struct {
	XMLName         xml.Name             `xml:"Response"`
	RequestId       string               `xml:"RequestId,omitempty"`
	TotalCount      int                  `xml:"TotalCount,omitempty"`
	PageNumber      int                  `xml:"PageNumber,omitempty"`
	PageSize        int                  `xml:"PageSize,omitempty"`
	MediaBucketList []MediaProcessBucket `xml:"MediaBucketList,omitempty"`
}
type MediaProcessBucket struct {
	BucketId   string `xml:"BucketId,omitempty"`
	Region     string `xml:"Region,omitempty"`
	CreateTime string `xml:"CreateTime,omitempty"`
}

// 媒体bucket接口 https://cloud.tencent.com/document/product/436/48988
func (s *CIService) DescribeMediaProcessBuckets(ctx context.Context, opt *DescribeMediaProcessBucketsOptions) (*DescribeMediaProcessBucketsResult, *Response, error) {
	var res DescribeMediaProcessBucketsResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.CIURL,
		uri:      "/mediabucket",
		optQuery: opt,
		method:   http.MethodGet,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type GetMediaInfoResult struct {
	XMLName   xml.Name `xml:"Response"`
	MediaInfo struct {
		Format struct {
			Bitrate        float32 `xml:"Bitrate"`
			Duration       float32 `xml:"Duration"`
			FormatLongName string  `xml:"FormatLongName"`
			FormatName     string  `xml:"FormatName"`
			NumProgram     int     `xml:"NumProgram"`
			NumStream      int     `xml:"NumStream"`
			Size           int     `xml:"Size"`
			StartTime      float32 `xml:"StartTime"`
		} `xml:"Format"`
		Stream struct {
			Audio []struct {
				Index          int     `xml:"Index"`
				CodecName      string  `xml:"CodecName"`
				CodecLongName  string  `xml:"CodecLongName"`
				CodecTimeBase  string  `xml:"CodecTimeBase"`
				CodecTagString string  `xml:"CodecTagString"`
				CodecTag       string  `xml:"CodecTag"`
				SampleFmt      string  `xml:"SampleFmt"`
				SampleRate     int     `xml:"SampleRate"`
				Channel        int     `xml:"Channel"`
				ChannelLayout  string  `xml:"ChannelLayout"`
				Timebase       string  `xml:"Timebase"`
				StartTime      float32 `xml:"StartTime"`
				Duration       float32 `xml:"Duration"`
				Bitrate        float32 `xml:"Bitrate"`
				Language       string  `xml:"Language"`
			} `xml:"Audio"`
			Subtitle struct {
				Index    int    `xml:"Index"`
				Language string `xml:"Language"`
			} `xml:"Subtitle"`
			Video struct {
				Index          int     `xml:"Index"`
				CodecName      string  `xml:"CodecName"`
				CodecLongName  string  `xml:"CodecLongName"`
				CodecTimeBase  string  `xml:"CodecTimeBase"`
				CodecTagString string  `xml:"CodecTagString"`
				CodecTag       string  `xml:"CodecTag"`
				Profile        string  `xml:"Profile"`
				Height         int     `xml:"Height"`
				Width          int     `xml:"Width"`
				HasBFrame      int     `xml:"HasBFrame"`
				RefFrames      int     `xml:"RefFrames"`
				Sar            string  `xml:"Sar"`
				Dar            string  `xml:"Dar"`
				PixFormat      string  `xml:"PixFormat"`
				FieldOrder     string  `xml:"FieldOrder"`
				Level          int     `xml:"Level"`
				Fps            float32 `xml:"Fps"`
				AvgFps         string  `xml:"AvgFps"`
				Timebase       string  `xml:"Timebase"`
				StartTime      float32 `xml:"StartTime"`
				Duration       float32 `xml:"Duration"`
				Bitrate        float32 `xml:"Bitrate"`
				NumFrames      int     `xml:"NumFrames"`
				Language       string  `xml:"Language"`
			} `xml:"Video"`
		} `xml:"Stream"`
	} `xml:"MediaInfo"`
}

// 媒体信息接口 https://cloud.tencent.com/document/product/436/55672
func (s *CIService) GetMediaInfo(ctx context.Context, name string, opt *ObjectGetOptions, id ...string) (*GetMediaInfoResult, *Response, error) {
	var u string
	if len(id) == 1 {
		u = fmt.Sprintf("/%s?versionId=%s&ci-process=videoinfo", encodeURIComponent(name), id[0])
	} else if len(id) == 0 {
		u = fmt.Sprintf("/%s?ci-process=videoinfo", encodeURIComponent(name))
	} else {
		return nil, nil, fmt.Errorf("wrong params")
	}

	var res GetMediaInfoResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       u,
		method:    http.MethodGet,
		optQuery:  opt,
		optHeader: opt,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type GetSnapshotOptions struct {
	Time   float32 `url:"time,omitempty"`
	Height int     `url:"height,omitempty"`
	Width  int     `url:"width,omitempty"`
	Format string  `url:"format,omitempty"`
	Rotate string  `url:"rotate,omitempty"`
	Mode   string  `url:"mode,omitempty"`
}

// 媒体截图接口 https://cloud.tencent.com/document/product/436/55671
func (s *CIService) GetSnapshot(ctx context.Context, name string, opt *GetSnapshotOptions, id ...string) (*Response, error) {
	var u string
	if len(id) == 1 {
		u = fmt.Sprintf("/%s?versionId=%s&ci-process=snapshot", encodeURIComponent(name), id[0])
	} else if len(id) == 0 {
		u = fmt.Sprintf("/%s?ci-process=snapshot", encodeURIComponent(name))
	} else {
		return nil, fmt.Errorf("wrong params")
	}

	sendOpt := sendOptions{
		baseURL:          s.client.BaseURL.BucketURL,
		uri:              u,
		method:           http.MethodGet,
		optQuery:         opt,
		disableCloseBody: true,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
