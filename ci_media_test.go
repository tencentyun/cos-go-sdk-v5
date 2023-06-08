package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
)

func TestCIService_CreateMultiMediaJobs(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Input><Object>test.mp4</Object></Input>" +
		"<Operation><Tag>Animation</Tag><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test.gif</Object></Output>" +
		"<TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId>" +
		"</Operation><Operation><Tag>Transcode</Tag><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test.mp4</Object></Output>" +
		"<TemplateId>t1995d523e42df4c5e858f244b4174360c</TemplateId>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreateMultiMediaJobsOptions{
		Input: &JobInput{
			Object: "test.mp4",
		},
		Operation: []MediaProcessJobOperation{
			{
				Tag: "Animation",
				Output: &JobOutput{
					Region: "ap-beijing",
					Bucket: "abc-1250000000",
					Object: "test.gif",
				},
				TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
			},
			{
				Tag: "Transcode",
				Output: &JobOutput{
					Region: "ap-beijing",
					Bucket: "abc-1250000000",
					Object: "test.mp4",
				},
				TemplateId: "t1995d523e42df4c5e858f244b4174360c",
			},
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreateMultiMediaJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMultiMediaJobs returned errors: %v", err)
	}
}

func TestCIService_CreateMediaJobs(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Animation</Tag><Input><Object>test.mp4</Object></Input>" +
		"<Operation><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test-trans.gif</Object></Output>" +
		"<TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaJobsOptions{
		Tag: "Animation",
		Input: &JobInput{
			Object: "test.mp4",
		},
		Operation: &MediaProcessJobOperation{
			Output: &JobOutput{
				Region: "ap-beijing",
				Bucket: "abc-1250000000",
				Object: "test-trans.gif",
			},
			TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreateMediaJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaJobs returned errors: %v", err)
	}
}

func TestCIService_CreatePicProcessJobs(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>PicProcess</Tag><Input><Object>input/deer.jpg</Object></Input>" +
		"<Operation><TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test-trans.gif</Object></Output>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/pic_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreatePicJobsOptions{
		Tag: "PicProcess",
		Input: &JobInput{
			Object: "input/deer.jpg",
		},
		Operation: &PicProcessJobOperation{
			Output: &JobOutput{
				Region: "ap-beijing",
				Bucket: "abc-1250000000",
				Object: "test-trans.gif",
			},
			TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreatePicProcessJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreatePicProcessJobs returned errors: %v", err)
	}
}

func TestCIService_CreateAIJobs(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Translation</Tag><Input><Object>input/en.pdf</Object><Lang>en</Lang><Type>pdf</Type><BasicType>pptx</BasicType></Input>" +
		"<Operation><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>output/zh.pdf</Object></Output><Translation><Lang>zh</Lang><Type>pdf</Type></Translation>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/ai_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreateAIJobsOptions{
		Tag: "Translation",
		Input: &JobInput{
			Object:    "input/en.pdf",
			Lang:      "en",
			Type:      "pdf",
			BasicType: "pptx",
		},
		Operation: &MediaProcessJobOperation{
			Output: &JobOutput{
				Region: "ap-beijing",
				Bucket: "abc-1250000000",
				Object: "output/zh.pdf",
			},
			Translation: &Translation{
				Lang: "zh",
				Type: "pdf",
			},
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreateAIJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateASRJobs returned errors: %v", err)
	}
}

func TestCIService_DescribeMediaJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "jabcsdssfeipplsdfwe"
	mux.HandleFunc("/jobs/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeMediaJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.DescribeMediaJob returned error: %v", err)
	}
}

func TestCIService_DescribePicProcessJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "ac7c990a00bf211ed946af9e0691f2b7a"
	mux.HandleFunc("/pic_jobs/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribePicProcessJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.DescribePicProcessJob returned error: %v", err)
	}
}

func TestCIService_DescribeAIJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "ac7c990a00bf211ed946af9e0691f2b7a"
	mux.HandleFunc("/ai_jobs/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeAIJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.DescribeAIJob returned error: %v", err)
	}
}

func TestCIService_DescribeMultiMediaJob(t *testing.T) {
	{
		setup()
		jobIDs := []string{}
		mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
		})

		_, _, err := client.CI.DescribeMultiMediaJob(context.Background(), jobIDs)
		if err == nil || err.Error() != "empty param jobids" {
			t.Fatalf("CI.DescribeMultiMediaJob returned error: %v", err)
		}
		teardown()
	}
	{
		setup()
		jobIDs := []string{"jc7c990a00bf211ed946af9e0691f2b7a", "jabcsdssfeipplsdfwe"}
		mux.HandleFunc("/jobs/"+"jc7c990a00bf211ed946af9e0691f2b7a,jabcsdssfeipplsdfwe", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
		})

		_, _, err := client.CI.DescribeMultiMediaJob(context.Background(), jobIDs)
		if err != nil {
			t.Fatalf("CI.DescribeMultiMediaJob returned error: %v", err)
		}
		teardown()
	}
}

func TestCIService_DescribeMediaJobs(t *testing.T) {
	setup()
	defer teardown()

	queueId := "aaaaaaaaaaa"
	tag := "Animation"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueId": queueId,
			"tag":     tag,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaJobsOptions{
		QueueId: queueId,
		Tag:     tag,
	}

	_, _, err := client.CI.DescribeMediaJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaJobs returned error: %v", err)
	}
}

func TestCIService_DescribeMediaProcessQueues(t *testing.T) {
	setup()
	defer teardown()

	queueIds := "A,B,C"
	mux.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueIds": queueIds,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaProcessQueuesOptions{
		QueueIds: queueIds,
	}

	_, _, err := client.CI.DescribeMediaProcessQueues(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaProcessQueues returned error: %v", err)
	}
}

func TestCIService_DescribePicProcessQueues(t *testing.T) {
	setup()
	defer teardown()

	queueIds := "A,B,C"
	mux.HandleFunc("/picqueue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueIds": queueIds,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribePicProcessQueuesOptions{
		QueueIds: queueIds,
	}

	_, _, err := client.CI.DescribePicProcessQueues(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribePicProcessQueues returned error: %v", err)
	}
}

func TestCIService_DescribeAIProcessQueues(t *testing.T) {
	setup()
	defer teardown()

	queueIds := "A,B,C"
	mux.HandleFunc("/ai_queue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueIds": queueIds,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaProcessQueuesOptions{
		QueueIds: queueIds,
	}

	_, _, err := client.CI.DescribeAIProcessQueues(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeAIProcessQueues returned error: %v", err)
	}
}

func TestCIService_DescribeASRProcessQueues(t *testing.T) {
	setup()
	defer teardown()

	queueIds := "A,B,C"
	mux.HandleFunc("/asrqueue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueIds": queueIds,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaProcessQueuesOptions{
		QueueIds: queueIds,
	}

	_, _, err := client.CI.DescribeASRProcessQueues(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeASRProcessQueues returned error: %v", err)
	}
}

func TestCIService_UpdateMediaProcessQueue(t *testing.T) {
	setup()
	defer teardown()

	queueID := "p8eb46b8cc1a94bc09512d16c5c4f4d3a"
	wantBody := "<Request><Name>QueueName</Name><QueueID>" + queueID + "</QueueID><State>Active</State>" +
		"<NotifyConfig><Url>test.com</Url><State>On</State><Type>Url</Type><Event>TransCodingFinish</Event>" +
		"</NotifyConfig></Request>"
	mux.HandleFunc("/queue/"+queueID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &UpdateMediaProcessQueueOptions{
		Name:    "QueueName",
		QueueID: queueID,
		State:   "Active",
		NotifyConfig: &MediaProcessQueueNotifyConfig{
			Url:   "test.com",
			State: "On",
			Type:  "Url",
			Event: "TransCodingFinish",
		},
	}

	_, _, err := client.CI.UpdateMediaProcessQueue(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.UpdateMediaProcessQueue returned error: %v", err)
	}
}

func TestCIService_DescribeMediaProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	regions := "ap-shanghai,ap-gaungzhou"
	bucketName := "testbucket-1250000000"
	mux.HandleFunc("/mediabucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"regions":    regions,
			"bucketName": bucketName,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaProcessBucketsOptions{
		Regions:    regions,
		BucketName: bucketName,
	}

	_, _, err := client.CI.DescribeMediaProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaProcessBuckets returned error: %v", err)
	}
}

func TestCIService_DescribePicProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	regions := "ap-shanghai,ap-gaungzhou"
	bucketName := "testbucket-1250000000"
	mux.HandleFunc("/picbucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"regions":    regions,
			"bucketName": bucketName,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribePicProcessBucketsOptions{
		Regions:    regions,
		BucketName: bucketName,
	}

	_, _, err := client.CI.DescribePicProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribePicProcessBuckets returned error: %v", err)
	}
}

func TestCIService_DescribeAIProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	regions := "ap-shanghai,ap-gaungzhou"
	bucketName := "testbucket-1250000000"
	mux.HandleFunc("/ai_bucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"regions":    regions,
			"bucketName": bucketName,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeAIProcessBucketsOptions{
		Regions:    regions,
		BucketName: bucketName,
	}

	_, _, err := client.CI.DescribeAIProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeAIProcessBuckets returned error: %v", err)
	}
}

func TestCIService_DescribeASRProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	regions := "ap-shanghai,ap-gaungzhou"
	bucketName := "testbucket-1250000000"
	mux.HandleFunc("/asrbucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"regions":    regions,
			"bucketName": bucketName,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeASRProcessBucketsOptions{
		Regions:    regions,
		BucketName: bucketName,
	}

	_, _, err := client.CI.DescribeASRProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeASRProcessBuckets returned error: %v", err)
	}
}

