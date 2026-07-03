package cos

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCIService_CreateDocProcessJobs(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<Request><Tag>DocProcess</Tag><Input><Object>1.doc</Object></Input>" +
		"<Operation><Output><Region>ap-chongqing</Region><Bucket>examplebucket-1250000000</Bucket>" +
		"<Object>big/test-${Number}</Object></Output><DocProcess>" +
		"<TgtType>png</TgtType><StartPage>1</StartPage><EndPage>-1</EndPage>" +
		"<ImageParams>watermark/1/image/aHR0cDovL3Rlc3QwMDUtMTI1MTcwNDcwOC5jb3MuYXAtY2hvbmdxaW5nLm15cWNsb3VkLmNvbS8xLmpwZw==/gravity/southeast</ImageParams>" +
		"</DocProcess></Operation><QueueId>p532fdead78444e649e1a4467c1cd19d3</QueueId></Request>"

	mux.HandleFunc("/doc_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	createJobOpt := &CreateDocProcessJobsOptions{
		Tag: "DocProcess",
		Input: &DocProcessJobInput{
			Object: "1.doc",
		},
		Operation: &DocProcessJobOperation{
			Output: &DocProcessJobOutput{
				Region: "ap-chongqing",
				Object: "big/test-${Number}",
				Bucket: "examplebucket-1250000000",
			},
			DocProcess: &DocProcessJobDocProcess{
				TgtType:     "png",
				StartPage:   1,
				EndPage:     -1,
				ImageParams: "watermark/1/image/aHR0cDovL3Rlc3QwMDUtMTI1MTcwNDcwOC5jb3MuYXAtY2hvbmdxaW5nLm15cWNsb3VkLmNvbS8xLmpwZw==/gravity/southeast",
			},
		},
		QueueId: "p532fdead78444e649e1a4467c1cd19d3",
	}

	_, _, err := client.CI.CreateDocProcessJobs(context.Background(), createJobOpt)
	if err != nil {
		t.Fatalf("CI.CreateDocProcessJobs returned error: %v", err)
	}
}

func TestCIService_DescribeDocProcessJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "d13cfd584cd9011ea820b597ad1785a2f"
	mux.HandleFunc("/doc_jobs"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeDocProcessJob(context.Background(), jobID)

	if err != nil {
		t.Fatalf("CI.DescribeDocProcessJob returned error: %v", err)
	}
}

func TestCIService_DescribeDocProcessJobs(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/doc_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueId": "QueueID",
			"tag":     "DocProcess",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeDocProcessJobsOptions{
		QueueId: "QueueID",
		Tag:     "DocProcess",
	}

	_, _, err := client.CI.DescribeDocProcessJobs(context.Background(), opt)

	if err != nil {
		t.Fatalf("CI.DescribeDocProcessJobs returned error: %v", err)
	}
}

func TestCIService_DescribeDocProcessQueues(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/docqueue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"pageNumber": "1",
			"pageSize":   "2",
			"queueIds":   "p111a8dd208104ce3b11c78398f658ca8,p4318f85d2aa14c43b1dba6f9b78be9b3,aacb2bb066e9c4478834d4196e76c49d3",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeDocProcessQueuesOptions{
		QueueIds:   "p111a8dd208104ce3b11c78398f658ca8,p4318f85d2aa14c43b1dba6f9b78be9b3,aacb2bb066e9c4478834d4196e76c49d3",
		PageNumber: 1,
		PageSize:   2,
	}

	_, _, err := client.CI.DescribeDocProcessQueues(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDocProcessQueues returned error: %v", err)
	}
}

