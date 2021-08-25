package cos

// Basic imports
import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	//"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type CosTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int

	// CI client
	Client *cos.Client

	// Copy source client
	CClient *cos.Client

	Region string

	Bucket string

	Appid string

	// test_object
	TestObject string

	// special_file_name
	SepFileName string
}

// 请替换成您的账号及存储桶信息
const (
	//uin
	kUin   = "100010805041"
	kAppid = 1259654469

	// 常规测试需要的存储桶
	kBucket = "cosgosdktest-1259654469"
	kRegion = "ap-guangzhou"

	// 跨区域复制需要的目标存储桶，地域不能与kBucket存储桶相同, 目的存储桶需要开启多版本
	kRepBucket = "cosgosdkreptest"
	kRepRegion = "ap-chengdu"

	// Batch测试需要的源存储桶和目标存储桶，目前只在成都、重庆地域公测
	kBatchBucket       = "cosgosdktest-1259654469"
	kTargetBatchBucket = "cosgosdktest-1259654469" //复用了存储桶
	kBatchRegion       = "ap-guangzhou"
)

func (s *CosTestSuite) SetupSuite() {
	fmt.Println("Set up test")
	// init
	s.TestObject = "test.txt"
	s.SepFileName = "中文" + "→↓←→↖↗↙↘! \"#$%&'()*+,-./0123456789:;<=>@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

	// CI client for test interface
	// URL like this http://test-1253846586.cos.ap-guangzhou.myqcloud.com
	u := "https://" + kBucket + ".cos." + kRegion + ".myqcloud.com"
	u2 := "https://" + kUin + ".cos-control." + kBatchRegion + ".myqcloud.com"

	// Get the region
	bucketurl, _ := url.Parse(u)
	batchurl, _ := url.Parse(u2)
	p := strings.Split(bucketurl.Host, ".")
	assert.Equal(s.T(), 5, len(p), "Bucket host is not right")
	s.Region = p[2]

	// Bucket name
	pi := strings.LastIndex(p[0], "-")
	s.Bucket = p[0][:pi]
	s.Appid = p[0][pi+1:]

	ib := &cos.BaseURL{BucketURL: bucketurl, BatchURL: batchurl}
	s.Client = cos.NewClient(ib, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	opt := &cos.BucketPutOptions{
		XCosACL: "public-read",
	}
	r, err := s.Client.Bucket.Put(context.Background(), opt)
	if err != nil && r != nil && r.StatusCode == 409 {
		fmt.Println("BucketAlreadyOwnedByYou")
	} else if err != nil {
		assert.Nil(s.T(), err, "PutBucket Failed")
	}
}

// Begin of api test

// Service API
func (s *CosTestSuite) TestGetService() {
	_, _, err := s.Client.Service.Get(context.Background())
	assert.Nil(s.T(), err, "GetService Failed")
}

func (s *CosTestSuite) TestGetRegionService() {
	u, _ := url.Parse("http://cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{ServiceURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	_, _, err := client.Service.Get(context.Background())
	assert.Nil(s.T(), err, "GetService Failed")
}

// Bucket API
func (s *CosTestSuite) TestPutHeadDeleteBucket() {
	// Notic sometimes the bucket host can not analyis, may has i/o timeout problem
	u := "http://" + "testgosdkbucket-create-head-del-" + s.Appid + ".cos." + kRegion + ".myqcloud.com"
	iu, _ := url.Parse(u)
	ib := &cos.BaseURL{BucketURL: iu}
	client := cos.NewClient(ib, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})
	r, err := client.Bucket.Put(context.Background(), nil)
	if err != nil && r != nil && r.StatusCode == 409 {
		fmt.Println("BucketAlreadyOwnedByYou")
	} else if err != nil {
		assert.Nil(s.T(), err, "PutBucket Failed")
	}

	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)

	_, err = client.Bucket.Head(context.Background())
	assert.Nil(s.T(), err, "HeadBucket Failed")
	if err == nil {
		_, err = client.Bucket.Delete(context.Background())
		assert.Nil(s.T(), err, "DeleteBucket Failed")
	}
}

func (s *CosTestSuite) TestPutBucketACLIllegal() {
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "public-read-writ",
		},
	}
	_, err := s.Client.Bucket.PutACL(context.Background(), opt)
	assert.NotNil(s.T(), err, "PutBucketACL illegal Failed")
}

func (s *CosTestSuite) TestPutGetBucketACLNormal() {
	// with header
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	_, err := s.Client.Bucket.PutACL(context.Background(), opt)
	assert.Nil(s.T(), err, "PutBucketACL normal Failed")

	v, _, err := s.Client.Bucket.GetACL(context.Background())
	assert.Nil(s.T(), err, "GetBucketACL normal Failed")
	assert.Equal(s.T(), 1, len(v.AccessControlList), "GetBucketACL normal Failed, must be private")

}

func (s *CosTestSuite) TestGetBucket() {
	opt := &cos.BucketGetOptions{
		Prefix:  "中文",
		MaxKeys: 3,
	}
	_, _, err := s.Client.Bucket.Get(context.Background(), opt)
	assert.Nil(s.T(), err, "GetBucket Failed")
}

