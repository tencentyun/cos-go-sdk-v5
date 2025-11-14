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
                    <HostName>examplebucket-1250000000.cos.ap-shanghai.myqcloud.com</HostName>
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
                    <HostName>examplebucket-1250000000.cos.ap-shanghai.myqcloud.com</HostName>
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
					HostInfo: []*BucketOriginHostInfo{
						&BucketOriginHostInfo{
							HostName:          "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname2"},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Host: "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Host: "2.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
						&BucketOriginHostInfo{
							HostName:          "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname2"},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Host: "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Host: "2.qq.com",
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
	t.Logf("HostInfo: %+v", res.Rule[0].OriginInfo.HostInfo)
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
					HostInfo: []*BucketOriginHostInfo{
						&BucketOriginHostInfo{
							HostName:          "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname2"},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Host: "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Host: "2.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
							},
						},
						&BucketOriginHostInfo{
							HostName:          "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
							Weight:            10,
							StandbyHostName_N: []string{"hostname1", "hostname2"},
							PrivateHost: &BucketOriginPrivateHost{
								Host: "www.qq.com",
								CredentialProvider: &BucketOriginCredentialProvider{
									Role: "qcs::cam::uin/123:roleName/name",
								},
							},
							PrivateStandbyHost_N: []*BucketOriginPrivateHost{
								&BucketOriginPrivateHost{
									Host: "1.qq.com",
									CredentialProvider: &BucketOriginCredentialProvider{
										Role: "qcs::cam::uin/123:roleName/name",
									},
								},
								&BucketOriginPrivateHost{
									Host: "2.qq.com",
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
		t.Logf("Bucket.PutOrigin request body: %+v\n", string(bs))
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
