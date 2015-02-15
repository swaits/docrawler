package main

import (
	"testing"
)

// TestURLCleaner tests cleaning up (regularization) of URLs
func TestURLCleaner(t *testing.T) {
	host, url := cleanURL("http://swaits:pass@someHOST.com:8765/blah/blah.html?x=y#foo")
	if host != "somehost.com" {
		t.Error("host extraction failed")
	}
	if url != "http://swaits:pass@somehost.com:8765/blah/blah.html?x=y#foo" {
		t.Error("url cleanup failed")
	}
}
