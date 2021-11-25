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

func TestCIService_ImageAuditing(t *testing.T) {
	setup()
	defer teardown()
	name := "test.jpg"
	opt := &ImageRecognitionOptions{
		CIProcess:  "sensitive-content-recognition",
		DetectType: "porn",
		MaxFrames:  1,
	}
	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"ci-process":  "sensitive-content-recognition",
			"detect-type": "porn",
			"max-frames":  "1",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<RecognitionResult>
    <Result>1</Result>
    <Label>Porn</Label>
    <SubLabel>SexBehavior</SubLabel>
    <Score>90</Score>
    <PornInfo>
        <Code>0</Code>
        <Msg>OK</Msg>
        <HitFlag>1</HitFlag>
        <Label>xxx</Label>
        <SubLabel>SexBehavior</SubLabel>
        <Score>100</Score>
    </PornInfo>
</RecognitionResult>`)
	})

	want := &ImageRecognitionResult{
		XMLName:  xml.Name{Local: "RecognitionResult"},
		Result:   1,
		Label:    "Porn",
		SubLabel: "SexBehavior",
		Score:    90,
		PornInfo: &RecognitionInfo{
			Code:     0,
			Msg:      "OK",
			HitFlag:  1,
			Label:    "xxx",
			SubLabel: "SexBehavior",
			Score:    100,
		},
	}

	res, _, err := client.CI.ImageAuditing(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.ImageAuditing error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageAuditing failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_BatchImageAuditing(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<Request><Input><Object>test.jpg</Object></Input><Conf><DetectType>Porn,Terrorism,Politics,Ads</DetectType></Conf></Request>"

	mux.HandleFunc("/image/auditing", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &BatchImageAuditingOptions{
		Input: []ImageAuditingInputOptions{
			ImageAuditingInputOptions{
				Object: "test.jpg",
			},
		},
		Conf: &ImageAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Ads",
		},
	}

	_, _, err := client.CI.BatchImageAuditing(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.BatchImageAuditing error: %v", err)
	}
}

func TestCIService_PutVideoAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	name := "test.mp4"
	wantBody := "<Request><Input><Object>test.mp4</Object></Input>" +
		"<Conf><DetectType>Porn,Terrorism,Politics,Ads</DetectType>" +
		"<Snapshot><Mode>Interval</Mode><Count>100</Count><TimeInterval>50</TimeInterval></Snapshot>" +
		"<Callback>http://callback.com/call_back_test</Callback></Conf></Request>"

	mux.HandleFunc("/video/auditing", func(writer http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &PutVideoAuditingJobOptions{
		InputObject: name,
		Conf: &VideoAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Ads",
			Snapshot: &PutVideoAuditingJobSnapshot{
				Mode:         "Interval",
				Count:        100,
				TimeInterval: 50,
			},
			Callback: "http://callback.com/call_back_test",
		},
	}

	_, _, err := client.CI.PutVideoAuditingJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PutVideoAuditingJob returned error: %v", err)
	}
}

func TestCIService_GetVideoAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "vab1ca9fc8a3ed11ea834c525400863904"

	mux.HandleFunc("/video/auditing"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetVideoAuditingJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetVideoAuditingJob returned error: %v", err)
	}
}

func TestCIService_PutAudioAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	name := "test.mp4"
	wantBody := "<Request><Input><Object>test.mp4</Object></Input>" +
		"<Conf><DetectType>Porn,Terrorism,Politics,Ads</DetectType>" +
		"<Callback>http://callback.com/call_back_test</Callback></Conf></Request>"

	mux.HandleFunc("/audio/auditing", func(writer http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &PutAudioAuditingJobOptions{
		InputObject: name,
		Conf: &AudioAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Ads",
			Callback:   "http://callback.com/call_back_test",
		},
	}

	_, _, err := client.CI.PutAudioAuditingJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PutAudioAuditingJob returned error: %v", err)
	}
}

func TestCIService_GetAudioAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "vab1ca9fc8a3ed11ea834c525400863904"

	mux.HandleFunc("/audio/auditing"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetAudioAuditingJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetAudioAuditingJob returned error: %v", err)
	}
}

func TestCIService_PutTextAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	wantBody := `<Request><Input><Object>a.txt</Object></Input><Conf><DetectType>Porn,Terrorism,Politics,Ads,Illegal,Abuse</DetectType><Callback>http://callback.com/</Callback><BizType>b81d45f94b91a683255e9a95******</BizType></Conf></Request>`

	mux.HandleFunc("/text/auditing", func(writer http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &PutTextAuditingJobOptions{
		InputObject: "a.txt",
		Conf: &TextAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Ads,Illegal,Abuse",
			Callback:   "http://callback.com/",
			BizType:    "b81d45f94b91a683255e9a95******",
		},
	}

	_, _, err := client.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PutTextAuditingJob returned error: %v", err)
	}
}

