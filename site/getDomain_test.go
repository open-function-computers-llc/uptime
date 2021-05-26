package site

import "testing"

func TestWeCanParseTheDomainFromAWebsiteURL(t *testing.T) {
	cases := map[string]string{
		"http://whatever.ofco.test/": "whatever.ofco.test",
	}

	for url, expected := range cases {
		s := Website{
			URL: url,
		}
		if s.GetDomain() != expected {
			t.Error("didn't work")
		}
	}
}
