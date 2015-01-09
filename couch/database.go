package couch

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// AllDBs fetches the list of databases on the server
func AllDBs(c *Client) (Databases, error) {
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

func NewDatabase(c *Client, name string) *Database {
	u := *c.Host
	u.Path += name
	return &Database{
		Name:   name,
		Client: c,
		URL:    &u,
	}
}

func DatabasesFromURLs(urls []string) (Databases, error) {
	// TODO: allow for relative urls, eg 'databasename'
	dbs := make(Databases, len(urls), len(urls))
	for i, s := range urls {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		if !u.IsAbs() {
			return nil, fmt.Errorf("invalid database url '%s'", s)
		}
		path := strings.Trim(u.Path, "/")
		if path == "" {
			return nil, fmt.Errorf("invalid database path from url '%s'", s)
		}
		u.Path = ""
		client, err := NewClient(u.String())
		if err != nil {
			return nil, err
		}
		dbs[i] = NewDatabase(client, path)
	}
	return dbs, nil
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
	return d.URL.String()
}

type Databases []*Database

func (a Databases) Len() int           { return len(a) }
func (a Databases) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Databases) Less(i, j int) bool { return a[i].Name < a[j].Name }
