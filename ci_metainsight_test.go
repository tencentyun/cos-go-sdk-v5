package cos

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestCIService_CreateDataset(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"TemplateId\":\"Official:COSBasicMeta\"}"

	mux.HandleFunc("/dataset", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"Dataset\":{\"BindCount\":0,\"CreateTime\":\"2024-05-07T18:36:24.838341549+08:00\",\"DatasetName\":\"dataset\",\"Description\":\"dataset test\",\"FileCount\":0,\"TemplateId\":\"Official:COSBasicMeta\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-05-07T18:36:24.838341633+08:00\"},\"RequestId\":\"NjYzYTA0MjhfM2FiNjI5MWVfNTQyMl8yZjM4ZTI=\"}")
	})

	client.MetaInsight.CreateDataset(context.Background(), nil)

	createJobOpt := &CreateDatasetOptions{
		DatasetName: "adataset",
		Description: "dataset test",
		TemplateId:  "Official:COSBasicMeta",
	}
	_, _, err := client.MetaInsight.CreateDataset(context.Background(), createJobOpt)
	if err != nil {
		t.Fatalf("CI.CreateDataSet returned error: %v", err)
	}
}

func TestCIService_UpdateDataset(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"TemplateId\":\"Official:COSBasicMeta\"}"

	mux.HandleFunc("/dataset", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"Dataset\":{\"BindCount\":0,\"CreateTime\":\"2024-05-07T18:36:24.838341549+08:00\",\"DatasetName\":\"dataset\",\"Description\":\"dataset test\",\"FileCount\":0,\"TemplateId\":\"Official:COSBasicMeta\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-05-07T18:36:24.838341633+08:00\"},\"RequestId\":\"NjYzYTA0MjhfM2FiNjI5MWVfNTQyMl8yZjM4ZTI=\"}")
	})

	client.MetaInsight.UpdateDataset(context.Background(), nil)

	opt := &UpdateDatasetOptions{
		DatasetName: "adataset",
		Description: "dataset test",
		TemplateId:  "Official:COSBasicMeta",
	}
	_, _, err := client.MetaInsight.UpdateDataset(context.Background(), opt)
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
		fmt.Fprint(w, "{\"Datasets\":[{\"BindCount\":0,\"CreateTime\":\"2024-05-06T19:49:17.49197866+08:00\",\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"FileCount\":0,\"TemplateId\":\"Official:COSBasicMeta\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-05-06T19:49:17.49197874+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-05-07T18:36:24.838341549+08:00\",\"DatasetName\":\"dataset\",\"Description\":\"dataset test\",\"FileCount\":0,\"TemplateId\":\"Official:COSBasicMeta\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-05-07T18:36:24.838341633+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-28T16:58:29.972112328+08:00\",\"DatasetName\":\"test111\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-28T16:58:29.972112399+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-26T19:58:28.71611987+08:00\",\"DatasetName\":\"test11111\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-26T19:58:28.716119968+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-30T19:01:00.603324265+08:00\",\"DatasetName\":\"test111111\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-30T19:01:00.603324346+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-30T19:01:28.324249664+08:00\",\"DatasetName\":\"test1111111\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-30T19:01:28.324249747+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-30T19:06:22.973049681+08:00\",\"DatasetName\":\"test11111111\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-30T19:06:22.973049766+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-28T16:41:59.766417255+08:00\",\"DatasetName\":\"test111112\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-28T16:41:59.766417337+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-28T21:34:37.469900633+08:00\",\"DatasetName\":\"test11111222\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-28T21:34:37.469900718+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-28T21:35:03.76822133+08:00\",\"DatasetName\":\"test111112222\",\"Description\":\"数据集描述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-28T21:35:03.768221411+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-28T14:11:26.88710993+08:00\",\"DatasetName\":\"test12\",\"Description\":\"\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-28T14:11:26.887110011+08:00\"},{\"BindCount\":0,\"CreateTime\":\"2024-04-30T19:06:34.117713669+08:00\",\"DatasetName\":\"test1asfdsasdfafsdfa1111111\",\"Description\":\"数asdfa据集asdf描asdfad述\",\"FileCount\":0,\"TemplateId\":\"Official:Empty\",\"TotalFileSize\":0,\"UpdateTime\":\"2024-04-30T19:06:34.117713747+08:00\"}],\"NextToken\":\"\",\"RequestId\":\"NjYzYTFlNTVfNTc2ODk0MGJfNjZkM18zNmUyZTA=\"}")
	})

	opt := &DescribeDatasetsOptions{
		Maxresults: 100,
	}
	_, _, err := client.MetaInsight.DescribeDatasets(context.Background(), opt)
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
		fmt.Fprintf(w, "{\"Code\":\"InvalidArgument\",\"Message\":\"dataset not empty\",\"RequestId\":\"NjYzYjZiNmRfM2FiNjI5MWVfNTQyMl8zMzVkZTY=\",\"TraceId\":\"\"}")
	})

	client.MetaInsight.DeleteDataset(context.Background(), nil)

	opt := &DeleteDatasetOptions{
		DatasetName: "adataset",
	}
	_, _, err := client.MetaInsight.DeleteDataset(context.Background(), opt)
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
		fmt.Fprint(w, "{\"Dataset\":{\"BindCount\":0,\"CreateTime\":\"2024-05-06T19:49:17.49197866+08:00\",\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"FileCount\":1,\"TemplateId\":\"Official:COSBasicMeta\",\"TotalFileSize\":495199,\"UpdateTime\":\"2024-05-06T19:49:17.49197874+08:00\"},\"RequestId\":\"NjYzYjZjYzFfNjg2ODk0MGJfNzI0M18zMmUzMzE=\"}")
	})

	opt := &DescribeDatasetOptions{
		Datasetname: "adataset",
		Statistics:  true,
	}
	_, _, err := client.MetaInsight.DescribeDataset(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDataset returned error: %v", err)
	}
}

