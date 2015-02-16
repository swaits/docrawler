package main

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

		// mark this page as complete
		crawled[dest.URL.String()] = dest

		// fetch page
		text, err := fetchPage(dest)
		if err != nil {
			continue // TODO handle
		}

		// parse links
		title, links, err := parseLinks(text)
		if err != nil {
			continue // TODO handle
		}
		dest.Title = title
		for _, l := range links {
			p, err := NewPage(dest, l)
			if err != nil {
				continue // TODO check error
			}

			// see if we already know about this page
			pCrawled, haveCrawled := crawled[p.URL.String()]
			pQueued, haveQueued := queued[p.URL.String()]
			if haveCrawled {
				// add the previously crawled page to children
				dest.Children = append(dest.Children, pCrawled)
			} else if haveQueued {
				// add the previously queued page to children
				dest.Children = append(dest.Children, pQueued)
			} else {
				// add this new page (never seen) to children
				dest.Children = append(dest.Children, p)
				// .. and queue it up for crawling
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
