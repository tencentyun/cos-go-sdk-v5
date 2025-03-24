package cos

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
			Listener:            &DefaultProgressListener{},
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
		t.Errorf("Object.GetToFile open file failed: %v\n", err)
	}
	defer fd.Close()
	bs, _ := ioutil.ReadAll(fd)
	if bytes.Compare(bs, data) != 0 {
		t.Errorf("Object.GetToFile data isn't consistent")
	}
}

func TestObjectService_GetRetry(t *testing.T) {
	setup()
	defer teardown()
	u, _ := url.Parse(server.URL)
	client := NewClient(&BaseURL{u, u, u, u, u, u}, &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
		},
	})
	name := "test/hello.txt"
	contentLength := 1024 * 1024 * 10
	data := make([]byte, contentLength)
	index := int32(0)
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"response-content-type": "text/html",
		}
		atomic.AddInt32(&index, 1)
		if atomic.LoadInt32(&index)%3 != 0 {
			if atomic.LoadInt32(&index) > 6 {
				w.WriteHeader(500)
				return
			}
			time.Sleep(time.Second * 2)
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
	if index != 9 {
		t.Errorf("retry time error, retry count: %v\n", index)
	}
}

func TestObjectService_GetObjectURL(t *testing.T) {
	setup()
	defer teardown()

	name := "test"
	uri, _ := url.Parse("/" + encodeURIComponent(name, []byte{'/'}))
	want := client.BaseURL.BucketURL.ResolveReference(uri)

	res := client.Object.GetObjectURL("test")
	if res.String() != client.BaseURL.BucketURL.ResolveReference(uri).String() {
		t.Errorf("GetObjectURL failed, want: %v, return: %v", want, res)
	}
}

func TestObjectService_GetPresignedURL(t *testing.T) {
	setup()
	defer teardown()

	exceptSign := "q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ*******&q-sign-time=1622702557%3B1622706157&q-key-time=1622702557%3B1622706157&q-header-list=&q-url-param-list=&q-signature=820975b5a8eccce9455b94d4ebed14d66654bf3c"
	exceptURL := &url.URL{
		Scheme:   "http",
		Host:     client.BaseURL.BucketURL.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}

	c := context.Background()
	name := "test.jpg"
	ak := "QmFzZTY0IGlzIGEgZ*******"
	sk := "ZfbOA78asKUYBcXFrJD0a1I*******"
	startTime := time.Unix(int64(1622702557), 0)
	endTime := time.Unix(int64(1622706157), 0)
	opt := &presignedURLTestingOptions{
		authTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}

	presignedURL, err := client.Object.GetPresignedURL(c, http.MethodPut, name, ak, sk, time.Hour, opt, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
	}

	exceptSign = "test=params&sign=q-sign-algorithm%3Dsha1%26q-ak%3DQmFzZTY0IGlzIGEgZ*******%26q-sign-time%3D1622702557%3B1622706157%26q-key-time%3D1622702557%3B1622706157%26q-header-list%3D%26q-url-param-list%3Dtest%26q-signature%3D7757e84ed5f8953eafc30afcd2a5d1ad68e00d67"
	exceptURL = &url.URL{
		Scheme:   "http",
		Host:     client.BaseURL.BucketURL.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}
	opt1 := &PresignedURLOptions{
		Query:      &url.Values{},
		SignMerged: true,
		AuthTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}
	opt1.Query.Add("test", "params")
	presignedURL, err = client.Object.GetPresignedURL(c, http.MethodPut, name, ak, sk, time.Hour, opt1, false)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
	}

	_, err = client.Object.GetPresignedURL(c, http.MethodPut, "", ak, sk, time.Hour, opt1)
	if err == nil {
		t.Errorf("GetPresignedURL expect err is not null")
	}

	_, err = client.Object.GetPresignedURL(c, http.MethodPut, "/", ak, sk, time.Hour, opt1)
	if err != nil {
		t.Errorf("GetPresignedURL return err: %v", err)
	}

}

/*
	func testobjectservice_getpresignedurl_authtime(t *testing.t) {
		setup()
		defer teardown()

		exceptsign := "q-sign-algorithm=sha1&q-ak=qmfzzty0iglzigegz*******&q-sign-time=1622702557;1622706157&q-key-time=1622702557;1622706157&q-header-list=&q-url-param-list=&q-signature=0f359fe9d29e7fa0c738ce6c8feaf4ed1e84f287"
		excepturl := &url.url{
			scheme:   "http",
			host:     client.ba,
			path:     "/test.jpg",
			rawquery: exceptsign,
		}

		c := context.background()
		name := "test.jpg"
		ak := "qmfzzty0iglzigegz*******"
		sk := "zfboa78askuybcxfrjd0a1i*******"
		starttime := time.unix(int64(1622702557), 0)
		endtime := time.unix(int64(1622706157), 0)
		opt := &presignedurloptions{
			authtime: &authtime{
				signstarttime: starttime,
				signendtime:   endtime,
				keystarttime:  starttime,
				keyendtime:    endtime,
			},
		}

		presignedurl, err := client.object.getpresignedurl(c, http.methodput, name, ak, sk, time.hour, opt)
		if err != nil {
			t.fatal(err)
		}

		if !reflect.deepequal(excepturl, presignedurl) {
			t.fatalf("wrong presignedurl!")
		}
	}
*/
func TestObjectService_GetPresignedURL2(t *testing.T) {
	setup()
	defer teardown()

	u, _ := url.Parse(server.URL)
	client := NewClient(&BaseURL{u, u, u, u, u, u}, &http.Client{
		Transport: &AuthorizationTransport{
			SecretID:  "QmFzZTY0IGlzIGEgZ*******",
			SecretKey: "ZfbOA78asKUYBcXFrJD0a1I*******",
		},
	})

	exceptSign := "q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ*******&q-sign-time=1622702557%3B1622706157&q-key-time=1622702557%3B1622706157&q-header-list=&q-url-param-list=&q-signature=820975b5a8eccce9455b94d4ebed14d66654bf3c"
	exceptURL := &url.URL{
		Scheme:   "http",
		Host:     client.BaseURL.BucketURL.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}

	c := context.Background()
	name := "test.jpg"
	startTime := time.Unix(int64(1622702557), 0)
	endTime := time.Unix(int64(1622706157), 0)
	opt := &presignedURLTestingOptions{
		authTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}

	presignedURL, err := client.Object.GetPresignedURL2(c, http.MethodPut, name, time.Hour, opt, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
	}

	exceptSign = "test=params&sign=q-sign-algorithm%3Dsha1%26q-ak%3DQmFzZTY0IGlzIGEgZ*******%26q-sign-time%3D1622702557%3B1622706157%26q-key-time%3D1622702557%3B1622706157%26q-header-list%3D%26q-url-param-list%3Dtest%26q-signature%3D7757e84ed5f8953eafc30afcd2a5d1ad68e00d67"
	exceptURL = &url.URL{
		Scheme:   "http",
		Host:     client.BaseURL.BucketURL.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}
	opt1 := &PresignedURLOptions{
		Query:      &url.Values{},
		SignMerged: true,
		AuthTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}
	opt1.Query.Add("test", "params")
	presignedURL, err = client.Object.GetPresignedURL2(c, http.MethodPut, name, time.Hour, opt1, false)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
	}

	_, err = client.Object.GetPresignedURL2(c, http.MethodPut, "", time.Hour, opt1)
	if err == nil {
		t.Errorf("GetPresignedURL expect err is not null")
	}

	_, err = client.Object.GetPresignedURL2(c, http.MethodPut, "/", time.Hour, opt1)
	if err != nil {
		t.Errorf("GetPresignedURL return err: %v", err)
	}

}

