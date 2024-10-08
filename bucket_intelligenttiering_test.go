package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_PutIntelligentTiering(t *testing.T) {
	setup()
	defer teardown()
	opt := &BucketPutIntelligentTieringOptions{
		Status: "Enabled",
		Transition: &BucketIntelligentTieringTransition{
			Days: 30,
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"intelligenttiering": "",
		}
		testFormValues(t, r, vs)

		body := &BucketPutIntelligentTieringOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "IntelligentTieringConfiguration"}
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PutIntelligentTiering request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutIntelligentTiering(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutIntelligentTiering failed, error: %v", err)
	}
}

func TestBucketService_GetIntelligentTiering(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"intelligenttiering": "",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "x-cos-meta-test", "test")

		fmt.Fprint(w, `<IntelligentTieringConfiguration>
            <Status>Enabled</Status>
            <Transition>
                <Days>30</Days>
            </Transition>
        </IntelligentTieringConfiguration>`)
	})
	opt := &BucketGetIntelligentTieringOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "test")
	res, _, err := client.Bucket.GetIntelligentTiering(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.GetIntelligentTiering failed, error: %v", err)
	}
	want := &BucketGetIntelligentTieringResult{
		XMLName: xml.Name{Local: "IntelligentTieringConfiguration"},
		Status:  "Enabled",
		Transition: &BucketIntelligentTieringTransition{
			Days: 30,
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetIntelligentTiering returned\n%+v, want\n%+v", res, want)
	}
}

func TestBucketService_PutIntelligentTieringV2(t *testing.T) {
	setup()
	defer teardown()
	opt := &BucketPutIntelligentTieringOptions{
		Status: "Enabled",
		Transition: &BucketIntelligentTieringTransition{
			Days: 30,
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"intelligent-tiering": "",
			"id":                  "test",
		}
		testFormValues(t, r, vs)

		body := &BucketPutIntelligentTieringOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "IntelligentTieringConfiguration"}
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PutIntelligentTieringV2 request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutIntelligentTieringV2(context.Background(), opt)
	if err == nil || err.Error() != "id is empty" {
		t.Fatalf("Bucket.PutIntelligentTieringV2 failed, error: %v", err)
	}
	opt.Id = "test"
	_, err = client.Bucket.PutIntelligentTieringV2(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutIntelligentTieringV2 failed, error: %v", err)
	}
}

func TestBucketService_GetIntelligentTieringV2(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"intelligent-tiering": "",
			"id":                  "test",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "x-cos-meta-test", "test")

		fmt.Fprint(w, `<IntelligentTieringConfiguration>
            <Status>Enabled</Status>
            <Transition>
                <Days>30</Days>
            </Transition>
        </IntelligentTieringConfiguration>`)
	})
	opt := &BucketGetIntelligentTieringOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("x-cos-meta-test", "test")
	res, _, err := client.Bucket.GetIntelligentTieringV2(context.Background(), "test", opt)
	if err != nil {
		t.Fatalf("Bucket.GetIntelligentTiering failed, error: %v", err)
	}
	want := &BucketGetIntelligentTieringResult{
		XMLName: xml.Name{Local: "IntelligentTieringConfiguration"},
		Status:  "Enabled",
		Transition: &BucketIntelligentTieringTransition{
			Days: 30,
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetIntelligentTiering returned\n%+v, want\n%+v", res, want)
	}
}

func TestBucketService_ListIntelligentTiering(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"intelligent-tiering": "",
		}
		testFormValues(t, r, vs)
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
                <Tag>
                    <Key>k2</Key>
                    <Value>v2</Value>
                </Tag>
            </And>
        </Filter>
        <Tiering>
            <AccessTier>ARCHIVE_ACCESS</AccessTier>
            <Days>30</Days>
        </Tiering>
        <Tiering>
            <AccessTier>DEEP_ARCHIVE_ACCESS</AccessTier>
            <Days>30</Days>
        </Tiering>
    </IntelligentTieringConfiguration>
</ListBucketIntelligentTieringConfigurationsOutput>`)
	})
	res, _, err := client.Bucket.ListIntelligentTiering(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetIntelligentTiering failed, error: %v", err)
	}
	want := &ListIntelligentTieringConfigurations{
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
							{
								Key:   "k2",
								Value: "v2",
							},
						},
					},
				},
				Tiering: []*BucketIntelligentTieringTransition{
					{
						AccessTier: "ARCHIVE_ACCESS",
						Days:       30,
					},
					{
						AccessTier: "DEEP_ARCHIVE_ACCESS",
						Days:       30,
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetIntelligentTiering returned\n%+v, want\n%+v", res, want)
	}

}

func TestBucketService_DeleteIntelligentTiering(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"intelligent-tiering": "",
			"id":                  "test",
		}
		testFormValues(t, r, vs)
	})
	_, err := client.Bucket.DeleteIntelligentTiering(context.Background(), "test")
	if err != nil {
		t.Fatalf("Bucket.GetIntelligentTiering failed, error: %v", err)
	}
}
