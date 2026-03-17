package cos

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type RetryError struct {
	Errs []error
}

func (r *RetryError) Error() string {
	var errStr []string
	for _, err := range r.Errs {
		errStr = append(errStr, err.Error())
	}
	return strings.Join(errStr, "; ")
}

func (r *RetryError) Add(err error) {
	r.Errs = append(r.Errs, err)
}

// ErrorResponse 包含 API 返回的错误信息
//
// https://www.qcloud.com/document/product/436/7730
type ErrorResponse struct {
	XMLName   xml.Name       `xml:"Error"`
	Response  *http.Response `xml:"-"`
	Code      string
	Message   string
	Resource  string
	RequestID string `header:"x-cos-request-id,omitempty" url:"-" xml:"RequestId,omitempty"`
	TraceID   string `xml:"TraceId,omitempty"`
}

// Error returns the error msg
func (r *ErrorResponse) Error() string {
	RequestID := r.RequestID
	if RequestID == "" {
		RequestID = r.Response.Header.Get("X-Cos-Request-Id")
	}
	TraceID := r.TraceID
	if TraceID == "" {
		TraceID = r.Response.Header.Get("X-Cos-Trace-Id")
	}
	decodeURL, err := decodeURIComponent(r.Response.Request.URL.String())
	if err != nil {
		decodeURL = r.Response.Request.URL.String()
	}
	return fmt.Sprintf("%v %v: %d %v(Message: %v, RequestId: %v, TraceId: %v)",
		r.Response.Request.Method, decodeURL,
		r.Response.StatusCode, r.Code, r.Message, RequestID, TraceID)
}

type jsonError struct {
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// 检查 response 是否是出错时的返回的 response
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		xml.Unmarshal(data, errorResponse)
	}
	// 是否为 json 格式
	if errorResponse.Code == "" {
		ctype := strings.TrimLeft(r.Header.Get("Content-Type"), " ")
		if strings.HasPrefix(ctype, "application/json") {
			var jerror jsonError
			json.Unmarshal(data, &jerror)
			errorResponse.Code = strconv.Itoa(jerror.Code)
			errorResponse.Message = jerror.Message
			errorResponse.RequestID = jerror.RequestID
		}

	}
	return errorResponse
}

func IsNotFoundError(e error) bool {
	if e == nil {
		return false
	}
	err, ok := e.(*ErrorResponse)
	if !ok {
		return false
	}
	if err.Response != nil && err.Response.StatusCode == 404 {
		return true
	}
	return false
}

func IsCOSError(e error) (*ErrorResponse, bool) {
	if e == nil {
		return nil, false
	}
	err, ok := e.(*ErrorResponse)
	return err, ok
}

// ==================== Vector 专用错误处理 ====================

// VectorValidateField 参数校验失败的字段信息
type VectorValidateField struct {
	Message string `json:"message"` // 参数校验失败原因
	Path    string `json:"path"`    // 参数字段在请求结构中的位置
}

// VectorErrorResponse 包含 Vector API 返回的 JSON 格式错误信息
// 与 COS 的 XML 错误格式不同，Vector 错误码通过响应头 X-Cos-Error-Code 返回
//
// 错误码列表:
//
//	ValidationException (400)        - 请求不合法
//	ServiceQuotaExceededException (402) - 请求超过服务配额
//	AccessDeniedException (403)      - 访问被拒绝
//	NotFoundException (404)          - 资源不存在
//	ConflictException (409)          - 资源冲突
//	TooManyRequestsException (429)   - 请求太多超过限制
//	InternalServerException (500)    - 服务内部错误
//	ServiceUnavailableException (503)- 服务不可用，请重试
//	KmsDisabledException (400)       - KMS 不可用
//	KmsInvalidKeyUsageException (400)- KMS key 不兼容
//	KmsInvalidStateException (400)   - KMS key 状态不合法
//	KmsNotFoundException (400)       - KMS key 不存在
//
// https://cloud.tencent.com/document/product/436/127703
type VectorErrorResponse struct {
	Response  *http.Response       `json:"-"`
	Code      string               `json:"-"`                        // 错误码，来自 X-Cos-Error-Code 头部
	Message   string               `json:"message"`                  // 错误信息
	FieldList []VectorValidateField `json:"fieldList,omitempty"`      // 参数校验失败详情
	RequestID string               `json:"-"`                        // 请求 ID
}

// Error 实现 error 接口
func (r *VectorErrorResponse) Error() string {
	requestID := r.RequestID
	if requestID == "" && r.Response != nil {
		requestID = r.Response.Header.Get("X-Cos-Request-Id")
	}
	var reqURL string
	if r.Response != nil && r.Response.Request != nil {
		var err error
		reqURL, err = decodeURIComponent(r.Response.Request.URL.String())
		if err != nil {
			// 解码失败时使用原始URL作为降级策略
			reqURL = r.Response.Request.URL.String()
		}
	}
	statusCode := 0
	method := ""
	if r.Response != nil {
		statusCode = r.Response.StatusCode
		if r.Response.Request != nil {
			method = r.Response.Request.Method
		}
	}
	msg := fmt.Sprintf("%v %v: %d %v(Message: %v, RequestId: %v)",
		method, reqURL, statusCode, r.Code, r.Message, requestID)
	if len(r.FieldList) > 0 {
		for _, f := range r.FieldList {
			msg += fmt.Sprintf(", Field(%v: %v)", f.Path, f.Message)
		}
	}
	return msg
}

// checkVectorResponse 检查 Vector API 的 HTTP 响应是否为错误
// 与 COS 的 checkResponse 不同:
// 1. 错误码从 X-Cos-Error-Code 响应头获取
// 2. 错误体是 JSON 格式，包含 message 和可选的 fieldList
func checkVectorResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	vectorError := &VectorErrorResponse{Response: r}
	// 从响应头获取错误码
	vectorError.Code = r.Header.Get("X-Cos-Error-Code")
	vectorError.RequestID = r.Header.Get("X-Cos-Request-Id")

	// 读取并解析 JSON 错误体
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, vectorError)
	}

	return vectorError
}

// IsVectorError 判断错误是否为 Vector 服务返回的错误
// 返回 VectorErrorResponse 和判断结果
func IsVectorError(e error) (*VectorErrorResponse, bool) {
	if e == nil {
		return nil, false
	}
	err, ok := e.(*VectorErrorResponse)
	return err, ok
}
