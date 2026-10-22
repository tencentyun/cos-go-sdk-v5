package cos

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
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

	_, err := client.Object.AbortMultipartUpload(context.Background(), name, uploadID)
	if err != nil {
		t.Fatalf("Object.AbortMultipartUpload returned error: %v", err)
	}
	_, err = client.Object.AbortMultipartUpload(context.Background(), "/", uploadID)
	if err == nil || err.Error() != "empty object name" {
		t.Fatalf("Object.AbortMultipartUpload expect error: %v", err)
	}
	_, err = client.Object.AbortMultipartUpload(context.Background(), "//", uploadID)
	if err != ObjectKeySimplifyCheckErr {
		t.Fatalf("Object.AbortMultipartUpload expect error: %v", err)
	}
	opt := &AbortMultipartUploadOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "value")
	_, err = client.Object.AbortMultipartUpload(context.Background(), name, uploadID, opt)
	if err != nil {
		t.Fatalf("Object.AbortMultipartUpload returned error: %v", err)
	}
}

func TestObjectService_InitiateMultipartUpload_XMLSyntaxError(t *testing.T) {
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
	eChar, _ := DecodeURIComponent("%07") //特殊字符
	eName := fmt.Sprintf("test/hello_%v.txt", eChar)
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
		fmt.Fprintf(w, `<InitiateMultipartUploadResult>
	<Bucket>test-1253846586</Bucket>
	<Key>%v</Key>
	<UploadId>149795166761115ef06e259b2fceb8ff34bf7dd840883d26a0f90243562dd398efa41718db</UploadId>
</InitiateMultipartUploadResult>`, eName)
	})

	ref, _, err := client.Object.InitiateMultipartUpload(context.Background(),
		name, opt)
	if err != nil {
		t.Fatalf("Object.InitiateMultipartUpload returned error: %v", err)
	}

	want := &InitiateMultipartUploadResult{
		XMLName:  xml.Name{Local: "InitiateMultipartUploadResult"},
		Bucket:   "test-1253846586",
		Key:      eName,
		UploadID: "149795166761115ef06e259b2fceb8ff34bf7dd840883d26a0f90243562dd398efa41718db",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.InitiateMultipartUpload returned %+v, want %+v", ref, want)
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

func TestObjectService_UploadPartWithEmptyData(t *testing.T) {
	setup()
	defer teardown()

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
		if r.Header.Get("Content-Length") != "0" {
			t.Fatalf("UploadPart Content-Length is not 0")
		}
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			t.Fatalf("data is not empty")
		}
	})

	_, err := client.Object.UploadPart(context.Background(),
		name, uploadID, partNumber, nil, nil)
	if err != nil {
		t.Fatalf("Object.UploadPart returned error: %v", err)
	}

	_, err = client.Object.UploadPart(context.Background(),
		name, uploadID, partNumber, http.NoBody, nil)
	if err != nil {
		t.Fatalf("Object.UploadPart returned error: %v", err)
	}

	_, err = client.Object.UploadPart(context.Background(),
		name, uploadID, partNumber, http.NoBody, &ObjectUploadPartOptions{
			ContentLength: 1,
		})
	if err == nil || err.Error() != "ContentLength must be 0 when reader is nil or http.NoBody." {
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

func TestObjectService_CompleteMultipartUpload_XMLSyntaxError(t *testing.T) {
	setup()
	defer teardown()

	eChar, _ := DecodeURIComponent("%07") //特殊字符
	eName := fmt.Sprintf("test/hello_%v.txt", eChar)
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
		fmt.Fprintf(w, `<CompleteMultipartUploadResult>
	<Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
	<Bucket>test-1253846586</Bucket>
	<Key>%v</Key>
	<ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`, eName)
	})

	ref, _, err := client.Object.CompleteMultipartUpload(context.Background(),
		name, uploadID, opt)
	if err != nil {
		t.Fatalf("Object.CompleteMultipartUpload returned error: %v", err)
	}

	want := &CompleteMultipartUploadResult{
		XMLName:  xml.Name{Local: "CompleteMultipartUploadResult"},
		Bucket:   "test-1253846586",
		Key:      eName,
		ETag:     "\"594f98b11c6901c0f0683de1085a6d0e-4\"",
		Location: "test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.CompleteMultipartUpload returned \n%#v, want \n%#v", ref, want)
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

