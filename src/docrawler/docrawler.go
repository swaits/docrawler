package main

import (
)

// docrawl begins crawling the site at "url"
func docrawl(url string) []*Page {
	// our result structure
	var sitemap []*Page

	// what we've already parsed, what we still need to parse
	crawled := make(map[string]*Page)
	queued := make(map[string]*Page)

	// create first page and add to the queue
	homepage, _ := NewPage(nil, url) // TODO check error
	queued[homepage.URL.String()] = homepage

	// main loop
	for len(queued) > 0 {
		// pop a page off of queued
		var dest *Page
		var destURL string
		for destURL, dest = range queued {
			break
		}
		delete(queued, destURL)

		// skip loop if we already did this URL
		if _, ok := crawled[dest.URL.String()]; ok {
			continue // TODO handle
		}

		// mark this page as complete
		crawled[dest.URL.String()] = dest

		// fetch page
		text, err := fetchPage(dest)
		if err != nil {
			continue // TODO handle
		}

		// parse links
		_, links, err := parseLinks(text)
		if err != nil {
			continue // TODO handle
		}
		for _, l := range links {
			p, _ := NewPage(dest, l) // TODO check error
			dest.Children = append(dest.Children, p)

			// see if we should queue this new page for crawling
			_, haveCrawled := crawled[p.URL.String()]
			_, haveQueued := queued[p.URL.String()]
			if !haveCrawled && !haveQueued {
				queued[p.URL.String()] = p
			}
		}

		// add page to list
		sitemap = append(sitemap, dest)
	}

	return sitemap
}

func main() {
}