func TestCIService_UpdateDocProcessQueue(t *testing.T) {
	setup()
	defer teardown()

	queueID := "p2505d57bdf4c4329804b58a6a5fb1572"
	wantBody := "<Request><Name>markjrzhang4</Name><QueueID>p2505d57bdf4c4329804b58a6a5fb1572</QueueID>" +
		"<State>Active</State>" +
		"<NotifyConfig><Url>http://google.com/</Url><State>On</State>" +
		"<Type>Url</Type><Event>TransCodingFinish</Event>" +
		"</NotifyConfig></Request>"

	mux.HandleFunc("/docqueue/"+queueID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &UpdateDocProcessQueueOptions{
		Name:    "markjrzhang4",
		QueueID: queueID,
		State:   "Active",
		NotifyConfig: &DocProcessQueueNotifyConfig{
			Url:   "http://google.com/",
			State: "On",
			Type:  "Url",
			Event: "TransCodingFinish",
		},
	}

	_, _, err := client.CI.UpdateDocProcessQueue(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDocProcessQueues returned error: %v", err)
	}
}

func TestCIService_DescribeDocProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/docbucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"pageNumber": "1",
			"pageSize":   "2",
			"regions":    "ap-shanghai",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeDocProcessBucketsOptions{
		Regions:    "ap-shanghai",
		PageNumber: 1,
		PageSize:   2,
	}

	_, _, err := client.CI.DescribeDocProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDocProcessBuckets returned error: %v", err)
	}
}

func TestCIService_DocPreview(t *testing.T) {
	setup()
	defer teardown()

	name := "sample.pdf"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":  "doc-preview",
			"page":        "1",
			"ImageParams": "imageMogr2/thumbnail/!50p|watermark/2/text/5pWw5o2u5LiH6LGh/fill/I0ZGRkZGRg==/fontsize/30/dx/20/dy/20",
		}
		testFormValues(t, r, v)
	})

	opt := &DocPreviewOptions{
		Page:        1,
		ImageParams: "imageMogr2/thumbnail/!50p|watermark/2/text/5pWw5o2u5LiH6LGh/fill/I0ZGRkZGRg==/fontsize/30/dx/20/dy/20",
	}

	_, err := client.CI.DocPreview(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.DocPreview returned error: %v", err)
	}
}

// ============================================================================
// 测试用例：DocPreview - srcType=ofd 支持
// ============================================================================

// 覆盖：ci_doc.go DocPreviewOptions.SrcType 走 srcType=ofd 分支
// 断言请求 URL 含 ci-process=doc-preview & srcType=ofd & dstType=jpg & page=1
func TestCIService_DocPreview_OFD(t *testing.T) {
	setup()
	defer teardown()

	name := "sample.ofd"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "doc-preview",
			"srcType":    "ofd",
			"dstType":    "jpg",
			"page":       "1",
		}
		testFormValues(t, r, v)
	})

	opt := &DocPreviewOptions{
		SrcType: "ofd",
		DstType: "jpg",
		Page:    1,
	}

	_, err := client.CI.DocPreview(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.DocPreview returned error: %v", err)
	}
}

func TestCIService_CIDocCompare(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/doccompare", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"object":      "doc/1.docx",
			"comparePath": "doc/2.docx",
		}
		testFormValues(t, r, v)
	})

	opt := &CIDocCompareOptions{
		Object:      "doc/1.docx",
		ComparePath: "doc/2.docx",
	}

	_, _, err := client.CI.CIDocCompare(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DocPreview returned error: %v", err)
	}
}

func TestCIService_DocPreviewHTML(t *testing.T) {
	setup()
	defer teardown()

	name := "sample.pdf"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "doc-preview",
			"dstType":    "html",
			"srcType":    "pdf",
			"htmlParams": `{"commonOptions":{"isShowTopArea":true,"isShowHeader":true,"isBrowserViewFullscreen":false,"isIframeViewFullscreen":false}}`,
		}
		testFormValues(t, r, v)
	})

	opt := &DocPreviewHTMLOptions{
		DstType: "html",
		SrcType: "pdf",
		HtmlParams: &HtmlParams{
			CommonOptions: &HtmlCommonParams{
				IsShowTopArea: true,
				IsShowHeader:  true,
			},
		},
	}

	_, err := client.CI.DocPreviewHTML(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.DocPreviewHTML returned error: %v", err)
	}
}

func TestCIService_CIDocCompareResultWrite(t *testing.T) {
	setup()
	defer teardown()
	result := &CIDocCompareResult{
		Code:       "success",
		ETag:       "1234567890",
		Msg:        "success",
		ResultPath: "abc/abc.jpg",
	}
	slice := make([]byte, 5)
	result.Write(slice)
}

func TestCIService_CreateDocProcessBucket(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/docbucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	opt := &CreateDocProcessBucketOptions{}

	_, _, err := client.CI.CreateDocProcessBucket(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateDocProcessBucket returned error: %v", err)
	}
}

