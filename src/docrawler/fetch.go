package main

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
)

var (
	errContentTypeNotFound = errors.New("no Content-Type header found")
	errFetchError          = errors.New("couldn't fetch page")
)

// fetchFiletype performs an http HEAD to get the media type, and sets it
// directly in httpItem.mediaType
func fetchFiletype(page *httpItem) error {
	resp, err := http.Head(page.url.String())
	if err != nil {
		return err
	}

	// check response code
	if resp.StatusCode != http.StatusOK {
		return errFetchError
	}

	// pull the content type out of the http header
	contentType, ok := resp.Header["Content-Type"]
	if !ok {
		// no Content-Type header found
		return errContentTypeNotFound
	}

	// parse mime type
	mediatype, _, err := mime.ParseMediaType(contentType[0])
	if err != nil {
		return err
	}

	// set and return
	page.mediaType = mediatype
	return nil
}

// fetchPage takes a httpItem, GETs it, and returns the body as a string
func fetchPage(page *httpItem) (string, error) {
	// figure out the file type
	err := fetchFiletype(page)
	if err != nil {
		return "", err
	}

	// we only want to fetch html
	if page.mediaType != "text/html" {
		return "", nil // but this isn't an error!
	}

	// GET the url
	resp, err := http.Get(page.url.String())
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
