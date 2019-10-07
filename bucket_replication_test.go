package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_PutReplication(t *testing.T) {
	setup()
	defer teardown()
	opt := &PutBucketReplicationOptions{
		Role: "qcs::cam::uin/100000000001:uin/100000000001",
		Rule: []BucketReplicationRule{
			{
				Status: "Disabled",
				Prefix: "prefix",
				Destination: &ReplicationDestination{
					Bucket: "qcs::cos:ap-beijing-1::examplebucket-1250000000",
				},
			},
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		vs := values{
			"replication": "",
		}
		testFormValues(t, r, vs)

		body := &PutBucketReplicationOptions{}
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "ReplicationConfiguration"}
		if !reflect.DeepEqual(want, body) {
			t.Fatalf("Bucket.PutReplication request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutBucketReplication(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutLogging failed, error: %v", err)
	}
}

func TestBucketService_GetReplication(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"replication": "",
		}
		testFormValues(t, r, vs)

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

	})
	res, _, err := client.Bucket.GetBucketReplication(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetReplication failed, error: %v", err)
	}
	want := &GetBucketReplicationResult{
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
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetBucketReplication\nres %+v\nwant %+v", res.Rule[0].Destination, want.Rule[0].Destination)
		t.Errorf("Bucket.GetBucketReplication\nres %+v\nwant %+v", res, want)
	}
}

func TestBucketService_DeleteReplication(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"replication": "",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteBucketReplication(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteBucketReplication returned error: %v", err)
	}
}