func TestCIService_DescribeFileProcessBuckets(t *testing.T) {
	setup()
	defer teardown()

	regions := "ap-shanghai,ap-gaungzhou"
	bucketName := "testbucket-1250000000"
	mux.HandleFunc("/file_bucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"regions":    regions,
			"bucketName": bucketName,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeFileProcessBucketsOptions{
		Regions:    regions,
		BucketName: bucketName,
	}

	_, _, err := client.CI.DescribeFileProcessBuckets(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeFileProcessBuckets returned error: %v", err)
	}
}

func TestCIService_GetMediaInfo(t *testing.T) {
	res_xml := `<Response>
	<MediaInfo>
		<Format>
			<Bitrate>1014.950000</Bitrate>
			<Duration>10.125000</Duration>
			<FormatLongName>QuickTime / MOV</FormatLongName>
			<FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName>
			<NumProgram>0</NumProgram>
			<NumStream>2</NumStream>
			<Size>1284547</Size>
			<StartTime>0.000000</StartTime>
		</Format>
		<Stream>
			<Audio>
				<Bitrate>70.451000</Bitrate>
				<Channel>1</Channel>
				<ChannelLayout>mono</ChannelLayout>
				<CodecLongName>AAC (Advanced Audio Coding)</CodecLongName>
				<CodecName>aac</CodecName>
				<CodecTag>0x6134706d</CodecTag>
				<CodecTagString>mp4a</CodecTagString>
				<CodecTimeBase>1/44100</CodecTimeBase>
				<Duration>0.440294</Duration>
				<Index>1</Index>
				<Language>und</Language>
				<SampleFmt>fltp</SampleFmt>
				<SampleRate>44100</SampleRate>
				<StartTime>0.000000</StartTime>
				<Timebase>1/44100</Timebase>
			</Audio>
			<Subtitle/>
			<Video>
				<AvgFps>24/1</AvgFps>
				<Bitrate>938.164000</Bitrate>
				<CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName>
				<CodecName>h264</CodecName>
				<CodecTag>0x31637661</CodecTag>
				<CodecTagString>avc1</CodecTagString>
				<CodecTimeBase>1/12288</CodecTimeBase>
				<Dar>40:53</Dar>
				<Duration>0.124416</Duration>
				<Fps>24.500000</Fps>
				<HasBFrame>2</HasBFrame>
				<Height>1280</Height>
				<Index>0</Index>
				<Language>und</Language>
				<Level>32</Level>
				<NumFrames>243</NumFrames>
				<PixFormat>yuv420p</PixFormat>
				<Profile>High</Profile>
				<RefFrames>1</RefFrames>
				<Sar>25600:25599</Sar>
				<StartTime>0.000000</StartTime>
				<Timebase>1/12288</Timebase>
				<Width>966</Width>
			</Video>
		</Stream>
	</MediaInfo>
</Response>`
	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "videoinfo",
			}
			testFormValues(t, r, v)
			fmt.Fprint(w, res_xml)
		})

		res, _, err := client.CI.GetMediaInfo(context.Background(), "test", nil)
		if err != nil {
			t.Fatalf("CI.GetMediaInfo returned error: %v", err)
		}
		want := &GetMediaInfoResult{}
		err = xml.Unmarshal([]byte(res_xml), want)
		if err != nil {
			t.Errorf("Bucket.GetMediaInfo Unmarshal returned error %+v", err)
		}
		if !reflect.DeepEqual(res, want) {
			t.Errorf("Bucket.GetMediaInfo returned %+v, want %+v", res, want)
		}
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "videoinfo",
				"versionId":  "1",
			}
			testFormValues(t, r, v)
			fmt.Fprint(w, res_xml)
		})

		res, _, err := client.CI.GetMediaInfo(context.Background(), "test", nil, "1")
		if err != nil {
			t.Fatalf("CI.GetMediaInfo returned error: %v", err)
		}
		want := &GetMediaInfoResult{}
		err = xml.Unmarshal([]byte(res_xml), want)
		if err != nil {
			t.Errorf("Bucket.GetMediaInfo Unmarshal returned error %+v", err)
		}
		if !reflect.DeepEqual(res, want) {
			t.Errorf("Bucket.GetMediaInfo returned %+v, want %+v", res, want)
		}
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "videoinfo",
			}
			testFormValues(t, r, v)
			fmt.Fprint(w, res_xml)
		})

		_, _, err := client.CI.GetMediaInfo(context.Background(), "test", nil, "1", "2")
		if err == nil || err.Error() != "wrong params" {
			t.Fatalf("CI.GetMediaInfo returned error: %v", err)
		}
		teardown()
	}
}

func TestCIService_GenerateMediaInfo(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Input><Object>input/test.mp4</Object></Input></Request>"

	mux.HandleFunc("/mediainfo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})
	opt := &GenerateMediaInfoOptions{
		Input: &JobInput{
			Object: "input/test.mp4",
		},
	}
	_, _, err := client.CI.GenerateMediaInfo(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.GenerateMediaInfo returned error: %v", err)
	}
}

func TestCIService_GetSnapshot(t *testing.T) {
	opt := &GetSnapshotOptions{
		Time:   1,
		Height: 100,
		Width:  100,
		Format: "jpg",
		Rotate: "auto",
		Mode:   "exactframe",
	}
	data := make([]byte, 1234*2)
	rand.Read(data)
	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "snapshot",
				"time":       "1",
				"format":     "jpg",
				"width":      "100",
				"height":     "100",
				"rotate":     "auto",
				"mode":       "exactframe",
			}
			testFormValues(t, r, v)
			w.Write(data)
		})

		resp, err := client.CI.GetSnapshot(context.Background(), "test", opt)
		if err != nil {
			t.Fatalf("CI.GetSnapshot returned error: %v", err)
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("CI.GetSnapshot ReadAll returned error: %v", err)
		}
		if bytes.Compare(bs, data) != 0 {
			t.Errorf("Bucket.GetSnapshot Compare failed")
		}
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "snapshot",
				"versionId":  "1.1",
				"time":       "1",
				"format":     "jpg",
				"width":      "100",
				"height":     "100",
				"rotate":     "auto",
				"mode":       "exactframe",
			}
			testFormValues(t, r, v)
			w.Write(data)
		})

		resp, err := client.CI.GetSnapshot(context.Background(), "test", opt, "1.1")
		if err != nil {
			t.Fatalf("CI.GetSnapshot returned error: %v", err)
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("CI.GetSnapshot ReadAll returned error: %v", err)
		}
		if bytes.Compare(bs, data) != 0 {
			t.Errorf("Bucket.GetSnapshot Compare failed")
		}
		teardown()
	}

	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "snapshot",
				"versionId":  "1.1",
				"time":       "1",
				"format":     "jpg",
				"width":      "100",
				"height":     "100",
				"rotate":     "auto",
				"mode":       "exactframe",
			}
			testFormValues(t, r, v)
			w.Write(data)
		})

		_, err := client.CI.GetSnapshot(context.Background(), "test", opt, "1.1", "1.2")
		if err == nil || err.Error() != "wrong params" {
			t.Fatalf("CI.GetMediaInfo returned error: %v", err)
		}
		teardown()
	}
}

func TestCIService_PostSnapshot(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Input><Object>test.mp4</Object></Input>" +
		"<Time>1.00</Time><Width>1024</Width><Height>768</Height><Mode>keyframe</Mode>" +
		"<Rotate>auto</Rotate><Format>jpg</Format>" +
		"<Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test.jpg</Object></Output></Request>"

	mux.HandleFunc("/snapshot", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})
	opt := &PostSnapshotOptions{
		Input: &JobInput{
			Object: "test.mp4",
		},
		Time:   "1.00",
		Width:  1024,
		Height: 768,
		Mode:   "keyframe",
		Rotate: "auto",
		Format: "jpg",
		Output: &JobOutput{
			Region: "ap-beijing",
			Bucket: "abc-1250000000",
			Object: "test.jpg",
		},
	}
	_, _, err := client.CI.PostSnapshot(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.PostSnapshot returned error: %v", err)
	}
}

func TestCIService_GetPrivateM3U8(t *testing.T) {
	opt := &GetPrivateM3U8Options{
		Expires: 3600,
	}
	data := make([]byte, 1234*2)
	rand.Read(data)

	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "pm3u8",
				"expires":    "3600",
			}
			testFormValues(t, r, v)
			w.Write(data)
		})

		resp, err := client.CI.GetPrivateM3U8(context.Background(), "test", opt)
		if err != nil {
			t.Fatalf("CI.GetPrivateM3U8 returned error: %v", err)
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("CI.GetSnapshot ReadAll returned error: %v", err)
		}
		if bytes.Compare(bs, data) != 0 {
			t.Errorf("Bucket.GetSnapshot Compare failed")
		}
		teardown()
	}
	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "pm3u8",
				"expires":    "3600",
				"versionId":  "1.1",
			}
			testFormValues(t, r, v)
			w.Write(data)
		})

		resp, err := client.CI.GetPrivateM3U8(context.Background(), "test", opt, "1.1")
		if err != nil {
			t.Fatalf("CI.GetPrivateM3U8 returned error: %v", err)
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("CI.GetSnapshot ReadAll returned error: %v", err)
		}
		if bytes.Compare(bs, data) != 0 {
			t.Errorf("Bucket.GetSnapshot Compare failed")
		}
		teardown()
	}
	{
		setup()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "pm3u8",
				"expires":    "3600",
			}
			testFormValues(t, r, v)
			w.Write(data)
		})

		_, err := client.CI.GetPrivateM3U8(context.Background(), "test", opt, "1.1", "1.2")
		if err == nil || err.Error() != "wrong params" {
			t.Fatalf("CI.GetMediaInfo returned error: %v", err)
		}
		teardown()
	}
}

func TestCIService_TriggerWorkflow(t *testing.T) {
	setup()
	defer teardown()

	res_xml := `<Response><RequestId>NjJmMWZmMzRfOTBmYTUwNjRfNjVmOF8x</RequestId><InstanceId>i6fc78ca77d6011eba0ac5254008618d9</InstanceId></Response>`

	mux.HandleFunc("/triggerworkflow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := values{
			"workflowId": "w1234567890",
			"object":     "test.mp4",
			"name":       "trigger",
		}
		testFormValues(t, r, v)
		fmt.Fprint(w, res_xml)
	})

	opt := &TriggerWorkflowOptions{
		WorkflowId: "w1234567890",
		Object:     "test.mp4",
		Name:       "trigger",
	}

	res, _, err := client.CI.TriggerWorkflow(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.TriggerWorkflow returned error: %v", err)
	}
	want := &TriggerWorkflowResult{}
	err = xml.Unmarshal([]byte(res_xml), want)
	if err != nil {
		t.Errorf("CI.TriggerWorkflow Unmarshal returned error %+v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.TriggerWorkflow returned %+v, want %+v", res, want)
	}
}

func TestCIService_DescribeWorkflowExecutions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/workflowexecution", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"workflowId":        "w1234567890",
			"name":              "test.mp4",
			"orderByTime":       "Asc",
			"size":              "50",
			"states":            "Failed",
			"startCreationTime": "2022-02-25T12:00:00z",
			"endCreationTime":   "2022-02-28T12:00:00z",
			"nextToken":         "123",
			"jobId":             "b1234567890",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeWorkflowExecutionsOptions{
		WorkflowId:        "w1234567890",
		Name:              "test.mp4",
		OrderByTime:       "Asc",
		Size:              50,
		States:            "Failed",
		StartCreationTime: "2022-02-25T12:00:00z",
		EndCreationTime:   "2022-02-28T12:00:00z",
		NextToken:         "123",
		JobId:             "b1234567890",
	}

	_, _, err := client.CI.DescribeWorkflowExecutions(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeWorkflowExecutions returned error: %v", err)
	}
}

