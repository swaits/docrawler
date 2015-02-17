package main

import (
	"regexp"
	"strings"
)

// a regexp which captures src/href/xhref="..." from text
var reURL = regexp.MustCompile(`(?i)(src|href|xhref)\s*=\s*"([^"]+)"`) // TODO make RFC compliant
var reTitle = regexp.MustCompile(`(?i)<\s*title\s*>([^<]*)<\s*\/\s*title`)

// parseLinks takes a string and attempts to parse any html title and all links out of it,
// and returns a slice of the captures found
func parseLinks(s string) (string, []string) {
	// find the title, or default to empty
	title := ""
	titleMatches := reTitle.FindStringSubmatch(s)
	if titleMatches != nil {
		// we found a title, replace our default title with the result
		title = strings.TrimSpace(titleMatches[1])
	}

	// find all of the URL matches (submatch = include capture groups)
	results := []string{}
	for _, match := range reURL.FindAllStringSubmatch(s, -1) {
		results = append(results, match[2])
	}
	return title, results
}
