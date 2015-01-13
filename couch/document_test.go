package couch

import "testing"

func TestNewDocumentFromURI(t *testing.T) {
	cases := []struct {
		Given          string
		ExpectedError  bool
		ExpectedDBName string
		ExpectedID     string
		ExpectedURI    string
	}{
		{"http://1.2.3.4:5984/db/doc1", false, "db", "doc1", "http://1.2.3.4:5984/db/doc1"},
		{"http://some.domain.com:5984/db/asdf-as", false, "db", "asdf-as", "http://some.domain.com:5984/db/asdf-as"},
		{"https://some.domain.com:5984/abc/test", false, "abc", "test", "https://some.domain.com:5984/abc/test"},
		{"abc/test", false, "abc", "test", "abc/test"},

		{"https://some.domain.com:5984/onlydb", true, "", "", ""},
		{"some.domain.com:5984/onlydb", true, "", "", ""},
	}

	for _, c := range cases {
		example, err := NewDocumentFromURI(c.Given)
		if err != nil && !c.ExpectedError {
			t.Fatal(err)
		}
		if c.ExpectedError {
			if err == nil {
				t.Fatal("Expected NewDocumentFromURI to return error, none given")
			}
			continue
		}

		if c.ExpectedDBName != example.Database.Name {
			t.Errorf("Unexpected Document.Database.Name. Expected %s, given %s", c.ExpectedDBName, example.Database.Name)
		}
		if c.ExpectedID != example.ID {
			t.Errorf("Unexpected Document.ID. Expected %s, given %s", c.ExpectedID, example.ID)
		}
		if c.ExpectedURI != example.URI() {
			t.Errorf("Unexpected Document#URI. Expected %s, given %s", c.ExpectedURI, example.URI())
		}
	}
}
