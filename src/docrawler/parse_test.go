package main

import (
	"testing"
)

// TestSimpleParse verifies that we can extract URLs (etc.) from a known document
func TestSimpleParse(t *testing.T) {
	doc := `<html>
	<head>
	</head>
	<body>
		<img src="/assets/image.png"/>
		<a href="/about.html">
	</body>
</html>`
	matches, err := parse(doc)
	if err != nil {
		t.Error(err)
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
