package couch

func NewEntityFromPath(path string) (Entity, error) {
	uri, err := NewURI(path)
	if err != nil {
		return Entity{}, err
	}
	return Entity{
		uri: uri,
	}, nil
}

type Entity struct {
	ID  string `json:"_id,omitempty"`
	REV string `json:"_rev,omitempty"`

	uri *URI
}

func (e Entity) Path() string {
	return e.uri.Path
}
