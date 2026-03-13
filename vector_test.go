package cos

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync/atomic"
	"testing"
)

// vectorSetup 为 Vector 测试创建独立的测试服务器和客户端
// 因为 Vector 不走 COS 的 newRequest/doAPI 流程，使用自己的请求链路
func vectorSetup() (mux *http.ServeMux, server *httptest.Server, client *Client) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	u, _ := url.Parse(server.URL)
	client = NewClient(&BaseURL{VectorURL: u}, nil)
	return
}

// ==================== NewVectorURL 测试 ====================

func TestNewVectorURL(t *testing.T) {
	u, err := NewVectorURL("ap-guangzhou", false)
	if err != nil {
		t.Fatalf("NewVectorURL returned error: %v", err)
	}
	want := "http://vectors.ap-guangzhou.coslake.com"
	if u.String() != want {
		t.Errorf("NewVectorURL returned %s, want %s", u.String(), want)
	}

	u, err = NewVectorURL("ap-shanghai", true)
	if err != nil {
		t.Fatalf("NewVectorURL returned error: %v", err)
	}
	want = "https://vectors.ap-shanghai.coslake.com"
	if u.String() != want {
		t.Errorf("NewVectorURL returned %s, want %s", u.String(), want)
	}
}

func TestNewVectorURL_EmptyRegion(t *testing.T) {
	_, err := NewVectorURL("", false)
	if err == nil {
		t.Error("Expected error for empty region")
	}
}

func TestNewVectorInternalURL(t *testing.T) {
	u, err := NewVectorInternalURL("ap-guangzhou", true)
	if err != nil {
		t.Fatalf("NewVectorInternalURL returned error: %v", err)
	}
	want := "https://vectors.ap-guangzhou.internal.tencentcos.com"
	if u.String() != want {
		t.Errorf("NewVectorInternalURL returned %s, want %s", u.String(), want)
	}

	u, err = NewVectorInternalURL("ap-beijing", false)
	if err != nil {
		t.Fatalf("NewVectorInternalURL returned error: %v", err)
	}
	want = "http://vectors.ap-beijing.internal.tencentcos.com"
	if u.String() != want {
		t.Errorf("NewVectorInternalURL returned %s, want %s", u.String(), want)
	}
}

func TestNewVectorInternalURL_EmptyRegion(t *testing.T) {
	_, err := NewVectorInternalURL("", true)
	if err == nil {
		t.Error("Expected error for empty region")
	}
}

func TestNewVectorEndpointURL(t *testing.T) {
	// 带 scheme 的 endpoint
	u, err := NewVectorEndpointURL("https://my-custom-vector.example.com")
	if err != nil {
		t.Fatalf("NewVectorEndpointURL returned error: %v", err)
	}
	want := "https://my-custom-vector.example.com"
	if u.String() != want {
		t.Errorf("NewVectorEndpointURL returned %s, want %s", u.String(), want)
	}

	// 不带 scheme 的 endpoint，自动加 https
	u, err = NewVectorEndpointURL("vectors.ap-guangzhou.coslake.com")
	if err != nil {
		t.Fatalf("NewVectorEndpointURL returned error: %v", err)
	}
	want = "https://vectors.ap-guangzhou.coslake.com"
	if u.String() != want {
		t.Errorf("NewVectorEndpointURL returned %s, want %s", u.String(), want)
	}

	// http endpoint
	u, err = NewVectorEndpointURL("http://vectors.ap-guangzhou.coslake.com")
	if err != nil {
		t.Fatalf("NewVectorEndpointURL returned error: %v", err)
	}
	want = "http://vectors.ap-guangzhou.coslake.com"
	if u.String() != want {
		t.Errorf("NewVectorEndpointURL returned %s, want %s", u.String(), want)
	}
}

func TestNewVectorEndpointURL_Empty(t *testing.T) {
	_, err := NewVectorEndpointURL("")
	if err == nil {
		t.Error("Expected error for empty endpoint")
	}
}

// ==================== 向量桶管理测试 ====================

