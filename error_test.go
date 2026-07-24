package cos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// func Test_checkResponse_error(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	mux.HandleFunc("/test_409", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusConflict)
// 		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
// <Error>
// 	<Code>BucketAlreadyExists</Code>
// 	<Message>The requested bucket name is not available.</Message>
// 	<Resource>testdelete-1253846586.cos.ap-guangzhou.myqcloud.com</Resource>
// 	<RequestId>NTk0NTRjZjZfNTViMjM1XzlkMV9hZTZh</RequestId>
// 	<TraceId>OGVmYzZiMmQzYjA2OWNhODk0NTRkMTBiOWVmMDAxODc0OWRkZjk0ZDM1NmI1M2E2MTRlY2MzZDhmNmI5MWI1OTBjYzE2MjAxN2M1MzJiOTdkZjMxMDVlYTZjN2FiMmI0NTk3NWFiNjAyMzdlM2RlMmVmOGNiNWIxYjYwNDFhYmQ=</TraceId>
// </Error>`)
// 	})

// 	req, _ := http.NewRequest("GET", client.BaseURL.ServiceURL.String()+"/test_409", nil)
// 	resp, _ := client.client.Do(req)
// 	err := checkResponse(resp)

// 	if e, ok := err.(*ErrorResponse); ok {
// 		if e.Error() == "" {
// 			t.Errorf("Expected e.Error() not empty, got %+v", e.Error())
// 		}
// 		if e.Code != "BucketAlreadyExists" {
// 			t.Errorf("Expected BucketAlreadyExists error, got %+v", e.Code)
// 		}
// 	} else {
// 		t.Errorf("Expected ErrorResponse error, got %+v", err)
// 	}
// }

func Test_checkResponse_no_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test_200", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `test`)
	})

	req, _ := http.NewRequest("GET", client.BaseURL.ServiceURL.String()+"/test_200", nil)
	resp, _ := client.client.Do(req)
	err := checkResponse(resp)

	if err != nil {
		t.Errorf("Expected error == nil, got %+v", err)
	}
}

func Test_checkResponse_with_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test_409", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<Error>
	<Code>BucketAlreadyExists</Code>
	<Message>The requested bucket name is not available.</Message>
	<Resource>testdelete-1253846586.cos.ap-guangzhou.myqcloud.com</Resource>
	<RequestId>NTk0NTRjZjZfNTViMjM1XzlkMV9hZTZh</RequestId>
	<TraceId>OGVmYzZiMmQzYjA2OWNhODk0NTRkMTBiOWVmMDAxODc0OWRkZjk0ZDM1NmI1M2E2MTRlY2MzZDhmNmI5MWI1OTBjYzE2MjAxN2M1MzJiOTdkZjMxMDVlYTZjN2FiMmI0NTk3NWFiNjAyMzdlM2RlMmVmOGNiNWIxYjYwNDFhYmQ=</TraceId>
</Error>`)
	})

	req, _ := http.NewRequest("GET", client.BaseURL.ServiceURL.String()+"/test_409", nil)
	resp, _ := client.client.Do(req)
	err := checkResponse(resp)

	if e, ok := err.(*ErrorResponse); ok {
		if e.Error() == "" {
			t.Errorf("Expected e.Error() not empty, got %+v", e.Error())
		}
		if e.Code != "BucketAlreadyExists" {
			t.Errorf("Expected BucketAlreadyExists error, got %+v", e.Code)
		}
	} else {
		t.Errorf("Expected ErrorResponse error, got %+v", err)
	}

}

func Test_IsNotFoundError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test_404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<Error>
    <Code>NoSuchKey</Code>
    <Message>The specified key does not exist.</Message>
    <Resource>examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com/test_404</Resource>
    <RequestId>NjA3OGY4NGFfNjJkMmMwYl8***</RequestId>
    <TraceId>OGVmYzZiMmQzYjA2OWNh***</TraceId>
</Error>`)
	})

	req, _ := http.NewRequest("GET", client.BaseURL.ServiceURL.String()+"/test_404", nil)
	resp, _ := client.client.Do(req)
	err := checkResponse(resp)

	e, ok := IsCOSError(err)
	if !ok {
		t.Errorf("IsCOSError Return Failed")
	}
	ok = IsNotFoundError(e)
	if !ok {
		t.Errorf("IsNotFoundError Return Failed")
	}
	if e.Code != "NoSuchKey" {
		t.Errorf("Expected NoSuchKey error, got %+v", e.Code)
	}
	_, ok = IsCOSError(nil)
	if ok {
		t.Errorf("IsCOSError Return Failed")
	}
	_, ok = IsCOSError(errors.New("test error"))
	if ok {
		t.Errorf("IsNotFoundError Return Failed")
	}
	ok = IsNotFoundError(nil)
	if ok {
		t.Errorf("IsNotFoundError Return Failed")
	}
	ok = IsNotFoundError(errors.New("test error"))
	if ok {
		t.Errorf("IsNotFoundError Return Failed")
	}
}