func TestCIService_DescribeWorkflowExecution(t *testing.T) {
	setup()
	defer teardown()

	res_xml := `<Response>
    <RequestId>NjJmMjA5MjZfZWM0YTYyNjRfN2U3ZF8yNzk1</RequestId>
    <WorkflowExecution>
        <WorkflowId>web6ac56c1ef54dbfa44d7f4103203be9</WorkflowId>
        <Name>workflow-1</Name>
        <RunId>i166ee19017b011eda8a5525400c540df</RunId>
        <CreateTime>2022-08-09T14:54:17+08:00</CreateTime>
        <Object>wk-test/game.mp4</Object>
        <State>Success</State>
        <Topology>
            <Dependencies>
                <Start>Transcode_1581665960537</Start>
                <Snapshot_1581665960536>End</Snapshot_1581665960536>
                <Transcode_1581665960537>Snapshot_1581665960536</Transcode_1581665960537>
            </Dependencies>
            <Nodes>
                <Start>
                    <Type>Start</Type>
                    <Input>
                        <QueueId>p09d709939fef48a0a5c247ef39d90cec</QueueId>
                        <ObjectPrefix>/wk-test</ObjectPrefix>
                        <ExtFilter>
                            <State>On</State>
                            <Video>false</Video>
                            <Audio>false</Audio>
                            <ContentType>false</ContentType>
                            <Custom>true</Custom>
                            <CustomExts>mp4</CustomExts>
                            <AllFile>false</AllFile>
                            <Image>false</Image>
                        </ExtFilter>
                        <PicProcessQueueId>p2911917386e148639319e13c285cc774</PicProcessQueueId>
                    </Input>
                </Start>
                <Snapshot_1581665960536>
                    <Type>Snapshot</Type>
                    <Operation>
                        <TemplateId>t07740e32081b44ad7a0aea03adcffd54a</TemplateId>
                        <Output>
                            <Region>ap-chongqing</Region>
                            <Bucket>test-1234567890</Bucket>
                            <Object>snapshot-${number}.jpg</Object>
                        </Output>
                    </Operation>
                </Snapshot_1581665960536>
                <Transcode_1581665960537>
                    <Type>Transcode</Type>
                    <Operation>
                        <TemplateId>t01e57db1c2d154d2fb57aa5de9313a897</TemplateId>
                        <Output>
                            <Region>ap-chongqing</Region>
                            <Bucket>test-1234567890</Bucket>
                            <Object>trans1.mp4</Object>
                        </Output>
                    </Operation>
                </Transcode_1581665960537>
            </Nodes>
        </Topology>
        <Tasks>
            <Type>Snapshot</Type>
            <JobId>j23c11e1e17b011edaab4ab15ec33d076</JobId>
            <CreateTime>2022-08-09T14:54:40+08:00</CreateTime>
            <Name>Snapshot_1581665960536</Name>
            <State>Success</State>
            <StartTime>2022-08-09T14:54:40+08:00</StartTime>
            <EndTime>2022-08-09T14:54:42+08:00</EndTime>
            <Code>Success</Code>
            <Message></Message>
        </Tasks>
        <Tasks>
            <Type>Transcode</Type>
            <JobId>j168668b217b011ed8efb27bb229e2d11</JobId>
            <CreateTime>2022-08-09T14:54:18+08:00</CreateTime>
            <Name>Transcode_1581665960537</Name>
            <State>Success</State>
            <StartTime>2022-08-09T14:54:18+08:00</StartTime>
            <EndTime>2022-08-09T14:54:39+08:00</EndTime>
            <Code>Success</Code>
            <Message>success</Message>
        </Tasks>
    </WorkflowExecution>
</Response>`
	want := &DescribeWorkflowExecutionResult{}
	err := xml.Unmarshal([]byte(res_xml), want)
	if err != nil {
		t.Errorf("CI.DescribeWorkflowExecution Unmarshal returned error %+v", err)
	}

	mux.HandleFunc("/workflowexecution/i1234567890", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, res_xml)
	})

	res, _, err := client.CI.DescribeWorkflowExecution(context.Background(), "i1234567890")
	if err != nil {
		t.Fatalf("CI.DescribeWorkflowExecution returned error: %v", err)
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("CI.DescribeWorkflowExecution returned %+v, want %+v", res, want)
	}
}

func TestCIService_CreateASRJobs(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>SpeechRecognition</Tag><Input><Object>test.mp3</Object></Input>" +
		"<Operation><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>out.txt</Object></Output>" +
		"<TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/asr_jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreateASRJobsOptions{
		Tag: "SpeechRecognition",
		Input: &JobInput{
			Object: "test.mp3",
		},
		Operation: &ASRJobOperation{
			Output: &JobOutput{
				Region: "ap-beijing",
				Bucket: "abc-1250000000",
				Object: "out.txt",
			},
			TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreateASRJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateASRJobs returned errors: %v", err)
	}
}

func TestCIService_DescribeMultiASRJob(t *testing.T) {
	{
		setup()
		jobIDs := []string{"ac7c990a00bf211ed946af9e0691f2b7a", "aabcsdssfeipplsdfwe"}
		mux.HandleFunc("/asr_jobs/"+"ac7c990a00bf211ed946af9e0691f2b7a,aabcsdssfeipplsdfwe", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
		})

		_, _, err := client.CI.DescribeMultiASRJob(context.Background(), jobIDs)
		if err != nil {
			t.Fatalf("CI.DescribeMultiASRJob returned error: %v", err)
		}
		teardown()
	}

	{
		setup()
		jobIDs := []string{}
		mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
		})

		_, _, err := client.CI.DescribeMultiASRJob(context.Background(), jobIDs)
		if err == nil || err.Error() != "empty param jobids" {
			t.Fatalf("CI.DescribeMultiASRJob returned error: %v", err)
		}
		teardown()
	}
}

func TestCIService_DescribeMediaTemplate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"tag":        "All",
			"category":   "Custom",
			"ids":        "t123456798",
			"name":       "test",
			"pageNumber": "1",
			"pageSize":   "20",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaTemplateOptions{
		Tag:        "All",
		Category:   "Custom",
		Ids:        "t123456798",
		Name:       "test",
		PageNumber: 1,
		PageSize:   20,
	}

	_, _, err := client.CI.DescribeMediaTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaTemplate returned error: %v", err)
	}
}