func TestVectorService_CreateVectorBucket(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/CreateVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := ioutil.ReadAll(r.Body)
		var req CreateVectorBucketOptions
		json.Unmarshal(body, &req)
		if req.VectorBucketName != "examplebucket-1250000000" {
			t.Errorf("Expected bucket name examplebucket-1250000000, got %s", req.VectorBucketName)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectorBucketQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000"
		}`)
	})

	opt := &CreateVectorBucketOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	res, _, err := client.Vector.CreateVectorBucket(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.CreateVectorBucket returned error: %v", err)
	}

	want := &CreateVectorBucketResult{
		VectorBucketQcs: "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000",
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Vector.CreateVectorBucket returned %+v, want %+v", res, want)
	}
}

func TestVectorService_CreateVectorBucketWithEncryption(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/CreateVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req CreateVectorBucketOptions
		json.Unmarshal(body, &req)
		if req.EncryptionConfiguration == nil || req.EncryptionConfiguration.SseType != "AES256" {
			t.Errorf("Expected AES256 encryption")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectorBucketQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000"
		}`)
	})

	opt := &CreateVectorBucketOptions{
		VectorBucketName: "examplebucket-1250000000",
		EncryptionConfiguration: &VectorEncryptionConfig{
			SseType: "AES256",
		},
	}
	res, _, err := client.Vector.CreateVectorBucket(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.CreateVectorBucket returned error: %v", err)
	}
	if res.VectorBucketQcs == "" {
		t.Error("Expected non-empty VectorBucketQcs")
	}
}

func TestVectorService_GetVectorBucket(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/GetVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectorBucket": {
				"creationTime": 1735445700,
				"encryptionConfiguration": {
					"sseType": "AES256"
				},
				"vectorBucketQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000",
				"vectorBucketName": "examplebucket-1250000000"
			}
		}`)
	})

	opt := &GetVectorBucketOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	res, _, err := client.Vector.GetVectorBucket(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.GetVectorBucket returned error: %v", err)
	}
	if res.VectorBucket == nil {
		t.Fatal("Expected non-nil VectorBucket")
	}
	if res.VectorBucket.VectorBucketName != "examplebucket-1250000000" {
		t.Errorf("Expected bucket name examplebucket-1250000000, got %s", res.VectorBucket.VectorBucketName)
	}
	if res.VectorBucket.CreationTime != 1735445700 {
		t.Errorf("Expected creationTime 1735445700, got %d", res.VectorBucket.CreationTime)
	}
	if res.VectorBucket.EncryptionConfiguration == nil || res.VectorBucket.EncryptionConfiguration.SseType != "AES256" {
		t.Error("Expected AES256 encryption")
	}
}

func TestVectorService_DeleteVectorBucket(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/DeleteVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req DeleteVectorBucketOptions
		json.Unmarshal(body, &req)
		if req.VectorBucketName != "examplebucket-1250000000" {
			t.Errorf("Expected bucket name examplebucket-1250000000, got %s", req.VectorBucketName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	opt := &DeleteVectorBucketOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	_, err := client.Vector.DeleteVectorBucket(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.DeleteVectorBucket returned error: %v", err)
	}
}

func TestVectorService_ListVectorBuckets(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/ListVectorBuckets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectorBuckets": [
				{
					"creationTime": 1735445700,
					"vectorBucketQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/bucket1-1250000000",
					"vectorBucketName": "bucket1-1250000000"
				},
				{
					"creationTime": 1735449900,
					"vectorBucketQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/bucket2-1250000000",
					"vectorBucketName": "bucket2-1250000000"
				}
			],
			"nextToken": "token123"
		}`)
	})

	opt := &ListVectorBucketsOptions{
		MaxResults: 10,
	}
	res, _, err := client.Vector.ListVectorBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.ListVectorBuckets returned error: %v", err)
	}
	if len(res.VectorBuckets) != 2 {
		t.Errorf("Expected 2 buckets, got %d", len(res.VectorBuckets))
	}
	if res.NextToken != "token123" {
		t.Errorf("Expected nextToken token123, got %s", res.NextToken)
	}
	if res.VectorBuckets[0].VectorBucketName != "bucket1-1250000000" {
		t.Errorf("Expected first bucket name bucket1-1250000000, got %s", res.VectorBuckets[0].VectorBucketName)
	}
}

// ==================== 向量桶策略管理测试 ====================

func TestVectorService_PutVectorBucketPolicy(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/PutVectorBucketPolicy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	policy := map[string]interface{}{
		"Statement": []interface{}{
			map[string]interface{}{
				"Effect":   "Allow",
				"Action":   "cos:GetVectors",
				"Resource": "*",
			},
		},
	}
	opt := &PutVectorBucketPolicyOptions{
		VectorBucketName: "examplebucket-1250000000",
		Policy:           policy,
	}
	_, err := client.Vector.PutVectorBucketPolicy(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.PutVectorBucketPolicy returned error: %v", err)
	}
}