func TestCIService_CreateFileMetaIndex(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"File\":{\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"Key\":\"\",\"Value\":\"\",\"MediaType\":\"image\",\"ContentType\":\"image/gif\",\"URI\":\"cos://test-125000000/12.gif\",\"MaxFaceNum\":0,\"Persons\":null}}"

	mux.HandleFunc("/filemeta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"EventId\":\"wi78e458510d3511ef95635254008dc19b\",\"RequestId\":\"NjYzYjZlM2VfNjg2ODk0MGJfNzIyMV8zMzA4MWE=\"}")
	})

	client.MetaInsight.CreateFileMetaIndex(context.Background(), nil)

	opt := &CreateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &File{
			URI:      "cos://test-125000000/12.gif",
			CustomId: "123",
			CustomLabels: &map[string]string{
				"age":   "18",
				"level": "18",
			},
			MediaType:   "image",
			ContentType: "image/gif",
		},
	}
	_, _, err := client.MetaInsight.CreateFileMetaIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateFileMetaIndex returned error: %v", err)
	}
}

func TestCIService_UpdateFileMetaIndex(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Callback\":\"\",\"File\":{\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"Key\":\"\",\"Value\":\"\",\"MediaType\":\"video\",\"ContentType\":\"video/gif\",\"URI\":\"cos://test-125000000/1.gif\",\"MaxFaceNum\":0,\"Persons\":null}}"

	mux.HandleFunc("/filemeta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"EventId\":\"wi78e458510d3511ef95635254008dc19b\",\"RequestId\":\"NjYzYjZlM2VfNjg2ODk0MGJfNzIyMV8zMzA4MWE=\"}")
	})

	client.MetaInsight.UpdateFileMetaIndex(context.Background(), nil)

	opt := &UpdateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &File{
			URI:      "cos://test-125000000/1.gif",
			CustomId: "123",
			CustomLabels: &map[string]string{
				"age":   "18",
				"level": "18",
			},
			MediaType:   "video",
			ContentType: "video/gif",
		},
	}
	_, _, err := client.MetaInsight.UpdateFileMetaIndex(context.Background(), opt)
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
		fmt.Fprint(w, "{\"Files\":[{\"COSCRC64\":\"447296710575197191\",\"COSStorageClass\":\"STANDARD\",\"ContentType\":\"image/gif\",\"CreateTime\":\"2024-05-08T20:21:18.766475412+08:00\",\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"DatasetName\":\"adataset\",\"ETag\":\"\\\"c3ad99087956ff0c3d8293ab35747030\\\"\",\"FileModifiedTime\":\"2024-05-06T20:54:07+08:00\",\"Filename\":\"1.gif\",\"MediaType\":\"video\",\"ObjectACL\":\"default\",\"ObjectId\":\"64992b92f79f8ffad132586c4ca26cd4d5dd19783b746e5f6b14dc773f1c0f20\",\"OwnerID\":\"2832742109\",\"Size\":495199,\"URI\":\"cos://test1-1250000000/1.gif\",\"UpdateTime\":\"2024-05-08T20:28:14.884074916+08:00\"}],\"RequestId\":\"NjYzYjcwOWVfMzliNjI5MWVfNmFiZV8zNGM4NDE=\"}")
	})

	opt := &DescribeFileMetaIndexOptions{
		Datasetname: "adataset",
		Uri:         "cos://test-1250000000/1.gif",
	}
	_, _, err := client.MetaInsight.DescribeFileMetaIndex(context.Background(), opt)
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
		fmt.Fprint(w, "{\"RequestId\":\"NjYzYjcxMTZfMmRiNjI5MWVfYWU1XzMxMjk3NQ==\"}")
	})

	client.MetaInsight.DeleteFileMetaIndex(context.Background(), nil)

	opt := &DeleteFileMetaIndexOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000/1.gif",
	}
	_, _, err := client.MetaInsight.DeleteFileMetaIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DeleteFileMetaIndex returned error: %v", err)
	}
}

