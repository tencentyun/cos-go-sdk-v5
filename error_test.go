package cos

import (
	"fmt"
	"net/http"
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