func TestObjectService_GetPresignedURL3(t *testing.T) {
	setup()
	defer teardown()

	u, _ := url.Parse(server.URL)
	client := NewClient(&BaseURL{u, u, u, u, u, u}, &http.Client{
		Transport: &AuthorizationTransport{
			SecretID:  "QmFzZTY0IGlzIGEgZ*******",
			SecretKey: "ZfbOA78asKUYBcXFrJD0a1I*******",
		},
	})

	exceptSign := "q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ*******&q-sign-time=1622702557%3B1622706157&q-key-time=1622702557%3B1622706157&q-header-list=&q-url-param-list=&q-signature=820975b5a8eccce9455b94d4ebed14d66654bf3c"
	exceptURL := &url.URL{
		Scheme:   "http",
		Host:     client.BaseURL.BucketURL.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}

	c := context.Background()
	name := "test.jpg"
	startTime := time.Unix(int64(1622702557), 0)
	endTime := time.Unix(int64(1622706157), 0)
	opt := &presignedURLTestingOptions{
		authTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}

	presignedURL, err := client.Object.GetPresignedURL3(c, http.MethodPut, name, time.Hour, opt, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
	}

	exceptSign = "test=params&sign=q-sign-algorithm%3Dsha1%26q-ak%3DQmFzZTY0IGlzIGEgZ*******%26q-sign-time%3D1622702557%3B1622706157%26q-key-time%3D1622702557%3B1622706157%26q-header-list%3D%26q-url-param-list%3Dtest%26q-signature%3D7757e84ed5f8953eafc30afcd2a5d1ad68e00d67"
	exceptURL = &url.URL{
		Scheme:   "http",
		Host:     client.BaseURL.BucketURL.Host,
		Path:     "/test.jpg",
		RawQuery: exceptSign,
	}
	opt1 := &PresignedURLOptions{
		Query:      &url.Values{},
		SignMerged: true,
		AuthTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}
	opt1.Query.Add("test", "params")
	presignedURL, err = client.Object.GetPresignedURL3(c, http.MethodPut, name, time.Hour, opt1, false)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(exceptURL, presignedURL) {
		t.Fatalf("Wrong PreSignedURL!")
	}

	opt1.EncodeDelimiter = true
	_, err = client.Object.GetPresignedURL3(c, http.MethodPut, "", time.Hour, opt1)
	if err == nil {
		t.Errorf("GetPresignedURL expect err is not null")
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

func TestObjectService_Put2(t *testing.T) {
	setup()
	defer teardown()

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
	data := []byte("hello")

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

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
	})

	// reader is nil
	_, err := client.Object.Put(context.Background(), name, nil, opt)
	if err == nil || err.Error() != "reader is nil" {
		t.Fatalf("CI.Put returned error: %v", err)
	}

	// 没法获取 totalBytes
	r := &tmpOtherReader{
		Reader: bytes.NewReader(data),
	}
	_, err = client.Object.Put(context.Background(), name, r, opt)
	if err == nil || err.Error() != "can't get reader content length, unkown reader type" {
		t.Fatalf("Object.Put returned error: %v", err)
	}
	// 成功，通过ContentLength获取totalBytes
	opt.ContentLength = int64(len(data))
	_, err = client.Object.Put(context.Background(), name, r, opt)
	if err != nil {
		t.Fatalf("Object.Put returned error: %v", err)
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
		opt.Listener = &DefaultProgressListener{}
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

	var withVersion bool
	versionId := "versionid"
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		if withVersion {
			vs := values{
				"versionId": versionId,
			}
			testFormValues(t, r, vs)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		t.Fatalf("Object.Delete returned error: %v", err)
	}

	_, err = client.Object.Delete(context.Background(), "/test/../")
	if err != ObjectKeySimplifyCheckErr {
		t.Fatalf("Object.Delete expect error: %v", err)
	}

	opt := &ObjectDeleteOptions{
		VersionId: versionId,
	}
	_, err = client.Object.Delete(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Delete return error: %v", err)
	}
}

func TestObjectService_Head(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	var appendType bool
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "HEAD")
		testHeader(t, r, "If-Modified-Since", "Mon, 12 Jun 2017 05:36:19 GMT")
		if appendType {
			w.Header().Add("X-Cos-Object-Type", "appendable")
			w.Header().Add("Content-Length", "100")
		}
	})

	opt := &ObjectHeadOptions{
		IfModifiedSince: "Mon, 12 Jun 2017 05:36:19 GMT",
	}

	_, err := client.Object.Head(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Head returned error: %v", err)
	}

	// err
	_, err = client.Object.Head(context.Background(), name, opt, "id1", "id2")
	if err == nil || err.Error() != "wrong params" {
		t.Fatalf("Object.Head expect error: %v", err)
	}

	appendType = true
	resp, err := client.Object.Head(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Head return error: %v", err)
	}
	if resp.Header.Get("Content-Length") != "100" {
		t.Errorf("Object.Head header error: %v", resp.Header)
	}
}