func (s *CosTestSuite) TestGetObjectVersions() {
	opt := &cos.BucketGetObjectVersionsOptions{
		Prefix:  "中文",
		MaxKeys: 3,
	}
	_, _, err := s.Client.Bucket.GetObjectVersions(context.Background(), opt)
	assert.Nil(s.T(), err, "GetObjectVersions Failed")
}

func (s *CosTestSuite) TestGetBucketLocation() {
	v, _, err := s.Client.Bucket.GetLocation(context.Background())
	assert.Nil(s.T(), err, "GetLocation Failed")
	assert.Equal(s.T(), s.Region, v.Location, "GetLocation wrong region")
}

func (s *CosTestSuite) TestPutGetDeleteCORS() {
	opt := &cos.BucketPutCORSOptions{
		Rules: []cos.BucketCORSRule{
			{
				AllowedOrigins: []string{"http://www.qq.com"},
				AllowedMethods: []string{"PUT", "GET"},
				AllowedHeaders: []string{"x-cos-meta-test", "x-cos-xx"},
				MaxAgeSeconds:  500,
				ExposeHeaders:  []string{"x-cos-meta-test1"},
			},
		},
	}
	_, err := s.Client.Bucket.PutCORS(context.Background(), opt)
	assert.Nil(s.T(), err, "PutBucketCORS Failed")

	v, _, err := s.Client.Bucket.GetCORS(context.Background())
	assert.Nil(s.T(), err, "GetBucketCORS Failed")
	assert.Equal(s.T(), 1, len(v.Rules), "GetBucketCORS wrong number rules")
}

func (s *CosTestSuite) TestVersionAndReplication() {
	opt := &cos.BucketPutVersionOptions{
		// Enabled or Suspended, the versioning once opened can not close.
		Status: "Enabled",
	}
	_, err := s.Client.Bucket.PutVersioning(context.Background(), opt)
	assert.Nil(s.T(), err, "PutVersioning Failed")
	time.Sleep(time.Second)
	v, _, err := s.Client.Bucket.GetVersioning(context.Background())
	assert.Nil(s.T(), err, "GetVersioning Failed")
	assert.Equal(s.T(), "Enabled", v.Status, "Get Wrong Version status")

	repOpt := &cos.PutBucketReplicationOptions{
		// qcs::cam::uin/[UIN]:uin/[Subaccount]
		Role: "qcs::cam::uin/" + kUin + ":uin/" + kUin,
		Rule: []cos.BucketReplicationRule{
			{
				ID: "1",
				// Enabled or Disabled
				Status: "Enabled",
				Destination: &cos.ReplicationDestination{
					// qcs::cos:[Region]::[Bucketname-Appid]
					Bucket: "qcs::cos:" + kRepRegion + "::" + kRepBucket + "-" + s.Appid,
				},
			},
		},
	}

	_, err = s.Client.Bucket.PutBucketReplication(context.Background(), repOpt)
	assert.Nil(s.T(), err, "PutBucketReplication Failed")
	time.Sleep(time.Second)
	vr, _, err := s.Client.Bucket.GetBucketReplication(context.Background())
	assert.Nil(s.T(), err, "GetBucketReplication Failed")
	for _, r := range vr.Rule {
		assert.Equal(s.T(), "Enabled", r.Status, "Get Wrong Version status")
		assert.Equal(s.T(), "qcs::cos:"+kRepRegion+"::"+kRepBucket+"-"+s.Appid, r.Destination.Bucket, "Get Wrong Version status")

	}
	_, err = s.Client.Bucket.DeleteBucketReplication(context.Background())
	assert.Nil(s.T(), err, "DeleteBucketReplication Failed")
}

func (s *CosTestSuite) TestBucketInventory() {
	id := "test1"
	dBucket := "qcs::cos:" + s.Region + "::" + s.Bucket + "-" + s.Appid
	opt := &cos.BucketPutInventoryOptions{
		ID: id,
		// True or False
		IsEnabled:              "True",
		IncludedObjectVersions: "All",
		Filter: &cos.BucketInventoryFilter{
			Prefix: "test",
		},
		OptionalFields: &cos.BucketInventoryOptionalFields{
			BucketInventoryFields: []string{
				"Size", "LastModifiedDate",
			},
		},
		Schedule: &cos.BucketInventorySchedule{
			// Weekly or Daily
			Frequency: "Daily",
		},
		Destination: &cos.BucketInventoryDestination{
			Bucket: dBucket,
			Format: "CSV",
		},
	}
	_, err := s.Client.Bucket.PutInventory(context.Background(), id, opt)
	assert.Nil(s.T(), err, "PutBucketInventory Failed")
	v, _, err := s.Client.Bucket.GetInventory(context.Background(), id)
	assert.Nil(s.T(), err, "GetBucketInventory Failed")
	assert.Equal(s.T(), "test1", v.ID, "Get Wrong inventory id")
	assert.Equal(s.T(), "true", v.IsEnabled, "Get Wrong inventory isenabled")
	assert.Equal(s.T(), dBucket, v.Destination.Bucket, "Get Wrong inventory isenabled")

	_, err = s.Client.Bucket.DeleteInventory(context.Background(), id)
	assert.Nil(s.T(), err, "DeleteBucketInventory Failed")
}

