package couch

import (
	"net/url"
	"sort"
	"strings"
)

func ListDesignDocs(db *Database) (DesignDocs, error) {
	// http://stackoverflow.com/questions/2814352/get-all-design-documents-in-couchdb
	path := db.Name + "/_all_docs"
	q := url.Values{}
	q.Set("startkey", `"_design/"`)
	q.Set("endkey", `"_design0"`)
	q.Set("include_docs", "true")

	var result struct {
		Rows []struct {
			ID  string `json:"id"`
			Key string `json:"key"`
			Doc struct {
				Views map[string]struct{} `json:"views"`
			} `json:"doc"`
		} `json:"rows"`
	}
	if err := db.Client.Get(path, q, &result); err != nil {
		return nil, err
	}

	docs := make(DesignDocs, len(result.Rows), len(result.Rows))
	for i, doc := range result.Rows {
		docs[i] = NewDesignDoc(db, doc.ID, doc.Key)
		docs[i].Views = make(Views, len(doc.Doc.Views), len(doc.Doc.Views))
		j := 0
		for view := range doc.Doc.Views {
			docs[i].Views[j] = NewView(docs[i], view)
			j++
		}
		sort.Sort(docs[i].Views)
	}
	sort.Sort(docs)
	return docs, nil
}

func NewDesignDoc(db *Database, id, key string) *DesignDoc {
	u := *db.URL
	u.Path += "/" + id
	components := strings.Split(id, "/")

	return &DesignDoc{
		ID:       id,
		Key:      key,
		Name:     components[1],
		Database: db,
		URL:      &u,
	}
}

type DesignDoc struct {
	ID    string
	Key   string
	Name  string
	Views Views

	Database *Database
	URL      *url.URL
}

func (d *DesignDoc) String() string {
	return d.Name
}

func (d *DesignDoc) URI() string {
	return d.URL.String()
}

func NewView(dd *DesignDoc, name string) *View {
	u := *dd.URL
	u.Path += "/_view/" + name
	return &View{
		Name:      name,
		DesignDoc: dd,
		URL:       &u,
	}
}

type View struct {
	Name string

	DesignDoc *DesignDoc
	URL       *url.URL
}

func (v *View) String() string {
	return v.Name
}

func (v *View) URI() string {
	return v.URL.String()
}

// sortable slices

type DesignDocs []*DesignDoc

func (a DesignDocs) Len() int           { return len(a) }
func (a DesignDocs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DesignDocs) Less(i, j int) bool { return a[i].ID < a[j].ID }

type Views []*View

func (a Views) Len() int           { return len(a) }
func (a Views) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Views) Less(i, j int) bool { return a[i].Name < a[j].Name }
