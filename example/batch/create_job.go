package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func main() {
	test_batch_bucket := "testcd-1259654469"
	target_batch_bucket := "cosgosdkreptest-1259654469"
	appid := 1259654469
	uin := "100010805041"
	region := "ap-chengdu"

	// bucket url：<Bucketname-Appid>.cos.<region>.mycloud.com
	bucketurl, _ := url.Parse("https://" + test_batch_bucket + ".cos." + region + ".myqcloud.com")
	// batch url：<uin>.cos-control.<region>.myqcloud.ccom
	batchurl, _ := url.Parse("https://" + uin + ".cos-control." + region + ".myqcloud.com")

	b := &cos.BaseURL{BucketURL: bucketurl, BatchURL: batchurl}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	// 创建需要复制的文件
	source_name := "test/1.txt"
	sf := strings.NewReader("batch test content")
	_, err := c.Object.Put(context.Background(), source_name, sf, nil)
	if err != nil {
		panic(err)
	}

	// 创建清单文件
	manifest_name := "test/manifest.csv"
	f := strings.NewReader(test_batch_bucket + "," + source_name)
	resp, err := c.Object.Put(context.Background(), manifest_name, f, nil)
	if err != nil {
		panic(err)
	}
	etag := resp.Header.Get("ETag")

	uuid_str := uuid.New().String()
	opt := &cos.BatchCreateJobOptions{
		ClientRequestToken:   uuid_str,
		ConfirmationRequired: "true",
		Description:          "test batch",
		Manifest: &cos.BatchJobManifest{
			Location: &cos.BatchJobManifestLocation{
				ETag:      etag,
				ObjectArn: "qcs::cos:" + region + "::" + test_batch_bucket + "/" + manifest_name,
			},
			Spec: &cos.BatchJobManifestSpec{
				Fields: []string{"Bucket", "Key"},
				Format: "COSBatchOperations_CSV_V1",
			},
		},
		Operation: &cos.BatchJobOperation{
			PutObjectCopy: &cos.BatchJobOperationCopy{
				TargetResource: "qcs::cos:" + region + ":uid/" + strconv.Itoa(appid) + ":" + target_batch_bucket,
			},
		},
		Priority: 1,
		Report: &cos.BatchJobReport{
			Bucket:      "qcs::cos:" + region + "::" + test_batch_bucket,
			Enabled:     "true",
			Format:      "Report_CSV_V1",
			Prefix:      "job-result",
			ReportScope: "AllTasks",
		},
		RoleArn: "qcs::cam::uin/" + uin + ":roleName/COSBatch_QcsRole",
	}
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}

	res, _, err := c.Batch.CreateJob(context.Background(), opt, headers)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