// TestObjectService_CompleteMultipartUpload_EmptyETag_RetrySuccess 模拟服务端 chunk 回包，
// 前两次返回空 ETag（空 body），第三次返回正常结果，验证重试后成功。
func TestObjectService_CompleteMultipartUpload_EmptyETag_RetrySuccess(t *testing.T) {
	setup()
	defer teardown()

	client.Conf.RetryOpt.Count = 3

	name := "test/hello.txt"
	uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"
	callCount := 0

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testFormValues(t, r, values{"uploadId": uploadID})
		callCount++
		if callCount < 3 {
			// 前两次：返回 200 但 body 为空（ETag 空）
			fmt.Fprint(w, `<CompleteMultipartUploadResult></CompleteMultipartUploadResult>`)
			return
		}
		// 第三次：返回完整正常结果
		fmt.Fprint(w, `<CompleteMultipartUploadResult>
	<Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
	<Bucket>test</Bucket>
	<Key>test/hello.txt</Key>
	<ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`)
	})

	opt := &CompleteMultipartUploadOptions{
		Parts: []Object{{PartNumber: 1, ETag: `"abc"`}},
	}
	ref, _, err := client.Object.CompleteMultipartUpload(context.Background(), name, uploadID, opt)
	if err != nil {
		t.Fatalf("CompleteMultipartUpload should succeed on retry, got error: %v", err)
	}
	if ref.ETag != `"594f98b11c6901c0f0683de1085a6d0e-4"` {
		t.Errorf("CompleteMultipartUpload ETag = %v, want %v", ref.ETag, `"594f98b11c6901c0f0683de1085a6d0e-4"`)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

// TestObjectService_CompleteMultipartUpload_EmptyETag_AllRetryFail 模拟所有重试均返回空 ETag，
// 验证最终返回 RetryError 并包含所有错误信息。
func TestObjectService_CompleteMultipartUpload_EmptyETag_AllRetryFail(t *testing.T) {
	setup()
	defer teardown()

	client.Conf.RetryOpt.Count = 3

	name := "test/hello.txt"
	uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"
	callCount := 0
	errBody := `<CompleteMultipartUploadResult><Location>loc</Location><Bucket>b</Bucket><Key>k</Key></CompleteMultipartUploadResult>`

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		callCount++
		// 每次都返回缺少 ETag 的 body
		fmt.Fprint(w, errBody)
	})

	opt := &CompleteMultipartUploadOptions{
		Parts: []Object{{PartNumber: 1, ETag: `"abc"`}},
	}
	_, _, err := client.Object.CompleteMultipartUpload(context.Background(), name, uploadID, opt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// 必须是 RetryError 类型
	retryErr, ok := err.(*RetryError)
	if !ok {
		t.Fatalf("expected *RetryError, got %T: %v", err, err)
	}
	// RetryOpt.Count=3，每次都失败，应收集到 3 条错误
	if len(retryErr.Errs) != 3 {
		t.Errorf("expected 3 errors in RetryError, got %d: %v", len(retryErr.Errs), retryErr)
	}
	// 每条错误信息应包含 body 内容（字符串形式）
	for i, e := range retryErr.Errs {
		if !strings.Contains(e.Error(), "body contains an error") {
			t.Errorf("error[%d] = %v, want contains 'body contains an error'", i, e)
		}
		if !strings.Contains(e.Error(), "<CompleteMultipartUploadResult>") {
			t.Errorf("error[%d] = %v, want contains XML string body", i, e)
		}
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

// TestObjectService_CompleteMultipartUpload_HTTPError 模拟服务端返回非 200 错误。
func TestObjectService_CompleteMultipartUpload_HTTPError(t *testing.T) {
	setup()
	defer teardown()

	name := "test/hello.txt"
	uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `<Error><Code>InternalError</Code><Message>server error</Message></Error>`)
	})

	opt := &CompleteMultipartUploadOptions{
		Parts: []Object{{PartNumber: 1, ETag: `"abc"`}},
	}
	_, _, err := client.Object.CompleteMultipartUpload(context.Background(), name, uploadID, opt)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
}

