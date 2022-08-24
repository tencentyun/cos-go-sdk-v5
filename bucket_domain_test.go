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

func TestBucketService_GetDomainCertificate(t *testing.T) {
	setup()
	defer teardown()

	rt := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"domaincertificate": "",
			"domainname":        "www.qq.com",
		}
		testFormValues(t, r, vs)
		rt++
		if rt < 3 {
			w.WriteHeader(http.StatusGatewayTimeout)
		}

		fmt.Fprint(w, `<DomainCertificate>
    	<Status>ENABLED</Status>
</DomainCertificate>`)
	})

	opt := &BucketGetDomainCertificateOptions{
		DomainName: "www.qq.com",
	}
	res, _, err := client.Bucket.GetDomainCertificate(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.GetDomainCertificate returned error %v", err)
	}

	want := &BucketGetDomainCertificateResult{
		XMLName: xml.Name{Local: "DomainCertificate"},
		Status:  "ENABLED",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetDomainCertificate returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutDomainCertificate(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutDomainCertificateOptions{
		XMLName: xml.Name{Local: "DomainCertificate"},
		CertificateInfo: &BucketDomainCertificateInfo{
			CertType: "CustomCert",
			CustomCert: &BucketDomainCustomCert{
				Cert:       "====certificate====",
				PrivateKey: "====PrivateKey====",
			},
		},
		DomainList: []string{"www.qq.com"},
	}

	rt := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"domaincertificate": "",
		}
		testFormValues(t, r, vs)
		rt++
		if rt < 3 {
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		body := new(BucketPutDomainCertificateOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "DomainCertificate"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutDomainCertificate request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutDomainCertificate(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutDomainCertificate returned error: %v", err)
	}
}

func TestBucketService_DeleteDomainCertificate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"domaincertificate": "",
			"domainname":        "www.qq.com",
		}
		testFormValues(t, r, vs)
		w.WriteHeader(http.StatusNoContent)
	})

	opt := &BucketDeleteDomainCertificateOptions{
		DomainName: "www.qq.com",
	}
	_, err := client.Bucket.DeleteDomainCertificate(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.DeleteDomainCertificate returned error: %v", err)
	}

}
