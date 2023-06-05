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
	{
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
	{
		res := EncodePicOperations(nil)
		if res != "" {
			t.Fatalf("EncodePicOperations Failed, returned:%v", res)
		}
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
		ProcessResults: []PicProcessObject{
			{
				Key:      "format.jpg",
				Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg",
				Format:   "PNG",
				Width:    103,
				Height:   99,
				Size:     21351,
				Quality:  100,
				ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
			},
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
			"ci-process": "sensitive-content-recognition",
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
		InputUrl: "http://www.example.com",
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
		ProcessResults: []PicProcessObject{
			{
				Key:      "format.jpg",
				Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg",
				Format:   "PNG",
				Width:    103,
				Height:   99,
				Size:     21351,
				Quality:  100,
				ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
			},
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
		ProcessResults: []PicProcessObject{
			{
				Key:      "format.jpg",
				Location: "example-1250000000.cos.ap-guangzhou.myqcloud.com/format.jpg",
				Format:   "PNG",
				Width:    103,
				Height:   99,
				Size:     21351,
				Quality:  100,
				ETag:     "\"8894dbe5e3ebfaf761e39b9d619c28f3327b8d85\"",
			},
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
	{
		setup()

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
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			vs := values{
				"imageMogr2/thumbnail/!50p": "",
				"versionId":                 "1.1",
			}
			testFormValues(t, r, vs)
		})
		_, err := client.CI.Get(context.Background(), "test.jpg", "imageMogr2/thumbnail/!50p", nil, "1.1")
		if err != nil {
			t.Fatalf("CI.Get returned error: %v", err)
		}
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			vs := values{
				"imageMogr2/thumbnail/!50p": "",
			}
			testFormValues(t, r, vs)
		})
		_, err := client.CI.Get(context.Background(), "test.jpg", "imageMogr2/thumbnail/!50p", nil, "1.1", "1.2")
		if err == nil || err.Error() != "wrong params" {
			t.Fatalf("CI.Get returned error: %v", err)
		}
		teardown()
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
	{
		setup()
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
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			vs := values{
				"ci-process": "QRcode",
				"cover":      "1",
				"versionId":  "1.1",
			}
			testFormValues(t, r, vs)
		})
		_, _, err := client.CI.GetQRcode(context.Background(), "test.jpg", 1, nil, "1.1")
		if err != nil {
			t.Fatalf("CI.GetQRcode returned error: %v", err)
		}
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			vs := values{
				"ci-process": "QRcode",
				"cover":      "1",
			}
			testFormValues(t, r, vs)
		})
		_, _, err := client.CI.GetQRcode(context.Background(), "test.jpg", 1, nil, "1.1", "1.2")
		if err == nil || err.Error() != "wrong params" {
			t.Fatalf("CI.GetQRcode returned error: %v", err)
		}
		teardown()
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

