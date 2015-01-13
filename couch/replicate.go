package couch

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const replicatorDB = "_replicator" // default replicator db. Possible to change via couchdb config.

func ListReplicators(c *Client) (Replicators, error) {
	var result struct {
		Rows []struct {
			Doc Replicator `json:"doc"`
		} `json:"rows"`
	}
	query := url.Values{}
	query.Set("include_docs", "true")
	if err := c.Get(replicatorDB+"/_all_docs", query, &result); err != nil {
		return nil, err
	}
	db := NewDatabase(c, replicatorDB)
	var replicators Replicators
	for _, row := range result.Rows {
		if strings.HasPrefix(row.Doc.ID, "_") {
			continue
		}
		replicator := row.Doc
		replicator.Document = NewDocument(db, row.Doc.ID)
		replicators = append(replicators, &replicator)
	}
	sort.Sort(replicators)
	return replicators, nil
}

func (c *Client) StartReplication(userCtx *UserCtx, id string, source *Database, target *Database, continuous bool, createTarget bool) (*Replicator, error) {
	// if source or target is a local database, use relative URL, else absolute
	var sourceURI string
	if source.Client == c {
		sourceURI = source.Name
	} else {
		sourceURI = source.URI()
	}
	var targetURI string
	if target.Client == c {
		targetURI = source.Name
	} else {
		targetURI = source.URI()
	}

	replicator := Replicator{
		Document:     NewDocument(NewDatabase(c, replicatorDB), id),
		Source:       sourceURI,
		Target:       targetURI,
		CreateTarget: createTarget,
		Continuous:   continuous,
		UserCtx:      *userCtx,
	}

	path := replicatorDB + "/" + replicator.ID
	var result struct {
		ID  string `json:"id"`
		REV string `json:"rev"`
		OK  bool   `json:"ok"`
	}
	if err := c.Put(path, nil, replicator, &result); err != nil {
		return nil, err
	}
	if !result.OK {
		return nil, fmt.Errorf("uknown error, replicator returned ok = %v", result.OK)
	}
	replicator.ID = result.ID
	replicator.REV = result.REV
	return &replicator, nil
}

type Replicator struct {
	*Document
	Source                 string  `json:"source"`
	Target                 string  `json:"target"`
	Cancel                 bool    `json:"cancel"`
	CreateTarget           bool    `json:"create_target"`
	Continuous             bool    `json:"continuous"`
	Owner                  string  `json:"owner,omitempty"`
	ReplicationID          string  `json:"_replication_id,omitempty"`
	ReplicationState       string  `json:"_replication_state,omitempty"`
	ReplicationStateReason string  `json:"_replication_state_reason,omitempty"`
	ReplicationStateTime   string  `json:"_replication_state_time,omitempty"`
	UserCtx                UserCtx `json:"user_ctx"` // see session
}

type Replicators []*Replicator

func (a Replicators) Len() int           { return len(a) }
func (a Replicators) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Replicators) Less(i, j int) bool { return a[i].ID < a[j].ID }
