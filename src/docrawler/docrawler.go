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

// docrawl begins crawling the site at "url"
func docrawl(url string) []*Page {
	// our result structure
	var sitemap []*Page

	// what we've already parsed, what we still need to parse
	parsed := make(map[string]*Page)
	var queued []*Page

	// create first page and add to the queue
	u, _ := resolveURL("", url) // TODO check error
	queued = append(queued, &Page{URL: u})

	// main loop
	for len(queued) > 0 {
		// pop a page off of queued
		var dest *Page
		dest, queued = queued[len(queued)-1], queued[:len(queued)-1]

		// skip loop if we already did this URL
		if _, ok := parsed[dest.URL.String()]; ok {
			continue // TODO handle
		}

		// mark this page as complete
		parsed[dest.URL.String()] = dest

		// fetch page
		text, err := fetch(dest.URL.String())
		if err != nil {
			continue // TODO handle
		}

		// parse it
		_, links, err := parse(text)
		if err != nil {
			continue // TODO handle
		}
		for _, l := range links {
			u, _ := resolveURL(dest.URL.String(), l) // TODO check error
			dest.Children = append(dest.Children, &Page{URL: u, Base: dest})
		}

		// add page to list
		sitemap = append(sitemap, dest)
	}

	return sitemap
}

func main() {
}
