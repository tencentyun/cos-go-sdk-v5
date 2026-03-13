package cos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VectorService 向量桶相关 API 服务
// 向量桶使用专用域名:
//   - 公网: vectors.<Region>.coslake.com
//   - 内网: vectors.<Region>.internal.tencentcos.com
//
// 所有接口均使用 JSON 格式进行请求和响应
type VectorService service

// NewVectorURL 生成 Vector 所需的基础 URL
//
//	region: 区域代码，如 ap-guangzhou, ap-shanghai, ap-beijing
//	secure: 是否使用 https
func NewVectorURL(region string, secure bool) (*url.URL, error) {
	schema := "https"
	if !secure {
		schema = "http"
	}
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	rawURL := fmt.Sprintf("%s://vectors.%s.coslake.com", schema, region)
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// NewVectorInternalURL 生成 Vector 内网访问的基础 URL
//
//	region: 区域代码，如 ap-guangzhou, ap-shanghai, ap-beijing
//	secure: 是否使用 https
func NewVectorInternalURL(region string, secure bool) (*url.URL, error) {
	schema := "https"
	if !secure {
		schema = "http"
	}
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	rawURL := fmt.Sprintf("%s://vectors.%s.internal.tencentcos.com", schema, region)
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// NewVectorEndpointURL 使用自定义 endpoint 生成 Vector 基础 URL
// 用户可以通过此方法使用任意自定义域名
//
//	endpoint: 完整的 endpoint URL，如 "https://vectors.ap-guangzhou.coslake.com"
func NewVectorEndpointURL(endpoint string) (*url.URL, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// vectorDoAPI 向量服务专用的 HTTP 请求执行方法
// 与 COS 的 doAPI 不同，Vector 服务：
// 1. 使用 JSON 格式而非 XML 格式
// 2. 错误响应体是 JSON 格式，错误码通过 X-Cos-Error-Code 头部返回
func (s *VectorService) vectorDoAPI(ctx context.Context, req *http.Request, result interface{}) (*Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := s.client.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	response := newResponse(resp)

	err = checkVectorResponse(resp)
	if err != nil {
		return response, err
	}

	if result != nil {
		if w, ok := result.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(result)
			if err == io.EOF {
				err = nil
			}
		}
	}

	return response, err
}

// vectorNewRequest 向量服务专用的 HTTP 请求构建方法
// 不做 COS 域名格式校验（Vector 使用独立域名），body 使用 JSON 序列化
func (s *VectorService) vectorNewRequest(ctx context.Context, uri, method string, body interface{}, isRetry bool) (*http.Request, error) {
	baseURL := s.client.BaseURL.VectorURL
	if baseURL == nil {
		return nil, fmt.Errorf("VectorURL is not set, please set BaseURL.VectorURL")
	}
	u, _ := url.Parse(uri)
	urlStr := baseURL.ResolveReference(u).String()

	var reader io.Reader
	contentLength := int64(0)
	if body != nil {
		if r, ok := body.(io.Reader); ok {
			reader = r
		} else {
			bs, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(bs)
			contentLength = int64(len(bs))
		}
	}

	req, err := http.NewRequest(method, urlStr, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if s.client.UserAgent != "" {
		req.Header.Set("User-Agent", s.client.UserAgent)
	}
	if contentLength > 0 {
		req.ContentLength = contentLength
	}
	if isRetry {
		req.Header.Set("X-Cos-Sdk-Retry", "true")
	}
	if s.client.Host != "" {
		req.Host = s.client.Host
	}
	if s.client.Conf.RequestBodyClose {
		req.Close = true
	}
	return req, nil
}

// vectorSend 向量服务专用的请求发送方法
// 认证信息由 http.Client.Transport（如 AuthorizationTransport）自动注入
func (s *VectorService) vectorSend(ctx context.Context, uri, method string, body interface{}, result interface{}, isRetry bool) (*Response, error) {
	req, err := s.vectorNewRequest(ctx, uri, method, body, isRetry)
	if err != nil {
		return nil, err
	}

	return s.vectorDoAPI(ctx, req, result)
}

// vectorCheckRetrieable 向量服务专用的重试判断
// 与 COS 不同：不切换域名，只在 5xx 或网络错误时重试
func (s *VectorService) vectorCheckRetrieable(resp *Response, err error) bool {
	if err == nil {
		return false
	}
	// VectorErrorResponse 类型的错误，根据状态码判断
	if verr, ok := err.(*VectorErrorResponse); ok {
		if verr.Response != nil && verr.Response.StatusCode >= 500 {
			return true
		}
		return false
	}
	// 网络错误等非业务错误，重试
	return true
}

// vectorDoRetry 向量服务专用的重试逻辑
// 不切换域名，仅在网络错误或 5xx 时进行同域名重试
func (s *VectorService) vectorDoRetry(ctx context.Context, uri, method string, body interface{}, result interface{}) (*Response, error) {
	// 如果 body 是 io.Reader（流式），不支持重试
	if body != nil {
		if _, ok := body.(io.Reader); ok {
			return s.vectorSend(ctx, uri, method, body, result, false)
		}
	}

	count := 1
	if s.client.Conf.RetryOpt.Count > 0 {
		count = s.client.Conf.RetryOpt.Count
	}

	retryErr := &RetryError{}
	var resp *Response
	var err error

	for nr := 0; nr < count; nr++ {
		if err != nil {
			retryErr.Add(err)
		}
		isRetry := nr > 0
		resp, err = s.vectorSend(ctx, uri, method, body, result, isRetry)
		if s.vectorCheckRetrieable(resp, err) {
			if s.client.Conf.RetryOpt.Interval > 0 && nr+1 < count {
				time.Sleep(s.client.Conf.RetryOpt.Interval)
			}
			continue
		}
		break
	}

	// 最后一次非 Vector 错误，输出所有重试结果
	if err != nil {
		if _, ok := err.(*VectorErrorResponse); !ok {
			retryErr.Add(err)
			err = retryErr
		}
	}
	return resp, err
}

// baseSend 向量服务统一的请求发送入口，支持重试
func (s *VectorService) baseSend(ctx context.Context, opt interface{}, uri string, method string) (*bytes.Buffer, *Response, error) {
	var buf bytes.Buffer
	resp, err := s.vectorDoRetry(ctx, uri, method, opt, &buf)
	return &buf, resp, err
}

// ==================== 向量桶管理 ====================

// CreateVectorBucketOptions 创建向量桶请求参数
type CreateVectorBucketOptions struct {
	VectorBucketName        string                  `json:"vectorBucketName"`                  // 向量桶名称，格式为 BucketName-APPID
	EncryptionConfiguration *VectorEncryptionConfig  `json:"encryptionConfiguration,omitempty"` // 加密配置
}

// VectorEncryptionConfig 加密配置
type VectorEncryptionConfig struct {
	SseType string `json:"sseType"` // 加密类型，当前仅支持 AES256
}

// CreateVectorBucketResult 创建向量桶响应
type CreateVectorBucketResult struct {
	VectorBucketQcs string `json:"vectorBucketQcs"` // 向量桶资源名称 (QCS)
}

// CreateVectorBucket 创建向量桶
// https://cloud.tencent.com/document/product/436/127725
func (s *VectorService) CreateVectorBucket(ctx context.Context, opt *CreateVectorBucketOptions) (*CreateVectorBucketResult, *Response, error) {
	var res CreateVectorBucketResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, "/CreateVectorBucket", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// GetVectorBucketOptions 查询向量桶请求参数
type GetVectorBucketOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
}

// VectorBucketInfo 向量桶信息
type VectorBucketInfo struct {
	CreationTime            int64                   `json:"creationTime"`                      // 创建时间戳
	EncryptionConfiguration *VectorEncryptionConfig `json:"encryptionConfiguration,omitempty"` // 加密配置
	VectorBucketQcs         string                  `json:"vectorBucketQcs"`                   // 资源名称
	VectorBucketName        string                  `json:"vectorBucketName"`                  // 向量桶名称
}

// GetVectorBucketResult 查询向量桶响应
type GetVectorBucketResult struct {
	VectorBucket *VectorBucketInfo `json:"vectorBucket"` // 向量桶信息
}

// GetVectorBucket 查询向量桶信息
// https://cloud.tencent.com/document/product/436/127726
func (s *VectorService) GetVectorBucket(ctx context.Context, opt *GetVectorBucketOptions) (*GetVectorBucketResult, *Response, error) {
	var res GetVectorBucketResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, "/GetVectorBucket", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// DeleteVectorBucketOptions 删除向量桶请求参数
type DeleteVectorBucketOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
}

// DeleteVectorBucket 删除向量桶
// https://cloud.tencent.com/document/product/436/127728
func (s *VectorService) DeleteVectorBucket(ctx context.Context, opt *DeleteVectorBucketOptions) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt param nil")
	}
	_, resp, err := s.baseSend(ctx, opt, "/DeleteVectorBucket", http.MethodPost)
	return resp, err
}

// ListVectorBucketsOptions 列出所有向量桶请求参数
type ListVectorBucketsOptions struct {
	MaxResults int    `json:"maxResults,omitempty"` // 最大返回数量，默认100，最大500
	NextToken  string `json:"nextToken,omitempty"`  // 分页标记
	Prefix     string `json:"prefix,omitempty"`     // 桶名前缀过滤
}

// ListVectorBucketsResult 列出所有向量桶响应
type ListVectorBucketsResult struct {
	NextToken     string             `json:"nextToken,omitempty"` // 下一页分页标记
	VectorBuckets []VectorBucketBrief `json:"vectorBuckets"`      // 向量桶列表
}

// VectorBucketBrief 向量桶简要信息
type VectorBucketBrief struct {
	CreationTime     int64  `json:"creationTime"`     // 创建时间戳
	VectorBucketQcs  string `json:"vectorBucketQcs"`  // 资源名称
	VectorBucketName string `json:"vectorBucketName"` // 向量桶名称
}

// ListVectorBuckets 列出所有向量桶
// https://cloud.tencent.com/document/product/436/127727
func (s *VectorService) ListVectorBuckets(ctx context.Context, opt *ListVectorBucketsOptions) (*ListVectorBucketsResult, *Response, error) {
	var res ListVectorBucketsResult
	buf, resp, err := s.baseSend(ctx, opt, "/ListVectorBuckets", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// ==================== 向量桶策略管理 ====================

// PutVectorBucketPolicyOptions 设置向量桶策略请求参数
type PutVectorBucketPolicyOptions struct {
	VectorBucketName string      `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string      `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
	Policy           interface{} `json:"policy"`                     // 策略内容
}

// PutVectorBucketPolicy 设置向量桶策略
func (s *VectorService) PutVectorBucketPolicy(ctx context.Context, opt *PutVectorBucketPolicyOptions) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt param nil")
	}
	_, resp, err := s.baseSend(ctx, opt, "/PutVectorBucketPolicy", http.MethodPost)
	return resp, err
}

// GetVectorBucketPolicyOptions 获取向量桶策略请求参数
type GetVectorBucketPolicyOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
}