func TestVectorService_GetVectorBucketPolicy(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/GetVectorBucketPolicy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"policy": {
				"Statement": [
					{
						"Effect": "Allow",
						"Action": "cos:GetVectors",
						"Resource": "*"
					}
				]
			}
		}`)
	})

	opt := &GetVectorBucketPolicyOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	res, _, err := client.Vector.GetVectorBucketPolicy(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.GetVectorBucketPolicy returned error: %v", err)
	}
	if res.Policy == nil {
		t.Error("Expected non-nil policy")
	}
}

func TestVectorService_DeleteVectorBucketPolicy(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/DeleteVectorBucketPolicy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	opt := &DeleteVectorBucketPolicyOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	_, err := client.Vector.DeleteVectorBucketPolicy(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.DeleteVectorBucketPolicy returned error: %v", err)
	}
}

// ==================== 向量索引管理测试 ====================

func TestVectorService_CreateIndex(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/CreateIndex", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		if req["indexName"] != "test-index" {
			t.Errorf("Expected index name test-index, got %v", req["indexName"])
		}
		if req["dimension"].(float64) != 128 {
			t.Errorf("Expected dimension 128, got %v", req["dimension"])
		}
		if req["vectorBucketName"] != "examplebucket-1250000000" {
			t.Errorf("Expected vectorBucketName examplebucket-1250000000, got %v", req["vectorBucketName"])
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"index": {
				"indexQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000:index/test-index",
				"indexName": "test-index",
				"vectorBucketName": "examplebucket-1250000000",
				"creationTime": 1735445700,
				"dimension": 128,
				"metric": "COSINE",
				"params": {
					"efConstruction": 200,
					"m": 16
				},
				"status": "Active"
			}
		}`)
	})

	opt := &CreateIndexOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	indexDef := &IndexDefinition{
		IndexName: "test-index",
		Dimension: 128,
		Metric:    "COSINE",
		Params: &IndexParams{
			EfConstruction: 200,
			M:              16,
		},
	}
	res, _, err := client.Vector.CreateIndex(context.Background(), opt, indexDef)
	if err != nil {
		t.Fatalf("Vector.CreateIndex returned error: %v", err)
	}
	if res.Index == nil {
		t.Fatal("Expected non-nil Index")
	}
	if res.Index.IndexName != "test-index" {
		t.Errorf("Expected index name test-index, got %s", res.Index.IndexName)
	}
	if res.Index.Dimension != 128 {
		t.Errorf("Expected dimension 128, got %d", res.Index.Dimension)
	}
	if res.Index.Status != "Active" {
		t.Errorf("Expected status Active, got %s", res.Index.Status)
	}
	if res.Index.Params == nil || res.Index.Params.EfConstruction != 200 {
		t.Error("Expected params with efConstruction=200")
	}
}

func TestVectorService_GetIndex(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/GetIndex", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"index": {
				"indexQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000:index/test-index",
				"indexName": "test-index",
				"vectorBucketName": "examplebucket-1250000000",
				"creationTime": 1735445700,
				"dimension": 128,
				"metric": "COSINE",
				"status": "Active"
			}
		}`)
	})

	opt := &GetIndexOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
	}
	res, _, err := client.Vector.GetIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.GetIndex returned error: %v", err)
	}
	if res.Index == nil {
		t.Fatal("Expected non-nil Index")
	}
	if res.Index.Metric != "COSINE" {
		t.Errorf("Expected metric COSINE, got %s", res.Index.Metric)
	}
}

func TestVectorService_ListIndexes(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/ListIndexes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"indexes": [
				{
					"creationTime": 1735449900,
					"indexName": "index1",
					"indexQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000:index/index1",
					"vectorBucketName": "examplebucket-1250000000"
				},
				{
					"creationTime": 1731657600,
					"indexName": "index2",
					"indexQcs": "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000:index/index2",
					"vectorBucketName": "examplebucket-1250000000"
				}
			],
			"nextToken": "nextpage"
		}`)
	})

	opt := &ListIndexesOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	res, _, err := client.Vector.ListIndexes(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.ListIndexes returned error: %v", err)
	}
	if len(res.Indexes) != 2 {
		t.Errorf("Expected 2 indexes, got %d", len(res.Indexes))
	}
	if res.NextToken != "nextpage" {
		t.Errorf("Expected nextToken nextpage, got %s", res.NextToken)
	}
}

