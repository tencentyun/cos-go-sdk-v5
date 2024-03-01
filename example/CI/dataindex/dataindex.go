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
	u, _ := url.Parse("https://test-125000000.cos.ap-beijing.myqcloud.com")
	cu, _ := url.Parse("https://ci.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
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
	opt := &cos.CreateDataSetOptions{
		DatasetName: "adataset",
		Description: "dataset test",
		TemplateId:  "Official:COSBasicMeta",
	}
	res, _, err := c.CI.CreateDataSet(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDatasets() {
	c := getClient()
	opt := &cos.DescribeDatasetsOptions{
		MaxResults: 100,
	}
	res, _, err := c.CI.DescribeDatasets(context.Background(), opt)
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
	res, _, err := c.CI.UpdateDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DeleteDataset() {
	c := getClient()
	opt := &cos.DeleteDatasetOptions{
		DatasetName: "adataset",
	}
	res, _, err := c.CI.DeleteDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDataset() {
	c := getClient()
	opt := &cos.DescribeDatasetOptions{
		DatasetName: "adataset",
		Statistics:  true,
	}
	res, _, err := c.CI.DescribeDataset(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func CreateFileMetaIndex() {
	c := getClient()
	opt := &cos.CreateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &cos.File{
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
	res, _, err := c.CI.CreateFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func UpdateFileMetaIndex() {
	c := getClient()
	opt := &cos.UpdateFileMetaIndexOptions{
		DatasetName: "adataset",
		File: &cos.File{
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
	res, _, err := c.CI.UpdateFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeFileMetaIndex() {
	c := getClient()
	opt := &cos.DescribeFileMetaIndexOptions{
		DatasetName: "adataset",
		Uri:         "cos://test-125000000/1.gif",
	}
	res, _, err := c.CI.DescribeFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DeleteFileMetaIndex() {
	c := getClient()
	opt := &cos.DeleteFileMetaIndexOptions{
		DatasetName: "adataset",
		Uri:         "cos://test-125000000/1.gif",
	}
	res, _, err := c.CI.DeleteFileMetaIndex(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DatasetSimpleQuery() {
	c := getClient()
	opt := &cos.DatasetSimpleQueryOptions{
		DatasetName: "adataset",
		Query: &cos.Query{
			Operation: "eq",
			Field:"ContentType",
			Value:"image/gif",
		},
	}
	res, _, err := c.CI.DatasetSimpleQuery(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DatasetSimpleQueryAggregations() {
	c := getClient()
	opt := &cos.DatasetSimpleQueryOptions{
		DatasetName: "adataset",
		Aggregations: []*cos.Aggregation{},
	}
	opt.Aggregations = append(opt.Aggregations, &cos.Aggregation{
		Field: "ContentType",
		Operation: "group",
	})
	res, _, err := c.CI.DatasetSimpleQuery(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func CreateDatasetBinding() {
	c := getClient()
	opt := &cos.CreateDatasetBindingOptions{
		DatasetName: "adataset",
		URI: "cos://wwj-bj1-1253960454",
	}
	res, _, err := c.CI.CreateDatasetBinding(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDatasetBinding() {
	c := getClient()
	opt := &cos.DescribeDatasetBindingOptions{
		DatasetName: "adataset",
		URI: "cos://test-125000000",
	}
	res, _, err := c.CI.DescribeDatasetBinding(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DescribeDatasetBindings() {
	c := getClient()
	opt := &cos.DescribeDatasetBindingsOptions{
		DatasetName: "adataset",
		// MaxResults: 3,
	}
	res, _, err := c.CI.DescribeDatasetBindings(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
}

func DeleteDatasetBinding() {
	c := getClient()
	opt := &cos.DeleteDatasetBindingOptions{
		DatasetName: "adataset",
		URI: "cos://test-125000000",
	}
	res, _, err := c.CI.DeleteDatasetBinding(context.Background(), opt)
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
}
