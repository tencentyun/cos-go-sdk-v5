package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetAccelerate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"accelerate": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<AccelerateConfiguration>
    <Status>Enabled</Status>
    <Type>COS</Type>
</AccelerateConfiguration>`)
	})

	res, _, err := client.Bucket.GetAccelerate(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetAccelerate returned error %v", err)
	}

	want := &BucketGetAccelerateResult{
		XMLName: xml.Name{Local: "AccelerateConfiguration"},
		Status:  "Enabled",
		Type:    "COS",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetAccelerate returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutAccelerate(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutAccelerateOptions{
		XMLName: xml.Name{Local: "AccelerateConfiguration"},
		Status:  "Enabled",
		Type:    "COS",
	}

	rt := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"accelerate": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutAccelerateOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "AccelerateConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutAccelerate request\n body: %+v\n, want %+v\n", body, want)
		}
		rt++
		if rt < 3 {
			w.WriteHeader(http.StatusBadGateway)
		}
	})

	_, err := client.Bucket.PutAccelerate(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutAccelerate returned error: %v", err)
	}
}
