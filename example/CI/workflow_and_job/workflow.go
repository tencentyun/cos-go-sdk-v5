package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
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

// DescribeWorkflow 查询工作流
// https://cloud.tencent.com/document/product/460/76857
func DescribeWorkflow() {
	c := getClient()
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
// https://cloud.tencent.com/document/product/460/76860
func DeleteWorkflow() {
	c := getClient()
	DescribeWorkflowRes, _, err := c.CI.DeleteMediaWorkflow(context.Background(), "w843779f0b22f49bbb7a189778d865059")
	log_status(err)
	fmt.Printf("%+v\n", DescribeWorkflowRes)
}

// CreateWorkflow 创建工作流
// https://cloud.tencent.com/document/product/460/76856
func CreateWorkflow() {
	c := getClient()
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
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "trans1.mp4"}}},
					"Snapshot_1581665960536": cos.Node{Type: "Snapshot", Operation: &cos.NodeOperation{TemplateId: "t07740e32081b44ad7a0aea03adcffd54a",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "snapshot-${number}.jpg"}}},
				},
			},
		},
	}
	createWorkflowRes, _, err := c.CI.CreateMediaWorkflow(context.Background(), createWorkflowOpt)
	log_status(err)
	fmt.Printf("%+v\n", createWorkflowRes.MediaWorkflow)
}

// UpdateWorkflow 更新工作流
// https://cloud.tencent.com/document/product/460/76861
func UpdateWorkflow() {
	c := getClient()
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
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "trans1.mp4"}}},
					"Snapshot_1581665960536": cos.Node{Type: "Snapshot", Operation: &cos.NodeOperation{TemplateId: "t07740e32081b44ad7a0aea03adcffd54a",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "snapshot-${number}.jpg"}}},
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
// https://cloud.tencent.com/document/product/460/76856
func CreateStreamWorkflow() {
	c := getClient()
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
						Output:               &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "${InputPath}/${InputName}._${RunId}.${ext}"},
						StreamPackConfigInfo: &cos.NodeStreamPackConfigInfo{PackType: "HLS", IgnoreFailedStream: true}}},
					"VideoStream_1581665960536": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t03e862f296fba4152a1dd186b4ad5f64b",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"VideoStream_1581665960537": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t09f9da59ed3c44ecd8ea1778e5ce5669c",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"StreamPack_1581665960538": cos.Node{Type: "StreamPack", Operation: &cos.NodeOperation{StreamPackInfo: hpi}},
				},
			},
		},
	}
	createWorkflowRes, _, err := c.CI.CreateMediaWorkflow(context.Background(), createWorkflowOpt)
	log_status(err)
	fmt.Printf("%+v\n", createWorkflowRes.MediaWorkflow)
}

// UpdatStreamWorkflow 更新自适应码流工作流
// https://cloud.tencent.com/document/product/460/76861
func UpdatStreamWorkflow() {
	c := getClient()
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
						Output:               &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "${InputPath}/${InputName}._${RunId}.${ext}"},
						StreamPackConfigInfo: &cos.NodeStreamPackConfigInfo{PackType: "HLS", IgnoreFailedStream: true}}},
					"VideoStream_1581665960536": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t03e862f296fba4152a1dd186b4ad5f64b",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"VideoStream_1581665960537": cos.Node{Type: "VideoStream", Operation: &cos.NodeOperation{TemplateId: "t09f9da59ed3c44ecd8ea1778e5ce5669c",
						Output: &cos.NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "${RunId}_Substream_1/video.m3u8"}}},
					"StreamPack_1581665960538": cos.Node{Type: "StreamPack", Operation: &cos.NodeOperation{StreamPackInfo: hpi}},
				},
			},
		},
	}
	WorkflowId := "w5fecd57f8a7745b3ac8143b211613789"
	updateWorkflowRes, _, err := c.CI.UpdateMediaWorkflow(context.Background(), updateWorkflowOpt, WorkflowId)
	log_status(err)
	fmt.Printf("%+v\n", updateWorkflowRes.MediaWorkflow)
}

