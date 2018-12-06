package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_Get(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketGetOptions{
		Prefix:  "test",
		MaxKeys: 2,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"prefix":   "test",
			"max-keys": "2",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<ListBucketResult>
	<Name>test-1253846586</Name>
	<Prefix>test</Prefix>
	<Marker/>
	<MaxKeys>2</MaxKeys>
	<IsTruncated>true</IsTruncated>
	<NextMarker>test/delete.txt</NextMarker>
	<Contents>
		<Key>test/</Key>
		<LastModified>2017-06-09T16:32:25.000Z</LastModified>
		<ETag>&quot;&quot;</ETag>
		<Size>0</Size>
		<Owner>
			<ID>1253846586</ID>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
	</Contents>
	<Contents>
		<Key>test/anonymous_get.go</Key>
		<LastModified>2017-06-17T15:09:26.000Z</LastModified>
		<ETag>&quot;5b7236085f08b3818bfa40b03c946dcc&quot;</ETag>
		<Size>8</Size>
		<Owner>
			<ID>1253846586</ID>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
	</Contents>
</ListBucketResult>`)
	})

	ref, _, err := client.Bucket.Get(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Get returned error: %v", err)
	}

	want := &BucketGetResult{
		XMLName:     xml.Name{Local: "ListBucketResult"},
		Name:        "test-1253846586",
		Prefix:      "test",
		MaxKeys:     2,
		IsTruncated: true,
		NextMarker:  "test/delete.txt",
		Contents: []Object{
			{
				Key:          "test/",
				LastModified: "2017-06-09T16:32:25.000Z",
				ETag:         "\"\"",
				Size:         0,
				Owner: &Owner{
					ID: "1253846586",
				},
				StorageClass: "STANDARD",
			},
			{
				Key:          "test/anonymous_get.go",
				LastModified: "2017-06-17T15:09:26.000Z",
				ETag:         "\"5b7236085f08b3818bfa40b03c946dcc\"",
				Size:         8,
				Owner: &Owner{
					ID: "1253846586",
				},
				StorageClass: "STANDARD",
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.Get returned %+v, want %+v", ref, want)
	}
}

func TestBucketService_Put(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutOptions{
		XCosACL: "public-read",
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		testHeader(t, r, "x-cos-acl", "public-read")
	})

	_, err := client.Bucket.Put(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Put returned error: %v", err)
	}

}

func TestBucketService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.Delete(context.Background())
	if err != nil {
		t.Fatalf("Bucket.Delete returned error: %v", err)
	}
}

func TestBucketService_Head(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodHead)
		w.WriteHeader(http.StatusOK)
	})

	_, err := client.Bucket.Head(context.Background())
	if err != nil {
		t.Fatalf("Bucket.Head returned error: %v", err)
	}
}
