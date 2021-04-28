package cos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetPolicy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"policy": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `{
            "Statement": [
            {
                "Principal": {
                    "qcs": [
                    "qcs::cam::uin/100000000001:uin/100000000011"
                    ]
                },
                "Effect": "allow",
                "Action": [
                "name/cos:GetBucket"
                ],
                "Resource": [
                "qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/*"
                ]
            }
            ],
            "version": "2.0"
        }`)
	})

	res, _, err := client.Bucket.GetPolicy(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetPolicy returned error %v", err)
	}

	want := &BucketGetPolicyResult{
		Statement: []BucketStatement{
			{
				Principal: map[string][]string{
					"qcs": []string{"qcs::cam::uin/100000000001:uin/100000000011"},
				},
				Effect:   "allow",
				Action:   []string{"name/cos:GetBucket"},
				Resource: []string{"qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/*"},
			},
		},
		Version: "2.0",
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetPolicy returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutPolicy(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutPolicyOptions{
		Statement: []BucketStatement{
			{
				Principal: map[string][]string{
					"qcs": []string{"qcs::cam::uin/100000000001:uin/100000000011"},
				},
				Effect:   "allow",
				Action:   []string{"name/cos:GetBucket"},
				Resource: []string{"qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/*"},
			},
		},
		Version: "2.0",
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"policy": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutPolicyOptions)
		json.NewDecoder(r.Body).Decode(body)
		want := opt
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutPolicy request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutPolicy(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutPolicy returned error: %v", err)
	}
}

func TestBucketService_DeletePolicy(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"policy": "",
		}
		testFormValues(t, r, vs)
		w.WriteHeader(http.StatusNoContent)
	})
	_, err := client.Bucket.DeletePolicy(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeletePolicy returned error: %v", err)
	}
}