func TestObjectService_IsExist(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "HEAD")
		w.WriteHeader(http.StatusNotFound)
	})

	isExisted, err := client.Object.IsExist(context.Background(), name)
	if err != nil {
		t.Fatalf("Object.Head returned error: %v", err)
	}
	if isExisted != false {
		t.Errorf("object IsExist failed")
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

	var withVersion bool
	versionId := "versionid"
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testHeader(t, r, "Content-Length", "106")
		testBody(t, r, wantBody)
		vs := values{
			"restore": "",
		}
		if withVersion {
			vs["versionId"] = versionId
		}
		testFormValues(t, r, vs)
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

	withVersion = true
	_, err = client.Object.PostRestore(context.Background(), name, opt, versionId)
	if err != nil {
		t.Fatalf("Object.PostRestore returned error: %v", err)
	}

	_, err = client.Object.PostRestore(context.Background(), name, opt, "id1", "id2")
	if err == nil || err.Error() != "wrong params" {
		t.Fatalf("Object.PostRestore expect error: %v", err)
	}
}

func TestObjectService_Append_Simple(t *testing.T) {
	setup()
	defer teardown()

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
	position := 0

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		vs := values{
			"append":   "",
			"position": "0",
		}
		testFormValues(t, r, vs)

		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Append request body: %#v, want %#v", v, want)
		}
		w.Header().Add("x-cos-content-sha1", hex.EncodeToString(calMD5Digest(b)))
		w.Header().Add("x-cos-next-append-position", strconv.FormatInt(int64(len(b)), 10))

	})

	r := bytes.NewReader([]byte("hello"))
	p, _, err := client.Object.Append(context.Background(), name, position, r, opt)
	if err != nil {
		t.Fatalf("Object.Append returned error: %v", err)
	}
	if p != len("hello") {
		t.Fatalf("Object.Append position error, want: %v, return: %v", len("hello"), p)
	}

	_, _, err = client.Object.Append(context.Background(), name, p, nil, opt)
	if err == nil || err.Error() != "reader is nil" {
		t.Errorf("Append expect error: %v", err)
	}
}

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

	var withVersion, returnErr bool
	versionId := "versionid"
	sourceURL := "test-1253846586.cos.ap-guangzhou.myqcloud.com/test.source"

	mux.HandleFunc("/test.go.copy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		if withVersion {
			source := r.Header.Get("X-Cos-Copy-Source")
			want := fmt.Sprintf("%v?versionId=%v", sourceURL, versionId)
			if source != want {
				t.Errorf("request copy-source: %v, want: %v", source, want)
			}
		}
		if returnErr {
			fmt.Fprint(w, `<Error>
		<Code>ErrorRequest</Code>
		<Message>Error Request</Message>
	</Error>`)
			return
		}
		fmt.Fprint(w, `<CopyObjectResult>
		<ETag>"098f6bcd4621d373cade4e832627b4f6"</ETag>
		<LastModified>2017-12-13T14:53:12</LastModified>
	</CopyObjectResult>`)
	})

	opt := &ObjectCopyOptions{
		&ObjectCopyHeaderOptions{
			ContentType: "application/xml",
		},
		&ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	wrongURL := "wrongURL"
	_, _, err := client.Object.Copy(context.Background(), "test.go.copy", wrongURL, opt)
	exceptedErr := errors.New(fmt.Sprintf("x-cos-copy-source format error: %s", wrongURL))
	if !reflect.DeepEqual(err, exceptedErr) {
		t.Fatalf("Object.Copy returned %#v, excepted %#v", err, exceptedErr)
	}

	ref, _, err := client.Object.Copy(context.Background(), "test.go.copy", sourceURL, opt)
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

	_, _, err = client.Object.Copy(context.Background(), "test.go.copy", "https://"+sourceURL, opt)
	if err == nil || err.Error() != "sourceURL format is invalid." {
		t.Errorf("Object.Copy returned failed: %v", err)
	}

	withVersion = true
	_, _, err = client.Object.Copy(context.Background(), "test.go.copy", sourceURL, opt, versionId)
	if err != nil {
		t.Errorf("Object.Copy returned failed: %v", err)
	}

	_, _, err = client.Object.Copy(context.Background(), "test.go.copy", sourceURL+"?versionId="+versionId, opt)
	if err != nil {
		t.Errorf("Object.Copy returned failed: %v", err)
	}

	// 响应错误
	returnErr = true
	_, _, err = client.Object.Copy(context.Background(), "test.go.copy", sourceURL+"?versionId="+versionId, opt)
	if err == nil {
		t.Errorf("Object.Copy expect error")
	}
	e, ok := err.(*ErrorResponse)
	if !ok || e.Code != "ErrorRequest" {
		t.Errorf("Object.Copy expect error: %v", err)
	}
}

