package couch

import (
	"fmt"
	"net/url"
	"strings"
)

// DocumentsFromURI processes a document URI (either absolute or relative) and creates a Document
func NewDocumentFromURI(uri string) (*Document, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	path := strings.Trim(u.Path, "/")
	if path == "" {
		return nil, fmt.Errorf("invalid document path from url '%s'", uri)
	}
	pathComponents := strings.Split(path, "/")
	if len(pathComponents) < 2 {
		return nil, fmt.Errorf("invalid document path from url '%s'", uri)
	}
	u.Path = pathComponents[0]
	db, err := NewDatabaseFromURI(u.String())
	if err != nil {
		return nil, err
	}
	return NewDocument(db, strings.Join(pathComponents[1:], "/")), nil
}

func NewDocument(db *Database, id string) *Document {
	u := *db.URL
	u.Path += "/" + id
	return &Document{
		ID:       id,
		Database: db,
		URL:      &u,
	}
}

// Document is a generic CouchDB document.
type Document struct {
	ID  string `json:"_id,omitempty"`
	REV string `json:"_rev,omitempty"`

	Database *Database `json:"-"`
	URL      *url.URL  `json:"-"`
}

func (r *Document) String() string {
	return r.ID
}

func (r *Document) URI() string {
	return r.URL.String()
}

type Documents []*Document

func (a Documents) Len() int           { return len(a) }
func (a Documents) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Documents) Less(i, j int) bool { return a[i].ID < a[j].ID }
