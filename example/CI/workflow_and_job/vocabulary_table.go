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
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
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

func CreateAsrVocabularyTable() {
	// 创建语音识别词表
	c := getClient()
	opt := &cos.CreateAsrVocabularyTableOptions{
		TableName:        "test",
		TableDescription: "test",
		VocabularyWeights: []cos.VocabularyWeight{
			{
				Vocabulary: "test",
				Weight:     10,
			},
		},
	}
	_, _, err := c.CI.CreateAsrVocabularyTable(context.Background(), opt)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", res)
}

func DescribeAsrVocabularyTables() {
	// 查询语音识别词表
	c := getClient()
	opt := &cos.DescribeAsrVocabularyTablesOptions{
		Offset: 0,
		Limit:  10,
	}
	_, _, err := c.CI.DescribeAsrVocabularyTables(context.Background(), opt)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", res)
}

func DeleteAsrVocabularyTable() {
	// 查询语音识别词表
	c := getClient()

	tableId := "c0398427aa1911eebe3c446a2eb5fd98"

	_, err := c.CI.DeleteAsrVocabularyTable(context.Background(), tableId)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", res)
}

func DescribeAsrVocabularyTable() {
	// 查询语音识别词表
	c := getClient()

	tableId := "fc6bd0ce320d11ef8484525400aec391"

	res, _, err := c.CI.DescribeAsrVocabularyTable(context.Background(), tableId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

func UpdateAsrVocabularyTable() {
	// 查询语音识别词表
	c := getClient()

	tableId := "fc6bd0ce320d11ef8484525400aec391"
	opt := &cos.UpdateAsrVocabularyTableOptions{
		TableName:        "test",
		TableDescription: "test",
		TableId:          tableId,
		VocabularyWeights: []cos.VocabularyWeight{
			{
				Vocabulary: "test",
				Weight:     10,
			},
		},
	}
	_, _, err := c.CI.UpdateAsrVocabularyTable(context.Background(), opt)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", res)

	res, _, err := c.CI.DescribeAsrVocabularyTable(context.Background(), tableId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

func main() {
	// CreateAsrVocabularyTable()
	// DescribeAsrVocabularyTables()
	// DeleteAsrVocabularyTable()
	// DescribeAsrVocabularyTable()
	UpdateAsrVocabularyTable()
}