func TestObjectService_Append(t *testing.T) {
	setup()
	defer teardown()
	size := 1111 * 1111 * 63
	b := make([]byte, size)
	p := int(math_rand.Int31n(int32(size)))
	var buf bytes.Buffer

	mux.HandleFunc("/test.append", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		bs, _ := ioutil.ReadAll(r.Body)
		buf.Write(bs)
		w.Header().Add("x-cos-content-sha1", hex.EncodeToString(calMD5Digest(bs)))
		w.Header().Add("x-cos-next-append-position", strconv.FormatInt(int64(buf.Len()), 10))
	})

	pos, _, err := client.Object.Append(context.Background(), "test.append", 0, bytes.NewReader(b[:p]), nil)
	if err != nil {
		t.Fatalf("Object.Append return error %v", err)
	}
	if pos != p {
		t.Fatalf("Object.Append pos error, returned:%v, wanted:%v", pos, p)
	}

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
			Listener:    &DefaultProgressListener{},
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}

	pos, _, err = client.Object.Append(context.Background(), "test.append", pos, bytes.NewReader(b[p:]), opt)
	if err != nil {
		t.Fatalf("Object.Append return error %v", err)
	}
	if pos != size {
		t.Fatalf("Object.Append pos error, returned:%v, wanted:%v", pos, size)
	}
	if bytes.Compare(b, buf.Bytes()) != 0 {
		t.Fatalf("Object.Append Compare failed")
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

	var mu sync.Mutex
	// 已上传内容, 10个分块
	rb := make([][]byte, 33)
	uploadid := "test-cos-multiupload-uploadid"
	partmap := make(map[int64]int)
	mux.HandleFunc("/test.go.upload", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
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
		OptIni: &InitiateMultipartUploadOptions{
			nil,
			&ObjectPutHeaderOptions{
				XCosMetaXXX: &http.Header{},
			},
		},
	}
	opt.OptIni.XCosMetaXXX.Add("x-cos-meta-test", "test")
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

	// 源文件内容
	totalBytes := int64(1024*1024*9 + 1230)
	b := make([]byte, totalBytes)
	_, err := rand.Read(b)
	tb := crc64.MakeTable(crc64.ECMA)
	localcrc := strconv.FormatUint(crc64.Update(0, tb, b), 10)

	var mu sync.Mutex
	retryMap := make(map[int64]int)
	mux.HandleFunc("/test.go.download", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
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
			// SDK 内部重试
			retryMap[start]++
			w.WriteHeader(http.StatusGatewayTimeout)
		} else if retryMap[start] == 1 {
			// SDK Download 做重试
			retryMap[start]++
			io.Copy(w, bytes.NewBuffer(b[start:end]))
		} else if retryMap[start] == 2 {
			// SDK Download 做重试
			retryMap[start]++
			st := start
			et := st + math_rand.Int63n(1024)
			io.Copy(w, bytes.NewBuffer(b[st:et+1]))
		} else {
			// SDK Download 成功
			io.Copy(w, bytes.NewBuffer(b[start:end+1]))
		}
	})

	opt := &MultiDownloadOptions{
		ThreadPoolSize: 3,
		PartSize:       1,
		Opt: &ObjectGetOptions{
			XCosSSECustomerAglo:   "AES256",
			XCosSSECustomerKey:    "MDEyMzQ1Njc4OUFCQ0RFRjAxMjM0NTY3ODlBQkNERUY=",
			XCosSSECustomerKeyMD5: "U5L61r7jcwdNvT7frmUG8g==",
		},
	}
	downPath := "down.file" + time.Now().Format(time.RFC3339)
	defer os.Remove(downPath)
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err != nil {
		t.Fatalf("Object.Upload returned error: %v", err)
	}
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err != nil {
		t.Fatalf("Object.Upload returned error: %v", err)
	}
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err != nil {
		t.Fatalf("Object.Upload returned error: %v", err)
	}

	totalBytes = 103
	name := "test/hello.txt"
	data := make([]byte, totalBytes)
	rand.Read(data)
	tb = crc64.MakeTable(crc64.ECMA)
	localcrc = strconv.FormatUint(crc64.Update(0, tb, data), 10)

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Add("Content-Length", strconv.FormatInt(totalBytes, 10))
			w.Header().Add("x-cos-hash-crc64ecma", localcrc)
			return
		}
		testMethod(t, r, "GET")
		w.Write(data)
	})
	fp := "test.file" + time.Now().Format(time.RFC3339)
	_, err = client.Object.Download(context.Background(), name, fp, opt)
	if err != nil {
		t.Fatalf("Object.Get returned error: %v", err)
	}
	defer os.Remove(fp)
	fd, err := os.Open(fp)
	if err != nil {
		t.Errorf("Object.GetToFile open file failed: %v\n", err)
	}
	defer fd.Close()
	bs, _ := ioutil.ReadAll(fd)
	if bytes.Compare(bs, data) != 0 {
		t.Errorf("Object.GetToFile data isn't consistent")
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
	var oddcount, evencount int32
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
				// 数据校验失败, Download做3次重试
				io.Copy(w, bytes.NewBuffer(b[start:end]))
			}
			atomic.AddInt32(&oddcount, 1)
		} else {
			io.Copy(w, bytes.NewBuffer(b[start:end+1]))
			atomic.AddInt32(&evencount, 1)
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

	if atomic.LoadInt32(&oddcount) != 15 || atomic.LoadInt32(&evencount) != 5 {
		t.Fatalf("Object.Download failed, odd:%v, even:%v", atomic.LoadInt32(&oddcount), atomic.LoadInt32(&evencount))
	}
	// 设置奇数块OK
	oddok = true
	_, err = client.Object.Download(context.Background(), "test.go.download", downPath, opt)
	if err != nil {
		// 下载成功
		t.Fatalf("Object.Download returned error: %v", err)
	}
	if atomic.LoadInt32(&oddcount) != 20 || atomic.LoadInt32(&evencount) != 5 {
		t.Fatalf("Object.Download failed, odd:%v, even:%v", atomic.LoadInt32(&oddcount), atomic.LoadInt32(&evencount))
	}
}

func TestObjectService_GetTagging(t *testing.T) {
	setup()
	defer teardown()

	var withVersion, withHeader bool
	versionId := "versionid"
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"tagging": "",
		}
		if withVersion {
			vs["versionId"] = versionId
		}
		testFormValues(t, r, vs)
		if withHeader {
			testHeader(t, r, "x-cos-meta-test", "test")
		}
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

	want := &ObjectGetTaggingResult{
		XMLName: xml.Name{Local: "Tagging"},
		TagSet: []ObjectTaggingTag{
			{"test_k2", "test_v2"},
			{"test_k3", "test_vv"},
		},
	}

	res, _, err := client.Object.GetTagging(context.Background(), "test", "id1", "id2", "id3")
	if err == nil || err.Error() != "wrong params" {
		t.Fatalf("Object.GetTagging expect error %v", err)
	}

	res, _, err = client.Object.GetTagging(context.Background(), "test")
	if err != nil {
		t.Fatalf("Object.GetTagging returned error %v", err)
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Object.GetTagging returned %+v, want %+v", res, want)
	}

	withVersion, withHeader = true, false
	res, _, err = client.Object.GetTagging(context.Background(), "test", versionId)
	if err != nil {
		t.Fatalf("Object.GetTagging returned error %v", err)
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Object.GetTagging returned %+v, want %+v", res, want)
	}

	withVersion, withHeader = false, true
	opt := &ObjectGetTaggingOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "test")
	res, _, err = client.Object.GetTagging(context.Background(), "test", opt)
	if err != nil {
		t.Fatalf("Object.GetTagging returned error %v", err)
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Object.GetTagging returned %+v, want %+v", res, want)
	}

	withVersion, withHeader = true, true
	res, _, err = client.Object.GetTagging(context.Background(), "test", versionId, opt)
	if err != nil {
		t.Fatalf("Object.GetTagging returned error %v", err)
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

	var withVersion bool
	versionId := "versionid"
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		v := new(ObjectPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		vs := values{
			"tagging": "",
		}
		if withVersion {
			vs["versionId"] = versionId
		}
		testFormValues(t, r, vs)

		want := opt
		want.XMLName = xml.Name{Local: "Tagging"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.PutTagging request body: %+v, want %+v", v, want)
		}

	})

	_, err := client.Object.PutTagging(context.Background(), "test", opt, "id", "id")
	if err == nil || err.Error() != "wrong params" {
		t.Errorf("Object.PutTagging expect error: %v", err)
	}

	_, err = client.Object.PutTagging(context.Background(), "test", opt)
	if err != nil {
		t.Fatalf("Object.PutTagging returned error: %v", err)
	}

	withVersion = true
	_, err = client.Object.PutTagging(context.Background(), "test", opt, versionId)
	if err != nil {
		t.Fatalf("Object.PutTagging returned error: %v", err)
	}

}

