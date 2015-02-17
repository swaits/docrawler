package main

import (
	"net/url"
)

// httpItem is a struct which defines a single page, which URLs (links and assets) it contains, etc.
type httpItem struct {
	url       *url.URL
	title     string
	mediaType string
	skipped   bool
	broken    bool
	children  []*httpItem
}

// newHTTPItem takes a referring httpItem + a URL and returns a new &httpItem{}
func newHTTPItem(referrer *httpItem, rawurl string) (*httpItem, error) {
	// determine the base URL so we can resulve this into a full URL
	baseURL := ""
	if referrer != nil {
		baseURL = referrer.url.String()
	}

	// resolve URL into a full absolute URL (non-relative)
	u, err := resolveURL(baseURL, rawurl)
	if err != nil {
		return nil, err
	}

	// create struct and return
	return &httpItem{url: u, skipped: false, broken: false}, nil
}
