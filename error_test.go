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
