package couch

import (
	"net/url"
	"sort"
	"strings"
)

func ListReplicators(c *Client) (Replicators, error) {
	var result struct {
		Rows []struct {
			Doc Replicator `json:"doc"`
		} `json:"rows"`
	}
	query := url.Values{}
	query.Set("include_docs", "true")
	if err := c.Get("_replicator/_all_docs", query, &result); err != nil {
		return nil, err
	}
	var replicators Replicators
	for _, row := range result.Rows {
		if strings.HasPrefix(row.Doc.ID, "_") {
			continue
		}
		replicator := row.Doc
		u := *c.Host
		u.Path += "_replicator/" + row.Doc.ID
		replicator.URL = &u
		replicators = append(replicators, &replicator)
	}
	sort.Sort(replicators)
	return replicators, nil
}

type Replicator struct {
	ID                     string `json:"_id,omitempty"`
	REV                    string `json:"_rev,omitempty"`
	Source                 string `json:"source"`
	Target                 string `json:"target"`
	Cancel                 bool   `json:"cancel"`
	CreateTarget           bool   `json:"create_target"`
	Continuous             bool   `json:"continuous"`
	Owner                  string `json:"owner,omitempty"`
	ReplicationID          string `json:"_replication_id,omitempty"`
	ReplicationState       string `json:"_replication_state,omitempty"`
	ReplicationStateReason string `json:"_replication_state_reason,omitempty"`
	ReplicationStateTime   string `json:"_replication_state_time,omitempty"`

	URL *url.URL
}

func (r *Replicator) String() string {
	return r.ID
}

func (r *Replicator) URI() string {
	return r.URL.String()
}

type Replicators []*Replicator

func (a Replicators) Len() int           { return len(a) }
func (a Replicators) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Replicators) Less(i, j int) bool { return a[i].ID < a[j].ID }
