package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func log_status(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Printf("ERROR: %v\n", err)
		// ERROR
	}
}
func getClient() *cos.Client {
	u, _ := url.Parse("https://test1-1250000000.cos.ap-beijing.myqcloud.com")
	metaInsight, _ := url.Parse("https://1250000000.ci.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, MetaInsightURL: metaInsight}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	return c
}

func CreateDataSet() {
	c := getClient()
	opt := &cos.CreateDatasetOptions{
		DatasetName: "dataset1",
		Description: "dataset test",
		TemplateId:  "Official:COSBasicMeta",
		Version:     "standard",
		Volume:      50,
		SceneType:   "general",
	}
	res, _, err := c.MetaInsight.CreateDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDatasets() {
	c := getClient()
	opt := &cos.DescribeDatasetsOptions{
		Maxresults: 100,
	}
	res, _, err := c.MetaInsight.DescribeDatasets(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func UpdateDataset() {
	c := getClient()
	opt := &cos.UpdateDatasetOptions{
		DatasetName: "adataset",
		Description: "adataset test",
		TemplateId:  "Official:COSBasicMeta",
	}
	res, _, err := c.MetaInsight.UpdateDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DeleteDataset() {
	c := getClient()
	opt := &cos.DeleteDatasetOptions{
		DatasetName: "dataset1",
	}
	res, _, err := c.MetaInsight.DeleteDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDataset() {
	c := getClient()
	opt := &cos.DescribeDatasetOptions{
		Datasetname: "adataset",
		Statistics:  true,
	}
	res, _, err := c.MetaInsight.DescribeDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func CreateFileMetaIndex() {
	c := getClient()
	opt := &cos.CreateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &cos.File{
			URI:      "cos://bj-test-1250000000/5.gif",
			CustomId: "123",
			CustomLabels: &map[string]string{
				"age":   "18",
				"level": "18",
			},
			MediaType:   "image",
			ContentType: "image/gif",
		},
	}
	res, _, err := c.MetaInsight.CreateFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func UpdateFileMetaIndex() {
	c := getClient()
	opt := &cos.UpdateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &cos.File{
			URI:      "cos://test1-1250000000/1.gif",
			CustomId: "123",
			CustomLabels: &map[string]string{
				"age":   "18",
				"level": "18",
			},
			MediaType:   "video",
			ContentType: "video/gif",
		},
	}
	res, _, err := c.MetaInsight.UpdateFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeFileMetaIndex() {
	c := getClient()
	opt := &cos.DescribeFileMetaIndexOptions{
		Datasetname: "adataset",
		Uri:         "cos://test1-1250000000/1.gif",
	}
	res, _, err := c.MetaInsight.DescribeFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DeleteFileMetaIndex() {
	c := getClient()
	opt := &cos.DeleteFileMetaIndexOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000/1.gif",
	}
	res, _, err := c.MetaInsight.DeleteFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DatasetSimpleQuery() {
	c := getClient()
	opt := &cos.DatasetSimpleQueryOptions{
		DatasetName: "adataset",
		Query: &cos.Query{
			Operation: "eq",
			Field:     "ContentType",
			Value:     "image/gif",
		},
	}
	res, _, err := c.MetaInsight.DatasetSimpleQuery(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DatasetSimpleQueryAggregations() {
	c := getClient()
	opt := &cos.DatasetSimpleQueryOptions{
		DatasetName:  "adataset",
		Aggregations: []*cos.Aggregations{},
	}
	opt.Aggregations = append(opt.Aggregations, &cos.Aggregations{
		Field:     "ContentType",
		Operation: "group",
	})
	res, _, err := c.MetaInsight.DatasetSimpleQuery(context.Background(), opt)
	log_status(err)
	for _, v := range res.Aggregations {
		for _, gp := range v.Groups {
			fmt.Printf("%+v\n", gp)
		}
	}
}

func CreateDatasetBinding() {
	c := getClient()
	opt := &cos.CreateDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
		Mode:        0,
	}
	res, _, err := c.MetaInsight.CreateDatasetBinding(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDatasetBinding() {
	c := getClient()
	opt := &cos.DescribeDatasetBindingOptions{
		Datasetname: "adataset",
		Uri:         "cos://test1-1250000000",
	}
	res, _, err := c.MetaInsight.DescribeDatasetBinding(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDatasetBindings() {
	c := getClient()
	opt := &cos.DescribeDatasetBindingsOptions{
		Datasetname: "adataset",
		Maxresults:  3,
	}
	res, _, err := c.MetaInsight.DescribeDatasetBindings(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DeleteDatasetBinding() {
	c := getClient()
	opt := &cos.DeleteDatasetBindingOptions{
		DatasetName: "adataset",
		URI:         "cos://test1-1250000000",
	}
	res, _, err := c.MetaInsight.DeleteDatasetBinding(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DatasetFaceSearch() {
	c := getClient()
	opt := &cos.DatasetFaceSearchOptions{
		DatasetName: "ci-sdk-face-search",
		URI:         "cos://bj-test-1250000000/face.jpeg",
	}
	res, _, err := c.MetaInsight.DatasetFaceSearch(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func SearchImage() {
	c := getClient()
	opt := &cos.SearchImageOptions{
		DatasetName: "ci-sdk-image-search",
		URI:         "cos://bj-test-1250000000/face.jpeg",
		Mode:        "pic",
	}
	res, _, err := c.MetaInsight.SearchImage(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func main() {
	// CreateDataSet()
	// DescribeDatasets()
	// UpdateDataset()
	// DeleteDataset()
	// DescribeDataset()
	// CreateFileMetaIndex()
	// UpdateFileMetaIndex()
	// DescribeFileMetaIndex()
	// DeleteFileMetaIndex()
	// DatasetSimpleQuery()
	// DatasetSimpleQueryAggregations()
	// CreateDatasetBinding()
	// DescribeDatasetBinding()
	// DescribeDatasetBindings()
	// DeleteDatasetBinding()
	// DatasetFaceSearch()
	// SearchImage()
}