func TestCIService_DeleteMediaTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "tc7c990a00bf211ed946af9e0691f2b7a"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, _, err := client.CI.DeleteMediaTemplate(context.Background(), tplId)
	if err != nil {
		t.Fatalf("CI.DeleteMediaTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaSnapshotTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Snapshot</Tag><Name>test-Snapshot</Name>" +
		"<Snapshot><Mode>Interval</Mode><Start>0</Start><TimeInterval>0.5</TimeInterval><Count>10</Count>" +
		"<Width>1280</Width><Height>960</Height><CIParam>imageMogr2/thumbnail/!50p</CIParam><IsCheckCount>true</IsCheckCount>" +
		"<IsCheckBlack>true</IsCheckBlack><BlackLevel>30</BlackLevel><PixelBlackThreshold>100</PixelBlackThreshold><SnapshotOutMode>SnapshotAndSprite</SnapshotOutMode>" +
		"<SpriteSnapshotConfig><CellHeight>1024</CellHeight><CellWidth>768</CellWidth><Color>Aquamarine</Color>" +
		"<Columns>9</Columns><Lines>9</Lines><Margin>64</Margin><Padding>32</Padding><ScaleMethod>MaxWHScale</ScaleMethod></SpriteSnapshotConfig></Snapshot></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSnapshotTemplateOptions{
		Tag:  "Snapshot",
		Name: "test-Snapshot",
		Snapshot: &Snapshot{
			Mode:                "Interval",
			Start:               "0",
			TimeInterval:        "0.5",
			Count:               "10",
			Width:               "1280",
			Height:              "960",
			CIParam:             "imageMogr2/thumbnail/!50p",
			IsCheckCount:        true,
			IsCheckBlack:        true,
			BlackLevel:          "30",
			PixelBlackThreshold: "100",
			SnapshotOutMode:     "SnapshotAndSprite",
			SpriteSnapshotConfig: &SpriteSnapshotConfig{
				CellHeight:  "1024",
				CellWidth:   "768",
				Color:       "Aquamarine",
				Columns:     "9",
				Lines:       "9",
				Margin:      "64",
				Padding:     "32",
				ScaleMethod: "MaxWHScale",
			},
		},
	}

	_, _, err := client.CI.CreateMediaSnapshotTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaSnapshotTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaSnapshotTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1f16e1dfbdc994105b31292d45710642a"

	wantBody := "<Request><Tag>Snapshot</Tag><Name>test-Snapshot-update</Name>" +
		"<Snapshot><Mode>Interval</Mode><Start>0</Start><TimeInterval>0.5</TimeInterval><Count>10</Count>" +
		"<Width>1280</Width><Height>960</Height><CIParam>imageMogr2/thumbnail/!50p</CIParam><IsCheckCount>true</IsCheckCount>" +
		"<IsCheckBlack>true</IsCheckBlack><BlackLevel>30</BlackLevel><PixelBlackThreshold>100</PixelBlackThreshold><SnapshotOutMode>SnapshotAndSprite</SnapshotOutMode>" +
		"<SpriteSnapshotConfig><CellHeight>1024</CellHeight><CellWidth>768</CellWidth><Color>Aquamarine</Color>" +
		"<Columns>9</Columns><Lines>9</Lines><Margin>64</Margin><Padding>32</Padding><ScaleMethod>MaxWHScale</ScaleMethod></SpriteSnapshotConfig></Snapshot></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSnapshotTemplateOptions{
		Tag:  "Snapshot",
		Name: "test-Snapshot-update",
		Snapshot: &Snapshot{
			Mode:                "Interval",
			Start:               "0",
			TimeInterval:        "0.5",
			Count:               "10",
			Width:               "1280",
			Height:              "960",
			CIParam:             "imageMogr2/thumbnail/!50p",
			IsCheckCount:        true,
			IsCheckBlack:        true,
			BlackLevel:          "30",
			PixelBlackThreshold: "100",
			SnapshotOutMode:     "SnapshotAndSprite",
			SpriteSnapshotConfig: &SpriteSnapshotConfig{
				CellHeight:  "1024",
				CellWidth:   "768",
				Color:       "Aquamarine",
				Columns:     "9",
				Lines:       "9",
				Margin:      "64",
				Padding:     "32",
				ScaleMethod: "MaxWHScale",
			},
		},
	}

	_, _, err := client.CI.UpdateMediaSnapshotTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaSnapshotTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaTranscodeTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Transcode</Tag><Name>test-Transcode</Name>" +
		"<Container><Format>hls</Format><ClipConfig><Duration>10</Duration></ClipConfig></Container>" +
		"<Video><Codec>h.264</Codec><Bitrate>25000</Bitrate></Video>" +
		"<Audio><Codec>mp3</Codec><Bitrate>64</Bitrate></Audio>" +
		"<TimeInterval><Start>5.5</Start><Duration>10.5</Duration></TimeInterval>" +
		"<TransConfig><DeleteMetadata>true</DeleteMetadata></TransConfig>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"</Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaTranscodeTemplateOptions{
		Tag:  "Transcode",
		Name: "test-Transcode",
		Container: &Container{
			Format: "hls",
			ClipConfig: &ClipConfig{
				Duration: "10",
			},
		},
		Video: &Video{
			Codec:   "h.264",
			Bitrate: "25000",
		},
		Audio: &Audio{
			Codec:   "mp3",
			Bitrate: "64",
		},
		TimeInterval: &TimeInterval{
			Start:    "5.5",
			Duration: "10.5",
		},
		TransConfig: &TransConfig{
			DeleteMetadata: "true",
		},
		AudioMixArray: []AudioMix{
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
		},
	}

	_, _, err := client.CI.CreateMediaTranscodeTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaTranscodeTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaTranscodeTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1f16e1dfbdc994105b31292d45710642a"

	wantBody := "<Request><Tag>Transcode</Tag><Name>test-Transcode-update</Name>" +
		"<Container><Format>hls</Format><ClipConfig><Duration>10</Duration></ClipConfig></Container>" +
		"<Video><Codec>h.264</Codec><Bitrate>25000</Bitrate></Video>" +
		"<Audio><Codec>mp3</Codec><Bitrate>64</Bitrate></Audio>" +
		"<TimeInterval><Start>5.5</Start><Duration>10.5</Duration></TimeInterval>" +
		"<TransConfig><DeleteMetadata>true</DeleteMetadata></TransConfig>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"</Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaTranscodeTemplateOptions{
		Tag:  "Transcode",
		Name: "test-Transcode-update",
		Container: &Container{
			Format: "hls",
			ClipConfig: &ClipConfig{
				Duration: "10",
			},
		},
		Video: &Video{
			Codec:   "h.264",
			Bitrate: "25000",
		},
		Audio: &Audio{
			Codec:   "mp3",
			Bitrate: "64",
		},
		TimeInterval: &TimeInterval{
			Start:    "5.5",
			Duration: "10.5",
		},
		TransConfig: &TransConfig{
			DeleteMetadata: "true",
		},
		AudioMixArray: []AudioMix{
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
		},
	}

	_, _, err := client.CI.UpdateMediaTranscodeTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaTranscodeTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaAnimationTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Animation</Tag><Name>test-Animation</Name>" +
		"<Container><Format>gif</Format></Container>" +
		"<Video><Codec>gif</Codec><Width>1024</Width><Height>768</Height><Fps>30</Fps><AnimateOnlyKeepKeyFrame>true</AnimateOnlyKeepKeyFrame><Quality>80</Quality></Video>" +
		"<TimeInterval><Start>2.5</Start><Duration>7.5</Duration></TimeInterval>" +
		"</Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaAnimationTemplateOptions{
		Tag:  "Animation",
		Name: "test-Animation",
		Container: &Container{
			Format: "gif",
		},
		Video: &AnimationVideo{
			Codec:                   "gif",
			Width:                   "1024",
			Height:                  "768",
			Fps:                     "30",
			AnimateOnlyKeepKeyFrame: "true",
			Quality:                 "80",
		},
		TimeInterval: &TimeInterval{
			Start:    "2.5",
			Duration: "7.5",
		},
	}

	_, _, err := client.CI.CreateMediaAnimationTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaAnimationTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaAnimationTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1f16e1dfbdc994105b31292d45710642a"

	wantBody := "<Request><Tag>Animation</Tag><Name>test-Animation</Name>" +
		"<Container><Format>gif</Format></Container>" +
		"<Video><Codec>gif</Codec><Width>1024</Width><Height>768</Height><Fps>30</Fps><AnimateOnlyKeepKeyFrame>true</AnimateOnlyKeepKeyFrame><Quality>80</Quality></Video>" +
		"<TimeInterval><Start>2.5</Start><Duration>7.5</Duration></TimeInterval>" +
		"</Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaAnimationTemplateOptions{
		Tag:  "Animation",
		Name: "test-Animation",
		Container: &Container{
			Format: "gif",
		},
		Video: &AnimationVideo{
			Codec:                   "gif",
			Width:                   "1024",
			Height:                  "768",
			Fps:                     "30",
			AnimateOnlyKeepKeyFrame: "true",
			Quality:                 "80",
		},
		TimeInterval: &TimeInterval{
			Start:    "2.5",
			Duration: "7.5",
		},
	}

	_, _, err := client.CI.UpdateMediaAnimationTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaAnimationTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaConcatTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Concat</Tag><Name>test-Concat</Name><ConcatTemplate>" +
		"<ConcatFragment><Url>http://bucket-1250000000.cos.ap-beijing.myqcloud.com/start.mp4</Url><Mode>Start</Mode></ConcatFragment>" +
		"<ConcatFragment><Url>http://bucket-1250000000.cos.ap-beijing.myqcloud.com/end.mp4</Url><Mode>End</Mode></ConcatFragment>" +
		"<Audio><Codec>mp3</Codec></Audio>" +
		"<Video><Codec>H.264</Codec><Width>1280</Width><Fps>30</Fps><Bitrate>1000</Bitrate></Video>" +
		"<Container><Format>mp4</Format></Container>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3</AudioSource><MixMode>Once</MixMode><Replace>true</Replace>" +
		"<EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime><EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.7</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3</AudioSource><MixMode>Once</MixMode><Replace>true</Replace>" +
		"<EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime><EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.7</BgmFadeTime></EffectConfig></AudioMixArray></ConcatTemplate></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaConcatTemplateOptions{
		Tag:  "Concat",
		Name: "test-Concat",
		ConcatTemplate: &ConcatTemplate{
			ConcatFragment: []ConcatFragment{
				{
					Url:  "http://bucket-1250000000.cos.ap-beijing.myqcloud.com/start.mp4",
					Mode: "Start",
				},
				{
					Url:  "http://bucket-1250000000.cos.ap-beijing.myqcloud.com/end.mp4",
					Mode: "End",
				},
			},
			Audio: &Audio{
				Codec: "mp3",
			},
			Video: &Video{
				Codec:   "H.264",
				Width:   "1280",
				Fps:     "30",
				Bitrate: "1000",
			},
			Container: &Container{
				Format: "mp4",
			},
			AudioMixArray: []AudioMix{
				{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3",
					MixMode: "Once",
					Replace: "true",
					EffectConfig: &EffectConfig{
						EnableStartFadein: "true",
						StartFadeinTime:   "3",
						EnableEndFadeout:  "false",
						EndFadeoutTime:    "0",
						EnableBgmFade:     "true",
						BgmFadeTime:       "1.7",
					}},
				{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3",
					MixMode: "Once",
					Replace: "true",
					EffectConfig: &EffectConfig{
						EnableStartFadein: "true",
						StartFadeinTime:   "3",
						EnableEndFadeout:  "false",
						EndFadeoutTime:    "0",
						EnableBgmFade:     "true",
						BgmFadeTime:       "1.7",
					}},
			},
		},
	}

	_, _, err := client.CI.CreateMediaConcatTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaConcatTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaConcatTemplate(t *testing.T) {
	setup()
	defer teardown()
	tplId := "t1f16e1dfbdc994105b31292d45710642a"
	wantBody := "<Request><Tag>Concat</Tag><Name>test-Concat-update</Name><ConcatTemplate>" +
		"<ConcatFragment><Url>http://bucket-1250000000.cos.ap-beijing.myqcloud.com/start.mp4</Url><Mode>Start</Mode></ConcatFragment>" +
		"<ConcatFragment><Url>http://bucket-1250000000.cos.ap-beijing.myqcloud.com/end.mp4</Url><Mode>End</Mode></ConcatFragment>" +
		"<Audio><Codec>mp3</Codec></Audio>" +
		"<Video><Codec>H.264</Codec><Width>1280</Width><Fps>30</Fps><Bitrate>1000</Bitrate></Video>" +
		"<Container><Format>mp4</Format></Container>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3</AudioSource><MixMode>Once</MixMode><Replace>true</Replace>" +
		"<EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime><EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.7</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3</AudioSource><MixMode>Once</MixMode><Replace>true</Replace>" +
		"<EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime><EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.7</BgmFadeTime></EffectConfig></AudioMixArray></ConcatTemplate></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaConcatTemplateOptions{
		Tag:  "Concat",
		Name: "test-Concat-update",
		ConcatTemplate: &ConcatTemplate{
			ConcatFragment: []ConcatFragment{
				{
					Url:  "http://bucket-1250000000.cos.ap-beijing.myqcloud.com/start.mp4",
					Mode: "Start",
				},
				{
					Url:  "http://bucket-1250000000.cos.ap-beijing.myqcloud.com/end.mp4",
					Mode: "End",
				},
			},
			Audio: &Audio{
				Codec: "mp3",
			},
			Video: &Video{
				Codec:   "H.264",
				Width:   "1280",
				Fps:     "30",
				Bitrate: "1000",
			},
			Container: &Container{
				Format: "mp4",
			},
			AudioMixArray: []AudioMix{
				{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3",
					MixMode: "Once",
					Replace: "true",
					EffectConfig: &EffectConfig{
						EnableStartFadein: "true",
						StartFadeinTime:   "3",
						EnableEndFadeout:  "false",
						EndFadeoutTime:    "0",
						EnableBgmFade:     "true",
						BgmFadeTime:       "1.7",
					}},
				{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3",
					MixMode: "Once",
					Replace: "true",
					EffectConfig: &EffectConfig{
						EnableStartFadein: "true",
						StartFadeinTime:   "3",
						EnableEndFadeout:  "false",
						EndFadeoutTime:    "0",
						EnableBgmFade:     "true",
						BgmFadeTime:       "1.7",
					}},
			},
		},
	}

	_, _, err := client.CI.UpdateMediaConcatTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaConcatTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaVideoProcessTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>VideoProcess</Tag><Name>test-VideoProcess</Name>" +
		"<ColorEnhance><Enable>true</Enable><Contrast>50</Contrast></ColorEnhance>" +
		"<MsSharpen><Enable>true</Enable><SharpenLevel>5</SharpenLevel></MsSharpen></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaVideoProcessTemplateOptions{
		Tag:  "VideoProcess",
		Name: "test-VideoProcess",
		ColorEnhance: &ColorEnhance{
			Enable:   "true",
			Contrast: "50",
		},
		MsSharpen: &MsSharpen{
			Enable:       "true",
			SharpenLevel: "5",
		},
	}

	_, _, err := client.CI.CreateMediaVideoProcessTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaVideoProcessTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaVideoProcessTemplate(t *testing.T) {
	setup()
	defer teardown()
	tplId := "t1f16e1dfbdc994105b31292d45710642a"
	wantBody := "<Request><Tag>VideoProcess</Tag><Name>test-VideoProcess-update</Name>" +
		"<ColorEnhance><Enable>true</Enable><Contrast>50</Contrast></ColorEnhance>" +
		"<MsSharpen><Enable>true</Enable><SharpenLevel>5</SharpenLevel></MsSharpen></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaVideoProcessTemplateOptions{
		Tag:  "VideoProcess",
		Name: "test-VideoProcess-update",
		ColorEnhance: &ColorEnhance{
			Enable:   "true",
			Contrast: "50",
		},
		MsSharpen: &MsSharpen{
			Enable:       "true",
			SharpenLevel: "5",
		},
	}

	_, _, err := client.CI.UpdateMediaVideoProcessTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaVideoProcessTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaVideoMontageTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Transcode</Tag><Name>test-Transcode</Name><Duration>10.5</Duration>" +
		"<Container><Format>mp4</Format></Container>" +
		"<Video><Codec>h.264</Codec><Bitrate>25000</Bitrate></Video>" +
		"<Audio><Codec>mp3</Codec><Bitrate>64</Bitrate></Audio>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"</Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaVideoMontageTemplateOptions{
		Tag:      "Transcode",
		Name:     "test-Transcode",
		Duration: "10.5",
		Container: &Container{
			Format: "mp4",
		},
		Video: &Video{
			Codec:   "h.264",
			Bitrate: "25000",
		},
		Audio: &Audio{
			Codec:   "mp3",
			Bitrate: "64",
		},
		AudioMixArray: []AudioMix{
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
		},
	}

	_, _, err := client.CI.CreateMediaVideoMontageTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaVideoMontageTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaVideoMontageTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1f16e1dfbdc994105b31292d45710642a"
	wantBody := "<Request><Tag>VideoMontage</Tag><Name>test-VideoMontage-update</Name><Duration>10.5</Duration>" +
		"<Container><Format>mp4</Format></Container>" +
		"<Video><Codec>h.264</Codec><Bitrate>25000</Bitrate></Video>" +
		"<Audio><Codec>mp3</Codec><Bitrate>64</Bitrate></Audio>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"<AudioMixArray><AudioSource>https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3</AudioSource><MixMode>Once</MixMode>" +
		"<Replace>true</Replace><EffectConfig><EnableStartFadein>true</EnableStartFadein><StartFadeinTime>3</StartFadeinTime>" +
		"<EnableEndFadeout>false</EnableEndFadeout><EndFadeoutTime>0</EndFadeoutTime>" +
		"<EnableBgmFade>true</EnableBgmFade><BgmFadeTime>1.5</BgmFadeTime></EffectConfig></AudioMixArray>" +
		"</Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaVideoMontageTemplateOptions{
		Tag:      "VideoMontage",
		Name:     "test-VideoMontage-update",
		Duration: "10.5",
		Container: &Container{
			Format: "mp4",
		},
		Video: &Video{
			Codec:   "h.264",
			Bitrate: "25000",
		},
		Audio: &Audio{
			Codec:   "mp3",
			Bitrate: "64",
		},
		AudioMixArray: []AudioMix{
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix1.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
			{AudioSource: "https://test-xxx.cos.ap-chongqing.myqcloud.com/mix2.mp3",
				MixMode: "Once",
				Replace: "true",
				EffectConfig: &EffectConfig{
					EnableStartFadein: "true",
					StartFadeinTime:   "3",
					EnableEndFadeout:  "false",
					EndFadeoutTime:    "0",
					EnableBgmFade:     "true",
					BgmFadeTime:       "1.5",
				}},
		},
	}

	_, _, err := client.CI.UpdateMediaVideoMontageTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaVideoMontageTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaVoiceSeparateTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>VoiceSeparate</Tag><Name>test-VoiceSeparate</Name><AudioMode>IsAudio</AudioMode>" +
		"<AudioConfig><Codec>aac</Codec><Samplerate>44100</Samplerate><Bitrate>128</Bitrate><Channels>4</Channels></AudioConfig></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaVoiceSeparateTemplateOptions{
		Tag:       "VoiceSeparate",
		Name:      "test-VoiceSeparate",
		AudioMode: "IsAudio",
		AudioConfig: &AudioConfig{
			Codec:      "aac",
			Samplerate: "44100",
			Bitrate:    "128",
			Channels:   "4",
		},
	}

	_, _, err := client.CI.CreateMediaVoiceSeparateTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaVoiceSeparateTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaVoiceSeparateTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	wantBody := "<Request><Tag>VoiceSeparate</Tag><Name>test-VoiceSeparate-update</Name><AudioMode>IsAudio</AudioMode>" +
		"<AudioConfig><Codec>aac</Codec><Samplerate>44100</Samplerate><Bitrate>128</Bitrate><Channels>4</Channels></AudioConfig></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaVoiceSeparateTemplateOptions{
		Tag:       "VoiceSeparate",
		Name:      "test-VoiceSeparate-update",
		AudioMode: "IsAudio",
		AudioConfig: &AudioConfig{
			Codec:      "aac",
			Samplerate: "44100",
			Bitrate:    "128",
			Channels:   "4",
		},
	}

	_, _, err := client.CI.UpdateMediaVoiceSeparateTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaVoiceSeparateTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaSuperResolutionTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>SuperResolution</Tag><Name>test-SuperResolution</Name><Resolution>sdtohd</Resolution>" +
		"<EnableScaleUp>true</EnableScaleUp><Version>Enhance</Version></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSuperResolutionTemplateOptions{
		Tag:           "SuperResolution",
		Name:          "test-SuperResolution",
		Resolution:    "sdtohd",
		EnableScaleUp: "true",
		Version:       "Enhance",
	}

	_, _, err := client.CI.CreateMediaSuperResolutionTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaSuperResolutionTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaSuperResolutionTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	wantBody := "<Request><Tag>SuperResolution</Tag><Name>test-SuperResolution-update</Name><Resolution>sdtohd</Resolution>" +
		"<EnableScaleUp>true</EnableScaleUp><Version>Enhance</Version></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSuperResolutionTemplateOptions{
		Tag:           "SuperResolution",
		Name:          "test-SuperResolution-update",
		Resolution:    "sdtohd",
		EnableScaleUp: "true",
		Version:       "Enhance",
	}

	_, _, err := client.CI.UpdateMediaSuperResolutionTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaSuperResolutionTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaPicProcessTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>PicProcess</Tag><Name>test-PicProcess</Name>" +
		"<PicProcess><IsPicInfo>true</IsPicInfo><ProcessRule>imageMogr2/rotate/90</ProcessRule></PicProcess></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaPicProcessTemplateOptions{
		Tag:  "PicProcess",
		Name: "test-PicProcess",
		PicProcess: &PicProcess{
			IsPicInfo:   "true",
			ProcessRule: "imageMogr2/rotate/90",
		},
	}

	_, _, err := client.CI.CreateMediaPicProcessTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaPicProcessTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaPicProcessTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	wantBody := "<Request><Tag>PicProcess</Tag><Name>test-PicProcess-update</Name>" +
		"<PicProcess><IsPicInfo>true</IsPicInfo><ProcessRule>imageMogr2/rotate/90</ProcessRule></PicProcess></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaPicProcessTemplateOptions{
		Tag:  "PicProcess",
		Name: "test-PicProcess-update",
		PicProcess: &PicProcess{
			IsPicInfo:   "true",
			ProcessRule: "imageMogr2/rotate/90",
		},
	}

	_, _, err := client.CI.UpdateMediaPicProcessTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaPicProcessTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaWatermarkTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Watermark</Tag><Name>test-Watermark</Name><Watermark><Type>Text</Type><Pos>TopRight</Pos>" +
		"<LocMode>Absolute</LocMode><Dx>128</Dx><Dy>128</Dy><StartTime>0</StartTime><EndTime>100.5</EndTime>" +
		"<Text><FontSize>30</FontSize><FontType>simfang.ttf</FontType><FontColor>0xaabbcc</FontColor><Transparency>30</Transparency><Text>水印内容</Text></Text></Watermark></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaWatermarkTemplateOptions{
		Tag:  "Watermark",
		Name: "test-Watermark",
		Watermark: &Watermark{
			Type:      "Text",
			Pos:       "TopRight",
			LocMode:   "Absolute",
			Dx:        "128",
			Dy:        "128",
			StartTime: "0",
			EndTime:   "100.5",
			Text: &Text{
				FontSize:     "30",
				FontType:     "simfang.ttf",
				FontColor:    "0xaabbcc",
				Transparency: "30",
				Text:         "水印内容",
			},
		},
	}

	_, _, err := client.CI.CreateMediaWatermarkTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaWatermarkTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaWatermarkTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	wantBody := "<Request><Tag>Watermark</Tag><Name>test-Watermark-update</Name><Watermark><Type>Text</Type><Pos>TopRight</Pos>" +
		"<LocMode>Absolute</LocMode><Dx>128</Dx><Dy>128</Dy><StartTime>0</StartTime><EndTime>100.5</EndTime>" +
		"<Text><FontSize>30</FontSize><FontType>simfang.ttf</FontType><FontColor>0xaabbcc</FontColor><Transparency>30</Transparency><Text>水印内容</Text></Text></Watermark></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaWatermarkTemplateOptions{
		Tag:  "Watermark",
		Name: "test-Watermark-update",
		Watermark: &Watermark{
			Type:      "Text",
			Pos:       "TopRight",
			LocMode:   "Absolute",
			Dx:        "128",
			Dy:        "128",
			StartTime: "0",
			EndTime:   "100.5",
			Text: &Text{
				FontSize:     "30",
				FontType:     "simfang.ttf",
				FontColor:    "0xaabbcc",
				Transparency: "30",
				Text:         "水印内容",
			},
		},
	}

	_, _, err := client.CI.UpdateMediaWatermarkTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaWatermarkTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaTranscodeProTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>TranscodePro</Tag><Name>test-TranscodePro</Name><Container><Format>mxf</Format></Container><Video><Codec>xavc</Codec>" +
		"<Profile>XAVC-HD_422_10bit</Profile><Width>1920</Width><Height>1080</Height><Interlaced>true</Interlaced><Fps>30000/1001</Fps><Bitrate>50000</Bitrate>" +
		"</Video><Audio><Codec>pcm_s24le</Codec></Audio><TimeInterval><Start>0</Start><Duration>60</Duration></TimeInterval><TransConfig><AdjDarMethod>scale</AdjDarMethod>" +
		"<IsCheckReso>true</IsCheckReso><ResoAdjMethod>1</ResoAdjMethod></TransConfig></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaTranscodeProTemplateOptions{
		Tag:  "TranscodePro",
		Name: "test-TranscodePro",
		Container: &Container{
			Format: "mxf",
		},
		Video: &TranscodeProVideo{
			Codec:      "xavc",
			Profile:    "XAVC-HD_422_10bit",
			Width:      "1920",
			Height:     "1080",
			Interlaced: "true",
			Fps:        "30000/1001",
			Bitrate:    "50000",
		},
		Audio: &TranscodeProAudio{
			Codec: "pcm_s24le",
		},
		TimeInterval: &TimeInterval{
			Start:    "0",
			Duration: "60",
		},
		TransConfig: &TransConfig{
			AdjDarMethod:  "scale",
			IsCheckReso:   "true",
			ResoAdjMethod: "1",
		},
	}

	_, _, err := client.CI.CreateMediaTranscodeProTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaTranscodeProTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaTranscodeProTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	wantBody := "<Request><Tag>TranscodePro</Tag><Name>test-TranscodePro-update</Name><Container><Format>mxf</Format></Container><Video><Codec>xavc</Codec>" +
		"<Profile>XAVC-HD_422_10bit</Profile><Width>1920</Width><Height>1080</Height><Interlaced>true</Interlaced><Fps>30000/1001</Fps><Bitrate>50000</Bitrate>" +
		"</Video><Audio><Codec>pcm_s24le</Codec></Audio><TimeInterval><Start>0</Start><Duration>60</Duration></TimeInterval><TransConfig><AdjDarMethod>scale</AdjDarMethod>" +
		"<IsCheckReso>true</IsCheckReso><ResoAdjMethod>1</ResoAdjMethod></TransConfig></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaTranscodeProTemplateOptions{
		Tag:  "TranscodePro",
		Name: "test-TranscodePro-update",
		Container: &Container{
			Format: "mxf",
		},
		Video: &TranscodeProVideo{
			Codec:      "xavc",
			Profile:    "XAVC-HD_422_10bit",
			Width:      "1920",
			Height:     "1080",
			Interlaced: "true",
			Fps:        "30000/1001",
			Bitrate:    "50000",
		},
		Audio: &TranscodeProAudio{
			Codec: "pcm_s24le",
		},
		TimeInterval: &TimeInterval{
			Start:    "0",
			Duration: "60",
		},
		TransConfig: &TransConfig{
			AdjDarMethod:  "scale",
			IsCheckReso:   "true",
			ResoAdjMethod: "1",
		},
	}

	_, _, err := client.CI.UpdateMediaTranscodeProTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaTranscodeProTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaTtsTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Tts</Tag><Name>test-Tts</Name><Mode>Sync</Mode><Codec>pcm</Codec><VoiceType>ruxue</VoiceType><Volume>2</Volume><Speed>200</Speed></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaTtsTemplateOptions{
		Tag:       "Tts",
		Name:      "test-Tts",
		Mode:      "Sync",
		Codec:     "pcm",
		VoiceType: "ruxue",
		Volume:    "2",
		Speed:     "200",
	}

	_, _, err := client.CI.CreateMediaTtsTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaTtsTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaTtsTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	wantBody := "<Request><Tag>Tts</Tag><Name>test-Tts-update</Name><Mode>Sync</Mode><Codec>pcm</Codec><VoiceType>ruxue</VoiceType><Volume>2</Volume><Speed>200</Speed></Request>"

	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaTtsTemplateOptions{
		Tag:       "Tts",
		Name:      "test-Tts-update",
		Mode:      "Sync",
		Codec:     "pcm",
		VoiceType: "ruxue",
		Volume:    "2",
		Speed:     "200",
	}

	_, _, err := client.CI.UpdateMediaTtsTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaTtsTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaSmartCoverTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>SmartCover</Tag><Name>test-SmartCover</Name><SmartCover><Format>jpg</Format>" +
		"<Width>1280</Width><Height>960</Height><Count>10</Count><DeleteDuplicates>true</DeleteDuplicates></SmartCover></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSmartCoverTemplateOptions{
		Tag:  "SmartCover",
		Name: "test-SmartCover",
		SmartCover: &NodeSmartCover{
			Format:           "jpg",
			Width:            "1280",
			Height:           "960",
			Count:            "10",
			DeleteDuplicates: "true",
		},
	}

	_, _, err := client.CI.CreateMediaSmartCoverTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaSmartCoverTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaSmartCoverTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>SmartCover</Tag><Name>test-SmartCover</Name><SmartCover><Format>jpg</Format>" +
		"<Width>1280</Width><Height>960</Height><Count>10</Count><DeleteDuplicates>true</DeleteDuplicates></SmartCover></Request>"

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSmartCoverTemplateOptions{
		Tag:  "SmartCover",
		Name: "test-SmartCover",
		SmartCover: &NodeSmartCover{
			Format:           "jpg",
			Width:            "1280",
			Height:           "960",
			Count:            "10",
			DeleteDuplicates: "true",
		},
	}

	_, _, err := client.CI.UpdateMediaSmartCoverTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaSmartCoverTemplate returned error: %v", err)
	}
}

