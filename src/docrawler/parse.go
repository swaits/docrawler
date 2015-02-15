package main

import (
	"regexp"
	"sort"
)

// a regexp which captures src/href/xhref="..." from text
var reURL = regexp.MustCompile(`(src|href|xhref)\s*=\s*"([^"]+)"`) // TODO make RFC nice

// parse takes a string and attempts to parse any html links out of it,
// and returns a sorted slice of the captures found
func parse(s string) ([]string, error) {
	// find all of the matches (submatch = include capture groups)
	results := []string{}
	for _, match := range reURL.FindAllStringSubmatch(s, -1) {
		results = append(results, match[2])
	}
	sort.Strings(results)
	return results, nil
}
