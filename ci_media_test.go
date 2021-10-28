package cos

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

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

func TestCIService_GetMediaInfo(t *testing.T) {
	setup()
	defer teardown()

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
}

func TestCIService_GetSnapshot(t *testing.T) {
	setup()
	defer teardown()

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
}