func (s *CosTestSuite) TestBucketLogging() {
	tBucket := s.Bucket + "-" + s.Appid
	opt := &cos.BucketPutLoggingOptions{
		LoggingEnabled: &cos.BucketLoggingEnabled{
			TargetBucket: tBucket,
		},
	}
	_, err := s.Client.Bucket.PutLogging(context.Background(), opt)
	assert.Nil(s.T(), err, "PutLogging Failed")
	v, _, err := s.Client.Bucket.GetLogging(context.Background())
	assert.Nil(s.T(), err, "GetLogging Failed")
	assert.Equal(s.T(), tBucket, v.LoggingEnabled.TargetBucket, "Get Wrong Version status")
}

func (s *CosTestSuite) TestBucketTagging() {
	opt := &cos.BucketPutTaggingOptions{
		TagSet: []cos.BucketTaggingTag{
			{
				Key:   "testk1",
				Value: "testv1",
			},
			{
				Key:   "testk2",
				Value: "testv2",
			},
		},
	}
	_, err := s.Client.Bucket.PutTagging(context.Background(), opt)
	assert.Nil(s.T(), err, "Put Tagging Failed")
	v, _, err := s.Client.Bucket.GetTagging(context.Background())
	assert.Nil(s.T(), err, "Get Tagging Failed")
	assert.Equal(s.T(), v.TagSet[0].Key, opt.TagSet[0].Key, "Get Wrong Tag key")
	assert.Equal(s.T(), v.TagSet[0].Value, opt.TagSet[0].Value, "Get Wrong Tag value")
	assert.Equal(s.T(), v.TagSet[1].Key, opt.TagSet[1].Key, "Get Wrong Tag key")
	assert.Equal(s.T(), v.TagSet[1].Value, opt.TagSet[1].Value, "Get Wrong Tag value")
}

func (s *CosTestSuite) TestPutGetDeleteLifeCycle() {
	lc := &cos.BucketPutLifecycleOptions{
		Rules: []cos.BucketLifecycleRule{
			{
				ID:     "1234",
				Filter: &cos.BucketLifecycleFilter{Prefix: "test"},
				Status: "Enabled",
				Transition: []cos.BucketLifecycleTransition{
					{
						Days:         10,
						StorageClass: "Standard",
					},
				},
			},
		},
	}
	_, err := s.Client.Bucket.PutLifecycle(context.Background(), lc)
	assert.Nil(s.T(), err, "PutBucketLifecycle Failed")
	_, r, err := s.Client.Bucket.GetLifecycle(context.Background())
	// Might cleaned by other case concrrent
	if err != nil && 404 != r.StatusCode {
		assert.Nil(s.T(), err, "GetBucketLifecycle Failed")
	}
	_, err = s.Client.Bucket.DeleteLifecycle(context.Background())
	assert.Nil(s.T(), err, "DeleteBucketLifecycle Failed")
}

func (s *CosTestSuite) TestPutGetDeleteWebsite() {
	opt := &cos.BucketPutWebsiteOptions{
		Index: "index.html",
		Error: &cos.ErrorDocument{"index_backup.html"},
		RoutingRules: &cos.WebsiteRoutingRules{
			[]cos.WebsiteRoutingRule{
				{
					ConditionErrorCode: "404",
					RedirectProtocol:   "https",
					RedirectReplaceKey: "404.html",
				},
				{
					ConditionPrefix:          "docs/",
					RedirectProtocol:         "https",
					RedirectReplaceKeyPrefix: "documents/",
				},
			},
		},
	}

	_, err := s.Client.Bucket.PutWebsite(context.Background(), opt)
	assert.Nil(s.T(), err, "PutBucketWebsite Failed")

	res, rsp, err := s.Client.Bucket.GetWebsite(context.Background())
	if err != nil && 404 != rsp.StatusCode {
		assert.Nil(s.T(), err, "GetBucketWebsite Failed")
	}
	assert.Equal(s.T(), opt.Index, res.Index, "GetBucketWebsite Failed")
	assert.Equal(s.T(), opt.Error, res.Error, "GetBucketWebsite Failed")
	assert.Equal(s.T(), opt.RedirectProtocol, res.RedirectProtocol, "GetBucketWebsite Failed")
	_, err = s.Client.Bucket.DeleteWebsite(context.Background())
	assert.Nil(s.T(), err, "DeleteBucketWebsite Failed")
}

func (s *CosTestSuite) TestListMultipartUploads() {
	// Create new upload
	name := "test_multipart" + time.Now().Format(time.RFC3339)
	flag := false
	v, _, err := s.Client.Object.InitiateMultipartUpload(context.Background(), name, nil)
	assert.Nil(s.T(), err, "InitiateMultipartUpload Failed")
	id := v.UploadID

	// List
	r, _, err := s.Client.Bucket.ListMultipartUploads(context.Background(), nil)
	assert.Nil(s.T(), err, "ListMultipartUploads Failed")
	for _, p := range r.Uploads {
		if p.Key == name {
			assert.Equal(s.T(), id, p.UploadID, "ListMultipartUploads wrong uploadid")
			flag = true
		}
	}
	assert.Equal(s.T(), true, flag, "ListMultipartUploads wrong key")

	// Abort
	_, err = s.Client.Object.AbortMultipartUpload(context.Background(), name, id)
	assert.Nil(s.T(), err, "AbortMultipartUpload Failed")
}

