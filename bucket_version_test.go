package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_PutVersioning(t *testing.T) {
	setup()
	defer teardown()
	opt := &BucketPutVersionOptions{
		Status: "Suspended",
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"versioning": "",
		}
		testFormValues(t, r, vs)

		body := &BucketPutVersionOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "VersioningConfiguration"}
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PutVersioning request\nbody: %+v\nwant %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutVersioning(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutVersioning failed, error: %v", err)
	}
}

func TestBucketService_GetVersioning(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"versioning": "",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<VersioningConfiguration>
    <Status>Suspended</Status>
</VersioningConfiguration>`)
	})
	res, _, err := client.Bucket.GetVersioning(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetVersioning failed, error: %v", err)
	}
	want := &BucketGetVersionResult{
		XMLName: xml.Name{Local: "VersioningConfiguration"},
		Status:  "Suspended",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetVersioning returned\n%+v, want\n%+v", res, want)
	}
}
