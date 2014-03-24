package main

import (
	"github.com/awilliams/couchdb-utils/util"
)

type couchdbutils struct {
	Name        string
	Version     string
	Description string
}

func (c couchdbutils) PP(printer util.Printer) {
	printer.Print("%s v%s\n%s", c.Name, c.Version, c.Description)
}

var CouchdbUtils *couchdbutils = &couchdbutils{
	"couchdb-utils",
	"0.0.2",
	"A fast, lightweight, and portable CouchDB utility.",
}

func main() {
	executeCli()
}
