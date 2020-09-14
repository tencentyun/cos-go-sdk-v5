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

		fmt.Fprint(w, `<IntelligentTieringConfiguration>
            <Status>Enabled</Status>
            <Transition>
                <Days>30</Days>
            </Transition>
        </IntelligentTieringConfiguration>`)
	})
	res, _, err := client.Bucket.GetIntelligentTiering(context.Background())
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