func TestCIService_DatasetSimpleQuery(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Query\":{\"Operation\":\"eq\",\"SubQueries\":null,\"Field\":\"ContentType\",\"Value\":\"image/gif\"},\"MaxResults\":0,\"NextToken\":\"\",\"Sort\":\"\",\"Order\":\"\",\"Aggregations\":null,\"WithFields\":null}"

	mux.HandleFunc("/datasetquery/simple", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"NextToken\":\"\",\"RequestId\":\"NjYzYjczNjdfNzQ2ODk0MGJfM2NlN18zMWY3YWU=\"}")
	})

	client.MetaInsight.DatasetSimpleQuery(context.Background(), nil)

	opt := &DatasetSimpleQueryOptions{
		DatasetName: "adataset",
		Query: &Query{
			Operation: "eq",
			Field:     "ContentType",
			Value:     "image/gif",
		},
	}
	_, _, err := client.MetaInsight.DatasetSimpleQuery(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DatasetSimpleQuery returned error: %v", err)
	}
}

func TestCIService_DatasetSimpleQueryAggregations(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Query\":null,\"MaxResults\":0,\"NextToken\":\"\",\"Sort\":\"\",\"Order\":\"\",\"Aggregations\":[{\"Operation\":\"group\",\"Field\":\"ContentType\"}],\"WithFields\":null}"

	mux.HandleFunc("/datasetquery/simple", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"Aggregations\":[{\"Field\":\"ContentType\",\"Operation\":\"group\"}],\"NextToken\":\"\",\"RequestId\":\"NjYzYjc0MDFfNTc2ODk0MGJfNjZkNl8zZDA2NjU=\"}")
	})

	client.MetaInsight.DatasetSimpleQuery(context.Background(), nil)

	opt := &DatasetSimpleQueryOptions{
		DatasetName:  "adataset",
		Aggregations: []*Aggregations{},
	}
	opt.Aggregations = append(opt.Aggregations, &Aggregations{
		Field:     "ContentType",
		Operation: "group",
	})
	_, _, err := client.MetaInsight.DatasetSimpleQuery(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DatasetSimpleQuery returned error: %v", err)
	}
}