// GetVectorBucketPolicyResult 获取向量桶策略响应
type GetVectorBucketPolicyResult struct {
	Policy interface{} `json:"policy"` // 策略内容
}

// GetVectorBucketPolicy 获取向量桶策略
func (s *VectorService) GetVectorBucketPolicy(ctx context.Context, opt *GetVectorBucketPolicyOptions) (*GetVectorBucketPolicyResult, *Response, error) {
	var res GetVectorBucketPolicyResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, "/GetVectorBucketPolicy", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// DeleteVectorBucketPolicyOptions 删除向量桶策略请求参数
type DeleteVectorBucketPolicyOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
}

// DeleteVectorBucketPolicy 删除向量桶策略
func (s *VectorService) DeleteVectorBucketPolicy(ctx context.Context, opt *DeleteVectorBucketPolicyOptions) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt param nil")
	}
	_, resp, err := s.baseSend(ctx, opt, "/DeleteVectorBucketPolicy", http.MethodPost)
	return resp, err
}

// ==================== 向量索引管理 ====================

// IndexParams 索引构建参数
type IndexParams struct {
	EfConstruction int `json:"efConstruction,omitempty"` // 构建索引时的搜索范围，默认200，范围 [1, 500]
	M              int `json:"m,omitempty"`              // 图的连接度，默认16，范围 [2, 100]
}

