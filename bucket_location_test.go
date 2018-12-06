package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetLocation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"location": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<LocationConstraint>ap-guangzhou</LocationConstraint>`)
	})

	ref, _, err := client.Bucket.GetLocation(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetLocation returned error: %v", err)
	}

	want := &BucketGetLocationResult{
		XMLName:  xml.Name{Local: "LocationConstraint"},
		Location: "ap-guangzhou",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.GetLocation returned %+v, want %+v", ref, want)
	}
}
