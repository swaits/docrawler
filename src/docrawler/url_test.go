package main

import (
	"net/url"
	"testing"
)

// TestURLCleaner tests cleaning up (regularization) of URLs
func TestURLCleaner(t *testing.T) {
	u, err := url.Parse("http://swaits:pass@someHOST.com:8765/blah/blah.html?x=y#foo")
	if err != nil {
		t.Error("error parsing test URL")
	}
	cleanURL(u)
	if u.Host != "somehost.com:8765" {
		t.Logf("   Got: %q\n", u.Host)
		t.Logf("Wanted: %q\n", "somehost.com:8765")
		t.Error("host extraction failed")
	}
	if u.String() != "http://swaits:pass@somehost.com:8765/blah/blah.html?x=y#foo" {
		t.Error("url cleanup failed")
	}
}

// send in a bad URL and make sure we get nil back
func TestBadURL(t *testing.T) {
	u, err := url.Parse("blah")
	if err != nil {
		t.Error("error parsing test URL")
	}
	if err := checkURL(u); err == nil {
		t.Error("parsing bogus URL still returned some values")
	}
}

// some tests for the base url resolution (absolute vs. relative links)
func doResolveTest(t *testing.T, wanted, referrer, current string) {
	u, err := resolveURL(referrer, current)
	if err != nil {
		t.Error("resolveURL failed (probably parsing)")
	}
	if u.String() != wanted {
		t.Logf("   Got: %q\n", u)
		t.Logf("Wanted: %q\n", wanted)
		t.Error("url resolution with base URL failed")
	}
}
func TestURLResolution(t *testing.T) {
	doResolveTest(t, "http://base.com/about.html", "http://BaSe.CoM/index.html", "about.html")
	doResolveTest(t, "http://base.com/some/crazy/path/test.png", "http://BaSe.CoM/some/crazy/path/index.html?parm=x&blah=foo", "test.png")
	doResolveTest(t, "http://S:D@base.com:1/test.png", "http://S:D@BaSe.CoM:1/some/crazy/path/index.html?parm=x&blah=foo", "/test.png")
	doResolveTest(t, "http://another.com/blah.html", "http://S:D@BaSe.CoM:1/some/crazy/path/index.html?parm=x&blah=foo", "http://aNOThER.CoM/blah.html")
	doResolveTest(t, "http://empty.com/blah.html", "", "http://empty.com/blah.html")
}

// TestBadResolve tries to resolve bogus URLs to be sure we're getting errors
/* TODO not sure how to make url.Parse fail in these cases!!!
func TestBadResolve(t* testing.T) {
	u, err := resolveURL("http://a.com/", "h%20t\tp://doesnte%20xist23492387492837492374982734.com/")
	if err == nil {
		t.Logf("got %q", u.String())
		t.Error("tried resolveURL with an invalid URL but didn't get an error")
	}
	u, err = resolveURL("h%20t\tp://doesnte%20xist23492387492837492374982734.com/", "http://a.com/")
	if err == nil {
		t.Logf("got %q", u.String())
		t.Error("tried resolveURL with an invalid URL but didn't get an error")
	}
}
*/

func TestURLStripAnchor(t *testing.T) {
	u, err := url.Parse("http://swaits:pass@somehost.com:8765/blah/blah.html?x=y#foo")
	if err != nil {
		t.Fatal("couldn't parse a URL")
	}
	stripped := stripAnchorFromURL(u)
	if stripped != "http://swaits:pass@somehost.com:8765/blah/blah.html?x=y" {
		t.Error("stripping anchor #foo from URL failed")
	}
}
