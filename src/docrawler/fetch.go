package main

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
)

// custom errors
var (
	errContentTypeNotFound = errors.New("no Content-Type header found")
	errFetchError          = errors.New("couldn't fetch item")
	errFileTypeUnknown     = errors.New("couldn't determine file type")
)

// fetchFiletype performs an http HEAD to get the media type, and sets it
// directly in httpItem.mediaType
func (item *httpItem) fetchFiletype() error {
	resp, err := http.Head(item.url.String())
	if err != nil {
		item.linkType = tBroken
		return err
	}

	// check response code
	if resp.StatusCode != http.StatusOK {
		item.linkType = tBroken
		return errFetchError
	}

	// pull the content type out of the http header
	contentType, ok := resp.Header["Content-Type"]
	if !ok {
		// no Content-Type header found
		item.linkType = tBroken
		return errContentTypeNotFound
	}

	// parse mime type
	mediatype, _, err := mime.ParseMediaType(contentType[0])
	if err != nil {
		item.linkType = tBroken
		return err
	}

	// success, set item type and return
	if mediatype == "text/html" {
		item.linkType = tHTMLPage
	} else {
		item.linkType = tAsset
	}
	return nil
}

// fetchPage takes a httpItem, GETs it, and returns the body as a string
func (item *httpItem) fetchItem() (string, error) {
	// figure out the file type
	err := item.fetchFiletype()
	if err != nil {
		return "", err
	}

	// make sure we got a file type!
	if item.linkType == tUnknown {
		return "", errFileTypeUnknown
	}

	// we only want to fetch html, so return if it's anything else
	if item.linkType != tHTMLPage {
		return "", nil // but this isn't an error!
	}

	// GET the url
	resp, err := http.Get(item.url.String())
	if err != nil {
		return "", err
	}

	// check response code
	if resp.StatusCode != http.StatusOK {
		return "", errFetchError
	}

	// read the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// success! close the body, return it as a string
	resp.Body.Close()
	return string(body), nil
}

// crawlItem crawls a single httpItem, fetching the header, hte page, parsing it,
// and filling out its structure as much as possible
func (item *httpItem) crawlItem() {
	// make sure this item is the same domain (i.e. URL "host part") as its referrer
	if item.refurl != nil && item.url.Host != item.refurl.Host {
		// skip URLs associated with other Hosts
		item.linkType = tRemote
		return
	}

	// fetch page
	text, err := item.fetchItem()
	if err != nil {
		return
	}

	// parse links
	title, links := parseLinks(text)
	item.title = title

	// walk links and add them as children to the current item
	for _, l := range links {
		newItem, err := newHTTPItem(item, l)
		if err != nil {
			continue // TODO bad item
		}
		item.children = append(item.children, newItem)
	}
}
