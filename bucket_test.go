package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_Get(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketGetOptions{
		Prefix:  "test",
		MaxKeys: 2,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"prefix":   "test",
			"max-keys": "2",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<ListBucketResult>
	<Name>test-1253846586</Name>
	<Prefix>test</Prefix>
	<Marker/>
	<MaxKeys>2</MaxKeys>
	<IsTruncated>true</IsTruncated>
	<NextMarker>test/delete.txt</NextMarker>
	<Contents>
		<Key>test/</Key>
		<LastModified>2017-06-09T16:32:25.000Z</LastModified>
		<ETag>&quot;&quot;</ETag>
		<Size>0</Size>
		<Owner>
			<ID>1253846586</ID>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
	</Contents>
	<Contents>
		<Key>test/anonymous_get.go</Key>
		<LastModified>2017-06-17T15:09:26.000Z</LastModified>
		<ETag>&quot;5b7236085f08b3818bfa40b03c946dcc&quot;</ETag>
		<Size>8</Size>
		<Owner>
			<ID>1253846586</ID>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
	</Contents>
</ListBucketResult>`)
	})

	ref, _, err := client.Bucket.Get(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Get returned error: %v", err)
	}

	want := &BucketGetResult{
		XMLName:     xml.Name{Local: "ListBucketResult"},
		Name:        "test-1253846586",
		Prefix:      "test",
		MaxKeys:     2,
		IsTruncated: true,
		NextMarker:  "test/delete.txt",
		Contents: []Object{
			{
				Key:          "test/",
				LastModified: "2017-06-09T16:32:25.000Z",
				ETag:         "\"\"",
				Size:         0,
				Owner: &Owner{
					ID: "1253846586",
				},
				StorageClass: "STANDARD",
			},
			{
				Key:          "test/anonymous_get.go",
				LastModified: "2017-06-17T15:09:26.000Z",
				ETag:         "\"5b7236085f08b3818bfa40b03c946dcc\"",
				Size:         8,
				Owner: &Owner{
					ID: "1253846586",
				},
				StorageClass: "STANDARD",
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.Get returned %+v, want %+v", ref, want)
	}
}

func TestBucketService_Put(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutOptions{
		XCosACL: "public-read",
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		testHeader(t, r, "x-cos-acl", "public-read")
	})

	_, err := client.Bucket.Put(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Put returned error: %v", err)
	}

	opt = &BucketPutOptions{
		XCosACL: "public-read",
		CreateBucketConfiguration: &CreateBucketConfiguration{
			BucketAZConfig: "MAZ",
		},
	}
	_, err = client.Bucket.Put(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Put returned error: %v", err)
	}
}

func TestBucketService_Delete(t *testing.T) {
	setup()
	defer teardown()

	var checkHeader bool
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		if checkHeader {
			testHeader(t, r, "x-cos-meta-test", "test")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.Delete(context.Background())
	if err != nil {
		t.Fatalf("Bucket.Delete returned error: %v", err)
	}

	checkHeader = true
	opt := &BucketDeleteOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "test")
	_, err = client.Bucket.Delete(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Delete returned error: %v", err)
	}

}

func TestBucketService_Head(t *testing.T) {
	setup()
	defer teardown()

	var checkHeader bool
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodHead)
		if checkHeader {
			testHeader(t, r, "x-cos-meta-test", "test")
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := client.Bucket.Head(context.Background())
	if err != nil {
		t.Fatalf("Bucket.Head returned error: %v", err)
	}

	checkHeader = true
	opt := &BucketHeadOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "test")
	_, err = client.Bucket.Head(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.Head returned error: %v", err)
	}

}

func TestBucketService_IsExist(t *testing.T) {
	setup()
	defer teardown()

	var exist, notfound, deny bool
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodHead)
		if exist {
			w.WriteHeader(http.StatusOK)
		}
		if notfound {
			w.WriteHeader(http.StatusNotFound)
		}
		if deny {
			w.WriteHeader(http.StatusForbidden)
		}
	})

	// 存在
	exist, notfound, deny = true, false, false
	isExisted, err := client.Bucket.IsExist(context.Background())
	if err != nil {
		t.Fatalf("Bucket.Head returned error: %v", err)
	}
	if isExisted != true {
		t.Errorf("bucket IsExist failed")
	}

	// 不存在
	exist, notfound, deny = false, true, false
	isExisted, err = client.Bucket.IsExist(context.Background())
	if err != nil {
		t.Fatalf("Bucket.Head returned error: %v", err)
	}
	if isExisted != false {
		t.Errorf("bucket IsExist failed")
	}

	// 报错
	exist, notfound, deny = false, false, true
	isExisted, err = client.Bucket.IsExist(context.Background())
	if err == nil {
		t.Fatalf("Bucket.Head returned error: %v", err)
	}
	if isExisted != false {
		t.Errorf("bucket IsExist failed")
	}

}

func TestBucketService_GetObjectVersions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		vs := values{
			"versions":  "",
			"delimiter": "/",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<ListVersionsResult>
    <Name>examplebucket-1250000000</Name>
    <Prefix/>
    <KeyMarker/>
    <VersionIdMarker/>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Delimiter>/</Delimiter>
    <CommonPrefixes>
        <Prefix>example-folder-1/</Prefix>
    </CommonPrefixes>
    <CommonPrefixes>
        <Prefix>example-folder-2/</Prefix>
    </CommonPrefixes>
    <Version>
        <Key>example-object-1.jpg</Key>
        <VersionId>MTg0NDUxNzgxMjEzNTU3NTk1Mjg</VersionId>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-16T10:45:53.000Z</LastModified>
        <ETag>&quot;5d1143df07a17b23320d0da161e2819e&quot;</ETag>
        <Size>30</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <DeleteMarker>
        <Key>example-object-1.jpg</Key>
        <VersionId>MTg0NDUxNzgxMjEzNjE1OTcxMzM</VersionId>
        <IsLatest>false</IsLatest>
        <LastModified>2019-08-16T10:45:47.000Z</LastModified>
        <Owner>
        <ID>1250000000</ID>
        <DisplayName>1250000000</DisplayName>
        </Owner>
    </DeleteMarker>
</ListVersionsResult>`)
	})

	want := &BucketGetObjectVersionsResult{
		XMLName:     xml.Name{Local: "ListVersionsResult"},
		Name:        "examplebucket-1250000000",
		MaxKeys:     1000,
		IsTruncated: false,
		Delimiter:   "/",
		CommonPrefixes: []string{
			"example-folder-1/",
			"example-folder-2/",
		},
		Version: []ListVersionsResultVersion{
			{
				Key:          "example-object-1.jpg",
				VersionId:    "MTg0NDUxNzgxMjEzNTU3NTk1Mjg",
				IsLatest:     true,
				LastModified: "2019-08-16T10:45:53.000Z",
				ETag:         "\"5d1143df07a17b23320d0da161e2819e\"",
				Size:         30,
				StorageClass: "STANDARD",
				Owner: &Owner{
					ID:          "1250000000",
					DisplayName: "1250000000",
				},
			},
		},
		DeleteMarker: []ListVersionsResultDeleteMarker{
			{
				Key:          "example-object-1.jpg",
				VersionId:    "MTg0NDUxNzgxMjEzNjE1OTcxMzM",
				IsLatest:     false,
				LastModified: "2019-08-16T10:45:47.000Z",
				Owner: &Owner{
					ID:          "1250000000",
					DisplayName: "1250000000",
				},
			},
		},
	}
	opt := &BucketGetObjectVersionsOptions{
		Delimiter: "/",
	}
	res, _, err := client.Bucket.GetObjectVersions(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.GetObjectVersions returned error: %v", err)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetObjectVersions returned\n%+v\nwant\n%+v", res, want)
	}

}

