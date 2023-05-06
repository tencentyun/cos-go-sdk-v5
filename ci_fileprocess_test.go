package cos

import (
	"context"
	"net/http"
	"testing"
)

func TestCIService_CreateFileProcessJob(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "<Request><Tag>FileHashCode</Tag><Input><Object>294028.zip</Object></Input>" +
		"<Operation><FileHashCodeConfig><Type>sha1</Type><AddToHeader>true</AddToHeader>" +
		"</FileHashCodeConfig></Operation><QueueId>pb6a88aead4dd4fa8bc953d4ca4e04430</QueueId></Request>"

	mux.HandleFunc("/file_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	createJobOpt := &FileProcessJobOptions{
		Tag: "FileHashCode",
		Input: &FileProcessInput{
			Object: "294028.zip",
		},
		Operation: &FileProcessJobOperation{
			FileHashCodeConfig: &FileHashCodeConfig{
				Type:        "sha1",
				AddToHeader: true,
			},
		},
		QueueId: "pb6a88aead4dd4fa8bc953d4ca4e04430",
	}

	_, _, err := client.CI.CreateFileProcessJob(context.Background(), createJobOpt)
	if err != nil {
		t.Fatalf("CI.CreateFileProcessJob returned error: %v", err)
	}
}

func TestCIService_DescribeFileProcessJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "f9640f1b0874211edb47e5fa2d6bd5e47"
	mux.HandleFunc("/file_jobs"+"/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeFileProcessJob(context.Background(), jobID)

	if err != nil {
		t.Fatalf("CI.DescribeFileProcessJob returned error: %v", err)
	}
}

func TestCIService_GetFileHash(t *testing.T) {
	setup()
	defer teardown()

	name := "sample.pdf"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process":  "filehash",
			"type":        "sha1",
			"addtoheader": "true",
		}
		testFormValues(t, r, v)
	})

	opt := &GetFileHashOptions{
		CIProcess:   "filehash",
		Type:        "sha1",
		AddToHeader: true,
	}

	_, _, err := client.CI.GetFileHash(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("CI.GetFileHash returned error: %v", err)
	}
}

func TestCIService_ZipPreview(t *testing.T) {
	setup()
	defer teardown()

	name := "test.zip"
	mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"ci-process": "zippreview",
		}
		testFormValues(t, r, v)
	})

	_, _, err := client.CI.ZipPreview(context.Background(), name)
	if err != nil {
		t.Fatalf("CI.ZipPreview returned error: %v", err)
	}
}
