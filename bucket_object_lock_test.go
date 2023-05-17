package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetObjectLockConfiguration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"object-lock": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<ObjectLockConfiguration>
	<ObjectLockEnabled>Enabled</ObjectLockEnabled> 
	<Rule> 
		<DefaultRetention>
			<Days>30</Days> 
		</DefaultRetention> 
	</Rule> 
</ObjectLockConfiguration>`)
	})

	res, _, err := client.Bucket.GetObjectLockConfiguration(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetObjectLockConfiguration returned error %v", err)
	}

	want := &BucketGetObjectLockResult{
		XMLName:           xml.Name{Local: "ObjectLockConfiguration"},
		ObjectLockEnabled: "Enabled",
		Rule: &ObjectLockRule{
			Days: 30,
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetObjectLockConfiguration returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutObjectLockConfiguration(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutObjectLockOptions{
		ObjectLockEnabled: "Enabled",
		Rule: &ObjectLockRule{
			Days: 30,
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"object-lock": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutObjectLockOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "ObjectLockConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutObjectLockConfiguration request\n body: %+v\nwant %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutObjectLockConfiguration(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutObjectLockConfiguration returned error: %v", err)
	}
}

func TestBucketService_GetRetention(t *testing.T) {
	setup()
	defer teardown()

	key := "example"
	mux.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"retention": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<Retention>
	<RetainUntilDate>2023-05-18T07:19:55.000Z</RetainUntilDate> 
</Retention>`)
	})

	res, _, err := client.Object.GetRetention(context.Background(), key, nil)
	if err != nil {
		t.Fatalf("Object.GetRetention returned error: %v", err)
	}

	want := &ObjectGetRetentionResult{
		XMLName:         xml.Name{Local: "Retention"},
		RetainUntilDate: "2023-05-18T07:19:55.000Z",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Object.GetRetention returned %+v, want %+v", res, want)
	}
}