// TestObjectService_CompleteMultipartUpload_200NonXML 覆盖 status=200 但 body 不是合法 XML 的场景。
// COS 规范：合并分块时服务端先返回 200，然后在 body 中流式写入结果；
// 若 body 完全不是 XML（如纯文本、HTML 错误页），需当作错误处理并触发重试。
func TestObjectService_CompleteMultipartUpload_200NonXML(t *testing.T) {
	cases := []struct {
		name string
		body string
		// 期望 lastErr.Error() 包含的关键字（按实际 xml 解析行为）
		wantContain string
		// 是否期望 lastErr 是 xml 解析相关错误（非 COS ErrorResponse）
		notCOSError bool
	}{
		{
			// 纯文本：xml.Unmarshal 遇到非 XML 内容抛 SyntaxError，
			// 再经 UnmarshalCompleteMultiUploadResult 仍失败，最终是 EOF
			name:        "plain text body",
			body:        "Internal Server Error",
			wantContain: "",
			notCOSError: true,
		},
		{
			// HTML 页面：xml 根标签是 <html>，不是 SyntaxError 而是解析错误
			name:        "html error page",
			body:        "<html><body>Bad Gateway</body></html>",
			wantContain: "html",
			notCOSError: true,
		},
		{
			// 空 body：xml.Unmarshal 返回 SyntaxError/EOF，
			// UnmarshalCompleteMultiUploadResult 也返回 EOF
			name:        "empty body",
			body:        "",
			wantContain: "",
			notCOSError: true,
		},
		{
			// 200 body 内嵌 XML Error（COS chunk 流式传输的真实错误格式）：
			// xml.Unmarshal 能解析但根元素不是 CompleteMultipartUploadResult，返回解析错误
			name:        "xml error embedded in 200",
			body:        `<?xml version="1.0" encoding="UTF-8"?><Error><Code>InternalError</Code><Message>We encountered an internal error.</Message></Error>`,
			wantContain: "Error",
			notCOSError: true,
		},
		{
			// XML 结构合法但 ETag 为空：触发 "body contains an error" 路径
			name:        "valid xml but empty ETag",
			body:        `<CompleteMultipartUploadResult><Location>loc</Location><Bucket>b</Bucket><Key>k</Key></CompleteMultipartUploadResult>`,
			wantContain: "body contains an error",
			notCOSError: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()

			uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"
			mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPost)
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, tc.body)
			})

			opt := &CompleteMultipartUploadOptions{
				Parts: []Object{{PartNumber: 1, ETag: `"abc"`}},
			}
			_, _, err := client.Object.CompleteMultipartUpload(
				context.Background(), "test/hello.txt", uploadID, opt)
			if err == nil {
				t.Fatal("expected error for 200 non-XML body, got nil")
			}

			// 错误应当被包在 RetryError 里
			retryErr, ok := err.(*RetryError)
			if !ok {
				t.Fatalf("expected *RetryError, got %T: %v", err, err)
			}
			if len(retryErr.Errs) == 0 {
				t.Fatal("RetryError.Errs should not be empty")
			}

			// 以最后一次重试的错误为准
			lastErr := retryErr.Errs[len(retryErr.Errs)-1]
			if tc.wantContain != "" && !strings.Contains(lastErr.Error(), tc.wantContain) {
				t.Errorf("last error = %q, want contains %q", lastErr.Error(), tc.wantContain)
			}

			// 200 非 XML body 不应被误判为 COS 服务端错误（ErrorResponse）
			if tc.notCOSError {
				if _, isCOS := IsCOSError(err); isCOS {
					t.Error("200 non-XML body should not be identified as COS ErrorResponse")
				}
			}
		})
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
	_, _, err = client.Object.CopyPart(context.Background(),
		name, uploadID, partNumber, "http://"+sourceUrl, opt)
	if err == nil {
		t.Fatalf("Object.CopyPart returned error is nil")
	}
	_, _, err = client.Object.CopyPart(context.Background(),
		name, uploadID, partNumber, "test-1253846586.cos.ap-guangzhou.myqcloud.com", opt)
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

	soruceURL = fmt.Sprintf("%s/%s?versionId=versionid", client.BaseURL.BucketURL.Host, "test.src.copy")
	_, _, err = client.Object.MultiCopy(context.Background(), dest, soruceURL, opt)
	if err != nil {
		t.Errorf("Object.MultiCopy failed %v", err)
	}

}

func TestObjectService_UploadPartRetry(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectUploadPartOptions{
		Listener: &DefaultProgressListener{},
	}
	name := "test/hello.txt"
	data := make([]byte, 1024*1024*3)
	_, err := rand.Read(data)
	tb := crc64.MakeTable(crc64.ECMA)
	realcrc := crc64.Update(0, tb, data)
	uploadID := "xxxxx"
	partNumber := 1

	nr, count := 0, 3
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"uploadId":   uploadID,
			"partNumber": "1",
		}
		testFormValues(t, r, vs)

		b, _ := ioutil.ReadAll(r.Body)
		crc := crc64.Update(0, tb, b)
		if !reflect.DeepEqual(crc, realcrc) {
			t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
		}
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
		nr++
		if nr < count {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	_, err = client.Object.UploadPart(context.Background(), name, uploadID, partNumber, bytes.NewReader(data), opt)
	if err != nil || nr != count {
		t.Errorf("Object.UploadPart failed: %v", err)
	}

	nr, count = 0, 3
	_, err = client.Object.UploadPart(context.Background(), name, uploadID, partNumber, strings.NewReader(string(data)), opt)
	if err != nil || nr != count {
		t.Errorf("Object.UploadPart failed: %v", err)
	}

	// 非io.Seeker不做重试
	nr, count = 0, 3
	_, err = client.Object.UploadPart(context.Background(), name, uploadID, partNumber, bytes.NewBuffer(data), opt)
	if err == nil || nr != 1 {
		t.Errorf("Object.UploadPart failed: %v", err)
	}

	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create tmp file failed")
	}
	defer os.Remove(filePath)
	// 源文件内容
	newfile.Write(data)
	newfile.Close()

	fd, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("open file failed: %v", err)
	}
	// 文件关闭, 不做重试
	nr, count = 0, 3
	_, err = client.Object.UploadPart(context.Background(), name, uploadID, partNumber, fd, opt)
	if err == nil || nr != 1 {
		t.Errorf("Object.UploadPart failed: %v", err)
	}
}

