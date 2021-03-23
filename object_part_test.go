package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestObjectService_AbortMultipartUpload(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"
	uploadID := "xxxxaabcc"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"uploadId": uploadID,
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Object.AbortMultipartUpload(context.Background(),
		name, uploadID)
	if err != nil {
		t.Fatalf("Object.AbortMultipartUpload returned error: %v", err)
	}
}

func TestObjectService_InitiateMultipartUpload(t *testing.T) {
	setup()
	defer teardown()

	opt := &InitiateMultipartUploadOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")
		vs := values{
			"uploads": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<InitiateMultipartUploadResult>
	<Bucket>test-1253846586</Bucket>
	<Key>test/hello.txt</Key>
	<UploadId>149795166761115ef06e259b2fceb8ff34bf7dd840883d26a0f90243562dd398efa41718db</UploadId>
</InitiateMultipartUploadResult>`)
	})

	ref, _, err := client.Object.InitiateMultipartUpload(context.Background(),
		name, opt)
	if err != nil {
		t.Fatalf("Object.InitiateMultipartUpload returned error: %v", err)
	}

	want := &InitiateMultipartUploadResult{
		XMLName:  xml.Name{Local: "InitiateMultipartUploadResult"},
		Bucket:   "test-1253846586",
		Key:      "test/hello.txt",
		UploadID: "149795166761115ef06e259b2fceb8ff34bf7dd840883d26a0f90243562dd398efa41718db",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.InitiateMultipartUpload returned %+v, want %+v", ref, want)
	}
}

func TestObjectService_UploadPart(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectUploadPartOptions{}
	name := "test/hello.txt"
	uploadID := "xxxxx"
	partNumber := 1

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"uploadId":   uploadID,
			"partNumber": "1",
		}
		testFormValues(t, r, vs)

		b, _ := ioutil.ReadAll(r.Body)
		tb := crc64.MakeTable(crc64.ECMA)
		crc := crc64.Update(0, tb, b)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.UploadPart request body: %#v, want %#v", v, want)
		}
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
	})

	r := bytes.NewReader([]byte("hello"))
	_, err := client.Object.UploadPart(context.Background(),
		name, uploadID, partNumber, r, opt)
	if err != nil {
		t.Fatalf("Object.UploadPart returned error: %v", err)
	}

}

func TestObjectService_ListParts(t *testing.T) {
	setup()
	defer teardown()

	name := "test/hello.txt"
	uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodGet)
		vs := values{
			"uploadId": uploadID,
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<ListPartsResult>
	<Bucket>test-1253846586</Bucket>
	<Encoding-type/>
	<Key>test/hello.txt</Key>
	<UploadId>149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387</UploadId>
	<Owner>
		<ID>1253846586</ID>
		<DisplayName>1253846586</DisplayName>
	</Owner>
	<PartNumberMarker>0</PartNumberMarker>
	<Initiator>
		<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
		<DisplayName>100000760461</DisplayName>
	</Initiator>
	<Part>
		<PartNumber>1</PartNumber>
		<LastModified>2017-06-20T09:45:49.000Z</LastModified>
		<ETag>&quot;fae3dba15f4d9b2d76cbaed5de3a08e3&quot;</ETag>
		<Size>6291456</Size>
	</Part>
	<Part>
		<PartNumber>2</PartNumber>
		<LastModified>2017-06-20T09:45:50.000Z</LastModified>
		<ETag>&quot;c81982550f2f965118d486176d9541d4&quot;</ETag>
		<Size>6391456</Size>
	</Part>
	<StorageClass>Standard</StorageClass>
	<MaxParts>1000</MaxParts>
	<IsTruncated>false</IsTruncated>
</ListPartsResult>`)
	})

	ref, _, err := client.Object.ListParts(context.Background(),
		name, uploadID, nil)
	if err != nil {
		t.Fatalf("Object.ListParts returned error: %v", err)
	}

	want := &ObjectListPartsResult{
		XMLName:  xml.Name{Local: "ListPartsResult"},
		Bucket:   "test-1253846586",
		UploadID: uploadID,
		Key:      name,
		Owner: &Owner{
			ID:          "1253846586",
			DisplayName: "1253846586",
		},
		PartNumberMarker: "0",
		Initiator: &Initiator{
			ID:          "qcs::cam::uin/100000760461:uin/100000760461",
			DisplayName: "100000760461",
		},
		Parts: []Object{
			{
				PartNumber:   1,
				LastModified: "2017-06-20T09:45:49.000Z",
				ETag:         "\"fae3dba15f4d9b2d76cbaed5de3a08e3\"",
				Size:         6291456,
			},
			{
				PartNumber:   2,
				LastModified: "2017-06-20T09:45:50.000Z",
				ETag:         "\"c81982550f2f965118d486176d9541d4\"",
				Size:         6391456,
			},
		},
		StorageClass: "Standard",
		MaxParts:     "1000",
		IsTruncated:  false,
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.ListParts returned \n%#v, want \n%#v", ref, want)
	}
}

func TestObjectService_CompleteMultipartUpload(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"
	uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"

	opt := &CompleteMultipartUploadOptions{
		Parts: []Object{
			{
				PartNumber: 1,
				ETag:       "\"fae3dba15f4d9b2d76cbaed5de3a08e3\"",
			},
			{
				PartNumber: 2,
				ETag:       "\"c81982550f2f965118d486176d9541d4\"",
			},
		},
	}

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		v := new(CompleteMultipartUploadOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPost)
		vs := values{
			"uploadId": uploadID,
		}
		testFormValues(t, r, vs)

		want := opt
		want.XMLName = xml.Name{Local: "CompleteMultipartUpload"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.CompleteMultipartUpload request body: %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `<CompleteMultipartUploadResult>
	<Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
	<Bucket>test</Bucket>
	<Key>test/hello.txt</Key>
	<ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`)
	})

	ref, _, err := client.Object.CompleteMultipartUpload(context.Background(),
		name, uploadID, opt)
	if err != nil {
		t.Fatalf("Object.ListParts returned error: %v", err)
	}

	want := &CompleteMultipartUploadResult{
		XMLName:  xml.Name{Local: "CompleteMultipartUploadResult"},
		Bucket:   "test",
		Key:      name,
		ETag:     "\"594f98b11c6901c0f0683de1085a6d0e-4\"",
		Location: "test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.CompleteMultipartUpload returned \n%#v, want \n%#v", ref, want)
	}
}

func TestObjectService_CopyPart(t *testing.T) {
	setup()
	defer teardown()

	sourceUrl := "test-1253846586.cos.ap-guangzhou.myqcloud.com/test.source"
	opt := &ObjectCopyPartOptions{}
	name := "test/hello.txt"
	uploadID := "xxxxx"
	partNumber := 1

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"uploadId":   uploadID,
			"partNumber": "1",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<CopyPartResult>
   <ETag>&quot;ba82b57cfdfda8bd17ad4e5879ebb4fe&quot;</ETag>
   <LastModified>2017-09-04T04:45:45</LastModified>
</CopyPartResult>`)
	})

	r, _, err := client.Object.CopyPart(context.Background(),
		name, uploadID, partNumber, sourceUrl, opt)
	if err != nil {
		t.Fatalf("Object.CopyPart returned error: %v", err)
	}

	want := &CopyPartResult{
		XMLName:      xml.Name{Local: "CopyPartResult"},
		ETag:         `"ba82b57cfdfda8bd17ad4e5879ebb4fe"`,
		LastModified: "2017-09-04T04:45:45",
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Object.Copy returned %+v, want %+v", r, want)
	}
}
