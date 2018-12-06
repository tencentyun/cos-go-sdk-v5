package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetCORS(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"cors": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<CORSConfiguration>
	<CORSRule>
		<AllowedOrigin>http://www.qq.com</AllowedOrigin>
		<AllowedMethod>PUT</AllowedMethod>
		<AllowedMethod>GET</AllowedMethod>
		<AllowedHeader>x-cos-meta-test</AllowedHeader>
		<AllowedHeader>x-cos-xx</AllowedHeader>
		<ExposeHeader>x-cos-meta-test1</ExposeHeader>
		<MaxAgeSeconds>500</MaxAgeSeconds>
	</CORSRule>
	<CORSRule>
		<ID>1234</ID>
		<AllowedOrigin>http://www.baidu.com</AllowedOrigin>
		<AllowedOrigin>twitter.com</AllowedOrigin>
		<AllowedMethod>PUT</AllowedMethod>
		<AllowedMethod>GET</AllowedMethod>
		<MaxAgeSeconds>500</MaxAgeSeconds>
	</CORSRule>
</CORSConfiguration>`)
	})

	ref, _, err := client.Bucket.GetCORS(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetCORS returned error: %v", err)
	}

	want := &BucketGetCORSResult{
		XMLName: xml.Name{Local: "CORSConfiguration"},
		Rules: []BucketCORSRule{
			{
				AllowedOrigins: []string{"http://www.qq.com"},
				AllowedMethods: []string{"PUT", "GET"},
				AllowedHeaders: []string{"x-cos-meta-test", "x-cos-xx"},
				MaxAgeSeconds:  500,
				ExposeHeaders:  []string{"x-cos-meta-test1"},
			},
			{
				ID:             "1234",
				AllowedOrigins: []string{"http://www.baidu.com", "twitter.com"},
				AllowedMethods: []string{"PUT", "GET"},
				MaxAgeSeconds:  500,
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.GetLifecycle returned %+v, want %+v", ref, want)
	}
}

func TestBucketService_PutCORS(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutCORSOptions{
		Rules: []BucketCORSRule{
			{
				AllowedOrigins: []string{"http://www.qq.com"},
				AllowedMethods: []string{"PUT", "GET"},
				AllowedHeaders: []string{"x-cos-meta-test", "x-cos-xx"},
				MaxAgeSeconds:  500,
				ExposeHeaders:  []string{"x-cos-meta-test1"},
			},
			{
				ID:             "1234",
				AllowedOrigins: []string{"http://www.baidu.com", "twitter.com"},
				AllowedMethods: []string{"PUT", "GET"},
				MaxAgeSeconds:  500,
			},
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutCORSOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPut)
		vs := values{
			"cors": "",
		}
		testFormValues(t, r, vs)

		want := opt
		want.XMLName = xml.Name{Local: "CORSConfiguration"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Bucket.PutCORS request body: %+v, want %+v", v, want)
		}

	})

	_, err := client.Bucket.PutCORS(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutCORS returned error: %v", err)
	}

}

func TestBucketService_DeleteCORS(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"cors": "",
		}
		testFormValues(t, r, vs)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteCORS(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteCORS returned error: %v", err)
	}

}
