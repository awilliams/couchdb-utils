package api

import (
	"github.com/awilliams/couchdb-utils/util"
	"net/url"
)

type Database struct {
	Name *string
}

func (d *Database) String() string {
	return *d.Name
}

func (d *Database) URL() string {
	return url.QueryEscape(d.String())
}

func (d Database) PP(printer util.Printer) {
	printer.Print(d.String())
}

func (d *Database) UnmarshalJSON(data []byte) error {
	if len(data) > 2 {
		name := string(data[1 : len(data)-1])
		d.Name = &name
	}
	return nil
}

type Databases []Database

func (d Databases) PP(printer util.Printer) {
	for _, db := range d {
		db.PP(printer)
	}
}

func (d *Databases) path() string {
	return "_all_dbs"
}

func (c Couchdb) GetDatabases() (Databases, error) {
	databases := new(Databases)
	err := c.getJson(databases, databases.path())
	return *databases, err
}
