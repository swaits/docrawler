package main

import (
	"net/url"
	"strings"
)

// cleanURL takes a URL, normalizes it by downcasing the host part, and reassembles it.
// returns the updated host and reassembled URL
func cleanURL(dirtyURL string) (string, string) {
	// parse the URL into its parts
	u, err := url.Parse(dirtyURL)
	if err != nil {
		return "", dirtyURL // TODO what should we really return here?
	}

	// convert host to all lower case
	u.Host = strings.ToLower(u.Host) 

	// return the updated host and reassembled URL
	return u.Host, u.String()
}
