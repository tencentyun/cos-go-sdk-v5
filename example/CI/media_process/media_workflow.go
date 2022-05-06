package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

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

// DescribeWorkflow 查询工作流
func DescribeWorkflow() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	opt := &cos.DescribeMediaWorkflowOptions{
		Ids:        "w93aa43ba105347169fa093ed857b2a90,abc,123",
		PageNumber: 1,
		PageSize:   5,
	}
	DescribeWorkflowRes, _, err := c.CI.DescribeMediaWorkflow(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

// DeleteWorkflow 删除工作流
func DeleteWorkflow() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	DescribeWorkflowRes, _, err := c.CI.DeleteMediaWorkflow(context.Background(), "w843779f0b22f49bbb7a189778d865059")
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

// CreateWorkflow 创建工作流
func CreateWorkflow() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// CreateMediaWorkflow
	rand.Seed(time.Now().UnixNano())
	createWorkflowOpt := &cos.CreateMediaWorkflowOptions{
		MediaWorkflow: &cos.MediaWorkflow{
			Name:  "workflow-" + strconv.Itoa(rand.Intn(100)),
			State: "Active",
			Topology: &cos.Topology{
				Dependencies: map[string]string{"Start": "Transcode_1581665960537", "Transcode_1581665960537": "Snapshot_1581665960536",
					"Snapshot_1581665960536": "End"},
				Nodes: map[string]cos.Node{"Start": cos.Node{Type: "Start", Input: &cos.NodeInput{QueueId: "p09d709939fef48a0a5c247ef39d90cec",
					ObjectPrefix: "wk-test", ExtFilter: &cos.ExtFilter{State: "On", Custom: "true", CustomExts: "mp4"}}},
					"Transcode_1581665960537": cos.Node{Type: "Transcode", Operation: &cos.NodeOperation{TemplateId: "t01e57db1c2d154d2fb57aa5de9313a897",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "trans1.mp4"}}},
					"Snapshot_1581665960536": cos.Node{Type: "Snapshot", Operation: &cos.NodeOperation{TemplateId: "t07740e32081b44ad7a0aea03adcffd54a",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "snapshot-${number}.jpg"}}},
				},
			},
		},
	}
	createWorkflowRes, _, err := c.CI.CreateMediaWorkflow(context.Background(), createWorkflowOpt)
	log_status(err)
	fmt.Printf("%+v\n", createWorkflowRes.MediaWorkflow)

	// DescribeMediaWorkflow
	if createWorkflowRes.MediaWorkflow != nil {
		opt := &cos.DescribeMediaWorkflowOptions{
			Ids: createWorkflowRes.MediaWorkflow.WorkflowId,
		}
		DescribeWorkflowRes, _, err := c.CI.DescribeMediaWorkflow(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeWorkflowRes)
	}
}

// UpdateWorkflow TODO
func UpdateWorkflow() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// UpdateMediaWorkflow
	rand.Seed(time.Now().UnixNano())
	updateWorkflowOpt := &cos.CreateMediaWorkflowOptions{
		MediaWorkflow: &cos.MediaWorkflow{
			Name:  "workflow-" + strconv.Itoa(rand.Intn(100)),
			State: "Paused",
			Topology: &cos.Topology{
				Dependencies: map[string]string{"Start": "Transcode_1581665960537", "Transcode_1581665960537": "Snapshot_1581665960536",
					"Snapshot_1581665960536": "End"},
				Nodes: map[string]cos.Node{"Start": cos.Node{Type: "Start", Input: &cos.NodeInput{QueueId: "p09d709939fef48a0a5c247ef39d90cec",
					ObjectPrefix: "wk-test", ExtFilter: &cos.ExtFilter{State: "On", Custom: "true", CustomExts: "mp4"}}},
					"Transcode_1581665960537": cos.Node{Type: "Transcode", Operation: &cos.NodeOperation{TemplateId: "t01e57db1c2d154d2fb57aa5de9313a897",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "trans1.mp4"}}},
					"Snapshot_1581665960536": cos.Node{Type: "Snapshot", Operation: &cos.NodeOperation{TemplateId: "t07740e32081b44ad7a0aea03adcffd54a",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "snapshot-${number}.jpg"}}},
				},
			},
		},
	}
	WorkflowId := "web6ac56c1ef54dbfa44d7f4103203be9"
	updateWorkflowRes, _, err := c.CI.UpdateMediaWorkflow(context.Background(), updateWorkflowOpt, WorkflowId)
	log_status(err)
	fmt.Printf("%+v\n", updateWorkflowRes.MediaWorkflow)

	opt := &cos.DescribeMediaWorkflowOptions{
		Ids: WorkflowId,
	}
	DescribeWorkflowRes, _, err := c.CI.DescribeMediaWorkflow(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

// CreateStreamWorkflow 创建自适应码流工作流
func CreateStreamWorkflow() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// CreateMediaWorkflow
	rand.Seed(time.Now().UnixNano())
	hpi := &cos.NodeHlsPackInfo{}
	vsc1 := &cos.VideoStreamConfig{VideoStreamName: "VideoStream_1581665960536", BandWidth: "200000"}
	vsc2 := &cos.VideoStreamConfig{VideoStreamName: "VideoStream_1581665960537", BandWidth: "500000"}
	hpi.VideoStreamConfig = append(hpi.VideoStreamConfig, *vsc1, *vsc2)
	createWorkflowOpt := &cos.CreateMediaWorkflowOptions{
		MediaWorkflow: &cos.MediaWorkflow{
			Name:  "workflow-" + strconv.Itoa(rand.Intn(100)),
			State: "Active",
			Topology: &cos.Topology{
				Dependencies: map[string]string{"Start": "StreamPackConfig_1581665960532", "StreamPackConfig_1581665960532": "VideoStream_1581665960536,VideoStream_1581665960537",
					"VideoStream_1581665960536": "StreamPack_1581665960538", "VideoStream_1581665960537": "StreamPack_1581665960538",
					"StreamPack_1581665960538": "End"},
				Nodes: map[string]cos.Node{"Start": cos.Node{Type: "Start", Input: &cos.NodeInput{QueueId: "p09d709939fef48a0a5c247ef39d90cec",
					ObjectPrefix: "wk-test", ExtFilter: &cos.ExtFilter{State: "On", Custom: "true", CustomExts: "mp4"}}},
					"StreamPackConfig_1581665960532": cos.Node{Type: "StreamPackConfig", Operation: &cos.NodeOperation{
						Output:               &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "${InputPath}/${InputName}._${RunId}.${ext}"},
						StreamPackConfigInfo: &cos.NodeStreamPackConfigInfo{PackType: "HLS", IgnoreFailedStream: true}}},
					"VideoStream_1581665960536": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t03e862f296fba4152a1dd186b4ad5f64b",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"VideoStream_1581665960537": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t09f9da59ed3c44ecd8ea1778e5ce5669c",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"StreamPack_1581665960538": cos.Node{Type: "StreamPack", Operation: &cos.NodeOperation{StreamPackInfo: hpi}},
				},
			},
		},
	}
	createWorkflowRes, _, err := c.CI.CreateMediaWorkflow(context.Background(), createWorkflowOpt)
	log_status(err)
	fmt.Printf("%+v\n", createWorkflowRes.MediaWorkflow)

	// DescribeMediaWorkflow
	if createWorkflowRes.MediaWorkflow != nil {
		opt := &cos.DescribeMediaWorkflowOptions{
			Ids: createWorkflowRes.MediaWorkflow.WorkflowId,
		}
		DescribeWorkflowRes, _, err := c.CI.DescribeMediaWorkflow(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeWorkflowRes)
	}
}