func Test_CheckReponse(t *testing.T) {
	setup()
	defer teardown()
	resp := &http.Response{
		StatusCode: 404,
		Header:     http.Header{},
		Body: ioutil.NopCloser(strings.NewReader(`{
			"code": 404,
			"message": "404 NotFound",
			"request_id": "requestid"}`)),
	}
	resp.Header.Add("Content-Type", "application/json")
	err := checkResponse(resp)
	e, _ := err.(*ErrorResponse)
	if e.Code != "404" || e.Message != "404 NotFound" || e.RequestID != "requestid" {
		t.Errorf("checkResponse failed: %v", e)
	}
}

func Test_RetryError(t *testing.T) {
	var errs RetryError
	errs.Add(errors.New("err1"))
	errs.Add(errors.New("err2"))
	errs.Add(errors.New("err3"))
	if errs.Error() != "err1; err2; err3" {
		t.Errorf("RetryError return err: %v", errs.Error())
	}
}

func Test_RetryError_Is(t *testing.T) {
	sentinel := errors.New("sentinel")
	other := errors.New("other")

	// 空 RetryError：Is 返回 false
	var empty RetryError
	if errors.Is(&empty, sentinel) {
		t.Error("empty RetryError.Is should return false")
	}

	// 最后一个错误匹配 → true
	r1 := &RetryError{}
	r1.Add(other)
	r1.Add(sentinel)
	if !errors.Is(r1, sentinel) {
		t.Error("RetryError.Is should match last error")
	}

	// 最后一个错误不匹配，即使前面有匹配的 → false
	r2 := &RetryError{}
	r2.Add(sentinel)
	r2.Add(other)
	if errors.Is(r2, sentinel) {
		t.Error("RetryError.Is should only check last error, not earlier ones")
	}

	// 最后一个错误通过 %w 包装了 sentinel → 能穿透
	r3 := &RetryError{}
	r3.Add(other)
	r3.Add(fmt.Errorf("wrapped: %w", sentinel))
	if !errors.Is(r3, sentinel) {
		t.Error("RetryError.Is should unwrap last error's chain")
	}
}

func Test_RetryError_As(t *testing.T) {
	makeResp := func(code string) *ErrorResponse {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
		return &ErrorResponse{
			Code:     code,
			Response: &http.Response{StatusCode: 404, Request: req, Header: http.Header{}},
		}
	}

	cosErr1 := makeResp("NoSuchKey")
	cosErr2 := makeResp("AccessDenied")
	plainErr := errors.New("plain error")

	// 空 RetryError：As 返回 false
	var empty RetryError
	var target *ErrorResponse
	if errors.As(&empty, &target) {
		t.Error("empty RetryError.As should return false")
	}

	// 最后一个错误匹配 *ErrorResponse → true，且拿到的是最后一个
	r1 := &RetryError{}
	r1.Add(cosErr1)
	r1.Add(cosErr2)
	target = nil
	if !errors.As(r1, &target) {
		t.Fatal("RetryError.As should match last error")
	}
	if target.Code != "AccessDenied" {
		t.Errorf("expected last error Code=AccessDenied, got %v", target.Code)
	}

	// 最后一个错误不是 *ErrorResponse，即使前面有 → false
	r2 := &RetryError{}
	r2.Add(cosErr1)
	r2.Add(plainErr)
	target = nil
	if errors.As(r2, &target) {
		t.Error("RetryError.As should only check last error")
	}

	// 最后一个错误通过 %w 包装了 *ErrorResponse → 能穿透
	r3 := &RetryError{}
	r3.Add(cosErr1)
	r3.Add(fmt.Errorf("wrapped: %w", cosErr2))
	target = nil
	if !errors.As(r3, &target) {
		t.Fatal("RetryError.As should unwrap last error's chain")
	}
	if target.Code != "AccessDenied" {
		t.Errorf("expected wrapped last error Code=AccessDenied, got %v", target.Code)
	}
}

