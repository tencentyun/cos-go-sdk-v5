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
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestObjectService_Get(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"response-content-type": "text/html",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "Range", "bytes=0-3")
		fmt.Fprint(w, `hello`)
	})

	opt := &ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}

	resp, err := client.Object.Get(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Get returned error: %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	ref := string(b)
	want := "hello"
	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.Get returned %+v, want %+v", ref, want)
	}

}

func TestObjectService_Put(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"

	retry := 0
	final := 10
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		if retry%2 == 0 {
			b, _ := ioutil.ReadAll(r.Body)
			tb := crc64.MakeTable(crc64.ECMA)
			crc := crc64.Update(0, tb, b)
			v := string(b)
			want := "hello"
			if !reflect.DeepEqual(v, want) {
				t.Errorf("Object.Put request body: %#v, want %#v", v, want)
			}
			realcrc := crc64.Update(0, tb, []byte("hello"))
			if !reflect.DeepEqual(crc, realcrc) {
				t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
			}
			w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
			if retry != final {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		} else {
			w.Header().Add("x-cos-hash-crc64ecma", "123456789")
		}
	})

	for retry <= final {
		r := bytes.NewReader([]byte("hello"))
		_, err := client.Object.Put(context.Background(), name, r, opt)
		if retry < final && err == nil {
			t.Fatalf("Error must not nil when retry < final")
		}
		if retry == final && err != nil {
			t.Fatalf("Put Error: %v", err)
		}
		retry++
	}
}

func TestObjectService_PutFromFile(t *testing.T) {
	setup()
	defer teardown()

	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create tmp file failed")
	}
	defer os.Remove(filePath)
	// 源文件内容
	b := make([]byte, 1024*1024*3)
	_, err = rand.Read(b)
	newfile.Write(b)
	newfile.Close()

	tb := crc64.MakeTable(crc64.ECMA)
	realcrc := crc64.Update(0, tb, b)
	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"
	retry := 0
	final := 4
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		if retry%2 == 0 {
			bs, _ := ioutil.ReadAll(r.Body)
			crc := crc64.Update(0, tb, bs)
			if !reflect.DeepEqual(bs, b) {
				t.Errorf("Object.Put request body Error")
			}
			if !reflect.DeepEqual(crc, realcrc) {
				t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
			}
			w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
			if retry != final {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		} else {
			w.Header().Add("x-cos-hash-crc64ecma", "123456789")
		}
	})

	for retry <= final {
		_, err := client.Object.PutFromFile(context.Background(), name, filePath, opt)
		if retry < final && err == nil {
			t.Fatalf("Error must not nil when retry < final")
		}
		if retry == final && err != nil {
			t.Fatalf("Put Error: %v", err)
		}
		retry++
	}
}

func TestObjectService_Delete(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		t.Fatalf("Object.Delete returned error: %v", err)
	}
}

func TestObjectService_Head(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "HEAD")
		testHeader(t, r, "If-Modified-Since", "Mon, 12 Jun 2017 05:36:19 GMT")
	})

	opt := &ObjectHeadOptions{
		IfModifiedSince: "Mon, 12 Jun 2017 05:36:19 GMT",
	}

	_, err := client.Object.Head(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Head returned error: %v", err)
	}
}

func TestObjectService_Options(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodOptions)
		testHeader(t, r, "Access-Control-Request-Method", "PUT")
		testHeader(t, r, "Origin", "www.qq.com")
	})

	opt := &ObjectOptionsOptions{
		Origin:                     "www.qq.com",
		AccessControlRequestMethod: "PUT",
	}

	_, err := client.Object.Options(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Options returned error: %v", err)
	}

}

// func TestObjectService_Append(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	opt := &ObjectPutOptions{
// 		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
// 			ContentType: "text/html",
// 		},
// 		ACLHeaderOptions: &ACLHeaderOptions{
// 			XCosACL: "private",
// 		},
// 	}
// 	name := "test/hello.txt"
// 	position := 0

// 	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
// 		vs := values{
// 			"append":   "",
// 			"position": "0",
// 		}
// 		testFormValues(t, r, vs)

// 		testMethod(t, r, http.MethodPost)
// 		testHeader(t, r, "x-cos-acl", "private")
// 		testHeader(t, r, "Content-Type", "text/html")

// 		b, _ := ioutil.ReadAll(r.Body)
// 		v := string(b)
// 		want := "hello"
// 		if !reflect.DeepEqual(v, want) {
// 			t.Errorf("Object.Append request body: %#v, want %#v", v, want)
// 		}
// 	})

// 	r := bytes.NewReader([]byte("hello"))
// 	_, err := client.Object.Append(context.Background(), name, position, r, opt)
// 	if err != nil {
// 		t.Fatalf("Object.Append returned error: %v", err)
// 	}
// }

