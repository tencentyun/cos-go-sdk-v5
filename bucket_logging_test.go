package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_PutLogging(t *testing.T) {
	setup()
	defer teardown()
	opt := &BucketPutLoggingOptions{
		LoggingEnabled: &BucketLoggingEnabled{
			TargetBucket: "logs",
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"logging": "",
		}
		testFormValues(t, r, vs)

		body := &BucketPutLoggingOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "BucketLoggingStatus"}
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PutLogging request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutLogging(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutLogging failed, error: %v", err)
	}
}

func TestBucketService_GetLogging(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"logging": "",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<BucketLoggingStatus>
    <LoggingEnabled>
        <TargetBucket>logs</TargetBucket>
        <TargetPrefix>mylogs</TargetPrefix>
    </LoggingEnabled>
</BucketLoggingStatus>`)
	})
	res, _, err := client.Bucket.GetLogging(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetLogging failed, error: %v", err)
	}
	want := &BucketGetLoggingResult{
		XMLName: xml.Name{Local: "BucketLoggingStatus"},
		LoggingEnabled: &BucketLoggingEnabled{
			TargetBucket: "logs",
			TargetPrefix: "mylogs",
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetLogging returned\n%+v, want\n%+v", res, want)
	}
}
