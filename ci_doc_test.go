package cos

import (
	"context"
	"net/http"
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
