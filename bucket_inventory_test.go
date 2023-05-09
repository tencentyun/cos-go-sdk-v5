package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_PutInventory(t *testing.T) {
	setup()
	defer teardown()
	opt := &BucketPutInventoryOptions{
		XMLName:                xml.Name{Local: "InventoryConfiguration"},
		ID:                     "list1",
		IsEnabled:              "True",
		IncludedObjectVersions: "All",
		Filter:                 &BucketInventoryFilter{
			Prefix: "myPrefix",
			Period: nil,
		},
		Schedule:               &BucketInventorySchedule{"Daily"},
		Destination: &BucketInventoryDestination{
			Bucket:     "qcs::cos:ap-guangzhou::examplebucket-1250000000",
			AccountId:  "100000000001",
			Prefix:     "list1",
			Format:     "CSV",
			Encryption: &BucketInventoryEncryption{},
		},
		OptionalFields: &BucketInventoryOptionalFields{
			BucketInventoryFields: []string{
				"Size",
				"LastModifiedDate",
				"ETag",
				"StorageClass",
				"IsMultipartUploaded",
				"ReplicationStatus",
			},
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"inventory": "",
			"id":        "list1",
		}
		testFormValues(t, r, vs)

		body := &BucketPutInventoryOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PutInventory request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutInventory(context.Background(), "list1", opt)
	if err != nil {
		t.Fatalf("Bucket.PutInventory failed, error: %v", err)
	}
}

func TestBucketService_GetInventory(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"inventory": "",
			"id":        "list1",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<InventoryConfiguration>
    <Id>list1</Id>
    <IsEnabled>True</IsEnabled>
    <Destination>
        <COSBucketDestination>
            <Format>CSV</Format>
            <Bucket>qcs::cos:ap-guangzhou::examplebucket-1250000000</Bucket>
            <Prefix>list1</Prefix>
            <AccountId>100000000001</AccountId>
        </COSBucketDestination>
    </Destination>
    <Schedule>
        <Frequency>Daily</Frequency>
    </Schedule>
    <Filter>
        <And>
			<Prefix>myPrefix</Prefix>
		</And>
    </Filter>
    <IncludedObjectVersions>All</IncludedObjectVersions>
    <OptionalFields>
        <Field>Size</Field>
        <Field>LastModifiedDate</Field>
        <Field>ETag</Field>
        <Field>StorageClass</Field>
        <Field>IsMultipartUploaded</Field>
        <Field>ReplicationStatus</Field>
    </OptionalFields>
</InventoryConfiguration>`)
	})
	res, _, err := client.Bucket.GetInventory(context.Background(), "list1")
	if err != nil {
		t.Fatalf("Bucket.GetInventory failed, error: %v", err)
	}
	want := &BucketGetInventoryResult{
		XMLName:                xml.Name{Local: "InventoryConfiguration"},
		ID:                     "list1",
		IsEnabled:              "True",
		IncludedObjectVersions: "All",
		Filter:                 &BucketInventoryFilter{
			Prefix: "myPrefix",
			Period: nil,
		},
		Schedule:               &BucketInventorySchedule{"Daily"},
		Destination: &BucketInventoryDestination{
			Bucket:    "qcs::cos:ap-guangzhou::examplebucket-1250000000",
			AccountId: "100000000001",
			Prefix:    "list1",
			Format:    "CSV",
		},
		OptionalFields: &BucketInventoryOptionalFields{
			BucketInventoryFields: []string{
				"Size",
				"LastModifiedDate",
				"ETag",
				"StorageClass",
				"IsMultipartUploaded",
				"ReplicationStatus",
			},
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetInventory returned\n%+v, want\n%+v", res, want)
	}
}

func TestBucketService_ListInventory(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"inventory": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<ListInventoryConfigurationResult>
    <InventoryConfiguration>
        <Id>list1</Id>
        <IsEnabled>True</IsEnabled>
        <Destination>
            <COSBucketDestination>
                <Format>CSV</Format>
                <AccountId>1250000000</AccountId>
                <Bucket>qcs::cos:ap-beijing::examplebucket-1250000000</Bucket>
                <Prefix>list1</Prefix>
                <Encryption>
                    <SSE-COS/>
                </Encryption>
            </COSBucketDestination>
        </Destination>
        <Schedule>
            <Frequency>Daily</Frequency>
        </Schedule>
        <Filter>
			<And>
				<Prefix>myPrefix</Prefix>
			</And>
        </Filter>
        <IncludedObjectVersions>All</IncludedObjectVersions>
        <OptionalFields>
            <Field>Size</Field>
            <Field>LastModifiedDate</Field>
            <Field>ETag</Field>
            <Field>StorageClass</Field>
            <Field>IsMultipartUpload</Field>
            <Field>ReplicationStatus</Field>
        </OptionalFields>
    </InventoryConfiguration>
    <InventoryConfiguration>
        <Id>list2</Id>
        <IsEnabled>True</IsEnabled>
        <Destination>
            <COSBucketDestination>
                <Format>CSV</Format>
                <AccountId>1250000000</AccountId>
                <Bucket>qcs::cos:ap-beijing::examplebucket-1250000000</Bucket>
                <Prefix>list2</Prefix>
            </COSBucketDestination>
        </Destination>
        <Schedule>
            <Frequency>Weekly</Frequency>
        </Schedule>
        <Filter>
            <And><Prefix>myPrefix2</Prefix></And>
        </Filter>
        <IncludedObjectVersions>All</IncludedObjectVersions>
        <OptionalFields>
            <Field>Size</Field>
            <Field>LastModifiedDate</Field>
            <Field>ETag</Field>
            <Field>StorageClass</Field>
        </OptionalFields>
    </InventoryConfiguration>
    <IsTruncated>false</IsTruncated>
    <ContinuationToken>...</ContinuationToken>
    <IsTruncated>true</IsTruncated>
    <NextContinuationToken>1ueSDFASDF1Tr/XDAFdadEADadf2J/wm36Hy4vbOwM=</NextContinuationToken>
</ListInventoryConfigurationResult>`)
	})

	res, _, err := client.Bucket.ListInventoryConfigurations(context.Background(), "")
	if err != nil {
		t.Fatalf("Bucket.ListInventory failed, error: %v", err)
	}
	want := &ListBucketInventoryConfigResult{
		XMLName:               xml.Name{Local: "ListInventoryConfigurationResult"},
		IsTruncated:           true,
		ContinuationToken:     "...",
		NextContinuationToken: "1ueSDFASDF1Tr/XDAFdadEADadf2J/wm36Hy4vbOwM=",
		InventoryConfigurations: []BucketListInventoryConfiguartion{
			BucketListInventoryConfiguartion{
				XMLName:                xml.Name{Local: "InventoryConfiguration"},
				ID:                     "list1",
				IsEnabled:              "True",
				IncludedObjectVersions: "All",
				Filter:                 &BucketInventoryFilter{
					Prefix: "myPrefix",
					Period: nil,
				},
				Schedule:               &BucketInventorySchedule{"Daily"},
				Destination: &BucketInventoryDestination{
					Bucket:     "qcs::cos:ap-beijing::examplebucket-1250000000",
					AccountId:  "1250000000",
					Prefix:     "list1",
					Format:     "CSV",
					Encryption: &BucketInventoryEncryption{},
				},
				OptionalFields: &BucketInventoryOptionalFields{
					BucketInventoryFields: []string{
						"Size",
						"LastModifiedDate",
						"ETag",
						"StorageClass",
						"IsMultipartUpload",
						"ReplicationStatus",
					},
				},
			},
			BucketListInventoryConfiguartion{
				XMLName:                xml.Name{Local: "InventoryConfiguration"},
				ID:                     "list2",
				IsEnabled:              "True",
				IncludedObjectVersions: "All",
				Filter:                 &BucketInventoryFilter{
					Prefix:"myPrefix2",
					Period: nil,
				},
				Schedule:               &BucketInventorySchedule{"Weekly"},
				Destination: &BucketInventoryDestination{
					Bucket:    "qcs::cos:ap-beijing::examplebucket-1250000000",
					AccountId: "1250000000",
					Prefix:    "list2",
					Format:    "CSV",
				},
				OptionalFields: &BucketInventoryOptionalFields{
					BucketInventoryFields: []string{
						"Size",
						"LastModifiedDate",
						"ETag",
						"StorageClass",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(res, want) {
		t.Fatalf("Bucket.ListInventory failed, \nwant: %+v\nres: %+v", want, res)
	}
}

func TestBucketService_DeleteInventory(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"inventory": "",
			"id":        "list1",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteInventory(context.Background(), "list1")
	if err != nil {
		t.Fatalf("Bucket.DeleteInventory returned error: %v", err)
	}
}

func TestBucketService_PostInventory(t *testing.T) {
	setup()
	defer teardown()
	opt := &BucketPostInventoryOptions{
		XMLName:                xml.Name{Local: "InventoryConfiguration"},
		ID:                     "list1",
		IncludedObjectVersions: "All",
		Filter:                 &BucketInventoryFilter{
			Prefix: "myPrefix",
			Period: nil,
		},
		Destination: &BucketInventoryDestination{
			Bucket:     "qcs::cos:ap-guangzhou::examplebucket-1250000000",
			AccountId:  "100000000001",
			Prefix:     "list1",
			Format:     "CSV",
			Encryption: &BucketInventoryEncryption{},
		},
		OptionalFields: &BucketInventoryOptionalFields{
			BucketInventoryFields: []string{
				"Size",
				"LastModifiedDate",
				"ETag",
				"StorageClass",
				"IsMultipartUploaded",
				"ReplicationStatus",
			},
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		vs := values{
			"inventory": "",
			"id":        "list1",
		}
		testFormValues(t, r, vs)

		body := &BucketPostInventoryOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PostInventory request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PostInventory(context.Background(), "list1", opt)
	if err != nil {
		t.Fatalf("Bucket.PostInventory failed, error: %v", err)
	}
}