// 启用工作流
func ActiveWorkflow() {
	c := getClient()
	WorkflowId := "w8d1f24d05b434b17b491555496acf11d"
	_, err := c.CI.ActiveMediaWorkflow(context.Background(), WorkflowId)
	log_status(err)
}

// 停用工作流
func PausedWorkflow() {
	c := getClient()
	WorkflowId := "w8d1f24d05b434b17b491555496acf11d"
	_, err := c.CI.PausedMediaWorkflow(context.Background(), WorkflowId)
	log_status(err)
}

// TriggerWorkflow 测试工作流
// https://cloud.tencent.com/document/product/460/76864
func TriggerWorkflow() {
	c := getClient()
	triggerWorkflowOpt := &cos.TriggerWorkflowOptions{
		WorkflowId: "w18fd791485904afba3ab07ed57d9cf1e",
		Object:     "100986-2999.mp4",
	}
	triggerWorkflowRes, _, err := c.CI.TriggerWorkflow(context.Background(), triggerWorkflowOpt)
	log_status(err)
	fmt.Printf("%+v\n", triggerWorkflowRes)
}

// DescribeWorkflowExecutions 获取工作流实例详情列表
// https://cloud.tencent.com/document/product/460/80050
func DescribeWorkflowExecutions() {
	c := getClient()
	describeWorkflowExecutionsOpt := &cos.DescribeWorkflowExecutionsOptions{
		WorkflowId: "w18fd791485904afba3ab07ed57d9cf1e",
	}
	describeWorkflowExecutionsRes, _, err := c.CI.DescribeWorkflowExecutions(context.Background(), describeWorkflowExecutionsOpt)
	log_status(err)
	fmt.Printf("%+v\n", describeWorkflowExecutionsRes)
}

// DescribeMultiWorkflowExecution 获取工作流实例详情
// https://cloud.tencent.com/document/product/460/80044
func DescribeMultiWorkflowExecution() {
	c := getClient()
	describeWorkflowExecutionsRes, _, err := c.CI.DescribeWorkflowExecution(context.Background(), "i00689df860ad11ec9c5952540019ee59")
	log_status(err)
	a, _ := json.Marshal(describeWorkflowExecutionsRes)
	fmt.Println(string(a))
	fmt.Printf("%+v\n", describeWorkflowExecutionsRes)
}

// WorkflowExecutionNotifyCallback TODO
func WorkflowExecutionNotifyCallback() {
	workflowExecutionBody := "<Response><EventName>WorkflowFinish</EventName><WorkflowExecution><RunId>i70ae991a152911ecb184525400a8700f</RunId><BucketId></BucketId><Object>62ddbc1245.mp4</Object><CosHeaders><Key>x-cos-meta-id</Key><Value>62ddbc1245</Value></CosHeaders><CosHeaders><Key>Content-Type</Key><Value>video/mp4</Value></CosHeaders><WorkflowId>w29ba54d02b7340dd9fb44eb5beb786b9</WorkflowId><WorkflowName></WorkflowName><CreateTime>2021-09-14 15:00:26+0800</CreateTime><State>Success</State><Tasks><Type>Transcode</Type><CreateTime>2021-09-14 15:00:27+0800</CreateTime><EndTime>2021-09-14 15:00:42+0800</EndTime><State>Success</State><JobId>j70bab192152911ecab79bba409874f7f</JobId><Name>Transcode_1607323983818</Name><TemplateId>t088613dea8d564a9ba7e6b02cbd5de877</TemplateId><TemplateName>HLS-FHD</TemplateName></Tasks></WorkflowExecution></Response>"
	var body cos.WorkflowExecutionNotifyBody
	err := xml.Unmarshal([]byte(workflowExecutionBody), &body)
	if err != nil {
		fmt.Println(fmt.Sprintf("err:%v", err))
	} else {
		fmt.Println(fmt.Sprintf("body:%v", body))
	}
}

func main() {
}
