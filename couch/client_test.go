package couch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type TestRespData struct {
	Hi    string `json:"hi"`
	Check bool   `json:"check"`
}

type TestPayload struct {
	A string
	B bool
}

func TestClientSend(t *testing.T) {
	path := "/a/b/c"
	query := url.Values{}
	query.Set("foo", "bar")
	query.Set("test", "1")

	body := TestRespData{"hola", true}

	srv, client := mockCouchServer(t, path, query, "GET", func(w http.ResponseWriter, req *http.Request) {
		b, _ := json.Marshal(body)
		w.Write(b)
	})
	defer srv.Close()

	req := client.NewRequest("GET", path, query)

	var respData TestRespData
	resp, err := req.Send(&respData)
	if err != nil {
		t.Error(err)
	}
	if 200 != resp.StatusCode {
		t.Errorf("Expected HTTP %d, recieved %d", 200, resp.StatusCode)
	}
	if body != respData {
		t.Errorf("Expected payload %#v, received %#v", body, respData)
	}
}

func TestClientSendPayload(t *testing.T) {
	path := "/a/b/c"
	query := url.Values{}
	query.Set("foo", "bar")
	query.Set("test", "1")

	body := TestRespData{"bye", true}
	payload := TestPayload{"something", true}

	srv, client := mockCouchServer(t, path, query, "GET", func(w http.ResponseWriter, req *http.Request) {
		var p TestPayload
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&p); err != nil {
			t.Error(err)
		}
		if payload != p {
			t.Errorf("Expected payload %#v, recieved %#v", payload, p)
		}
		b, _ := json.Marshal(body)
		w.Write(b)
	})
	defer srv.Close()

	req := client.NewRequest("GET", path, query)
	req.Payload = payload

	var respData TestRespData
	resp, err := req.Send(&respData)
	if err != nil {
		t.Error(err)
	}
	if 200 != resp.StatusCode {
		t.Errorf("Expected HTTP %d, recieved %d", 200, resp.StatusCode)
	}
	if body != respData {
		t.Errorf("Expected payload %#v, received %#v", body, respData)
	}
}

func TestClientSendJSONError(t *testing.T) {
	path := "/a/b/c"
	query := url.Values{}
	query.Set("foo", "bar")
	query.Set("test", "1")
	respStatus := 400

	respError := Error{
		Err:    "failed",
		Reason: "no idea",
		Status: respStatus,
	}

	srv, client := mockCouchServer(t, path, query, "GET", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(respStatus)
		b, _ := json.Marshal(respError)
		w.Write(b)
	})
	defer srv.Close()

	req := client.NewRequest("GET", path, query)

	resp, err := req.Send(nil)
	if err == nil {
		t.Error("Expected error, none given")
	}
	if respStatus != resp.StatusCode {
		t.Errorf("Expected HTTP %d, recieved %d", 200, resp.StatusCode)
	}
	if respError != err {
		t.Errorf("Expected error payload %#v, received %#v", respError, err)
	}
}

func TestClientSendError(t *testing.T) {
	path := "/a/b/c"
	query := url.Values{}
	respStatus := 400

	respError := Error{
		Err:    fmt.Sprintf("%d %s", respStatus, http.StatusText(respStatus)),
		Status: respStatus,
	}

	srv, client := mockCouchServer(t, path, query, "GET", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(respStatus)
	})
	defer srv.Close()

	req := client.NewRequest("GET", path, query)

	resp, err := req.Send(nil)
	if err == nil {
		t.Error("Expected error, none given")
	}
	if respStatus != resp.StatusCode {
		t.Errorf("Expected HTTP %d, recieved %d", 200, resp.StatusCode)
	}
	if respError != err {
		t.Errorf("Expected error payload %#v, received %#v", respError, err)
	}
}

func mockCouchServer(t *testing.T, path string, query url.Values, method string, handler http.HandlerFunc) (*httptest.Server, *Client) {
	basicAuthUser := "user"
	basicAuthPw := "secret"
	f := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		m := strings.ToUpper(req.Method)
		if method != m {
			t.Errorf("Expected method %s, received %v", method, m)
		}
		p := req.URL.Path
		if path != p {
			t.Errorf("Expected path '%s', received '%s'", path, p)
		}
		q := req.URL.RawQuery
		if query.Encode() != q {
			t.Errorf("Expected query '%s', received '%s'", query.Encode(), q)
		}
		if user, pw, ok := req.BasicAuth(); !ok {
			t.Error("Expected basic auth, none provided")
		} else if basicAuthUser != user {
			t.Errorf("Expected basic auth user '%s', received '%s'", basicAuthUser, user)
		} else if basicAuthPw != pw {
			t.Errorf("Expected basic auth password '%s', received '%s'", basicAuthPw, pw)
		}
		req.Header.Add("Content-Type", "application/json")
		handler(w, req)
	})
	srv := httptest.NewServer(http.HandlerFunc(f))
	host := fmt.Sprintf("http://%s:%s@%s", basicAuthUser, basicAuthPw, srv.Listener.Addr().String())
	client, err := NewClient(host)
	if err != nil {
		t.Error(err)
	}

	return srv, client
}