// makeCOSErr 构造一个带 HTTP 响应的 *ErrorResponse，方便复用。
func makeCOSErr(code string, statusCode int) *ErrorResponse {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/obj", nil)
	return &ErrorResponse{
		Code:    code,
		Message: code,
		Response: &http.Response{
			StatusCode: statusCode,
			Request:    req,
			Header:     http.Header{},
		},
	}
}

// makeVectorErr 构造一个带 HTTP 响应的 *VectorErrorResponse，方便复用。
func makeVectorErr(code string, statusCode int) *VectorErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/vector", nil)
	return &VectorErrorResponse{
		Code:    code,
		Message: code,
		Response: &http.Response{
			StatusCode: statusCode,
			Request:    req,
			Header:     http.Header{},
		},
	}
}

func Test_IsCOSError(t *testing.T) {
	cosErr := makeCOSErr("NoSuchKey", 404)

	// nil → false
	if _, ok := IsCOSError(nil); ok {
		t.Error("IsCOSError(nil) should return false")
	}

	// 非 COS 错误 → false
	if _, ok := IsCOSError(errors.New("plain")); ok {
		t.Error("IsCOSError(plain) should return false")
	}

	// 直接 *ErrorResponse → true，Code 正确
	e, ok := IsCOSError(cosErr)
	if !ok {
		t.Fatal("IsCOSError(*ErrorResponse) should return true")
	}
	if e.Code != "NoSuchKey" {
		t.Errorf("expected Code=NoSuchKey, got %v", e.Code)
	}

	// fmt.Errorf %w 包装后 → true（errors.As 穿透）
	wrapped := fmt.Errorf("outer: %w", cosErr)
	e, ok = IsCOSError(wrapped)
	if !ok {
		t.Fatal("IsCOSError(wrapped *ErrorResponse) should return true")
	}
	if e.Code != "NoSuchKey" {
		t.Errorf("expected Code=NoSuchKey through wrap, got %v", e.Code)
	}

	// RetryError 最后一个是 *ErrorResponse → true
	r := &RetryError{}
	r.Add(errors.New("attempt1"))
	r.Add(cosErr)
	e, ok = IsCOSError(r)
	if !ok {
		t.Fatal("IsCOSError(RetryError last=*ErrorResponse) should return true")
	}
	if e.Code != "NoSuchKey" {
		t.Errorf("expected Code=NoSuchKey from RetryError, got %v", e.Code)
	}

	// RetryError 最后一个不是 *ErrorResponse → false
	r2 := &RetryError{}
	r2.Add(cosErr)
	r2.Add(errors.New("last is plain"))
	if _, ok := IsCOSError(r2); ok {
		t.Error("IsCOSError(RetryError last=plain) should return false")
	}
}

