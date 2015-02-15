package main

import (
	"regexp"
)

func parse(s string) ([]string, error) {
	// hacky RE to capture any kind of URL link in html ('a', 'img', 'script', 'video', etc.)
	re := regexp.MustCompile(`(src|href|xhref)\s*=\s*"([^"]+)"`)

	// find all of the matches (submatch = include capture groups)
	results := []string{}
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		results = append(results, match[2])
	}
	return results, nil
}
