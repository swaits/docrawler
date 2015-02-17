package main

import (
	"fmt"
	"os"
)

type itemSlice []*httpItem
type itemMap   map[string]*httpItem


// docrawl begins crawling the site at "url"
func docrawl(url string) itemSlice {
	// our result structure
	var sitemap itemSlice

	// what we've already parsed, what we still need to parse
	crawled := make(itemMap)
	queued := make(itemMap)

	// create first page and add to the queue
	homepage, _ := newHTTPItem(nil, url) // TODO check error
	queued[homepage.url.String()] = homepage

	// main loop
	for len(queued) > 0 {
		// pop a page off of queued
		var dest *httpItem
		var destURL string
		for destURL, dest = range queued {
			break
		}
		delete(queued, destURL)

		// mark this page as complete
		crawled[dest.url.String()] = dest

		// add page to list
		sitemap = append(sitemap, dest)

		// make sure this URL has the same hostname as our first page
		if dest.url.Host != homepage.url.Host {
			// skip URLs associated with other Hosts
			dest.skipped = true
			continue
		}

		// fetch page
		text, err := fetchPage(dest)
		if err != nil {
			dest.broken = true
			continue
		}

		// parse links
		title, links, err := parseLinks(text)
		if err != nil {
			dest.broken = true
			continue
		}
		dest.title = title
		for _, l := range links {
			p, err := newHTTPItem(dest, l)
			if err != nil {
				continue // TODO check error
			}

			// see if we already know about this page
			pCrawled, haveCrawled := crawled[p.url.String()]
			pQueued, haveQueued := queued[p.url.String()]
			if haveCrawled {
				// add the previously crawled page to children
				dest.children = append(dest.children, pCrawled)
			} else if haveQueued {
				// add the previously queued page to children
				dest.children = append(dest.children, pQueued)
			} else {
				// add this new page (never seen) to children
				dest.children = append(dest.children, p)
				// .. and queue it up for crawling
				queued[p.url.String()] = p
			}
		}
	}

	return sitemap
}

func main() {
	// crawl each URL on the command line
	for _, u := range os.Args[1:] {
		pages := docrawl(u)
		l := sitemapToLocations(pages)
		j, _ := locationsToJSON(l)
		fmt.Println(j)
	}
}
