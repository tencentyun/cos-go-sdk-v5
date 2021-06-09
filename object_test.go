package cos

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestObjectService_Get(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"
	contentLength := 1024 * 1024 * 10
	data := make([]byte, contentLength)

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"response-content-type": "text/html",
		}
		testFormValues(t, r, vs)
		strRange := r.Header.Get("Range")
		slice1 := strings.Split(strRange, "=")
		slice2 := strings.Split(slice1[1], "-")
		start, _ := strconv.ParseInt(slice2[0], 10, 64)
		end, _ := strconv.ParseInt(slice2[1], 10, 64)
		io.Copy(w, bytes.NewBuffer(data[start:end+1]))
	})
	for i := 0; i < 3; i++ {
		math_rand.Seed(time.Now().UnixNano())
		rangeStart := math_rand.Intn(contentLength)
		rangeEnd := rangeStart + math_rand.Intn(contentLength-rangeStart)
		if rangeEnd == rangeStart || rangeStart >= contentLength-1 {
			continue
		}
		opt := &ObjectGetOptions{
			ResponseContentType: "text/html",
			Range:               fmt.Sprintf("bytes=%v-%v", rangeStart, rangeEnd),
		}
		resp, err := client.Object.Get(context.Background(), name, opt)
		if err != nil {
			t.Fatalf("Object.Get returned error: %v", err)
		}

		b, _ := ioutil.ReadAll(resp.Body)
		if bytes.Compare(b, data[rangeStart:rangeEnd+1]) != 0 {
			t.Errorf("Object.Get Failed")
		}
	}
}

func TestObjectService_GetToFile(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"
	data := make([]byte, 1024*1024*10)
	rand.Read(data)

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"response-content-type": "text/html",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "Range", "bytes=0-3")
		io.Copy(w, bytes.NewReader(data))
	})
	opt := &ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}
	filePath := "test.file" + time.Now().Format(time.RFC3339)
	_, err := client.Object.GetToFile(context.Background(), name, filePath, opt)
	if err != nil {
		t.Fatalf("Object.Get returned error: %v", err)
	}
	defer os.Remove(filePath)
	fd, err := os.Open(filePath)
	if err != nil {
	}
	defer fd.Close()
	bs, _ := ioutil.ReadAll(fd)
	if bytes.Compare(bs, data) != 0 {
		t.Errorf("Object.GetToFile data isn't consistent")
	}
}

