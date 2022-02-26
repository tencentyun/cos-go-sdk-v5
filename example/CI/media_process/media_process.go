package main

import (
	"context"
	"encoding/json"
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
	taskBody := "<Response><JobsDetail><Code>Success</Code><CreationTime>2022-02-09T11:25:43+0800</CreationTime><EndTime>2022-02-09T11:25:47+0800</EndTime><Input><BucketId>wwj-cq-1253960454</BucketId><Object>input/117374C.mp4</Object><Region>ap-chongqing</Region></Input><JobId>jf6717076895711ecafdd594be6cca70c</JobId><Message/><Operation><MediaInfo><Format><Bitrate>215.817000</Bitrate><Duration>96.931000</Duration><FormatLongName>QuickTime / MOV</FormatLongName><FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName><NumProgram>0</NumProgram><NumStream>1</NumStream><Size>2614921</Size><StartTime>0.000000</StartTime></Format><Stream><Audio/><Subtitle/><Video><AvgFps>29.970030</AvgFps><Bitrate>212.875000</Bitrate><CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName><CodecName>h264</CodecName><CodecTag>0x31637661</CodecTag><CodecTagString>avc1</CodecTagString><CodecTimeBase>1/30000</CodecTimeBase><Dar>16:9</Dar><Duration>96.930167</Duration><Fps>29.970030</Fps><HasBFrame>2</HasBFrame><Height>400</Height><Index>0</Index><Language>und</Language><Level>30</Level><NumFrames>2905</NumFrames><PixFormat>yuv420p</PixFormat><Profile>High</Profile><RefFrames>1</RefFrames><Rotation>0.000000</Rotation><Sar>640:639</Sar><StartTime>0.000000</StartTime><Timebase>1/30000</Timebase><Width>710</Width></Video></Stream></MediaInfo><MediaResult><OutputFile><Bucket>wwj-cq-1253960454</Bucket><Md5Info><Md5>38f0b40c78562f819421137541043f09</Md5><ObjectName>output/go_117374C.mp4</ObjectName></Md5Info><ObjectName>output/go_117374C.mp4</ObjectName><ObjectPrefix/><Region>ap-chongqing</Region><SpriteOutputFile><Bucket/><Md5Info/><ObjectName/><ObjectPrefix/><Region/></SpriteOutputFile></OutputFile></MediaResult><Output><Bucket>wwj-cq-1253960454</Bucket><Object>output/go_117374C.mp4</Object><Region>ap-chongqing</Region></Output><Transcode><Audio><Bitrate/><Channels/><Codec>AAC</Codec><KeepTwoTracks>false</KeepTwoTracks><Profile/><Remove>false</Remove><SampleFormat/><Samplerate>44100</Samplerate><SwitchTrack>false</SwitchTrack></Audio><Container><Format>mp4</Format></Container><TimeInterval><Duration/><Start>10</Start></TimeInterval><TransConfig><AdjDarMethod/><AudioBitrateAdjMethod/><DeleteMetadata>false</DeleteMetadata><IsCheckAudioBitrate>false</IsCheckAudioBitrate><IsCheckReso>false</IsCheckReso><IsCheckVideoBitrate>false</IsCheckVideoBitrate><IsHdr2Sdr>false</IsHdr2Sdr><IsStreamCopy>false</IsStreamCopy><ResoAdjMethod/><VideoBitrateAdjMethod/></TransConfig><Video><AnimateFramesPerSecond/><AnimateOnlyKeepKeyFrame>false</AnimateOnlyKeepKeyFrame><AnimateTimeIntervalOfFrame/><Bitrate/><Bufsize/><Codec>H.264</Codec><Crf>25</Crf><Crop/><Fps/><Gop/><Height/><Interlaced>false</Interlaced><LongShortMode>false</LongShortMode><Maxrate/><Pad/><Pixfmt/><Preset>medium</Preset><Profile>high</Profile><Quality/><Remove>false</Remove><ScanMode/><SliceTime>5</SliceTime><Width/></Video></Transcode></Operation><Progress>100</Progress><QueueId>paaf4fce5521a40888a3034a5de80f6ca</QueueId><StartTime>2022-02-09T11:25:43+0800</StartTime><State>Success</State><Tag>Transcode</Tag></JobsDetail></Response>"
	var body cos.MediaProcessJobsNotifyBody
	err := xml.Unmarshal([]byte(taskBody), &body)
	if err != nil {
		fmt.Println(fmt.Sprintf("err:%v", err))
	} else {
		fmt.Println(fmt.Sprintf("body:%+v", body))
		fmt.Println(fmt.Sprintf("mediaInfo:%+v", body.JobsDetail.Operation.MediaInfo))
		fmt.Println(fmt.Sprintf("mediaResult:%+v", body.JobsDetail.Operation.MediaResult))
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

func TriggerWorkflow() {
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
	triggerWorkflowOpt := &cos.TriggerWorkflowOptions{
		WorkflowId: "w18fd791485904afba3ab07ed57d9cf1e",
		Object:     "100986-2999.mp4",
	}
	triggerWorkflowRes, _, err := c.CI.TriggerWorkflow(context.Background(), triggerWorkflowOpt)
	log_status(err)
	fmt.Printf("%+v\n", triggerWorkflowRes)
}

func DescribeWorkflowExecutions() {
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
	describeWorkflowExecutionsOpt := &cos.DescribeWorkflowExecutionsOptions{
		WorkflowId: "w18fd791485904afba3ab07ed57d9cf1e",
	}
	describeWorkflowExecutionsRes, _, err := c.CI.DescribeWorkflowExecutions(context.Background(), describeWorkflowExecutionsOpt)
	log_status(err)
	fmt.Printf("%+v\n", describeWorkflowExecutionsRes)
}

func DescribeMultiWorkflowExecution() {
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
	describeWorkflowExecutionsRes, _, err := c.CI.DescribeWorkflowExecution(context.Background(), "i00689df860ad11ec9c5952540019ee59")
	log_status(err)
	a, _ := json.Marshal(describeWorkflowExecutionsRes)
	fmt.Println(string(a))
	fmt.Printf("%+v\n", describeWorkflowExecutionsRes)
}

func InvokeASRJob() {
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
	// CreateMediaJobs
	createJobOpt := &cos.CreateASRJobsOptions{
		Tag: "SpeechRecognition",
		Input: &cos.JobInput{
			Object: "abc.mp3",
		},
		Operation: &cos.ASRJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "music.txt",
				Bucket: "wwj-cq-1253960454",
			},
			SpeechRecognition: &cos.SpeechRecognition{
				ChannelNum:      "1",
				EngineModelType: "8k_zh",
			},
		},
		QueueId: "p1db6a1a76ff04806b6af0d96e9bc80ab",
	}
	createJobRes, _, err := c.CI.CreateASRJobs(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
	DescribeJobRes, _, err := c.CI.DescribeMultiASRJob(context.Background(), []string{createJobRes.JobsDetail.JobId})
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

func DescribeASRJob() {
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
	DescribeJobRes, _, err := c.CI.DescribeMultiASRJob(context.Background(), []string{"sa59de0a06a4711ec9b81df13272c69a9"})
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail[0].Operation.SpeechRecognitionResult)
}

func DescribeJob() {
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
	// DescribeMediaJobs
	DescribeJobRes, _, err := c.CI.DescribeMediaJob(context.Background(), "jf6717076895711ecafdd594be6cca70c")
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail.Operation.MediaInfo)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail.Operation.MediaResult)
}

func GenerateMediaInfo() {
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
	opt := &cos.GenerateMediaInfoOptions{
		Input: &cos.JobInput{
			Object: "video02.mp4",
		},
	}
	// DescribeMediaJobs
	res, _, err := c.CI.GenerateMediaInfo(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", res)
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
	// InvokeSuperResolutionJob()
	// TriggerWorkflow()
	// DescribeWorkflowExecutions()
	// DescribeMultiWorkflowExecution()
	// InvokeASRJob()
	// DescribeASRJob()
	// DescribeJob()
	GenerateMediaInfo()
}