func TestCIService_CreateMediaSpeechRecognitionTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>SpeechRecognition</Tag><Name>test-SpeechRecognition</Name><SpeechRecognition><ChannelNum>2</ChannelNum><ConvertNumMode>0</ConvertNumMode>" +
		"<EngineModelType>16k_zh</EngineModelType><FilterDirty>0</FilterDirty><FilterModal>1</FilterModal><ResTextFormat>1</ResTextFormat><SpeakerDiarization>1</SpeakerDiarization>" +
		"<SpeakerNumber>0</SpeakerNumber><FilterPunc>0</FilterPunc><OutputFileType>txt</OutputFileType></SpeechRecognition></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSpeechRecognitionTemplateOptions{
		Tag:  "SpeechRecognition",
		Name: "test-SpeechRecognition",
		SpeechRecognition: &SpeechRecognition{
			ChannelNum:         "2",
			ConvertNumMode:     "0",
			EngineModelType:    "16k_zh",
			FilterDirty:        "0",
			FilterModal:        "1",
			ResTextFormat:      "1",
			SpeakerDiarization: "1",
			SpeakerNumber:      "0",
			FilterPunc:         "0",
			OutputFileType:     "txt",
		},
	}

	_, _, err := client.CI.CreateMediaSpeechRecognitionTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaSpeechRecognitionTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateMediaSpeechRecognitionTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>SpeechRecognition</Tag><Name>test-SpeechRecognition</Name><SpeechRecognition><ChannelNum>2</ChannelNum><ConvertNumMode>0</ConvertNumMode>" +
		"<EngineModelType>16k_zh</EngineModelType><FilterDirty>0</FilterDirty><FilterModal>1</FilterModal><ResTextFormat>1</ResTextFormat><SpeakerDiarization>1</SpeakerDiarization>" +
		"<SpeakerNumber>0</SpeakerNumber><FilterPunc>0</FilterPunc><OutputFileType>txt</OutputFileType></SpeechRecognition></Request>"

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaSpeechRecognitionTemplateOptions{
		Tag:  "SpeechRecognition",
		Name: "test-SpeechRecognition",
		SpeechRecognition: &SpeechRecognition{
			ChannelNum:         "2",
			ConvertNumMode:     "0",
			EngineModelType:    "16k_zh",
			FilterDirty:        "0",
			FilterModal:        "1",
			ResTextFormat:      "1",
			SpeakerDiarization: "1",
			SpeakerNumber:      "0",
			FilterPunc:         "0",
			OutputFileType:     "txt",
		},
	}

	_, _, err := client.CI.UpdateMediaSpeechRecognitionTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaSpeechRecognitionTemplate returned error: %v", err)
	}
}