// CreateIndexOptions 创建索引请求选项（指定目标向量桶）
type CreateIndexOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
}

// IndexDefinition 索引定义（索引的业务配置数据）
type IndexDefinition struct {
	IndexName   string       `json:"indexName"`             // 索引名称
	Dimension   int          `json:"dimension"`             // 向量维度
	Metric      string       `json:"metric"`                // 距离度量方式: L2, INNER_PRODUCT, COSINE
	Params      *IndexParams `json:"params,omitempty"`      // 索引构建参数
	Description string       `json:"description,omitempty"` // 索引描述信息
}

// createIndexRequest 创建索引的内部请求体（合并 opt + indexDef）
type createIndexRequest struct {
	VectorBucketName string       `json:"vectorBucketName,omitempty"`
	VectorBucketQcs  string       `json:"vectorBucketQcs,omitempty"`
	IndexName        string       `json:"indexName"`
	Dimension        int          `json:"dimension"`
	Metric           string       `json:"metric"`
	Params           *IndexParams `json:"params,omitempty"`
	Description      string       `json:"description,omitempty"`
}

// IndexInfo 索引信息
type IndexInfo struct {
	IndexQcs         string       `json:"indexQcs"`                   // 索引资源名称 (QCS)
	IndexName        string       `json:"indexName"`                  // 索引名称
	VectorBucketName string       `json:"vectorBucketName"`           // 向量桶名称
	CreationTime     int64        `json:"creationTime"`               // 创建时间戳
	Dimension        int          `json:"dimension"`                  // 向量维度
	Metric           string       `json:"metric"`                     // 距离度量方式
	Params           *IndexParams `json:"params,omitempty"`           // 索引参数
	Description      string       `json:"description,omitempty"`      // 描述信息
	Status           string       `json:"status,omitempty"`           // 索引状态
}

