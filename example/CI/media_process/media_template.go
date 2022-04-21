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

// DescribeTemplate 搜索模板
func DescribeTemplate() {
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
	opt := &cos.DescribeMediaTemplateOptions{
		Tag:        "Transcode",
		PageNumber: 1,
		PageSize:   5,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// DeleteTemplate 删除模板
func DeleteTemplate() {
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
	DescribeTemplateRes, _, err := c.CI.DeleteMediaTemplate(context.Background(), "t11c1b0a3fb304463096e828a40a013579")
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateTranscodeTemplate TODO
func CreateTranscodeTemplate() {
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
	// CreateMediatemplate
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

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateTranscodeTemplate TODO
func UpdateTranscodeTemplate() {
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
	// UpdateMediatemplate
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
	templateId := "t1d31d58d8a4204d7396087f56a448abd5"
	createTplRes, _, err := c.CI.UpdateMediaTranscodeTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateSnapshotTemplate TODO
func CreateSnapshotTemplate() {
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
	// CreateMediatemplate
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

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateSnapshotTemplate TODO
func UpdateSnapshotTemplate() {
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
	// UpdateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaSnapshotTemplateOptions{
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
	createTplRes, _, err := c.CI.UpdateMediaSnapshotTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateHighSpeedHdTemplate TODO
func CreateHighSpeedHdTemplate() {
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
	// CreateMediatemplate
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

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateHighSpeedHdTemplate TODO
func UpdateHighSpeedHdTemplate() {
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
	// UpdateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaTranscodeTemplateOptions{
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
	createTplRes, _, err := c.CI.UpdateMediaTranscodeTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateAnimationTemplate TODO
func CreateAnimationTemplate() {
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
	// CreateMediatemplate
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

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateAnimationTemplate TODO
func UpdateAnimationTemplate() {
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
	// UpdateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaAnimationTemplateOptions{
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
	createTplRes, _, err := c.CI.UpdateMediaAnimationTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateConcatTemplate TODO
func CreateConcatTemplate() {
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
	// CreateMediatemplate
	rand.Seed(time.Now().UnixNano())
	concatFragmentStart := cos.ConcatFragment{
		Url:  "https://lilang-1253960454.cos.ap-chongqing.myqcloud.com/start.mp4",
		Mode: "Start",
	}
	concatFragmentEnd := cos.ConcatFragment{
		Url:  "https://lilang-1253960454.cos.ap-chongqing.myqcloud.com/end.mp4",
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

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateConcatTemplate TODO
func UpdateConcatTemplate() {
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
	// CreateMediatemplate
	rand.Seed(time.Now().UnixNano())
	concatFragmentEnd := cos.ConcatFragment{
		Url:  "https://lilang-1253960454.cos.ap-chongqing.myqcloud.com/end.mp4",
		Mode: "End",
	}
	var concatFragment []cos.ConcatFragment
	concatFragment = append(concatFragment, concatFragmentEnd)
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
	templateId := "t12a4e410d78fd48e9a999bb682831fc79"
	createTplRes, _, err := c.CI.UpdateMediaConcatTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	// DescribeMediaTemplate
	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

// CreateVideoProcessTemplate TODO
func CreateVideoProcessTemplate() {
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
	// CreateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaVideoProcessTemplateOptions{
		Tag:  "VideoProcess",
		Name: "VideoProcess-" + strconv.Itoa(rand.Intn(100)),
		ColorEnhance: &cos.ColorEnhance{
			Enable:     "true",
			Contrast:   "50",
			Correction: "30",
			Saturation: "20",
		},
		MsSharpen: &cos.MsSharpen{
			Enable:       "true",
			SharpenLevel: "5",
		},
	}
	createTplRes, _, err := c.CI.CreateMediaVideoProcessTemplate(context.Background(), createTplOpt)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	// DescribeMediaTemplate
	if createTplRes.Template != nil {
		opt := &cos.DescribeMediaTemplateOptions{
			Ids: createTplRes.Template.TemplateId,
		}
		DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
		log_status(err)
		fmt.Printf("%+v\n", DescribeTemplateRes)
	}
}

// UpdateVideoProcessTemplate TODO
func UpdateVideoProcessTemplate() {
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
	// CreateMediatemplate
	rand.Seed(time.Now().UnixNano())
	createTplOpt := &cos.CreateMediaVideoProcessTemplateOptions{
		Tag:  "VideoProcess",
		Name: "VideoProcess-" + strconv.Itoa(rand.Intn(100)),
		ColorEnhance: &cos.ColorEnhance{
			Enable:     "true",
			Contrast:   "45",
			Correction: "30",
			Saturation: "20",
		},
		MsSharpen: &cos.MsSharpen{
			Enable:       "true",
			SharpenLevel: "5",
		},
	}
	templateId := "t10af0e373be4d46df9b643a82c779eb10"
	createTplRes, _, err := c.CI.UpdateMediaVideoProcessTemplate(context.Background(), createTplOpt, templateId)
	log_status(err)
	fmt.Printf("%+v\n", createTplRes.Template)

	// DescribeMediaTemplate
	opt := &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	}
	DescribeTemplateRes, _, err := c.CI.DescribeMediaTemplate(context.Background(), opt)
	log_status(err)
	fmt.Printf("%+v\n", DescribeTemplateRes)
}

func main() {
	// DescribeTemplate()
	// DeleteTemplate()
	// CreateTranscodeTemplate()
	// UpdateTranscodeTemplate()
	// CreateSnapshotTemplate()
	// UpdateSnapshotTemplate()
	// CreateHighSpeedHdTemplate()
	// UpdateHighSpeedHdTemplate()
	// CreateAnimationTemplate()
	// UpdateAnimationTemplate()
	// CreateConcatTemplate()
	// UpdateConcatTemplate()
	// CreateVideoProcessTemplate()
	UpdateVideoProcessTemplate()
}
