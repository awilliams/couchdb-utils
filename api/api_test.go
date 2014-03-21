package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestingServer(status int, responses ...string) (*httptest.Server, *Couchdb) {
	if status == 0 {
		status = 200
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		for _, line := range responses {
			fmt.Fprintln(w, line)
		}
	}))
	couchdb, _ := New(ts.URL)
	return ts, couchdb
}

func TestNewUrl(t *testing.T) {
	host := "1.2.3.4:5544"
	couchdb, err := New(host)
	if err != nil {
		t.Fatal(err)
	}
	expectedUrl := fmt.Sprintf("http://%s", host)
	actualUrl := couchdb.url()
	if actualUrl.String() != expectedUrl {
		t.Fatalf("Expected: %s, Received: %s", expectedUrl, actualUrl.String())
	}

	expectedUrl = fmt.Sprintf("http://%s/one/two", host)
	actualUrl = couchdb.url("one", "two")
	if actualUrl.String() != expectedUrl {
		t.Fatalf("Expected: %s, Received: %s", expectedUrl, actualUrl.String())
	}
	
	// test weird database names
	expectedUrl = fmt.Sprintf("http://%s/weird%%2Fdata%%2Bbase%%2Fname", host)
	actualUrl = couchdb.url("weird%2Fdata%2Bbase%2Fname")
	if actualUrl.String() != expectedUrl {
		t.Fatalf("Expected: %s, Received: %s", expectedUrl, actualUrl.String())
	}
}

func TestGet(t *testing.T) {
	body := `{"OK":"Computer"}`
	ts, couchdb := newTestingServer(200, body)
	defer ts.Close()
	resp, err := couchdb.get("")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	respbody, err := ioutil.ReadAll(resp)
	if err != nil {
		t.Fatal(err)
	}
	if string(respbody) != body+"\n" {
		t.Fatalf("Expected: %s, Actual: %s", body, respbody)
	}
}

func TestGet404(t *testing.T) {
	body := `{"error": "no", "reason":"happens"}`
	ts, couchdb := newTestingServer(404, body)
	defer ts.Close()
	resp, err := couchdb.get("")
	if err == nil {
		t.Fatal("Error was expected")
	}
	if resp != nil {
		defer resp.Close()
	}
	switch et := err.(type) {
	case CouchdbError:
		couchErr := err.(CouchdbError)
		if couchErr.Status != 404 {
			t.Fatalf("Status incorrect. Expected: %v, Actual: %v", 404, couchErr.Status)
		}
		if couchErr.ErrorMsg != "no" {
			t.Fatalf("ErrorMsg incorrect. Expected: %v, Actual: %v", "no", couchErr.ErrorMsg)
		}
		if couchErr.Reason != "happens" {
			t.Fatalf("Reason incorrect. Expected: %v, Actual: %v", "happens", couchErr.Reason)
		}
	default:
		t.Fatalf("error of type CouchdbError expected, given %v", et)
	}
}

type testJson struct {
	OK       string
	Computer bool
}

func TestGetJson(t *testing.T) {
	body := `{"OK":"Bill", "computer": false}`
	ts, couchdb := newTestingServer(200, body)
	defer ts.Close()
	var json testJson
	err := couchdb.getJson(&json, "")
	if err != nil {
		t.Fatal(err)
	}
	expected := testJson{"Bill", false}
	if json != expected {
		t.Fatalf("Expected: %#v, Actual: %#v", expected, json)
	}
}

func TestGetInvalidJson(t *testing.T) {
	body := `{OK:"Bill" "computer": false}`
	ts, couchdb := newTestingServer(200, body)
	defer ts.Close()
	var json testJson
	err := couchdb.getJson(&json, "")
	if err == nil {
		t.Fatalf("Should have errored parsing crazy json")
	}
}

func TestResultHandler(t *testing.T) {
	body := ""
	ts, couchdb := newTestingServer(200, body)
	done := false
	couchdb.ResultHandler = func(r *Result) {
		done = true
		if r.Method != "GET" {
			t.Fatalf("Method incorrect, Expected: %s, Actual %s", "GET", r.Method)
		}
		if r.StatusCode != 200 {
			t.Fatalf("StatusCode incorrect, Expected: %s, Actual %s", 200, r.StatusCode)
		}
	}
	couchdb.get("")
	ts.Close()
	if !done {
		t.Fatal("Did not call ResultHandler function!")
	}
}