func TestObjectService_DeleteMulti(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		vs := values{
			"delete": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<DeleteResult>
	<Deleted>
		<Key>test1</Key>
	</Deleted>
	<Deleted>
		<Key>test3</Key>
	</Deleted>
	<Deleted>
		<Key>test2</Key>
	</Deleted>
</DeleteResult>`)
	})

	opt := &ObjectDeleteMultiOptions{
		Objects: []Object{
			{
				Key: "test1",
			},
			{
				Key: "test3",
			},
			{
				Key: "test2",
			},
		},
	}

	ref, _, err := client.Object.DeleteMulti(context.Background(), opt)
	if err != nil {
		t.Fatalf("Object.DeleteMulti returned error: %v", err)
	}

	want := &ObjectDeleteMultiResult{
		XMLName: xml.Name{Local: "DeleteResult"},
		DeletedObjects: []Object{
			{
				Key: "test1",
			},
			{
				Key: "test3",
			},
			{
				Key: "test2",
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.DeleteMulti returned %+v, want %+v", ref, want)
	}

}

func TestObjectService_Copy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.go.copy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		fmt.Fprint(w, `<CopyObjectResult>
		<ETag>"098f6bcd4621d373cade4e832627b4f6"</ETag>
		<LastModified>2017-12-13T14:53:12</LastModified>
	</CopyObjectResult>`)
	})

	sourceURL := "test-1253846586.cos.ap-guangzhou.myqcloud.com/test.source"
	ref, _, err := client.Object.Copy(context.Background(), "test.go.copy", sourceURL, nil)
	if err != nil {
		t.Fatalf("Object.Copy returned error: %v", err)
	}

	want := &ObjectCopyResult{
		XMLName:      xml.Name{Local: "CopyObjectResult"},
		ETag:         `"098f6bcd4621d373cade4e832627b4f6"`,
		LastModified: "2017-12-13T14:53:12",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.Copy returned %+v, want %+v", ref, want)
	}
}

func TestObjectService_Upload(t *testing.T) {
	setup()
	defer teardown()

	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create tmp file failed")
	}
	defer os.Remove(filePath)
	// 源文件内容
	b := make([]byte, 1024*1024*33)
	_, err = rand.Read(b)
	newfile.Write(b)
	newfile.Close()

	// 已上传内容, 10个分块
	rb := make([][]byte, 33)
	uploadid := "test-cos-multiupload-uploadid"
	partmap := make(map[int64]int)
	mux.HandleFunc("/test.go.upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut { // 分块上传
			r.ParseForm()
			part, _ := strconv.ParseInt(r.Form.Get("partNumber"), 10, 64)
			if partmap[part] == 0 {
				// 重试检验1
				partmap[part]++
				ioutil.ReadAll(r.Body)
				w.WriteHeader(http.StatusGatewayTimeout)
			} else if partmap[part] == 1 {
				// 重试校验2
				partmap[part]++
				w.Header().Add("x-cos-hash-crc64ecma", "123456789")
			} else { // 正确上传
				bs, _ := ioutil.ReadAll(r.Body)
				rb[part-1] = bs
				md := hex.EncodeToString(calMD5Digest(bs))
				crc := crc64.Update(0, crc64.MakeTable(crc64.ECMA), bs)
				w.Header().Add("ETag", md)
				w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
			}
		} else {
			testMethod(t, r, http.MethodPost)
			initreq := url.Values{}
			initreq.Set("uploads", "")
			compreq := url.Values{}
			compreq.Set("uploadId", uploadid)
			r.ParseForm()
			if reflect.DeepEqual(r.Form, initreq) {
				// 初始化分块上传
				fmt.Fprintf(w, `<InitiateMultipartUploadResult>
                    <Bucket></Bucket>
                    <Key>%v</Key>
                    <UploadId>%v</UploadId>
                </InitiateMultipartUploadResult>`, "test.go.upload", uploadid)
			} else if reflect.DeepEqual(r.Form, compreq) {
				// 完成分块上传
				tb := crc64.MakeTable(crc64.ECMA)
				crc := uint64(0)
				for _, v := range rb {
					crc = crc64.Update(crc, tb, v)
				}
				w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
				fmt.Fprintf(w, `<CompleteMultipartUploadResult>
                    <Location>/test.go.upload</Location>
                    <Bucket></Bucket>
                    <Key>test.go.upload</Key>
                    <ETag>&quot;%v&quot;</ETag>
                </CompleteMultipartUploadResult>`, hex.EncodeToString(calMD5Digest(b)))
			} else {
				t.Errorf("TestObjectService_Upload Unknown Request")
			}
		}
	})

	opt := &MultiUploadOptions{
		ThreadPoolSize: 3,
		PartSize:       1,
	}
	_, _, err = client.Object.Upload(context.Background(), "test.go.upload", filePath, opt)
	if err != nil {
		t.Fatalf("Object.Upload returned error: %v", err)
	}
}

func TestObjectService_Upload2(t *testing.T) {
	setup()
	defer teardown()

	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create tmp file failed")
	}
	defer os.Remove(filePath)
	// 源文件内容
	b := make([]byte, 1024*1024*3)
	_, err = rand.Read(b)
	newfile.Write(b)
	newfile.Close()

	tb := crc64.MakeTable(crc64.ECMA)
	realcrc := crc64.Update(0, tb, b)
	name := "test/hello.txt"
	retry := 0
	final := 4
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		if retry%2 == 0 {
			bs, _ := ioutil.ReadAll(r.Body)
			crc := crc64.Update(0, tb, bs)
			if !reflect.DeepEqual(bs, b) {
				t.Errorf("Object.Put request body Error")
			}
			if !reflect.DeepEqual(crc, realcrc) {
				t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
			}
			w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
			if retry != final {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		} else {
			w.Header().Add("x-cos-hash-crc64ecma", "123456789")
		}
	})

	mopt := &MultiUploadOptions{
		OptIni: &InitiateMultipartUploadOptions{
			ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
				ContentType: "text/html",
			},
			ACLHeaderOptions: &ACLHeaderOptions{
				XCosACL: "private",
			},
		},
	}
	for retry <= final {
		_, _, err := client.Object.Upload(context.Background(), name, filePath, mopt)
		if retry < final && err == nil {
			t.Fatalf("Error must not nil when retry < final")
		}
		if retry == final && err != nil {
			t.Fatalf("Put Error: %v", err)
		}
		retry++
	}
}
