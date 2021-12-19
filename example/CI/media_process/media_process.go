package main

import (
	"context"
	"encoding/xml"
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

func InvokeSnapshotJob() {
	u, _ := url.Parse("https://wwj-bj-1253960454.cos.ap-beijing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-bj-1253960454.ci.ap-beijing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Snapshot",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-beijing",
				Object: "output/abc-${Number}.jpg",
				Bucket: "wwj-bj-1253960454",
			},
			Snapshot: &cos.Snapshot{
				Mode:  "Interval",
				Start: "0",
				Count: "1",
			},
		},
		QueueId: DescribeQueueRes.QueueList[0].QueueId,
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeConcatJob() {
	u, _ := url.Parse("https://wwj-bj-1253960454.cos.ap-beijing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-bj-1253960454.ci.ap-beijing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	concatFragment := make([]cos.ConcatFragment, 0)
	concatFragment = append(concatFragment, cos.ConcatFragment{
		Url:       "https://wwj-bj-1253960454.cos.ap-beijing.myqcloud.com/input/117374C.mp4",
		StartTime: "0",
		EndTime:   "10",
	})
	concatFragment = append(concatFragment, cos.ConcatFragment{
		Url:       "https://wwj-bj-1253960454.cos.ap-beijing.myqcloud.com/input/117374C.mp4",
		StartTime: "20",
		EndTime:   "30",
	})
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Concat",
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-beijing",
				Object: "output/go_117374C.mp4",
				Bucket: "wwj-bj-1253960454",
			},
			ConcatTemplate: &cos.ConcatTemplate{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.265",
				},
				Audio: &cos.Audio{
					//Codec: "AAC",
				},
				ConcatFragment: concatFragment,
			},
		},
		QueueId: DescribeQueueRes.QueueList[0].QueueId,
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeTranscodeJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Transcode",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/go_117374C.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Transcode: &cos.Transcode{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.264",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
				TimeInterval: &cos.TimeInterval{
					Start:    "10",
					Duration: "",
				},
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeMultiJobs() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMultiMediaJobsOptions{
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: []cos.MediaProcessJobOperation{
			cos.MediaProcessJobOperation{
				Tag: "Snapshot",
				Output: &cos.JobOutput{
					Region: "ap-chongqing",
					Object: "output/go_${Number}.mp4",
					Bucket: "wwj-cq-1253960454",
				},
				Snapshot: &cos.Snapshot{
					Mode:  "Interval",
					Start: "0",
					Count: "1",
				},
			},
			cos.MediaProcessJobOperation{
				Tag: "Transcode",
				Output: &cos.JobOutput{
					Region: "ap-chongqing",
					Object: "output/go_117374C.mp4",
					Bucket: "wwj-cq-1253960454",
				},
				Transcode: &cos.Transcode{
					Container: &cos.Container{
						Format: "mp4",
					},
					Video: &cos.Video{
						Codec: "H.264",
					},
					Audio: &cos.Audio{
						Codec: "AAC",
					},
					TimeInterval: &cos.TimeInterval{
						Start:    "10",
						Duration: "",
					},
				},
			},
			cos.MediaProcessJobOperation{
				Tag: "Animation",
				Output: &cos.JobOutput{
					Region: "ap-chongqing",
					Object: "output/go_117374C.gif",
					Bucket: "wwj-cq-1253960454",
				},
				Animation: &cos.Animation{
					Container: &cos.Container{
						Format: "gif",
					},
					Video: &cos.AnimationVideo{
						Codec:                   "gif",
						AnimateOnlyKeepKeyFrame: "true",
					},
					TimeInterval: &cos.TimeInterval{
						Start:    "0",
						Duration: "",
					},
				},
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMultiMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	for k, job := range createJobRes.JobsDetail {
		fmt.Printf("job:%d, %+v\n", k, job)
	}
}

func JobNotifyCallback() {
	taskBody := "<Response><EventName>TaskFinish</EventName><JobsDetail><Code>Success</Code><CreationTime>2021-09-14T14:38:59+0800</CreationTime><EndTime>2021-09-14T14:39:59+0800</EndTime><Input><BucketId></BucketId><Object>e28ef88b61ed41b7a6ea76a8147f4f9e</Object><Region>ap-guangzhou</Region></Input><JobId>j7123e78c152611ec9e8c1ba6632b91a8</JobId><Message/><Operation><MediaInfo><Format><Bitrate>3775.738000</Bitrate><Duration>143.732000</Duration><FormatLongName>QuickTime / MOV</FormatLongName><FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName><NumProgram>0</NumProgram><NumStream>2</NumStream><Size>67836813</Size><StartTime>0.000000</StartTime></Format><Stream><Audio><Bitrate>125.049000</Bitrate><Channel>2</Channel><ChannelLayout>stereo</ChannelLayout><CodecLongName>AAC (Advanced Audio Coding)</CodecLongName><CodecName>aac</CodecName><CodecTag>0x6134706d</CodecTag><CodecTagString>mp4a</CodecTagString><CodecTimeBase>1/44100</CodecTimeBase><Duration>143.730998</Duration><Index>1</Index><Language>und</Language><SampleFmt>fltp</SampleFmt><SampleRate>44100</SampleRate><StartTime>0.000000</StartTime><Timebase>1/44100</Timebase></Audio><Subtitle/><Video><AvgFps>25/1</AvgFps><Bitrate>3645.417000</Bitrate><CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName><CodecName>h264</CodecName><CodecTag>0x31637661</CodecTag><CodecTagString>avc1</CodecTagString><CodecTimeBase>1/12800</CodecTimeBase><Dar>9:16</Dar><Duration>143.680000</Duration><Fps>25.500000</Fps><HasBFrame>2</HasBFrame><Height>3412</Height><Index>0</Index><Language>und</Language><Level>51</Level><NumFrames>3592</NumFrames><PixFormat>yuv420p</PixFormat><Profile>High</Profile><RefFrames>1</RefFrames><Rotation>0.000000</Rotation><Sar>2559:2560</Sar><StartTime>0.000000</StartTime><Timebase>1/12800</Timebase><Width>1920</Width></Video></Stream></MediaInfo><MediaResult><OutputFile><Bucket></Bucket><ObjectName>f89f22f7a8be4f478434da58eb11d7da.mp4</ObjectName><ObjectPrefix/><Region>ap-guangzhou</Region></OutputFile></MediaResult><Output><Bucket></Bucket><Object>f89f22f7a8be4f478434da58eb11d7da.mp4</Object><Region>ap-guangzhou</Region></Output><TemplateId>t064fb9214850f49aaac44b5561a7b0b3b</TemplateId><TemplateName>MP4-FHD</TemplateName></Operation><QueueId>p6f358a37bf9442ad8f859db055cd0edb</QueueId><StartTime>2021-09-14T14:38:59+0800</StartTime><State>Success</State><Tag>Transcode</Tag></JobsDetail></Response>"
	var body cos.MediaProcessJobsNotifyBody
	err := xml.Unmarshal([]byte(taskBody), &body)
	if err != nil {
		fmt.Println(fmt.Sprintf("err:%v", err))
	} else {
		fmt.Println(fmt.Sprintf("body:%v", body))
	}
}

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

func InvokeSpriteSnapshotJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Snapshot",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region:       "ap-chongqing",
				Object:       "output/abc-${Number}.jpg",
				Bucket:       "wwj-cq-1253960454",
				SpriteObject: "output/sprite-${Number}.jpg",
			},
			Snapshot: &cos.Snapshot{
				Mode:            "Interval",
				Start:           "0",
				Count:           "100",
				SnapshotOutMode: "SnapshotAndSprite", // OnlySnapshot OnlySprite
				SpriteSnapshotConfig: &cos.SpriteSnapshotConfig{
					CellHeight: "128",
					CellWidth:  "128",
					Color:      "Black",
					Columns:    "3",
					Lines:      "10",
					Margin:     "2",
				},
			},
		},
		QueueId: DescribeQueueRes.QueueList[0].QueueId,
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeSegmentJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Segment",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/abc-${Number}.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Segment: &cos.Segment{
				Format:   "mp4",
				Duration: "10",
			},
		},
		QueueId: DescribeQueueRes.QueueList[0].QueueId,
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func DescribeMultiMediaJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "Segment",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/abc-${Number}.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Segment: &cos.Segment{
				Format:   "mp4",
				Duration: "10",
			},
		},
		QueueId: DescribeQueueRes.QueueList[0].QueueId,
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	jobids := make([]string, 0)
	jobids = append(jobids, createJobRes.JobsDetail.JobId)
	jobids = append(jobids, "a")
	jobids = append(jobids, "b")
	jobids = append(jobids, "c")
	DescribeJobRes, _, err := c.CI.DescribeMultiMediaJob(context.Background(), jobids)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func GetPrivateM3U8() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	getPrivateM3U8Opt := &cos.GetPrivateM3U8Options{
		Expires: 3600,
	}
	getPrivateM3U8Res, err := c.CI.GetPrivateM3U8(context.Background(), "output/linkv.m3u8", getPrivateM3U8Opt)
	log_status(err)
	fmt.Printf("%+v\n", getPrivateM3U8Res)
}

func InvokeVideoMontageJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "VideoMontage",
		Input: &cos.JobInput{
			Object: "input/117374C.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/go_117374C.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			VideoMontage: &cos.VideoMontage{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.VideoMontageVideo{
					Codec: "H.264",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeVoiceSeparateJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "VoiceSeparate",
		Input: &cos.JobInput{
			Object: "example.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/go_example.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			VoiceSeparate: &cos.VoiceSeparate{
				AudioMode: "AudioAndBackground",
				AudioConfig: &cos.AudioConfig{
					Codec: "AAC",
				},
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeVideoProcessJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "VideoProcess",
		Input: &cos.JobInput{
			Object: "example.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/vp_example.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Transcode: &cos.Transcode{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.264",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
			},
			VideoProcess: &cos.VideoProcess{
				ColorEnhance: &cos.ColorEnhance{
					Enable: "true",
				},
				MsSharpen: &cos.MsSharpen{
					Enable: "true",
				},
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeSDRtoHDRJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "SDRtoHDR",
		Input: &cos.JobInput{
			Object: "linkv.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/sdrtohdr_linkv.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Transcode: &cos.Transcode{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.265",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
			},
			SDRtoHDR: &cos.SDRtoHDR{
				HdrMode: "HLG",
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func InvokeSuperResolutionJob() {
	u, _ := url.Parse("https://wwj-cq-1253960454.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://wwj-cq-1253960454.ci.ap-chongqing.myqcloud.com")
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
	// DescribeMediaProcessQueues
	DescribeQueueOpt := &cos.DescribeMediaProcessQueuesOptions{
		QueueIds:   "",
		PageNumber: 1,
		PageSize:   2,
	}
	DescribeQueueRes, _, err := c.CI.DescribeMediaProcessQueues(context.Background(), DescribeQueueOpt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeQueueRes)
	// CreateMediaJobs
	createJobOpt := &cos.CreateMediaJobsOptions{
		Tag: "SuperResolution",
		Input: &cos.JobInput{
			Object: "100986-2999.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/sp-100986-2999.mp4",
				Bucket: "wwj-cq-1253960454",
			},
			Transcode: &cos.Transcode{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.264",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
			},
			SuperResolution: &cos.SuperResolution{
				Resolution:    "hdto4k",
				EnableScaleUp: "true",
			},
		},
		QueueId: "paaf4fce5521a40888a3034a5de80f6ca",
	}
	createJobRes, _, err := c.CI.CreateMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)

	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), createJobRes.JobsDetail.JobId)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func main() {
	// InvokeSnapshotJob()
	// InvokeConcatJob()
	// InvokeTranscodeJob()
	// InvokeMultiJobs()
	// JobNotifyCallback()
	// WorkflowExecutionNotifyCallback()
	// InvokeSpriteSnapshotJob()
	// InvokeSegmentJob()
	// DescribeMultiMediaJob()
	// GetPrivateM3U8()
	// InvokeVideoMontageJob()
	// InvokeVoiceSeparateJob()
	// InvokeVideoProcessJob()
	// InvokeSDRtoHDRJob()
	InvokeSuperResolutionJob()
}