// 测试 MultiCopy 所有分块持续失败（耗尽重试）后返回的 error
func TestObjectService_MultiCopyPartFailed(t *testing.T) {
	setup()
	defer teardown()

	totalBytes := 1024 * 1024 * 35

	mux.HandleFunc("/test.src.copyfail", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "HEAD")
		w.Header().Add("Content-Length", strconv.FormatInt(int64(totalBytes), 10))
	})

	uploadid := "test-multicopy-partfailed-id"
	mux.HandleFunc("/test.dst.copyfail", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Method == http.MethodPost {
			r.ParseForm()
			if r.Form.Get("uploadId") != "" {
				// CompleteMultipartUpload - 不应该到达这里
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// InitiateMultipartUpload
			fmt.Fprintf(w, `<InitiateMultipartUploadResult>
	<Bucket></Bucket>
	<Key>test.dst.copyfail</Key>
	<UploadId>%v</UploadId>
</InitiateMultipartUploadResult>`, uploadid)
		}
	})

	opt := &MultiCopyOptions{
		PartSize:       1,
		ThreadPoolSize: 3,
		useMulti:       true,
	}
	dest := "test.dst.copyfail"
	sourceURL := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, "test.src.copyfail")
	_, _, err := client.Object.MultiCopy(context.Background(), dest, sourceURL, opt)
	if err == nil {
		t.Fatalf("MultiCopy should fail when all parts return 500")
	}
	if !strings.Contains(err.Error(), "failed to get resp content") {
		t.Fatalf("MultiCopy error should contain 'failed to get resp content', got: %v", err)
	}
}

// =============================================================================
// CompleteMultipartUploadWithHead
// =============================================================================

const (
	cwhTestName     = "test/hello.txt"
	cwhTestUploadID = "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"
	cwhTestETag     = `"594f98b11c6901c0f0683de1085a6d0e-4"`
	cwhTestCRC      = uint64(1234567890123)
)

// --- helpers ---

// logServerReq 在 mock handler 入口打印收到的请求信息
func logServerReq(t *testing.T, r *http.Request) {
	t.Helper()
	t.Logf("[mock-server] recv: method=%s path=%s query=%q", r.Method, r.URL.Path, r.URL.RawQuery)
}

// logServerResp 在 mock handler 响应分支打印将要返回的响应信息
func logServerResp(t *testing.T, status int, note string) {
	t.Helper()
	t.Logf("[mock-server] send: status=%d %s", status, note)
}

// logClientResp 在客户端调用后打印收到的响应与错误
func logClientResp(t *testing.T, label string, result interface{}, err error) {
	t.Helper()
	if err != nil {
		t.Logf("[client] %s: err=%v", label, err)
		return
	}
	t.Logf("[client] %s: ok, result=%+v", label, result)
}

// writeNoSuchUpload 模拟服务端在 uploadId 已失效时返回 404 NoSuchUpload
func writeNoSuchUpload(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, `<Error><Code>NoSuchUpload</Code><Message>The specified upload does not exist.</Message></Error>`)
	logServerResp(t, http.StatusNotFound, "NoSuchUpload")
}

// writeCompleteSuccess 模拟 CompleteMultipartUpload 成功响应
func writeCompleteSuccess(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	fmt.Fprint(w, `<CompleteMultipartUploadResult>
	<Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
	<Bucket>test</Bucket>
	<Key>test/hello.txt</Key>
	<ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`)
	logServerResp(t, http.StatusOK, "CompleteMultipartUploadResult")
}

func cwhOptions() *CompleteMultipartUploadOptions {
	return &CompleteMultipartUploadOptions{
		Parts: []Object{{PartNumber: 1, ETag: `"abc"`}},
	}
}

func cwhOptionsVerifyCRC(crc uint64) *CompleteMultipartUploadOptions {
	opt := cwhOptions()
	opt.Crc64CheckInHead = &crc
	return opt
}

// --- T1: Complete 首次成功 → 不发起 HEAD ---

func TestCompleteMultipartUploadWithHead_FirstSuccess(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok (no crc header)")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		writeCompleteSuccess(t, w)
	})

	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("CompleteMultipartUploadWithHead failed: %v", err)
	}
	if ref.ETag != cwhTestETag {
		t.Errorf("ETag = %v, want %v", ref.ETag, cwhTestETag)
	}
	if postCount != 1 {
		t.Errorf("post count = %d, want 1", postCount)
	}
	if headCount != 0 {
		t.Errorf("head count = %d, want 0 when complete succeeds", headCount)
	}
}

// --- T2: 首次失败超时 → 重试成功 → 不发起 HEAD ---

func TestCompleteMultipartUploadWithHead_RetryToSuccess(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code><Message>server error</Message></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeCompleteSuccess(t, w)
	})

	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("should succeed on retry: %v", err)
	}
	if ref.ETag != cwhTestETag {
		t.Errorf("ETag = %v, want %v", ref.ETag, cwhTestETag)
	}
	if headCount != 0 {
		t.Errorf("head count = %d, want 0 when retry succeeds", headCount)
	}
}

// --- T3: 重试中命中 NoSuchUpload + HEAD 成功 + Crc64CheckInHead=nil → 判定成功 ---

