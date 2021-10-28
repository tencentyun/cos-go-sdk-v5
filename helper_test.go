package cos

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func Test_calSHA1Digest(t *testing.T) {
	want := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	got := fmt.Sprintf("%x", calSHA1Digest([]byte("test")))
	if got != want {

		t.Errorf("calSHA1Digest request sha1: %+v, want %+v", got, want)
	}
}

func Test_calMD5Digest(t *testing.T) {
	want := "098f6bcd4621d373cade4e832627b4f6"
	got := fmt.Sprintf("%x", calMD5Digest([]byte("test")))
	if got != want {

		t.Errorf("calMD5Digest request md5: %+v, want %+v", got, want)
	}
}

func Test_cloneHeader(t *testing.T) {
	ori := http.Header{}
	opt := &ori
	opt.Add("TestHeader1", "h1")
	opt.Add("TestHeader1", "h2")
	res := cloneHeader(opt)
	if !reflect.DeepEqual(res, opt) {
		t.Errorf("cloneHeader, returned:%+v, want:%+v", res, opt)
	}
	if !reflect.DeepEqual(ori, *opt) {
		t.Errorf("cloneHeader, returned:%+v, want:%+v", *opt, ori)
	}
	res.Add("cloneHeader1", "c1")
	res.Add("cloneHeader2", "c2")
	if v := opt.Get("cloneHeader1"); v != "" {
		t.Errorf("cloneHeader, returned:%+v, want:%+v", res, opt)
	}
	if v := opt.Get("cloneHeader2"); v != "" {
		t.Errorf("cloneHeader, returned:%+v, want:%+v", res, opt)
	}
	opt = &http.Header{}
	res = cloneHeader(opt)
	if !reflect.DeepEqual(res, opt) {
		t.Errorf("cloneHeader, returned:%+v, want:%+v", res, opt)
	}
}

func Test_CloneCompleteMultipartUploadOptions(t *testing.T) {
	ori := CompleteMultipartUploadOptions{
		XMLName: xml.Name{Local: "CompleteMultipartUploadResult"},
		Parts: []Object{
			{
				Key:  "Key1",
				ETag: "Etag1",
			},
			{
				Key:  "Key2",
				ETag: "Etag2",
			},
		},
		XOptionHeader: &http.Header{},
	}
	ori.XOptionHeader.Add("Test", "value")
	opt := &ori
	res := CloneCompleteMultipartUploadOptions(opt)
	if !reflect.DeepEqual(res, opt) {
		t.Errorf("CloneCompleteMultipartUploadOptions, returned:%+v,want:%+v", res, opt)
	}
	if !reflect.DeepEqual(ori, *opt) {
		t.Errorf("CloneCompleteMultipartUploadOptions, returned:%+v,want:%+v", *opt, ori)
	}
	res.XOptionHeader.Add("TestClone", "value")
	if v := opt.XOptionHeader.Get("TestClone"); v != "" {
		t.Errorf("CloneCompleteMultipartUploadOptions, returned:%+v,want:%+v", res, opt)
	}
	opt = &CompleteMultipartUploadOptions{}
	res = CloneCompleteMultipartUploadOptions(opt)
	if !reflect.DeepEqual(res, opt) {
		t.Errorf("CloneCompleteMultipartUploadOptions, returned:%+v,want:%+v", res, opt)
	}
	res.Parts = append(res.Parts, Object{Key: "K", ETag: "T"})
	if len(opt.Parts) > 0 {
		t.Errorf("CloneCompleteMultipartUploadOptions Failed")
	}
	if reflect.DeepEqual(res, opt) {
		t.Errorf("CloneCompleteMultipartUploadOptions, returned:%+v,want:%+v", res, opt)
	}
}

func Test_CopyOptionsToMulti(t *testing.T) {
	opt := &ObjectCopyOptions{
		&ObjectCopyHeaderOptions{
			CacheControl:    "max-age=1",
			ContentEncoding: "gzip",
			ContentType:     "text/html",
		},
		nil,
	}
	mul := CopyOptionsToMulti(opt)
	if opt.ContentType != mul.ContentType {
		t.Errorf("CopyOptionsToMulti, returned:%+v,want:%+v", mul, opt)
	}
	if opt.CacheControl != mul.CacheControl {
		t.Errorf("CopyOptionsToMulti, returned:%+v,want:%+v", mul, opt)
	}
	if opt.ContentEncoding != mul.ContentEncoding {
		t.Errorf("CopyOptionsToMulti, returned:%+v,want:%+v", mul, opt)
	}
}

func Test_CloneInitiateMultipartUploadOptions(t *testing.T) {
	opt := &InitiateMultipartUploadOptions{
		&ACLHeaderOptions{},
		&ObjectPutHeaderOptions{
			CacheControl:    "max-age=1",
			ContentEncoding: "gzip",
			ContentType:     "text/html",
		},
	}
	res := CloneInitiateMultipartUploadOptions(opt)
	if !reflect.DeepEqual(opt, res) {
		t.Errorf("CloneInitiateMultipartUploadOptions, returned:%+v,want:%+v", res, opt)
	}
}
