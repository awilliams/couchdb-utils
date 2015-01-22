package couch

import "sort"

// AllDBs fetches the list of databases on the server
func (c *Client) AllDBs() (Databases, error) {
	var result []string
	var err error
	if err = c.Get("_all_dbs", nil, &result); err != nil {
		return nil, err
	}
	dbs := make(Databases, len(result), len(result))
	for i, db := range result {
		dbs[i], err = NewDatabase(c, db)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(dbs)
	return dbs, nil
}

// NewDatabase receives a client and database name and returns a Database
func NewDatabase(c *Client, name string) (*Database, error) {
	entity, err := NewEntityFromPath(c, name)
	if err != nil {
		return nil, err
	}
	return &Database{
		Name:   name,
		Entity: entity,
	}, nil
}

type Database struct {
	Name string
	Entity
}

func (d *Database) String() string {
	return d.Name
}

type Databases []*Database

func (a Databases) Len() int           { return len(a) }
func (a Databases) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Databases) Less(i, j int) bool { return a[i].Name < a[j].Name }