func TestCIService_CreateDatasetBinding(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\",\"Mode\":0}"

	mux.HandleFunc("/datasetbinding", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"Binding\":{\"CreateTime\":\"2024-05-08T20:47:20.632182296+08:00\",\"DatasetName\":\"adataset\",\"Detail\":\"\",\"State\":\"Running\",\"URI\":\"cos://test1-1250000000\",\"UpdateTime\":\"2024-05-08T20:47:20.632182375+08:00\"},\"RequestId\":\"NjYzYjc0NThfNmQ2ODk0MGJfYmUyXzMyZWE3ZA==\"}")
	})

	opt := &CreateDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
		Mode:        0,
	}
	_, _, err := client.MetaInsight.CreateDatasetBinding(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateDatasetBinding returned error: %v", err)
	}
}

func TestCIService_DescribeDatasetBinding(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/datasetbinding", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
		fmt.Fprint(w, "{\"Binding\":{\"CreateTime\":\"2024-05-08T20:47:20.632182296+08:00\",\"DatasetName\":\"adataset\",\"Detail\":\"\",\"State\":\"Running\",\"URI\":\"cos://test1-1250000000\",\"UpdateTime\":\"2024-05-08T20:47:20.632182375+08:00\"},\"RequestId\":\"NjYzYjc0YTRfNTc2ODk0MGJfNjZkN18zYzcyNjY=\"}")
	})

	opt := &DescribeDatasetBindingOptions{
		Datasetname: "adataset",
		Uri:         "cos://test1-1250000000",
	}
	_, _, err := client.MetaInsight.DescribeDatasetBinding(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeDatasetBinding returned error: %v", err)
	}
}

func TestCIService_DescribeDatasetBindings(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/datasetbindings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Content-Type", "application/json")
		fmt.Fprint(w, "{\"Bindings\":[{\"CreateTime\":\"2024-05-08T20:47:20.632182296+08:00\",\"DatasetName\":\"adataset\",\"Detail\":\"\",\"State\":\"Running\",\"URI\":\"cos://test1-1250000000\",\"UpdateTime\":\"2024-05-08T20:47:20.632182375+08:00\"}],\"NextToken\":\"\",\"RequestId\":\"NjYzYjc1MDBfNmQ2ODk0MGJfYmUyXzMyZWRlNQ==\"}")
	})

	opt := &DescribeDatasetBindingsOptions{
		Datasetname: "adataset",
		// Maxresults: 3,
	}
	_, _, err := client.MetaInsight.DescribeDatasetBindings(context.Background(), opt)
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
		fmt.Fprint(w, "{\"RequestId\":\"NjYzYjc1NDZfNTc2ODk0MGJfNjZkNF8zYzA2MTI=\"}")
	})

	client.MetaInsight.DeleteDatasetBinding(context.Background(), nil)

	opt := &DeleteDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
	}
	_, _, err := client.MetaInsight.DeleteDatasetBinding(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DeleteDatasetBinding returned error: %v", err)
	}
}

func TestCIService_DatasetFaceSearch(t *testing.T) {
	setup()
	defer teardown()
	// wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"facesearch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		// testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"RequestId\":\"NjYzYjc1NDZfNTc2ODk0MGJfNjZkNF8zYzA2MTI=\"}")
	})

	client.MetaInsight.DatasetFaceSearch(context.Background(), nil)

	opt := &DatasetFaceSearchOptions{
		DatasetName:    "test",
		URI:            "cos://examplebucket-1250000000/test.jpg",
		MaxFaceNum:     1,
		Limit:          10,
		MatchThreshold: 10,
	}
	_, _, err := client.MetaInsight.DatasetFaceSearch(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DatasetFaceSearch returned error: %v", err)
	}
}

func TestCIService_SearchImage(t *testing.T) {
	setup()
	defer teardown()
	// wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"imagesearch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		// testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"RequestId\":\"NjYzYjc1NDZfNTc2ODk0MGJfNjZkNF8zYzA2MTI=\"}")
	})

	client.MetaInsight.SearchImage(context.Background(), nil)

	opt := &SearchImageOptions{
		DatasetName:    "ImageSearch001",
		Mode:           "pic",
		URI:            "cos://facesearch-1258726280/huge_base.jpg",
		Limit:          10,
		MatchThreshold: 1,
	}
	_, _, err := client.MetaInsight.SearchImage(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.SearchImage returned error: %v", err)
	}
}
