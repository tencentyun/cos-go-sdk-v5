package cos

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// 非IsLenReader
type tmpOtherReader struct {
	Reader io.Reader
}

func (this *tmpOtherReader) Read(p []byte) (int, error) {
	return this.Reader.Read(p)
}

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

func Test_CloneObjectGetOptions(t *testing.T) {
	opt := ObjectGetOptions{
		Range: "bytes=1-100",
	}
	res := CloneObjectGetOptions(&opt)
	if opt.Range != res.Range {
		t.Errorf("CloneObjectGetOptions failed")
	}
	ro, _ := GetRangeOptions(&opt)
	if FormatRangeOptions(ro) != "bytes=1-100" {
		t.Errorf("FormatRangeOptions failed")
	}
	ro.HasStart = true
	ro.HasEnd = false
	if FormatRangeOptions(ro) != "bytes=1-" {
		t.Errorf("FormatRangeOptions failed")
	}
	ro.HasStart = false
	ro.HasEnd = true
	if FormatRangeOptions(ro) != "bytes=-100" {
		t.Errorf("FormatRangeOptions failed")
	}
	ro.HasStart = false
	ro.HasEnd = false
	if FormatRangeOptions(ro) != "" {
		t.Errorf("FormatRangeOptions failed")
	}
}

func Test_progress(t *testing.T) {
	listener := &DefaultProgressListener{}
	listener.ProgressChangedCallback(&ProgressEvent{
		EventType:  ProgressStartedEvent,
		TotalBytes: 1,
	})
	listener.ProgressChangedCallback(&ProgressEvent{
		EventType:  ProgressDataEvent,
		TotalBytes: 1,
	})
	listener.ProgressChangedCallback(&ProgressEvent{
		EventType:  ProgressCompletedEvent,
		TotalBytes: 1,
	})
	listener.ProgressChangedCallback(&ProgressEvent{
		EventType:  ProgressFailedEvent,
		TotalBytes: 1,
	})
	listener.ProgressChangedCallback(&ProgressEvent{
		EventType:  -1,
		TotalBytes: 1,
	})
}

func Test_UnmarshalInitMultiUploadResult_Error(t *testing.T) {
	data := []byte(`
<InitiateMultipartUploadResult>
    <Key>test</Key>
    <UploadId>149795166761115ef06e259b2fceb8ff34bf7dd840883d26a0f90243562dd398efa41718db</UploadId>
</InitiateMultipartUploadResult>`)
	var res InitiateMultipartUploadResult
	err := UnmarshalInitMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}

	data = []byte(`
<InitiateMultipartUploadResult>
    <Bucket>test-1253846586</Bucket>
    <UploadId>149795166761115ef06e259b2fceb8ff34bf7dd840883d26a0f90243562dd398efa41718db</UploadId>
</InitiateMultipartUploadResult>`)
	err = UnmarshalInitMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}

	data = []byte(`
<InitiateMultipartUploadResult>
    <Bucket>test-1253846586</Bucket>
    <Key>test</Key>
</InitiateMultipartUploadResult>`)
	err = UnmarshalInitMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}
}

func Test_UnmarshalCompleteMultiUploadResult_Error(t *testing.T) {
	data := []byte(`
<CompleteMultipartUploadResult>
    <Bucket>test</Bucket>
    <Key>test/hello.txt</Key>
    <ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`)
	var res CompleteMultipartUploadResult
	err := UnmarshalCompleteMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}

	data = []byte(`
<CompleteMultipartUploadResult>
    <Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
    <Key>test/hello.txt</Key>
    <ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`)
	err = UnmarshalCompleteMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}

	data = []byte(`
<CompleteMultipartUploadResult>
    <Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
    <Bucket>test</Bucket>
    <ETag>&quot;594f98b11c6901c0f0683de1085a6d0e-4&quot;</ETag>
</CompleteMultipartUploadResult>`)
	err = UnmarshalCompleteMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}

	data = []byte(`
<CompleteMultipartUploadResult>
    <Location>test-1253846586.cos.ap-guangzhou.myqcloud.com/test/hello.txt</Location>
    <Bucket>test</Bucket>
    <Key>test/hello.txt</Key>
</CompleteMultipartUploadResult>`)
	err = UnmarshalCompleteMultiUploadResult(data, &res)
	if err == nil {
		t.Errorf("UnmarshalInitMultiUploadResult except return error")
	}
}

func Test_EncodeURIComponent(t *testing.T) {
	res := EncodeURIComponent("/")
	if res != "%2F" {
		t.Errorf("EncodeURIComponent failed, res: %v", res)
	}
}

