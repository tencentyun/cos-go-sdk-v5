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

// CancelJob 取消一个未执行的任务
// https://cloud.tencent.com/document/product/460/84803
func CancelJob() {
	c := getClient()
	_, err := c.CI.CancelJob(context.Background(), "j9334ff26044611eebf2565013e042dc9")
	log_status(err)
}

// DescribeJobs 获取符合条件的任务列表
// https://cloud.tencent.com/document/product/460/84766
func DescribeJobs() {
	c := getClient()
	opt := &cos.DescribeJobsOptions{
		QueueId: "pa27b2bd96bef43b6baba820175485532",
		Tag:     "Transcode",
	}
	DescribeJobRes, _, err := c.CI.DescribeJobs(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

// DescribeJob 查询指定任务
// https://cloud.tencent.com/document/product/460/84765
func DescribeJob() {
	c := getClient()
	DescribeJobRes, _, err := c.CI.DescribeJob(context.Background(), "ja507fb3413f711eebccc9dd62ab48c0e")
	log_status(err)
	fmt.Printf("%+v\n", DescribeJobRes.JobsDetail)
}

// InvokeTranscodeJob 提交一个转码任务
// https://cloud.tencent.com/document/product/460/84790
func InvokeTranscodeJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Transcode",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/test.mp4",
				Bucket: "test-1234567890",
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
			UserData: "hello world",
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeVideoEnhanceJob 提交一个画质增强任务
// https://cloud.tencent.com/document/product/460/84775
func InvokeVideoEnhanceJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "VideoEnhance",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/test.mp4",
				Bucket: "test-1234567890",
			},
			VideoEnhance: &cos.VideoEnhance{
				Transcode: &cos.Transcode{
					Container: &cos.Container{
						Format: "mp4",
					},
					Video: &cos.Video{
						Codec:   "H.264",
						Bitrate: "1000",
						Width:   "1280",
						Fps:     "30",
					},
					Audio: &cos.Audio{
						Codec:      "aac",
						Bitrate:    "128",
						Samplerate: "44100",
						Channels:   "4",
					},
				},
				SuperResolution: &cos.SuperResolution{
					Resolution:    "sdtohd",
					EnableScaleUp: "true",
					Version:       "Enhance",
				},
				ColorEnhance: &cos.ColorEnhance{
					Contrast:   "50",
					Correction: "100",
					Saturation: "100",
				},
				MsSharpen: &cos.MsSharpen{
					SharpenLevel: "5",
				},
				SDRtoHDR: &cos.SDRtoHDR{
					HdrMode: "HDR10",
				},
				FrameEnhance: &cos.FrameEnhance{
					FrameDoubling: "true",
				},
			},
			UserData: "This is my VideoEnhance job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeMediaInfoJob 提交一个获取媒体信息任务
// https://cloud.tencent.com/document/product/460/84776
func InvokeMediaInfoJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "MediaInfo",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeVideoMontageJob 提交一个精彩集锦任务
// https://cloud.tencent.com/document/product/460/84778
func InvokeVideoMontageJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "VideoMontage",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/test.mp4",
				Bucket: "test-1234567890",
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
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeVideoTagJob 提交一个视频标签任务
// https://cloud.tencent.com/document/product/460/84779
func InvokeVideoTagJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "VideoTag",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			VideoTag: &cos.VideoTag{
				Scenario: "Stream",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSnapshotJob 提交一个截图任务
// https://cloud.tencent.com/document/product/460/84780
func InvokeSnapshotJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Snapshot",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/abc-${Number}.jpg",
				Bucket: "test-1234567890",
			},
			Snapshot: &cos.Snapshot{
				Mode:  "Interval",
				Start: "0",
				Count: "1",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSnapshotJob 提交一个截图任务(雪碧图)
// https://cloud.tencent.com/document/product/460/84780
func InvokeSpriteSnapshotJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Snapshot",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region:       "ap-chongqing",
				Object:       "output/abc-${Number}.jpg",
				Bucket:       "test-1234567890",
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
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeWatermarkJob 提交一个视频明水印任务
// https://cloud.tencent.com/document/product/460/84781
func InvokeWatermarkJob() {
	c := getClient()
	w := cos.Watermark{
		Type:      "Text",
		LocMode:   "Absolute",
		Dx:        "20",
		Dy:        "20",
		Pos:       "TopRight",
		StartTime: "5",
		EndTime:   "20",
		Text: &cos.Text{
			Text:         "helloworld",
			FontSize:     "12",
			FontType:     "simfang.ttf",
			FontColor:    "0xff0000",
			Transparency: "100",
		},
	}
	ws := []cos.Watermark{}
	ws = append(ws, w)
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Watermark",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/test.mp4",
				Bucket: "test-1234567890",
			},
			Watermark: ws,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeQualityEstimateJob 提交一个视频质量分析任务
// https://cloud.tencent.com/document/product/460/84783
func InvokeQualityEstimateJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "QualityEstimate",
		Input: &cos.JobInput{
			Object: "input/dog.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			UserData: "This is my QualityEstimate job",
			JobLevel: 1,
			QualityEstimateConfig: &cos.QualityEstimateConfig{
				Rotate: "90", // 只支持0 90 180 270
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeAnimationJob 提交一个动图任务
// https://cloud.tencent.com/document/product/460/84784
func InvokeAnimationJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Animation",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/game.jpg",
				Bucket: "test-12345678900",
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
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeDigitalWatermarkJob 提交一个添加数字水印任务
// https://cloud.tencent.com/document/product/460/84785
func InvokeDigitalWatermarkJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "DigitalWatermark",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/test.mp4",
				Bucket: "test-1234567890",
			},
			DigitalWatermark: &cos.DigitalWatermark{
				Message: "HelloWorld",
				Type:    "Text",
				Version: "V1",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeExtractDigitalWatermarkJob 提交一个提取数字水印任务
// https://cloud.tencent.com/document/product/460/84786
func InvokeExtractDigitalWatermarkJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "ExtractDigitalWatermark",
		Input: &cos.JobInput{
			Object: "output/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			ExtractDigitalWatermark: &cos.ExtractDigitalWatermark{
				Type:    "Text",
				Version: "V1",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeStreamExtractJob 提交一个音视频流分离任务
// https://cloud.tencent.com/document/product/460/84787
func InvokeStreamExtractJob() {
	c := getClient()
	streamEtract := make([]cos.StreamExtract, 0)
	streamEtract = append(streamEtract, cos.StreamExtract{
		Index:  "1",
		Object: "stream/video02_1.mp4",
	})

	createJobOpt := &cos.CreateJobsOptions{
		Tag: "StreamExtract",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region:        "ap-chongqing",
				Bucket:        "test-1234567890",
				StreamExtract: streamEtract,
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeConcatJob 提交一个拼接任务
// https://cloud.tencent.com/document/product/460/84788
func InvokeConcatJob() {
	c := getClient()
	concatFragment := make([]cos.ConcatFragment, 0)
	concatFragment = append(concatFragment, cos.ConcatFragment{
		Url:           "https://test-1234567890.cos.ap-beijing.myqcloud.com/input/test1.mp4",
		StartTime:     "0",
		EndTime:       "10",
		FragmentIndex: "0",
	})
	concatFragment = append(concatFragment, cos.ConcatFragment{
		Url:           "https://test-1234567890.cos.ap-beijing.myqcloud.com/input/test2.mp4",
		StartTime:     "20",
		EndTime:       "30",
		FragmentIndex: "1",
	})
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Concat",
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-beijing",
				Object: "output/go_test.mp4",
				Bucket: "test-1234567890",
			},
			ConcatTemplate: &cos.ConcatTemplate{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec: "H.264",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
				ConcatFragment: concatFragment,
				SceneChangeInfo: &cos.SceneChangeInfo{
					Mode: "GRADIENT", // Default：不添加转场特效; FADE：淡入淡出; GRADIENT：渐变
					Time: "4.5",      // 取值范围：(0, 5], 支持小数, 默认值3
				},
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSegmentJob 提交一个转封装任务
// https://cloud.tencent.com/document/product/460/84789
func InvokeSegmentJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Segment",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/m3u8/a",
				Bucket: "test-1234567890",
			},
			Segment: &cos.Segment{
				Format:   "hls",
				Duration: "10",
				HlsEncrypt: &cos.HlsEncrypt{
					IsHlsEncrypt: true,
					UriKey:       "http://abc.com/",
				},
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSmartCoverJob 提交一个智能封面任务
// https://cloud.tencent.com/document/product/460/84791
func InvokeSmartCoverJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "SmartCover",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			SmartCover: &cos.NodeSmartCover{
				Format:           "jpg",
				Height:           "1280",
				Width:            "720",
				Count:            "2",
				DeleteDuplicates: "true",
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "output/mc-${number}.jpg",
				Bucket: "test-1234567890",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokePicProcessJob 提交一个图片处理任务
// https://cloud.tencent.com/document/product/460/84793
func InvokePicProcessJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "PicProcess",
		Input: &cos.JobInput{
			Object: "pic/cup.jpeg",
		},
		Operation: &cos.MediaProcessJobOperation{
			PicProcess: &cos.PicProcess{
				IsPicInfo:   "true",
				ProcessRule: "imageMogr2/format/jpg/interlace/0/quality/100",
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "test.jpg",
			},
		},
		CallBack: "https://demo.org/callback",
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeTranslationJob 提交一个翻译任务
// https://cloud.tencent.com/document/product/460/84799
func InvokeTranslationJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Translation",
		Input: &cos.JobInput{
			Object:    "input/translate-en.txt",
			Lang:      "en",
			Type:      "txt",
			BasicType: "",
		},
		Operation: &cos.MediaProcessJobOperation{
			Translation: &cos.Translation{
				Lang: "zh",
				Type: "txt",
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "output/out.txt",
			},
			UserData: "This is my Translation job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeWordsGeneralizeJob 提交一个分词任务
// https://cloud.tencent.com/document/product/460/84800
func InvokeWordsGeneralizeJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "WordsGeneralize",
		Input: &cos.JobInput{
			Object: "input/WordsGeneralize.txt",
		},
		Operation: &cos.MediaProcessJobOperation{
			WordsGeneralize: &cos.WordsGeneralize{
				NerMethod: "DL",
				SegMethod: "MIX",
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "output/out.txt",
			},
			UserData: "This is my WordsGeneralize job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeVideoTargetRecJob 提交一个视频目标检测任务
// https://cloud.tencent.com/document/product/460/84801
func InvokeVideoTargetRecJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "VideoTargetRec",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			VideoTargetRec: &cos.VideoTargetRec{
				Body: "true",
				Pet:  "true",
				Car:  "true",
			},
			UserData: "This is my VideoTargetRec job",
			JobLevel: 1,
		},
	}
	// 注意这里是AI
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSplitVideoPartsJob 提交一个视频拆条任务
// https://cloud.tencent.com/document/product/460/90888
func InvokeSplitVideoPartsJob() {
	c := getClient()
	// CreateJob
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "SplitVideoParts",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			SplitVideoParts: &cos.SplitVideoParts{
				Mode: "SHOTDETECT",
			},
			UserData: "This is my SplitVideoParts job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSegmentVideoBodyJob 提交一个视频人像抠图任务
// https://cloud.tencent.com/document/product/460/84802
func InvokeSegmentVideoBodyJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "SegmentVideoBody",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			SegmentVideoBody: &cos.SegmentVideoBody{
				Mode: "Mask",
			},
			UserData: "This is my SegmentVideoBody job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeVoiceSeparateJob 提交一个人声分离任务
// https://cloud.tencent.com/document/product/460/84794
func InvokeVoiceSeparateJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "VoiceSeparate",
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region:   "ap-chongqing",
				Object:   "output/voice.mp3",
				AuObject: "output/au.mp4",
				Bucket:   "test-1234567890",
			},
			VoiceSeparate: &cos.VoiceSeparate{
				AudioMode: "AudioAndBackground",
				AudioConfig: &cos.AudioConfig{
					Codec: "AAC",
				},
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSoundHoundJob 提交一个听歌识曲任务
// https://cloud.tencent.com/document/product/460/84795
func InvokeSoundHoundJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "SoundHound",
		Operation: &cos.MediaProcessJobOperation{
			UserData: "This is my SoundHound job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeNoiseReductionJob 提交一个音频降噪任务
// https://cloud.tencent.com/document/product/460/84796
func InvokeNoiseReductionJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "NoiseReduction",
		Input: &cos.JobInput{
			Object: "input/zhanghuimei_wen.mp3",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "output/out.mp3",
			},
			UserData: "This is my NoiseReduction job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeTtsJob 提交一个语音合成任务
// https://cloud.tencent.com/document/product/460/84797
func InvokeTtsJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "Tts",
		Operation: &cos.MediaProcessJobOperation{
			TtsTpl: &cos.TtsTpl{
				Mode:      "Sync",
				Codec:     "mp3",
				VoiceType: "aixiaonan",
				Volume:    "5",
				Speed:     "150",
			},
			TtsConfig: &cos.TtsConfig{
				Input:     "床前明月光，疑是地上霜",
				InputType: "Text",
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "output/out.mp3",
			},
			UserData: "This is my Tts job",
			JobLevel: 1,
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeSpeechRecognitionJob 提交一个语音识别任务
// https://cloud.tencent.com/document/product/460/84798
func InvokeSpeechRecognitionJob() {
	c := getClient()
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "SpeechRecognition",
		Input: &cos.JobInput{
			Object: "abc.mp3",
		},
		Operation: &cos.MediaProcessJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "music.txt",
				Bucket: "test-1234567890",
			},
			SpeechRecognition: &cos.SpeechRecognition{
				ChannelNum:      "1",
				EngineModelType: "8k_zh",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeMultiJobs 提交多任务处理
// https://cloud.tencent.com/document/product/460/84764
func InvokeMultiJobs() {
	c := getClient()
	createJobOpt := &cos.CreateMultiMediaJobsOptions{
		Input: &cos.JobInput{
			Object: "input/test.mp4",
		},
		Operation: []cos.MediaProcessJobOperation{
			cos.MediaProcessJobOperation{
				Tag: "Snapshot",
				Output: &cos.JobOutput{
					Region: "ap-chongqing",
					Object: "output/go_${Number}.mp4",
					Bucket: "test-1234567890",
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
					Object: "output/go_test.mp4",
					Bucket: "test-1234567890",
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
					Bucket: "test-1234567890",
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
	}
	createJobRes, _, err := c.CI.CreateMultiMediaJobs(context.Background(), createJobOpt)
	log_status(err)
	for k, job := range createJobRes.JobsDetail {
		fmt.Printf("job:%d, %+v\n", k, job)
	}
}

// InvokeFillConcatJob 提交一个填充拼接任务
func InvokeFillConcatJob() {
	c := getClient()
	FillConcat := make([]cos.FillConcatInput, 0)
	FillConcat = append(FillConcat, cos.FillConcatInput{
		Url: "https://test-1234567890.cos.ap-chongqing.myqcloud.com/input/car.mp4",
	})
	FillConcat = append(FillConcat, cos.FillConcatInput{
		FillTime: "5.5",
	})
	FillConcat = append(FillConcat, cos.FillConcatInput{
		Url: "https://test-1234567890.cos.ap-chongqing.myqcloud.com/input/game.mp4",
	})
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "FillConcat",
		Operation: &cos.MediaProcessJobOperation{
			FillConcat: &cos.FillConcat{
				Format:    "mp4",
				FillInput: FillConcat,
			},
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "fill_concat.mp4",
				Bucket: "test-1234567890",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// InvokeVideoSynthesisJob 提交一个视频合成任务
func InvokeVideoSynthesisJob() {
	c := getClient()
	SpliceInfo := make([]cos.VideoSynthesisSpliceInfo, 0)
	SpliceInfo = append(SpliceInfo, cos.VideoSynthesisSpliceInfo{
		Url:   "https://test-1234567890.cos.ap-chongqing.myqcloud.com/input/car.mp4",
		Width: "640",
	})
	SpliceInfo = append(SpliceInfo, cos.VideoSynthesisSpliceInfo{
		Url:   "https://test-1234567890.cos.ap-chongqing.myqcloud.com/input/game.mp4",
		X:     "640",
		Width: "640",
	})
	w := cos.Watermark{
		Type:    "Text",
		LocMode: "Absolute",
		Dx:      "640",
		Pos:     "TopLeft",
		Text: &cos.Text{
			Text:         "helloworld",
			FontSize:     "25",
			FontType:     "simfang.ttf",
			FontColor:    "0xff0000",
			Transparency: "100",
		},
	}
	ws := []cos.Watermark{}
	ws = append(ws, w)
	createJobOpt := &cos.CreateJobsOptions{
		Tag: "VideoSynthesis",
		Operation: &cos.MediaProcessJobOperation{
			VideoSynthesis: &cos.VideoSynthesis{
				KeepAudioTrack: "false",
				SpliceInfo:     SpliceInfo,
			},
			Transcode: &cos.Transcode{
				Container: &cos.Container{
					Format: "mp4",
				},
				Video: &cos.Video{
					Codec:  "H.264",
					Width:  "1280",
					Height: "960",
				},
				Audio: &cos.Audio{
					Codec: "AAC",
				},
			},
			Watermark: ws,
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Object: "video_synthesis.mp4",
				Bucket: "test-1234567890",
			},
		},
	}
	createJobRes, _, err := c.CI.CreateJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

// JobNotifyCallback 解析任务回调
func JobNotifyCallback() {
	taskBody := "<Response><EventName>TaskFinish</EventName><JobsDetail><Code>Success</Code><CreationTime>2022-06-30T19:30:20+0800</CreationTime><EndTime>2022-06-30T19:31:56+0800</EndTime><Input><BucketId>test-123456789</BucketId><Object>input/demo.mp4</Object><Region>ap-chongqing</Region><CosHeaders><Key>Content-Type</Key><Value>video/mp4</Value></CosHeaders><CosHeaders><Key>x-cos-request-id</Key><Value>NjJiZDYwYTFfNjUzYTYyNjRfZjEwZl8xMmZhYzY5</Value></CosHeaders><CosHeaders><Key>EventName</Key><Value>cos:ObjectCreated:Put</Value></CosHeaders><CosHeaders><Key>Size</Key><Value>1424687</Value></CosHeaders></Input><JobId>j06668dc0f86811ecb90d0b03267ce0e5</JobId><Message/><Operation><DigitalWatermark><IgnoreError>false</IgnoreError><Message>123456789ab</Message><State>Failed</State><Type>Text</Type><Version>V1</Version></DigitalWatermark><MediaInfo><Format><Bitrate>8867.172000</Bitrate><Duration>13.654000</Duration><FormatLongName>QuickTime / MOV</FormatLongName><FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName><NumProgram>0</NumProgram><NumStream>2</NumStream><Size>15134046</Size><StartTime>0.000000</StartTime></Format><Stream><Audio><Bitrate>128.726000</Bitrate><Channel>2</Channel><ChannelLayout>stereo</ChannelLayout><CodecLongName>AAC (Advanced Audio Coding)</CodecLongName><CodecName>aac</CodecName><CodecTag>0x6134706d</CodecTag><CodecTagString>mp4a</CodecTagString><CodecTimeBase>1/44100</CodecTimeBase><Duration>13.652993</Duration><Index>1</Index><Language>und</Language><SampleFmt>fltp</SampleFmt><SampleRate>44100</SampleRate><StartTime>0.000000</StartTime><Timebase>1/44100</Timebase></Audio><Subtitle/><Video><AvgFps>25.000000</AvgFps><Bitrate>9197.180000</Bitrate><CodecLongName>H.265 / HEVC (High Efficiency Video Coding)</CodecLongName><CodecName>hevc</CodecName><CodecTag>0x31766568</CodecTag><CodecTagString>hev1</CodecTagString><CodecTimeBase>1/12800</CodecTimeBase><ColorPrimaries>bt470bg</ColorPrimaries><ColorRange>tv</ColorRange><ColorTransfer>smpte170m</ColorTransfer><Duration>12.960000</Duration><FieldOrder>progressive</FieldOrder><Fps>25.000000</Fps><HasBFrame>2</HasBFrame><Height>1920</Height><Index>0</Index><Language>und</Language><Level>120</Level><NumFrames>324</NumFrames><PixFormat>yuv420p</PixFormat><Profile>Main</Profile><RefFrames>1</RefFrames><Rotation>0.000000</Rotation><StartTime>0.000000</StartTime><Timebase>1/12800</Timebase><Width>1088</Width></Video></Stream></MediaInfo><MediaResult><OutputFile><Bucket>test-123456789</Bucket><Md5Info><Md5>852883012a6ba726e6ed8d9b984edfdf</Md5><ObjectName>output/super_resolution.mp4</ObjectName></Md5Info><ObjectName>output/super_resolution.mp4</ObjectName><ObjectPrefix/><Region>ap-chongqing</Region></OutputFile></MediaResult><Output><Bucket>test-123456789</Bucket><Object>output/super_resolution.${ext}</Object><Region>ap-chongqing</Region></Output><TemplateId>t1f1ae1dfsdc9ds41dsb31632d45710642a</TemplateId><TemplateName>template_superresolution</TemplateName><TranscodeTemplateId>t156c107210e7243c5817354565d81b578</TranscodeTemplateId><UserData>This is my SuperResolution job.</UserData><JobLevel>0</JobLevel><WatermarkTemplateId>t143ae6e040af6431aa772c9ec3f0a3f36</WatermarkTemplateId><WatermarkTemplateId>t12a74d11687d444deba8a6cc52051ac27</WatermarkTemplateId></Operation><QueueId>p2242ab62c7c94486915508540933a2c6</QueueId><StartTime>2022-06-30T19:30:21+0800</StartTime><State>Success</State><Progress>100</Progress><SubTag>DigitalWatermark</SubTag><Tag>SuperResolution</Tag><Workflow><Name>SuperResolution_1581665960537</Name><RunId>ic90edd59f84f11ec9d4f525400a3c59f</RunId><WorkflowId>web6ac56c1ef54dbfa44d7f4103203be9</WorkflowId><WorkflowName>workflow-test</WorkflowName></Workflow></JobsDetail></Response>"
	var body cos.JobsNotifyBody
	err := xml.Unmarshal([]byte(taskBody), &body)
	if err != nil {
		fmt.Println(fmt.Sprintf("err:%v", err))
	} else {
		fmt.Println(fmt.Sprintf("body:%+v", body))
		fmt.Println(fmt.Sprintf("mediaInfo:%+v", body.JobsDetail[0].Operation.MediaInfo))
		fmt.Println(fmt.Sprintf("mediaResult:%+v", body.JobsDetail[0].Operation.MediaResult))
	}
}

func main() {
	InvokeTranscodeJob()
}