func TestObjectService_GetPresignedURL(t *testing.T) {
	setup()
	defer teardown()

	exceptSign := "q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ*******&q-sign-time=1622702557;1622706157&q-key-time=1622702557;1622706157&q-header-list=&q-url-param-list=&q-signature=0f359fe9d29e7fa0c738ce6c8feaf4ed1e84f287"
	exceptURL := &url.URL{
		Scheme:   "http",
		Host:     client.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}

	c := context.Background()
	name := "test.jpg"
	ak := "QmFzZTY0IGlzIGEgZ*******"
	sk := "ZfbOA78asKUYBcXFrJD0a1I*******"
	startTime := time.Unix(int64(1622702557), 0)
	endTime := time.Unix(int64(1622706157), 0)
	opt := presignedURLTestingOptions{
		authTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}

	presignedURL, err := client.Object.GetPresignedURL(c, http.MethodPut, name, ak, sk, time.Hour, opt)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
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
			Listener:    &DefaultProgressListener{},
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

func TestObjectService_PostRestore(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"
	wantBody := "<RestoreRequest><Days>3</Days><CASJobParameters><Tier>Expedited</Tier></CASJobParameters></RestoreRequest>"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testHeader(t, r, "Content-Length", "106")
		//b, _ := ioutil.ReadAll(r.Body)
		//fmt.Printf("%s", string(b))
		testBody(t, r, wantBody)
	})

	opt := &ObjectRestoreOptions{
		Days: 3,
		Tier: &CASJobParameters{
			Tier: "Expedited",
		},
	}

	_, err := client.Object.PostRestore(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.PostRestore returned error: %v", err)
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

func TestObiectService_Read_and_Close(t *testing.T) {
	data := make([]byte, 1024*10)
	rand.Read(data)
	body := bytes.NewReader(data)
	r, _ := http.NewRequest(http.MethodGet, "test", body)

	drc := DiscardReadCloser{
		RC:      r.Body,
		Discard: 10,
	}

	res := make([]byte, 1024*10)
	readLen, err := drc.Read(res)
	if err != nil {
		t.Fatalf("Object.Read returned %v", err)
	}
	if readLen != 10230 {
		t.Fatalf("Object.Read returned %#v, excepted %#v", readLen, 10230)
	}
	if drc.Discard != 0 {
		t.Fatalf("Object.Read: drc.Discard = %v, excepted %v", drc.Discard, 0)
	}
	if !reflect.DeepEqual(res[:10230], data[10:]) {
		t.Fatalf("Object.Read: Wrong data!")
	}

	err = drc.Close()
	if err != nil {
		t.Fatal(err)
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

	wrongURL := "wrongURL"
	_, _, err := client.Object.Copy(context.Background(), "test.go.copy", wrongURL, nil)
	exceptedErr := errors.New(fmt.Sprintf("x-cos-copy-source format error: %s", wrongURL))
	if !reflect.DeepEqual(err, exceptedErr) {
		t.Fatalf("Object.Copy returned %#v, excepted %#v", err, exceptedErr)
	}

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

func TestObjectService_Download(t *testing.T) {
	setup()
	defer teardown()

	filePath := "rsp.file" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create tmp file failed")
	}
	defer os.Remove(filePath)
	// 源文件内容
	totalBytes := int64(1024*1024*9 + 123)
	b := make([]byte, totalBytes)
	_, err = rand.Read(b)
	newfile.Write(b)
	newfile.Close()
	tb := crc64.MakeTable(crc64.ECMA)
	localcrc := strconv.FormatUint(crc64.Update(0, tb, b), 10)

	retryMap := make(map[int64]int)
	mux.HandleFunc("/test.go.download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Add("Content-Length", strconv.FormatInt(totalBytes, 10))
			w.Header().Add("x-cos-hash-crc64ecma", localcrc)
			return
		}
		strRange := r.Header.Get("Range")
		slice1 := strings.Split(strRange, "=")
		slice2 := strings.Split(slice1[1], "-")
		start, _ := strconv.ParseInt(slice2[0], 10, 64)
		end, _ := strconv.ParseInt(slice2[1], 10, 64)
		if retryMap[start] == 0 {
			// 重试校验1
			retryMap[start]++
			w.WriteHeader(http.StatusGatewayTimeout)
		} else if retryMap[start] == 1 {
			// 重试检验2
			retryMap[start]++
			io.Copy(w, bytes.NewBuffer(b[start:end]))
		} else if retryMap[start] == 2 {
			// 重试检验3
			retryMap[start]++
			st := math_rand.Int63n(totalBytes - 1024*1024)
			et := st + end - start
			io.Copy(w, bytes.NewBuffer(b[st:et+1]))
		} else {
			io.Copy(w, bytes.NewBuffer(b[start:end+1]))
		}
	})

	opt := &MultiDownloadOptions{
		ThreadPoolSize: 3,
		PartSize:       1,
	}
	downPath := "down.file" + time.Now().Format(time.RFC3339)
	defer os.Remove(downPath)
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err == nil {
		// 长度不一致 Failed
		t.Fatalf("Object.Upload returned error: %v", err)
	}
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err == nil {
		// CRC不一致
		t.Fatalf("Object.Upload returned error: %v", err)
	}
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err != nil {
		// 正确
		t.Fatalf("Object.Upload returned error: %v", err)
	}
}

