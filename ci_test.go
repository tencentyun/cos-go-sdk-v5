package cos

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCIService_EncodePicOperations(t *testing.T) {
	opt := &PicOperations{
		IsPicInfo: 1,
		Rules: []PicOperationsRules{
			{
				FileId: "example.jpg",
				Rule:   "imageView2/format/png",
			},
		},
	}
	res := EncodePicOperations(opt)
	jsonStr := `{"is_pic_info":1,"rules":[{"fileid":"example.jpg","rule":"imageView2/format/png"}]}`
	if jsonStr != res {
		t.Fatalf("EncodePicOperations Failed, returned:%v, want:%v", res, jsonStr)
	}
}

func TestCIService_ImageProcess(t *testing.T) {
	setup()
	defer teardown()
	name := "test.jpg"

	opt := &ImageProcessOptions{
		IsPicInfo: 1,
		Rules: []PicOperationsRules{
			{
				FileId: "format.jpg",
				Rule:   "imageView2/format/png",
			},
		},
	}
	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		vs := values{
			"image_process": "",
		}
		testFormValues(t, r, vs)
		header := r.Header.Get("Pic-Operations")
		body := new(ImageProcessOptions)
		err := json.Unmarshal([]byte(header), body)
		want := opt
		if err != nil {
			t.Errorf("CI.ImageProcess Failed: %v", err)
		}
		if !reflect.DeepEqual(want, body) {
			t.Errorf("CI.ImageProcess Failed, wanted:%v, body:%v", want, body)
		}
		fmt.Fprint(w, `<UploadResult>
    <OriginalInfo>
        <Key>test.jpg</Key>
        <Location>example-1250000000.cos.ap-guangzhou.myqcloud.com/test.jpg</Location>
        <ETag>&quot;8894dbe5e3ebfaf761e39b9d619c28f3327b8d85&quot;</ETag>
        <ImageInfo>
            <Format>PNG</Format>
            <Width>103</Width>
            <Height>99</Height>
            <Quality>100</Quality>
            <Ave>0xa08162</Ave>
            <Orientation>0</Orientation>
        </ImageInfo>
    </OriginalInfo>
    <ProcessResults>
        <Object>
            <Key>format.jpg</Key>
            <Location>example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg</Location>
            <Format>PNG</Format>
            <Width>103</Width>
            <Height>99</Height>
            <Size>21351</Size>
            <Quality>100</Quality>
            <ETag>&quot;8894dbe5e3ebfaf761e39b9d619c28f3327b8d85&quot;</ETag>
        </Object>
    </ProcessResults>
</UploadResult>`)
	})

	want := &ImageProcessResult{
		XMLName: xml.Name{Local: "UploadResult"},
		OriginalInfo: &PicOriginalInfo{
			Key:      "test.jpg",
			Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/test.jpg",
			ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
			ImageInfo: &PicImageInfo{
				Format:      "PNG",
				Width:       103,
				Height:      99,
				Quality:     100,
				Ave:         "0xa08162",
				Orientation: 0,
			},
		},
		ProcessResults: &PicProcessObject{
			Key:      "format.jpg",
			Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg",
			Format:   "PNG",
			Width:    103,
			Height:   99,
			Size:     21351,
			Quality:  100,
			ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
		},
	}

	res, _, err := client.CI.ImageProcess(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.ImageProcess returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageProcess failed, return:%v, want:%v", res, want)
	}
}

