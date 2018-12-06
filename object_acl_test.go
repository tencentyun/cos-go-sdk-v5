package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestObjectService_GetACL(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"acl": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<AccessControlPolicy>
	<Owner>
		<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
		<DisplayName>qcs::cam::uin/100000760461:uin/100000760461</DisplayName>
	</Owner>
	<AccessControlList>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="RootAccount">
				<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
				<DisplayName>qcs::cam::uin/100000760461:uin/100000760461</DisplayName>
			</Grantee>
			<Permission>FULL_CONTROL</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="RootAccount">
				<ID>qcs::cam::uin/100000760461:uin/100000760461</ID>
				<DisplayName>qcs::cam::uin/100000760461:uin/100000760461</DisplayName>
			</Grantee>
			<Permission>READ</Permission>
		</Grant>
	</AccessControlList>
</AccessControlPolicy>`)
	})

	ref, _, err := client.Object.GetACL(context.Background(), name)
	if err != nil {
		t.Fatalf("Object.GetACL returned error: %v", err)
	}

	want := &ObjectGetACLResult{
		XMLName: xml.Name{Local: "AccessControlPolicy"},
		Owner: &Owner{
			ID:          "qcs::cam::uin/100000760461:uin/100000760461",
			DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
		},
		AccessControlList: []ACLGrant{
			{
				Grantee: &ACLGrantee{
					Type:        "RootAccount",
					ID:          "qcs::cam::uin/100000760461:uin/100000760461",
					DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
				},
				Permission: "FULL_CONTROL",
			},
			{
				Grantee: &ACLGrantee{
					Type:        "RootAccount",
					ID:          "qcs::cam::uin/100000760461:uin/100000760461",
					DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
				},
				Permission: "READ",
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.GetACL returned %+v, want %+v", ref, want)
	}
}

func TestObjectService_PutACL_with_header_opt(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutACLOptions{
		Header: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, http.MethodPut)
		vs := values{
			"acl": "",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "x-cos-acl", "private")

		want := 0
		v, _ := r.Body.Read([]byte{})
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.PutACL request body: %#v, want %#v", v, want)
		}
	})

	_, err := client.Object.PutACL(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.PutACL returned error: %v", err)
	}

}

func TestObjectService_PutACL_with_body_opt(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutACLOptions{
		Body: &ACLXml{
			Owner: &Owner{
				ID:          "qcs::cam::uin/100000760461:uin/100000760461",
				DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
			},
			AccessControlList: []ACLGrant{
				{
					Grantee: &ACLGrantee{
						Type:        "RootAccount",
						ID:          "qcs::cam::uin/100000760461:uin/100000760461",
						DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
					},

					Permission: "FULL_CONTROL",
				},
				{
					Grantee: &ACLGrantee{
						Type:        "RootAccount",
						ID:          "qcs::cam::uin/100000760461:uin/100000760461",
						DisplayName: "qcs::cam::uin/100000760461:uin/100000760461",
					},
					Permission: "READ",
				},
			},
		},
	}
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		v := new(ACLXml)
		xml.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPut)
		vs := values{
			"acl": "",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "x-cos-acl", "")

		want := opt.Body
		want.XMLName = xml.Name{Local: "AccessControlPolicy"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.PutACL request body: %+v, want %+v", v, want)
		}

	})

	_, err := client.Object.PutACL(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.PutACL returned error: %v", err)
	}

}