func TestCompleteMultipartUploadWithHead_HeadMatchNoCRCCheck(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.Header().Set("ETag", `"head-etag"`)
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code><Message>server error</Message></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeNoSuchUpload(t, w)
	})

	// Crc64CheckInHead 为 nil
	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptions())
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("should treat as upload completed: %v", err)
	}
	if ref.ETag != `"head-etag"` {
		t.Errorf("ETag = %v, want %v", ref.ETag, `"head-etag"`)
	}
	if ref.Key != cwhTestName {
		t.Errorf("Key = %v, want %v", ref.Key, cwhTestName)
	}
	// 首次 500 + 第 1 次重试 404 → 中断，不再进行第 2 次重试
	if postCount != 2 {
		t.Errorf("post count = %d, want 2 (stop retry after NoSuchUpload)", postCount)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T4: 重试中命中 NoSuchUpload + HEAD 成功 + CRC64 一致 → 判定成功 ---

func TestCompleteMultipartUploadWithHead_HeadMatchCRCEqual(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.Header().Set("x-cos-hash-crc64ecma", strconv.FormatUint(cwhTestCRC, 10))
			w.Header().Set("ETag", `"head-etag"`)
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code><Message>server error</Message></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeNoSuchUpload(t, w)
	})

	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("should treat as upload completed: %v", err)
	}
	if ref.ETag != `"head-etag"` {
		t.Errorf("ETag = %v, want %v", ref.ETag, `"head-etag"`)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T5: 重试中命中 NoSuchUpload + HEAD 成功 + CRC64 不一致 → 返回 Complete 错误 ---

func TestCompleteMultipartUploadWithHead_CRCMismatch(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.Header().Set("x-cos-hash-crc64ecma", "999")
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeNoSuchUpload(t, w)
	})

	_, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
	if err == nil {
		t.Fatal("expected error when crc64 mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "crc64 mismatch") {
		t.Errorf("error = %v, want contains 'crc64 mismatch'", err)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T6: 重试中命中 NoSuchUpload + HEAD 404 → 返回 Complete 错误 ---

func TestCompleteMultipartUploadWithHead_HeadNotFound(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusNotFound)
			logServerResp(t, http.StatusNotFound, "HEAD 404 NotFound")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeNoSuchUpload(t, w)
	})

	_, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptions())
	logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
	if err == nil {
		t.Fatal("expected error when head returns 404, got nil")
	}
	if !strings.Contains(err.Error(), "NoSuchUpload") {
		t.Errorf("error = %v, want contains NoSuchUpload", err)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T7: 重试中命中 NoSuchUpload + HEAD 其他错误 → 返回 Complete 错误并附带 HEAD 错误 ---

func TestCompleteMultipartUploadWithHead_HeadInternalError(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount := 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code><Message>head failed</Message></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeNoSuchUpload(t, w)
	})

	_, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptions())
	logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "NoSuchUpload") {
		t.Errorf("error = %v, want contains NoSuchUpload", err)
	}
	if !strings.Contains(err.Error(), "head object failed") {
		t.Errorf("error = %v, want contains 'head object failed'", err)
	}
}

// --- T8: Crc64CheckInHead != nil + 服务端未返回 crc64 → 返回 Complete 错误 ---

func TestCompleteMultipartUploadWithHead_NoCRCHeader(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK) // 不设置 x-cos-hash-crc64ecma
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		if postCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
			return
		}
		writeNoSuchUpload(t, w)
	})

	_, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
	if err == nil {
		t.Fatal("expected error when crc64 header missing, got nil")
	}
	if !strings.Contains(err.Error(), "x-cos-hash-crc64ecma") {
		t.Errorf("error = %v, want contains 'x-cos-hash-crc64ecma'", err)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T9: Complete 首次返回 NoSuchUpload（nr==0）→ 不发起 HEAD，继续重试 ---

func TestCompleteMultipartUploadWithHead_FirstNoSuchUploadNotHead(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 2

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.Header().Set("x-cos-hash-crc64ecma", strconv.FormatUint(cwhTestCRC, 10))
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		writeNoSuchUpload(t, w)
	})

	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("should succeed on retry fallback: %v", err)
	}
	if ref.Key != cwhTestName {
		t.Errorf("Key = %v, want %v", ref.Key, cwhTestName)
	}
	// 首次 404（nr=0）不 HEAD，第 2 次（即重试）命中才触发
	if postCount != 2 {
		t.Errorf("post count = %d, want 2 (first failure must not trigger head)", postCount)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T10: 重试返回非 NoSuchUpload（5xx）→ 按原规则继续重试，次数耗尽不触发 HEAD ---

func TestCompleteMultipartUploadWithHead_NonNoSuchUploadRetryExhausted(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
		logServerResp(t, http.StatusInternalServerError, "InternalError")
	})

	_, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if headCount != 0 {
		t.Errorf("head count = %d, want 0 for non-NoSuchUpload errors", headCount)
	}
	if postCount != 3 {
		t.Errorf("post count = %d, want 3 (retry exhausted)", postCount)
	}
	if _, ok := err.(*RetryError); !ok {
		t.Errorf("expected *RetryError, got %T", err)
	}
}

