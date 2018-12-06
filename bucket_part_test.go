package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_ListMultipartUploads(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"uploads": "",
			"prefix":  "t",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<ListMultipartUploadsResult>
	<Bucket>test-1253846586</Bucket>
	<Encoding-Type/>
	<KeyMarker/>
	<UploadIdMarker/>
	<MaxUploads>1000</MaxUploads>
	<Prefix>t</Prefix>
	<Delimiter>/</Delimiter>
	<IsTruncated>false</IsTruncated>
	<CommonPrefixs>
		<Prefix>test/</Prefix>
	</CommonPrefixs>
	<Upload>
		<Key>test_multipart.txt</Key>
		<UploadId>14972623850a5de3f4f10605ab9f339c8bdf1b77e06f03fb981e7e76c86554b7bdb6072b36</UploadId>
		<Initiator>
			<ID>100000760461/100000760461</ID>
			<DisplayName/>
		</Initiator>
		<Owner>
			<ID>100000760461/100000760461</ID>
			<DisplayName/>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
		<Initiated>2017-06-12T10:13:05.000Z</Initiated>
	</Upload>
	<Upload>
		<Key>test_multipar2t.txt</Key>
		<UploadId>1497515958744e899fc341bfbb995ebd57b395f63930411d855aaac1b5cd7d834a15442831</UploadId>
		<Initiator>
			<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
			<DisplayName>100000760461</DisplayName>
		</Initiator>
		<Owner>
			<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
			<DisplayName>100000760461</DisplayName>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
		<Initiated>2017-06-15T08:39:18.000Z</Initiated>
	</Upload>
</ListMultipartUploadsResult>`)
	})

	opt := &ListMultipartUploadsOptions{
		Prefix: "t",
	}
	ref, _, err := client.Bucket.ListMultipartUploads(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.ListMultipartUploads returned error: %v", err)
	}

	want := &ListMultipartUploadsResult{
		XMLName:     xml.Name{Local: "ListMultipartUploadsResult"},
		Bucket:      "test-1253846586",
		MaxUploads:  1000,
		IsTruncated: false,
		Uploads: []struct {
			Key          string
			UploadID     string `xml:"UploadId"`
			StorageClass string
			Initiator    *Initiator
			Owner        *Owner
			Initiated    string
		}{
			{
				Key:      "test_multipart.txt",
				UploadID: "14972623850a5de3f4f10605ab9f339c8bdf1b77e06f03fb981e7e76c86554b7bdb6072b36",
				Initiator: &Initiator{
					ID: "100000760461/100000760461",
				},
				Owner: &Owner{
					ID: "100000760461/100000760461",
				},
				StorageClass: "STANDARD",
				Initiated:    "2017-06-12T10:13:05.000Z",
			},
			{
				Key:      "test_multipar2t.txt",
				UploadID: "1497515958744e899fc341bfbb995ebd57b395f63930411d855aaac1b5cd7d834a15442831",
				Initiator: &Initiator{
					ID:          "qcs::cam::uin/100000760461:uin/100000760461",
					DisplayName: "100000760461",
				},
				Owner: &Owner{
					ID:          "qcs::cam::uin/100000760461:uin/100000760461",
					DisplayName: "100000760461",
				},
				StorageClass: "STANDARD",
				Initiated:    "2017-06-15T08:39:18.000Z",
			},
		},
		Prefix:         "t",
		CommonPrefixes: []string{"test/"},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.ListMultipartUploads returned \n%+v, want \n%+v", ref, want)
	}
}