func TestBucketService_GetMeta(t *testing.T) {
	setup()
	defer teardown()
	want := &BucketGetMetadataResult{
        BucketUrl: "https://test-125000000.cos.ap-guangzhou.myqcloud.com",
        BucketName: "test-125000000",
        Location: "ap-guangzhou",
        MAZ: true,
        OFS: true,
		Encryption: &BucketGetEncryptionResult{
			XMLName: xml.Name{Local: "ServerSideEncryptionConfiguration"},
			Rule: &BucketEncryptionConfiguration{
				SSEAlgorithm: "AES256",
			},
		},
		ACL: &BucketGetACLResult{
			XMLName: xml.Name{Local: "AccessControlPolicy"},
			Owner: &Owner{
				ID:          "qcs::cam::uin/100000760461:uin/100000760461",
				DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
			},
			AccessControlList: []ACLGrant{
				{
					Grantee: &ACLGrantee{
						Type:        "RootAccount",
						ID:          "qcs::cam::uin/100000760461:uin/100000760461",
						DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
					},
					Permission: "FULL_CONTROL",
				},
				{
					Grantee: &ACLGrantee{
						Type:        "RootAccount",
						ID:          "qcs::cam::uin/100000760461:uin/100000760461",
						DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
					},
					Permission: "READ",
				},
			},
		},
		Website: &BucketGetWebsiteResult{
			XMLName: xml.Name{Local: "WebsiteConfiguration"},
			Index:   "index.html",
			RedirectProtocol: &RedirectRequestsProtocol{
				"https",
			},
			RoutingRules: &WebsiteRoutingRules{
				Rules: []WebsiteRoutingRule{
					{
						ConditionErrorCode: "404",
						RedirectProtocol:   "https",
						RedirectReplaceKey: "404.html",
					},
				},
			},
		},
		Logging: &BucketGetLoggingResult{
			XMLName: xml.Name{Local: "BucketLoggingStatus"},
			LoggingEnabled: &BucketLoggingEnabled{
				TargetBucket: "logs",
				TargetPrefix: "mylogs",
			},
		},
		CORS: &BucketGetCORSResult{
			XMLName: xml.Name{Local: "CORSConfiguration"},
			Rules: []BucketCORSRule{
				{
					AllowedOrigins: []string{"http://www.qq.com"},
					AllowedMethods: []string{"PUT", "GET"},
					AllowedHeaders: []string{"x-cos-meta-test", "x-cos-xx"},
					MaxAgeSeconds:  500,
					ExposeHeaders:  []string{"x-cos-meta-test1"},
				},
			},
		},
		Versioning: &BucketGetVersionResult{
			XMLName: xml.Name{Local: "VersioningConfiguration"},
			Status:  "Suspended",
		},
		Lifecycle: &BucketGetLifecycleResult{
			XMLName: xml.Name{Local: "LifecycleConfiguration"},
			Rules: []BucketLifecycleRule{
				{
					ID: "1234",
					Filter: &BucketLifecycleFilter{
						And: &BucketLifecycleAndOperator{
							Prefix: "test",
							Tag: []BucketTaggingTag{
								{Key: "key", Value: "value"},
							},
						},
					},
					Status: "Enabled",
					Transition: []BucketLifecycleTransition{
						{Days: 10, StorageClass: "Standard"},
					},
					Expiration: &BucketLifecycleExpiration{Days: 10},
					NoncurrentVersionExpiration: &BucketLifecycleNoncurrentVersion{
						NoncurrentDays: 360,
					},
					NoncurrentVersionTransition: []BucketLifecycleNoncurrentVersion{
						{
							NoncurrentDays: 90,
							StorageClass:   "ARCHIVE",
						},
					},
				},
				{
					ID:         "123422",
					Filter:     &BucketLifecycleFilter{Prefix: "gg"},
					Status:     "Disabled",
					Expiration: &BucketLifecycleExpiration{Days: 10},
				},
			},
		},
		IntelligentTiering: &ListIntelligentTieringConfigurations{
			XMLName: xml.Name{Local: "ListBucketIntelligentTieringConfigurationsOutput"},
			Configurations: []*IntelligentTieringConfiguration{
				{
					Id:     "default",
					Status: "Enabled",
					Tiering: []*BucketIntelligentTieringTransition{
						{
							AccessTier:      "INFREQUENT",
							Days:            30,
							RequestFrequent: 1,
						},
					},
				},
				{
					Id:     "1",
					Status: "Enabled",
					Filter: &BucketIntelligentTieringFilter{
						And: &BucketIntelligentTieringFilterAnd{
							Prefix: "test",
							Tag: []*BucketTaggingTag{
								{
									Key:   "k1",
									Value: "v1",
								},
							},
						},
					},
					Tiering: []*BucketIntelligentTieringTransition{
						{
							AccessTier: "ARCHIVE_ACCESS",
							Days:       30,
						},
					},
				},
			},
		},
		Tagging: &BucketGetTaggingResult{
			XMLName: xml.Name{Local: "Tagging"},
			TagSet: []BucketTaggingTag{
				{"test_k2", "test_v2"},
			},
		},
		ObjectLock: &BucketGetObjectLockResult{
			XMLName:           xml.Name{Local: "ObjectLockConfiguration"},
			ObjectLockEnabled: "Enabled",
			Rule: &ObjectLockRule{
				Days: 30,
			},
		},
		Replication: &GetBucketReplicationResult{
			XMLName: xml.Name{Local: "ReplicationConfiguration"},
			Role:    "qcs::cam::uin/100000000001:uin/100000000001",
			Rule: []BucketReplicationRule{
				{
					Status: "Disabled",
					Prefix: "prefix",
					Destination: &ReplicationDestination{
						Bucket: "qcs::cos:ap-beijing-1::examplebucket-1250000000",
					},
				},
			},
		},
	}

	//var head, encryption, acl, website, logging, cors, versioning, lifecycle, intelligenttiering, tagging, lock, replication int
	actionMap := map[string]func(http.ResponseWriter, *http.Request){
		"encryption": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<ServerSideEncryptionConfiguration>
                <Rule>
                    <ApplyServerSideEncryptionByDefault>
                        <SSEAlgorithm>AES256</SSEAlgorithm>
                    </ApplyServerSideEncryptionByDefault>
                </Rule>
            </ServerSideEncryptionConfiguration>`)
		},
		"acl": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<AccessControlPolicy>
	<Owner>
		<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
		<DisplayName>qcs::cam::uin/100000760461:uin/100000760461</DisplayName>
	</Owner>
	<AccessControlList>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="RootAccount">
				<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
				<DisplayName>qcs::cam::uin/100000760461:uin/100000760461</DisplayName>
			</Grantee>
			<Permission>FULL_CONTROL</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="RootAccount">
				<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
				<DisplayName>qcs::cam::uin/100000760461:uin/100000760461</DisplayName>
			</Grantee>
			<Permission>READ</Permission>
		</Grant>
	</AccessControlList>
</AccessControlPolicy>`)
		},
		"website": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<WebsiteConfiguration>
	<IndexDocument>
		<Suffix>index.html</Suffix>
	</IndexDocument>
	<RedirectAllRequestsTo>
		<Protocol>https</Protocol>
	</RedirectAllRequestsTo>
	<RoutingRules>
		<RoutingRule>
			<Condition>
				<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>
			</Condition>
			<Redirect>
				<Protocol>https</Protocol>
				<ReplaceKeyWith>404.html</ReplaceKeyWith>
			</Redirect>
		</RoutingRule>
	</RoutingRules>
