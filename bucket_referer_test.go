package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetReferer(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"referer": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<RefererConfiguration>
    <Status>Enabled</Status>
    <RefererType>White-List</RefererType>
    <DomainList>
        <Domain>*.qq.com</Domain>
        <Domain>*.qcloud.com</Domain>
    </DomainList>
    <EmptyReferConfiguration>Allow</EmptyReferConfiguration>
</RefererConfiguration>`)
	})

	res, _, err := client.Bucket.GetReferer(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetReferer returned error %v", err)
	}

	want := &BucketGetRefererResult{
		XMLName:     xml.Name{Local: "RefererConfiguration"},
		Status:      "Enabled",
		RefererType: "White-List",
		DomainList: []string{
			"*.qq.com",
			"*.qcloud.com",
		},
		EmptyReferConfiguration: "Allow",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetReferer returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutReferer(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutRefererOptions{
		Status:      "Enabled",
		RefererType: "White-List",
		DomainList: []string{
			"*.qq.com",
			"*.qcloud.com",
		},
		EmptyReferConfiguration: "Allow",
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"referer": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutRefererOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "RefererConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutReferer request\n body: %+v\nwant %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutReferer(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutReferer returned error: %v", err)
	}
}