func TestObjectService_DownloadWithCheckPoint(t *testing.T) {
	setup()
	defer teardown()

	filePath := "rsp.file" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create tmp file failed")
	}
	defer os.Remove(filePath)
	// 源文件内容
	totalBytes := int64(1024*1024*9 + 123)
	partSize := 1024 * 1024
	b := make([]byte, totalBytes)
	_, err = rand.Read(b)
	newfile.Write(b)
	newfile.Close()
	tb := crc64.MakeTable(crc64.ECMA)
	localcrc := strconv.FormatUint(crc64.Update(0, tb, b), 10)

	oddok := false
	var oddcount, evencount int
	mux.HandleFunc("/test.go.download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Add("Content-Length", strconv.FormatInt(totalBytes, 10))
			w.Header().Add("x-cos-hash-crc64ecma", localcrc)
			return
		}
		strRange := r.Header.Get("Range")
		slice1 := strings.Split(strRange, "=")
		slice2 := strings.Split(slice1[1], "-")
		start, _ := strconv.ParseInt(slice2[0], 10, 64)
		end, _ := strconv.ParseInt(slice2[1], 10, 64)
		if (start/int64(partSize))%2 == 1 {
			if oddok {
				io.Copy(w, bytes.NewBuffer(b[start:end+1]))
			} else {
				// 数据校验失败, Download不会做重试
				io.Copy(w, bytes.NewBuffer(b[start:end]))
			}
			oddcount++
		} else {
			io.Copy(w, bytes.NewBuffer(b[start:end+1]))
			evencount++
		}
	})

	opt := &MultiDownloadOptions{
		ThreadPoolSize: 3,
		PartSize:       1,
		CheckPoint:     true,
	}
	downPath := "down.file" + time.Now().Format(time.RFC3339)
	defer os.Remove(downPath)
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err == nil {
		// 偶数块下载完成，奇数块下载失败
		t.Fatalf("Object.Download returned error: %v", err)
	}
	fd, err := os.Open(downPath)
	if err != nil {
		t.Fatalf("Object Download Open File Failed:%v", err)
	}
	offset := 0
	for i := 0; i < 10; i++ {
		bs, _ := ioutil.ReadAll(io.LimitReader(fd, int64(partSize)))
		offset += len(bs)
		if i%2 == 1 {
			bs[len(bs)-1] = b[offset-1]
		}
		if bytes.Compare(bs, b[i*partSize:offset]) != 0 {
			t.Fatalf("Compare Error, index:%v, len:%v, offset:%v", i, len(bs), offset)
		}
	}
	fd.Close()

	if oddcount != 5 || evencount != 5 {
		t.Fatalf("Object.Download failed, odd:%v, even:%v", oddcount, evencount)
	}
	// 设置奇数块OK
	oddok = true
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err != nil {
		// 下载成功
		t.Fatalf("Object.Download returned error: %v", err)
	}
	if oddcount != 10 || evencount != 5 {
		t.Fatalf("Object.Download failed, odd:%v, even:%v", oddcount, evencount)
	}
}
func TestObjectService_GetTagging(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"tagging": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<Tagging>
	<TagSet>
		<Tag>
			<Key>test_k2</Key>
			<Value>test_v2</Value>
		</Tag>
		<Tag>
			<Key>test_k3</Key>
			<Value>test_vv</Value>
		</Tag>
	</TagSet>
</Tagging>`)
	})

	res, _, err := client.Object.GetTagging(context.Background(), "test")
	if err != nil {
		t.Fatalf("Object.GetTagging returned error %v", err)
	}

	want := &ObjectGetTaggingResult{
		XMLName: xml.Name{Local: "Tagging"},
		TagSet: []ObjectTaggingTag{
			{"test_k2", "test_v2"},
			{"test_k3", "test_vv"},
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Object.GetTagging returned %+v, want %+v", res, want)
	}
}

func TestObjectService_PutTagging(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutTaggingOptions{
		TagSet: []ObjectTaggingTag{
			{
				Key:   "test_k2",
				Value: "test_v2",
			},
			{
				Key:   "test_k3",
				Value: "test_v3",
			},
		},
	}

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		v := new(ObjectPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		vs := values{
			"tagging": "",
		}
		testFormValues(t, r, vs)

		want := opt
		want.XMLName = xml.Name{Local: "Tagging"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.PutTagging request body: %+v, want %+v", v, want)
		}

	})

	_, err := client.Object.PutTagging(context.Background(), "test", opt)
	if err != nil {
		t.Fatalf("Object.PutTagging returned error: %v", err)
	}

}

func TestObjectService_DeleteTagging(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"tagging": "",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Object.DeleteTagging(context.Background(), "test")
	if err != nil {
		t.Fatalf("Object.DeleteTagging returned error: %v", err)
	}

}
