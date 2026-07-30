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
	wantBody := "{\"DatasetName\":\"adataset\",\"Description\":\"dataset test\",\"TemplateId\":\"Official:COSBasicMeta\",\"SceneType\":\"general\"}"

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
		SceneType:   "general",
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
	wantBody := "{\"DatasetName\":\"adataset\",\"File\":{\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"MediaType\":\"image\",\"ContentType\":\"image/gif\",\"URI\":\"cos://test-125000000/12.gif\"}}"

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
	wantBody := "{\"DatasetName\":\"adataset\",\"File\":{\"CustomId\":\"123\",\"CustomLabels\":{\"age\":\"18\",\"level\":\"18\"},\"MediaType\":\"video\",\"ContentType\":\"video/gif\",\"URI\":\"cos://test-125000000/1.gif\"}}"

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
	wantBody := "{\"DatasetName\":\"adataset\",\"Query\":{\"Operation\":\"eq\",\"Field\":\"ContentType\",\"Value\":\"image/gif\"}}"

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
	wantBody := "{\"DatasetName\":\"adataset\",\"Aggregations\":[{\"Operation\":\"group\",\"Field\":\"ContentType\"}]}"

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
	wantBody := "{\"DatasetName\":\"adataset\",\"URI\":\"cos://test1-1250000000\"}"

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

func TestCIService_HybridSearch(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"adataset\",\"Mode\":\"pic\",\"Templates\":\"ImageSearch\",\"SearchURIs\":[\"cos://examplebucket-1250000000/test.jpg\"],\"Limit\":10,\"MatchThreshold\":80,\"Filter\":{\"$and\":[{\"MediaType\":{\"$in\":[\"image\",\"document\"]}},{\"Size\":{\"$gt\":123}}]}}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"hybridsearch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"ImageResult\":[{\"URI\":\"cos://examplebucket-1250000000/result.jpg\",\"Score\":95}],\"DocResult\":[{\"URI\":\"cos://examplebucket-1250000000/doc.pdf\",\"TextPage\":3,\"Score\":88,\"Text\":\"包含一棵大树的图片 {Image_0}\",\"ImageUrls\":{\"Image_0\":\"cos://examplebucket-1250000000/tree.jpg\"}}],\"RequestId\":\"NjYzYjc1NDZfNTc2ODk0MGJfNjZkNF8zYzA2MTI=\"}")
	})

	client.MetaInsight.HybridSearch(context.Background(), nil)

	opt := &HybridSearchOptions{
		DatasetName:    "adataset",
		Mode:           "pic",
		Templates:      "ImageSearch",
		SearchURIs:     []string{"cos://examplebucket-1250000000/test.jpg"},
		Limit:          10,
		MatchThreshold: 80,
		Filter: map[string]interface{}{
			"$and": []map[string]interface{}{
				{
					"MediaType": map[string]interface{}{
						"$in": []string{"image", "document"},
					},
				},
				{
					"Size": map[string]interface{}{
						"$gt": 123,
					},
				},
			},
		},
	}
	_, _, err := client.MetaInsight.HybridSearch(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.HybridSearch returned error: %v", err)
	}
}

func TestCIService_HybridSearch_VideoSearch(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"videosearch\",\"Mode\":\"text\",\"Templates\":\"VideoSearch\",\"SearchText\":\"熊猫在草地上玩耍的片段\",\"Limit\":10}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"hybridsearch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"VideoResult\":[{\"URI\":\"cos://examplebucket-1250000000/panda_play_001.mp4\",\"From\":8.333333,\"To\":14.916667,\"Score\":43},{\"URI\":\"cos://examplebucket-1250000000/panda_play_002.mp4\",\"From\":0,\"To\":8.333333,\"Score\":42}],\"RequestId\":\"NjYwYzEwYjhfNGQ2ODk0MGJfMjcxxxx\"}")
	})

	client.MetaInsight.HybridSearch(context.Background(), nil)

	opt := &HybridSearchOptions{
		DatasetName:    "videosearch",
		Mode:           "text",
		Templates:      "VideoSearch",
		SearchText:     "熊猫在草地上玩耍的片段",
		Limit:          10,
		MatchThreshold: 0,
	}
	res, _, err := client.MetaInsight.HybridSearch(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.HybridSearch VideoSearch returned error: %v", err)
	}
	if len(res.VideoResult) != 2 {
		t.Fatalf("CI.HybridSearch VideoSearch expected 2 results, got %d", len(res.VideoResult))
	}
	if res.VideoResult[0].Score != 43 {
		t.Fatalf("CI.HybridSearch VideoSearch expected score 43, got %d", res.VideoResult[0].Score)
	}
}

