package cos

import (
	"context"
	"net/http"
	"testing"
)

func TestCIService_CreateDataSet(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"TemplateId\":\"Official:COSBasicMeta\"}"

	mux.HandleFunc("/dataset/create", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	createJobOpt := &CreateDataSetOptions{
		DatasetName: "adataset",
		Description: "dataset test",
		TemplateId:  "Official:COSBasicMeta",
	}
	_, _, err := client.CI.CreateDataSet(context.Background(), createJobOpt)
	if err != nil {
		t.Fatalf("CI.CreateDataSet returned error: %v", err)
	}
}

func TestCIService_UpdateDataset(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"TemplateId\":\"Official:COSBasicMeta\"}"

	mux.HandleFunc("/dataset/update", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	createJobOpt := &UpdateDatasetOptions{
		DatasetName: "adataset",
		Description: "dataset test",
		TemplateId:  "Official:COSBasicMeta",
	}
	_, _, err := client.CI.UpdateDataset(context.Background(), createJobOpt)
	if err != nil {
		t.Fatalf("CI.UpdateDataset returned error: %v", err)
	}
}

func TestCIService_DescribeDatasets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/datasets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
	})

	opt := &DescribeDatasetsOptions{
		MaxResults: 100,
	}
	_, _, err := client.CI.DescribeDatasets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDatasets returned error: %v", err)
	}
}

func TestCIService_DeleteDataset(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\"}"

	mux.HandleFunc("/dataset", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &DeleteDatasetOptions{
		DatasetName: "adataset",
	}
	_, _, err := client.CI.DeleteDataset(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DeleteDataset returned error: %v", err)
	}
}

func TestCIService_DescribeDataset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/dataset", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
	})

	opt := &DescribeDatasetOptions{
		DatasetName: "adataset",
		Statistics:  true,
	}
	_, _, err := client.CI.DescribeDataset(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDataset returned error: %v", err)
	}
}

func TestCIService_CreateFileMetaIndex(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"File\":{\"URI\":\"cos://test-125000000/12.gif\",\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"MediaType\":\"image\",\"contenttype\":\"image/gif\"}}"

	mux.HandleFunc("/filemeta/create", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &CreateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &File{
			URI:      "cos://test-125000000/12.gif",
			CustomID: "123",
			CustomLabels: &map[string]string{
				"age":   "18",
				"level": "18",
			},
			MediaType:   "image",
			ContentType: "image/gif",
		},
	}
	_, _, err := client.CI.CreateFileMetaIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateFileMetaIndex returned error: %v", err)
	}
}

func TestCIService_UpdateFileMetaIndex(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"File\":{\"URI\":\"cos://test-125000000/1.gif\",\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"MediaType\":\"video\",\"contenttype\":\"video/gif\"}}"

	mux.HandleFunc("/filemeta/update", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &UpdateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &File{
			URI:      "cos://test-125000000/1.gif",
			CustomID: "123",
			CustomLabels: &map[string]string{
				"age":   "18",
				"level": "18",
			},
			MediaType:   "video",
			ContentType: "video/gif",
		},
	}
	_, _, err := client.CI.UpdateFileMetaIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.UpdateFileMetaIndex returned error: %v", err)
	}
}

func TestCIService_DescribeFileMetaIndex(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/filemeta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
	})

	opt := &DescribeFileMetaIndexOptions{
		DatasetName: "adataset",
		Uri:         "cos://test-1250000000/1.gif",
	}
	_, _, err := client.CI.DescribeFileMetaIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeFileMetaIndex returned error: %v", err)
	}
}

func TestCIService_DeleteFileMetaIndex(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/filemeta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testHeader(t, r, "Content-Type", "application/json")
	})

	opt := &DeleteFileMetaIndexOptions{
		DatasetName: "adataset",
		Uri:         "cos://test1-1250000000/1.gif",
	}
	_, _, err := client.CI.DeleteFileMetaIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DeleteFileMetaIndex returned error: %v", err)
	}
}

func TestCIService_DatasetSimpleQuery(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Query\":{\"Operation\":\"eq\",\"Field\":\"ContentType\",\"Value\":\"image/gif\"}}"

	mux.HandleFunc("/datasetquery/simple", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &DatasetSimpleQueryOptions{
		DatasetName: "adataset",
		Query: &Query{
			Operation: "eq",
			Field:     "ContentType",
			Value:     "image/gif",
		},
	}
	_, _, err := client.CI.DatasetSimpleQuery(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DatasetSimpleQuery returned error: %v", err)
	}
}

func TestCIService_DatasetSimpleQueryAggregations(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Aggregations\":[{\"Field\":\"ContentType\",\"Operation\":\"group\"}]}"

	mux.HandleFunc("/datasetquery/simple", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &DatasetSimpleQueryOptions{
		DatasetName:  "adataset",
		Aggregations: []*Aggregation{},
	}
	opt.Aggregations = append(opt.Aggregations, &Aggregation{
		Field:     "ContentType",
		Operation: "group",
	})
	_, _, err := client.CI.DatasetSimpleQuery(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DatasetSimpleQuery returned error: %v", err)
	}
}

func TestCIService_CreateDatasetBinding(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

	mux.HandleFunc("/datasetbinding/create", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &CreateDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
	}
	_, _, err := client.CI.CreateDatasetBinding(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateDatasetBinding returned error: %v", err)
	}
}

func TestCIService_DescribeDatasetBinding(t *testing.T) {
	setup()
	defer teardown()
	// wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

	mux.HandleFunc("/datasetbinding", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
		// testBody(t, r, wantBody)
	})

	opt := &DescribeDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
	}
	_, _, err := client.CI.DescribeDatasetBinding(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDatasetBinding returned error: %v", err)
	}
}

func TestCIService_DescribeDatasetBindings(t *testing.T) {
	setup()
	defer teardown()
	// wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

	mux.HandleFunc("/datasetbindings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
		// testBody(t, r, wantBody)
	})

	opt := &DescribeDatasetBindingsOptions{
		DatasetName: "adataset",
		// MaxResults: 3,
	}
	_, _, err := client.CI.DescribeDatasetBindings(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDatasetBindings returned error: %v", err)
	}
}

func TestCIService_DeleteDatasetBinding(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

	mux.HandleFunc("/datasetbinding", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
	})

	opt := &DeleteDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
	}
	_, _, err := client.CI.DeleteDatasetBinding(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DeleteDatasetBinding returned error: %v", err)
	}
}
