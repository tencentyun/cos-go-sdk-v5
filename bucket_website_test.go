package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBucketService_GetWebsite(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"website": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<WebsiteConfiguration>
	<IndexDocument>
		<Suffix>index.html</Suffix>
	</IndexDocument>
	<RedirectAllRequestsTo>
		<Protocol>https</Protocol>
	</RedirectAllRequestsTo>
	<RoutingRules>
		<RoutingRule>
			<Condition>
				<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>
			</Condition>
			<Redirect>
				<Protocol>https</Protocol>
				<ReplaceKeyWith>404.html</ReplaceKeyWith>
			</Redirect>
		</RoutingRule>
		<RoutingRule>
			<Condition>
				<KeyPrefixEquals>docs/</KeyPrefixEquals>
			</Condition>
			<Redirect>
				<Protocol>https</Protocol>
				<ReplaceKeyPrefixWith>documents/</ReplaceKeyPrefixWith>
			</Redirect>
		</RoutingRule>
		<RoutingRule>
			<Condition>
				<KeyPrefixEquals>img/</KeyPrefixEquals>
			</Condition>
			<Redirect>
				<Protocol>https</Protocol>
				<ReplaceKeyWith>demo.jpg</ReplaceKeyWith>
			</Redirect>
		</RoutingRule>
	</RoutingRules>
</WebsiteConfiguration>`)
	})

	res, _, err := client.Bucket.GetWebsite(context.Background())
	if err != nil {
		t.Fatalf("Bucket.GetWebsite returned error %v", err)
	}

	want := &BucketGetWebsiteResult{
		XMLName: xml.Name{Local: "WebsiteConfiguration"},
		Index:   "index.html",
		RedirectProtocol: &RedirectRequestsProtocol{
			"https",
		},
		RoutingRules: &WebsiteRoutingRules{
			Rules: []WebsiteRoutingRule{
				{
					ConditionErrorCode: "404",
					RedirectProtocol:   "https",
					RedirectReplaceKey: "404.html",
				},
				{
					ConditionPrefix:          "docs/",
					RedirectProtocol:         "https",
					RedirectReplaceKeyPrefix: "documents/",
				},
				{
					ConditionPrefix:    "img/",
					RedirectProtocol:   "https",
					RedirectReplaceKey: "demo.jpg",
				},
			},
		},
	}

	if !reflect.DeepEqual(res, want) {
		t.Errorf("Bucket.GetWebsite returned %+v, want %+v", res, want)
	}
}

func TestBucketService_PutWebsite(t *testing.T) {
	setup()
	defer teardown()

	opt := &BucketPutWebsiteOptions{
		Index: "index.html",
		RedirectProtocol: &RedirectRequestsProtocol{
			"https",
		},
		Error: &ErrorDocument{
			"Error.html",
		},
		RoutingRules: &WebsiteRoutingRules{
			[]WebsiteRoutingRule{
				{
					ConditionErrorCode: "404",
					RedirectProtocol:   "https",
					RedirectReplaceKey: "404.html",
				},
				{
					ConditionPrefix:          "docs/",
					RedirectProtocol:         "https",
					RedirectReplaceKeyPrefix: "documents/",
				},
				{
					ConditionPrefix:    "img/",
					RedirectProtocol:   "https",
					RedirectReplaceKey: "demo.jpg",
				},
			},
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		vs := values{
			"website": "",
		}
		testFormValues(t, r, vs)

		body := new(BucketPutWebsiteOptions)
		xml.NewDecoder(r.Body).Decode(body)
		want := opt
		want.XMLName = xml.Name{Local: "WebsiteConfiguration"}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("Bucket.PutWebsite request\n body: %+v\n, want %+v\n", body, want)
		}
	})

	_, err := client.Bucket.PutWebsite(context.Background(), opt)
	if err != nil {
		t.Fatalf("Bucket.PutWebsite returned error: %v", err)
	}
}

func TestBucketService_DeleteWebsite(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		vs := values{
			"website": "",
		}
		testFormValues(t, r, vs)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Bucket.DeleteWebsite(context.Background())
	if err != nil {
		t.Fatalf("Bucket.DeleteWebsite returned error: %v", err)
	}

}
