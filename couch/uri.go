package couch

import (
	"net/url"
	"strings"
)

// URI is used to identify CouchDB databases, documents, views, etc
type URI url.URL

const delimiter = "/"

// NewURI creates a URI from given string (absolute or relative)
func NewURI(s string) (*URI, error) {
	u, err := url.Parse(strings.Trim(s, delimiter))
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		u.Path = delimiter + u.Path
	}
	return (*URI)(u), nil
}

// Client returns a Client and bool indicating if it was successful or not
func (u *URI) Client(logHTTP bool) (*Client, bool) {
	if (*url.URL)(u).IsAbs() {
		c := *u
		c.Path = ""
		return Client{
			Host:    &c,
			LogHTTP: logHTTP,
		}, true
	} else {
		return nil, false
	}
}

// Extend returns a copy of u with addtional path components appended
func (u *URI) Extend(components ...string) *URI {
	// strip any leading or trailing slashes from components
	for i, component := range components {
		components[i] = strings.Trim(component, delimiter)
	}
	// make a copy to return
	c := *u
	if !strings.HasSuffix(c.Path, delimiter) {
		c.Path += delimiter
	}
	c.Path += strings.Join(components, delimiter)
	return &c
}

// Reduce returns a copy of u with suffix stripped from u's path
func (u *URI) Reduce(suffix string) *URI {
	// make a copy to return
	c := *u
	if len(suffix) == 0 || suffix == delimiter {
		return &c
	}
	c.Path = strings.TrimSuffix(c.Path, strings.Trim(suffix, delimiter))
	c.Path = strings.TrimSuffix(c.Path, delimiter)
	return &c
}

// Components returns a slice of path components
func (u *URI) Components() []string {
	if len(u.Path) == 0 || u.Path == delimiter {
		return []string{}
	}
	components := strings.Split(strings.Trim(u.Path, delimiter), delimiter)
	for i, component := range components {
		components[i] = strings.Trim(component, delimiter)
	}
	return components
}

// String returns an absolute or relative URI string
func (u *URI) String() string {
	return (*url.URL)(u).String()
}