// --- T11: 重试返回 403（不可重试错误）→ 原样返回，不发起 HEAD ---
// 说明：Complete 内部的循环对任何 err != nil 都 retry.Add 并继续循环，403 也会走完 count 次。
// 这里主要验证「不触发 HEAD」。

func TestCompleteMultipartUploadWithHead_Returns403NoHead(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 2

	headCount := 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `<Error><Code>AccessDenied</Code><Message>denied</Message></Error>`)
		logServerResp(t, http.StatusForbidden, "AccessDenied")
	})

	_, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if headCount != 0 {
		t.Errorf("head count = %d, want 0", headCount)
	}
}

// --- T12: 重试次数为 3，第 2 次重试命中 NoSuchUpload → 中断重试，不执行第 3 次 ---

func TestCompleteMultipartUploadWithHead_BreakRemainingRetries(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.Header().Set("x-cos-hash-crc64ecma", strconv.FormatUint(cwhTestCRC, 10))
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		// 请求 1（首次）: 500；请求 2（第 1 次重试）: 500；请求 3（第 2 次重试）: 404
		switch postCount {
		case 1, 2:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
			logServerResp(t, http.StatusInternalServerError, "InternalError")
		default:
			writeNoSuchUpload(t, w)
		}
	})

	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID,
		cwhOptionsVerifyCRC(cwhTestCRC))
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("should succeed via head fallback: %v", err)
	}
	if ref.Key != cwhTestName {
		t.Errorf("Key = %v, want %v", ref.Key, cwhTestName)
	}
	// 必须恰好 3 次 POST（首次 + 2 次重试），兜底后不再进行"第 3 次重试"
	if postCount != 3 {
		t.Errorf("post count = %d, want exactly 3 (break remaining retries)", postCount)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// --- T13: 原 CompleteMultipartUpload 调用时设置 Crc64CheckInHead → 字段被忽略，不发 HEAD、不进请求体 ---

func TestCompleteMultipartUpload_IgnoresCrc64CheckInHead(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	postCount, headCount := 0, 0
	bodyBytes := []byte{}
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		postCount++
		b, _ := io.ReadAll(r.Body)
		bodyBytes = append(bodyBytes, b...)
		writeNoSuchUpload(t, w)
	})

	opt := cwhOptionsVerifyCRC(cwhTestCRC)
	_, _, err := client.Object.CompleteMultipartUpload(context.Background(), cwhTestName, cwhTestUploadID, opt)
	logClientResp(t, "CompleteMultipartUpload", nil, err)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// 字段必须被忽略：不发起 HEAD、按原重试策略跑完 3 次
	if headCount != 0 {
		t.Errorf("head count = %d, want 0 (CompleteMultipartUpload must ignore Crc64CheckInHead)", headCount)
	}
	if postCount != 3 {
		t.Errorf("post count = %d, want 3 (retry exhausted)", postCount)
	}
	// 字段必须不进入 HTTP body
	if strings.Contains(string(bodyBytes), "Crc64CheckInHead") {
		t.Errorf("request body must not contain Crc64CheckInHead, got: %s", bodyBytes)
	}
}

// --- T14: Drop-in 替换：不设 Crc64CheckInHead 调用新方法 → 行为与原方法一致 ---

func TestCompleteMultipartUploadWithHead_DropInReplacement(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 2

	headCount := 0
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		if r.Method == http.MethodHead {
			headCount++
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
			return
		}
		testMethod(t, r, http.MethodPost)
		writeCompleteSuccess(t, w)
	})

	// 完全使用与 CompleteMultipartUpload 相同的调用形式
	ref, _, err := client.Object.CompleteMultipartUploadWithHead(context.Background(), cwhTestName, cwhTestUploadID, cwhOptions())
	logClientResp(t, "CompleteMultipartUploadWithHead", ref, err)
	if err != nil {
		t.Fatalf("drop-in call failed: %v", err)
	}
	if ref.ETag != cwhTestETag {
		t.Errorf("ETag = %v, want %v", ref.ETag, cwhTestETag)
	}
	if headCount != 0 {
		t.Errorf("head count = %d, want 0 for drop-in (no Crc64CheckInHead)", headCount)
	}
}

// --- 辅助：isNoSuchUploadError 单测 ---

func TestIsNoSuchUploadError(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 1

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		testMethod(t, r, http.MethodPost)
		writeNoSuchUpload(t, w)
	})

	_, _, err := client.Object.CompleteMultipartUpload(context.Background(), cwhTestName, cwhTestUploadID, cwhOptions())
	logClientResp(t, "CompleteMultipartUpload", nil, err)
	if !isNoSuchUploadError(err) {
		t.Errorf("isNoSuchUploadError = false, want true for 404 NoSuchUpload, err: %v", err)
	}
	if isNoSuchUploadError(nil) {
		t.Error("isNoSuchUploadError(nil) should be false")
	}
}

