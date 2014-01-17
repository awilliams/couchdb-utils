package api

import (
	"fmt"
	"github.com/awilliams/couchdb-utils/util"
	"strings"
	"sync"
)

type viewsJson struct {
	Rows []struct {
		Key string
		Doc struct {
			Views map[string]struct{}
		}
	}
}

type View struct {
	Database  Database
	DesignDoc DesignDoc
	Name      string
}

func (v View) String() string {
	return v.Name
}

func (v View) refreshPath() string {
	return fmt.Sprintf(`%s/%s/_view/%s?limit=0&stale=update_after`, v.Database.String(), v.DesignDoc.ID, v.Name)
}

func (v View) PP(printer util.Printer) {
	printer.Print("%s\t%s\t%s", v.Database.String(), v.DesignDoc.String(), v.String())
}

type DesignDoc struct {
	Database Database
	ID       string
}

func (d DesignDoc) String() string {
	s := strings.Split(fmt.Sprintf("%s", d.ID), "/")
	return s[len(s)-1]
}

func (d *DesignDoc) UnmarshalJSON(data []byte) error {
	if len(data) > 2 {
		d.ID = string(data[1 : len(data)-1])
	}
	return nil
}

func (d DesignDoc) PP(printer util.Printer) {
	printer.Print(d.String())
}

type Views map[DesignDoc][]View

func (v Views) path(db Database) string {
	return fmt.Sprintf(`%s/_all_docs?startkey="_design/"&endkey="_design0"&include_docs=true`, db.String())
}

func (v Views) PP(printer util.Printer) {
	for _, views := range v {
		for _, view := range views {
			view.PP(printer)
		}
	}
}

func (c Couchdb) GetViews(db Database) (Views, error) {
	j := new(viewsJson)
	views := make(Views)
	err := c.getJson(j, views.path(db))
	if err != nil {
		return views, err
	}
	for _, row := range j.Rows {
		for viewName := range row.Doc.Views {
			var designDoc DesignDoc = DesignDoc{ID: row.Key, Database: db}
			views[designDoc] = append(views[designDoc], View{Name: viewName, Database: db, DesignDoc: designDoc})
		}
	}
	return views, nil
}

func (c Couchdb) RefreshView(view View) error {
	reader, err := c.get(view.refreshPath())
	if err != nil {
		return err
	}
	defer reader.Close()
	return nil
}

func (c Couchdb) RefreshViews(views Views) (Views, []error) {
	refreshedViews := make(Views)
	var errors []error
	var w sync.WaitGroup
	for designDoc, viewArr := range views {
		w.Add(1)
		var view View = viewArr[0]
		refreshedViews[designDoc] = append(refreshedViews[designDoc], view)
		go func(v View) {
			err := c.RefreshView(v)
			if err != nil {
				errors = append(errors, err)
			}
			w.Done()
		}(view)
	}
	w.Wait()
	return refreshedViews, errors
}
