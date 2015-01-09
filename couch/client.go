package couch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// NewClient makes a new couch client given the CouchDB server URL
func NewClient(host string) (*Client, error) {
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &Client{
		Host: u,
	}, nil
}

type Client struct {
	Host *url.URL
}

func (c *Client) Get(path string, query url.Values, result interface{}) error {
	_, err := c.NewRequest("GET", path, query).Send(result)
	return err
}

func (c *Client) Head(path string, query url.Values, result interface{}) error {
	_, err := c.NewRequest("HEAD", path, query).Send(result)
	return err
}

func (c *Client) Delete(path string, query url.Values, result interface{}) error {
	_, err := c.NewRequest("DELETE", path, query).Send(result)
	return err
}

func (c *Client) Post(path string, query url.Values, payload interface{}, result interface{}) error {
	req := c.NewRequest("POST", path, query)
	req.Payload = payload
	_, err := req.Send(result)
	return err
}

func (c *Client) Put(path string, query url.Values, payload interface{}, result interface{}) error {
	req := c.NewRequest("PUT", path, query)
	req.Payload = payload
	_, err := req.Send(result)
	return err
}

// NewRequest creates a Request given the HTTP verb, path, and optional query values. Additional data
// can be added to the Request after creation, such as Payload
func (c *Client) NewRequest(method string, path string, query url.Values) *Request {
	u := *c.Host
	u.Path = path
	u.RawQuery = query.Encode()

	return &Request{
		Method: method,
		URL:    &u,
	}
}

// Request is an intermediate object for creating CouchDB requests
type Request struct {
	Method  string
	URL     *url.URL
	Payload interface{}
}

// Send performs the HTTP action and parses the JSON body into result
func (r *Request) Send(result interface{}) (*http.Response, error) {
	req, err := r.HTTPRequest()
	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		err := Error{
			Err:    resp.Status,
			Status: resp.StatusCode,
		}
		dec.Decode(&err) // ignore any decoding errors, in case there's no JSON error
		return resp, err
	}
	if result != nil {
		err = dec.Decode(result)
	}
	return resp, err
}

// HTTPRequest creates a http.Request suitable for use with http.Client.Do
func (r *Request) HTTPRequest() (*http.Request, error) {
	var body io.Reader
	if r.Payload != nil {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(r.Payload); err != nil {
			return nil, err
		}
		body = &buf
	}
	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.URL.User != nil {
		pw, _ := r.URL.User.Password()
		req.SetBasicAuth(r.URL.User.Username(), pw)
	}
	return req, nil
}

// Error is a CouchDB JSON error message
// http://couchdb.readthedocs.org/en/latest/api/basics.html#request-format-and-responses
type Error struct {
	Err    string `json:"error"`
	Reason string `json:"reason"`
	Status int    `json:"-"`
}

// Error returns a string, conforming to the error interface
func (e Error) Error() string {
	return fmt.Sprintf("error: %s; reason: %s; HTTP status: %d", e.Err, e.Reason, e.Status)
}