func TestCIService_GetAIMediaInfo(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"your-dataset-name-001\",\"URI\":\"cos://examplebucket-1250000000/test.mp4\"}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"getaimediainfo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"RequestId\":\"2FDE2411-DB8D-4A9A-875B-275798F14A5E\",\"MediaInfo\":{\"URI\":\"cos://examplebucket-1250000000/test.mp4\",\"ModifiedTime\":\"2020-12-26T04:11:10Z\",\"FileInfo\":{\"FileBasicInfo\":{\"FileName\":\"example.mp4\",\"FileSize\":30611502,\"FileUrl\":\"cos://examplebucket-1250000000/test.mp4\",\"Region\":\"ap-shanghai\",\"ModifiedTime\":\"2020-12-26T04:11:10Z\",\"FormatName\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Duration\":216.206667,\"Bitrate\":1132.68,\"Width\":960,\"Height\":540},\"AiData\":{\"AiLabelInfo\":[{\"From\":1.4,\"To\":2.5,\"Timestamp\":1.4,\"LabelDetail\":[{\"Confidence\":\"high\",\"LabelInfos\":{\"Type\":\"物体\",\"Category\":\"交通工具\",\"Name\":\"车\"}}]}],\"AsrInfo\":[{\"From\":1.4,\"To\":2.5,\"ClipId\":\"5FE19530C7A422197535FE74F5DB\",\"Content\":\"欢迎观看本期节目\",\"Timestamp\":1.4}],\"OcrInfo\":[{\"From\":1.4,\"To\":2.5,\"ClipId\":\"5FE19530C7A422197535FE74F5DB\",\"Content\":\"城市交通\",\"BoxPosition\":{\"Left\":517,\"Top\":409,\"Width\":128,\"Height\":168},\"Score\":0.86,\"Timestamp\":1.4}],\"AiRoughData\":{\"AiCategory\":\"生活\",\"Description\":\"这是一段关于城市道路交通场景的视频\"},\"FaceInfo\":[{\"Score\":99.04,\"LabelName\":\"张三\",\"Category\":\"custom\",\"OccurrencesInfos\":[{\"From\":61.066353,\"To\":69.06635,\"TrackData\":[{\"Timestamp\":62.03302,\"BoxPosition\":{\"Left\":517,\"Top\":409,\"Width\":128,\"Height\":168}}]}]}]}}}}")
	})

	client.MetaInsight.GetAIMediaInfo(context.Background(), nil)

	opt := &GetAIMediaInfoOptions{
		DatasetName: "your-dataset-name-001",
		URI:         "cos://examplebucket-1250000000/test.mp4",
	}
	_, _, err := client.MetaInsight.GetAIMediaInfo(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.GetAIMediaInfo returned error: %v", err)
	}
}