// Object API
func (s *CosTestSuite) TestPutHeadGetDeleteObject_10MB() {
	name := "test/objectPut" + time.Now().Format(time.RFC3339)
	b := make([]byte, 1024*1024*10)
	_, err := rand.Read(b)
	content := fmt.Sprintf("%X", b)
	f := strings.NewReader(content)

	_, err = s.Client.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	_, err = s.Client.Object.Head(context.Background(), name, nil)
	assert.Nil(s.T(), err, "HeadObject Failed")

	_, err = s.Client.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObjectByFile_10MB() {
	// Create tmp file
	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	assert.Nil(s.T(), err, "create tmp file Failed")
	defer newfile.Close()

	name := "test/objectPutByFile" + time.Now().Format(time.RFC3339)
	b := make([]byte, 1024*1024*10)
	_, err = rand.Read(b)

	newfile.Write(b)
	_, err = s.Client.Object.PutFromFile(context.Background(), name, filePath, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	// Over write tmp file
	_, err = s.Client.Object.GetToFile(context.Background(), name, filePath, nil)
	assert.Nil(s.T(), err, "HeadObject Failed")

	_, err = s.Client.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")

	// remove the local tmp file
	err = os.Remove(filePath)
	assert.Nil(s.T(), err, "remove local file Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObjectByUpload_10MB() {
	// Create tmp file
	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	assert.Nil(s.T(), err, "create tmp file Failed")
	defer newfile.Close()

	name := "test/objectUpload" + time.Now().Format(time.RFC3339)
	b := make([]byte, 1024*1024*10)
	_, err = rand.Read(b)

	newfile.Write(b)
	opt := &cos.MultiUploadOptions{
		PartSize:       1,
		ThreadPoolSize: 3,
	}
	_, _, err = s.Client.Object.Upload(context.Background(), name, filePath, opt)
	assert.Nil(s.T(), err, "PutObject Failed")

	// Over write tmp file
	_, err = s.Client.Object.GetToFile(context.Background(), name, filePath, nil)
	assert.Nil(s.T(), err, "HeadObject Failed")

	_, err = s.Client.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")

	// remove the local tmp file
	err = os.Remove(filePath)
	assert.Nil(s.T(), err, "remove local file Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObjectByUploadAndDownload_10MB() {
	// Create tmp file
	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newfile, err := os.Create(filePath)
	assert.Nil(s.T(), err, "create tmp file Failed")
	defer newfile.Close()

	name := "test/objectUpload" + time.Now().Format(time.RFC3339)
	b := make([]byte, 1024*1024*10)
	_, err = rand.Read(b)

	newfile.Write(b)
	opt := &cos.MultiUploadOptions{
		PartSize:       1,
		ThreadPoolSize: 3,
	}
	_, _, err = s.Client.Object.Upload(context.Background(), name, filePath, opt)
	assert.Nil(s.T(), err, "PutObject Failed")

	// Over write tmp file
	_, err = s.Client.Object.Download(context.Background(), name, filePath, nil)
	assert.Nil(s.T(), err, "DownloadObject Failed")

	_, err = s.Client.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")

	// remove the local tmp file
	err = os.Remove(filePath)
	assert.Nil(s.T(), err, "remove local file Failed")
}

func (s *CosTestSuite) TestPutGetDeleteObjectSpecialName() {
	f := strings.NewReader("test")
	name := s.SepFileName + time.Now().Format(time.RFC3339)
	_, err := s.Client.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	resp, err := s.Client.Object.Get(context.Background(), name, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(s.T(), "test", string(bs), "GetObject failed content wrong")

	_, err = s.Client.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutObjectToNonExistBucket() {
	u := "http://gosdknonexistbucket-" + s.Appid + ".cos." + s.Region + ".myqcloud.com"
	iu, _ := url.Parse(u)
	ib := &cos.BaseURL{BucketURL: iu}
	client := cos.NewClient(ib, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})
	name := "test/objectPut.go"
	f := strings.NewReader("test")
	r, err := client.Object.Put(context.Background(), name, f, nil)
	assert.NotNil(s.T(), err, "PutObject ToNonExistBucket Failed")
	assert.Equal(s.T(), 404, r.StatusCode, "PutObject ToNonExistBucket, not 404")
}

func (s *CosTestSuite) TestPutGetObjectACL() {
	name := "test/objectACL.go" + time.Now().Format(time.RFC3339)
	f := strings.NewReader("test")
	_, err := s.Client.Object.Put(context.Background(), name, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	// Put acl
	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	_, err = s.Client.Object.PutACL(context.Background(), name, opt)
	assert.Nil(s.T(), err, "PutObjectACL Failed")
	v, _, err := s.Client.Object.GetACL(context.Background(), name)
	assert.Nil(s.T(), err, "GetObjectACL Failed")
	assert.Equal(s.T(), 2, len(v.AccessControlList), "GetLifecycle wrong number rules")

	_, err = s.Client.Object.Delete(context.Background(), name)
	assert.Nil(s.T(), err, "DeleteObject Failed")
}

func (s *CosTestSuite) TestPutObjectRestore() {
	name := "archivetest"
	putOpt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosStorageClass: "ARCHIVE",
		},
	}
	f := strings.NewReader("test")
	_, err := s.Client.Object.Put(context.Background(), name, f, putOpt)
	assert.Nil(s.T(), err, "PutObject Archive faild")
	opt := &cos.ObjectRestoreOptions{
		Days: 2,
		Tier: &cos.CASJobParameters{
			// Standard, Exepdited and Bulk
			Tier: "Expedited",
		},
	}
	resp, _ := s.Client.Object.PostRestore(context.Background(), name, opt)
	retCode := resp.StatusCode
	if retCode != 200 && retCode != 202 && retCode != 409 {
		right := false
		fmt.Println("PutObjectRestore get code is:", retCode)
		assert.Equal(s.T(), true, right, "PutObjectRestore Failed")
	}

}

func (s *CosTestSuite) TestCopyObject() {
	u := "http://" + kRepBucket + "-" + s.Appid + ".cos." + kRepRegion + ".myqcloud.com"
	iu, _ := url.Parse(u)
	ib := &cos.BaseURL{BucketURL: iu}
	c := cos.NewClient(ib, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	opt := &cos.BucketPutOptions{
		XCosACL: "public-read",
	}

	// Notice in intranet the bucket host sometimes has i/o timeout problem
	r, err := c.Bucket.Put(context.Background(), opt)
	if err != nil && r != nil && r.StatusCode == 409 {
		fmt.Println("BucketAlreadyOwnedByYou")
	} else if err != nil {
		assert.Nil(s.T(), err, "PutBucket Failed")
	}

	source := "test/objectMove1" + time.Now().Format(time.RFC3339)
	expected := "test"
	f := strings.NewReader(expected)

	r, err = c.Object.Put(context.Background(), source, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")
	var version_id string
	if r.Header["X-Cos-Version-Id"] != nil {
		version_id = r.Header.Get("X-Cos-Version-Id")
	}

	time.Sleep(3 * time.Second)
	// Copy file
	soruceURL := fmt.Sprintf("%s/%s", iu.Host, source)
	dest := "test/objectMove1" + time.Now().Format(time.RFC3339)
	//opt := &cos.ObjectCopyOptions{}
	if version_id == "" {
		_, _, err = s.Client.Object.Copy(context.Background(), dest, soruceURL, nil)
	} else {
		_, _, err = s.Client.Object.Copy(context.Background(), dest, soruceURL, nil, version_id)
	}
	assert.Nil(s.T(), err, "PutObjectCopy Failed")

	// Check content
	resp, err := s.Client.Object.Get(context.Background(), dest, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	result := string(bs)
	assert.Equal(s.T(), expected, result, "PutObjectCopy Failed, wrong content")
}

func (s *CosTestSuite) TestCreateAbortMultipartUpload() {
	name := "test_multipart" + time.Now().Format(time.RFC3339)
	v, _, err := s.Client.Object.InitiateMultipartUpload(context.Background(), name, nil)
	assert.Nil(s.T(), err, "InitiateMultipartUpload Failed")

	_, err = s.Client.Object.AbortMultipartUpload(context.Background(), name, v.UploadID)
	assert.Nil(s.T(), err, "AbortMultipartUpload Failed")
}

func (s *CosTestSuite) TestCreateCompleteMultipartUpload() {
	name := "test/test_complete_upload" + time.Now().Format(time.RFC3339)
	v, _, err := s.Client.Object.InitiateMultipartUpload(context.Background(), name, nil)
	uploadID := v.UploadID
	blockSize := 1024 * 1024 * 3

	opt := &cos.CompleteMultipartUploadOptions{}
	for i := 1; i < 3; i++ {
		b := make([]byte, blockSize)
		_, err := rand.Read(b)
		content := fmt.Sprintf("%X", b)
		f := strings.NewReader(content)

		resp, err := s.Client.Object.UploadPart(
			context.Background(), name, uploadID, i, f, nil,
		)
		assert.Nil(s.T(), err, "UploadPart Failed")
		etag := resp.Header.Get("Etag")
		opt.Parts = append(opt.Parts, cos.Object{
			PartNumber: i, ETag: etag},
		)
	}

	_, _, err = s.Client.Object.CompleteMultipartUpload(
		context.Background(), name, uploadID, opt,
	)

	assert.Nil(s.T(), err, "CompleteMultipartUpload Failed")
}

func (s *CosTestSuite) TestSSE_C() {
	name := "test/TestSSE_C"
	content := "test sse-c " + time.Now().Format(time.RFC3339)
	f := strings.NewReader(content)
	putOpt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
			//XCosServerSideEncryption: "AES256",
			XCosSSECustomerAglo:   "AES256",
			XCosSSECustomerKey:    "MDEyMzQ1Njc4OUFCQ0RFRjAxMjM0NTY3ODlBQkNERUY=",
			XCosSSECustomerKeyMD5: "U5L61r7jcwdNvT7frmUG8g==",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
			//XCosACL: "private",
		},
	}
	_, err := s.Client.Object.Put(context.Background(), name, f, putOpt)
	assert.Nil(s.T(), err, "PutObject with SSE failed")

	headOpt := &cos.ObjectHeadOptions{
		XCosSSECustomerAglo:   "AES256",
		XCosSSECustomerKey:    "MDEyMzQ1Njc4OUFCQ0RFRjAxMjM0NTY3ODlBQkNERUY=",
		XCosSSECustomerKeyMD5: "U5L61r7jcwdNvT7frmUG8g==",
	}
	_, err = s.Client.Object.Head(context.Background(), name, headOpt)
	assert.Nil(s.T(), err, "HeadObject with SSE failed")

	getOpt := &cos.ObjectGetOptions{
		XCosSSECustomerAglo:   "AES256",
		XCosSSECustomerKey:    "MDEyMzQ1Njc4OUFCQ0RFRjAxMjM0NTY3ODlBQkNERUY=",
		XCosSSECustomerKeyMD5: "U5L61r7jcwdNvT7frmUG8g==",
	}
	var resp *cos.Response
	resp, err = s.Client.Object.Get(context.Background(), name, getOpt)
	assert.Nil(s.T(), err, "GetObject with SSE failed")

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyContent := string(bodyBytes)
	assert.Equal(s.T(), content, bodyContent, "GetObject with SSE failed, want: %+v, res: %+v", content, bodyContent)

	copyOpt := &cos.ObjectCopyOptions{
		&cos.ObjectCopyHeaderOptions{
			XCosCopySourceSSECustomerAglo:   "AES256",
			XCosCopySourceSSECustomerKey:    "MDEyMzQ1Njc4OUFCQ0RFRjAxMjM0NTY3ODlBQkNERUY=",
			XCosCopySourceSSECustomerKeyMD5: "U5L61r7jcwdNvT7frmUG8g==",
		},
		&cos.ACLHeaderOptions{},
	}
	copySource := s.Bucket + "-" + s.Appid + ".cos." + s.Region + ".myqcloud.com/" + name
	_, _, err = s.Client.Object.Copy(context.Background(), "test/TestSSE_C_Copy", copySource, copyOpt)
	assert.Nil(s.T(), err, "CopyObject with SSE failed")

	partIni := &cos.MultiUploadOptions{
		OptIni: &cos.InitiateMultipartUploadOptions{
			&cos.ACLHeaderOptions{},
			&cos.ObjectPutHeaderOptions{
				XCosSSECustomerAglo:   "AES256",
				XCosSSECustomerKey:    "MDEyMzQ1Njc4OUFCQ0RFRjAxMjM0NTY3ODlBQkNERUY=",
				XCosSSECustomerKeyMD5: "U5L61r7jcwdNvT7frmUG8g==",
			},
		},
		PartSize: 1,
	}
	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newFile, err := os.Create(filePath)
	assert.Nil(s.T(), err, "create tmp file Failed")
	defer newFile.Close()
	b := make([]byte, 1024*10)
	_, err = rand.Read(b)
	newFile.Write(b)

	_, _, err = s.Client.Object.MultiUpload(context.Background(), "test/TestSSE_C_MultiUpload", filePath, partIni)
	assert.Nil(s.T(), err, "MultiUpload with SSE failed")

	err = os.Remove(filePath)
	assert.Nil(s.T(), err, "remove local file Failed")
}

func (s *CosTestSuite) TestMultiUpload() {
	filePath := "tmpfile" + time.Now().Format(time.RFC3339)
	newFile, err := os.Create(filePath)
	assert.Nil(s.T(), err, "create tmp file Failed")
	defer newFile.Close()
	b := make([]byte, 1024*1024*10)
	_, err = rand.Read(b)
	newFile.Write(b)

	partIni := &cos.MultiUploadOptions{}

	_, _, err = s.Client.Object.MultiUpload(context.Background(), "test/Test_MultiUpload", filePath, partIni)

	err = os.Remove(filePath)
	assert.Nil(s.T(), err, "remove tmp file failed")
}

func (s *CosTestSuite) TestAppend() {
	name := "append" + time.Now().Format(time.RFC3339)
	b1 := make([]byte, 1024*1024*10)
	_, err := rand.Read(b1)
	pos, _, err := s.Client.Object.Append(context.Background(), name, 0, bytes.NewReader(b1), nil)
	assert.Nil(s.T(), err, "append object failed")
	assert.Equal(s.T(), len(b1), pos, "append object pos error")

	b2 := make([]byte, 12345)
	rand.Read(b2)
	pos, _, err = s.Client.Object.Append(context.Background(), name, pos, bytes.NewReader(b2), nil)
	assert.Nil(s.T(), err, "append object failed")
	assert.Equal(s.T(), len(b1)+len(b2), pos, "append object pos error")
}

/*
func (s *CosTestSuite) TestBatch() {
	client := cos.NewClient(s.Client.BaseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	source_name := "test/1.txt"
	sf := strings.NewReader("batch test content")
	_, err := client.Object.Put(context.Background(), source_name, sf, nil)
	assert.Nil(s.T(), err, "object put Failed")

	manifest_name := "test/manifest.csv"
	f := strings.NewReader(kBatchBucket + "," + source_name)
	resp, err := client.Object.Put(context.Background(), manifest_name, f, nil)
	assert.Nil(s.T(), err, "object put Failed")
	etag := resp.Header.Get("ETag")

	uuid_str := uuid.New().String()
	opt := &cos.BatchCreateJobOptions{
		ClientRequestToken:   uuid_str,
		ConfirmationRequired: "true",
		Description:          "test batch",
		Manifest: &cos.BatchJobManifest{
			Location: &cos.BatchJobManifestLocation{
				ETag:      etag,
				ObjectArn: "qcs::cos:" + kBatchRegion + ":uid/" + s.Appid + ":" + kBatchBucket + "/" + manifest_name,
			},
			Spec: &cos.BatchJobManifestSpec{
				Fields: []string{"Bucket", "Key"},
				Format: "COSBatchOperations_CSV_V1",
			},
		},
		Operation: &cos.BatchJobOperation{
			PutObjectCopy: &cos.BatchJobOperationCopy{
				TargetResource: "qcs::cos:" + kBatchRegion + ":uid/" + s.Appid + ":" + kTargetBatchBucket,
			},
		},
		Priority: 1,
		Report: &cos.BatchJobReport{
			Bucket:      "qcs::cos:" + kBatchRegion + ":uid/" + s.Appid + ":" + kBatchBucket,
			Enabled:     "true",
			Format:      "Report_CSV_V1",
			Prefix:      "job-result",
			ReportScope: "AllTasks",
		},
		RoleArn: "qcs::cam::uin/" + kUin + ":roleName/COSBatch_QcsRole",
	}
	headers := &cos.BatchRequestHeaders{
		XCosAppid: kAppid,
	}

	res1, _, err := client.Batch.CreateJob(context.Background(), opt, headers)
	assert.Nil(s.T(), err, "create job Failed")

	jobid := res1.JobId

	res2, _, err := client.Batch.DescribeJob(context.Background(), jobid, headers)
	assert.Nil(s.T(), err, "describe job Failed")
	assert.Equal(s.T(), res2.Job.ConfirmationRequired, "true", "ConfirmationRequired not right")
	assert.Equal(s.T(), res2.Job.Description, "test batch", "Description not right")
	assert.Equal(s.T(), res2.Job.JobId, jobid, "jobid not right")
	assert.Equal(s.T(), res2.Job.Priority, 1, "priority not right")
	assert.Equal(s.T(), res2.Job.RoleArn, "qcs::cam::uin/"+kUin+":roleName/COSBatch_QcsRole", "priority not right")

	_, _, err = client.Batch.ListJobs(context.Background(), nil, headers)
	assert.Nil(s.T(), err, "list jobs failed")

	up_opt := &cos.BatchUpdatePriorityOptions{
		JobId:    jobid,
		Priority: 3,
	}
	res3, _, err := client.Batch.UpdateJobPriority(context.Background(), up_opt, headers)
	assert.Nil(s.T(), err, "list jobs failed")
	assert.Equal(s.T(), res3.JobId, jobid, "jobid failed")
	assert.Equal(s.T(), res3.Priority, 3, "priority not right")

	// 等待状态变成Suspended
	for i := 0; i < 50; i = i + 1 {
		res, _, err := client.Batch.DescribeJob(context.Background(), jobid, headers)
		assert.Nil(s.T(), err, "describe job Failed")
		assert.Equal(s.T(), res2.Job.ConfirmationRequired, "true", "ConfirmationRequired not right")
		assert.Equal(s.T(), res2.Job.Description, "test batch", "Description not right")
		assert.Equal(s.T(), res2.Job.JobId, jobid, "jobid not right")
		assert.Equal(s.T(), res2.Job.Priority, 1, "priority not right")
		assert.Equal(s.T(), res2.Job.RoleArn, "qcs::cam::uin/"+kUin+":roleName/COSBatch_QcsRole", "priority not right")
		if res.Job.Status == "Suspended" {
			break
		}
		if i == 9 {
			assert.Error(s.T(), errors.New("Job status is not Suspended or timeout"))
		}
		time.Sleep(time.Second * 2)
	}
	us_opt := &cos.BatchUpdateStatusOptions{
		JobId:              jobid,
		RequestedJobStatus: "Ready", // 允许状态转换见 https://cloud.tencent.com/document/product/436/38604
		StatusUpdateReason: "to test",
	}
	res4, _, err := client.Batch.UpdateJobStatus(context.Background(), us_opt, headers)
	assert.Nil(s.T(), err, "list jobs failed")
	assert.Equal(s.T(), res4.JobId, jobid, "jobid failed")
	assert.Equal(s.T(), res4.Status, "Ready", "status failed")
	assert.Equal(s.T(), res4.StatusUpdateReason, "to test", "StatusUpdateReason failed")
}
*/

func (s *CosTestSuite) TestEncryption() {
	opt := &cos.BucketPutEncryptionOptions{
		Rule: &cos.BucketEncryptionConfiguration{
			SSEAlgorithm: "AES256",
		},
	}

	_, err := s.Client.Bucket.PutEncryption(context.Background(), opt)
	assert.Nil(s.T(), err, "PutEncryption Failed")

	time.Sleep(time.Second * 2)
	res, _, err := s.Client.Bucket.GetEncryption(context.Background())
	assert.Nil(s.T(), err, "GetEncryption Failed")
	assert.Equal(s.T(), opt.Rule.SSEAlgorithm, res.Rule.SSEAlgorithm, "GetEncryption Failed")

	_, err = s.Client.Bucket.DeleteEncryption(context.Background())
	assert.Nil(s.T(), err, "DeleteEncryption Failed")
}

func (s *CosTestSuite) TestReferer() {
	opt := &cos.BucketPutRefererOptions{
		Status:      "Enabled",
		RefererType: "White-List",
		DomainList: []string{
			"*.qq.com",
			"*.qcloud.com",
		},
		EmptyReferConfiguration: "Allow",
	}

	_, err := s.Client.Bucket.PutReferer(context.Background(), opt)
	assert.Nil(s.T(), err, "PutReferer Failed")

	res, _, err := s.Client.Bucket.GetReferer(context.Background())
	assert.Nil(s.T(), err, "GetReferer Failed")
	assert.Equal(s.T(), opt.Status, res.Status, "GetReferer Failed")
	assert.Equal(s.T(), opt.RefererType, res.RefererType, "GetReferer Failed")
	assert.Equal(s.T(), opt.DomainList, res.DomainList, "GetReferer Failed")
	assert.Equal(s.T(), opt.EmptyReferConfiguration, res.EmptyReferConfiguration, "GetReferer Failed")
}

func (s *CosTestSuite) TestAccelerate() {
	opt := &cos.BucketPutAccelerateOptions{
		Status: "Enabled",
		Type:   "COS",
	}
	_, err := s.Client.Bucket.PutAccelerate(context.Background(), opt)
	assert.Nil(s.T(), err, "PutAccelerate Failed")

	time.Sleep(time.Second)
	res, _, err := s.Client.Bucket.GetAccelerate(context.Background())
	assert.Nil(s.T(), err, "GetAccelerate Failed")
	assert.Equal(s.T(), opt.Status, res.Status, "GetAccelerate Failed")
	assert.Equal(s.T(), opt.Type, res.Type, "GetAccelerate Failed")

	opt.Status = "Suspended"
	_, err = s.Client.Bucket.PutAccelerate(context.Background(), opt)
	assert.Nil(s.T(), err, "PutAccelerate Failed")

	time.Sleep(time.Second)
	res, _, err = s.Client.Bucket.GetAccelerate(context.Background())
	assert.Nil(s.T(), err, "GetAccelerate Failed")
	assert.Equal(s.T(), opt.Status, res.Status, "GetAccelerate Failed")
	assert.Equal(s.T(), opt.Type, res.Type, "GetAccelerate Failed")
}

func (s *CosTestSuite) TestMultiCopy() {
	u := "http://" + kRepBucket + "-" + s.Appid + ".cos." + kRepRegion + ".myqcloud.com"
	iu, _ := url.Parse(u)
	ib := &cos.BaseURL{BucketURL: iu}
	c := cos.NewClient(ib, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	opt := &cos.BucketPutOptions{
		XCosACL: "public-read",
	}

	// Notice in intranet the bucket host sometimes has i/o timeout problem
	r, err := c.Bucket.Put(context.Background(), opt)
	if err != nil && r != nil && r.StatusCode == 409 {
		fmt.Println("BucketAlreadyOwnedByYou")
	} else if err != nil {
		assert.Nil(s.T(), err, "PutBucket Failed")
	}

	source := "test/objectMove1" + time.Now().Format(time.RFC3339)
	expected := "test"
	f := strings.NewReader(expected)

	r, err = c.Object.Put(context.Background(), source, f, nil)
	assert.Nil(s.T(), err, "PutObject Failed")

	time.Sleep(3 * time.Second)
	// Copy file
	soruceURL := fmt.Sprintf("%s/%s", iu.Host, source)
	dest := "test/objectMove1" + time.Now().Format(time.RFC3339)
	_, _, err = s.Client.Object.MultiCopy(context.Background(), dest, soruceURL, nil)
	assert.Nil(s.T(), err, "MultiCopy Failed")

	// Check content
	resp, err := s.Client.Object.Get(context.Background(), dest, nil)
	assert.Nil(s.T(), err, "GetObject Failed")
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	result := string(bs)
	assert.Equal(s.T(), expected, result, "MultiCopy Failed, wrong content")
}

// End of api test

// All methods that begin with "Test" are run as tests within a
// suite.
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCosTestSuite(t *testing.T) {
	suite.Run(t, new(CosTestSuite))
}

func (s *CosTestSuite) TearDownSuite() {
	// Clean the file in bucket
	// r, _, err := s.Client.Bucket.ListMultipartUploads(context.Background(), nil)
	// assert.Nil(s.T(), err, "ListMultipartUploads Failed")
	// for _, p := range r.Uploads {
	// 	// Abort
	// 	_, err = s.Client.Object.AbortMultipartUpload(context.Background(), p.Key, p.UploadID)
	// 	assert.Nil(s.T(), err, "AbortMultipartUpload Failed")
	// }

	// // Delete objects
	// opt := &cos.BucketGetOptions{
	// 	MaxKeys: 500,
	// }
	// v, _, err := s.Client.Bucket.Get(context.Background(), opt)
	// assert.Nil(s.T(), err, "GetBucket Failed")
	// for _, c := range v.Contents {
	// 	_, err := s.Client.Object.Delete(context.Background(), c.Key)
	// 	assert.Nil(s.T(), err, "DeleteObject Failed")
	// }

	// When clean up these infos, can not solve the concurrent test problem

	fmt.Println("tear down~")

}
