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

// DescribeTemplate 搜索模板
// https://cloud.tencent.com/document/product/460/84739
func DescribeTemplate() {
	c := getClient()
	opt := &cos.DescribeTemplateOptions{
		Tag:        "Transcode",
		PageNumber: 1,
		PageSize:   5,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// DeleteTemplate 删除模板
// https://cloud.tencent.com/document/product/460/84738
func DeleteTemplate() {
	c := getClient()
	DescribeTemplateRes, _, err := c.CI.DeleteTemplate(context.Background(), "t11c1b0a3fb304463096e828a40a013579")
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateTranscodeTemplate 创建音视频转码模板
// https://cloud.tencent.com/document/product/460/84733
func CreateTranscodeTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
		Tag:  "Transcode",
		Name: "transtpl-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
		},
		Audio: &cos.Audio{
			Codec: "aac",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaTranscodeTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateTranscodeTemplate 更新音视频转码模板
// https://cloud.tencent.com/document/product/460/84754
func UpdateTranscodeTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
		Tag:  "Transcode",
		Name: "transtpl-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
			Crf:   "26",
		},
		Audio: &cos.Audio{
			Codec: "aac",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	templateId := "t139d04d903fee41dd88572cf56b8449fc"
	updateTplRes, _, err := c.CI.UpdateMediaTranscodeTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateVideoEnhanceTemplate 创建画质增强模板
// https://cloud.tencent.com/document/product/460/84722
func CreateVideoEnhanceTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateVideoEnhanceTemplateOptions{
		Tag:  "VideoEnhance",
		Name: "VideoEnhance-" + strconv.Itoa(rand.Intn(100)),
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
	}
	createTplRes, _, err := c.CI.CreateVideoEnhanceTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateVideoEnhanceTemplate 更新画质增强模板
// https://cloud.tencent.com/document/product/460/84745
func UpdateVideoEnhanceTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateVideoEnhanceTemplateOptions{
		Tag:  "VideoEnhance",
		Name: "VideoEnhance-" + strconv.Itoa(rand.Intn(100)),
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
	}
	templateId := "t1e3af4d2467474ebd9e65782b909a7b8b"
	updateTplRes, _, err := c.CI.UpdateVideoEnhanceTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateHighSpeedHdTemplate 创建极速高清转码模板
// https://cloud.tencent.com/document/product/460/84723
func CreateHighSpeedHdTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
		Tag:  "HighSpeedHd",
		Name: "HighSpeedHd-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.265",
			Width: "1280",
			Fps:   "30",
		},
		Audio: &cos.Audio{
			Codec: "aac",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
		TransConfig: &cos.TransConfig{
			IsHdr2Sdr: "true",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaTranscodeTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateHighSpeedHdTemplate 更新极速高清转码模板
// https://cloud.tencent.com/document/product/460/84746
func UpdateHighSpeedHdTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
		Tag:  "HighSpeedHd",
		Name: "HighSpeedHd-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
		},
		Audio: &cos.Audio{
			Codec: "aac",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
		TransConfig: &cos.TransConfig{
			IsHdr2Sdr: "true",
		},
	}
	templateId := "t143d74628378645ed843201ce56b0796a"
	updateTplRes, _, err := c.CI.UpdateMediaTranscodeTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateVideoMontageTemplate 创建精彩集锦模板
// https://cloud.tencent.com/document/product/460/84724
func CreateVideoMontageTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaVideoMontageTemplateOptions{
		Tag:      "VideoMontage",
		Name:     "VideoMontage-" + strconv.Itoa(rand.Intn(100)),
		Duration: "120",
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
		},
		Audio: &cos.Audio{
			Codec:    "aac",
			Channels: "4",
		},
		AudioMix: &cos.AudioMix{
			AudioSource: "https://test-123456789.cos.ap-chongqing.myqcloud.com/src.mp3",
			MixMode:     "Repeat",
			Replace:     "true",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaVideoMontageTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateVideoMontageTemplate 更新精彩集锦模板
// https://cloud.tencent.com/document/product/460/84747
func UpdateVideoMontageTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaVideoMontageTemplateOptions{
		Tag:  "VideoMontage",
		Name: "VideoMontage-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mp4",
		},
		Video: &cos.Video{
			Codec: "h.264",
			Width: "1280",
			Fps:   "30",
		},
		Audio: &cos.Audio{
			Codec:    "aac",
			Channels: "4",
		},
		AudioMix: &cos.AudioMix{
			AudioSource: "https://test-123456789.cos.ap-chongqing.myqcloud.com/src.mp3",
			MixMode:     "Once",
			Replace:     "true",
		},
	}
	templateId := "t188cb6223ca48420f9cd15ca9855e8a9b"
	updateTplRes, _, err := c.CI.UpdateMediaVideoMontageTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateWatermarkTemplate 创建明水印模板
// https://cloud.tencent.com/document/product/460/84725
func CreateWatermarkTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaWatermarkTemplateOptions{
		Tag:  "Watermark",
		Name: "Watermark-" + strconv.Itoa(rand.Intn(100)),
		Watermark: &cos.Watermark{
			Type:      "Text",
			LocMode:   "Absolute",
			Dx:        "20",
			Dy:        "20",
			Pos:       "TopRight",
			StartTime: "5",
			EndTime:   "20",
			Text: &cos.Text{
				Text:         "tencent",
				FontSize:     "12",
				FontType:     "simfang.ttf",
				FontColor:    "0xff0000",
				Transparency: "100",
			},
		},
	}
	createTplRes, _, err := c.CI.CreateMediaWatermarkTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateWatermarkTemplate 更新明水印模板
// https://cloud.tencent.com/document/product/460/84748
func UpdateWatermarkTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaWatermarkTemplateOptions{
		Tag:  "Watermark",
		Name: "Watermark-" + strconv.Itoa(rand.Intn(100)),
		Watermark: &cos.Watermark{
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
		},
	}
	templateId := "t1740baca715ad4ec2b5fbc02c76987025"
	updateTplRes, _, err := c.CI.UpdateMediaWatermarkTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateSnapshotTemplate 创建视频截图模板
// https://cloud.tencent.com/document/product/460/84727
func CreateSnapshotTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaSnapshotTemplateOptions{
		Tag:  "Snapshot",
		Name: "Snapshot-" + strconv.Itoa(rand.Intn(100)),
		Snapshot: &cos.Snapshot{
			Width:           "1280",
			Height:          "960",
			Start:           "0",
			TimeInterval:    "5",
			Count:           "10",
			SnapshotOutMode: "SnapshotAndSprite",
			SpriteSnapshotConfig: &cos.SpriteSnapshotConfig{
				Color:   "AliceBlue",
				Columns: "3",
				Lines:   "3",
			},
		},
	}
	createTplRes, _, err := c.CI.CreateMediaSnapshotTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateSnapshotTemplate 更新视频截图模板
// https://cloud.tencent.com/document/product/460/84749
func UpdateSnapshotTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaSnapshotTemplateOptions{
		Tag:  "Snapshot",
		Name: "Snapshot-" + strconv.Itoa(rand.Intn(100)),
		Snapshot: &cos.Snapshot{
			Width:           "720",
			Height:          "480",
			Start:           "0",
			TimeInterval:    "5",
			Count:           "10",
			SnapshotOutMode: "SnapshotAndSprite",
			SpriteSnapshotConfig: &cos.SpriteSnapshotConfig{
				Color:   "AliceBlue",
				Columns: "3",
				Lines:   "3",
			},
		},
	}
	templateId := "t1bc84403414784c9d969037b96cef9cf9"
	updateTplRes, _, err := c.CI.UpdateMediaSnapshotTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateAnimationTemplate 创建视频转动图模板
// https://cloud.tencent.com/document/product/460/84729
func CreateAnimationTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaAnimationTemplateOptions{
		Tag:  "Animation",
		Name: "Animation-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "hgif",
		},
		Video: &cos.AnimationVideo{
			Codec: "gif",
			Width: "1280",
			Fps:   "30",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaAnimationTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateAnimationTemplate 更新视频转动图模板
// https://cloud.tencent.com/document/product/460/84751
func UpdateAnimationTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaAnimationTemplateOptions{
		Tag:  "Animation",
		Name: "Animation-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "gif",
		},
		Video: &cos.AnimationVideo{
			Codec: "gif",
			Width: "1280",
			Fps:   "50",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	templateId := "t10a23d5cf28ee453eb7982d4709587ecf"
	updateTplRes, _, err := c.CI.UpdateMediaAnimationTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateConcatTemplate 创建音视频拼接模板
// https://cloud.tencent.com/document/product/460/84730
func CreateConcatTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	concatFragmentStart := cos.ConcatFragment{
		Url:  "https://test-123456789.cos.ap-chongqing.myqcloud.com/start.mp4",
		Mode: "Start",
	}
	concatFragmentEnd := cos.ConcatFragment{
		Url:  "https://test-123456789.cos.ap-chongqing.myqcloud.com/end.mp4",
		Mode: "End",
	}
	var concatFragment []cos.ConcatFragment
	concatFragment = append(concatFragment, concatFragmentStart, concatFragmentEnd)
	createTplOpt := &cos.CreateMediaConcatTemplateOptions{
		Tag:  "Concat",
		Name: "Concat-" + strconv.Itoa(rand.Intn(100)),
		ConcatTemplate: &cos.ConcatTemplate{
			Container: &cos.Container{
				Format: "mp4",
			},
			Video: &cos.Video{
				Codec: "h.264",
				Width: "1280",
				Fps:   "30",
			},
			Audio: &cos.Audio{
				Codec:    "aac",
				Channels: "4",
			},
			ConcatFragment: concatFragment,
		},
	}
	createTplRes, _, err := c.CI.CreateMediaConcatTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateConcatTemplate 更新音视频拼接模板
// https://cloud.tencent.com/document/product/460/84752
func UpdateConcatTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	concatFragmentEnd := cos.ConcatFragment{
		Url:  "https://test-123456789.cos.ap-chongqing.myqcloud.com/end.mp4",
		Mode: "End",
	}
	var concatFragment []cos.ConcatFragment
	concatFragment = append(concatFragment, concatFragmentEnd)
	updateTplOpt := &cos.CreateMediaConcatTemplateOptions{
		Tag:  "Concat",
		Name: "Concat-" + strconv.Itoa(rand.Intn(100)),
		ConcatTemplate: &cos.ConcatTemplate{
			Container: &cos.Container{
				Format: "mp4",
			},
			Video: &cos.Video{
				Codec: "h.264",
				Width: "1280",
				Fps:   "30",
			},
			Audio: &cos.Audio{
				Codec:    "aac",
				Channels: "4",
			},
			ConcatFragment: concatFragment,
		},
	}
	templateId := "t12a4e410d78fd48e9a999bb682831fc79"
	updateTplRes, _, err := c.CI.UpdateMediaConcatTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateTranscodeProTemplate 创建音视频转码 pro 模板
// https://cloud.tencent.com/document/product/460/84732
func CreateTranscodeProTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTranscodeProTemplateOptions{
		Tag:  "TranscodePro",
		Name: "TranscodePro-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mxf",
		},
		Video: &cos.TranscodeProVideo{
			Codec:      "xavc",
			Profile:    "XAVC-HD_intra_420_10bit_class50",
			Width:      "1440",
			Height:     "1080",
			Interlaced: "false",
			Fps:        "30000/1001",
		},
		Audio: &cos.TranscodeProAudio{
			Codec: "pcm_s24le",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaTranscodeProTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateWatermarkTemplate 更新音视频转码 pro 模板
// https://cloud.tencent.com/document/product/460/84753
func UpdateTranscodeProTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaTranscodeProTemplateOptions{
		Tag:  "TranscodePro",
		Name: "TranscodePro-" + strconv.Itoa(rand.Intn(100)),
		Container: &cos.Container{
			Format: "mxf",
		},
		Video: &cos.TranscodeProVideo{
			Codec:      "xavc",
			Profile:    "XAVC-HD_intra_420_10bit_class50",
			Width:      "1440",
			Height:     "1080",
			Interlaced: "false",
			Fps:        "24000/1001",
		},
		Audio: &cos.TranscodeProAudio{
			Codec: "pcm_s24le",
		},
		TimeInterval: &cos.TimeInterval{
			Start:    "0",
			Duration: "",
		},
	}
	templateId := "t11837976491864248885b037453466e49"
	updateTplRes, _, err := c.CI.UpdateMediaTranscodeProTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateSmartCoverTemplate 创建智能封面模板
// https://cloud.tencent.com/document/product/460/84734
func CreateSmartCoverTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaSmartCoverTemplateOptions{
		Tag:  "SmartCover",
		Name: "SmartCover-" + strconv.Itoa(rand.Intn(100)),
		SmartCover: &cos.NodeSmartCover{
			Format:           "jpg",
			Width:            "1280",
			Height:           "960",
			Count:            "10",
			DeleteDuplicates: "true",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaSmartCoverTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateSmartCoverTemplate 更新智能封面模板
// https://cloud.tencent.com/document/product/460/84755
func UpdateSmartCoverTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaSmartCoverTemplateOptions{
		Tag:  "SmartCover",
		Name: "SmartCover-" + strconv.Itoa(rand.Intn(100)),
		SmartCover: &cos.NodeSmartCover{
			Format:           "jpg",
			Width:            "1280",
			Height:           "960",
			Count:            "5",
			DeleteDuplicates: "true",
		},
	}
	templateId := "t17fcea6bf45f44fa1a76f3b11b1f4523f"
	updateTplRes, _, err := c.CI.UpdateMediaSmartCoverTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreatePicProcessTemplate 创建图片处理模板
// https://cloud.tencent.com/document/product/460/84735
func CreatePicProcessTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaPicProcessTemplateOptions{
		Tag:  "PicProcess",
		Name: "PicProcess-" + strconv.Itoa(rand.Intn(100)),
		PicProcess: &cos.PicProcess{
			IsPicInfo:   "true",
			ProcessRule: "imageMogr2/thumbnail/!50p",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaPicProcessTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdatePicProcessTemplate 更新图片处理模板
// https://cloud.tencent.com/document/product/460/84756
func UpdatePicProcessTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaPicProcessTemplateOptions{
		Tag:  "PicProcess",
		Name: "PicProcess-" + strconv.Itoa(rand.Intn(100)),
		PicProcess: &cos.PicProcess{
			IsPicInfo:   "true",
			ProcessRule: "imageMogr2/thumbnail/!55p",
		},
	}
	templateId := "t12db15e06bf504228951b2fa62b1b7b90"
	updateTplRes, _, err := c.CI.UpdateMediaPicProcessTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateVideoTargetRecTemplate 创建视频目标检测模板
// https://cloud.tencent.com/document/product/460/84736
func CreateVideoTargetRecTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateVideoTargetRecTemplateOptions{
		Tag:  "VideoTargetRec",
		Name: "VideoTargetRec-" + strconv.Itoa(rand.Intn(100)),
		VideoTargetRec: &cos.VideoTargetRec{
			Body: "true",
			Pet:  "true",
			Car:  "true",
		},
	}
	createTplRes, _, err := c.CI.CreateVideoTargetRecTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateVideoTargetRecTemplate 更新视频目标检测模板
// https://cloud.tencent.com/document/product/460/84756
func UpdateVideoTargetRecTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateVideoTargetRecTemplateOptions{
		Tag:  "VideoTargetRec",
		Name: "VideoTargetRec-" + strconv.Itoa(rand.Intn(100)),
		VideoTargetRec: &cos.VideoTargetRec{
			Body: "true",
			Pet:  "true",
			Car:  "true",
		},
	}
	templateId := "t10d7cdebcea61426e9b7bd701fb2f2fdc"
	updateTplRes, _, err := c.CI.UpdateVideoTargetRecTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateSpeechRecognitionTemplate 创建语音识别模板
// https://cloud.tencent.com/document/product/460/84498
func CreateSpeechRecognitionTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaSpeechRecognitionTemplateOptions{
		Tag:  "SpeechRecognition",
		Name: "SpeechRecognition-" + strconv.Itoa(rand.Intn(100)),
		SpeechRecognition: &cos.SpeechRecognition{
			ChannelNum:         "1",
			EngineModelType:    "16k_zh",
			ResTextFormat:      "1",
			FilterDirty:        "0",
			FilterModal:        "1",
			ConvertNumMode:     "0",
			SpeakerDiarization: "1",
			SpeakerNumber:      "0",
			FilterPunc:         "0",
			OutputFileType:     "txt",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaSpeechRecognitionTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateSpeechRecognitionTemplate 更新语音识别模板
// https://cloud.tencent.com/document/product/460/84759
func UpdateSpeechRecognitionTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaSpeechRecognitionTemplateOptions{
		Tag:  "SpeechRecognition",
		Name: "SpeechRecognition-" + strconv.Itoa(rand.Intn(100)),
		SpeechRecognition: &cos.SpeechRecognition{
			EngineModelType:    "16k_zh",
			ResTextFormat:      "1",
			FilterDirty:        "0",
			FilterModal:        "1",
			ConvertNumMode:     "0",
			SpeakerDiarization: "1",
			SpeakerNumber:      "0",
			FilterPunc:         "0",
			OutputFileType:     "txt",
		},
	}
	templateId := "t1a883a072103f440fa7b9b54b744a2fbf"
	updateTplRes, _, err := c.CI.UpdateMediaSpeechRecognitionTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateTtsTemplate 创建语音合成模板
// https://cloud.tencent.com/document/product/460/84499
func CreateTtsTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTtsTemplateOptions{
		Tag:       "Tts",
		Name:      "Tts-" + strconv.Itoa(rand.Intn(100)),
		Mode:      "Sync",
		Codec:     "mp3",
		VoiceType: "aixiaoxing",
		Volume:    "5",
		Speed:     "150",
		Emotion:   "arousal",
	}
	createTplRes, _, err := c.CI.CreateMediaTtsTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateTtsTemplate 更新语音合成模板
// https://cloud.tencent.com/document/product/460/84758
func UpdateTtsTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaTtsTemplateOptions{
		Tag:       "Tts",
		Name:      "Tts-" + strconv.Itoa(rand.Intn(100)),
		Mode:      "Sync",
		Codec:     "mp3",
		VoiceType: "aixiaonan",
		Volume:    "5",
		Speed:     "150",
	}
	templateId := "t174f96537bae547c785042ecdbb228e6e"
	updateTplRes, _, err := c.CI.UpdateMediaTtsTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateVoiceSeparateTemplate 创建人声分离模板
// https://cloud.tencent.com/document/product/460/84500
func CreateVoiceSeparateTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaVoiceSeparateTemplateOptions{
		Tag:       "VoiceSeparate",
		Name:      "VoiceSeparate-" + strconv.Itoa(rand.Intn(100)),
		AudioMode: "IsAudio",
		AudioConfig: &cos.AudioConfig{
			Codec:      "aac",
			Samplerate: "32000",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaVoiceSeparateTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateVoiceSeparateTemplate 更新人声分离模板
// https://cloud.tencent.com/document/product/460/84757
func UpdateVoiceSeparateTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateMediaVoiceSeparateTemplateOptions{
		Tag:       "VoiceSeparate",
		Name:      "VoiceSeparate-" + strconv.Itoa(rand.Intn(100)),
		AudioMode: "IsAudio",
		AudioConfig: &cos.AudioConfig{
			Codec:      "mp3",
			Samplerate: "32000",
		},
	}
	templateId := "t169f9ad6166e24695a7de413c646f9e77"
	updateTplRes, _, err := c.CI.UpdateMediaVoiceSeparateTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

// CreateNoiseReductionTemplate 创建降噪模板
func CreateNoiseReductionTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateNoiseReductionTemplateOptions{
		Tag:  "NoiseReduction",
		Name: "NoiseReduction-" + strconv.Itoa(rand.Intn(100)),
		NoiseReduction: &cos.NoiseReduction{
			Format:     "wav",
			Samplerate: "16000",
		},
	}
	createTplRes, _, err := c.CI.CreateNoiseReductionTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)
}

// UpdateNoiseReductionTemplate 更新降噪模板
func UpdateNoiseReductionTemplate() {
	c := getClient()
	rand.Seed(time.Now().UnixNano())
	updateTplOpt := &cos.CreateNoiseReductionTemplateOptions{
		Tag:  "NoiseReduction",
		Name: "NoiseReduction-" + strconv.Itoa(rand.Intn(100)),
		NoiseReduction: &cos.NoiseReduction{
			Format:     "mp3",
			Samplerate: "16000",
		},
	}
	templateId := "t178bbee7296b3412db24a274980d5eb1a"
	updateTplRes, _, err := c.CI.UpdateNoiseReductionTemplate(context.Background(), updateTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", updateTplRes.Template)
}

func main() {
}