func TestCIService_ImageRecognition(t *testing.T) {
	setup()
	defer teardown()
	name := "test.jpg"

	detectType := "porn,terrorist,politics"
	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"ci-process":  "sensitive-content-recognition",
			"detect-type": "porn,terrorist,politics",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<RecognitionResult>
    <PornInfo>
        <Code>0</Code>
        <Msg>OK</Msg>
        <HitFlag>0</HitFlag>
        <Score>0</Score>
        <Label/>
    </PornInfo>
    <TerroristInfo>
        <Code>0</Code>
        <Msg>OK</Msg>
        <HitFlag>0</HitFlag>
        <Score>0</Score>
        <Label/>
    </TerroristInfo>
    <PoliticsInfo>
        <Code>0</Code>
        <Msg>OK</Msg>
        <HitFlag>0</HitFlag>
        <Score>0</Score>
        <Label/>
    </PoliticsInfo>
</RecognitionResult>`)
	})

	want := &ImageRecognitionResult{
		XMLName: xml.Name{Local: "RecognitionResult"},
		PornInfo: &RecognitionInfo{
			Code:    0,
			Msg:     "OK",
			HitFlag: 0,
			Score:   0,
		},
		TerroristInfo: &RecognitionInfo{
			Code:    0,
			Msg:     "OK",
			HitFlag: 0,
			Score:   0,
		},
		PoliticsInfo: &RecognitionInfo{
			Code:    0,
			Msg:     "OK",
			HitFlag: 0,
			Score:   0,
		},
	}

	res, _, err := client.CI.ImageRecognition(context.Background(), name, detectType)
	if err != nil {
		t.Fatalf("CI.ImageRecognitionreturned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageRecognition failed, return:%v, want:%v", res, want)
	}
}

func TestCIService_Put(t *testing.T) {
	setup()
	defer teardown()
	name := "test.jpg"
	data := make([]byte, 1024*1024*3)
	rand.Read(data)

	pic := &ImageProcessOptions{
		IsPicInfo: 1,
		Rules: []PicOperationsRules{
			{
				FileId: "format.jpg",
				Rule:   "imageView2/format/png",
			},
		},
	}
	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		header := r.Header.Get("Pic-Operations")
		body := new(ImageProcessOptions)
		err := json.Unmarshal([]byte(header), body)
		want := pic
		if err != nil {
			t.Errorf("CI.Put Failed: %v", err)
		}
		if !reflect.DeepEqual(want, body) {
			t.Errorf("CI.Put Failed, wanted:%v, body:%v", want, body)
		}
		tb := crc64.MakeTable(crc64.ECMA)
		ht := crc64.New(tb)
		tr := TeeReader(r.Body, ht, 0, nil)
		bs, err := ioutil.ReadAll(tr)
		if err != nil {
			t.Errorf("CI.Put ReadAll Failed: %v", err)
		}
		if bytes.Compare(bs, data) != 0 {
			t.Errorf("CI.Put Failed, data isn't consistent")
		}
		crc := tr.Crc64()
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
		fmt.Fprint(w, `<UploadResult>
    <OriginalInfo>
        <Key>test.jpg</Key>
        <Location>example-1250000000.cos.ap-guangzhou.myqcloud.com/test.jpg</Location>
        <ETag>&quot;8894dbe5e3ebfaf761e39b9d619c28f3327b8d85&quot;</ETag>
        <ImageInfo>
            <Format>PNG</Format>
            <Width>103</Width>
            <Height>99</Height>
            <Quality>100</Quality>
            <Ave>0xa08162</Ave>
            <Orientation>0</Orientation>
        </ImageInfo>
    </OriginalInfo>
    <ProcessResults>
        <Object>
            <Key>format.jpg</Key>
            <Location>example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg</Location>
            <Format>PNG</Format>
            <Width>103</Width>
            <Height>99</Height>
            <Size>21351</Size>
            <Quality>100</Quality>
            <ETag>&quot;8894dbe5e3ebfaf761e39b9d619c28f3327b8d85&quot;</ETag>
        </Object>
    </ProcessResults>
</UploadResult>`)
	})

	want := &ImageProcessResult{
		XMLName: xml.Name{Local: "UploadResult"},
		OriginalInfo: &PicOriginalInfo{
			Key:      "test.jpg",
			Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/test.jpg",
			ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
			ImageInfo: &PicImageInfo{
				Format:      "PNG",
				Width:       103,
				Height:      99,
				Quality:     100,
				Ave:         "0xa08162",
				Orientation: 0,
			},
		},
		ProcessResults: &PicProcessObject{
			Key:      "format.jpg",
			Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg",
			Format:   "PNG",
			Width:    103,
			Height:   99,
			Size:     21351,
			Quality:  100,
			ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
		},
	}

	f := bytes.NewReader(data)
	opt := &ObjectPutOptions{
		nil,
		&ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", EncodePicOperations(pic))
	res, _, err := client.CI.Put(context.Background(), name, f, opt)
	if err != nil {
		t.Fatalf("CI.Put returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageProcess failed, return:%v, want:%v", res, want)
	}
}

func TestCIService_PutFromFile(t *testing.T) {
	setup()
	defer teardown()
	name := "test.jpg"
	filePath := "test.file" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("creat tmp file failed")
	}
	defer os.Remove(filePath)
	data := make([]byte, 1024*1024*3)
	rand.Read(data)
	newfile.Write(data)
	newfile.Close()

	pic := &ImageProcessOptions{
		IsPicInfo: 1,
		Rules: []PicOperationsRules{
			{
				FileId: "format.jpg",
				Rule:   "imageView2/format/png",
			},
		},
	}
	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		header := r.Header.Get("Pic-Operations")
		body := new(ImageProcessOptions)
		err := json.Unmarshal([]byte(header), body)
		want := pic
		if err != nil {
			t.Errorf("CI.Put Failed: %v", err)
		}
		if !reflect.DeepEqual(want, body) {
			t.Errorf("CI.Put Failed, wanted:%v, body:%v", want, body)
		}
		tb := crc64.MakeTable(crc64.ECMA)
		ht := crc64.New(tb)
		tr := TeeReader(r.Body, ht, 0, nil)
		bs, err := ioutil.ReadAll(tr)
		if err != nil {
			t.Errorf("CI.Put ReadAll Failed: %v", err)
		}
		if bytes.Compare(bs, data) != 0 {
			t.Errorf("CI.Put Failed, data isn't consistent")
		}
		crc := tr.Crc64()
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))
		fmt.Fprint(w, `<UploadResult>
    <OriginalInfo>
        <Key>test.jpg</Key>
        <Location>example-1250000000.cos.ap-guangzhou.myqcloud.com/test.jpg</Location>
        <ETag>&quot;8894dbe5e3ebfaf761e39b9d619c28f3327b8d85&quot;</ETag>
        <ImageInfo>
            <Format>PNG</Format>
            <Width>103</Width>
            <Height>99</Height>
            <Quality>100</Quality>
            <Ave>0xa08162</Ave>
            <Orientation>0</Orientation>
        </ImageInfo>
    </OriginalInfo>
    <ProcessResults>
        <Object>
            <Key>format.jpg</Key>
            <Location>example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg</Location>
            <Format>PNG</Format>
            <Width>103</Width>
            <Height>99</Height>
            <Size>21351</Size>
            <Quality>100</Quality>
            <ETag>&quot;8894dbe5e3ebfaf761e39b9d619c28f3327b8d85&quot;</ETag>
        </Object>
    </ProcessResults>
</UploadResult>`)
	})

	want := &ImageProcessResult{
		XMLName: xml.Name{Local: "UploadResult"},
		OriginalInfo: &PicOriginalInfo{
			Key:      "test.jpg",
			Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/test.jpg",
			ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
			ImageInfo: &PicImageInfo{
				Format:      "PNG",
				Width:       103,
				Height:      99,
				Quality:     100,
				Ave:         "0xa08162",
				Orientation: 0,
			},
		},
		ProcessResults: &PicProcessObject{
			Key:      "format.jpg",
			Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg",
			Format:   "PNG",
			Width:    103,
			Height:   99,
			Size:     21351,
			Quality:  100,
			ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
		},
	}

	opt := &ObjectPutOptions{
		nil,
		&ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", EncodePicOperations(pic))
	res, _, err := client.CI.PutFromFile(context.Background(), name, filePath, opt)
	if err != nil {
		t.Fatalf("CI.Put returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageProcess failed, return:%v, want:%v", res, want)
	}
}

func TestBucketService_GetGuetzli(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"guetzli": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<GuetzliStatus>on</GuetzliStatus>`)
	})

	res, _, err := client.CI.GetGuetzli(context.Background())
	if err != nil {
		t.Fatalf("CI.GetGuetzli returned error %v", err)
	}

	want := &GetGuetzliResult{
		XMLName:       xml.Name{Local: "GuetzliStatus"},
		GuetzliStatus: "on",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetGuetzli %+v, want %+v", res, want)
	}
}

func TestBucketService_PutGuetzli(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"guetzli": "",
		}
		testFormValues(t, r, vs)
	})

	_, err := client.CI.PutGuetzli(context.Background())
	if err != nil {
		t.Fatalf("CI.PutGuetzli returned error: %v", err)
	}
}

func TestBucketService_DeleteGuetzli(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		vs := values{
			"guetzli": "",
		}
		testFormValues(t, r, vs)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.CI.DeleteGuetzli(context.Background())
	if err != nil {
		t.Fatalf("CI.PutGuetzli returned error: %v", err)
	}
}