func Test_IsNotFoundError_Extended(t *testing.T) {
	cosErr404 := makeCOSErr("NoSuchKey", 404)
	cosErr403 := makeCOSErr("AccessDenied", 403)

	// nil → false
	if IsNotFoundError(nil) {
		t.Error("IsNotFoundError(nil) should return false")
	}

	// 非 COS 错误 → false
	if IsNotFoundError(errors.New("plain")) {
		t.Error("IsNotFoundError(plain) should return false")
	}

	// 404 *ErrorResponse → true
	if !IsNotFoundError(cosErr404) {
		t.Error("IsNotFoundError(404 cosErr) should return true")
	}

	// 非 404 *ErrorResponse → false
	if IsNotFoundError(cosErr403) {
		t.Error("IsNotFoundError(403 cosErr) should return false")
	}

	// fmt.Errorf %w 包装 404 → true（穿透）
	if !IsNotFoundError(fmt.Errorf("wrap: %w", cosErr404)) {
		t.Error("IsNotFoundError(wrapped 404) should return true")
	}

	// RetryError 最后一个是 404 → true
	r := &RetryError{}
	r.Add(cosErr403)
	r.Add(cosErr404)
	if !IsNotFoundError(r) {
		t.Error("IsNotFoundError(RetryError last=404) should return true")
	}

	// RetryError 最后一个是 403 → false（即使前面有 404）
	r2 := &RetryError{}
	r2.Add(cosErr404)
	r2.Add(cosErr403)
	if IsNotFoundError(r2) {
		t.Error("IsNotFoundError(RetryError last=403) should return false")
	}
}

func Test_IsVectorError(t *testing.T) {
	vecErr := makeVectorErr("NotFoundException", 404)

	// nil → false
	if _, ok := IsVectorError(nil); ok {
		t.Error("IsVectorError(nil) should return false")
	}

	// 非 Vector 错误 → false
	if _, ok := IsVectorError(errors.New("plain")); ok {
		t.Error("IsVectorError(plain) should return false")
	}

	// *ErrorResponse（COS 错误）→ false
	if _, ok := IsVectorError(makeCOSErr("NoSuchKey", 404)); ok {
		t.Error("IsVectorError(*ErrorResponse) should return false")
	}

	// 直接 *VectorErrorResponse → true，Code 正确
	e, ok := IsVectorError(vecErr)
	if !ok {
		t.Fatal("IsVectorError(*VectorErrorResponse) should return true")
	}
	if e.Code != "NotFoundException" {
		t.Errorf("expected Code=NotFoundException, got %v", e.Code)
	}

	// fmt.Errorf %w 包装后 → true（errors.As 穿透）
	wrapped := fmt.Errorf("outer: %w", vecErr)
	e, ok = IsVectorError(wrapped)
	if !ok {
		t.Fatal("IsVectorError(wrapped *VectorErrorResponse) should return true")
	}
	if e.Code != "NotFoundException" {
		t.Errorf("expected Code=NotFoundException through wrap, got %v", e.Code)
	}

	// RetryError 最后一个是 *VectorErrorResponse → true
	r := &RetryError{}
	r.Add(errors.New("attempt1"))
	r.Add(vecErr)
	e, ok = IsVectorError(r)
	if !ok {
		t.Fatal("IsVectorError(RetryError last=*VectorErrorResponse) should return true")
	}
	if e.Code != "NotFoundException" {
		t.Errorf("expected Code=NotFoundException from RetryError, got %v", e.Code)
	}

	// RetryError 最后一个不是 *VectorErrorResponse → false
	r2 := &RetryError{}
	r2.Add(vecErr)
	r2.Add(errors.New("last is plain"))
	if _, ok := IsVectorError(r2); ok {
		t.Error("IsVectorError(RetryError last=plain) should return false")
	}
}

func Test_ErrorResponse(t *testing.T) {
	setup()
	defer teardown()
	request, _ := http.NewRequest(http.MethodGet, client.BaseURL.BucketURL.String(), nil)
	request.Header.Add("X-Cos-Request-Id", "requestid")
	request.Header.Add("X-Cos-Trace-Id", "traceid")
	resp := &http.Response{
		StatusCode: 404,
		Header:     http.Header{},
		Request:    request,
		Body: ioutil.NopCloser(strings.NewReader(`{
			"code": 404,
			"message": "404 NotFound",
			"request_id": "requestid"}`)),
	}
	resp.Header.Add("Content-Type", "application/json")
	err := &ErrorResponse{
		Response: resp,
	}
	except := fmt.Sprintf("%v %v: %d %v(Message: %v, RequestId: %v, TraceId: %v)", resp.Request.Method, client.BaseURL.BucketURL.String(), resp.StatusCode, err.Code, err.Message, err.RequestID, err.TraceID)
	if err.Error() != except {
		t.Errorf("error message is invalid, return: %v, except: %v", err.Error(), except)
	}
}
