package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

// logStatus 打印错误信息，支持 COS 错误和 Vector 错误
func logStatus(err error) {
	if err == nil {
		return
	}
	// 优先判断是否是 Vector 错误
	if verr, ok := cos.IsVectorError(err); ok {
		fmt.Printf("ERROR: Vector Error Code: %v\n", verr.Code)
		fmt.Printf("ERROR: Message: %v\n", verr.Message)
		fmt.Printf("ERROR: RequestId: %v\n", verr.RequestID)
		if verr.Response != nil {
			fmt.Printf("ERROR: StatusCode: %d\n", verr.Response.StatusCode)
		}
		for _, f := range verr.FieldList {
			fmt.Printf("ERROR: Field(%s): %s\n", f.Path, f.Message)
		}
	} else if cos.IsNotFoundError(err) {
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: COS Error Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
	} else {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func main() {
	region := "ap-guangzhou"
	secretID := os.Getenv("SECRETID")
	secretKey := os.Getenv("SECRETKEY")

	// ===================== 创建客户端的三种方式 =====================

	// 方式一: 使用公网域名（推荐）
	// 自动生成域名: vectors.<Region>.coslake.com
	vectorURL, err := cos.NewVectorURL(region, true)
	if err != nil {
		fmt.Printf("创建 Vector URL 失败: %v\n", err)
		return
	}

	// 方式二: 使用内网域名（VPC 内访问，可减少流量费用）
	// 自动生成域名: vectors.<Region>.internal.tencentcos.com
	// vectorURL, err := cos.NewVectorInternalURL(region, true)

	// 方式三: 使用自定义 endpoint（完全自定义域名）
	// vectorURL, err := cos.NewVectorEndpointURL("https://my-custom-endpoint.example.com")

	// 创建 COS 客户端
	c := cos.NewClient(&cos.BaseURL{VectorURL: vectorURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	// 可选: 配置重试参数
	// Vector 重试不会切换域名，仅在 5xx 或网络错误时重试
	// c.Conf.RetryOpt.Count = 3         // 最大重试次数（默认3）
	// c.Conf.RetryOpt.Interval = time.Second  // 重试间隔

	bucketName := "examplebucket-1250000000"
	indexName := "my-index"

	// =============== 向量桶管理 ===============

	// 1. 创建向量桶
	fmt.Println("=== 创建向量桶 ===")
	createBucketResult, _, err := c.Vector.CreateVectorBucket(context.Background(), &cos.CreateVectorBucketOptions{
		VectorBucketName: bucketName,
	})
	logStatus(err)
	if err == nil {
		fmt.Printf("向量桶创建成功，QCS: %s\n", createBucketResult.VectorBucketQcs)
	}

	// 2. 查询向量桶信息
	fmt.Println("\n=== 查询向量桶信息 ===")
	getBucketResult, _, err := c.Vector.GetVectorBucket(context.Background(), &cos.GetVectorBucketOptions{
		VectorBucketName: bucketName,
	})
	logStatus(err)
	if err == nil && getBucketResult.VectorBucket != nil {
		fmt.Printf("向量桶名称: %s\n", getBucketResult.VectorBucket.VectorBucketName)
		fmt.Printf("创建时间: %d\n", getBucketResult.VectorBucket.CreationTime)
	}

	// 3. 列出所有向量桶
	fmt.Println("\n=== 列出所有向量桶 ===")
	listBucketsResult, _, err := c.Vector.ListVectorBuckets(context.Background(), &cos.ListVectorBucketsOptions{
		MaxResults: 100,
	})
	logStatus(err)
	if err == nil {
		for _, b := range listBucketsResult.VectorBuckets {
			fmt.Printf("  - %s (创建时间: %d)\n", b.VectorBucketName, b.CreationTime)
		}
	}

	// =============== 向量索引管理 ===============

	// 4. 创建索引
	fmt.Println("\n=== 创建索引 ===")
	createIndexResult, _, err := c.Vector.CreateIndex(context.Background(),
		&cos.CreateIndexOptions{
			VectorBucketName: bucketName,
		},
		&cos.IndexDefinition{
			IndexName:   indexName,
			Dimension:   4,
			Metric:      "COSINE",
			Params: &cos.IndexParams{
				EfConstruction: 200,
				M:              16,
			},
			Description: "示例向量索引",
		},
	)
	logStatus(err)
	if err == nil {
		fmt.Printf("索引创建成功: %s, 状态: %s\n", createIndexResult.Index.IndexName, createIndexResult.Index.Status)
	}

	// 5. 查询索引信息
	fmt.Println("\n=== 查询索引信息 ===")
	getIndexResult, _, err := c.Vector.GetIndex(context.Background(), &cos.GetIndexOptions{
		VectorBucketName: bucketName,
		IndexName:        indexName,
	})
	logStatus(err)
	if err == nil && getIndexResult.Index != nil {
		fmt.Printf("索引名称: %s, 维度: %d, 度量: %s\n",
			getIndexResult.Index.IndexName, getIndexResult.Index.Dimension, getIndexResult.Index.Metric)
	}

	// 6. 列出所有索引
	fmt.Println("\n=== 列出所有索引 ===")
	listIndexesResult, _, err := c.Vector.ListIndexes(context.Background(), &cos.ListIndexesOptions{
		VectorBucketName: bucketName,
	})
	logStatus(err)
	if err == nil {
		for _, idx := range listIndexesResult.Indexes {
			fmt.Printf("  - %s (创建时间: %d)\n", idx.IndexName, idx.CreationTime)
		}
	}

	// =============== 向量数据操作 ===============

	// 7. 插入向量数据
	fmt.Println("\n=== 插入向量数据 ===")
	_, err = c.Vector.PutVectors(context.Background(),
		&cos.PutVectorsOptions{
			VectorBucketName: bucketName,
			IndexName:        indexName,
		},
		[]cos.Vector{
			{
				Key:  "doc-001",
				Data: &cos.VectorData{Float32: []float32{0.1, 0.2, 0.3, 0.4}},
				Metadata: map[string]interface{}{
					"title":    "Go语言入门",
					"category": "programming",
				},
			},
			{
				Key:  "doc-002",
				Data: &cos.VectorData{Float32: []float32{0.5, 0.6, 0.7, 0.8}},
				Metadata: map[string]interface{}{
					"title":    "Python机器学习",
					"category": "AI",
				},
			},
			{
				Key:  "doc-003",
				Data: &cos.VectorData{Float32: []float32{0.9, 0.1, 0.2, 0.3}},
				Metadata: map[string]interface{}{
					"title":    "深度学习基础",
					"category": "AI",
				},
			},
		},
	)
	logStatus(err)
	if err == nil {
		fmt.Println("向量数据插入成功")
	}

	// 8. 获取指定向量
	fmt.Println("\n=== 获取指定向量 ===")
	getVectorsResult, _, err := c.Vector.GetVectors(context.Background(),
		&cos.GetVectorsOptions{
			VectorBucketName: bucketName,
			IndexName:        indexName,
			ReturnData:       true,
			ReturnMetadata:   true,
		},
		[]string{"doc-001", "doc-002"},
	)
	logStatus(err)
	if err == nil {
		for _, v := range getVectorsResult.Vectors {
			fmt.Printf("  Key: %s, Metadata: %v\n", v.Key, v.Metadata)
		}
	}

	// 9. 列出向量
	fmt.Println("\n=== 列出向量 ===")
	listVectorsResult, _, err := c.Vector.ListVectors(context.Background(), &cos.ListVectorsOptions{
		VectorBucketName: bucketName,
		IndexName:        indexName,
		MaxResults:       10,
	})
	logStatus(err)
	if err == nil {
		fmt.Printf("共 %d 个向量\n", len(listVectorsResult.Vectors))
	}

	// 10. 相似度搜索
	fmt.Println("\n=== 相似度搜索 ===")
	queryResult, _, err := c.Vector.QueryVectors(context.Background(),
		&cos.QueryVectorsOptions{
			VectorBucketName: bucketName,
			IndexName:        indexName,
			ReturnData:       true,
			ReturnMetadata:   true,
			ReturnDistance:    true,
		},
		&cos.VectorData{Float32: []float32{0.5, 0.5, 0.6, 0.7}},
		3,
	)
	logStatus(err)
	if err == nil {
		fmt.Printf("找到 %d 个最相似的向量:\n", len(queryResult.Vectors))
		for _, v := range queryResult.Vectors {
			fmt.Printf("  Key: %s, Distance: %f, Metadata: %v\n", v.Key, v.Distance, v.Metadata)
		}
	}

	// 11. 带过滤条件的搜索
	fmt.Println("\n=== 带过滤条件的搜索 ===")
	queryResult2, _, err := c.Vector.QueryVectors(context.Background(),
		&cos.QueryVectorsOptions{
			VectorBucketName: bucketName,
			IndexName:        indexName,
			Filter: map[string]interface{}{
				"category": map[string]interface{}{
					"$eq": "AI",
				},
			},
			ReturnMetadata: true,
			ReturnDistance:  true,
		},
		&cos.VectorData{Float32: []float32{0.5, 0.5, 0.6, 0.7}},
		3,
	)
	logStatus(err)
	if err == nil {
		fmt.Printf("过滤后找到 %d 个向量\n", len(queryResult2.Vectors))
	}

	// 12. 删除向量
	fmt.Println("\n=== 删除向量 ===")
	_, err = c.Vector.DeleteVectors(context.Background(),
		&cos.DeleteVectorsOptions{
			VectorBucketName: bucketName,
			IndexName:        indexName,
		},
		[]string{"doc-001", "doc-002", "doc-003"},
	)
	logStatus(err)
	if err == nil {
		fmt.Println("向量删除成功")
	}

	// =============== 清理资源 ===============

	// 13. 删除索引
	fmt.Println("\n=== 删除索引 ===")
	_, err = c.Vector.DeleteIndex(context.Background(), &cos.DeleteIndexOptions{
		VectorBucketName: bucketName,
		IndexName:        indexName,
	})
	logStatus(err)
	if err == nil {
		fmt.Println("索引删除成功")
	}

	// 14. 删除向量桶
	fmt.Println("\n=== 删除向量桶 ===")
	_, err = c.Vector.DeleteVectorBucket(context.Background(), &cos.DeleteVectorBucketOptions{
		VectorBucketName: bucketName,
	})
	logStatus(err)
	if err == nil {
		fmt.Println("向量桶删除成功")
	}
}
