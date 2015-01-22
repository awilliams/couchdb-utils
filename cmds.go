package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/awilliams/couchdb-utils/couch"
)

// listDatabase writes a list of database to output
func listDatabases(client *couch.Client, output *Output, uri bool) error {
	dbs, err := client.AllDBs()
	if err != nil {
		return err
	}
	var s string
	for _, db := range dbs {
		if uri {
			s = db.URI()
		} else {
			s = db.String()
		}
		output.Println(s)
	}
	return nil
}

// listReplicators writes a list of replicator documents to output
func listReplicators(client *couch.Client, output *Output, uri bool) error {
	replicators, err := couch.ListReplicators(client)
	if err != nil {
		return err
	}
	for _, replicator := range replicators {
		if uri {
			output.Println(replicator.URI())
		} else {
			var cont string
			if replicator.Continuous {
				cont = "continuous"
			} else {
				cont = "noncontinuous"
			}
			output.Println(replicator.String(), replicator.Source, replicator.Target, replicator.ReplicationState, cont)
		}
	}
	return nil
}

// listViews writes a list of all database views for given databases to output
func listViews(client *couch.Client, output *Output, databaseNames []string, uri bool) error {
	// create databases from input, otherwise if empty, fetch all databases
	var dbs couch.Databases
	var err error
	if len(databaseNames) == 0 {
		dbs, err = client.AllDBs()
		if err != nil {
			return err
		}
	} else {
		dbs = make(couch.Databases, len(databaseNames), len(databaseNames))
		for i, name := range databaseNames {
			dbs[i] = couch.NewDatabase(client, name)
		}
	}

	for _, db := range dbs {
		designDocs, err := couch.ListDesignDocs(db)
		if err != nil {
			return err
		}
		for _, doc := range designDocs {
			for _, view := range doc.Views {
				if uri {
					output.Println(view.URI())
				} else {
					output.Println(view.String(), doc.String(), db.String())
				}
			}
		}
	}
	return nil
}

// deleteDocs receives a list of document URLs and issues a DELETE request for each
func deleteDocs(client *couch.Client, output *Output, docURIs []string) error {
	var result struct {
		OK bool `json:"ok"`
	}
	for _, uri := range docURIs {
		doc, err := couch.NewDocumentFromURI(uri)
		if err != nil {
			return err
		}
		var curDoc couch.Document
		if err = client.Get(doc.URL.Path, nil, &curDoc); err != nil {
			return err
		}
		q := url.Values{}
		q.Set("rev", curDoc.REV)
		result.OK = false
		if err = client.Delete(doc.URL.Path, q, &result); err != nil {
			return err
		}
		if !result.OK {
			return fmt.Errorf("unknown error, DELETE operation returned OK = %v", result.OK)
		}
	}
	return nil
}

// pull receives a list of external database URLs and starts pull replication for each
func pull(client *couch.Client, output *Output, sourceURIs []string, continuous bool, createTarget bool) error {
	sourceDBs := make(couch.Databases, len(sourceURIs), len(sourceURIs))
	for i, uri := range sourceURIs {
		db, err := couch.NewDatabaseFromURI(uri)
		if err != nil {
			return err
		}
		sourceDBs[i] = db
	}
	userCtx, err := couch.GetUserCtx(client)
	if err != nil {
		return err
	}
	replicators, err := couch.ListReplicators(client)
	if err != nil {
		return err
	}

	replicatorIDs := make(map[string]*couch.Replicator, len(replicators))
	for _, replicator := range replicators {
		replicatorIDs[replicator.ID] = replicator
	}

	for _, source := range sourceDBs {
		target := couch.NewDatabase(client, source.Name)
		// create a hopefully unique replicator doc ID
		var id string
		if source.URL.IsAbs() {
			id = fmt.Sprintf("%s:%s->%s", source.URL.Host, source.Name, target.Name)
		} else {
			id = fmt.Sprintf("%s->%s", source.Name, target.Name)
		}
		if existing, found := replicatorIDs[id]; found {
			if existing.ReplicationStateReason != "" {
				return fmt.Errorf("replicator %s exists with state '%s' and reason: '%s'", id, existing.ReplicationState, existing.ReplicationStateReason)
			}
			// skip existing replicators
			continue
		}
		if _, err = client.StartReplication(userCtx, id, source, target, continuous, createTarget); err != nil {
			return err
		}
	}
	return nil
}

// refreshViews receives a list of views (either absolute or in component form) and does a HEAD request for each with state=update_after
func refreshViews(client *couch.Client, output *Output, input []string, uri bool) error {
	views := make(couch.Views, len(input), len(input))
	if uri {
		for i, uri := range input {
			view, err := couch.NewViewFromURI(uri)
			if err != nil {
				return err
			}
			views[i] = view
		}
	} else {
		for i, s := range input {
			fields := strings.Fields(s)
			if len(fields) != 3 {
				return fmt.Errorf("invalid view from components: %s", s)
			}
			db := couch.NewDatabase(client, fields[2])
			design := couch.NewDesignDoc(db, fields[1])
			views[i] = couch.NewView(design, fields[0])
		}
	}
	q := url.Values{}
	q.Set("stale", "update_after")
	var err error
	for _, view := range views {
		if err = client.Head(view.URL.Path, q, nil); err != nil {
			return err
		}
	}
	return nil
}