func TestObjectService_DeleteTagging(t *testing.T) {
	setup()
	defer teardown()

	var withVersion, withHeader bool
	versionId := "versionid"
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"tagging": "",
		}
		if withVersion {
			vs["versionId"] = versionId
		}
		testFormValues(t, r, vs)
		if withHeader {
			testHeader(t, r, "x-cos-meta-test", "test")
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Object.DeleteTagging(context.Background(), "/test/..")
	if err == nil {
		t.Errorf("DeleteTagging Expect error")
	}
	_, err = client.Object.DeleteTagging(context.Background(), "test")
	if err != nil {
		t.Fatalf("Object.DeleteTagging returned error: %v", err)
	}

	withVersion, withHeader = true, false
	_, err = client.Object.DeleteTagging(context.Background(), "test", versionId)
	if err != nil {
		t.Fatalf("Object.DeleteTagging returned error %v", err)
	}

	withVersion, withHeader = false, true
	opt := &ObjectGetTaggingOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "test")
	_, err = client.Object.DeleteTagging(context.Background(), "test", opt)
	if err != nil {
		t.Fatalf("Object.DeleteTagging returned error %v", err)
	}

	withVersion, withHeader = true, true
	_, err = client.Object.DeleteTagging(context.Background(), "test", versionId, opt)
	if err != nil {
		t.Fatalf("Object.DeleteTagging returned error %v", err)
	}
}

func TestObjectService_PutFetchTask(t *testing.T) {
	setup()
	defer teardown()

	opt := &PutFetchTaskOptions{
		Url:                "http://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com/exampleobject",
		Key:                "exampleobject",
		MD5:                "MD5",
		OnKeyExist:         "OnKeyExist",
		IgnoreSameKey:      true,
		SuccessCallbackUrl: "SuccessCallbackUrl",
		FailureCallbackUrl: "FailureCallbackUrl",
		XOptionHeader:      &http.Header{},
	}
	opt.XOptionHeader.Add("Content-Type", "application/json")
	opt.XOptionHeader.Add("Content-Type", "application/xml")
	opt.XOptionHeader.Add("Cache-Control", "max-age=10")
	opt.XOptionHeader.Add("Cache-Control", "max-stale=10")
	res := &PutFetchTaskResult{
		Code:      0,
		Message:   "SUCCESS",
		RequestId: "NjE0ZGMxMDhfMmZjMjNiMGFfNWY2N18yOTRjYWE=",
		Data: struct {
			TaskId string `json:"taskId,omitempty"`
		}{
			TaskId: "NjE0ZGMxMDhfMmZjMjNiMGFfNWY2N18yOTRjYWE=",
		},
	}
	mux.HandleFunc("/examplebucket-1250000000/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		opt.XOptionHeader.Set("Content-Type", "application/json")
		for k, v := range *opt.XOptionHeader {
			if k != "Content-Type" {
				if !reflect.DeepEqual(r.Header[k], v) {
					t.Errorf("Object.PutFetchTask request header: %+v, want %+v", r.Header[k], v)
				}
				continue
			}
			if r.Header.Get(k) != "application/json" || len(r.Header[k]) != 1 {
				t.Errorf("Object.PutFetchTask request header: %+v, want %+v", r.Header[k], v)
			}
		}
		v := new(PutFetchTaskOptions)
		json.NewDecoder(r.Body).Decode(v)
		want := opt
		v.XOptionHeader = opt.XOptionHeader
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.PutFetchTask request body: %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{
            "code":0,
            "message":"SUCCESS",
            "request_id":"NjE0ZGMxMDhfMmZjMjNiMGFfNWY2N18yOTRjYWE=",
            "data":{"taskid":"NjE0ZGMxMDhfMmZjMjNiMGFfNWY2N18yOTRjYWE="}
        }`)
	})

	r, _, err := client.Object.PutFetchTask(context.Background(), "examplebucket-1250000000", opt)
	if err != nil {
		t.Fatalf("Object.PutFetchTask returned error: %v", err)
	}
	if !reflect.DeepEqual(r, res) {
		t.Errorf("object.PutFetchTask res: %+v, want: %+v", r, res)
	}
}

func TestObjectService_GetFetchTask(t *testing.T) {
	setup()
	defer teardown()

	res := &GetFetchTaskResult{
		Code:      0,
		Message:   "SUCCESS",
		RequestId: "NjE0ZGNiMDVfMmZjMjNiMGFfNWY2N18yOTRjYWM=",
		Data: struct {
			Code    string `json:"code,omitempty"`
			Message string `json:"msg,omitempty"`
			Percent int    `json:"percent,omitempty"`
			Status  string `json:"status,omitempty"`
		}{
			Code:    "Forbidden",
			Message: "The specified download can not be allowed.",
			Percent: 0,
			Status:  "TASK_FAILED",
		},
	}
	mux.HandleFunc("/examplebucket-1250000000/NjE0ZGMxMDhfMmZjMjNiMGFfNWY2N18yOTRjYWE=", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
            "code":0,
            "message":"SUCCESS",
            "request_id":"NjE0ZGNiMDVfMmZjMjNiMGFfNWY2N18yOTRjYWM=",
            "data": {
                "code":"Forbidden",
                "msg":"The specified download can not be allowed.",
                "percent":0,
                "status":"TASK_FAILED"
            }
        }`)
	})

	r, _, err := client.Object.GetFetchTask(context.Background(), "examplebucket-1250000000", "NjE0ZGMxMDhfMmZjMjNiMGFfNWY2N18yOTRjYWE=")
	if err != nil {
		t.Fatalf("Object.GetFetchTask returned error: %v", err)
	}
	if !reflect.DeepEqual(r, res) {
		t.Errorf("object.GetFetchTask res: %+v, want: %+v", r, res)
	}
}

