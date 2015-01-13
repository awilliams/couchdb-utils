package couch

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// AllDBs fetches the list of databases on the server
func (c *Client) AllDBs() (Databases, error) {
	var result []string
	if err := c.Get("_all_dbs", nil, &result); err != nil {
		return nil, err
	}
	dbs := make(Databases, len(result), len(result))
	for i, db := range result {
		dbs[i] = NewDatabase(c, db)
	}
	sort.Sort(dbs)
	return dbs, nil
}

// NewDatabase receives a client and database name and returns a Database
func NewDatabase(c *Client, name string) *Database {
	u := *c.Host
	u.Path += name
	return &Database{
		Name:   name,
		Client: c,
		URL:    &u,
	}
}

// DatabasesFromURI processes a database URI (either absolute or relative) and creates a Database
func NewDatabaseFromURI(uri string) (*Database, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	path := strings.Trim(u.Path, "/")
	if path == "" {
		return nil, fmt.Errorf("invalid database path from url '%s'", uri)
	}
	u.Path = ""
	client, err := NewClient(u.String())
	if err != nil {
		return nil, err
	}
	return NewDatabase(client, path), nil
}

type Database struct {
	Name   string
	Client *Client
	URL    *url.URL
}

func (d *Database) String() string {
	return d.Name
}

func (d *Database) URI() string {
	if d.URL.IsAbs() {
		return d.URL.String()
	} else {
		return d.URL.Path
	}
}

type Databases []*Database

func (a Databases) Len() int           { return len(a) }
func (a Databases) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Databases) Less(i, j int) bool { return a[i].Name < a[j].Name }