</WebsiteConfiguration>`)

		},
		"logging": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<BucketLoggingStatus>
    <LoggingEnabled>
        <TargetBucket>logs</TargetBucket>
        <TargetPrefix>mylogs</TargetPrefix>
    </LoggingEnabled>
</BucketLoggingStatus>`)

		},
		"cors": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<CORSConfiguration>
	<CORSRule>
		<AllowedOrigin>http://www.qq.com</AllowedOrigin>
		<AllowedMethod>PUT</AllowedMethod>
		<AllowedMethod>GET</AllowedMethod>
		<AllowedHeader>x-cos-meta-test</AllowedHeader>
		<AllowedHeader>x-cos-xx</AllowedHeader>
		<ExposeHeader>x-cos-meta-test1</ExposeHeader>
		<MaxAgeSeconds>500</MaxAgeSeconds>
	</CORSRule>
</CORSConfiguration>`)

		},
		"versioning": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<VersioningConfiguration>
    <Status>Suspended</Status>
</VersioningConfiguration>`)

		},
		"lifecycle": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<LifecycleConfiguration>
	<Rule>
		<ID>1234</ID>
		<Filter>
            <And>
                <Prefix>test</Prefix>
                <Tag>
                    <Key>key</Key>
                    <Value>value</Value>
                </Tag>
            </And>
		</Filter>
		<Status>Enabled</Status>
		<Transition>
			<Days>10</Days>
			<StorageClass>Standard</StorageClass>
		</Transition>
		<Expiration>
			<Days>10</Days>
		</Expiration>
		<NoncurrentVersionTransition>
			<NoncurrentDays>90</NoncurrentDays>
			<StorageClass>ARCHIVE</StorageClass>
		</NoncurrentVersionTransition>
		<NoncurrentVersionExpiration>
			<NoncurrentDays>360</NoncurrentDays>
		</NoncurrentVersionExpiration>
	</Rule>
	<Rule>
		<ID>123422</ID>
		<Filter>
			<Prefix>gg</Prefix>
		</Filter>
		<Status>Disabled</Status>
		<Expiration>
			<Days>10</Days>
		</Expiration>
	</Rule>