func TestVectorService_DeleteIndex(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/DeleteIndex", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req DeleteIndexOptions
		json.Unmarshal(body, &req)
		if req.IndexName != "test-index" {
			t.Errorf("Expected index name test-index, got %s", req.IndexName)
		}
		w.WriteHeader(http.StatusOK)
	})

	opt := &DeleteIndexOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
	}
	_, err := client.Vector.DeleteIndex(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.DeleteIndex returned error: %v", err)
	}
}

// ==================== 向量数据操作测试 ====================

func TestVectorService_PutVectors(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/PutVectors", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		vectors := req["vectors"].([]interface{})
		if len(vectors) != 2 {
			t.Errorf("Expected 2 vectors, got %d", len(vectors))
		}
		first := vectors[0].(map[string]interface{})
		if first["key"] != "doc-001" {
			t.Errorf("Expected key doc-001, got %v", first["key"])
		}

		w.WriteHeader(http.StatusOK)
	})

	opt := &PutVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
	}
	vectors := []Vector{
		{
			Key:  "doc-001",
			Data: &VectorData{Float32: []float32{0.1, 0.2, 0.3, 0.4}},
			Metadata: map[string]interface{}{
				"title":    "文档标题",
				"category": "AI",
			},
		},
		{
			Key:  "doc-002",
			Data: &VectorData{Float32: []float32{0.5, 0.6, 0.7, 0.8}},
		},
	}
	_, err := client.Vector.PutVectors(context.Background(), opt, vectors)
	if err != nil {
		t.Fatalf("Vector.PutVectors returned error: %v", err)
	}
}

func TestVectorService_GetVectors(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/GetVectors", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		keys := req["keys"].([]interface{})
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectors": [
				{
					"key": "doc-001",
					"data": {"float32": [1.0, 2.0]},
					"metadata": {"color": "red", "count": 10}
				},
				{
					"key": "doc-002",
					"data": {"float32": [3.0, 4.0]},
					"metadata": {"color": "blue", "count": 20}
				}
			]
		}`)
	})

	opt := &GetVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
		ReturnData:       true,
		ReturnMetadata:   true,
	}
	keys := []string{"doc-001", "doc-002"}
	res, _, err := client.Vector.GetVectors(context.Background(), opt, keys)
	if err != nil {
		t.Fatalf("Vector.GetVectors returned error: %v", err)
	}
	if len(res.Vectors) != 2 {
		t.Errorf("Expected 2 vectors, got %d", len(res.Vectors))
	}
	if res.Vectors[0].Key != "doc-001" {
		t.Errorf("Expected key doc-001, got %s", res.Vectors[0].Key)
	}
	if res.Vectors[0].Data == nil || len(res.Vectors[0].Data.Float32) != 2 {
		t.Error("Expected vector data with 2 floats")
	}
}

func TestVectorService_ListVectors(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/ListVectors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectors": [
				{"key": "doc-001"},
				{"key": "doc-002"},
				{"key": "doc-003"}
			],
			"nextToken": "abc"
		}`)
	})

	opt := &ListVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
		MaxResults:       10,
	}
	res, _, err := client.Vector.ListVectors(context.Background(), opt)
	if err != nil {
		t.Fatalf("Vector.ListVectors returned error: %v", err)
	}
	if len(res.Vectors) != 3 {
		t.Errorf("Expected 3 vectors, got %d", len(res.Vectors))
	}
	if res.NextToken != "abc" {
		t.Errorf("Expected nextToken abc, got %s", res.NextToken)
	}
}

func TestVectorService_DeleteVectors(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/DeleteVectors", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		keys := req["keys"].([]interface{})
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}

		w.WriteHeader(http.StatusOK)
	})

	opt := &DeleteVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
	}
	keys := []string{"doc-001", "doc-002"}
	_, err := client.Vector.DeleteVectors(context.Background(), opt, keys)
	if err != nil {
		t.Fatalf("Vector.DeleteVectors returned error: %v", err)
	}
}