// ============================================================================
// 测试用例：DocProcessJobInput.Url — 文档转码任务支持第三方 URL 输入源
// ============================================================================

// 覆盖：ci_doc.go DocProcessJobInput.Url（新增字段）+ XML 正常序列化
// 断言 XML 含 <Url> 且不含 <Object>（因为 Object 空，omitempty 应生效）
func TestDocProcessJobInput_Url_MarshalXML(t *testing.T) {
	in := &DocProcessJobInput{
		Url: "https://x/a.docx",
	}
	out, err := xml.Marshal(in)
	if err != nil {
		t.Fatalf("xml.Marshal DocProcessJobInput returned error: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "<Url>https://x/a.docx</Url>") {
		t.Errorf("XML should contain <Url> tag, got: %s", got)
	}
	if strings.Contains(got, "<Object>") {
		t.Errorf("XML should NOT contain <Object> tag when Object is empty, got: %s", got)
	}
}

// 覆盖：ci_doc.go DocProcessJobInput.Object 老用法回归 + Url omitempty 分支
// 断言 XML 含 <Object> 且不含 <Url>
func TestDocProcessJobInput_Object_MarshalXML(t *testing.T) {
	in := &DocProcessJobInput{
		Object: "a.docx",
	}
	out, err := xml.Marshal(in)
	if err != nil {
		t.Fatalf("xml.Marshal DocProcessJobInput returned error: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "<Object>a.docx</Object>") {
		t.Errorf("XML should contain <Object> tag, got: %s", got)
	}
	if strings.Contains(got, "<Url>") {
		t.Errorf("XML should NOT contain <Url> tag when Url is empty, got: %s", got)
	}
}

// 覆盖：ci_doc.go 端到端 — 通过 CreateDocProcessJobs 发送含 Url 的请求，
// 断言 HTTP 请求 Body 包含 <Url>https://example.com/a.docx</Url>
func TestCIService_CreateDocProcessJobs_WithUrl(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/doc_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body error: %v", err)
		}
		if !strings.Contains(string(body), "<Url>https://example.com/a.docx</Url>") {
			t.Errorf("request body should contain <Url> tag, got: %s", string(body))
		}
	})

	createJobOpt := &CreateDocProcessJobsOptions{
		Tag: "DocProcess",
		Input: &DocProcessJobInput{
			Url: "https://example.com/a.docx",
		},
		Operation: &DocProcessJobOperation{
			Output: &DocProcessJobOutput{
				Region: "ap-chongqing",
				Object: "big/test-${Number}",
				Bucket: "examplebucket-1250000000",
			},
			DocProcess: &DocProcessJobDocProcess{
				TgtType:   "png",
				StartPage: 1,
				EndPage:   -1,
			},
		},
		QueueId: "p532fdead78444e649e1a4467c1cd19d3",
	}

	_, _, err := client.CI.CreateDocProcessJobs(context.Background(), createJobOpt)
	if err != nil {
		t.Fatalf("CI.CreateDocProcessJobs returned error: %v", err)
	}
}

// ============================================================================
// 测试用例：DocWatermark - 平铺三字段（Batch/HorizontalSpacing/VerticalSpacing）
// ============================================================================

// 覆盖：ci_doc.go DocWatermark 新增三字段 XML 序列化
// 断言：赋值后三个字段都出现在 XML 输出
func TestDocWatermark_Batch_MarshalXML(t *testing.T) {
	wm := DocWatermark{
		SrcType:           "pdf",
		Type:              "Text",
		Batch:             "1",
		HorizontalSpacing: "20",
		VerticalSpacing:   "30",
	}
	buf, err := xml.Marshal(&wm)
	if err != nil {
		t.Fatalf("xml.Marshal error: %v", err)
	}
	got := string(buf)
	for _, want := range []string{
		"<Batch>1</Batch>",
		"<HorizontalSpacing>20</HorizontalSpacing>",
		"<VerticalSpacing>30</VerticalSpacing>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expect XML contains %q, got: %s", want, got)
		}
	}
}

