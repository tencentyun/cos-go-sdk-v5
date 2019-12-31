package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

type BatchService service

type BatchRequestHeaders struct {
	XCosAppid     int          `header:"x-cos-appid"`
	ContentLength string       `header:"Content-Length,omitempty"`
	ContentType   string       `header:"Content-Type,omitempty"`
	Headers       *http.Header `header:"-"`
}

// BatchProgressSummary
type BatchProgressSummary struct {
	NumberOfTasksFailed    int `xml:"NumberOfTasksFailed"`
	NumberOfTasksSucceeded int `xml:"NumberOfTasksSucceeded"`
	TotalNumberOfTasks     int `xml:"TotalNumberOfTasks"`
}

// BatchJobReport
type BatchJobReport struct {
	Bucket      string `xml:"Bucket"`
	Enabled     string `xml:"Enabled"`
	Format      string `xml:"Format"`
	Prefix      string `xml:"Prefix,omitempty"`
	ReportScope string `xml:"ReportScope"`
}

// BatchJobOperationCopy
type BatchMetadata struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}
type BatchNewObjectMetadata struct {
	CacheControl       string          `xml:"CacheControl,omitempty"`
	ContentDisposition string          `xml:"ContentDisposition,omitempty"`
	ContentEncoding    string          `xml:"ContentEncoding,omitempty"`
	ContentType        string          `xml:"ContentType,omitempty"`
	HttpExpiresDate    string          `xml:"HttpExpiresDate,omitempty"`
	SSEAlgorithm       string          `xml:"SSEAlgorithm,omitempty"`
	UserMetadata       []BatchMetadata `xml:"UserMetadata>member,omitempty"`
}
type BatchGrantee struct {
	DisplayName    string `xml:"DisplayName,omitempty"`
	Identifier     string `xml:"Identifier"`
	TypeIdentifier string `xml:"TypeIdentifier"`
}
type BatchCOSGrant struct {
	Grantee    *BatchGrantee `xml:"Grantee"`
	Permission string        `xml:"Permission"`
}
type BatchAccessControlGrants struct {
	COSGrants *BatchCOSGrant `xml:"COSGrant,omitempty"`
}
type BatchJobOperationCopy struct {
	AccessControlGrants       *BatchAccessControlGrants `xml:"AccessControlGrants,omitempty"`
	CannedAccessControlList   string                    `xml:"CannedAccessControlList,omitempty"`
	MetadataDirective         string                    `xml:"MetadataDirective,omitempty"`
	ModifiedSinceConstraint   int64                     `xml:"ModifiedSinceConstraint,omitempty"`
	UnModifiedSinceConstraint int64                     `xml:"UnModifiedSinceConstraint,omitempty"`
	NewObjectMetadata         *BatchNewObjectMetadata   `xml:"NewObjectMetadata,omitempty"`
	StorageClass              string                    `xml:"StorageClass,omitempty"`
	TargetResource            string                    `xml:"TargetResource"`
}

// BatchJobOperation
type BatchJobOperation struct {
	PutObjectCopy *BatchJobOperationCopy `xml:"COSPutObjectCopy,omitempty" header:"-"`
}

// BatchJobManifest
type BatchJobManifestLocation struct {
	ETag            string `xml:"ETag" header:"-"`
	ObjectArn       string `xml:"ObjectArn" header:"-"`
	ObjectVersionId string `xml:"ObjectVersionId,omitempty" header:"-"`
}
type BatchJobManifestSpec struct {
	Fields []string `xml:"Fields>member,omitempty" header:"-"`
	Format string   `xml:"Format" header:"-"`
}
type BatchJobManifest struct {
	Location *BatchJobManifestLocation `xml:"Location" header:"-"`
	Spec     *BatchJobManifestSpec     `xml:"Spec" header:"-"`
}

type BatchCreateJobOptions struct {
	XMLName              xml.Name           `xml:"CreateJobRequest" header:"-"`
	ClientRequestToken   string             `xml:"ClientRequestToken" header:"-"`
	ConfirmationRequired string             `xml:"ConfirmationRequired,omitempty" header:"-"`
	Description          string             `xml:"Description,omitempty" header:"-"`
	Manifest             *BatchJobManifest  `xml:"Manifest" header:"-"`
	Operation            *BatchJobOperation `xml:"Operation" header:"-"`
	Priority             int                `xml:"Priority" header:"-"`
	Report               *BatchJobReport    `xml:"Report" header:"-"`
	RoleArn              string             `xml:"RoleArn" header:"-"`
}

type BatchCreateJobResult struct {
	XMLName xml.Name `xml:"CreateJobResult"`
	JobId   string   `xml:"JobId"`
}

func processETag(opt *BatchCreateJobOptions) *BatchCreateJobOptions {
	if opt != nil && opt.Manifest != nil && opt.Manifest.Location != nil {
		opt.Manifest.Location.ETag = "<ETag>" + opt.Manifest.Location.ETag + "</ETag>"
	}
	return opt
}

func (s *BatchService) CreateJob(ctx context.Context, opt *BatchCreateJobOptions, headers *BatchRequestHeaders) (*BatchCreateJobResult, *Response, error) {
	var res BatchCreateJobResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       "/jobs",
		method:    http.MethodPost,
		optHeader: headers,
		body:      opt,
		result:    &res,
	}

	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchJobFailureReasons struct {
	FailureCode   string `xml:"FailureCode"`
	FailureReason string `xml:"FailureReason"`
}