func Test_GetRange(t *testing.T) {
	_, err := GetRange("byte=0-1")
	if err == nil {
		t.Errorf("GetRange expect failed")
	}

	_, err = GetRange("bytes=10")
	if err == nil {
		t.Errorf("GetRange expect failed")
	}

	_, err = GetRange("bytes=a-10")
	if err == nil {
		t.Errorf("GetRange expect failed")
	}

	_, err = GetRange("bytes=0-b")
	if err == nil {
		t.Errorf("GetRange expect failed")
	}

	_, err = GetRange("bytes=0-10")
	if err != nil {
		t.Errorf("GetRange failed: %v", err)
	}

	_, err = GetRange("bytes=0-")
	if err != nil {
		t.Errorf("GetRange failed: %v", err)
	}

	_, err = GetRange("bytes=-10")
	if err != nil {
		t.Errorf("GetRange failed: %v", err)
	}

	res, err := GetRangeOptions(&ObjectGetOptions{})
	if err != nil || res != nil {
		t.Errorf("GetRangeOptions failed: %v, %v", res, err)
	}
}

func Test_GetBucketRegionFromUrl(t *testing.T) {
	bucket, region := GetBucketRegionFromUrl(nil)
	if bucket != "" || region != "" {
		t.Errorf("GetBucketRegionFromUrl failed")
	}
	u, _ := url.Parse("http://test-125000000.cos.ap-guangzhou.myqcloud.com")
	bucket, region = GetBucketRegionFromUrl(u)
	if bucket != "test-125000000" && region != "ap-guangzhou" {
		t.Errorf("GetBucketRegionFromUrl failed")
	}
	u, _ = url.Parse("http://test-125000000.com")
	bucket, region = GetBucketRegionFromUrl(u)
	if bucket != "" || region != "" {
		t.Errorf("GetBucketRegionFromUrl failed")
	}
}

// ============================================================================
// 测试用例：EncodeUserData - S1 边转边播 user-data helper
// ============================================================================

// 覆盖：helper.go EncodeUserData nil 分支
func TestEncodeUserData_Nil(t *testing.T) {
	got, err := EncodeUserData(nil)
	if err != nil {
		t.Fatalf("EncodeUserData(nil) unexpected err: %v", err)
	}
	if got != "" {
		t.Errorf("EncodeUserData(nil) = %q, want empty string", got)
	}
}

// 覆盖：helper.go EncodeUserData 空 map 分支
// 空 map JSON 为 "{}", base64 = "e30"
func TestEncodeUserData_Empty(t *testing.T) {
	got, err := EncodeUserData(map[string]string{})
	if err != nil {
		t.Fatalf("EncodeUserData(empty) unexpected err: %v", err)
	}
	if got != "e30" {
		t.Errorf("EncodeUserData(empty) = %q, want %q", got, "e30")
	}
}

// 覆盖：helper.go EncodeUserData 稳定排序（同键值不同插入顺序应产生相同输出）
func TestEncodeUserData_Deterministic(t *testing.T) {
	a, err := EncodeUserData(map[string]string{"a": "1", "b": "2"})
	if err != nil {
		t.Fatalf("EncodeUserData a err: %v", err)
	}
	b, err := EncodeUserData(map[string]string{"b": "2", "a": "1"})
	if err != nil {
		t.Fatalf("EncodeUserData b err: %v", err)
	}
	if a != b {
		t.Errorf("EncodeUserData not deterministic: %q vs %q", a, b)
	}
	if a == "" {
		t.Errorf("EncodeUserData returned empty for non-empty map")
	}
}

// 覆盖：helper.go EncodeUserData 特殊字符 + 中文
// 断言：编码后可 base64 URL 解码 + JSON 反序列化，还原原值
func TestEncodeUserData_SpecialChars(t *testing.T) {
	in := map[string]string{
		"slash":   "a/b+c=d",
		"chinese": "测试值",
		"empty":   "",
	}
	got, err := EncodeUserData(in)
	if err != nil {
		t.Fatalf("EncodeUserData err: %v", err)
	}
	if got == "" {
		t.Fatalf("EncodeUserData returned empty for non-empty input")
	}
	// base64 URL 解码 + 反序列化
	raw, err := base64.RawURLEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("base64 decode err: %v, encoded: %q", err, got)
	}
	out := map[string]string{}
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("json unmarshal err: %v, raw: %q", err, string(raw))
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("roundtrip mismatch: in=%v out=%v", in, out)
	}
}