func TestCIService_CreateNoiseReductionTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>NoiseReduction</Tag><Name>NoiseReduction-1</Name><NoiseReduction><Format>wav</Format><Samplerate>16000</Samplerate></NoiseReduction></Request>"

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateNoiseReductionTemplateOptions{
		Tag:  "NoiseReduction",
		Name: "NoiseReduction-1",
		NoiseReduction: &NoiseReduction{
			Format:     "wav",
			Samplerate: "16000",
		},
	}

	_, _, err := client.CI.CreateNoiseReductionTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateNoiseReductionTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateNoiseReductionTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>NoiseReduction</Tag><Name>NoiseReduction-1</Name><NoiseReduction><Format>mp3</Format><Samplerate>16000</Samplerate></NoiseReduction></Request>"

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateNoiseReductionTemplateOptions{
		Tag:  "NoiseReduction",
		Name: "NoiseReduction-1",
		NoiseReduction: &NoiseReduction{
			Format:     "mp3",
			Samplerate: "16000",
		},
	}

	_, _, err := client.CI.UpdateNoiseReductionTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateNoiseReductionTemplate returned error: %v", err)
	}
}

func TestCIService_CreateVideoEnhanceTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>VideoEnhance</Tag><Name>VideoEnhance-test1</Name><VideoEnhance>" +
		"<Transcode><Container><Format>mp4</Format></Container>" +
		"<Video><Codec>H.264</Codec><Width>1280</Width><Fps>30</Fps><Bitrate>1000</Bitrate></Video>" +
		"<Audio><Codec>aac</Codec><Samplerate>44100</Samplerate><Bitrate>128</Bitrate><Channels>4</Channels></Audio></Transcode>" +
		"<SuperResolution><Resolution>sdtohd</Resolution><EnableScaleUp>true</EnableScaleUp><Version>Enhance</Version></SuperResolution>" +
		"<ColorEnhance><Contrast>50</Contrast><Correction>100</Correction><Saturation>100</Saturation></ColorEnhance>" +
		"<MsSharpen><SharpenLevel>5</SharpenLevel></MsSharpen><SDRtoHDR><HdrMode>HDR10</HdrMode></SDRtoHDR><FrameEnhance><FrameDoubling>true</FrameDoubling></FrameEnhance></VideoEnhance></Request>"
	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateVideoEnhanceTemplateOptions{
		Tag:  "VideoEnhance",
		Name: "VideoEnhance-test1",
		VideoEnhance: &VideoEnhance{
			Transcode: &Transcode{
				Container: &Container{
					Format: "mp4",
				},
				Video: &Video{
					Codec:   "H.264",
					Bitrate: "1000",
					Width:   "1280",
					Fps:     "30",
				},
				Audio: &Audio{
					Codec:      "aac",
					Bitrate:    "128",
					Samplerate: "44100",
					Channels:   "4",
				},
			},
			SuperResolution: &SuperResolution{
				Resolution:    "sdtohd",
				EnableScaleUp: "true",
				Version:       "Enhance",
			},
			ColorEnhance: &ColorEnhance{
				Contrast:   "50",
				Correction: "100",
				Saturation: "100",
			},
			MsSharpen: &MsSharpen{
				SharpenLevel: "5",
			},
			SDRtoHDR: &SDRtoHDR{
				HdrMode: "HDR10",
			},
			FrameEnhance: &FrameEnhance{
				FrameDoubling: "true",
			},
		},
	}

	_, _, err := client.CI.CreateVideoEnhanceTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateVideoEnhanceTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateVideoEnhanceTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>VideoEnhance</Tag><Name>VideoEnhance-test2</Name><VideoEnhance>" +
		"<Transcode><Container><Format>mp4</Format></Container>" +
		"<Video><Codec>H.264</Codec><Width>1280</Width><Fps>30</Fps><Bitrate>1000</Bitrate></Video>" +
		"<Audio><Codec>aac</Codec><Samplerate>44100</Samplerate><Bitrate>128</Bitrate><Channels>4</Channels></Audio></Transcode>" +
		"<SuperResolution><Resolution>sdtohd</Resolution><EnableScaleUp>true</EnableScaleUp><Version>Enhance</Version></SuperResolution>" +
		"<ColorEnhance><Contrast>50</Contrast><Correction>100</Correction><Saturation>100</Saturation></ColorEnhance>" +
		"<MsSharpen><SharpenLevel>5</SharpenLevel></MsSharpen><SDRtoHDR><HdrMode>HDR10</HdrMode></SDRtoHDR><FrameEnhance><FrameDoubling>true</FrameDoubling></FrameEnhance></VideoEnhance></Request>"
	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateVideoEnhanceTemplateOptions{
		Tag:  "VideoEnhance",
		Name: "VideoEnhance-test2",
		VideoEnhance: &VideoEnhance{
			Transcode: &Transcode{
				Container: &Container{
					Format: "mp4",
				},
				Video: &Video{
					Codec:   "H.264",
					Bitrate: "1000",
					Width:   "1280",
					Fps:     "30",
				},
				Audio: &Audio{
					Codec:      "aac",
					Bitrate:    "128",
					Samplerate: "44100",
					Channels:   "4",
				},
			},
			SuperResolution: &SuperResolution{
				Resolution:    "sdtohd",
				EnableScaleUp: "true",
				Version:       "Enhance",
			},
			ColorEnhance: &ColorEnhance{
				Contrast:   "50",
				Correction: "100",
				Saturation: "100",
			},
			MsSharpen: &MsSharpen{
				SharpenLevel: "5",
			},
			SDRtoHDR: &SDRtoHDR{
				HdrMode: "HDR10",
			},
			FrameEnhance: &FrameEnhance{
				FrameDoubling: "true",
			},
		},
	}

	_, _, err := client.CI.UpdateVideoEnhanceTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateVideoEnhanceTemplate returned error: %v", err)
	}
}