func TestObjectService_Select(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectSelectOptions{
		Expression:     "Select * from COSObject",
		ExpressionType: "SQL",
		InputSerialization: &SelectInputSerialization{
			CSV: &CSVInputSerialization{
				FileHeaderInfo: "IGNORE",
			},
		},
		OutputSerialization: &SelectOutputSerialization{
			CSV: &CSVOutputSerialization{
				RecordDelimiter: "\n",
			},
		},
		RequestProgress: "TRUE",
	}

	// send a frame
	sendAFrame := func(header map[string]string, payload []byte) []byte {
		buf := bytes.NewBuffer([]byte{})

		var totalFrameLength, totalHeaderLength int
		for k, v := range header {
			totalHeaderLength += len(k) + len(v) + 4
		}
		totalFrameLength = 12 + totalHeaderLength + len(payload) + 4

		// 预响应
		binary.Write(buf, binary.BigEndian, int32(totalFrameLength))
		binary.Write(buf, binary.BigEndian, int32(totalHeaderLength))
		c := crc32.ChecksumIEEE(buf.Bytes())
		binary.Write(buf, binary.BigEndian, c)

		// headers
		for k, v := range header {
			binary.Write(buf, binary.BigEndian, int8(len(k)))
			buf.Write([]byte(k))
			binary.Write(buf, binary.BigEndian, int8(7))
			binary.Write(buf, binary.BigEndian, int16(len(v)))
			buf.Write([]byte(v))
		}

		// payload
		buf.Write(payload)

		// crc
		c32 := crc32.ChecksumIEEE(buf.Bytes())
		binary.Write(buf, binary.BigEndian, c32)

		return buf.Bytes()
	}

	dataSize := 1222 * 3
	data := make([]byte, dataSize)
	rand.Read(data)
	result := ObjectSelectResult{
		ProgressFrame: ProgressFrame{
			XMLName:        xml.Name{Local: "Progress"},
			BytesScanned:   dataSize,
			BytesProcessed: dataSize,
			BytesReturned:  dataSize,
		},
		StatsFrame: StatsFrame{
			XMLName:        xml.Name{Local: "Stats"},
			BytesScanned:   dataSize,
			BytesProcessed: dataSize,
			BytesReturned:  dataSize,
		},
		ErrorFrame: &ErrorFrame{
			Code:    "InternalError",
			Message: "We encounted an internal error, Please try again",
		},
	}

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// 检查请求
		testMethod(t, r, "POST")
		vs := values{
			"select":      "",
			"select-type": "2",
		}
		testFormValues(t, r, vs)
		want := opt
		want.XMLName = xml.Name{Local: "SelectRequest"}
		v := new(ObjectSelectOptions)
		xml.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Select request body: %+v, want %+v", v, want)
		}

		// Continue Message
		w.Write(sendAFrame(map[string]string{
			":message-type": "event",
			":event-type":   "Cont"}, []byte{}))

		// Records Message
		w.Write(sendAFrame(map[string]string{
			":message-type": "event",
			":event-type":   "Records"}, data))

		// Progress Message
		pframe, _ := xml.Marshal(result.ProgressFrame)
		w.Write(sendAFrame(map[string]string{
			":message-type": "event",
			":event-type":   "Progress"}, pframe))

		// Stat Message
		sframe, _ := xml.Marshal(result.StatsFrame)
		w.Write(sendAFrame(map[string]string{
			":message-type": "event",
			":event-type":   "Stats"}, sframe))

		// End Message
		w.Write(sendAFrame(map[string]string{
			":message-type": "event",
			":event-type":   "End"}, []byte{}))
	})

	mux.HandleFunc("/test_error", func(w http.ResponseWriter, r *http.Request) {
		// Records Message
		w.Write(sendAFrame(map[string]string{
			":message-type": "event",
			":event-type":   "Records"}, data))

		// Error Message
		w.Write(sendAFrame(map[string]string{
			":message-type":  "error",
			":error-code":    result.ErrorFrame.Code,
			":error-message": result.ErrorFrame.Message}, []byte{}))
	})

	// 测试正常情况
	filePath := "test.file" + time.Now().Format(time.RFC3339)
	res, err := client.Object.SelectToFile(context.Background(), "test", filePath, opt)
	if err != nil {
		t.Errorf("Object.Select failed: %v\n", err)
	}
	defer os.Remove(filePath)
	fd, err := os.Open(filePath)
	if err != nil {
		t.Errorf("Object.Select open file failed: %v\n", err)
	}
	defer fd.Close()
	bs, err := ioutil.ReadAll(fd)
	if err != nil {
		t.Errorf("Object.Select read failed: %v\n", err)
	}
	if bytes.Compare(bs, data) != 0 {
		t.Errorf("Object.Select compare failed\n")
	}
	if !reflect.DeepEqual(result.StatsFrame, res.Frame.StatsFrame) {
		t.Errorf("Object.Select stat frame failed, return: %+v, want: %+v\n", res.Frame.StatsFrame, result.StatsFrame)
	}
	if !reflect.DeepEqual(result.ProgressFrame, res.Frame.ProgressFrame) {
		t.Errorf("Object.Select progress frame failed, return: %+v, want: %+v\n", res.Frame.ProgressFrame, result.ProgressFrame)
	}

	// 测试错误情况
	resp, err := client.Object.Select(context.Background(), "test_error", opt)
	if err != nil {
		t.Errorf("Object.Select failed: %v\n", err)
	}
	_, err = ioutil.ReadAll(resp)
	ef, ok := err.(*ErrorFrame)
	if !ok {
		t.Errorf("Object.Select error is not ErrorFrame, %v", err)
	}
	if !reflect.DeepEqual(ef, result.ErrorFrame) {
		t.Errorf("Object.Select error frame failed, return: %+v, want: %+v\n", res.Frame.ErrorFrame, result.ErrorFrame)
	}
}

func TestObjectService_GetSignature(t *testing.T) {
	setup()
	defer teardown()

	timekey := "q-key-time="
	secretID := "ak"
	secretKey := "sk"
	name := "exampleobject"
	sign := client.Object.GetSignature(context.Background(), http.MethodGet, name, secretID, secretKey, time.Hour, nil)
	if sign == "" || len(sign) <= len(timekey) {
		t.Errorf("GetSignature sign is invalid: %v", sign)
		return
	}
	pos := strings.Index(sign, timekey) + len(timekey)
	st_et := strings.SplitN(sign[pos:pos+strings.Index(sign[pos:], "&")], ";", 2)
	if len(st_et) != 2 {
		t.Errorf("GetSignature sign is invalid: %v", sign)
		return
	}
	startTime, _ := strconv.ParseInt(st_et[0], 10, 64)
	endTime, _ := strconv.ParseInt(st_et[1], 10, 64)
	authTime := &AuthTime{
		SignStartTime: time.Unix(startTime, 0),
		SignEndTime:   time.Unix(endTime, 0),
		KeyStartTime:  time.Unix(startTime, 0),
		KeyEndTime:    time.Unix(endTime, 0),
	}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String()+"/"+name, nil)
	wanted := newAuthorization(secretID, secretKey, req, authTime, true)

	if sign != wanted {
		t.Errorf("GetSignature error, return: %+v, want: %+v\n", sign, wanted)
	}
}

