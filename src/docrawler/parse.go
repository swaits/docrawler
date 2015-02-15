package main

import (
	"regexp"
	"sort"
	"strings"
)

// a regexp which captures src/href/xhref="..." from text
var reURL = regexp.MustCompile(`(src|href|xhref)\s*=\s*"([^"]+)"`) // TODO make RFC compliant
var reTitle = regexp.MustCompile(`<\s*title\s*>([^<]*)`)

// parse takes a string and attempts to parse any html title and all links out of it,
// and returns a sorted slice of the captures found
func parse(s string) (string, []string, error) {
	// find the title, or default to empty
	title := ""
	titleMatches := reTitle.FindStringSubmatch(s)
	if titleMatches != nil {
		title = strings.TrimSpace(titleMatches[1])
	}

	// find all of the URL matches (submatch = include capture groups)
	results := []string{}
	for _, match := range reURL.FindAllStringSubmatch(s, -1) {
		results = append(results, match[2])
	}
	sort.Strings(results)
	return title, results, nil
}