func TestVectorService_QueryVectors(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/QueryVectors", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		if req["topK"].(float64) != 5 {
			t.Errorf("Expected topK 5, got %v", req["topK"])
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectors": [
				{
					"key": "doc-001",
					"data": {"float32": [1.0, 2.0]},
					"metadata": {"color": "red"},
					"distance": 0.0
				},
				{
					"key": "doc-002",
					"data": {"float32": [3.0, 4.0]},
					"metadata": {"color": "blue"},
					"distance": 8.0
				}
			]
		}`)
	})

	opt := &QueryVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
		ReturnData:       true,
		ReturnMetadata:   true,
		ReturnDistance:    true,
	}
	queryVector := &VectorData{Float32: []float32{1.0, 2.0}}
	res, _, err := client.Vector.QueryVectors(context.Background(), opt, queryVector, 5)
	if err != nil {
		t.Fatalf("Vector.QueryVectors returned error: %v", err)
	}
	if len(res.Vectors) != 2 {
		t.Errorf("Expected 2 vectors, got %d", len(res.Vectors))
	}
	if res.Vectors[0].Key != "doc-001" {
		t.Errorf("Expected key doc-001, got %s", res.Vectors[0].Key)
	}
	if res.Vectors[0].Distance != 0.0 {
		t.Errorf("Expected distance 0.0, got %f", res.Vectors[0].Distance)
	}
	if res.Vectors[1].Distance != 8.0 {
		t.Errorf("Expected distance 8.0, got %f", res.Vectors[1].Distance)
	}
}

func TestVectorService_QueryVectorsWithFilter(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/QueryVectors", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		if req["filter"] == nil {
			t.Error("Expected non-nil filter")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"vectors": [
				{
					"key": "doc-001",
					"distance": 0.5
				}
			]
		}`)
	})

	opt := &QueryVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
		Filter: map[string]interface{}{
			"category": map[string]interface{}{
				"$eq": "AI",
			},
		},
		ReturnDistance: true,
	}
	queryVector := &VectorData{Float32: []float32{1.0, 2.0}}
	res, _, err := client.Vector.QueryVectors(context.Background(), opt, queryVector, 5)
	if err != nil {
		t.Fatalf("Vector.QueryVectors returned error: %v", err)
	}
	if len(res.Vectors) != 1 {
		t.Errorf("Expected 1 vector, got %d", len(res.Vectors))
	}
}

// ==================== 错误处理测试 ====================

func TestVectorService_Error_ValidationException(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/CreateVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "ValidationException")
		w.Header().Set("X-Cos-Request-Id", "NjM3ZmI5YTlfOTBm")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"message": "VectorBucketName is invalid",
			"fieldList": [
				{
					"message": "VectorBucketName should match pattern",
					"path": "/vectorBucketName"
				}
			]
		}`)
	})

	opt := &CreateVectorBucketOptions{
		VectorBucketName: "invalid",
	}
	_, resp, err := client.Vector.CreateVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error for 400 response")
	}
	if resp == nil {
		t.Fatal("Expected non-nil response even on error")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T: %v", err, err)
	}
	if verr.Code != "ValidationException" {
		t.Errorf("Expected code ValidationException, got %s", verr.Code)
	}
	if verr.Message != "VectorBucketName is invalid" {
		t.Errorf("Expected message 'VectorBucketName is invalid', got %s", verr.Message)
	}
	if verr.RequestID != "NjM3ZmI5YTlfOTBm" {
		t.Errorf("Expected requestId NjM3ZmI5YTlfOTBm, got %s", verr.RequestID)
	}
	if len(verr.FieldList) != 1 {
		t.Fatalf("Expected 1 field error, got %d", len(verr.FieldList))
	}
	if verr.FieldList[0].Path != "/vectorBucketName" {
		t.Errorf("Expected path /vectorBucketName, got %s", verr.FieldList[0].Path)
	}

	// 检查 Error() 输出
	errStr := verr.Error()
	if errStr == "" {
		t.Error("Expected non-empty error string")
	}
}

func TestVectorService_Error_NotFoundException(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/GetVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "NotFoundException")
		w.Header().Set("X-Cos-Request-Id", "req-123456")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{
			"message": "VectorBucket not found"
		}`)
	})

	opt := &GetVectorBucketOptions{
		VectorBucketName: "nonexistent-1250000000",
	}
	_, _, err := client.Vector.GetVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error for 404 response")
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T", err)
	}
	if verr.Code != "NotFoundException" {
		t.Errorf("Expected code NotFoundException, got %s", verr.Code)
	}
	if verr.Message != "VectorBucket not found" {
		t.Errorf("Expected message 'VectorBucket not found', got %s", verr.Message)
	}
}