func TestCIService_GetTextAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "vab1ca9fc8a3ed11ea834c525400863904"

	mux.HandleFunc("/text/auditing"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetTextAuditingJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetTextAuditingJob returned error: %v", err)
	}
}

func TestCIService_PutDocumentAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	wantBody := `<Request><Input><Url>http://www.example.com/doctest.doc</Url><Type>doc</Type></Input><Conf><DetectType>Porn,Terrorism,Politics,Ads</DetectType><Callback>http://www.example.com/</Callback><BizType>b81d45f94b91a683255e9a9506f4****</BizType></Conf></Request>`

	mux.HandleFunc("/document/auditing", func(writer http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &PutDocumentAuditingJobOptions{
		InputUrl:  "http://www.example.com/doctest.doc",
		InputType: "doc",
		Conf: &DocumentAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Ads",
			Callback:   "http://www.example.com/",
			BizType:    "b81d45f94b91a683255e9a9506f4****",
		},
	}

	_, _, err := client.CI.PutDocumentAuditingJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PutDocumentAuditingJob returned error: %v", err)
	}
}

func TestCIService_GetDocumentAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "vab1ca9fc8a3ed11ea834c525400863904"

	mux.HandleFunc("/document/auditing"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetDocumentAuditingJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetDocumentAuditingJob returned error: %v", err)
	}
}

func TestCIService_PutWebpageAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	wantBody := `<Request><Input><Url>http://www.example.com</Url></Input><Conf><DetectType>Porn,Terrorism,Politics,Ads</DetectType><Callback>http://www.example.com/</Callback></Conf></Request>`

	mux.HandleFunc("/webpage/auditing", func(writer http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &PutWebpageAuditingJobOptions{
		InputUrl:  "http://www.example.com",
		Conf: &WebpageAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Ads",
			Callback:   "http://www.example.com/",
		},
	}

	_, _, err := client.CI.PutWebpageAuditingJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PutWebpageAuditingJob returned error: %v", err)
	}
}

func TestCIService_GetWebpageAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "shb1ca9fc8a3ed11ea834c525400863904"

	mux.HandleFunc("/webpage/auditing/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetWebpageAuditingJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetWebpageAuditingJob returned error: %v", err)
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

func TestCIService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"imageMogr2/thumbnail/!50p": "",
		}
		testFormValues(t, r, vs)
	})

	_, err := client.CI.Get(context.Background(), "test.jpg", "imageMogr2/thumbnail/!50p", nil)
	if err != nil {
		t.Fatalf("CI.Get returned error: %v", err)
	}
}

func TestCIService_GetToFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"imageMogr2/thumbnail/!50p": "",
		}
		testFormValues(t, r, vs)
	})

	filepath := "test.jpg." + time.Now().Format(time.RFC3339)
	defer os.Remove(filepath)

	_, err := client.CI.GetToFile(context.Background(), "test.jpg", filepath, "imageMogr2/thumbnail/!50p", nil)
	if err != nil {
		t.Fatalf("CI.GetToFile returned error: %v", err)
	}
}

func TestCIService_GetQRcode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"ci-process": "QRcode",
			"cover":      "1",
		}
		testFormValues(t, r, vs)
	})

	_, _, err := client.CI.GetQRcode(context.Background(), "test.jpg", 1, nil)
	if err != nil {
		t.Fatalf("CI.GetQRcode returned error: %v", err)
	}
}

func TestCIService_GenerateQRcode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"ci-process":     "qrcode-generate",
			"qrcode-content": "<https://www.example.com>",
			"mode":           "1",
			"width":          "200",
		}
		testFormValues(t, r, vs)
	})

	opt := &GenerateQRcodeOptions{
		QRcodeContent: "<https://www.example.com>",
		Mode:          1,
		Width:         200,
	}

	_, _, err := client.CI.GenerateQRcode(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.GenerateQRcode returned error: %v", err)
	}
}

func TestCIService_GenerateQRcodeToFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"ci-process":     "qrcode-generate",
			"qrcode-content": "<https://www.example.com>",
			"mode":           "1",
			"width":          "200",
		}
		testFormValues(t, r, vs)
	})

	opt := &GenerateQRcodeOptions{
		QRcodeContent: "<https://www.example.com>",
		Mode:          1,
		Width:         200,
	}

	filepath := "test.file." + time.Now().Format(time.RFC3339)
	defer os.Remove(filepath)

	_, _, err := client.CI.GenerateQRcodeToFile(context.Background(), filepath, opt)
	if err != nil {
		t.Fatalf("CI.GenerateQRcode returned error: %v", err)
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
