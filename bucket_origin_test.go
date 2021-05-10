package cos

import (
	"context"
	"encoding/xml"
	"fmt"
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
                <FollowQueryString>true</FollowQueryString>
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
					FollowQueryString: true,
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
					FollowRedirection: true,
					HttpRedirectCode:  "302",
				},
				OriginInfo: &BucketOriginInfo{
					HostInfo: "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
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
					FollowQueryString: true,
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
					FollowRedirection: true,
					HttpRedirectCode:  "302",
				},
				OriginInfo: &BucketOriginInfo{
					HostInfo: "examplebucket-1250000000.cos.ap-shanghai.myqcloud.com",
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
		xml.NewDecoder(r.Body).Decode(body)
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
