package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetEncryption(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"encryption": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<ServerSideEncryptionConfiguration>
    <Rule>
        <ApplyServerSideEncryptionByDefault>
            <SSEAlgorithm>AES256</SSEAlgorithm>
        </ApplyServerSideEncryptionByDefault>
    </Rule>
</ServerSideEncryptionConfiguration>`)

	})

	res, _, err := client.Bucket.GetEncryption(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetEncryption returned error %v", err)
	}

	want := &BucketGetEncryptionResult{
		XMLName: xml.Name{Local: "ServerSideEncryptionConfiguration"},
		Rule: &BucketEncryptionConfiguration{
			SSEAlgorithm: "AES256",
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetEncryption returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutEncryption(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutEncryptionOptions{
		Rule: &BucketEncryptionConfiguration{
			SSEAlgorithm: "AES256",
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"encryption": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutEncryptionOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "ServerSideEncryptionConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutEncryption request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutEncryption(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutEncryption returned error: %v", err)
	}
}

func TestBucketService_DeleteEncryption(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"encryption": "",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteEncryption(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteEncryption returned error: %v", err)
	}

}