// 覆盖：ci_doc.go DocWatermark 新字段 omitempty 行为
// 断言：空值时三个新字段都不出现在 XML
func TestDocWatermark_OmitEmpty(t *testing.T) {
	wm := DocWatermark{SrcType: "pdf", Type: "Text"}
	buf, err := xml.Marshal(&wm)
	if err != nil {
		t.Fatalf("xml.Marshal error: %v", err)
	}
	got := string(buf)
	for _, notWant := range []string{"<Batch>", "<HorizontalSpacing>", "<VerticalSpacing>"} {
		if strings.Contains(got, notWant) {
			t.Errorf("expect XML NOT contains %q, got: %s", notWant, got)
		}
	}
}

// 覆盖：ci_doc.go 端到端 — CreateDocProcessJobs 请求 Body 含 DocWatermark 新字段
// 断言：Body 内 <DocWatermark> 包含 Batch/HorizontalSpacing/VerticalSpacing
func TestCIService_CreateDocProcessJobs_WithDocWatermarkBatch(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/doc_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body error: %v", err)
		}
		s := string(body)
		for _, want := range []string{
			"<Batch>1</Batch>",
			"<HorizontalSpacing>20</HorizontalSpacing>",
			"<VerticalSpacing>30</VerticalSpacing>",
		} {
			if !strings.Contains(s, want) {
				t.Errorf("body should contain %q, got: %s", want, s)
			}
		}
	})

	opt := &CreateDocProcessJobsOptions{
		Tag: "DocProcess",
		Input: &DocProcessJobInput{
			Object: "1.pdf",
		},
		Operation: &DocProcessJobOperation{
			Output: &DocProcessJobOutput{
				Region: "ap-chongqing",
				Object: "out/1.png",
				Bucket: "examplebucket-1250000000",
			},
			DocWatermark: &DocWatermark{
				SrcType:           "pdf",
				Type:              "Text",
				Batch:             "1",
				HorizontalSpacing: "20",
				VerticalSpacing:   "30",
			},
		},
		QueueId: "p532fdead78444e649e1a4467c1cd19d3",
	}

	_, _, err := client.CI.CreateDocProcessJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CreateDocProcessJobs error: %v", err)
	}
}

// ============================================================================
// 测试用例：DocPreview - PDF 加水印场景（dstType=watermark，复用 DocPreview 方法）
// ============================================================================

// 覆盖：ci_doc.go DocPreviewOptions 新增水印字段（Type/Text/Batch/HorizontalSpacing/VerticalSpacing）
// 断言：URL query 含 ci-process=doc-preview & dstType=watermark & 平铺三参数
func TestCIService_DocPreview_Watermark_BatchMode(t *testing.T) {
	setup()
	defer teardown()

	name := "sample.pdf"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":         "doc-preview",
			"srcType":            "pdf",
			"dstType":            "watermark",
			"type":               "Text",
			"text":               "confidential",
			"batch":              "1",
			"horizontal-spacing": "20",
			"vertical-spacing":   "30",
		}
		testFormValues(t, r, v)
	})

	opt := &DocPreviewOptions{
		SrcType:           "pdf",
		DstType:           "watermark",
		Type:              "Text",
		Text:              "confidential",
		Batch:             1,
		HorizontalSpacing: 20,
		VerticalSpacing:   30,
	}

	_, err := client.CI.DocPreview(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.DocPreview returned error: %v", err)
	}
}

// 覆盖：ci_doc.go DocPreviewOptions 图片水印 + Pos/Dx/Dy/Password 字段拼接
// 断言：URL query 含 image / pos / dx / dy / password 五个字段
func TestCIService_DocPreview_Watermark_ImageAndPos(t *testing.T) {
	setup()
	defer teardown()

	name := "sample.pdf"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "doc-preview",
			"dstType":    "watermark",
			"type":       "Image",
			"image":      "https://example.com/wm.png",
			"pos":        "TopRight",
			"dx":         "10",
			"dy":         "15",
			"password":   "p@ss",
		}
		testFormValues(t, r, v)
	})

	opt := &DocPreviewOptions{
		DstType:  "watermark",
		Type:     "Image",
		Image:    "https://example.com/wm.png",
		Pos:      "TopRight",
		Dx:       10,
		Dy:       15,
		Password: "p@ss",
	}

	_, err := client.CI.DocPreview(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.DocPreview returned error: %v", err)
	}
}