// CreateIndexResult 创建索引响应
type CreateIndexResult struct {
	Index *IndexInfo `json:"index"` // 索引信息
}

// CreateIndex 创建向量索引
//
//	opt: 请求选项（指定目标向量桶）
//	indexDef: 索引定义（名称、维度、度量方式等）
//
// https://cloud.tencent.com/document/product/436/127715
func (s *VectorService) CreateIndex(ctx context.Context, opt *CreateIndexOptions, indexDef *IndexDefinition) (*CreateIndexResult, *Response, error) {
	var res CreateIndexResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	if indexDef == nil {
		return nil, nil, fmt.Errorf("indexDef param nil")
	}
	reqBody := &createIndexRequest{
		VectorBucketName: opt.VectorBucketName,
		VectorBucketQcs:  opt.VectorBucketQcs,
		IndexName:        indexDef.IndexName,
		Dimension:        indexDef.Dimension,
		Metric:           indexDef.Metric,
		Params:           indexDef.Params,
		Description:      indexDef.Description,
	}
	buf, resp, err := s.baseSend(ctx, reqBody, "/CreateIndex", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// GetIndexOptions 查询索引请求参数
type GetIndexOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string `json:"indexName,omitempty"`        // 索引名称
}

// GetIndexResult 查询索引响应
type GetIndexResult struct {
	Index *IndexInfo `json:"index"` // 索引信息
}

// GetIndex 查询索引信息
func (s *VectorService) GetIndex(ctx context.Context, opt *GetIndexOptions) (*GetIndexResult, *Response, error) {
	var res GetIndexResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, "/GetIndex", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// ListIndexesOptions 列出索引请求参数
type ListIndexesOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	VectorBucketQcs  string `json:"vectorBucketQcs,omitempty"`  // 向量桶资源名称 (QCS)
	MaxResults       int    `json:"maxResults,omitempty"`       // 最大返回数量
	NextToken        string `json:"nextToken,omitempty"`        // 分页标记
	Prefix           string `json:"prefix,omitempty"`           // 索引名前缀过滤
}

// IndexBrief 索引简要信息
type IndexBrief struct {
	CreationTime     int64  `json:"creationTime"`     // 创建时间戳
	IndexQcs         string `json:"indexQcs"`          // 索引资源名称 (QCS)
	IndexName        string `json:"indexName"`         // 索引名称
	VectorBucketName string `json:"vectorBucketName"`  // 向量桶名称
}

// ListIndexesResult 列出索引响应
type ListIndexesResult struct {
	Indexes   []IndexBrief `json:"indexes"`             // 索引列表
	NextToken string       `json:"nextToken,omitempty"` // 下一页分页标记
}

// ListIndexes 列出所有索引
// https://cloud.tencent.com/document/product/436/127729
func (s *VectorService) ListIndexes(ctx context.Context, opt *ListIndexesOptions) (*ListIndexesResult, *Response, error) {
	var res ListIndexesResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, "/ListIndexes", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// DeleteIndexOptions 删除索引请求参数
type DeleteIndexOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string `json:"indexName"`                  // 索引名称
}

// DeleteIndex 删除索引
// https://cloud.tencent.com/document/product/436/127716
func (s *VectorService) DeleteIndex(ctx context.Context, opt *DeleteIndexOptions) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt param nil")
	}
	_, resp, err := s.baseSend(ctx, opt, "/DeleteIndex", http.MethodPost)
	return resp, err
}

