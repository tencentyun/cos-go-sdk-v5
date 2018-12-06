package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetTagging(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"tagging": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<Tagging>
	<TagSet>
		<Tag>
			<Key>test_k2</Key>
			<Value>test_v2</Value>
		</Tag>
		<Tag>
			<Key>test_k3</Key>
			<Value>test_vv</Value>
		</Tag>
	</TagSet>
</Tagging>`)
	})

	ref, _, err := client.Bucket.GetTagging(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetTagging returned error: %v", err)
	}

	want := &BucketGetTaggingResult{
		XMLName: xml.Name{Local: "Tagging"},
		TagSet: []BucketTaggingTag{
			{"test_k2", "test_v2"},
			{"test_k3", "test_vv"},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Bucket.GetTagging returned %+v, want %+v", ref, want)
	}
}

func TestBucketService_PutTagging(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutTaggingOptions{
		TagSet: []BucketTaggingTag{
			{
				Key:   "test_k2",
				Value: "test_v2",
			},
			{
				Key:   "test_k3",
				Value: "test_v3",
			},
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := new(BucketPutTaggingOptions)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		vs := values{
			"tagging": "",
		}
		testFormValues(t, r, vs)

		want := opt
		want.XMLName = xml.Name{Local: "Tagging"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Bucket.PutTagging request body: %+v, want %+v", v, want)
		}

	})

	_, err := client.Bucket.PutTagging(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutTagging returned error: %v", err)
	}

}

func TestBucketService_DeleteTagging(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"tagging": "",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteTagging(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteTagging returned error: %v", err)
	}

}