</LifecycleConfiguration>`)

		},
		"intelligent-tiering": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<ListBucketIntelligentTieringConfigurationsOutput>
    <IntelligentTieringConfiguration>
        <Id>default</Id>
        <Status>Enabled</Status>
        <Tiering>
            <AccessTier>INFREQUENT</AccessTier>
            <Days>30</Days>
            <RequestFrequent>1</RequestFrequent>
        </Tiering>
    </IntelligentTieringConfiguration>
    <IntelligentTieringConfiguration>
        <Id>1</Id>
        <Status>Enabled</Status>
        <Filter>
            <And>
                <Prefix>test</Prefix>
                <Tag>
                    <Key>k1</Key>
                    <Value>v1</Value>
                </Tag>
            </And>
        </Filter>
        <Tiering>
            <AccessTier>ARCHIVE_ACCESS</AccessTier>
            <Days>30</Days>
        </Tiering>
    </IntelligentTieringConfiguration>
</ListBucketIntelligentTieringConfigurationsOutput>`)
		},
		"tagging": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<Tagging>
	<TagSet>
		<Tag>
			<Key>test_k2</Key>
			<Value>test_v2</Value>
		</Tag>
	</TagSet>
</Tagging>`)

		},
		"object-lock": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<ObjectLockConfiguration>
	<ObjectLockEnabled>Enabled</ObjectLockEnabled> 
	<Rule> 
		<DefaultRetention>
			<Days>30</Days> 
		</DefaultRetention> 
	</Rule> 
</ObjectLockConfiguration>`)

		},
		"replication": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `<ReplicationConfiguration>
    <Role>qcs::cam::uin/100000000001:uin/100000000001</Role>
    <Rule>
        <Status>Disabled</Status>
        <ID></ID>
        <Prefix>prefix</Prefix>
        <Destination>
            <Bucket>qcs::cos:ap-beijing-1::examplebucket-1250000000</Bucket>
        </Destination>
    </Rule>
</ReplicationConfiguration>`)
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Add("X-Cos-Bucket-Az-Type", "MAZ")
			w.Header().Add("X-Cos-Bucket-Arch", "OFS")
			w.Header().Add("X-Cos-Bucket-Region", "ap-guangzhou")
		} else {
			r.ParseForm()
			for key, _ := range r.Form {
				if fn, ok := actionMap[key]; ok {
					fn(w, r)
				}
			}
		}
	})
	// BucketURL为空
	tmpUrl := client.BaseURL.BucketURL
	client.BaseURL.BucketURL = nil
	_, _, err := client.Bucket.GetMeta(context.Background())
	if err == nil || err.Error() != "BucketURL is empty" {
		t.Fatalf("Bucket.GetMeta returned error: %v", err)
	}
	client.BaseURL.BucketURL = tmpUrl
	// 没有提供bucketname
	_, _, err = client.Bucket.GetMeta(context.Background())
	if err == nil || err.Error() != "you must provide bucket-appid param in using custom domain" {
		t.Fatalf("Bucket.GetMeta returned error: %v", err)
	}
    // 成功
	res, _, err := client.Bucket.GetMeta(context.Background(), "test-125000000")
    if err != nil {
        t.Fatalf("Bucket.GetMeta returned error: %v", err)
    }
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetObjectVersions returned\n%+v\nwant\n%+v", res, want)
	}
}

