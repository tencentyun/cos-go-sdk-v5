package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetLifecycle(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		vs := values{
			"lifecycle": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<LifecycleConfiguration>
	<Rule>
		<ID>1234</ID>
		<Filter>
			<Prefix>test</Prefix>
		</Filter>
		<Status>Enabled</Status>
		<Transition>
			<Days>10</Days>
			<StorageClass>Standard</StorageClass>
		</Transition>
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
	})

	ref, _, err := client.Bucket.GetLifecycle(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetLifecycle returned error: %v", err)
	}

	want := &BucketGetLifecycleResult{
		XMLName: xml.Name{Local: "LifecycleConfiguration"},
		Rules: []BucketLifecycleRule{
			{
				ID:         "1234",
				Filter:     &BucketLifecycleFilter{Prefix: "test"},
				Status:     "Enabled",
				Transition: &BucketLifecycleTransition{Days: 10, StorageClass: "Standard"},
			},
			{
				ID:         "123422",
				Filter:     &BucketLifecycleFilter{Prefix: "gg"},
				Status:     "Disabled",
				Expiration: &BucketLifecycleExpiration{Days: 10},
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.GetLifecycle returned %+v, want %+v", ref, want)
	}
}

func TestBucketService_PutLifecycle(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutLifecycleOptions{
		Rules: []BucketLifecycleRule{
			{
				ID:         "1234",
				Filter:     &BucketLifecycleFilter{Prefix: "test"},
				Status:     "Enabled",
				Transition: &BucketLifecycleTransition{Days: 10, StorageClass: "Standard"},
			},
			{
				ID:         "123422",
				Filter:     &BucketLifecycleFilter{Prefix: "gg"},
				Status:     "Disabled",
				Expiration: &BucketLifecycleExpiration{Days: 10},
			},
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutLifecycleOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPut)
		vs := values{
			"lifecycle": "",
		}
		testFormValues(t, r, vs)

		want := opt
		want.XMLName = xml.Name{Local: "LifecycleConfiguration"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Bucket.PutLifecycle request body: %+v, want %+v", v, want)
		}

	})

	_, err := client.Bucket.PutLifecycle(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutLifecycle returned error: %v", err)
	}

}

func TestBucketService_DeleteLifecycle(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"lifecycle": "",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteLifecycle(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteLifecycle returned error: %v", err)
	}

}