func TestObjectService_GetSignature2(t *testing.T) {
	setup()
	defer teardown()

	secretID := "ak"
	secretKey := "sk"
	name := "exampleobject"
	startTime := time.Unix(int64(1622702557), 0)
	endTime := time.Unix(int64(1622706157), 0)
	opt := &PresignedURLOptions{
		Query: &url.Values{},
		AuthTime: &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		},
	}
	opt.Query.Add("key", "value")
	client.Object.GetSignature(context.Background(), http.MethodGet, name, secretID, secretKey, time.Hour, opt, true)
}

func TestObjectService_getResumableUploadID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"prefix":        "Object",
			"encoding-type": "url",
			"uploads":       "",
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
</ListMultipartUploadsResult>`)
	})

	id, err := client.Object.getResumableUploadID(context.Background(), "Object")
	if err != nil {
		t.Errorf("getResumableUploadID failed: %v", err)
	}
	if id != "1484727158f2b8034e5407d18cbf28e84f754b791ecab607d25a2e52de9fee641e5f60707c" {
		t.Errorf("getResumableUploadID failed: %v", id)
	}
}

func TestObjectService_checkUploadedParts(t *testing.T) {
	setup()
	defer teardown()

	name := "test/hello.txt"
	uploadID := "149795194893578fd83aceef3a88f708f81f00e879fda5ea8a80bf15aba52746d42d512387"
	bs := make([]byte, 2048)
	rand.Read(bs)
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodGet)
		vs := values{
			"uploadId":      uploadID,
			"encoding-type": "url",
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
		<ETag>&quot;`+hex.EncodeToString(calMD5Digest(bs[:1024]))+`&quot;</ETag>
		<Size>6291456</Size>
	</Part>
	<Part>
		<PartNumber>2</PartNumber>
		<LastModified>2017-06-20T09:45:50.000Z</LastModified>
		<ETag>&quot;`+hex.EncodeToString(calMD5Digest(bs[1024:]))+`&quot;</ETag>
		<Size>6391456</Size>
	</Part>
	<StorageClass>Standard</StorageClass>
	<MaxParts>1000</MaxParts>
	<IsTruncated>false</IsTruncated>
	</ListPartsResult>`)
	})
	filePath := "test.file" + time.Now().Format(time.RFC3339)
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	defer os.Remove(filePath)
	_, err = fd.Write(bs)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	fd.Close()

	chunks := []Chunk{
		{
			OffSet: 0,
			Size:   1024,
		},
		{
			OffSet: 1024,
			Size:   1024,
		},
	}
	err = client.Object.checkUploadedParts(context.Background(),
		name, uploadID, filePath, chunks, 2)
	if err != nil {
		t.Fatalf("Object.checkUploadedParts returned error: %v", err)
	}
	chunks = []Chunk{
		{
			OffSet: 0,
			Size:   1024,
		},
		{
			OffSet: 1024,
			Size:   1023, // 特意校验失败
		},
	}
	err = client.Object.checkUploadedParts(context.Background(),
		name, uploadID, filePath, chunks, 2)
	if err == nil {
		t.Fatalf("Object.checkUploadedParts should return err: %v", err)
	}
}

func TestObjectService_PutSymlink(t *testing.T) {
	setup()
	defer teardown()

	name := "symlink"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"symlink": "",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "x-cos-symlink-target", "target")
	})

	opt := &ObjectPutSymlinkOptions{
		SymlinkTarget: "target",
	}
	_, err := client.Object.PutSymlink(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.PutSymlink  returned error %v", err)
	}
}

func TestObjectService_GetSymlink(t *testing.T) {
	setup()
	defer teardown()

	name := "symlink"
	want := "target"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"symlink": "",
		}
		testFormValues(t, r, vs)
		w.Header().Set("x-cos-symlink-target", want)
	})

	res, _, err := client.Object.GetSymlink(context.Background(), name, nil)
	if err != nil {
		t.Fatalf("Object.GetSymlink returned error %v", err)
	}
	if res != want {
		t.Fatalf("Object.GetSymlink, target is invalid, return: %v, want: %v", res, want)
	}
}