func TestBucketService_GetMeta404(t *testing.T) {
	setup()
	defer teardown()
	want := &BucketGetMetadataResult{
        BucketUrl: "https://test-125000000.cos.ap-guangzhou.myqcloud.com",
        BucketName: "test-125000000",
        Location: "ap-guangzhou",
        MAZ: true,
        OFS: true,
	}

	//var head, encryption, acl, website, logging, cors, versioning, lifecycle, intelligenttiering, tagging, lock, replication int
	actionMap := map[string]func(http.ResponseWriter, *http.Request){
		"encryption": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"acl": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"website": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"logging": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"cors": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"versioning": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"lifecycle": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"intelligent-tiering": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"tagging": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"object-lock": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
		"replication": func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNotFound)
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Add("X-Cos-Bucket-Az-Type", "MAZ")
			w.Header().Add("X-Cos-Bucket-Arch", "OFS")
			w.Header().Add("X-Cos-Bucket-Region", "ap-guangzhou")
		} else {
			r.ParseForm()
			for key, _ := range r.Form {
				if fn, ok := actionMap[key]; ok {
					fn(w, r)
				}
			}
		}
	})
	// BucketURL为空
	tmpUrl := client.BaseURL.BucketURL
	client.BaseURL.BucketURL = nil
	_, _, err := client.Bucket.GetMeta(context.Background())
	if err == nil || err.Error() != "BucketURL is empty" {
		t.Fatalf("Bucket.GetMeta returned error: %v", err)
	}
	client.BaseURL.BucketURL = tmpUrl
	// 没有提供bucketname
	_, _, err = client.Bucket.GetMeta(context.Background())
	if err == nil || err.Error() != "you must provide bucket-appid param in using custom domain" {
		t.Fatalf("Bucket.GetMeta returned error: %v", err)
	}
    // 成功
	res, _, err := client.Bucket.GetMeta(context.Background(), "test-125000000")
    if err != nil {
        t.Fatalf("Bucket.GetMeta returned error: %v", err)
    }
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetObjectVersions returned\n%+v\nwant\n%+v", res, want)
	}
}