func TestVectorService_Error_ConflictException(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/CreateVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "ConflictException")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{
			"message": "VectorBucket already exists"
		}`)
	})

	opt := &CreateVectorBucketOptions{
		VectorBucketName: "existing-1250000000",
	}
	_, _, err := client.Vector.CreateVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error for 409 response")
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T", err)
	}
	if verr.Code != "ConflictException" {
		t.Errorf("Expected code ConflictException, got %s", verr.Code)
	}
}

func TestVectorService_Error_AccessDeniedException(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/GetVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "AccessDeniedException")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"message": "Access denied"
		}`)
	})

	opt := &GetVectorBucketOptions{
		VectorBucketName: "private-1250000000",
	}
	_, _, err := client.Vector.GetVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error for 403 response")
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T", err)
	}
	if verr.Code != "AccessDeniedException" {
		t.Errorf("Expected code AccessDeniedException, got %s", verr.Code)
	}
}

func TestVectorService_Error_TooManyRequestsException(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	mux.HandleFunc("/PutVectors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "TooManyRequestsException")
		w.WriteHeader(429)
		fmt.Fprint(w, `{
			"message": "Too many requests"
		}`)
	})

	opt := &PutVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
	}
	vectors := []Vector{{Key: "doc-001"}}
	_, err := client.Vector.PutVectors(context.Background(), opt, vectors)
	if err == nil {
		t.Fatal("Expected error for 429 response")
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T", err)
	}
	if verr.Code != "TooManyRequestsException" {
		t.Errorf("Expected code TooManyRequestsException, got %s", verr.Code)
	}
}

func TestVectorService_Error_InternalServerException(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	// 设置为不重试，方便测试
	client.Conf.RetryOpt.Count = 1

	mux.HandleFunc("/GetVectors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "InternalServerException")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{
			"message": "Internal server error"
		}`)
	})

	opt := &GetVectorsOptions{
		VectorBucketName: "examplebucket-1250000000",
		IndexName:        "test-index",
	}
	keys := []string{"doc-001"}
	_, _, err := client.Vector.GetVectors(context.Background(), opt, keys)
	if err == nil {
		t.Fatal("Expected error for 500 response")
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T: %v", err, err)
	}
	if verr.Code != "InternalServerException" {
		t.Errorf("Expected code InternalServerException, got %s", verr.Code)
	}
}

func TestVectorService_Error_EmptyBody(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	client.Conf.RetryOpt.Count = 1

	mux.HandleFunc("/DeleteVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Cos-Error-Code", "InternalServerException")
		w.WriteHeader(http.StatusInternalServerError)
	})

	opt := &DeleteVectorBucketOptions{
		VectorBucketName: "examplebucket-1250000000",
	}
	_, err := client.Vector.DeleteVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error for 500 response")
	}

	verr, ok := IsVectorError(err)
	if !ok {
		t.Fatalf("Expected VectorErrorResponse, got %T", err)
	}
	if verr.Code != "InternalServerException" {
		t.Errorf("Expected code InternalServerException, got %s", verr.Code)
	}
	// Message 应该为空因为没有 body
	if verr.Message != "" {
		t.Errorf("Expected empty message, got %s", verr.Message)
	}
}

// ==================== IsVectorError 测试 ====================

func TestIsVectorError_Nil(t *testing.T) {
	verr, ok := IsVectorError(nil)
	if ok || verr != nil {
		t.Error("Expected false for nil error")
	}
}

func TestIsVectorError_NonVector(t *testing.T) {
	err := fmt.Errorf("not a vector error")
	verr, ok := IsVectorError(err)
	if ok || verr != nil {
		t.Error("Expected false for non-VectorErrorResponse")
	}
}

func TestIsVectorError_CosError(t *testing.T) {
	// COS ErrorResponse 不应该被识别为 VectorError
	err := &ErrorResponse{Code: "NoSuchKey", Message: "Key not found"}
	verr, ok := IsVectorError(err)
	if ok || verr != nil {
		t.Error("Expected false for COS ErrorResponse")
	}
}

// ==================== 重试逻辑测试 ====================

func TestVectorService_Retry_On500(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	client.Conf.RetryOpt.Count = 3

	var callCount int32

	mux.HandleFunc("/ListVectorBuckets", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count < 3 {
			// 前两次返回 500
			w.Header().Set("X-Cos-Error-Code", "InternalServerException")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message": "server error"}`)
		} else {
			// 第三次成功
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"vectorBuckets": []}`)
		}
	})

	res, _, err := client.Vector.ListVectorBuckets(context.Background(), &ListVectorBucketsOptions{})
	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}
	if res == nil {
		t.Fatal("Expected non-nil result")
	}

	finalCount := atomic.LoadInt32(&callCount)
	if finalCount != 3 {
		t.Errorf("Expected 3 calls (2 retries + 1 success), got %d", finalCount)
	}
}

func TestVectorService_NoRetry_On400(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	client.Conf.RetryOpt.Count = 3

	var callCount int32

	mux.HandleFunc("/CreateVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cos-Error-Code", "ValidationException")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message": "invalid request"}`)
	})

	opt := &CreateVectorBucketOptions{
		VectorBucketName: "invalid",
	}
	_, _, err := client.Vector.CreateVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error")
	}

	finalCount := atomic.LoadInt32(&callCount)
	if finalCount != 1 {
		t.Errorf("Expected 1 call (no retry on 400), got %d", finalCount)
	}
}