func TestCIService_AIBodyRecognition(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"ci-process": "AIBodyRecognition",
			"detect-url": "http://test123.com/test.jpg",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<RecognitionResult>
	<Status>1</Status>
	<PedestrianInfo>
		<Name>person</Name>
		<Score>91</Score>
		<Location>
			<Point>77,37</Point>
			<Point>77,346</Point>
			<Point>522,346</Point>
			<Point>522,37</Point>
		</Location>
	</PedestrianInfo>
</RecognitionResult>`)
	})

	opt := &AIBodyRecognitionOptions{
		DetectUrl: "http://test123.com/test.jpg",
	}

	want := &AIBodyRecognitionResult{
		XMLName: xml.Name{Local: "RecognitionResult"},
		Status:  1,
		PedestrianInfo: []PedestrianInfo{
			{
				Name:  "person",
				Score: 91,
				Location: &PedestrianLocation{
					Point: []string{
						"77,37",
						"77,346",
						"522,346",
						"522,37",
					},
				},
			},
		},
	}

	res, _, err := client.CI.AIBodyRecognition(context.Background(), "", opt)
	if err != nil {
		t.Fatalf("CI.AIBodyRecognition returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.AIBodyRecognition failed, return:%v, want:%v", res, want)
	}
}

func TestCIService_GetImageAuditingJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "si8e6586b0a6e411edb8db525400a28986"

	mux.HandleFunc("/image/auditing/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `<Response>
  <JobsDetail>
    <JobId>si8e6586b0a6e411edb8db525400a28986</JobId>
    <State>Success</State>
    <DataId>99999999999999977777777777777777</DataId>
    <Object>202106300295654374239be721f.jpg</Object>
    <Text>xxxyyyzzz</Text>
    <Label>Politics</Label>
    <Result>1</Result>
    <Score>99</Score>
    <Category>PolityOCR</Category>
    <PornInfo></PornInfo>
    <TerrorismInfo></TerrorismInfo>
    <PoliticsInfo>
      <HitFlag>1</HitFlag>
      <Score>99</Score>
      <Label>xxx</Label>
      <Category>PolityOCR</Category>
      <OcrResults>
        <Text>xxxyyyzzz</Text>
        <Keywords>xxx</Keywords>
        <Keywords>yyy</Keywords>
        <Keywords>zzz</Keywords>
        <Location></Location>
      </OcrResults>
    </PoliticsInfo>
    <AdsInfo></AdsInfo>
  </JobsDetail>
  <RequestId>NjNlMjQ2YzJfNDE3OTgyMDlfNGI5NV9jZjU=</RequestId>
</Response>
`)
	})

	want := &GetImageAuditingJobResult{
		XMLName: xml.Name{Local: "Response"},
		JobsDetail: &ImageAuditingResult{
			JobId:         "si8e6586b0a6e411edb8db525400a28986",
			State:         "Success",
			DataId:        "99999999999999977777777777777777",
			Object:        "202106300295654374239be721f.jpg",
			Text:          "xxxyyyzzz",
			Label:         "Politics",
			Result:        1,
			Score:         99,
			Category:      "PolityOCR",
			PornInfo:      &RecognitionInfo{},
			TerrorismInfo: &RecognitionInfo{},
			PoliticsInfo: &RecognitionInfo{
				HitFlag:  1,
				Score:    99,
				Label:    "xxx",
				Category: "PolityOCR",
				OcrResults: []OcrResult{OcrResult{
					Text: "xxxyyyzzz",
					Keywords: []string{
						"xxx",
						"yyy",
						"zzz",
					},
					Location: &Location{},
				}},
			},
			AdsInfo: &RecognitionInfo{},
		},
		RequestId: "NjNlMjQ2YzJfNDE3OTgyMDlfNGI5NV9jZjU=",
	}

	res, _, err := client.CI.GetImageAuditingJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetImageAuditingJob returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetImageAuditingJob failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_PostVideoAuditingCancelJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "av8e6586b0a6e411edb8db525400a28986"

	mux.HandleFunc("/video/cancel_auditing/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, _, err := client.CI.PostVideoAuditingCancelJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.PostVideoAuditingCancelJob returned error: %v", err)
	}
}

func TestCIService_ReportBadcase(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<Request><ContentType>1</ContentType><Text>abc</Text>" +
		"<Label>Ad</Label><SuggestedLabel>Normal</SuggestedLabel></Request>"

	mux.HandleFunc("/report/badcase", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &ReportBadcaseOptions{
		ContentType:    1,
		Text:           "abc",
		Label:          "Ad",
		SuggestedLabel: "Normal",
	}

	_, _, err := client.CI.ReportBadcase(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.ReportBadcase returned error: %v", err)
	}
}

func TestCIService_PutVirusDetectJob(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<Request><Input><Object>a.exe</Object></Input>" +
		"<Conf><DetectType>Virus</DetectType>" +
		"<Callback>http://callback.com/call_back_test</Callback></Conf></Request>"

	mux.HandleFunc("/virus/detect", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &PutVirusDetectJobOptions{
		InputObject: "a.exe",
		Conf: &VirusDetectJobConf{
			DetectType: "Virus",
			Callback:   "http://callback.com/call_back_test",
		},
	}

	_, _, err := client.CI.PutVirusDetectJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PutVirusDetectJob returned error: %v", err)
	}
}

func TestCIService_GetVirusDetectJob(t *testing.T) {
	setup()
	defer teardown()
	jobID := "ssb1ca9fc8a3ed11ea834c525400863904"

	mux.HandleFunc("/virus/detect"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetVirusDetectJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.GetVirusDetectJob returned error: %v", err)
	}
}

func TestCIService_AddStyle(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<AddStyle><StyleName>style1</StyleName>" +
		"<StyleBody>imageMogr2/thumbnail/!50px</StyleBody></AddStyle>"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		v := values{
			"style": "",
		}
		testFormValues(t, r, v)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &AddStyleOptions{
		StyleName: "style1",
		StyleBody: "imageMogr2/thumbnail/!50px",
	}

	_, err := client.CI.AddStyle(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.AddStyle returned error: %v", err)
	}
}

func TestCIService_GetStyle(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<GetStyle><StyleName>style1</StyleName></GetStyle>"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"style": "",
		}
		testFormValues(t, r, v)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
		fmt.Fprint(w, `<StyleList>
  <StyleRule>
    <StyleName>style1</StyleName>
    <StyleBody>imageMogr2/thumbnail/!50px</StyleBody>
  </StyleRule>
</StyleList>
`)
	})

	opt := &GetStyleOptions{
		StyleName: "style1",
	}

	want := &GetStyleResult{
		XMLName: xml.Name{Local: "StyleList"},
		StyleRule: []StyleRule{StyleRule{
			StyleName: "style1",
			StyleBody: "imageMogr2/thumbnail/!50px",
		},
		},
	}

	res, _, err := client.CI.GetStyle(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.GetStyle returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetStyle failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_DeleteStyle(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<DeleteStyle><StyleName>style1</StyleName></DeleteStyle>"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		v := values{
			"style": "",
		}
		testFormValues(t, r, v)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &DeleteStyleOptions{
		StyleName: "style1",
	}

	_, err := client.CI.DeleteStyle(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DeleteStyle returned error: %v", err)
	}
}

func TestCIService_ImageQuality(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AssessQuality",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <LongImage>TRUE</LongImage>
  <BlackAndWhite>TRUE</BlackAndWhite>
  <SmallImage>TRUE</SmallImage>
  <BigImage>FALSE</BigImage>
  <PureImage>FALSE</PureImage>
  <ClarityScore>50</ClarityScore>
  <AestheticScore>50</AestheticScore>
  <RequestId>xxx</RequestId>
</Response>
`)
	})

	want := &ImageQualityResult{
		XMLName:        xml.Name{Local: "Response"},
		LongImage:      true,
		BlackAndWhite:  true,
		SmallImage:     true,
		BigImage:       false,
		PureImage:      false,
		ClarityScore:   50,
		AestheticScore: 50,
		RequestId:      "xxx",
	}

	res, _, err := client.CI.ImageQuality(context.Background(), "test.jpg")
	if err != nil {
		t.Fatalf("CI.ImageQuality returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageQuality failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_ImageQualityWithOpt(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AssessQuality",
			"detect-url": "http://test123.com/test.jpg",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <LongImage>TRUE</LongImage>
  <BlackAndWhite>TRUE</BlackAndWhite>
  <SmallImage>TRUE</SmallImage>
  <BigImage>FALSE</BigImage>
  <PureImage>FALSE</PureImage>
  <ClarityScore>50</ClarityScore>
  <AestheticScore>50</AestheticScore>
  <RequestId>xxx</RequestId>
</Response>
`)
	})

	opt := &ImageQualityOptions{
		DetectUrl: "http://test123.com/test.jpg",
	}

	want := &ImageQualityResult{
		XMLName:        xml.Name{Local: "Response"},
		LongImage:      true,
		BlackAndWhite:  true,
		SmallImage:     true,
		BigImage:       false,
		PureImage:      false,
		ClarityScore:   50,
		AestheticScore: 50,
		RequestId:      "xxx",
	}

	res, _, err := client.CI.ImageQualityWithOpt(context.Background(), "", opt)
	if err != nil {
		t.Fatalf("CI.ImageQualityWithOpt returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.ImageQualityWithOpt failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_OcrRecognition(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":    "OCR",
			"type":          "general",
			"language-type": "zh",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <Angel>359.99</Angel>
  <Language>zh</Language>
  <PdfPageSize>0</PdfPageSize>
  <RequestId>xxx</RequestId>
  <TextDetections>
      <Confidence>99</Confidence>
      <DetectedText>你好</DetectedText>
      <ItemPolygon>
          <Height>64</Height>
          <Width>123</Width>
          <X>140</X>
          <Y>167</Y>
      </ItemPolygon>
      <Polygon>
          <X>140</X>
          <Y>167</Y>
      </Polygon>
      <Polygon>
          <X>263</X>
          <Y>167</Y>
      </Polygon>
      <Polygon>
          <X>263</X>
          <Y>231</Y>
      </Polygon>
      <Polygon>
          <X>140</X>
          <Y>231</Y>
      </Polygon>
      <Words>
          <Character>你</Character>
          <Confidence>99</Confidence>
          <WordCoordPoint>
              <WordCoordinate>
                  <X>212</X>
                  <Y>167</Y>
              </WordCoordinate>
              <WordCoordinate>
                  <X>341</X>
                  <Y>167</Y>
              </WordCoordinate>
              <WordCoordinate>
                  <X>341</X>
                  <Y>231</Y>
              </WordCoordinate>
              <WordCoordinate>
                  <X>212</X>
                  <Y>231</Y>
              </WordCoordinate>
          </WordCoordPoint>
      </Words>
      <Words>
          <Character>好</Character>
          <Confidence>99</Confidence>
          <WordCoordPoint>
              <WordCoordinate>
                  <X>341</X>
                  <Y>167</Y>
              </WordCoordinate>
              <WordCoordinate>
                  <X>263</X>
                  <Y>167</Y>
              </WordCoordinate>
              <WordCoordinate>
                  <X>263</X>
                  <Y>231</Y>
              </WordCoordinate>
              <WordCoordinate>
                  <X>341</X>
                  <Y>230</Y>
              </WordCoordinate>
          </WordCoordPoint>
      </Words>
  </TextDetections>
</Response>
`)
	})

	opt := &OcrRecognitionOptions{
		Type:         "general",
		LanguageType: "zh",
		Ispdf:        false,
		Isword:       false,
	}

	want := &OcrRecognitionResult{
		XMLName: xml.Name{Local: "Response"},
		TextDetections: []TextDetections{
			TextDetections{
				DetectedText: "你好",
				Confidence:   99,
				Polygon: []Polygon{
					Polygon{
						X: 140,
						Y: 167,
					},
					Polygon{
						X: 263,
						Y: 167,
					},
					Polygon{
						X: 263,
						Y: 231,
					},
					Polygon{
						X: 140,
						Y: 231,
					},
				},
				ItemPolygon: []ItemPolygon{
					ItemPolygon{
						X:      140,
						Y:      167,
						Width:  123,
						Height: 64,
					},
				},
				Words: []Words{
					Words{
						Confidence: 99,
						Character:  "你",
						WordCoordPoint: &WordCoordPoint{
							WordCoordinate: []Polygon{
								Polygon{
									X: 212,
									Y: 167,
								},
								Polygon{
									X: 341,
									Y: 167,
								},
								Polygon{
									X: 341,
									Y: 231,
								},
								Polygon{
									X: 212,
									Y: 231,
								},
							},
						},
					},
					Words{
						Confidence: 99,
						Character:  "好",
						WordCoordPoint: &WordCoordPoint{
							WordCoordinate: []Polygon{
								Polygon{
									X: 341,
									Y: 167,
								},
								Polygon{
									X: 263,
									Y: 167,
								},
								Polygon{
									X: 263,
									Y: 231,
								},
								Polygon{
									X: 341,
									Y: 230,
								},
							},
						},
					},
				},
			},
		},
		Language:    "zh",
		Angel:       359.99,
		PdfPageSize: 0,
		RequestId:   "xxx",
	}

	res, _, err := client.CI.OcrRecognition(context.Background(), "test.jpg", opt)
	if err != nil {
		t.Fatalf("CI.OcrRecognition returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.OcrRecognition failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_DetectCar(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "DetectCar",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <RequestId>xxx</RequestId>
  <CarTags>
      <Serial>五菱宏光</Serial>
      <Brand>五菱</Brand>
      <Type>面包车</Type>
      <Color>白</Color>
      <Confidence>0</Confidence>
      <Year>0</Year>
      <CarLocation>
          <X>0</X>
          <Y>228</Y>
      </CarLocation>
      <CarLocation>
          <X>0</X>
          <Y>81</Y>
      </CarLocation>
      <CarLocation>
          <X>73</X>
          <Y>81</Y>
      </CarLocation>
      <CarLocation>
          <X>73</X>
          <Y>228</Y>
      </CarLocation>
  </CarTags>
</Response>
`)
	})

	want := &DetectCarResult{
		XMLName: xml.Name{Local: "Response"},
		CarTags: []CarTags{
			CarTags{
				Serial:     "五菱宏光",
				Brand:      "五菱",
				Type:       "面包车",
				Color:      "白",
				Confidence: 0,
				Year:       0,
				CarLocation: []CarLocation{
					CarLocation{
						X: 0,
						Y: 228,
					},
					CarLocation{
						X: 0,
						Y: 81,
					},
					CarLocation{
						X: 73,
						Y: 81,
					},
					CarLocation{
						X: 73,
						Y: 228,
					},
				},
			},
		},
		RequestId: "xxx",
	}

	res, _, err := client.CI.DetectCar(context.Background(), "test.jpg")
	if err != nil {
		t.Fatalf("CI.DetectCar returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.DetectCar failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_OpenCIService(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.CI.OpenCIService(context.Background())
	if err != nil {
		t.Fatalf("CI.OpenCIService returned error: %v", err)
	}
}

func TestCIService_GetCIService(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `<CIStatus>on</CIStatus>`)
	})

	want := &CIServiceResult{
		XMLName:  xml.Name{Local: "CIStatus"},
		CIStatus: "on",
	}

	res, _, err := client.CI.GetCIService(context.Background())
	if err != nil {
		t.Fatalf("CI.GetCIService returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetCIService failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_CloseCIService(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		v := values{
			"unbind": "",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.CloseCIService(context.Background())
	if err != nil {
		t.Fatalf("CI.CloseCIService returned error: %v", err)
	}
}

func TestCIService_SetHotLink(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<Hotlink><Url>www.example.com</Url>" +
		"<Type>black</Type></Hotlink>"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		v := values{
			"hotlink": "",
		}
		testFormValues(t, r, v)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &HotLinkOptions{
		Url: []string{
			"www.example.com",
		},
		Type: "black",
	}

	_, err := client.CI.SetHotLink(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.SetHotLink returned error: %v", err)
	}
}

func TestCIService_GetHotLink(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"hotlink": "",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Hotlink>
  <Status>on</Status>
  <Type>white</Type>
  <Url>xxx</Url>
  <Url>yyy</Url>
</Hotlink>`)
	})

	want := &HotLinkResult{
		XMLName: xml.Name{Local: "Hotlink"},
		Status:  "on",
		Type:    "white",
		Url: []string{
			"xxx",
			"yyy",
		},
	}

	res, _, err := client.CI.GetHotLink(context.Background())
	if err != nil {
		t.Fatalf("CI.GetHotLink returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetHotLink failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_OpenOriginProtect(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		v := values{
			"origin-protect": "",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.OpenOriginProtect(context.Background())
	if err != nil {
		t.Fatalf("CI.OpenOriginProtect returned error: %v", err)
	}
}

func TestCIService_GetOriginProtect(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"origin-protect": "",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<OriginProtectStatus>on</OriginProtectStatus>`)
	})

	want := &OriginProtectResult{
		XMLName:             xml.Name{Local: "OriginProtectStatus"},
		OriginProtectStatus: "on",
	}

	res, _, err := client.CI.GetOriginProtect(context.Background())
	if err != nil {
		t.Fatalf("CI.GetOriginProtect returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetOriginProtect failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_CloseOriginProtect(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		v := values{
			"origin-protect": "",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.CloseOriginProtect(context.Background())
	if err != nil {
		t.Fatalf("CI.CloseOriginProtect returned error: %v", err)
	}
}

func TestCIService_PicTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "detect-label",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<RecognitionResult>
  <Labels>
    <Confidence>88</Confidence>
    <Name>玩具</Name>
  </Labels>
  <Labels>
    <Confidence>87</Confidence>
    <Name>毛绒玩具</Name>
  </Labels>
  <Labels>
    <Confidence>77</Confidence>
    <Name>泰迪熊</Name>
  </Labels>
  <Labels>
    <Confidence>74</Confidence>
    <Name>纺织品</Name>
  </Labels>
  <Labels>
    <Confidence>15</Confidence>
    <Name>艺术</Name>
  </Labels>
</RecognitionResult>
`)
	})

	want := &PicTagResult{
		XMLName: xml.Name{Local: "RecognitionResult"},
		Labels: []PicTag{
			PicTag{
				Confidence: 88,
				Name:       "玩具",
			},
			PicTag{
				Confidence: 87,
				Name:       "毛绒玩具",
			},
			PicTag{
				Confidence: 77,
				Name:       "泰迪熊",
			},
			PicTag{
				Confidence: 74,
				Name:       "纺织品",
			},
			PicTag{
				Confidence: 15,
				Name:       "艺术",
			},
		},
	}

	res, _, err := client.CI.PicTag(context.Background(), "test.jpg")
	if err != nil {
		t.Fatalf("CI.PicTag returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.PicTag failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_DetectFace(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":   "DetectFace",
			"max-face-num": "1",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <ImageWidth>616</ImageWidth>
  <ImageHeight>442</ImageHeight>
  <FaceModelVersion>3.0</FaceModelVersion>
  <RequestId>xxx</RequestId>
  <FaceInfos>
    <X>312</X>
    <Y>-5</Y>
    <Width>117</Width>
    <Height>173</Height>
  </FaceInfos>
  <FaceInfos>
    <X>600</X>
    <Y>-5</Y>
    <Width>117</Width>
    <Height>173</Height>
  </FaceInfos>
</Response>
`)
	})

	opt := &DetectFaceOptions{
		MaxFaceNum: 1,
	}

	want := &DetectFaceResult{
		XMLName:          xml.Name{Local: "Response"},
		ImageWidth:       616,
		ImageHeight:      442,
		FaceModelVersion: "3.0",
		RequestId:        "xxx",
		FaceInfos: []FaceInfos{
			FaceInfos{
				X:      312,
				Y:      -5,
				Width:  117,
				Height: 173,
			},
			FaceInfos{
				X:      600,
				Y:      -5,
				Width:  117,
				Height: 173,
			},
		},
	}

	res, _, err := client.CI.DetectFace(context.Background(), "test.jpg", opt)
	if err != nil {
		t.Fatalf("CI.DetectFace returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.DetectFace failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_FaceEffect(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":   "face-effect",
			"type":         "face-beautify",
			"whitening":    "70",
			"smoothing":    "80",
			"faceLifting":  "70",
			"eyeEnlarging": "70",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <ResultImage>xxx</ResultImage>
</Response>
`)
	})

	opt := &FaceEffectOptions{
		Type:         "face-beautify",
		Whitening:    70,
		Smoothing:    80,
		FaceLifting:  70,
		EyeEnlarging: 70,
	}

	want := &FaceEffectResult{
		XMLName:     xml.Name{Local: "Response"},
		ResultImage: "xxx",
	}

	res, _, err := client.CI.FaceEffect(context.Background(), "test.jpg", opt)
	if err != nil {
		t.Fatalf("CI.FaceEffect returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.FaceEffect failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_IdCardOCRWhenCloud(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "IDCardOCR",
			"CardSide":   "FRONT",
			"Config":     `{"CropIdCard":true}`,
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <IdInfo>
    <Name>李明</Name>
    <Sex>男</Sex>
    <Nation>汉</Nation>
    <Birth>1987/1/1</Birth>
    <Address>北京市xx大楼</Address>
    <IdNum>440524198701010014</IdNum>
  </IdInfo>
  <AdvancedInfo>
    <IdCard>xxx</IdCard>
    <Portrait>yyy</Portrait>
  </AdvancedInfo>
</Response>
`)
	})

	opt := &IdCardOCROptions{
		CardSide: "FRONT",
		Config: &IdCardOCROptionsConfig{
			CropIdCard: true,
		},
	}

	want := &IdCardOCRResult{
		XMLName: xml.Name{Local: "Response"},
		IdInfo: &IdCardInfo{
			Name:    "李明",
			Sex:     "男",
			Nation:  "汉",
			Birth:   "1987/1/1",
			Address: "北京市xx大楼",
			IdNum:   "440524198701010014",
		},
		AdvancedInfo: &IdCardAdvancedInfo{
			IdCard:   "xxx",
			Portrait: "yyy",
		},
	}

	res, _, err := client.CI.IdCardOCRWhenCloud(context.Background(), "test.jpg", opt)
	if err != nil {
		t.Fatalf("CI.IdCardOCRWhenCloud returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.IdCardOCRWhenCloud failed, return:%+v, want:%+v", res, want)
	}
}

func TestObjectService_IdCardOCRWhenUpload(t *testing.T) {
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
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		bs, _ := ioutil.ReadAll(r.Body)
		crc := crc64.Update(0, tb, bs)
		if !reflect.DeepEqual(bs, b) {
			t.Errorf("Object.Put request body Error")
		}
		if !reflect.DeepEqual(crc, realcrc) {
			t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
		}
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))

		w.Header().Add("Content-Type", "application/xml")
		fmt.Fprint(w, `<Response>
  <IdInfo>
    <Name>李明</Name>
    <Sex>男</Sex>
    <Nation>汉</Nation>
    <Birth>1987/1/1</Birth>
    <Address>北京市xx大楼</Address>
    <IdNum>440524198701010014</IdNum>
  </IdInfo>
  <AdvancedInfo>
    <IdCard>xxx</IdCard>
    <Portrait>yyy</Portrait>
  </AdvancedInfo>
</Response>
`)
	})

	qopt := &IdCardOCROptions{
		CardSide: "FRONT",
		Config: &IdCardOCROptionsConfig{
			CropIdCard: true,
		},
	}

	hopt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
			Listener:    &DefaultProgressListener{},
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}

	want := &IdCardOCRResult{
		XMLName: xml.Name{Local: "Response"},
		IdInfo: &IdCardInfo{
			Name:    "李明",
			Sex:     "男",
			Nation:  "汉",
			Birth:   "1987/1/1",
			Address: "北京市xx大楼",
			IdNum:   "440524198701010014",
		},
		AdvancedInfo: &IdCardAdvancedInfo{
			IdCard:   "xxx",
			Portrait: "yyy",
		},
	}

	res, _, err := client.CI.IdCardOCRWhenUpload(context.Background(), name, filePath, qopt, hopt)
	if err != nil {
		t.Fatalf("CI.IdCardOCRWhenUpload returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.IdCardOCRWhenUpload failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_GetLiveCode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "GetLiveCode",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <LiveCode>0521</LiveCode>
</Response>
`)
	})

	want := &GetLiveCodeResult{
		XMLName:  xml.Name{Local: "Response"},
		LiveCode: "0521",
	}

	res, _, err := client.CI.GetLiveCode(context.Background())
	if err != nil {
		t.Fatalf("CI.GetLiveCode returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetLiveCode failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_GetActionSequence(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "GetActionSequence",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <ActionSequence>2,1</ActionSequence>
</Response>
`)
	})

	want := &GetActionSequenceResult{
		XMLName:        xml.Name{Local: "Response"},
		ActionSequence: "2,1",
	}

	res, _, err := client.CI.GetActionSequence(context.Background())
	if err != nil {
		t.Fatalf("CI.GetActionSequence returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.GetActionSequence failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_LivenessRecognitionWhenCloud(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":   "LivenessRecognition",
			"IdCard":       "11204416541220243X",
			"Name":         "韦小宝",
			"LivenessType": "SILENT",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, `<Response>
  <BestFrameBase64>/9j/4AAQSkZJRgABAQAAAQA</BestFrameBase64>
  <Sim>100</Sim>
</Response>
`)
	})

	opt := &LivenessRecognitionOptions{
		IdCard:       "11204416541220243X",
		Name:         "韦小宝",
		LivenessType: "SILENT",
	}

	want := &LivenessRecognitionResult{
		XMLName:         xml.Name{Local: "Response"},
		BestFrameBase64: "/9j/4AAQSkZJRgABAQAAAQA",
		Sim:             100,
	}

	res, _, err := client.CI.LivenessRecognitionWhenCloud(context.Background(), "test.jpg", opt)
	if err != nil {
		t.Fatalf("CI.LivenessRecognitionWhenCloud returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.LivenessRecognitionWhenCloud failed, return:%+v, want:%+v", res, want)
	}
}

func TestObjectService_LivenessRecognitionWhenUpload(t *testing.T) {
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
	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		bs, _ := ioutil.ReadAll(r.Body)
		crc := crc64.Update(0, tb, bs)
		if !reflect.DeepEqual(bs, b) {
			t.Errorf("Object.Put request body Error")
		}
		if !reflect.DeepEqual(crc, realcrc) {
			t.Errorf("Object.Put crc: %v, want: %v", crc, realcrc)
		}
		w.Header().Add("x-cos-hash-crc64ecma", strconv.FormatUint(crc, 10))

		w.Header().Add("Content-Type", "application/xml")
		fmt.Fprint(w, `<Response>
  <BestFrameBase64>/9j/4AAQSkZJRgABAQAAAQA</BestFrameBase64>
  <Sim>100</Sim>
</Response>
`)
	})

	qopt := &LivenessRecognitionOptions{
		IdCard:       "11204416541220243X",
		Name:         "韦小宝",
		LivenessType: "SILENT",
	}

	hopt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
			Listener:    &DefaultProgressListener{},
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}

	want := &LivenessRecognitionResult{
		XMLName:         xml.Name{Local: "Response"},
		BestFrameBase64: "/9j/4AAQSkZJRgABAQAAAQA",
		Sim:             100,
	}

	res, _, err := client.CI.LivenessRecognitionWhenUpload(context.Background(), name, filePath, qopt, hopt)
	if err != nil {
		t.Fatalf("CI.LivenessRecognitionWhenUpload returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.LivenessRecognitionWhenUpload failed, return:%+v, want:%+v", res, want)
	}
}

func TestCIService_GoodsMatting(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "GoodsMatting",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.GoodsMatting(context.Background(), "test.jpg")
	if err != nil {
		t.Fatalf("CI.GoodsMatting returned error: %v", err)
	}
}

func TestCIService_GoodsMattingWithOpt(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.jpg", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "GoodsMatting",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.GoodsMattingWithOpt(context.Background(), "test.jpg", nil)
	if err != nil {
		t.Fatalf("CI.GoodsMattingWithOpt returned error: %v", err)
	}
}

func TestCIService_PutPosterproductionTemplatet(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Input><Object>input/sample.psd</Object></Input><Name>test</Name></Request>"

	mux.HandleFunc("/posterproduction/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	PosterproductionTemplate := &PosterproductionTemplateOptions{
		Input: &PosterproductionInput{
			Object: "input/sample.psd",
		},
		Name: "test",
	}

	_, _, err := client.CI.PutPosterproductionTemplate(context.Background(), PosterproductionTemplate)
	if err != nil {
		t.Fatalf("CI.PutPosterproductionTemplate returned error: %v", err)
	}
}

func TestCIService_GetPosterproductionTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "1234567890"

	mux.HandleFunc("/posterproduction/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.GetPosterproductionTemplate(context.Background(), tplId)
	if err != nil {
		t.Fatalf("CI.GetPosterproductionTemplate returned error: %v", err)
	}
}

func TestCIService_GetPosterproductionTemplates(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/posterproduction/template/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"pageNumber": "1",
			"pageSize":   "10",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribePosterproductionTemplateOptions{
		PageNumber: 1,
		PageSize:   10,
	}

	_, _, err := client.CI.GetPosterproductionTemplates(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.GetPosterproductionTemplates returned error: %v", err)
	}
}

func TestCIService_GetOriginImage(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "originImage",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.GetOriginImage(context.Background(), key)
	if err != nil {
		t.Fatalf("CI.GetOriginImage returned error: %v", err)
	}
}

func TestCIService_GetAIImageColoring(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AIImageColoring",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.GetAIImageColoring(context.Background(), key)
	if err != nil {
		t.Fatalf("CI.GetAIImageColoring returned error: %v", err)
	}
}

func TestCIService_GetAISuperResolution(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AISuperResolution",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.GetAISuperResolution(context.Background(), key)
	if err != nil {
		t.Fatalf("CI.GetAISuperResolution returned error: %v", err)
	}
}

func TestCIService_GetAIEnhanceImage(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AIEnhanceImage",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.GetAIEnhanceImage(context.Background(), key)
	if err != nil {
		t.Fatalf("CI.GetAIEnhanceImage returned error: %v", err)
	}
}

func TestCIService_GetAIImageCrop(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":   "AIImageCrop",
			"width":        "128",
			"height":       "96",
			"fixed":        "1",
			"ignore-error": "1",
		}
		testFormValues(t, r, v)
	})

	opt := &AIImageCropOptions{
		Width:       128,
		Height:      96,
		Fixed:       1,
		IgnoreError: 1,
	}
	_, err := client.CI.GetAIImageCrop(context.Background(), key, opt)
	if err != nil {
		t.Fatalf("CI.GetAIImageCrop returned error: %v", err)
	}
}

func TestCIService_GetAutoTranslationBlock(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AutoTranslationBlock",
			"InputText":  "你好",
			"SourceLang": "zh",
			"TargetLang": "en",
			"TextDomain": "ecommerce",
			"TextStyle":  "sentence",
		}
		testFormValues(t, r, v)
	})

	opt := &AutoTranslationBlockOptions{
		InputText:  "你好",
		SourceLang: "zh",
		TargetLang: "en",
		TextDomain: "ecommerce",
		TextStyle:  "sentence",
	}
	_, _, err := client.CI.GetAutoTranslationBlock(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.GetAutoTranslationBlock returned error: %v", err)
	}
}

func TestCIService_GetImageRepair(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "ImageRepair",
			"MaskPic":    "xxx",
		}
		testFormValues(t, r, v)
	})

	opt := &ImageRepairOptions{
		MaskPic: "xxx",
	}
	_, err := client.CI.GetImageRepair(context.Background(), key, opt)
	if err != nil {
		t.Fatalf("CI.GetImageRepair returned error: %v", err)
	}
}

func TestCIService_GetRecognizeLogo(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":   "RecognizeLogo",
			"ignore-error": "1",
		}
		testFormValues(t, r, v)
	})

	opt := &RecognizeLogoOptions{
		IgnoreError: 1,
	}
	_, _, err := client.CI.GetRecognizeLogo(context.Background(), key, opt)
	if err != nil {
		t.Fatalf("CI.GetRecognizeLogo returned error: %v", err)
	}
}

func TestCIService_GetAssessQuality(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "AssessQuality",
		}
		testFormValues(t, r, v)
	})

	_, _, err := client.CI.GetAssessQuality(context.Background(), key)
	if err != nil {
		t.Fatalf("CI.GetAssessQuality returned error: %v", err)
	}
}

func TestCIService_TDCRefresh(t *testing.T) {
	setup()
	defer teardown()

	key := "pic/cup.jpeg"
	mux.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := values{
			"TDCRefresh": "",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.TDCRefresh(context.Background(), key)
	if err != nil {
		t.Fatalf("CI.TDCRefresh returned error: %v", err)
	}
}
