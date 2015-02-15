package main

import (
	"testing"
)

// TestSimpleParse verifies that we can extract the page Title and all URLs (etc.) from a known document
func TestSimpleParse(t *testing.T) {
	doc := `<html>
	<head>
		<title> Test Page  </title>
	</head>
	<body>
		<img src="/assets/image.png"/>
		<a href="/about.html">
	</body>
</html>`
	title, matches, err := parse(doc)
	if err != nil {
		t.Error(err)
	}
	if title != "Test Page" {
		t.Error("got wrong title")
	}
	if len(matches) != 2 {
		t.Error("invalid number of matches in parse")
	}
	if matches[0] != "/about.html" {
		t.Error("match text is invalid")
	}
	if matches[1] != "/assets/image.png" {
		t.Error("match text is invalid")
	}
}