// ==================== 向量数据操作 ====================

// VectorData 向量数据
type VectorData struct {
	Float32 []float32 `json:"float32,omitempty"` // float32 类型的向量数据
}

// Vector 完整的向量信息
type Vector struct {
	Key      string                 `json:"key"`                // 向量主键
	Data     *VectorData            `json:"data,omitempty"`     // 向量数据
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据
}

// PutVectorsOptions 插入/更新向量请求选项（指定目标索引）
type PutVectorsOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string `json:"indexName,omitempty"`        // 索引名称
}

// putVectorsRequest 插入向量的内部请求体（合并 opt + vectors）
type putVectorsRequest struct {
	VectorBucketName string   `json:"vectorBucketName,omitempty"`
	IndexQcs         string   `json:"indexQcs,omitempty"`
	IndexName        string   `json:"indexName,omitempty"`
	Vectors          []Vector `json:"vectors"`
}

// PutVectors 插入或更新向量数据
//
//	opt: 请求选项（指定目标索引）
//	vectors: 向量数据列表，最大500条
//
// https://cloud.tencent.com/document/product/436/127719
func (s *VectorService) PutVectors(ctx context.Context, opt *PutVectorsOptions, vectors []Vector) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt param nil")
	}
	if len(vectors) == 0 {
		return nil, fmt.Errorf("vectors param is empty")
	}
	reqBody := &putVectorsRequest{
		VectorBucketName: opt.VectorBucketName,
		IndexQcs:         opt.IndexQcs,
		IndexName:        opt.IndexName,
		Vectors:          vectors,
	}
	_, resp, err := s.baseSend(ctx, reqBody, "/PutVectors", http.MethodPost)
	return resp, err
}