func TestCIService_MediaFaceSearch(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"your-dataset-name-001\",\"URI\":\"cos://examplebucket-1250000000/face_query.jpg\"}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"mediafacesearch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"RequestId\":\"7CA7D615-CFB1-5437-9A12-2D185C3EE6CB\",\"MediaInfoList\":[{\"FaceId\":\"face_20260206_0001\",\"UriList\":[\"cos://examplebucket-1250000000/photo_001.jpg\",\"cos://examplebucket-1250000000/photo_002.jpg\"]},{\"FaceId\":\"face_20260206_0002\",\"UriList\":[\"cos://examplebucket-1250000000/photo_003.jpg\",\"cos://examplebucket-1250000000/photo_004.jpg\"]}]}")
	})

	client.MetaInsight.MediaFaceSearch(context.Background(), nil)

	opt := &MediaFaceSearchOptions{
		DatasetName: "your-dataset-name-001",
		URI:         "cos://examplebucket-1250000000/face_query.jpg",
	}
	res, _, err := client.MetaInsight.MediaFaceSearch(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.MediaFaceSearch returned error: %v", err)
	}
	if len(res.MediaInfoList) != 2 {
		t.Fatalf("CI.MediaFaceSearch expected 2 face results, got %d", len(res.MediaInfoList))
	}
	if res.MediaInfoList[0].FaceId != "face_20260206_0001" {
		t.Fatalf("CI.MediaFaceSearch expected FaceId face_20260206_0001, got %s", res.MediaInfoList[0].FaceId)
	}
	if len(res.MediaInfoList[0].UriList) != 2 {
		t.Fatalf("CI.MediaFaceSearch expected 2 URIs for first face, got %d", len(res.MediaInfoList[0].UriList))
	}
}

func TestCIService_MediaFaceClipSearch(t *testing.T) {
	setup()
	defer teardown()
	wantBody := "{\"DatasetName\":\"your-dataset-name-001\",\"URI\":\"cos://examplebucket-1250000000/videos/interview_001.mp4\",\"FaceId\":\"face_20260206_0001\"}"

	mux.HandleFunc("/"+"datasetquery"+"/"+"mediafaceclipsearch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, wantBody)
		fmt.Fprint(w, "{\"RequestId\":\"E44FFACD-9E90-555A-A09A-6FD3B7335E39\",\"MediaClipList\":[{\"Score\":99.04,\"LabelName\":\"张三\",\"Category\":\"celebrity\",\"OccurrencesInfos\":[{\"From\":61.066353,\"To\":69.06635,\"TrackData\":[{\"Timestamp\":62.03302,\"BoxPosition\":{\"Left\":517,\"Top\":409,\"Width\":128,\"Height\":168}}]}]}]}")
	})

	client.MetaInsight.MediaFaceClipSearch(context.Background(), nil)

	opt := &MediaFaceClipSearchOptions{
		DatasetName: "your-dataset-name-001",
		URI:         "cos://examplebucket-1250000000/videos/interview_001.mp4",
		FaceId:      "face_20260206_0001",
	}
	res, _, err := client.MetaInsight.MediaFaceClipSearch(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.MediaFaceClipSearch returned error: %v", err)
	}
	if len(res.MediaClipList) != 1 {
		t.Fatalf("CI.MediaFaceClipSearch expected 1 clip result, got %d", len(res.MediaClipList))
	}
	if res.MediaClipList[0].Score != 99.04 {
		t.Fatalf("CI.MediaFaceClipSearch expected Score 99.04, got %f", res.MediaClipList[0].Score)
	}
	if res.MediaClipList[0].Category != "celebrity" {
		t.Fatalf("CI.MediaFaceClipSearch expected Category celebrity, got %s", res.MediaClipList[0].Category)
	}
	if len(res.MediaClipList[0].OccurrencesInfos) != 1 {
		t.Fatalf("CI.MediaFaceClipSearch expected 1 occurrence, got %d", len(res.MediaClipList[0].OccurrencesInfos))
	}
	if res.MediaClipList[0].OccurrencesInfos[0].From != 61.066353 {
		t.Fatalf("CI.MediaFaceClipSearch expected From 61.066353, got %f", res.MediaClipList[0].OccurrencesInfos[0].From)
	}
	if len(res.MediaClipList[0].OccurrencesInfos[0].TrackData) != 1 {
		t.Fatalf("CI.MediaFaceClipSearch expected 1 track data, got %d", len(res.MediaClipList[0].OccurrencesInfos[0].TrackData))
	}
	if res.MediaClipList[0].OccurrencesInfos[0].TrackData[0].BoxPosition.Left != 517 {
		t.Fatalf("CI.MediaFaceClipSearch expected BoxPosition.Left 517, got %d", res.MediaClipList[0].OccurrencesInfos[0].TrackData[0].BoxPosition.Left)
	}
}