type BatchDescribeJob struct {
	ConfirmationRequired string                  `xml:"ConfirmationRequired,omitempty"`
	CreationTime         string                  `xml:"CreationTime,omitempty"`
	Description          string                  `xml:"Description,omitempty"`
	FailureReasons       *BatchJobFailureReasons `xml:"FailureReasons>JobFailure,omitempty"`
	JobId                string                  `xml:"JobId"`
	Manifest             *BatchJobManifest       `xml:"Manifest"`
	Operation            *BatchJobOperation      `xml:"Operation"`
	Priority             int                     `xml:"Priority"`
	ProgressSummary      *BatchProgressSummary   `xml:"ProgressSummary"`
	Report               *BatchJobReport         `xml:"Report,omitempty"`
	RoleArn              string                  `xml:"RoleArn,omitempty"`
	Status               string                  `xml:"Status,omitempty"`
	StatusUpdateReason   string                  `xml:"StatusUpdateReason,omitempty"`
	SuspendedCause       string                  `xml:"SuspendedCause,omitempty"`
	SuspendedDate        string                  `xml:"SuspendedDate,omitempty"`
	TerminationDate      string                  `xml:"TerminationDate,omitempty"`
}
type BatchDescribeJobResult struct {
	XMLName xml.Name          `xml:"DescribeJobResult"`
	Job     *BatchDescribeJob `xml:"Job"`
}

func (s *BatchService) DescribeJob(ctx context.Context, id string, headers *BatchRequestHeaders) (*BatchDescribeJobResult, *Response, error) {
	var res BatchDescribeJobResult
	u := fmt.Sprintf("/jobs/%s", id)
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       u,
		method:    http.MethodGet,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchListJobsOptions struct {
	JobStatuses string `url:"jobStatuses,omitempty" header:"-"`
	MaxResults  int    `url:"maxResults,omitempty" header:"-"`
	NextToken   string `url:"nextToken,omitempty" header:"-"`
}

type BatchListJobsMember struct {
	CreationTime    string                `xml:"CreationTime,omitempty"`
	Description     string                `xml:"Description,omitempty"`
	JobId           string                `xml:"JobId,omitempty"`
	Operation       string                `xml:"Operation,omitempty"`
	Priority        int                   `xml:"Priority,omitempty"`
	ProgressSummary *BatchProgressSummary `xml:"ProgressSummary,omitempty"`
	Status          string                `xml:"Status,omitempty"`
	TerminationDate string                `xml:"TerminationDate,omitempty"`
}
type BatchListJobs struct {
	Members []BatchListJobsMember `xml:"member,omitempty"`
}
type BatchListJobsResult struct {
	XMLName   xml.Name       `xml:"ListJobsResult"`
	Jobs      *BatchListJobs `xml:"Jobs"`
	NextToken string         `xml:"NextToken,omitempty"`
}

func (s *BatchService) ListJobs(ctx context.Context, opt *BatchListJobsOptions, headers *BatchRequestHeaders) (*BatchListJobsResult, *Response, error) {
	var res BatchListJobsResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       "/jobs",
		method:    http.MethodGet,
		optQuery:  opt,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchUpdatePriorityOptions struct {
	JobId    string `url:"-" header:"-" xml:"-"`
	Priority int    `url:"priority" header:"-" xml:"-"`
}
type BatchUpdatePriorityResult struct {
	XMLName  xml.Name `xml:"UpdateJobPriorityResult"`
	JobId    string   `xml:"JobId,omitempty"`
	Priority int      `xml:"Priority,omitempty"`
}

func (s *BatchService) UpdateJobPriority(ctx context.Context, opt *BatchUpdatePriorityOptions, headers *BatchRequestHeaders) (*BatchUpdatePriorityResult, *Response, error) {
	u := fmt.Sprintf("/jobs/%s/priority", opt.JobId)
	var res BatchUpdatePriorityResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       u,
		method:    http.MethodPost,
		optQuery:  opt,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type BatchUpdateStatusOptions struct {
	JobId              string `header:"-" url:"-" xml:"-"`
	RequestedJobStatus string `url:"requestedJobStatus" header:"-" xml:"-"`
	StatusUpdateReason string `url:"statusUpdateReason,omitempty" header:"-", xml:"-"`
}
type BatchUpdateStatusResult struct {
	XMLName            xml.Name `xml:"UpdateJobStatusResult"`
	JobId              string   `xml:"JobId,omitempty"`
	Status             string   `xml:"Status,omitempty"`
	StatusUpdateReason string   `xml:"StatusUpdateReason,omitempty"`
}

func (s *BatchService) UpdateJobStatus(ctx context.Context, opt *BatchUpdateStatusOptions, headers *BatchRequestHeaders) (*BatchUpdateStatusResult, *Response, error) {
	u := fmt.Sprintf("/jobs/%s/status", opt.JobId)
	var res BatchUpdateStatusResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BatchURL,
		uri:       u,
		method:    http.MethodPost,
		optQuery:  opt,
		optHeader: headers,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}
