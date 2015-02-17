package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// docrawl begins crawling the site at "homeurl"
func doCrawl(homeurl string) itemSlice {
	// set of what we have already crawled, our results
	crawled := make(itemMap)

	// a set of crawled URLs which have been cleaned
	// useful so that we don't crawl http://a.com/index.html#about if we've
	// already crawled http://a.com/index.html, or vice versa
	crawledStripped := make(itemMap)

	// a channel to receive our results on
	rchan := make(chan *httpItem)

	// start a ticker which we'll use to output status and check for completion
	ticker := time.Tick(250 * time.Millisecond)

	// create first page's httpItem
	homeitem, err := newHTTPItem(nil, homeurl)
	if err != nil {
		return nil
	}

	// set our number of outstanding pages (initially 1 to account for first page)
	crawlingCount := 1

	// start the home page crawl
	crawled[homeitem.url.String()] = homeitem
	crawledStripped[stripURL(homeitem.url)] = homeitem
	go crawlItem(homeitem, rchan)

	// wait for results
	for {
		select {
		case r := <-rchan: // new results?
			// add result to our results map
			crawled[r.url.String()] = r

			// decrease the outstanding page count by 1
			crawlingCount--

			// start crawly any new child pages we haven't yet crawled
			for i, c := range r.children {
				// see if we already have a result for this page
				// if so, point to that result
				if existing, ok := crawled[c.url.String()]; ok {
					r.children[i] = existing
					continue
				}

				// see if we're already crawling this page (but maybe don't have results yet)
				// if so, point to that item (we will have it as a result later!)
				if existing, ok := crawled[c.url.String()]; ok {
					r.children[i] = existing
					continue
				}
				if existing, ok := crawledStripped[stripURL(c.url)]; ok {
					// we crawled a different version of this same page
					// i.e. same page, different anchor. we don't need to crawl it again
					// but we do need to keep it in our results, so copy the existing
					// struct over for everything except the URLs
					r.children[i].title = existing.title
					r.children[i].linkType = existing.linkType
					r.children[i].children = existing.children
					continue
				}

				// haven't crawled this one yet, do so now
				crawlingCount++
				crawled[c.url.String()] = c
				crawledStripped[stripURL(c.url)] = c
				go crawlItem(c, rchan)
			}

		case <-ticker: // our regular ticker. for status output and checking for completion.
			// output status to console
			log.Printf("Crawled %v links, have %v left.\n", len(crawled), crawlingCount)

			// see if we're finished
			if crawlingCount == 0 {
				// finished! convert results map to a slice and return it
				rslice := itemSlice{}
				for _, v := range crawled {
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
	title, links := parseLinks(text) // TODO remove error from parseLinks, not needed
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
	// see if we've got no arguments
	if len(os.Args) == 1 {
		fmt.Printf("\nD.O. Crawler 1.0  Copyright (c) 2015 Stephen Waits <steve@waits.net>  2015-02-17\n\n")
		fmt.Printf("usage: %v <URLs...>\n\n", os.Args[0])
		os.Exit(1)
	}
	// crawl each URL on the command line
	for _, u := range os.Args[1:] {
		pages := doCrawl(u)
		if pages == nil {
			log.Fatalf("unable to crawl %q, may be an invalid URL\n", u)
		}
		l := sitemapToLocations(pages)
		j, _ := locationsToJSON(l)
		fmt.Println(j)
	}
}
