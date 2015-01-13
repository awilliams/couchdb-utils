package couch

import "testing"

func TestNewDatabaseFromURIs(t *testing.T) {
	cases := []struct {
		Given         string
		ExpectedError bool
		ExpectedName  string
		ExpectedURI   string
	}{
		{"http://1.2.3.4:5984/db", false, "db", "http://1.2.3.4:5984/db"},
		{"http://some.domain.com:5984/db", false, "db", "http://some.domain.com:5984/db"},
		{"https://some.domain.com:5984/db/test", false, "db/test", "https://some.domain.com:5984/db/test"},
		{"some.domain.com/db/test", false, "some.domain.com/db/test", "some.domain.com/db/test"},
		{"/db", false, "db", "db"},
		{"/db/test", false, "db/test", "db/test"},

		{"http://1.2.3.4:5984/", true, "", ""},
		{"http://somehost.com/", true, "", ""},
	}

	for _, c := range cases {
		example, err := NewDatabaseFromURI(c.Given)
		if err != nil && !c.ExpectedError {
			t.Fatal(err)
		}
		if c.ExpectedError {
			if err == nil {
				t.Fatal("Expected NewDatabaseFromURI to return error, none given")
			}
			continue
		}

		if c.ExpectedName != example.Name {
			t.Fatalf("Unexpected Database.Name. Expected %s, given %s", c.ExpectedName, example.Name)
		}
		if c.ExpectedURI != example.URI() {
			t.Fatalf("Unexpected Database#URI. Expected %s, given %s", c.ExpectedURI, example.URI())
		}
	}
}
