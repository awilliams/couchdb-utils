package couch

import "testing"

func TestNewURI(t *testing.T) {
	cases := []struct {
		Given    string
		Expected string
	}{
		{"http://www.test.org", ""},
		{"http://www.test.org/", ""},
		{"http://www.test.org/a", "/a"},
		{"http://www.test.org/a/", "/a"},

		{"a", "/a"},
		{"a/", "/a"},
		{"/a/", "/a"},
		{"/a/b/c/", "/a/b/c"},
	}

	for _, example := range cases {
		u, _ := NewURI(example.Given)
		given := u.Path
		if example.Expected != given {
			t.Errorf("Unexpected URI.Path for input %q. Expected %s, given %s", example.Given, example.Expected, given)
		}
	}
}

func TestURIExtend(t *testing.T) {
	cases := []struct {
		GivenBase       string
		GivenComponents []string
		Expected        string
	}{
		{"http://www.test.org", []string{"a", "b"}, "http://www.test.org/a/b"},
		{"http://www.test.org/", []string{"/a", "/b"}, "http://www.test.org/a/b"},
		{"http://www.test.org/", []string{"/a/", "/b/"}, "http://www.test.org/a/b"},
		{"http://www.test.org/a", []string{"/b", "/c"}, "http://www.test.org/a/b/c"},
		{"http://www.test.org/a/", []string{"/b", "/c"}, "http://www.test.org/a/b/c"},
		{"http://www.test.org/a/", []string{"/b/", "/c/"}, "http://www.test.org/a/b/c"},

		{"http://www.test.org/a/", []string{"/b/c/"}, "http://www.test.org/a/b/c"},

		{"a", []string{"/b/", "/c/"}, "/a/b/c"},
		{"a/", []string{"/b/", "/c/"}, "/a/b/c"},
		{"/a/", []string{"/b/", "/c/"}, "/a/b/c"},
	}

	for _, example := range cases {
		u, _ := NewURI(example.GivenBase)
		given := u.Extend(example.GivenComponents...).String()
		if example.Expected != given {
			t.Errorf("Unexpected URI#Extend for input '%q:%q'. Expected %s, given %s", example.GivenBase, example.GivenComponents, example.Expected, given)
		}
	}
}

func TestURIReduce(t *testing.T) {
	cases := []struct {
		GivenBase   string
		GivenSuffix string
		Expected    string
	}{
		{"http://www.test.org/", "", "http://www.test.org"},
		{"http://www.test.org", "/", "http://www.test.org"},
		{"http://www.test.org/a", "a", "http://www.test.org"},
		{"http://www.test.org/a/", "a", "http://www.test.org"},
		{"http://www.test.org/a", "/a/", "http://www.test.org"},
		{"http://www.test.org/a/", "/a/", "http://www.test.org"},
		{"http://www.test.org/a", "/a", "http://www.test.org"},
		{"http://www.test.org/a", "a/", "http://www.test.org"},

		{"http://www.test.org/a/b", "a/b", "http://www.test.org"},
		{"http://www.test.org/a/b", "b", "http://www.test.org/a"},

		{"http://www.test.org/a/b", "a", "http://www.test.org/a/b"},
		{"http://www.test.org/a/b", "z/a/b", "http://www.test.org/a/b"},

		{"a", "a", ""},
		{"/a/b/c/", "c", "/a/b"},
	}

	for _, example := range cases {
		u, _ := NewURI(example.GivenBase)
		given := u.Reduce(example.GivenSuffix).String()
		if example.Expected != given {
			t.Errorf("Unexpected URI#Reduce for input '%q:%q'. Expected %s, given %s", example.GivenBase, example.GivenSuffix, example.Expected, given)
		}
	}
}

func TestURIComponents(t *testing.T) {
	cases := []struct {
		Given    string
		Expected []string
	}{
		{"http://www.test.org", []string{}},
		{"http://www.test.org/", []string{}},
		{"http://www.test.org/a", []string{"a"}},
		{"http://www.test.org/a/b/", []string{"a", "b"}},
		{"http://www.test.org/a/b/c", []string{"a", "b", "c"}},

		{"", []string{}},
		{"a", []string{"a"}},
		{"/a", []string{"a"}},
		{"/a/", []string{"a"}},
		{"/a/b", []string{"a", "b"}},
		{"somedomain.com/a/b", []string{"somedomain.com", "a", "b"}},
	}

	for _, example := range cases {
		u, _ := NewURI(example.Given)
		given := u.Components()
		if len(example.Expected) != len(given) {
			t.Errorf("Unexpected URI#Components for input %q. Expected %q, given %q", example.Given, example.Expected, given)
		}
		for i, expected := range example.Expected {
			if expected != given[i] {
				t.Errorf("Unexpected URI#Components for input %q. Expected %q, given %q", example.Given, example.Expected, given)
			}
		}
	}
}