func TestVectorService_NoRetry_On404(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	client.Conf.RetryOpt.Count = 3

	var callCount int32

	mux.HandleFunc("/GetVectorBucket", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.Header().Set("X-Cos-Error-Code", "NotFoundException")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message": "not found"}`)
	})

	opt := &GetVectorBucketOptions{
		VectorBucketName: "nonexistent-1250000000",
	}
	_, _, err := client.Vector.GetVectorBucket(context.Background(), opt)
	if err == nil {
		t.Fatal("Expected error")
	}

	finalCount := atomic.LoadInt32(&callCount)
	if finalCount != 1 {
		t.Errorf("Expected 1 call (no retry on 404), got %d", finalCount)
	}
}

func TestVectorService_Retry_RetryHeader(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	client.Conf.RetryOpt.Count = 2

	var retryHeaders []string

	mux.HandleFunc("/ListVectorBuckets", func(w http.ResponseWriter, r *http.Request) {
		retryHeaders = append(retryHeaders, r.Header.Get("X-Cos-Sdk-Retry"))
		w.Header().Set("X-Cos-Error-Code", "ServiceUnavailableException")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"message": "service unavailable"}`)
	})

	client.Vector.ListVectorBuckets(context.Background(), &ListVectorBucketsOptions{})

	if len(retryHeaders) != 2 {
		t.Fatalf("Expected 2 calls, got %d", len(retryHeaders))
	}
	// 第一次请求不带 retry header
	if retryHeaders[0] != "" {
		t.Errorf("First request should not have retry header, got %s", retryHeaders[0])
	}
	// 重试请求应该带 retry header
	if retryHeaders[1] != "true" {
		t.Errorf("Retry request should have X-Cos-Sdk-Retry=true, got %s", retryHeaders[1])
	}
}

func TestVectorService_Retry_NoDomainSwitch(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()

	client.Conf.RetryOpt.Count = 3
	client.Conf.RetryOpt.AutoSwitchHost = true // 即使开启了域名切换，Vector 也不应切换

	var requestHosts []string

	mux.HandleFunc("/ListVectorBuckets", func(w http.ResponseWriter, r *http.Request) {
		requestHosts = append(requestHosts, r.Host)
		w.Header().Set("X-Cos-Error-Code", "InternalServerException")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "server error"}`)
	})

	client.Vector.ListVectorBuckets(context.Background(), &ListVectorBucketsOptions{})

	// 所有请求应该使用同一个 host，不切换
	for i := 1; i < len(requestHosts); i++ {
		if requestHosts[i] != requestHosts[0] {
			t.Errorf("Expected same host on retry, got %s vs %s", requestHosts[0], requestHosts[i])
		}
	}
}

// ==================== nil 参数测试 ====================

