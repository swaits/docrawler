package main

import (
	"net/url"
)

// Page is a struct which defines a single page, which URLs (links and assets) it contains, etc.
type Page struct {
	URL       *url.URL
	Base      *Page
	Title     string
	MediaType string
	Children  []*Page
}

// NewPage takes a referring Page + a URL and returns a new &Page{}
func NewPage(referrer *Page, rawurl string) (*Page, error) {
	// determine the base URL so we can resulve this into a full URL
	baseURL := ""
	if referrer != nil {
		baseURL = referrer.URL.String()
	}

	// resolve URL into a full absolute URL (non-relative)
	u, err := resolveURL(baseURL, rawurl)
	if err != nil {
		return nil, err
	}

	// create struct and return
	return &Page{URL: u}, nil
}
