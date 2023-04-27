package cos

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
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

func TestObjectService_ListUploads(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"uploads":     "",
			"max-uploads": "1000",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<ListMultipartUploadsResult>
    <Bucket>examplebucket-1250000000</Bucket>
    <Encoding-Type/>
    <KeyMarker/>
    <UploadIdMarker/>
    <MaxUploads>1000</MaxUploads>
    <Prefix/>
    <Delimiter>/</Delimiter>
    <IsTruncated>false</IsTruncated>
    <Upload>
        <Key>Object</Key>
        <UploadId>1484726657932bcb5b17f7a98a8cad9fc36a340ff204c79bd2f51e7dddf0b6d1da6220520c</UploadId>
        <Initiator>
            <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            <DisplayName>100000000001</DisplayName>
        </Initiator>
        <Owner>
            <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            <DisplayName>100000000001</DisplayName>
        </Owner>
        <StorageClass>Standard</StorageClass>
        <Initiated>Wed Jan 18 16:04:17 2017</Initiated>
    </Upload>
    <Upload>
        <Key>Object</Key>
        <UploadId>1484727158f2b8034e5407d18cbf28e84f754b791ecab607d25a2e52de9fee641e5f60707c</UploadId>
        <Initiator>
            <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            <DisplayName>100000000001</DisplayName>
        </Initiator>
        <Owner>
            <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            <DisplayName>100000000001</DisplayName>
        </Owner>
        <StorageClass>Standard</StorageClass>
        <Initiated>Wed Jan 18 16:12:38 2017</Initiated>
    </Upload>
    <Upload>
        <Key>exampleobject</Key>
        <UploadId>1484727270323ddb949d528c629235314a9ead80f0ba5d993a3d76b460e6a9cceb9633b08e</UploadId>
        <Initiator>
            <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            <DisplayName>100000000001</DisplayName>
        </Initiator>
        <Owner>
            <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            <DisplayName>100000000001</DisplayName>
        </Owner>
        <StorageClass>Standard</StorageClass>
        <Initiated>Wed Jan 18 16:14:30 2017</Initiated>
    </Upload>
</ListMultipartUploadsResult>`)
	})

	opt := &ObjectListUploadsOptions{
		MaxUploads: 1000,
	}
	ref, _, err := client.Object.ListUploads(context.Background(), opt)
	if err != nil {
		t.Fatalf("Object.ListParts returned error: %v", err)
	}

	want := &ObjectListUploadsResult{
		XMLName:     xml.Name{Local: "ListMultipartUploadsResult"},
		Bucket:      "examplebucket-1250000000",
		MaxUploads:  "1000",
		IsTruncated: false,
		Delimiter:   "/",
		Upload: []ListUploadsResultUpload{
			{
				Key:      "Object",
				UploadID: "1484726657932bcb5b17f7a98a8cad9fc36a340ff204c79bd2f51e7dddf0b6d1da6220520c",
				Initiator: &Initiator{
					ID:          "qcs::cam::uin/100000000001:uin/100000000001",
					DisplayName: "100000000001",
				},
				Owner: &Owner{
					ID:          "qcs::cam::uin/100000000001:uin/100000000001",
					DisplayName: "100000000001",
				},
				StorageClass: "Standard",
				Initiated:    "Wed Jan 18 16:04:17 2017",
			},
			{
				Key:      "Object",
				UploadID: "1484727158f2b8034e5407d18cbf28e84f754b791ecab607d25a2e52de9fee641e5f60707c",
				Initiator: &Initiator{
					ID:          "qcs::cam::uin/100000000001:uin/100000000001",
					DisplayName: "100000000001",
				},
				Owner: &Owner{
					ID:          "qcs::cam::uin/100000000001:uin/100000000001",
					DisplayName: "100000000001",
				},
				StorageClass: "Standard",
				Initiated:    "Wed Jan 18 16:12:38 2017",
			},
			{
				Key:      "exampleobject",
				UploadID: "1484727270323ddb949d528c629235314a9ead80f0ba5d993a3d76b460e6a9cceb9633b08e",
				Initiator: &Initiator{
					ID:          "qcs::cam::uin/100000000001:uin/100000000001",
					DisplayName: "100000000001",
				},
				Owner: &Owner{
					ID:          "qcs::cam::uin/100000000001:uin/100000000001",
					DisplayName: "100000000001",
				},
				StorageClass: "Standard",
				Initiated:    "Wed Jan 18 16:14:30 2017",
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.ListUploads returned \n%+v, want \n%+v", ref, want)
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
		fmt.Fprint(w, `<Error>
		            <Code>string</Code>
					            <Message>string</Message>
								            <Resource>string</Resource>
											            <RequestId>string</RequestId>
														            <TraceId>string</TraceId>

</Error>`)
	})

	_, _, err := client.Object.CopyPart(context.Background(),
		name, uploadID, partNumber, sourceUrl, opt)
	if err == nil {
		t.Fatalf("Object.CopyPart returned error is nil")
	}
}

func TestObjectService_MultiCopy(t *testing.T) {
	setup()
	defer teardown()

	totalBytes := 1024 * 1024 * 35
	b := make([]byte, totalBytes)
	rand.Read(b)

	opt := &MultiCopyOptions{
		PartSize:       1,
		ThreadPoolSize: 3,
		useMulti:       true,
		OptCopy: &ObjectCopyOptions{
			&ObjectCopyHeaderOptions{},
			&ACLHeaderOptions{},
		},
	}
	uploadid := "test_uploadid"
	optcom := &CompleteMultipartUploadOptions{}
	for i := 1; i <= 35; i += 1 {
		optcom.Parts = append(optcom.Parts, Object{
			PartNumber: i,
			ETag:       hex.EncodeToString(calMD5Digest(b[(i-1)*1024*1024 : i*1024*1024])),
		})
	}

	mux.HandleFunc("/test.src.copy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "HEAD")
		w.Header().Add("Content-Length", strconv.FormatInt(int64(totalBytes), 10))
	})

	mux.HandleFunc("/test.copy", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// Upload Part
		if r.Method == http.MethodPut {
			str := r.Form.Get("partNumber")
			partNumber, _ := strconv.ParseInt(str, 10, 64)
			ranger, err := GetRange(r.Header.Get("x-cos-copy-source-range"))
			if err != nil {
				t.Errorf("Object.MultiCopy GetRange failed: %v", err)
			}
			if ranger.Start != (partNumber-1)*1024*1024 || ranger.End != partNumber*1024*1024-1 {
				t.Errorf("Object.MultiCopy range error, range: %+v, partNumber: %+v\n", r.Header.Get("x-cos-copy-source-range"), partNumber)
			}
			if r.Form.Get("uploadId") != uploadid {
				t.Errorf("Object.MultiCopy PartCopy returned %+v, want %+v", r.Form.Get("uploadId"), uploadid)
			}

			fmt.Fprintf(w, `<CopyPartResult>
    <ETag>%v</ETag>
    <LastModified></LastModified>
</CopyPartResult>`, optcom.Parts[partNumber-1].ETag)
			return
		}

		testMethod(t, r, http.MethodPost)
		// Complete MultiPart
		if r.Form.Get("uploadId") == uploadid {
			v := &CompleteMultipartUploadOptions{}
			xml.NewDecoder(r.Body).Decode(v)
			vs := values{
				"uploadId": uploadid,
			}
			testFormValues(t, r, vs)
			want := optcom
			want.XMLName = xml.Name{Local: "CompleteMultipartUpload"}
			if !reflect.DeepEqual(v, want) {
				t.Errorf("Object.MultiCopy Complete request body: %+v, want %+v", v, want)
			}
			fmt.Fprint(w, `<CompleteMultipartUploadResult>
	<Location></Location>
	<Bucket></Bucket>
	<Key></Key>
	<ETag>etag</ETag>
</CompleteMultipartUploadResult>`)
			return
		}

		// Init MultiPart
		fmt.Fprintf(w, `<InitiateMultipartUploadResult>
	<Bucket></Bucket>
	<Key></Key>
	<UploadId>%v</UploadId>
</InitiateMultipartUploadResult>`, uploadid)

	})

	dest := "test.copy"
	soruceURL := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, "test.src.copy")
	_, _, err := client.Object.MultiCopy(context.Background(), dest, soruceURL, opt)

	if err != nil {
		t.Errorf("Object.MultiCopy failed %v", err)
	}

}
