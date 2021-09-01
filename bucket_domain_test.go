package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetDomain(t *testing.T) {
	setup()
	defer teardown()

	rt := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"domain": "",
		}
		testFormValues(t, r, vs)
		rt++
		if rt < 3 {
			w.WriteHeader(http.StatusGatewayTimeout)
		}

		fmt.Fprint(w, `<DomainConfiguration>
  	<DomainRule>
    	<Status>ENABLED</Status>
	    <Name>www.abc.com</Name>
		<Type>REST</Type>
		<ForcedReplacement>CNAME</ForcedReplacement>
	</DomainRule>
</DomainConfiguration>`)
	})

	res, _, err := client.Bucket.GetDomain(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetDomain returned error %v", err)
	}

	want := &BucketGetDomainResult{
		XMLName: xml.Name{Local: "DomainConfiguration"},
		Rules: []BucketDomainRule{
			{
				Status:            "ENABLED",
				Name:              "www.abc.com",
				Type:              "REST",
				ForcedReplacement: "CNAME",
			},
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetDomain returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutDomain(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutDomainOptions{
		XMLName: xml.Name{Local: "DomainConfiguration"},
		Rules: []BucketDomainRule{
			{
				Status:            "ENABLED",
				Name:              "www.abc.com",
				Type:              "REST",
				ForcedReplacement: "CNAME",
			},
		},
	}

	rt := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"domain": "",
		}
		testFormValues(t, r, vs)
		rt++
		if rt < 3 {
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		body := new(BucketPutDomainOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "DomainConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutDomain request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutDomain(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutDomain returned error: %v", err)
	}
}

func TestBucketService_DeleteDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"domain": "",
		}
		testFormValues(t, r, vs)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteDomain(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteDomain returned error: %v", err)
	}

}
