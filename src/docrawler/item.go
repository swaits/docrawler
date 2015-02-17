package main

import (
	"net/url"
)

// itemSlice is a convenience type for a slice of items
type itemSlice []*httpItem

// itemMap convenience type for map of string's (URL) to items
type itemMap map[string]*httpItem

// itemType is an enum so we know how to classify each item
type itemType int

// enums for itemType
const (
	tUnknown itemType = iota
	tHTMLPage
	tAsset
	tRemote
	tBroken
)

// httpItem is a struct which defines a single page, which URLs (links and assets) it contains, etc.
type httpItem struct {
	url      *url.URL
	refurl   *url.URL
	title    string
	linkType itemType
	children itemSlice
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

	// figure out referrer url
	var rurl *url.URL
	if referrer != nil {
		rurl = referrer.url
	}

	// create struct and return
	return &httpItem{url: u, refurl: rurl}, nil
}