// GetVectorsOptions 获取指定向量请求选项（指定目标索引和返回控制）
type GetVectorsOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string `json:"indexName,omitempty"`        // 索引名称
	ReturnData       bool   `json:"returnData,omitempty"`       // 是否返回向量数据
	ReturnMetadata   bool   `json:"returnMetadata,omitempty"`   // 是否返回元数据
}

// getVectorsRequest 获取向量的内部请求体（合并 opt + keys）
type getVectorsRequest struct {
	VectorBucketName string   `json:"vectorBucketName,omitempty"`
	IndexQcs         string   `json:"indexQcs,omitempty"`
	IndexName        string   `json:"indexName,omitempty"`
	Keys             []string `json:"keys"`
	ReturnData       bool     `json:"returnData,omitempty"`
	ReturnMetadata   bool     `json:"returnMetadata,omitempty"`
}

// GetVectorsResult 获取向量响应
type GetVectorsResult struct {
	Vectors []Vector `json:"vectors"` // 向量列表
}

// GetVectors 获取指定向量
//
//	opt: 请求选项（指定目标索引和返回控制）
//	keys: 向量主键列表，最大100个
//
// https://cloud.tencent.com/document/product/436/127720
func (s *VectorService) GetVectors(ctx context.Context, opt *GetVectorsOptions, keys []string) (*GetVectorsResult, *Response, error) {
	var res GetVectorsResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	if len(keys) == 0 {
		return nil, nil, fmt.Errorf("keys param is empty")
	}
	reqBody := &getVectorsRequest{
		VectorBucketName: opt.VectorBucketName,
		IndexQcs:         opt.IndexQcs,
		IndexName:        opt.IndexName,
		Keys:             keys,
		ReturnData:       opt.ReturnData,
		ReturnMetadata:   opt.ReturnMetadata,
	}
	buf, resp, err := s.baseSend(ctx, reqBody, "/GetVectors", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// ListVectorsOptions 列出向量请求参数
type ListVectorsOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string `json:"indexName,omitempty"`        // 索引名称
	MaxResults       int    `json:"maxResults,omitempty"`       // 最大返回数量
	NextToken        string `json:"nextToken,omitempty"`        // 分页标记
	ReturnData       bool   `json:"returnData,omitempty"`       // 是否返回向量数据
	ReturnMetadata   bool   `json:"returnMetadata,omitempty"`   // 是否返回元数据
	SegmentCount     int    `json:"segmentCount,omitempty"`     // 分段总数
	SegmentIndex     int    `json:"segmentIndex,omitempty"`     // 分段索引（从0开始）
}

// ListVectorsResult 列出向量响应
type ListVectorsResult struct {
	Vectors   []Vector `json:"vectors"`             // 向量列表
	NextToken string   `json:"nextToken,omitempty"` // 下一页分页标记
}

// ListVectors 列出向量列表
// https://cloud.tencent.com/document/product/436/127721
func (s *VectorService) ListVectors(ctx context.Context, opt *ListVectorsOptions) (*ListVectorsResult, *Response, error) {
	var res ListVectorsResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	buf, resp, err := s.baseSend(ctx, opt, "/ListVectors", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}

// DeleteVectorsOptions 删除向量请求选项（指定目标索引）
type DeleteVectorsOptions struct {
	VectorBucketName string `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string `json:"indexName,omitempty"`        // 索引名称
}

// deleteVectorsRequest 删除向量的内部请求体（合并 opt + keys）
type deleteVectorsRequest struct {
	VectorBucketName string   `json:"vectorBucketName,omitempty"`
	IndexQcs         string   `json:"indexQcs,omitempty"`
	IndexName        string   `json:"indexName,omitempty"`
	Keys             []string `json:"keys"`
}

// DeleteVectors 删除指定向量
//
//	opt: 请求选项（指定目标索引）
//	keys: 要删除的向量主键列表，最大500个
//
// https://cloud.tencent.com/document/product/436/127722
func (s *VectorService) DeleteVectors(ctx context.Context, opt *DeleteVectorsOptions, keys []string) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt param nil")
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("keys param is empty")
	}
	reqBody := &deleteVectorsRequest{
		VectorBucketName: opt.VectorBucketName,
		IndexQcs:         opt.IndexQcs,
		IndexName:        opt.IndexName,
		Keys:             keys,
	}
	_, resp, err := s.baseSend(ctx, reqBody, "/DeleteVectors", http.MethodPost)
	return resp, err
}

// QueryVectorsOptions 相似度搜索请求选项（指定目标索引和搜索控制）
type QueryVectorsOptions struct {
	VectorBucketName string      `json:"vectorBucketName,omitempty"` // 向量桶名称
	IndexQcs         string      `json:"indexQcs,omitempty"`         // 索引资源名称 (QCS)
	IndexName        string      `json:"indexName,omitempty"`        // 索引名称
	Filter           interface{} `json:"filter,omitempty"`           // 过滤条件
	ReturnData       bool        `json:"returnData,omitempty"`       // 是否返回向量数据
	ReturnMetadata   bool        `json:"returnMetadata,omitempty"`   // 是否返回元数据
	ReturnDistance   bool        `json:"returnDistance,omitempty"`   // 是否返回距离值
}

// queryVectorsRequest 相似度搜索的内部请求体（合并 opt + queryVector + topK）
type queryVectorsRequest struct {
	VectorBucketName string      `json:"vectorBucketName,omitempty"`
	IndexQcs         string      `json:"indexQcs,omitempty"`
	IndexName        string      `json:"indexName,omitempty"`
	QueryVector      *VectorData `json:"queryVector"`
	TopK             int         `json:"topK"`
	Filter           interface{} `json:"filter,omitempty"`
	ReturnData       bool        `json:"returnData,omitempty"`
	ReturnMetadata   bool        `json:"returnMetadata,omitempty"`
	ReturnDistance   bool        `json:"returnDistance,omitempty"`
}

// QueryVector 查询结果向量（包含距离）
type QueryVector struct {
	Key      string                 `json:"key"`                // 向量主键
	Data     *VectorData            `json:"data,omitempty"`     // 向量数据
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据
	Distance float64                `json:"distance,omitempty"` // 距离值
}

// QueryVectorsResult 相似度搜索响应
type QueryVectorsResult struct {
	Vectors []QueryVector `json:"vectors"` // 结果向量列表
}

// QueryVectors 向量相似度搜索
//
//	opt: 请求选项（指定目标索引和搜索控制）
//	queryVector: 查询向量
//	topK: 返回最相似的 K 个结果，范围 1~30
//
// https://cloud.tencent.com/document/product/436/127723
func (s *VectorService) QueryVectors(ctx context.Context, opt *QueryVectorsOptions, queryVector *VectorData, topK int) (*QueryVectorsResult, *Response, error) {
	var res QueryVectorsResult
	if opt == nil {
		return nil, nil, fmt.Errorf("opt param nil")
	}
	if queryVector == nil {
		return nil, nil, fmt.Errorf("queryVector param nil")
	}
	if topK <= 0 {
		return nil, nil, fmt.Errorf("topK must be greater than 0")
	}
	reqBody := &queryVectorsRequest{
		VectorBucketName: opt.VectorBucketName,
		IndexQcs:         opt.IndexQcs,
		IndexName:        opt.IndexName,
		QueryVector:      queryVector,
		TopK:             topK,
		Filter:           opt.Filter,
		ReturnData:       opt.ReturnData,
		ReturnMetadata:   opt.ReturnMetadata,
		ReturnDistance:    opt.ReturnDistance,
	}
	buf, resp, err := s.baseSend(ctx, reqBody, "/QueryVectors", http.MethodPost)
	if err != nil {
		return nil, resp, err
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &res)
	}
	return &res, resp, err
}
