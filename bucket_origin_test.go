package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetOrigin(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"origin": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<OriginConfiguration>
        <OriginRule>
            <RulePriority>1</RulePriority>
            <OriginType>Mirror</OriginType>
            <OriginCondition>
                <HTTPStatusCode>404</HTTPStatusCode>
                <Prefix></Prefix>
            </OriginCondition>
            <OriginParameter>
                <Protocol>HTTP</Protocol>
                <FollowQueryString>false</FollowQueryString>
                <HttpHeader>
                    <NewHttpHeaders>
                        <Header>
                            <Key>x-cos</Key>
                            <Value>exampleHeader</Value>
                        </Header>
                    </NewHttpHeaders>
                    <FollowHttpHeaders>
                        <Header>
                            <Key>exampleHeaderKey</Key>
                        </Header>
                    </FollowHttpHeaders>
                </HttpHeader>
                <FollowRedirection>true</FollowRedirection>
                <HttpRedirectCode>302</HttpRedirectCode>
            </OriginParameter>
            <OriginInfo>
                <HostInfo>
                    <HostName>examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com</HostName>
                    <Weight>10</Weight>
                    <StandbyHostName_1>hostname1</StandbyHostName_1>
                    <StandbyHostName_2>hostname2</StandbyHostName_2>
					<PrivateHost>
						<Host>www.qq.com</Host>
						<CredentialProvider>
							<Role>qcs::cam::uin/123:roleName/name</Role>
						</CredentialProvider>
					</PrivateHost>
					<PrivateStandbyHost_1>
						<Host>1.qq.com</Host>
						<CredentialProvider>
							<Role>qcs::cam::uin/123:roleName/name</Role>
						</CredentialProvider>
					</PrivateStandbyHost_1>
					<PrivateStandbyHost_2>
						<Host>2.qq.com</Host>
						<CredentialProvider>
							<Role>qcs::cam::uin/123:roleName/name</Role>
						</CredentialProvider>
					</PrivateStandbyHost_2>
                </HostInfo>
				<HostInfo>
                    <HostName>examplebucket2-1250000000.cos.ap-shanghai.myqcloud.com</HostName>
                    <Weight>10</Weight>
                    <StandbyHostName_3>hostname3</StandbyHostName_3>
                    <StandbyHostName_4>hostname4</StandbyHostName_4>
					<PrivateHost>
						<Host>www.qq1.com</Host>
						<CredentialProvider>
							<Role>qcs::cam::uin/123:roleName/name</Role>
						</CredentialProvider>
					</PrivateHost>
					<PrivateStandbyHost_1>
						<Host>1.qq1.com</Host>
						<CredentialProvider>
							<Role>qcs::cam::uin/123:roleName/name</Role>
						</CredentialProvider>
					</PrivateStandbyHost_1>
					<PrivateStandbyHost_2>
						<Host>2.qq1.com</Host>
						<CredentialProvider>
							<Role>qcs::cam::uin/123:roleName/name</Role>
						</CredentialProvider>
					</PrivateStandbyHost_2>
                </HostInfo>
            </OriginInfo>
        </OriginRule>
        </OriginConfiguration>
        `)
	})

	res, _, err := client.Bucket.GetOrigin(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetOrigin returned error %v", err)
	}

	want := &BucketGetOriginResult{
		XMLName: xml.Name{Local: "OriginConfiguration"},
		Rule: []BucketOriginRule{
			{
				OriginType:   "Mirror",
				RulePriority: 1,
				OriginCondition: &BucketOriginCondition{
					HTTPStatusCode: "404",
				},
				OriginParameter: &BucketOriginParameter{
					Protocol:          "HTTP",
					FollowQueryString: Bool(false),
					HttpHeader: &BucketOriginHttpHeader{
						FollowHttpHeaders: []OriginHttpHeader{
							{
								Key: "exampleHeaderKey",
							},
						},
						NewHttpHeaders: []OriginHttpHeader{
							{
								Key:   "x-cos",
								Value: "exampleHeader",
							},
						},
					},
					FollowRedirection: Bool(true),
					HttpRedirectCode:  "302",
				},
				OriginInfo: &BucketOriginInfo{
					HostInfo: &BucketOriginHostInfo{
						HostName:          "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
						Weight:            10,
						StandbyHostName_N: []string{"hostname1", "hostname2"},
						StandbyHostName: []*BucketOriginStandbyHost{
							&BucketOriginStandbyHost{
								Index:    1,
								HostName: "hostname1",
							},
							&BucketOriginStandbyHost{
								Index:    2,
								HostName: "hostname2",
							},
						},
						PrivateHost: &BucketOriginPrivateHost{
							Host: "www.qq.com",
							CredentialProvider: &BucketOriginCredentialProvider{
								Role: "qcs::cam::uin/123:roleName/name",
							},
						},
						PrivateStandbyHost_N: []*BucketOriginPrivateHost{
							&BucketOriginPrivateHost{
								Index: 1,
								Host:  "1.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							&BucketOriginPrivateHost{
								Index: 2,
								Host:  "2.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
						},
					},
					HostInfos: []*BucketOriginHostInfo{
						&BucketOriginHostInfo{
							HostName:          "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname2"},
							StandbyHostName: []*BucketOriginStandbyHost{
								&BucketOriginStandbyHost{
									Index:    1,
									HostName: "hostname1",
								},
								&BucketOriginStandbyHost{
									Index:    2,
									HostName: "hostname2",
								},
							},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Index: 1,
									Host:  "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Index: 2,
									Host:  "2.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
						&BucketOriginHostInfo{
							HostName:          "examplebucket2-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname3", "hostname4"},
							StandbyHostName: []*BucketOriginStandbyHost{
								&BucketOriginStandbyHost{
									Index:    3,
									HostName: "hostname3",
								},
								&BucketOriginStandbyHost{
									Index:    4,
									HostName: "hostname4",
								},
							},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq1.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Index: 1,
									Host:  "1.qq1.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Index: 2,
									Host:  "2.qq1.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetOrigin returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutOrigin(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutOriginOptions{
		XMLName: xml.Name{Local: "OriginConfiguration"},
		Rule: []BucketOriginRule{
			{
				OriginType:   "Mirror",
				RulePriority: 1,
				OriginCondition: &BucketOriginCondition{
					HTTPStatusCode: "404",
				},
				OriginParameter: &BucketOriginParameter{
					Protocol:          "HTTP",
					FollowQueryString: Bool(true),
					HttpHeader: &BucketOriginHttpHeader{
						FollowHttpHeaders: []OriginHttpHeader{
							{
								Key: "exampleHeaderKey",
							},
						},
						NewHttpHeaders: []OriginHttpHeader{
							{
								Key:   "x-cos",
								Value: "exampleHeader",
							},
						},
					},
					FollowRedirection: Bool(true),
					HttpRedirectCode:  "302",
				},
				OriginInfo: &BucketOriginInfo{
					HostInfo: &BucketOriginHostInfo{
						HostName: "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
						Weight:   10,
						StandbyHostName: []*BucketOriginStandbyHost{
							&BucketOriginStandbyHost{
								Index:    1,
								HostName: "hostname1",
							},
							&BucketOriginStandbyHost{
								Index:    4,
								HostName: "hostname4",
							},
						},
						StandbyHostName_N: []string{"hostname1", "hostname4"},
						PrivateHost: &BucketOriginPrivateHost{
							Host: "www.qq.com",
							CredentialProvider: &BucketOriginCredentialProvider{
								Role: "qcs::cam::uin/123:roleName/name",
							},
						},
						PrivateStandbyHost_N: []*BucketOriginPrivateHost{
							&BucketOriginPrivateHost{
								Index: 2,
								Host:  "2.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							&BucketOriginPrivateHost{
								Index: 3,
								Host:  "1.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
						},
					},
					HostInfos: []*BucketOriginHostInfo{
						&BucketOriginHostInfo{
							HostName:          "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname4"},
							StandbyHostName: []*BucketOriginStandbyHost{
								&BucketOriginStandbyHost{
									Index:    1,
									HostName: "hostname1",
								},
								&BucketOriginStandbyHost{
									Index:    4,
									HostName: "hostname4",
								},
							},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Index: 2,
									Host:  "2.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Index: 3,
									Host:  "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
						&BucketOriginHostInfo{
							HostName:          "examplebucket2-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname3", "hostname4"},
							StandbyHostName: []*BucketOriginStandbyHost{
								&BucketOriginStandbyHost{
									Index:    3,
									HostName: "hostname3",
								},
								&BucketOriginStandbyHost{
									Index:    4,
									HostName: "hostname4",
								},
							},

							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq1.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Index: 1,
									Host:  "1.qq1.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Index: 2,
									Host:  "2.qq1.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"origin": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutOriginOptions)
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Bucket.PutOrigin request body read error: %v", err)
		}
		err = xml.Unmarshal(bs, body)
		if err != nil {
			t.Fatalf("Bucket.PutOrigin request body xml unmarshal error: %v", err)
		}
		want := opt
		want.XMLName = xml.Name{Local: "OriginConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutOrigin request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutOrigin(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutOrigin returned error: %v", err)
	}
}

func TestBucketService_DeleteOrigin(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"origin": "",
		}
		testFormValues(t, r, vs)
		w.WriteHeader(http.StatusNoContent)
	})
	_, err := client.Bucket.DeleteOrigin(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteOrigin returned error: %v", err)
	}
}

func TestBucketService_OriginXml(t *testing.T) {
	xmlBody := `<OriginConfiguration><OriginRule><OriginInfo><HostInfo><HostName>examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com</HostName><Weight>10</Weight><StandbyHostName_1>hostname1</StandbyHostName_1><StandbyHostName_2>hostname2</StandbyHostName_2><PrivateHost><Host>www.qq.com</Host><CredentialProvider><Role>qcs::cam::uin/123:roleName/name</Role></CredentialProvider></PrivateHost><PrivateStandbyHost_1><Host>1.qq.com</Host><CredentialProvider><Role>qcs::cam::uin/123:roleName/name</Role></CredentialProvider></PrivateStandbyHost_1><PrivateStandbyHost_2><Host>2.qq.com</Host><CredentialProvider><Role>qcs::cam::uin/123:roleName/name</Role></CredentialProvider></PrivateStandbyHost_2></HostInfo><HostInfo><HostName>examplebucket2-1250000000.cos.ap-shanghai.myqcloud.com</HostName><Weight>10</Weight><StandbyHostName_1>hostname3</StandbyHostName_1><StandbyHostName_2>hostname4</StandbyHostName_2><PrivateHost><Host>www.qq1.com</Host><CredentialProvider><Role>qcs::cam::uin/123:roleName/name</Role></CredentialProvider></PrivateHost><PrivateStandbyHost_1><Host>1.qq1.com</Host><CredentialProvider><Role>qcs::cam::uin/123:roleName/name</Role></CredentialProvider></PrivateStandbyHost_1><PrivateStandbyHost_2><Host>2.qq1.com</Host><CredentialProvider><Role>qcs::cam::uin/123:roleName/name</Role></CredentialProvider></PrivateStandbyHost_2></HostInfo></OriginInfo></OriginRule></OriginConfiguration>`

	want := &BucketGetOriginResult{
		XMLName: xml.Name{Local: "OriginConfiguration"},
		Rule: []BucketOriginRule{
			{
				OriginInfo: &BucketOriginInfo{
					// 兼容HostInfo
					HostInfo: &BucketOriginHostInfo{
						HostName:          "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
						Weight:            10,
						StandbyHostName_N: []string{"hostname1", "hostname2"},
						StandbyHostName: []*BucketOriginStandbyHost{
							&BucketOriginStandbyHost{
								Index:    1,
								HostName: "hostname1",
							},
							&BucketOriginStandbyHost{
								Index:    2,
								HostName: "hostname2",
							},
						},
						PrivateHost: &BucketOriginPrivateHost{
							Host: "www.qq.com",
							CredentialProvider: &BucketOriginCredentialProvider{
								Role: "qcs::cam::uin/123:roleName/name",
							},
						},
						PrivateStandbyHost_N: []*BucketOriginPrivateHost{
							&BucketOriginPrivateHost{
								Index: 1,
								Host:  "1.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							&BucketOriginPrivateHost{
								Index: 2,
								Host:  "2.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
						},
					},
					HostInfos: []*BucketOriginHostInfo{
						&BucketOriginHostInfo{
							HostName:          "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname2"},
							StandbyHostName: []*BucketOriginStandbyHost{
								&BucketOriginStandbyHost{
									Index:    1,
									HostName: "hostname1",
								},
								&BucketOriginStandbyHost{
									Index:    2,
									HostName: "hostname2",
								},
							},

							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Index: 1,
									Host:  "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Index: 2,
									Host:  "2.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
						&BucketOriginHostInfo{
							HostName:          "examplebucket2-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname3", "hostname4"},
							StandbyHostName: []*BucketOriginStandbyHost{
								&BucketOriginStandbyHost{
									Index:    1,
									HostName: "hostname3",
								},
								&BucketOriginStandbyHost{
									Index:    2,
									HostName: "hostname4",
								},
							},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq1.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Index: 1,
									Host:  "1.qq1.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Index: 2,
									Host:  "2.qq1.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	var res BucketGetOriginResult
	err := xml.Unmarshal([]byte(xmlBody), &res)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(want, &res) {
		t.Errorf("BucketGetOriginResult = %+v, want %+v", res, want)
	}

}

func TestBucketService_OriginXml2(t *testing.T) {
	want := &BucketPutOriginOptions{
		XMLName: xml.Name{Local: "OriginConfiguration"},
		Rule: []BucketOriginRule{
			{
				OriginInfo: &BucketOriginInfo{
					// 兼容HostInfo
					HostInfo: &BucketOriginHostInfo{
						HostName:          "examplebucket1-1250000000.cos.ap-shanghai.myqcloud.com",
						Weight:            10,
						StandbyHostName_N: []string{"hostname1", "hostname2"},
						StandbyHostName: []*BucketOriginStandbyHost{
							&BucketOriginStandbyHost{
								Index:    1,
								HostName: "hostname1",
							},
							&BucketOriginStandbyHost{
								Index:    2,
								HostName: "hostname2",
							},
						},
						PrivateHost: &BucketOriginPrivateHost{
							Host: "www.qq.com",
							CredentialProvider: &BucketOriginCredentialProvider{
								Role: "qcs::cam::uin/123:roleName/name",
							},
						},
						PrivateStandbyHost_N: []*BucketOriginPrivateHost{
							&BucketOriginPrivateHost{
								Index: 1,
								Host:  "1.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							&BucketOriginPrivateHost{
								Index: 2,
								Host:  "2.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
						},
					},
				},
			},
		},
	}
	body, err := xml.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	t.Logf("body: %s", string(body))
}

func TestBucketService_BucketOriginHostInfoUnmarshal(t *testing.T) {
	str := `
<HostInfo>
	<Weight>1</Weight>
	<PrivateHost>
		<Host>hostname1</Host>
		<CredentialProvider>
			<SecretId>1</SecretId>
			<EncryptedSecretKey>key</EncryptedSecretKey>
			<Region>1</Region>
			<AuthorizationAlgorithm>S3</AuthorizationAlgorithm>
		</CredentialProvider>
	</PrivateHost>
</HostInfo>`
	var info BucketOriginHostInfo
	err := xml.Unmarshal([]byte(str), &info)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}
	want := BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "hostname1",
			CredentialProvider: &BucketOriginCredentialProvider{
				SecretId:               "1",
				EncryptedSecretKey:     "key",
				Region:                 "1",
				AuthorizationAlgorithm: "S3",
			},
		},
	}
	if !reflect.DeepEqual(info, want) {
		t.Fatalf("BucketOriginHostInfo unmarshal err, res: %v, want: %v", info, want)
		return
	}

	str = `
<HostInfo>
	<Weight>1</Weight>
	<PrivateHost>
		<Host>test-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateHost>
	<PrivateStandbyHost_1>
		<Host>test2-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateStandbyHost_1>
	<StandbyHostName_4>hostname4</StandbyHostName_4>
	<PrivateStandbyHost_3>
		<Host>test3-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateStandbyHost_3>
	<StandbyHostName_2>hostname2</StandbyHostName_2>
</HostInfo>`

	info = BucketOriginHostInfo{}
	err = xml.Unmarshal([]byte(str), &info)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}

	want = BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "test-125000000.cos.ap-shanghai.myqcloud.com",
			CredentialProvider: &BucketOriginCredentialProvider{
				Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
			},
		},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Index: 1,
				Host:  "test2-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
			&BucketOriginPrivateHost{
				Index: 3,
				Host:  "test3-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
		},
		StandbyHostName_N: []string{"hostname2", "hostname4"},
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				Index:    2,
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				Index:    4,
				HostName: "hostname4",
			},
		},
	}
	if !reflect.DeepEqual(info, want) {
		t.Fatalf("BucketOriginHostInfo unmarshal err, res: %v, want: %v", info, want)
		return
	}

	str = `
<HostInfo>
	<Weight>1</Weight>
	<PrivateHost>
		<Host>test-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateHost>
	<PrivateStandbyHost_1>
		<Host>hostname1</Host>
		<CredentialProvider>
			<SecretId>test</SecretId>
		</CredentialProvider>
	</PrivateStandbyHost_1>
	<PrivateStandbyHost_2>
		<Host>hostname2</Host>
		<CredentialProvider>
			<SecretId>id</SecretId>
		</CredentialProvider>
	</PrivateStandbyHost_2>
</HostInfo>`
	info = BucketOriginHostInfo{}
	err = xml.Unmarshal([]byte(str), &info)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}
	want = BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "test-125000000.cos.ap-shanghai.myqcloud.com",
			CredentialProvider: &BucketOriginCredentialProvider{
				Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
			},
		},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Index: 1,
				Host:  "hostname1",
				CredentialProvider: &BucketOriginCredentialProvider{
					SecretId: "test",
				},
			},
			&BucketOriginPrivateHost{
				Index: 2,
				Host:  "hostname2",
				CredentialProvider: &BucketOriginCredentialProvider{
					SecretId: "id",
				},
			},
		},
	}
	if !reflect.DeepEqual(info, want) {
		t.Fatalf("BucketOriginHostInfo unmarshal err, res: %v, want: %v", info, want)
		return
	}

	str = `
<HostInfo>
	<Weight>1</Weight>
	<PrivateHost>
		<Host>test-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateHost>
	<PrivateStandbyHost_1>
		<Host>test2-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateStandbyHost_1>
	<StandbyHostName_4>hostname4</StandbyHostName_4>
	<PrivateStandbyHost_3>
		<Host>test3-125000000.cos.ap-shanghai.myqcloud.com</Host>
		<CredentialProvider>
			<Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role>
		</CredentialProvider>
	</PrivateStandbyHost_3>
	<StandbyHostName_2>hostname2</StandbyHostName_2>
	<StandbyHostName_5>hostname5</StandbyHostName_5>
	<PrivateStandbyHost_5>
		<Host>test5-125000000.cos.ap-shanghai.myqcloud.com</Host>
	</PrivateStandbyHost_5>
</HostInfo>`

	info = BucketOriginHostInfo{}
	err = xml.Unmarshal([]byte(str), &info)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}

	want = BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "test-125000000.cos.ap-shanghai.myqcloud.com",
			CredentialProvider: &BucketOriginCredentialProvider{
				Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
			},
		},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Index: 1,
				Host:  "test2-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
			&BucketOriginPrivateHost{
				Index: 3,
				Host:  "test3-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
			&BucketOriginPrivateHost{
				Index: 5,
				Host:  "test5-125000000.cos.ap-shanghai.myqcloud.com",
			},
		},
		StandbyHostName_N: []string{"hostname2", "hostname4", "hostname5"},
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				Index:    2,
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				Index:    4,
				HostName: "hostname4",
			},
			&BucketOriginStandbyHost{
				Index:    5,
				HostName: "hostname5",
			},
		},
	}
	if !reflect.DeepEqual(info, want) {
		t.Fatalf("BucketOriginHostInfo unmarshal err, res: %v, want: %v", info, want)
		return
	}

	str = `
<HostInfo>
	<Weight>1</Weight>
	<StandbyHostName_4>hostname4</StandbyHostName_4>
	<StandbyHostName_2>hostname2</StandbyHostName_2>
	<StandbyHostName_5>hostname5</StandbyHostName_5>
</HostInfo>`

	info = BucketOriginHostInfo{}
	err = xml.Unmarshal([]byte(str), &info)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}

	want = BucketOriginHostInfo{
		Weight:            1,
		StandbyHostName_N: []string{"hostname2", "hostname4", "hostname5"},
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				Index:    2,
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				Index:    4,
				HostName: "hostname4",
			},
			&BucketOriginStandbyHost{
				Index:    5,
				HostName: "hostname5",
			},
		},
	}
	if !reflect.DeepEqual(info, want) {
		t.Fatalf("BucketOriginHostInfo unmarshal err, res: %v, want: %v", info, want)
		return
	}

	str = `
<HostInfo>
	<Weight>1</Weight>
	<StandbyHostName_x>hostname4</StandbyHostName_x>
</HostInfo>`

	info = BucketOriginHostInfo{}
	err = xml.Unmarshal([]byte(str), &info)
	if err == nil || err.Error() != "StandbyHostName Parse failed, node: StandbyHostName_x" {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}

	str = `
<HostInfo>
	<PrivateStandbyHost_x>
		<Host>test2-125000000.cos.ap-shanghai.myqcloud.com</Host>
	</PrivateStandbyHost_x>
</HostInfo>`

	info = BucketOriginHostInfo{}
	err = xml.Unmarshal([]byte(str), &info)
	if err == nil || err.Error() != "PrivateStandbyHost Parse failed, node: PrivateStandbyHost_x" {
		t.Fatalf("Unmarshal error: %v", err)
		return
	}

}

func TestBucketService_BucketOriginHostInfoMarshal(t *testing.T) {
	opt := &BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "hostname1",
			CredentialProvider: &BucketOriginCredentialProvider{
				SecretId:               "1",
				EncryptedSecretKey:     "key",
				Region:                 "1",
				AuthorizationAlgorithm: "S3",
			},
		},
	}
	want := `<BucketOriginHostInfo><Weight>1</Weight><PrivateHost><Host>hostname1</Host><CredentialProvider><AuthorizationAlgorithm>S3</AuthorizationAlgorithm><Region>1</Region><SecretId>1</SecretId><EncryptedSecretKey>key</EncryptedSecretKey></CredentialProvider></PrivateHost></BucketOriginHostInfo>`
	bs, err := xml.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	res := string(bs)
	if res != want {
		t.Fatalf("BucketOriginHostInfo marshal err, res: %v, want: %v", res, want)
		return
	}

	want = `<BucketOriginHostInfo><Weight>1</Weight><PrivateHost><Host>test-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateHost><StandbyHostName_2>hostname2</StandbyHostName_2><StandbyHostName_4>hostname4</StandbyHostName_4><PrivateStandbyHost_1><Host>test2-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateStandbyHost_1><PrivateStandbyHost_3><Host>test3-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateStandbyHost_3></BucketOriginHostInfo>`

	opt = &BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "test-125000000.cos.ap-shanghai.myqcloud.com",
			CredentialProvider: &BucketOriginCredentialProvider{
				Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
			},
		},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Index: 1,
				Host:  "test2-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
			&BucketOriginPrivateHost{
				Index: 3,
				Host:  "test3-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
		},
		StandbyHostName_N: []string{"hostname2", "hostname4"},
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				Index:    2,
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				Index:    4,
				HostName: "hostname4",
			},
		},
	}
	bs, err = xml.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	res = string(bs)
	if res != want {
		t.Fatalf("BucketOriginHostInfo marshal err, res: %v, want: %v", res, want)
	}

	want = `<BucketOriginHostInfo><Weight>1</Weight><PrivateHost><Host>test-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateHost><PrivateStandbyHost_1><Host>hostname1</Host><CredentialProvider><SecretId>test</SecretId></CredentialProvider></PrivateStandbyHost_1><PrivateStandbyHost_2><Host>hostname2</Host><CredentialProvider><SecretId>id</SecretId></CredentialProvider></PrivateStandbyHost_2></BucketOriginHostInfo>`
	opt = &BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "test-125000000.cos.ap-shanghai.myqcloud.com",
			CredentialProvider: &BucketOriginCredentialProvider{
				Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
			},
		},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Index: 1,
				Host:  "hostname1",
				CredentialProvider: &BucketOriginCredentialProvider{
					SecretId: "test",
				},
			},
			&BucketOriginPrivateHost{
				Index: 2,
				Host:  "hostname2",
				CredentialProvider: &BucketOriginCredentialProvider{
					SecretId: "id",
				},
			},
		},
	}
	bs, err = xml.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	res = string(bs)
	if res != want {
		t.Fatalf("BucketOriginHostInfo marshal err, res: %v, want: %v", res, want)
	}

	want = `<BucketOriginHostInfo><Weight>1</Weight><PrivateHost><Host>test-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateHost><StandbyHostName_2>hostname2</StandbyHostName_2><StandbyHostName_4>hostname4</StandbyHostName_4><StandbyHostName_5>hostname5</StandbyHostName_5><PrivateStandbyHost_1><Host>test2-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateStandbyHost_1><PrivateStandbyHost_3><Host>test3-125000000.cos.ap-shanghai.myqcloud.com</Host><CredentialProvider><Role>qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole</Role></CredentialProvider></PrivateStandbyHost_3><PrivateStandbyHost_5><Host>test5-125000000.cos.ap-shanghai.myqcloud.com</Host></PrivateStandbyHost_5></BucketOriginHostInfo>`

	opt = &BucketOriginHostInfo{
		Weight: 1,
		PrivateHost: &BucketOriginPrivateHost{
			Host: "test-125000000.cos.ap-shanghai.myqcloud.com",
			CredentialProvider: &BucketOriginCredentialProvider{
				Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
			},
		},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Index: 1,
				Host:  "test2-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
			&BucketOriginPrivateHost{
				Index: 3,
				Host:  "test3-125000000.cos.ap-shanghai.myqcloud.com",
				CredentialProvider: &BucketOriginCredentialProvider{
					Role: "qcs::cam::uin/10000001:roleName/COSOrigin_QCSRole",
				},
			},
			&BucketOriginPrivateHost{
				Index: 5,
				Host:  "test5-125000000.cos.ap-shanghai.myqcloud.com",
			},
		},
		StandbyHostName_N: []string{"hostname2", "hostname4", "hostname5"},
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				Index:    2,
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				Index:    4,
				HostName: "hostname4",
			},
			&BucketOriginStandbyHost{
				Index:    5,
				HostName: "hostname5",
			},
		},
	}
	bs, err = xml.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	res = string(bs)
	if res != want {
		t.Fatalf("BucketOriginHostInfo marshal err, res: %v, want: %v", res, want)
	}

	want = `<BucketOriginHostInfo><Weight>1</Weight><StandbyHostName_2>hostname2</StandbyHostName_2><StandbyHostName_4>hostname4</StandbyHostName_4><StandbyHostName_5>hostname5</StandbyHostName_5></BucketOriginHostInfo>`

	opt = &BucketOriginHostInfo{
		Weight: 1,
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				Index:    2,
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				Index:    4,
				HostName: "hostname4",
			},
			&BucketOriginStandbyHost{
				Index:    5,
				HostName: "hostname5",
			},
		},
	}
	bs, err = xml.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	res = string(bs)
	if res != want {
		t.Fatalf("BucketOriginHostInfo marshal err, res: %v, want: %v", res, want)
	}

	want = `<BucketOriginHostInfo><Weight>1</Weight><StandbyHostName_1>hostname2</StandbyHostName_1><StandbyHostName_2>hostname4</StandbyHostName_2><StandbyHostName_3>hostname5</StandbyHostName_3></BucketOriginHostInfo>`

	opt = &BucketOriginHostInfo{
		Weight:            1,
		StandbyHostName_N: []string{"hostname2", "hostname4", "hostname5"},
	}
	bs, err = xml.Marshal(opt)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
		return
	}
	res = string(bs)
	if res != want {
		t.Fatalf("BucketOriginHostInfo marshal err, res: %v, want: %v", res, want)
	}

	opt = &BucketOriginHostInfo{
		Weight: 1,
		StandbyHostName: []*BucketOriginStandbyHost{
			&BucketOriginStandbyHost{
				HostName: "hostname2",
			},
			&BucketOriginStandbyHost{
				HostName: "hostname4",
			},
			&BucketOriginStandbyHost{
				HostName: "hostname5",
			},
		},
	}
	_, err = xml.Marshal(opt)
	if err == nil || err.Error() != "The parameter Index must be set in StandbyHostName" {
		t.Fatalf("BucketOriginHostInfo marshal expect err")
		return
	}

	opt = &BucketOriginHostInfo{
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Host: "test2-125000000.cos.ap-shanghai.myqcloud.com",
			},
			&BucketOriginPrivateHost{
				Host: "test3-125000000.cos.ap-shanghai.myqcloud.com",
			},
			&BucketOriginPrivateHost{
				Host: "test5-125000000.cos.ap-shanghai.myqcloud.com",
			},
		},
	}
	_, err = xml.Marshal(opt)
	if err == nil || err.Error() != "The parameter Index must be set in PrivateStandbyHost_N" {
		t.Fatalf("BucketOriginHostInfo marshal expect err")
		return
	}

	opt = &BucketOriginHostInfo{
		Weight:            1,
		StandbyHostName_N: []string{"hostname2", "hostname4", "hostname5"},
		PrivateStandbyHost_N: []*BucketOriginPrivateHost{
			&BucketOriginPrivateHost{
				Host: "test2-125000000.cos.ap-shanghai.myqcloud.com",
			},
		},
	}
	_, err = xml.Marshal(opt)
	if err == nil || err.Error() != "StandbyHostName_N and PrivateStandbyHost_N can not be both set, use StandbyHostName and PrivateStandbyHost_N instand" {
		t.Fatalf("BucketOriginHostInfo marshal expect err")
		return
	}
}