// UpdatStreamWorkflow 更新自适应码流工作流
func UpdatStreamWorkflow() {
	u, _ := url.Parse("https://lilang-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://lilang-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// UpdateMediaWorkflow
	rand.Seed(time.Now().UnixNano())
	hpi := &cos.NodeHlsPackInfo{}
	vsc1 := &cos.VideoStreamConfig{VideoStreamName: "VideoStream_1581665960536", BandWidth: "200000"}
	vsc2 := &cos.VideoStreamConfig{VideoStreamName: "VideoStream_1581665960537", BandWidth: "500000"}
	hpi.VideoStreamConfig = append(hpi.VideoStreamConfig, *vsc1, *vsc2)
	updateWorkflowOpt := &cos.CreateMediaWorkflowOptions{
		MediaWorkflow: &cos.MediaWorkflow{
			Name:  "workflow-" + strconv.Itoa(rand.Intn(100)),
			State: "Active",
			Topology: &cos.Topology{
				Dependencies: map[string]string{"Start": "StreamPackConfig_1581665960532", "StreamPackConfig_1581665960532": "VideoStream_1581665960536,VideoStream_1581665960537",
					"VideoStream_1581665960536": "StreamPack_1581665960538", "VideoStream_1581665960537": "StreamPack_1581665960538",
					"StreamPack_1581665960538": "End"},
				Nodes: map[string]cos.Node{"Start": cos.Node{Type: "Start", Input: &cos.NodeInput{QueueId: "p09d709939fef48a0a5c247ef39d90cec",
					ObjectPrefix: "wk-test", ExtFilter: &cos.ExtFilter{State: "On", Custom: "true", CustomExts: "mp4"}}},
					"StreamPackConfig_1581665960532": cos.Node{Type: "StreamPackConfig", Operation: &cos.NodeOperation{
						Output:               &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "${InputPath}/${InputName}._${RunId}.${ext}"},
						StreamPackConfigInfo: &cos.NodeStreamPackConfigInfo{PackType: "HLS", IgnoreFailedStream: true}}},
					"VideoStream_1581665960536": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t03e862f296fba4152a1dd186b4ad5f64b",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"VideoStream_1581665960537": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t09f9da59ed3c44ecd8ea1778e5ce5669c",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "lilang-1253960454", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"StreamPack_1581665960538": cos.Node{Type: "StreamPack", Operation: &cos.NodeOperation{StreamPackInfo: hpi}},
				},
			},
		},
	}
	WorkflowId := "w5fecd57f8a7745b3ac8143b211613789"
	updateWorkflowRes, _, err := c.CI.UpdateMediaWorkflow(context.Background(), updateWorkflowOpt, WorkflowId)
	log_status(err)
	fmt.Printf("%+v\n", updateWorkflowRes.MediaWorkflow)

	opt := &cos.DescribeMediaWorkflowOptions{
		Ids: WorkflowId,
	}
	DescribeWorkflowRes, _, err := c.CI.DescribeMediaWorkflow(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

func main() {
	// DescribeWorkflow()
	// DeleteWorkflow()
	// CreateWorkflow()
	// UpdateWorkflow()
	// CreateStreamWorkflow()
	UpdatStreamWorkflow()
}
