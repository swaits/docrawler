package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// docrawl begins crawling the site at "homeurl"
func docrawl(homeurl string) itemSlice {
	// set of what we have already crawled
	crawled := make(map[string]struct{})

	// map (of urlString -> httpItem) of our results
	results := make(itemMap)

	// a channel to receive our results on
	rchan := make(chan *httpItem)

	// start a ticker which we'll use to output status and check for completion
	ticker := time.Tick(250 * time.Millisecond)

	// create first page's httpItem
	homeitem, err := newHTTPItem(nil, homeurl)
	if err != nil {
		log.Fatal("unable to create httpItem for homeurl")
	}

	// set our number of outstanding pages (initially 1 to account for first page)
	crawlingCount := 1

	// start the home page crawl
	crawled[homeitem.url.String()] = struct{}{}
	go crawlItem(homeitem, rchan)

	// wait for results
	for {
		select {
		case r := <-rchan: // new results?
			// add result to our results map
			results[r.url.String()] = r

			// decrease the outstanding page count by 1
			crawlingCount--

			// start crawly any new child pages we haven't yet crawled
			for i, c := range r.children {
				// see if we already have a result for this page, if so, point to that result
				if existing, ok := results[c.url.String()]; ok {
					r.children[i] = existing
					continue
				}

				// see if we're already crawling this page (but maybe don't have results yet)
				if _, ok := crawled[c.url.String()]; !ok {
					// haven't crawled this one yet, do so now
					crawlingCount++
					crawled[c.url.String()] = struct{}{}
					go crawlItem(c, rchan)
				}
			}

		case <-ticker: // our regular ticker. for status output and checking for completion.
			// output status to console
			log.Printf("Crawled %v links, have %v left.\n", len(results), crawlingCount)

			// see if we're finished
			if crawlingCount == 0 {
				// finished! convert results map to a slice and return it
				rslice := itemSlice{}
				for _, v := range results {
					rslice = append(rslice, v)
				}
				return rslice
			}
		}
	}
}

// crawlItem crawls a single httpItem, fetching the header, hte page, parsing it,
// and filling out its structure as much as possible
func crawlItem(item *httpItem, rchan chan<- *httpItem) {
	// make sure this item is the same domain (i.e. URL "host part") as its referrer
	if item.refurl != nil && item.url.Host != item.refurl.Host {
		// skip URLs associated with other Hosts
		item.linkType = tRemote
		rchan <- item
		return
	}

	// fetch page
	text, err := fetchPage(item)
	if err != nil {
		rchan <- item
		return
	}

	// parse links
	title, links, _ := parseLinks(text) // TODO remove error from parseLinks, not needed
	if err != nil {
		item.linkType = tBroken
		rchan <- item
		return
	}
	item.title = title

	// walk links and add them as children to the current item
	for _, l := range links {
		newItem, err := newHTTPItem(item, l)
		if err != nil {
			continue // TODO bad item
		}
		item.children = append(item.children, newItem)
	}

	// send back our item struct now that it's all filled out
	rchan <- item
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
