package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// doCrawl begins crawling the site at "homeurl"
func doCrawl(homeurl string, nWorkers int) itemSlice {
	// set of what we have already crawled, our results
	crawled := make(itemMap)

	// a set of crawled URLs which have been cleaned
	// useful so that we don't crawl http://a.com/index.html#about if we've
	// already crawled http://a.com/index.html, or vice versa
	crawledStripped := make(itemMap)

	// channels to send work to and receive our results on
	rxchan := make(chan *httpItem)
	txchan := make(chan *httpItem)

	// spin up our crawler workers
	for i := 0; i < nWorkers; i++ {
		go crawlWorker(txchan, rxchan)
	}

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
	txchan <- homeitem

	// wait for results
	for {
		select {
		case r := <-rxchan: // new results?
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
				go func(newItem *httpItem) {
					txchan <- newItem
				}(c)
			}

		case <-ticker: // our regular ticker. for status output and checking for completion.
			// output status to console
			log.Printf("Crawled %v links, have %v left.\n", len(crawled), crawlingCount)

			// see if we're finished
			if crawlingCount == 0 {
				// close our channels, signalling any workers to exit
				close(txchan)
				close(rxchan)

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

// crawlWorker is a goroutine'ized wrapper around crawlItem that listens
// for new jobs and sends them off to crawlItem, returning the results in rxchan
func crawlWorker(txchan <-chan *httpItem, rxchan chan<- *httpItem) {
	// read off the incoming item channel forever
	for newJob := range txchan {
		// perform the crawl
		crawlItem(newJob)
		// return the result
		rxchan <- newJob
	}
}

// main is our program's entry point
func main() {
	// parse our flags
	nWorkers := *(flag.Uint("num", 100, "number of workers"))
	// see if we've got no arguments
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Printf("\nerror: Please specify at least one URL to crawl.\n")
		fmt.Printf("\nD.O. Crawler 1.0  Copyright (c) 2015 Stephen Waits <steve@waits.net>  2015-02-17\n\n")
		fmt.Printf("usage: %v [-num=100] <URLs...>\n\n", os.Args[0])
		fmt.Printf("  -num=100: number of workers\n")
		fmt.Printf("  URLs:     URLs to crawl\n\n")
		os.Exit(1)
	}

	// crawl each URL on the command line
	for _, u := range flag.Args() {
		pages := doCrawl(u, int(nWorkers))
		if pages == nil {
			log.Fatalf("unable to crawl %q, may be an invalid URL\n", u)
		}
		l := sitemapToLocations(pages)
		j, _ := locationsToJSON(l)
		fmt.Println(j)
	}
}