func TestObjectService_PutRetry(t *testing.T) {
	setup()
	defer teardown()
	name := "test/retry"
	data := make([]byte, 1024*1024*3)
	_, err := rand.Read(data)
	tb := crc64.MakeTable(crc64.ECMA)
	realcrc := crc64.Update(0, tb, data)

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			Listener: &DefaultProgressListener{},
		},
	}

	nr, count := 0, 3
	mux.HandleFunc("/test/retry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		bs, _ := ioutil.ReadAll(r.Body)
		crc := crc64.Update(0, tb, bs)
		if !reflect.DeepEqual(crc, realcrc) {
			t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
		}
		nr++
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
		if nr < count {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	_, err = client.Object.Put(context.Background(), name, bytes.NewReader(data), opt)
	if err != nil || nr != count {
		t.Errorf("Object.Put failed: %v", err)
	}
	nr, count = 0, 3
	_, err = client.Object.Put(context.Background(), name, strings.NewReader(string(data)), opt)
	if err != nil || nr != count {
		t.Errorf("Object.Put failed: %v", err)
	}
	// 非io.Seeker不做重试
	nr, count = 0, 3
	_, err = client.Object.Put(context.Background(), name, bytes.NewBuffer(data), opt)
	if err == nil || nr != 1 {
		t.Errorf("Object.Put failed: %v", err)
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

	nr, count = 0, 3
	_, err = client.Object.PutFromFile(context.Background(), name, filePath, opt)
	if err != nil || nr != count {
		t.Errorf("PutFromFile failed: %v", err)
	}
}

func TestObjectService_UploadRetry(t *testing.T) {
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

	var mu sync.Mutex
	// 已上传内容, 10个分块
	rb := make([][]byte, 33)
	uploadid := "test-cos-multiupload-uploadid"
	partmap := make(map[int64]int)
	mux.HandleFunc("/test.go.upload", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.Method == http.MethodPut { // 分块上传
			r.ParseForm()
			part, _ := strconv.ParseInt(r.Form.Get("partNumber"), 10, 64)
			if partmap[part] == 0 {
				// 重试检验1
				partmap[part]++
				ioutil.ReadAll(r.Body)
				w.WriteHeader(http.StatusInternalServerError)
			} else if partmap[part] == 1 {
				// 重试校验2
				partmap[part]++
				ioutil.ReadAll(r.Body)
				w.WriteHeader(http.StatusBadRequest)
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
		OptIni: &InitiateMultipartUploadOptions{
			nil,
			&ObjectPutHeaderOptions{
				XCosMetaXXX: &http.Header{},
			},
		},
	}
	_, _, err = client.Object.Upload(context.Background(), "test.go.upload", filePath, opt)
	if err != nil {
		t.Fatalf("Object.Upload returned error: %v", err)
	}
}

func TestObjectKeyErr(t *testing.T) {
	names := []string{"../", "/../", "abc/..", "/abc/../123////..", "/bcd/../abc/../", "/bcd/./abc/../123/../..", "////./123/////..", "/./123//////bcd/../.."}

	for _, name := range names {
		_, err := client.Object.Get(context.Background(), name, nil)
		if err != ObjectKeySimplifyCheckErr {
			t.Fatalf("Get Err, want: %v, return: %v, name: %v", ObjectKeySimplifyCheckErr, err, name)
		}
		_, err = client.Object.GetToFile(context.Background(), name, "", nil)
		if err != ObjectKeySimplifyCheckErr {
			t.Fatalf("GetToFile Err, want: %v, return: %v, name: %v", ObjectKeySimplifyCheckErr, err, name)
		}
		_, err = client.Object.Download(context.Background(), name, "", nil)
		if err != ObjectKeySimplifyCheckErr {
			t.Fatalf("Download Err, want: %v, return: %v, name: %v", ObjectKeySimplifyCheckErr, err, name)
		}
	}
}

func TestObjectService_PutFromURL(t *testing.T) {
	setup()
	defer teardown()

	source := "source"
	dest := "dest"
	partSize := 8

	data := make([]byte, 1024*1024*33+133)
	_, err := rand.Read(data)
	tb := crc64.MakeTable(crc64.ECMA)
	realcrc := crc64.Update(0, tb, data)

	var finalcrc uint64
	var gtable *crc64.Table
	var partNumber int64
	var testResourceFailed, testInitFailed, testPartFailed, testCompleteFailed bool
	initTest := func() {
		finalcrc = 0
		gtable = crc64.MakeTable(crc64.ECMA)
		partNumber = 0
		testResourceFailed, testInitFailed, testPartFailed, testCompleteFailed = false, false, false, false
	}
	// 源文件内容
	mux.HandleFunc("/"+source, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		if testResourceFailed {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		n, err := io.Copy(w, bytes.NewReader(data))
		if err != nil && err != io.EOF || n != int64(len(data)) {
			t.Errorf("io.copy failed: %v", err)
		}
	})

	// 全局crc
	mux.HandleFunc("/"+dest, func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		switch r.Method {
		case http.MethodPost:
			// init
			init := url.Values{}
			init.Set("uploads", "")
			if reflect.DeepEqual(r.Form, init) {
				if testInitFailed {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(w, `<InitiateMultipartUploadResult>
                    <Bucket></Bucket>
                    <Key>dest</Key>
                    <UploadId>putfromurl_uploadid</UploadId>
                </InitiateMultipartUploadResult>`)
				return
			}
			// complete
			complete := url.Values{}
			complete.Set("uploadId", "putfromurl_uploadid")
			if !reflect.DeepEqual(r.Form, complete) {
				t.Errorf("complete check query failed, get: %v, want %v", r.Form, complete)
			}
			// 校验crc
			if realcrc != finalcrc {
				t.Errorf("crc64ecma mismatch, want: %v, return: %v", realcrc, finalcrc)
			}
			if testCompleteFailed {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(realcrc, 10))
			fmt.Fprintf(w, `<CompleteMultipartUploadResult>
            	<Location>/dest</Location>
                <Bucket></Bucket>
                <Key>dest</Key>
                <ETag>&quot;etag&quot;</ETag>
            </CompleteMultipartUploadResult>`)
		case http.MethodPut:
			// 分块
			partNumber++
			vs := values{
				"uploadId":   "putfromurl_uploadid",
				"partNumber": strconv.FormatInt(partNumber, 10),
			}
			testFormValues(t, r, vs)
			bs, _ := io.ReadAll(r.Body)
			finalcrc = crc64.Update(finalcrc, gtable, bs)
			st, ed := partSize*1024*1024*int(partNumber-1), len(bs)
			// 比较数据
			if bytes.Compare(bs, data[st:st+ed]) != 0 {
				t.Errorf("data mismatch: %v", partNumber)
			}
			tb := crc64.MakeTable(crc64.ECMA)
			partcrc := crc64.Update(0, tb, bs)
			if testPartFailed {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(partcrc, 10))
			w.Header().Add("Etag", "\""+hex.EncodeToString(calMD5Digest(bs))+"\"")
		case http.MethodDelete:
			vs := values{
				"uploadId": "putfromurl_uploadid",
			}
			testFormValues(t, r, vs)
		}
	})
	downloadUrl := client.BaseURL.BucketURL.String() + "/" + source

	opt := &ObjectPutFromURLOptions{
		PartSize:  partSize,
		QueueSize: 1,
	}
	initTest()
	_, _, err = client.Object.PutFromURL(context.Background(), dest, downloadUrl, opt)
	if err != nil {
		t.Errorf("Object.PutFromURL returned error: %v", err)
	}
	initTest()
	_, _, err = client.Object.PutFromURL(context.Background(), dest, downloadUrl, nil)
	if err != nil {
		t.Errorf("Object.PutFromURL returned error: %v", err)
	}

	initTest()
	testInitFailed = true
	_, _, err = client.Object.PutFromURL(context.Background(), dest, downloadUrl, nil)
	if err == nil {
		t.Errorf("Object.PutFromURL expect error")
	}
	initTest()
	testResourceFailed = true
	_, _, err = client.Object.PutFromURL(context.Background(), dest, downloadUrl, nil)
	if err == nil {
		t.Errorf("Object.PutFromURL expect error")
	}
	initTest()
	testPartFailed = true
	_, _, err = client.Object.PutFromURL(context.Background(), dest, downloadUrl, nil)
	if err == nil {
		t.Errorf("Object.PutFromURL expect error")
	}
	initTest()
	testCompleteFailed = true
	_, _, err = client.Object.PutFromURL(context.Background(), dest, downloadUrl, nil)
	if err == nil {
		t.Errorf("Object.PutFromURL expect error")
	}
}
