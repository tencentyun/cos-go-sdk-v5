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

func TestCIService_DescribeMediaJobs(t *testing.T) {
	setup()
	defer teardown()

	queueId := "aaaaaaaaaaa"
	tag := "Animation"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueId": queueId,
			"tag":     tag,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaJobsOptions{
		QueueId: queueId,
		Tag:     tag,
	}

	_, _, err := client.CI.DescribeMediaJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaJobs returned error: %v", err)
	}
}

func TestCIService_DescribeMediaProcessQueues(t *testing.T) {
	setup()
	defer teardown()

	queueIds := "A,B,C"
	mux.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueIds": queueIds,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaProcessQueuesOptions{
		QueueIds: queueIds,
	}

	_, _, err := client.CI.DescribeMediaProcessQueues(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaProcessQueues returned error: %v", err)
	}
}

func TestCIService_UpdateMediaProcessQueue(t *testing.T) {
	setup()
	defer teardown()

	queueID := "p8eb46b8cc1a94bc09512d16c5c4f4d3a"
	wantBody := "<Request><Name>QueueName</Name><QueueID>" + queueID + "</QueueID><State>Active</State>" +
		"<NotifyConfig><Url>test.com</Url><State>On</State><Type>Url</Type><Event>TransCodingFinish</Event>" +
		"</NotifyConfig></Request>"
	mux.HandleFunc("/queue/"+queueID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &UpdateMediaProcessQueueOptions{
		Name:    "QueueName",
		QueueID: queueID,
		State:   "Active",
		NotifyConfig: &MediaProcessQueueNotifyConfig{
			Url:   "test.com",
			State: "On",
			Type:  "Url",
			Event: "TransCodingFinish",
		},
	}

	_, _, err := client.CI.UpdateMediaProcessQueue(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.UpdateMediaProcessQueue returned error: %v", err)
	}
}

func TestCIService_DescribeMediaProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	regions := "ap-shanghai,ap-gaungzhou"
	bucketName := "testbucket-1250000000"
	mux.HandleFunc("/mediabucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"regions":    regions,
			"bucketName": bucketName,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaProcessBucketsOptions{
		Regions:     regions,
		BucketName:  bucketName,
	}

	_, _, err := client.CI.DescribeMediaProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaProcessBuckets returned error: %v", err)
	}
}
