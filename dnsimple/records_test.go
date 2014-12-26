package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestRecord_recordPath(t *testing.T) {
	var pathTest = []struct {
		domainInput interface{}
		recordInput interface{}
		expected    string
	}{
		{"example.com", nil, "domains/example.com/records"},
		{"example.com", 2, "domains/example.com/records/2"},
		{23, nil, "domains/23/records"},
		{23, 2, "domains/23/records/2"},
	}

	for _, pt := range pathTest {
		actual := recordPath(pt.domainInput, pt.recordInput)
		if actual != pt.expected {
			t.Errorf("recordPath(%+v, %+v): expected %s, actual %s", pt.domainInput, pt.recordInput, pt.expected, actual)
		}
	}
}

func TestDomainsService_ListRecords_all(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"record":{"id":1, "name":"foo.example.com"}}]`)
	})

	records, _, err := client.Domains.ListRecords("example.com", "", "")

	if err != nil {
		t.Errorf("Domains.ListRecords returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Domains.ListRecords returned %+v, want %+v", records, want)
	}
}

func TestDomainsService_ListRecords_subdomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "foo"})

		fmt.Fprint(w, `[{"record":{"id":1, "name":"foo.example.com"}}]`)
	})

	records, _, err := client.Domains.ListRecords("example.com", "foo", "")

	if err != nil {
		t.Errorf("Domains.ListRecords returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Domains.ListRecords returned %+v, want %+v", records, want)
	}
}

func TestDomainsService_ListRecords_type(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "foo", "type": "CNAME"})

		fmt.Fprint(w, `[{"record":{"id":1, "name":"foo.example.com"}}]`)
	})

	records, _, err := client.Domains.ListRecords("example.com", "foo", "CNAME")

	if err != nil {
		t.Errorf("Domains.ListRecords returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Domains.ListRecords returned %+v, want %+v", records, want)
	}
}

func TestDomainsService_CreateRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["record"] = map[string]interface{}{"name": "foo", "content": "192.168.0.10", "record_type": "A"}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprintf(w, `{"record":{"id":42, "name":"foo"}}`)
	})

	recordValues := Record{Name: "foo", Content: "192.168.0.10", Type: "A"}
	record, _, err := client.Domains.CreateRecord("example.com", recordValues)

	if err != nil {
		t.Errorf("Domains.CreateRecord returned error: %v", err)
	}

	want := Record{Id: 42, Name: "foo"}
	if !reflect.DeepEqual(record, want) {
		t.Errorf("Domains.CreateRecord returned %+v, want %+v", record, want)
	}
}

func TestDomainsService_GetRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records/1539", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
			{
			  "record": {
				"id": 1539,
				"domain_id": 227,
				"parent_id": null,
				"name": "",
				"content": "mail.example.com",
				"ttl": 3600,
				"prio": 1,
				"record_type": "MX",
				"system_record": null,
				"created_at": "2014-01-15T22:03:03Z",
				"updated_at": "2014-01-15T22:03:04Z"
			  }
			}
		`)
	})

	record, _, err := client.Domains.GetRecord("example.com", 1539)

	if err != nil {
		t.Errorf("Domains.GetRecord returned error: %v", err)
	}

	want := Record{Id: 1539, DomainId: 227, Name: "", Content: "mail.example.com", TTL: 3600, Priority: 1, Type: "MX", CreatedAt: "2014-01-15T22:03:03Z", UpdatedAt: "2014-01-15T22:03:04Z"}
	if !reflect.DeepEqual(record, want) {
		t.Errorf("Domains.GetRecord returned %+v, want %+v", record, want)
	}
}

func TestDomainsService_DeleteRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		// fmt.Fprint(w, `{}`)
	})

	_, err := client.Domains.DeleteRecord("example.com", 2)

	if err != nil {
		t.Errorf("Domains.DeleteRecord returned error: %v", err)
	}
}

func TestDomainsService_DeleteRecord_failed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"Invalid request"}`)
	})

	_, err := client.Domains.DeleteRecord("example.com", 2)
	if err == nil {
		t.Errorf("Domains.DeleteRecord expected error to be returned")
	}

	if match := "400 Invalid request"; !strings.Contains(err.Error(), match) {
		t.Errorf("Records.Delete returned %+v, should match %+v", err, match)
	}
}

func TestRecord_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/24/records/42", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		// TODO: there's a problem when verifying adding prio and ttl integers. Why?
		want["record"] = map[string]interface{}{"content": "192.168.0.10", "name": "bar"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"record":{"id":24, "domain_id":42}}`)
	})

	record := Record{Id: 42, DomainId: 24, Name: "foo"}
	recordAttributes := Record{Name: "bar", Content: "192.168.0.10"}

	_, err := record.Update(client, recordAttributes)

	if err != nil {
		t.Errorf("Update returned error: %v", err)
	}
}

func TestRecord_UpdateIP(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/24/records/42", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["record"] = map[string]interface{}{"name": "foo", "content": "192.168.0.1"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"record":{"id":24, "domain_id":42}}`)
	})

	record := Record{Id: 42, DomainId: 24, Name: "foo"}
	err := record.UpdateIP(client, "192.168.0.1")

	if err != nil {
		t.Errorf("UpdateIP returned error: %v", err)
	}
}