func TestCIService_CreateVideoTargetRecTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>VideoTargetRec</Tag><Name>VideoTargetRec-0</Name><VideoTargetRec><Body>true</Body><Pet>true</Pet><Car>true</Car></VideoTargetRec></Request>"
	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateVideoTargetRecTemplateOptions{
		Tag:  "VideoTargetRec",
		Name: "VideoTargetRec-0",
		VideoTargetRec: &VideoTargetRec{
			Body: "true",
			Pet:  "true",
			Car:  "true",
		},
	}

	_, _, err := client.CI.CreateVideoTargetRecTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateVideoTargetRecTemplate returned error: %v", err)
	}
}

func TestCIService_UpdateVideoTargetRecTemplate(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>VideoTargetRec</Tag><Name>VideoTargetRec-1</Name><VideoTargetRec><Body>true</Body><Pet>true</Pet><Car>true</Car></VideoTargetRec></Request>"

	tplId := "t1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateVideoTargetRecTemplateOptions{
		Tag:  "VideoTargetRec",
		Name: "VideoTargetRec-1",
		VideoTargetRec: &VideoTargetRec{
			Body: "true",
			Pet:  "true",
			Car:  "true",
		},
	}

	_, _, err := client.CI.UpdateVideoTargetRecTemplate(context.Background(), opt, tplId)
	if err != nil {
		t.Fatalf("CI.UpdateVideoTargetRecTemplate returned error: %v", err)
	}
}

func TestCIService_UnmarshalXML(t *testing.T) {
	responseBody := "<Request><MediaWorkflow><Name>workflow-1</Name><State>Active</State><Topology><Dependencies>" +
		"<Start>Transcode_1581665960537</Start><Transcode_1581665960537>Snapshot_1581665960536</Transcode_1581665960537>" +
		"<Snapshot_1581665960536>End</Snapshot_1581665960536></Dependencies><Nodes><Start><Type>Start</Type><Input><QueueId>p09d709939fef48a0a5c247ef39d90cec</QueueId>" +
		"<ObjectPrefix>wk-test</ObjectPrefix><ExtFilter><State>On</State><Custom>true</Custom><CustomExts>mp4</CustomExts></ExtFilter></Input></Start>" +
		"<Transcode_1581665960537><Type>Transcode</Type><Operation><TemplateId>t01e57db1c2d154d2fb57aa5de9313a897</TemplateId><Output><Region>ap-chongqing</Region>" +
		"<Bucket>test-123456789</Bucket><Object>trans1.mp4</Object></Output></Operation></Transcode_1581665960537><Snapshot_1581665960536><Type>Snapshot</Type>" +
		"<Operation><TemplateId>t07740e32081b44ad7a0aea03adcffd54a</TemplateId><Output><Region>ap-chongqing</Region><Bucket>test-123456789</Bucket>" +
		"<Object>snapshot-${number}.jpg</Object></Output></Operation></Snapshot_1581665960536></Nodes></Topology></MediaWorkflow></Request>"

	opt := &CreateMediaWorkflowOptions{
		MediaWorkflow: &MediaWorkflow{
			Topology: &Topology{},
		},
	}
	err := xml.Unmarshal([]byte(responseBody), opt)
	if err != nil {
		t.Fatalf("CI.UnmarshalXML returned error: %v", err)
	}

}
func TestCIService_CreateMediaWorkflow(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><MediaWorkflow><Name>workflow-1</Name><State>Active</State><Topology><Dependencies>" +
		"<Snapshot_1581665960536>End</Snapshot_1581665960536><Start>Transcode_1581665960537</Start><Transcode_1581665960537>Snapshot_1581665960536</Transcode_1581665960537>" +
		"</Dependencies><Nodes><Snapshot_1581665960536><Type>Snapshot</Type><Operation><TemplateId>t07740e32081b44ad7a0aea03adcffd54a</TemplateId><Output><Region>ap-chongqing</Region><Bucket>test-123456789</Bucket><Object>snapshot-${number}.jpg</Object></Output></Operation></Snapshot_1581665960536><Start><Type>Start</Type><Input><QueueId>p09d709939fef48a0a5c247ef39d90cec</QueueId><ObjectPrefix>wk-test</ObjectPrefix><ExtFilter><State>On</State><Custom>true</Custom><CustomExts>mp4</CustomExts></ExtFilter></Input></Start><Transcode_1581665960537><Type>Transcode</Type><Operation><TemplateId>t01e57db1c2d154d2fb57aa5de9313a897</TemplateId><Output><Region>ap-chongqing</Region><Bucket>test-123456789</Bucket><Object>trans1.mp4</Object></Output></Operation></Transcode_1581665960537></Nodes></Topology></MediaWorkflow></Request>"

	mux.HandleFunc("/workflow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaWorkflowOptions{
		MediaWorkflow: &MediaWorkflow{
			Name:  "workflow-1",
			State: "Active",
			Topology: &Topology{
				Dependencies: map[string]string{"Start": "Transcode_1581665960537", "Transcode_1581665960537": "Snapshot_1581665960536",
					"Snapshot_1581665960536": "End"},
				Nodes: map[string]Node{"Start": {Type: "Start", Input: &NodeInput{QueueId: "p09d709939fef48a0a5c247ef39d90cec",
					ObjectPrefix: "wk-test", ExtFilter: &ExtFilter{State: "On", Custom: "true", CustomExts: "mp4"}}},
					"Transcode_1581665960537": {Type: "Transcode", Operation: &NodeOperation{TemplateId: "t01e57db1c2d154d2fb57aa5de9313a897",
						Output: &NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "trans1.mp4"}}},
					"Snapshot_1581665960536": {Type: "Snapshot", Operation: &NodeOperation{TemplateId: "t07740e32081b44ad7a0aea03adcffd54a",
						Output: &NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "snapshot-${number}.jpg"}}},
				},
			},
		},
	}

	_, _, err := client.CI.CreateMediaWorkflow(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateMediaWorkflow returned error: %v", err)
	}
}

func TestCIService_UpdateMediaWorkflow(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><MediaWorkflow><Name>workflow-1</Name><State>Active</State><Topology><Dependencies>" +
		"<Snapshot_1581665960536>End</Snapshot_1581665960536><Start>Transcode_1581665960537</Start><Transcode_1581665960537>Snapshot_1581665960536</Transcode_1581665960537>" +
		"</Dependencies><Nodes><Snapshot_1581665960536><Type>Snapshot</Type><Operation><TemplateId>t07740e32081b44ad7a0aea03adcffd54a</TemplateId><Output><Region>ap-chongqing</Region><Bucket>test-123456789</Bucket><Object>snapshot-${number}.jpg</Object></Output></Operation></Snapshot_1581665960536><Start><Type>Start</Type><Input><QueueId>p09d709939fef48a0a5c247ef39d90cec</QueueId><ObjectPrefix>wk-test</ObjectPrefix><ExtFilter><State>On</State><Custom>true</Custom><CustomExts>mp4</CustomExts></ExtFilter></Input></Start><Transcode_1581665960537><Type>Transcode</Type><Operation><TemplateId>t01e57db1c2d154d2fb57aa5de9313a897</TemplateId><Output><Region>ap-chongqing</Region><Bucket>test-123456789</Bucket><Object>trans1.mp4</Object></Output></Operation></Transcode_1581665960537></Nodes></Topology></MediaWorkflow></Request>"

	wId := "w1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/workflow/"+wId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testBody(t, r, wantBody)
	})

	opt := &CreateMediaWorkflowOptions{
		MediaWorkflow: &MediaWorkflow{
			Name:  "workflow-1",
			State: "Active",
			Topology: &Topology{
				Dependencies: map[string]string{"Start": "Transcode_1581665960537", "Transcode_1581665960537": "Snapshot_1581665960536",
					"Snapshot_1581665960536": "End"},
				Nodes: map[string]Node{"Start": {Type: "Start", Input: &NodeInput{QueueId: "p09d709939fef48a0a5c247ef39d90cec",
					ObjectPrefix: "wk-test", ExtFilter: &ExtFilter{State: "On", Custom: "true", CustomExts: "mp4"}}},
					"Transcode_1581665960537": {Type: "Transcode", Operation: &NodeOperation{TemplateId: "t01e57db1c2d154d2fb57aa5de9313a897",
						Output: &NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "trans1.mp4"}}},
					"Snapshot_1581665960536": {Type: "Snapshot", Operation: &NodeOperation{TemplateId: "t07740e32081b44ad7a0aea03adcffd54a",
						Output: &NodeOutput{Region: "ap-chongqing", Bucket: "test-123456789", Object: "snapshot-${number}.jpg"}}},
				},
			},
		},
	}

	_, _, err := client.CI.UpdateMediaWorkflow(context.Background(), opt, wId)
	if err != nil {
		t.Fatalf("CI.UpdateMediaWorkflow returned error: %v", err)
	}
}

func TestCIService_ActiveMediaWorkflow(t *testing.T) {
	setup()
	defer teardown()

	wId := "w1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/workflow/"+wId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		v := values{
			"active": "",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.ActiveMediaWorkflow(context.Background(), wId)
	if err != nil {
		t.Fatalf("CI.ActiveMediaWorkflow returned error: %v", err)
	}
}

func TestCIService_PausedMediaWorkflow(t *testing.T) {
	setup()
	defer teardown()

	wId := "w1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/workflow/"+wId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		v := values{
			"paused": "",
		}
		testFormValues(t, r, v)
	})

	_, err := client.CI.PausedMediaWorkflow(context.Background(), wId)
	if err != nil {
		t.Fatalf("CI.PausedMediaWorkflow returned error: %v", err)
	}
}

func TestCIService_DescribeMediaWorkflow(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/workflow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"pageNumber": "2",
			"pageSize":   "5",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeMediaWorkflowOptions{
		Ids:        "",
		Name:       "",
		PageNumber: 2,
		PageSize:   5,
	}

	_, _, err := client.CI.DescribeMediaWorkflow(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeMediaWorkflow returned error: %v", err)
	}
}

func TestCIService_DeleteMediaWorkflow(t *testing.T) {
	setup()
	defer teardown()

	wId := "w1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/workflow/"+wId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, _, err := client.CI.DeleteMediaWorkflow(context.Background(), wId)
	if err != nil {
		t.Fatalf("CI.DeleteMediaWorkflow returned error: %v", err)
	}
}