// TestCompleteMultipartUploadWithHead_FallbackErrIsCOSError 验证兜底失败时返回的
// error 仍保留原始 Complete 的 *ErrorResponse 类型信息，调用方可用 IsCOSError 识别。
// 覆盖三条用 %w 包装的失败分支：HEAD 自身失败 / CRC header 缺失 / CRC 不匹配。
func TestCompleteMultipartUploadWithHead_FallbackErrIsCOSError(t *testing.T) {
	cases := []struct {
		name     string
		headResp func(w http.ResponseWriter)
		opt      func() *CompleteMultipartUploadOptions
	}{
		{
			name: "head_internal_error",
			headResp: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
			},
			opt: cwhOptions,
		},
		{
			name: "head_missing_crc_header",
			headResp: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
			},
			opt: func() *CompleteMultipartUploadOptions { return cwhOptionsVerifyCRC(cwhTestCRC) },
		},
		{
			name: "crc_mismatch",
			headResp: func(w http.ResponseWriter) {
				w.Header().Set("x-cos-hash-crc64ecma", "999")
				w.WriteHeader(http.StatusOK)
			},
			opt: func() *CompleteMultipartUploadOptions { return cwhOptionsVerifyCRC(cwhTestCRC) },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()
			client.Conf.RetryOpt.Count = 3

			postCount := 0
			mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
				logServerReq(t, r)
				if r.Method == http.MethodHead {
					tc.headResp(w)
					logServerResp(t, 0, "HEAD (per case)")
					return
				}
				postCount++
				if postCount == 1 {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, `<Error><Code>InternalError</Code></Error>`)
					logServerResp(t, http.StatusInternalServerError, "InternalError")
					return
				}
				writeNoSuchUpload(t, w)
			})

			_, _, err := client.Object.CompleteMultipartUploadWithHead(
				context.Background(), cwhTestName, cwhTestUploadID, tc.opt())
			logClientResp(t, "CompleteMultipartUploadWithHead", nil, err)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			// 关键断言：%w 包装后仍能被 IsCOSError 识别为原始 Complete 的 NoSuchUpload
			cosErr, ok := IsCOSError(err)
			if !ok {
				t.Fatalf("IsCOSError = false, want true; err type=%T, err=%v", err, err)
			}
			if cosErr.Code != "NoSuchUpload" {
				t.Errorf("cosErr.Code = %q, want NoSuchUpload", cosErr.Code)
			}
			if cosErr.Response == nil || cosErr.Response.StatusCode != http.StatusNotFound {
				t.Errorf("cosErr.Response.StatusCode = %v, want 404", cosErr.Response)
			}
			// isNoSuchUploadError 也应穿透
			if !isNoSuchUploadError(err) {
				t.Errorf("isNoSuchUploadError = false, want true, err=%v", err)
			}
		})
	}
}

// =============================================================================
// Upload / MultiUploadOptions.EnableCompleteWithHead
// =============================================================================

// uploadTestFile 构造一份随机内容的本地文件（≥ 1MB，分块路径），返回路径与 CRC64。
func uploadTestFile(t *testing.T) (string, uint64) {
	t.Helper()
	f, err := ioutil.TempFile("", "cos-upload-complete-with-head")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	// 2.5MB 内容，确保 partSize=1MB 时分块数 >= 2
	block := bytes.Repeat([]byte("a"), 1024*1024)
	for i := 0; i < 2; i++ {
		if _, err := f.Write(block); err != nil {
			t.Fatalf("write temp file: %v", err)
		}
	}
	if _, err := f.Write(bytes.Repeat([]byte("b"), 512*1024)); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	name := f.Name()
	fd, err := os.Open(name)
	if err != nil {
		t.Fatalf("open temp file: %v", err)
	}
	defer fd.Close()
	crc, err := calCRC64(fd)
	if err != nil {
		t.Fatalf("calCRC64: %v", err)
	}
	return name, crc
}

// uploadNoSuchUploadHandler 模拟整条 Upload 链路：Init 成功、UploadPart 回传匹配 CRC 的 ETag、
// Complete 固定返回 404 NoSuchUpload；HEAD 回传指定 CRC64。
func uploadNoSuchUploadHandler(t *testing.T, headCRC uint64, headETag string, completeCount, headCount *int) {
	t.Helper()
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		logServerReq(t, r)
		switch r.Method {
		case http.MethodPut:
			// UploadPart：计算分块 body 的 CRC64 回传，避免 SDK 的 part CRC 校验失败
			crc, err := calCRC64(r.Body)
			if err != nil {
				t.Errorf("calCRC64 of part body: %v", err)
			}
			w.Header().Set("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
			w.Header().Set("ETag", `"part-etag"`)
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "UploadPart ok")
		case http.MethodPost:
			if r.URL.Query().Get("uploadId") == "" {
				fmt.Fprintf(w, `<InitiateMultipartUploadResult><Bucket>test</Bucket><Key>test/hello.txt</Key><UploadId>%s</UploadId></InitiateMultipartUploadResult>`, cwhTestUploadID)
				logServerResp(t, http.StatusOK, "InitiateMultipartUploadResult")
				return
			}
			*completeCount++
			writeNoSuchUpload(t, w)
		case http.MethodHead:
			*headCount++
			if headCRC != 0 {
				w.Header().Set("x-cos-hash-crc64ecma", strconv.FormatUint(headCRC, 10))
			}
			w.Header().Set("ETag", headETag)
			w.WriteHeader(http.StatusOK)
			logServerResp(t, http.StatusOK, "HEAD ok")
		}
	})
}

