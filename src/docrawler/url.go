package main

import (
	"errors"
	"net/url"
	"strings"
)

// custom errors
var (
	errInvalidURL = errors.New("this URL can't be parsed successfully")
)

// cleanURL takes a URL and normalizes it by downcasing the host part
func cleanURL(u *url.URL) {
	// convert host to all lower case, return updated url.URL
	u.Host = strings.ToLower(u.Host)
}

// checkURL looks to see if our URL meets our minimum requirements
func checkURL(u *url.URL) error {
	// check for non-empty Host and Scheme
	if len(u.Host) == 0 || len(u.Scheme) == 0 {
		// something about this URL doesn't meet our spec
		return errInvalidURL
	}
	return nil // success!
}

// resolveURL takes a URL and its referring URL and tries to parse it into an acceptable complete URL
// returns the final URL and any error
func resolveURL(referringURL string, currentURL string) (*url.URL, error) {
	// begin by parsing both URLs
	uCurrent, err := url.Parse(currentURL)
	if err != nil {
		return nil, err
	}
	uReferrer, err := url.Parse(referringURL) // many times this will be an extraneous parse
	if err != nil {
		return nil, err
	}

	// normalize the URLs
	cleanURL(uCurrent)
	cleanURL(uReferrer)

	// try to resolve it with the referrer
	uResolved := uReferrer.ResolveReference(uCurrent)

	// return the URL we ended up with, and the error from checkURL
	return uResolved, checkURL(uResolved)
}

// stripURL returns a version of the URL without the "Fragment" part,
// which is anything after the '#' character, so:
// http://a.com/blah.html#anchor becomes http://a.com/blah.html
// ... and without files like "index.htm" in the Path (since that's
// equivalent to "/"
func stripURL(u *url.URL) string {
	// copy the url locally
	ucopy := *u

	// blank out the anchor, return URI string
	ucopy.Fragment = ""

	// remove index.html if it exists
	suffixes := []string{"index.html", "index.htm"}
	for _, s := range suffixes {
		if strings.HasSuffix(strings.ToLower(ucopy.Path), s) {
			ucopy.Path = ucopy.Path[0 : len(ucopy.Path)-len(s)]
		}
	}

	// return new stripped down URL string
	return ucopy.String()
}
