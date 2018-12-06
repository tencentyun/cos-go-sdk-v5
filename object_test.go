package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
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

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Put request body: %#v, want %#v", v, want)
		}
	})

	r := bytes.NewReader([]byte("hello"))
	_, err := client.Object.Put(context.Background(), name, r, opt)
	if err != nil {
		t.Fatalf("Object.Put returned error: %v", err)
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
		Origin: "www.qq.com",
		AccessControlRequestMethod: "PUT",
	}

	_, err := client.Object.Options(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Options returned error: %v", err)
	}

}

func TestObjectService_Append(t *testing.T) {
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
	})

	r := bytes.NewReader([]byte("hello"))
	_, err := client.Object.Append(context.Background(), name, position, r, opt)
	if err != nil {
		t.Fatalf("Object.Append returned error: %v", err)
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