// TestObjectService_Upload_EnableCompleteWithHead_SuccessViaHead 开启开关，
// Complete 重试过程中命中 NoSuchUpload，HEAD + CRC64 一致，Upload 应判定成功。
func TestObjectService_Upload_EnableCompleteWithHead_SuccessViaHead(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	filepath, localcrc := uploadTestFile(t)
	defer os.Remove(filepath)

	completeCount, headCount := 0, 0
	uploadNoSuchUploadHandler(t, localcrc, `"head-etag"`, &completeCount, &headCount)

	v, _, err := client.Object.Upload(context.Background(), "test/hello.txt", filepath, &MultiUploadOptions{
		PartSize:               1,
		ThreadPoolSize:         1,
		EnableCompleteWithHead: true,
	})
	logClientResp(t, "Upload", v, err)
	if err != nil {
		t.Fatalf("Upload should succeed via head fallback: %v", err)
	}
	if v.ETag != `"head-etag"` {
		t.Errorf("ETag = %v, want %v", v.ETag, `"head-etag"`)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
	if completeCount != 2 {
		// 首次 404（nr=0）不触发 HEAD，第 1 次重试命中后中断
		t.Errorf("complete count = %d, want 2 (first + one retry)", completeCount)
	}
}

// TestObjectService_Upload_EnableCompleteWithHead_Disabled 开关关闭（默认），
// Complete 全部命中 NoSuchUpload 但 Upload 不发起 HEAD、按原策略跑完重试返回错误。
func TestObjectService_Upload_EnableCompleteWithHead_Disabled(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	filepath, localcrc := uploadTestFile(t)
	defer os.Remove(filepath)

	completeCount, headCount := 0, 0
	uploadNoSuchUploadHandler(t, localcrc, `"head-etag"`, &completeCount, &headCount)

	_, _, err := client.Object.Upload(context.Background(), "test/hello.txt", filepath, &MultiUploadOptions{
		PartSize:       1,
		ThreadPoolSize: 1,
		// EnableCompleteWithHead 默认 false
	})
	logClientResp(t, "Upload", nil, err)
	if err == nil {
		t.Fatal("expected error when switch off, got nil")
	}
	if headCount != 0 {
		t.Errorf("head count = %d, want 0 when switch off", headCount)
	}
	if completeCount != 3 {
		t.Errorf("complete count = %d, want 3 (retry exhausted)", completeCount)
	}
}

// TestObjectService_Upload_EnableCompleteWithHead_CRCMismatch 开启开关，
// HEAD 到对象但 CRC64 与本地不一致 → Upload 返回错误（符合安全预期：宁可失败不可误判）。
func TestObjectService_Upload_EnableCompleteWithHead_CRCMismatch(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	filepath, _ := uploadTestFile(t)
	defer os.Remove(filepath)

	completeCount, headCount := 0, 0
	// HEAD 回传一个不等于 localcrc 的 CRC
	uploadNoSuchUploadHandler(t, 1, `"head-etag"`, &completeCount, &headCount)

	_, _, err := client.Object.Upload(context.Background(), "test/hello.txt", filepath, &MultiUploadOptions{
		PartSize:               1,
		ThreadPoolSize:         1,
		EnableCompleteWithHead: true,
	})
	logClientResp(t, "Upload", nil, err)
	if err == nil {
		t.Fatal("expected error when crc64 mismatch, got nil")
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}

// TestObjectService_Upload_EnableCompleteWithHead_DisableChecksum 开启开关且 DisableChecksum=true，
// Upload 内部不会计算 localcrc，兜底应传 Crc64CheckInHead=nil（仅凭对象存在判定），HEAD 到对象即成功。
func TestObjectService_Upload_EnableCompleteWithHead_DisableChecksum(t *testing.T) {
	setup()
	defer teardown()
	client.Conf.RetryOpt.Count = 3

	filepath, _ := uploadTestFile(t)
	defer os.Remove(filepath)

	completeCount, headCount := 0, 0
	// HEAD 不返回 CRC header（headCRC=0 表示不设置），同时也不需要
	uploadNoSuchUploadHandler(t, 0, `"head-etag"`, &completeCount, &headCount)

	v, _, err := client.Object.Upload(context.Background(), "test/hello.txt", filepath, &MultiUploadOptions{
		PartSize:               1,
		ThreadPoolSize:         1,
		DisableChecksum:        true,
		EnableCompleteWithHead: true,
	})
	logClientResp(t, "Upload", v, err)
	if err != nil {
		t.Fatalf("Upload should succeed via head fallback (no crc check): %v", err)
	}
	if v.ETag != `"head-etag"` {
		t.Errorf("ETag = %v, want %v", v.ETag, `"head-etag"`)
	}
	if headCount != 1 {
		t.Errorf("head count = %d, want 1", headCount)
	}
}
