package domain

import "testing"

func TestWeCanParseTheDomainFromAWebsiteURL(t *testing.T) {
	cases := map[string]string{
		"http://whatever.ofco.test/": "ofco.test",
		"https://www.google.com":     "google.com",
	}

	for url, expected := range cases {
		if GetDomainFromURL(url) != expected {
			t.Error("didn't work.\nexpected: " + expected + "\ngot:      " + GetDomainFromURL(url))
		}
	}
}
