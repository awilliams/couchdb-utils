package api

import (
	"encoding/json"
	"fmt"
	"github.com/awilliams/couchdb-utils/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	HEAD            = "HEAD"
	PUT             = "PUT"
	POST            = "POST"
	GET             = "GET"
	DELETE          = "DELETE"
	JSONCONTENTTYPE = "application/json"
)

type Couchdb struct {
	_url          string
	ResultHandler func(*Result)
}

type CouchdbError struct {
	ErrorMsg string `json:"error"`
	Status   int
	Reason   string
}

func (e CouchdbError) Error() string {
	var err string
	if e.ErrorMsg != "" && e.Reason != "" {
		err = fmt.Sprintf("%s: %s", e.ErrorMsg, e.Reason)
	}
	if e.Status != 0 {
		err += fmt.Sprintf("\nHTTP %d", e.Status)
	}
	return err
}

func (e CouchdbError) IsConflict() bool {
	return e.Status == 409
}

func (e CouchdbError) IsNotFound() bool {
	return e.Status == 404
}

type Result struct {
	Path       string
	Method     string
	StatusCode int
}

func (r Result) PP(printer util.Printer) {
	printer.Print("# %s %s %d", r.Method, r.Path, r.StatusCode)
}

func parseHost(rawurl string) (string, error) {
	if !strings.HasPrefix(rawurl, "http") {
		rawurl = "http://" + rawurl
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return u.String() + "/", nil
}

func New(url string) (*Couchdb, error) {
	host, err := parseHost(url)
	if err != nil {
		return nil, err
	}
	return &Couchdb{_url: host}, nil
}

func (c *Couchdb) url(pathComponents ...string) string {
	if len(pathComponents) == 0 {
		return c._url
	} else {
		return c._url + strings.Join(pathComponents, "/")
	}
}

func (c *Couchdb) _perform(method string, bodyType string, body io.Reader, path string) (io.ReadCloser, error) {
	path = c.url(path)
	result := Result{Method: method, Path: sanitizePath(path)}
	if c.ResultHandler != nil {
		defer c.ResultHandler(&result)
	}
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	result.StatusCode = resp.StatusCode
	return handleResponse(resp)
}

func (c *Couchdb) get(path string) (io.ReadCloser, error) {
	return c._perform(GET, "", nil, path)
}

func (c *Couchdb) del(path string) (io.ReadCloser, error) {
	return c._perform(DELETE, "", nil, path)
}

func (c *Couchdb) post(bodyType string, body io.Reader, path string) (io.ReadCloser, error) {
	return c._perform(POST, bodyType, body, path)
}

func (c *Couchdb) put(bodyType string, body io.Reader, path string) (io.ReadCloser, error) {
	return c._perform(PUT, bodyType, body, path)
}

func (c *Couchdb) head(path string) (io.ReadCloser, error) {
	return c._perform(HEAD, "", nil, path)
}

func (c *Couchdb) getJson(jsontype interface{}, path string) error {
	body, err := c.get(path)
	if err != nil {
		return err
	}
	return parseJson(body, jsontype)
}

func (c *Couchdb) deleteJson(jsontype interface{}, path string) error {
	body, err := c.del(path)
	if err != nil {
		return err
	}
	return parseJson(body, jsontype)
}

func (c *Couchdb) postJson(jsontype interface{}, b io.Reader, path string) error {
	body, err := c.post(JSONCONTENTTYPE, b, path)
	if err != nil {
		return err
	}
	return parseJson(body, jsontype)
}

func (c *Couchdb) putJson(jsontype interface{}, b io.Reader, path string) error {
	body, err := c.put(JSONCONTENTTYPE, b, path)
	if err != nil {
		return err
	}
	return parseJson(body, jsontype)
}

func parseJson(body io.ReadCloser, jsontype interface{}) error {
	if body != nil {
		defer body.Close()
	}
	decoder := json.NewDecoder(body)
	return decoder.Decode(jsontype)
}

func handleResponse(resp *http.Response) (io.ReadCloser, error) {
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		var errObj CouchdbError = CouchdbError{Status: resp.StatusCode}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&errObj) // ignore errors
		return nil, errObj
	}
	return resp.Body, nil
}

func sanitizePath(path string) string {
	return sanitizePathRegex.ReplaceAllString(path, "$1:***@$3")
}

var sanitizePathRegex *regexp.Regexp

func init() {
	sanitizePathRegex = regexp.MustCompile("^(.*):(.)+@(.*)$")
}