func TestCIService_CreateInventoryTriggerJob(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Name>demo</Name><Type>Job</Type><Input><Prefix>input</Prefix></Input><Operation><TimeInterval><Start>2022-02-01T12:00:00+0800</Start>" +
		"<End>2022-05-01T12:00:00+0800</End></TimeInterval><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId><UserData>this is my inventorytriggerjob</UserData><JobLevel>1</JobLevel>" +
		"<CallBack>https://www.callback.com</CallBack><Tag>Transcode</Tag><JobParam><TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId></JobParam>" +
		"<Output><Region>ap-chongqing</Region><Bucket>test-1234567890</Bucket><Object>output/${InventoryTriggerJobId}/out.mp4</Object></Output>" +
		"</Operation></Request>"

	mux.HandleFunc("/inventorytriggerjob", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateInventoryTriggerJobOptions{
		Name: "demo",
		Type: "Job",
		Input: &InventoryTriggerJobInput{
			Prefix: "input",
		},
		Operation: &InventoryTriggerJobOperation{
			TimeInterval: InventoryTriggerJobOperationTimeInterval{
				Start: "2022-02-01T12:00:00+0800",
				End:   "2022-05-01T12:00:00+0800",
			},
			QueueId:  "p893bcda225bf4945a378da6662e81a89",
			UserData: "this is my inventorytriggerjob",
			JobLevel: 1,
			CallBack: "https://www.callback.com",
			Tag:      "Transcode",
			JobParam: &InventoryTriggerJobOperationJobParam{
				TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
			},
			Output: &JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1234567890",
				Object: "output/${InventoryTriggerJobId}/out.mp4",
			},
		},
	}

	_, _, err := client.CI.CreateInventoryTriggerJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateInventoryTriggerJob returned error: %v", err)
	}
}

func TestCIService_DescribeInventoryTriggerJob(t *testing.T) {
	setup()
	defer teardown()

	id := "b1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/inventorytriggerjob/"+id, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeInventoryTriggerJob(context.Background(), id)
	if err != nil {
		t.Fatalf("CI.DescribeInventoryTriggerJob returned error: %v", err)
	}
}

func TestCIService_DescribeInventoryTriggerJobs(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inventorytriggerjob", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"size":              "10",
			"type":              "Job",
			"orderByTime":       "Asc",
			"states":            "Running",
			"startCreationTime": "2022-02-25T12:00:00z",
			"endCreationTime":   "2022-02-28T12:00:00z",
			"workflowId":        "w123456789",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeInventoryTriggerJobsOptions{
		NextToken:         "",
		Size:              "10",
		Type:              "Job",
		OrderByTime:       "Asc",
		States:            "Running",
		StartCreationTime: "2022-02-25T12:00:00z",
		EndCreationTime:   "2022-02-28T12:00:00z",
		WorkflowId:        "w123456789",
		JobId:             "",
		Name:              "",
	}

	_, _, err := client.CI.DescribeInventoryTriggerJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeInventoryTriggerJobs returned error: %v", err)
	}
}

func TestCIService_CancelInventoryTriggerJob(t *testing.T) {
	setup()
	defer teardown()

	id := "b1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/inventorytriggerjob/"+id, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.CI.CancelInventoryTriggerJob(context.Background(), id)
	if err != nil {
		t.Fatalf("CI.CancelInventoryTriggerJob returned error: %v", err)
	}
}

func TestCIService_CreateImageSearchBucket(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><MaxCapacity>10000</MaxCapacity><MaxQps>10</MaxQps></Request>"
	mux.HandleFunc("/ImageSearchBucket", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, wantBody)
	})

	opt := &CreateImageSearchBucketOptions{
		MaxCapacity: "10000",
		MaxQps:      "10",
	}

	_, err := client.CI.CreateImageSearchBucket(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateImageSearchBucket returned error: %v", err)
	}
}

func TestCIService_AddImage(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><EntityId>car</EntityId><CustomContent>my car</CustomContent><Tags>{&#34;type&#34;: &#34;car&#34;}</Tags></Request>"
	mux.HandleFunc("/pic/car.jpg", func(w http.ResponseWriter, r *http.Request) {
		v := values{
			"ci-process": "ImageSearch",
			"action":     "AddImage",
		}
		testMethod(t, r, http.MethodPost)
		testFormValues(t, r, v)
		testBody(t, r, wantBody)
	})

	opt := &AddImageOptions{
		EntityId:      "car",
		CustomContent: "my car",
		Tags:          "{\"type\": \"car\"}",
	}

	_, err := client.CI.AddImage(context.Background(), "pic/car.jpg", opt)
	if err != nil {
		t.Fatalf("CI.AddImage returned error: %v", err)
	}
}

func TestCIService_ImageSearch(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/pic/car.jpg", func(w http.ResponseWriter, r *http.Request) {
		v := values{
			"ci-process":     "ImageSearch",
			"action":         "SearchImage",
			"MatchThreshold": "60",
			"Offset":         "1",
			"Limit":          "5",
			"Filter":         "{\"type\": \"car\"}",
		}
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, v)
	})

	opt := &ImageSearchOptions{
		MatchThreshold: 60,
		Offset:         1,
		Limit:          5,
		Filter:         "{\"type\": \"car\"}",
	}

	_, _, err := client.CI.ImageSearch(context.Background(), "pic/car.jpg", opt)
	if err != nil {
		t.Fatalf("CI.ImageSearch returned error: %v", err)
	}
}

func TestCIService_DelImage(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><EntityId>car</EntityId></Request>"
	mux.HandleFunc("/pic/car.jpg", func(w http.ResponseWriter, r *http.Request) {
		v := values{
			"ci-process": "ImageSearch",
			"action":     "DeleteImage",
		}
		testMethod(t, r, http.MethodPost)
		testFormValues(t, r, v)
		testBody(t, r, wantBody)
	})

	opt := &DelImageOptions{
		EntityId: "car",
	}

	_, err := client.CI.DelImage(context.Background(), "pic/car.jpg", opt)
	if err != nil {
		t.Fatalf("CI.DelImage returned error: %v", err)
	}
}

func TestCIService_CreateJob(t *testing.T) {
	setup()
	defer teardown()

	wantBody := "<Request><Tag>Animation</Tag><Input><Object>test.mp4</Object></Input>" +
		"<Operation><Output><Region>ap-beijing</Region><Bucket>abc-1250000000</Bucket>" +
		"<Object>test-trans.gif</Object></Output>" +
		"<TemplateId>t1460606b9752148c4ab182f55163ba7cd</TemplateId>" +
		"</Operation><QueueId>p893bcda225bf4945a378da6662e81a89</QueueId>" +
		"<CallBack>https://www.callback.com</CallBack></Request>"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/xml")
		testBody(t, r, wantBody)
	})

	opt := &CreateJobsOptions{
		Tag: "Animation",
		Input: &JobInput{
			Object: "test.mp4",
		},
		Operation: &MediaProcessJobOperation{
			Output: &JobOutput{
				Region: "ap-beijing",
				Bucket: "abc-1250000000",
				Object: "test-trans.gif",
			},
			TemplateId: "t1460606b9752148c4ab182f55163ba7cd",
		},
		QueueId:  "p893bcda225bf4945a378da6662e81a89",
		CallBack: "https://www.callback.com",
	}

	_, _, err := client.CI.CreateJob(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.CreateJob returned errors: %v", err)
	}
}

func TestCIService_CancelJob(t *testing.T) {
	setup()
	defer teardown()

	jobId := "j1460606b9752148c4ab182f55163ba7cd"
	mux.HandleFunc("/jobs/"+jobId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.CI.CancelJob(context.Background(), jobId)
	if err != nil {
		t.Fatalf("CI.CancelJob returned errors: %v", err)
	}
}

func TestCIService_DescribeJobs(t *testing.T) {
	setup()
	defer teardown()

	queueId := "aaaaaaaaaaa"
	tag := "Animation"

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"queueId": queueId,
			"tag":     tag,
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeJobsOptions{
		QueueId: queueId,
		Tag:     tag,
	}

	_, _, err := client.CI.DescribeJobs(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeJobs returned error: %v", err)
	}
}

func TestCIService_DescribeJob(t *testing.T) {
	setup()
	defer teardown()

	jobID := "jabcsdssfeipplsdfwe"
	mux.HandleFunc("/jobs/"+jobID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
	})

	_, _, err := client.CI.DescribeJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("CI.DescribeJobs returned error: %v", err)
	}
}

func TestCIService_ModifyM3U8Token(t *testing.T) {
	name := "test.m3u8"
	{
		setup()
		mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "modifym3u8token",
				"token":      "abc",
			}
			testFormValues(t, r, v)
		})

		opt := &ModifyM3U8TokenOptions{
			Token: "abc",
		}

		_, err := client.CI.ModifyM3U8Token(context.Background(), name, opt)
		if err != nil {
			t.Fatalf("CI.ModifyM3U8Token returned error: %v", err)
		}
		teardown()
	}
	{
		setup()
		mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "modifym3u8token",
				"token":      "abc",
				"versionId":  "1",
			}
			testFormValues(t, r, v)
		})

		opt := &ModifyM3U8TokenOptions{
			Token: "abc",
		}

		_, err := client.CI.ModifyM3U8Token(context.Background(), name, opt, "1")
		if err != nil {
			t.Fatalf("CI.ModifyM3U8Token returned error: %v", err)
		}
		teardown()
	}
	{
		setup()
		mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			testMethod(t, r, http.MethodGet)
			v := values{
				"ci-process": "modifym3u8token",
				"token":      "abc",
				"versionId":  "1",
			}
			testFormValues(t, r, v)
		})

		opt := &ModifyM3U8TokenOptions{
			Token: "abc",
		}

		_, err := client.CI.ModifyM3U8Token(context.Background(), name, opt, "1", "2")
		if err == nil || err.Error() != "wrong params" {
			t.Fatalf("CI.ModifyM3U8Token returned error: %v", err)
		}
		teardown()
	}
}

func TestCIService_DescribeTemplate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		v := values{
			"tag":        "All",
			"category":   "Custom",
			"ids":        "t123456798",
			"name":       "test",
			"pageNumber": "1",
			"pageSize":   "20",
		}
		testFormValues(t, r, v)
	})

	opt := &DescribeTemplateOptions{
		Tag:        "All",
		Category:   "Custom",
		Ids:        "t123456798",
		Name:       "test",
		PageNumber: 1,
		PageSize:   20,
	}

	_, _, err := client.CI.DescribeTemplate(context.Background(), opt)
	if err != nil {
		t.Fatalf("CI.DescribeTemplate returned error: %v", err)
	}
}

func TestCIService_DeleteTemplate(t *testing.T) {
	setup()
	defer teardown()

	tplId := "tc7c990a00bf211ed946af9e0691f2b7a"
	mux.HandleFunc("/template/"+tplId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, _, err := client.CI.DeleteTemplate(context.Background(), tplId)
	if err != nil {
		t.Fatalf("CI.DeleteTemplate returned error: %v", err)
	}
}
