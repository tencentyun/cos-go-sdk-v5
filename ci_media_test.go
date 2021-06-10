package cos

import (
	"context"
	"net/http"
	"testing"
)

func TestCIService_CreateMediaJobs(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Animation</Tag><Input><Object>test.mp4</Object></Input>" +
		"<Operation><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test-trans.gif</Object></Output>" +
		"<TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaJobsOptions{
		Tag: "Animation",
		Input: &JobInput{
			Object: "test.mp4",
		},
		Operation: &MediaProcessJobOperation{
			Output: &JobOutput{
				Region: "ap-beijing",
				Bucket: "abc-1250000000",
				Object: "test-trans.gif",
			},
			TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreateMediaJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaJobs returned errors: %v", err)
	}
}

func TestCIService_DescribeMediaJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "jabcsdssfeipplsdfwe"
	mux.HandleFunc("/jobs/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeMediaJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.DescribeMediaJob returned error: %v", err)
	}
}
