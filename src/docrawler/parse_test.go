package main

import (
	"testing"
)

// TestSimpleParse verifies that we can extract the page Title and all URLs (etc.) from a known document
func TestSimpleParse(t *testing.T) {
	doc := `<html>
	<head>
		<TITle> Test Page  </title>
	</head>
	<body>
		<img src="/assets/image.png"/>
		<a href="/about.html">fooey!</a>
		<script href="scripts/blah.js"/>
	</body>
</html>`
	title, matches := parseLinks(doc)
	if title != "Test Page" {
		t.Error("got wrong title")
	}
	if len(matches) != 3 {
		t.Error("invalid number of matches in parse")
	}
	if matches[0] != "/assets/image.png" {
		t.Error("match text is invalid")
	}
	if matches[1] != "/about.html" {
		t.Error("match text is invalid")
	}
	if matches[2] != "scripts/blah.js" {
		t.Error("match text is invalid")
	}
}