func TestVectorService_NilParams(t *testing.T) {
	mux, server, client := vectorSetup()
	defer server.Close()
	_ = mux

	_, _, err := client.Vector.CreateVectorBucket(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil CreateVectorBucket options")
	}

	_, _, err = client.Vector.GetVectorBucket(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil GetVectorBucket options")
	}

	_, err = client.Vector.DeleteVectorBucket(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil DeleteVectorBucket options")
	}

	_, err = client.Vector.PutVectorBucketPolicy(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil PutVectorBucketPolicy options")
	}

	_, _, err = client.Vector.GetVectorBucketPolicy(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil GetVectorBucketPolicy options")
	}

	_, err = client.Vector.DeleteVectorBucketPolicy(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil DeleteVectorBucketPolicy options")
	}

	_, _, err = client.Vector.CreateIndex(context.Background(), nil, nil)
	if err == nil {
		t.Error("Expected error for nil CreateIndex options")
	}

	// CreateIndex: opt 非 nil 但 indexDef 为 nil
	_, _, err = client.Vector.CreateIndex(context.Background(), &CreateIndexOptions{}, nil)
	if err == nil {
		t.Error("Expected error for nil CreateIndex indexDef")
	}

	_, _, err = client.Vector.GetIndex(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil GetIndex options")
	}

	_, _, err = client.Vector.ListIndexes(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil ListIndexes options")
	}

	_, err = client.Vector.DeleteIndex(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil DeleteIndex options")
	}

	_, err = client.Vector.PutVectors(context.Background(), nil, nil)
	if err == nil {
		t.Error("Expected error for nil PutVectors options")
	}

	// PutVectors: opt 非 nil 但 vectors 为空
	_, err = client.Vector.PutVectors(context.Background(), &PutVectorsOptions{}, nil)
	if err == nil {
		t.Error("Expected error for empty PutVectors vectors")
	}

	_, _, err = client.Vector.GetVectors(context.Background(), nil, nil)
	if err == nil {
		t.Error("Expected error for nil GetVectors options")
	}

	// GetVectors: opt 非 nil 但 keys 为空
	_, _, err = client.Vector.GetVectors(context.Background(), &GetVectorsOptions{}, nil)
	if err == nil {
		t.Error("Expected error for empty GetVectors keys")
	}

	_, _, err = client.Vector.ListVectors(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil ListVectors options")
	}

	_, err = client.Vector.DeleteVectors(context.Background(), nil, nil)
	if err == nil {
		t.Error("Expected error for nil DeleteVectors options")
	}

	// DeleteVectors: opt 非 nil 但 keys 为空
	_, err = client.Vector.DeleteVectors(context.Background(), &DeleteVectorsOptions{}, nil)
	if err == nil {
		t.Error("Expected error for empty DeleteVectors keys")
	}

	_, _, err = client.Vector.QueryVectors(context.Background(), nil, nil, 0)
	if err == nil {
		t.Error("Expected error for nil QueryVectors options")
	}

	// QueryVectors: opt 非 nil 但 queryVector 为 nil
	_, _, err = client.Vector.QueryVectors(context.Background(), &QueryVectorsOptions{}, nil, 5)
	if err == nil {
		t.Error("Expected error for nil QueryVectors queryVector")
	}

	// QueryVectors: topK <= 0
	_, _, err = client.Vector.QueryVectors(context.Background(), &QueryVectorsOptions{}, &VectorData{Float32: []float32{1.0}}, 0)
	if err == nil {
		t.Error("Expected error for topK <= 0")
	}
}

// ==================== VectorURL 未设置测试 ====================

func TestVectorService_NilVectorURL(t *testing.T) {
	client := NewClient(&BaseURL{}, nil)

	_, _, err := client.Vector.CreateVectorBucket(context.Background(), &CreateVectorBucketOptions{
		VectorBucketName: "test-1250000000",
	})
	if err == nil {
		t.Fatal("Expected error when VectorURL is nil")
	}
}

// ==================== VectorErrorResponse.Error() 格式测试 ====================

func TestVectorErrorResponse_Error(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://vectors.ap-guangzhou.coslake.com/CreateVectorBucket", nil)
	resp := &http.Response{
		StatusCode: 400,
		Request:    req,
		Header:     http.Header{},
	}

	verr := &VectorErrorResponse{
		Response:  resp,
		Code:      "ValidationException",
		Message:   "VectorBucketName is invalid",
		RequestID: "req-123",
		FieldList: []VectorValidateField{
			{Message: "name should match pattern", Path: "/vectorBucketName"},
		},
	}

	errStr := verr.Error()
	if errStr == "" {
		t.Error("Expected non-empty error string")
	}

	// 检查格式包含关键信息
	for _, expected := range []string{"POST", "CreateVectorBucket", "400", "ValidationException", "VectorBucketName is invalid", "req-123", "/vectorBucketName"} {
		found := false
		if len(errStr) > 0 {
			for i := 0; i <= len(errStr)-len(expected); i++ {
				if errStr[i:i+len(expected)] == expected {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("Error string missing expected content '%s': %s", expected, errStr)
		}
	}
}

// ==================== checkVectorResponse 测试 ====================

func TestCheckVectorResponse_Success(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
	}

	err := checkVectorResponse(resp)
	if err != nil {
		t.Errorf("Expected no error for 200 response, got %v", err)
	}
}

func TestCheckVectorResponse_Success299(t *testing.T) {
	resp := &http.Response{
		StatusCode: 201,
		Header:     http.Header{},
	}

	err := checkVectorResponse(resp)
	if err != nil {
		t.Errorf("Expected no error for 201 response, got %v", err)
	}
}
