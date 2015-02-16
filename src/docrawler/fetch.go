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
// directly in Page.MediaType
func fetchFiletype(page *Page) error {
	resp, err := http.Head(page.URL.String())
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
	page.MediaType = mediatype
	return nil
}

// fetchPage takes a Page, GETs it, and returns the body as a string
func fetchPage(page *Page) (string, error) {
	// figure out the file type
	err := fetchFiletype(page)
	if err != nil {
		return "", err
	}

	// we only want to fetch html
	if page.MediaType != "text/html" {
		return "", nil // but this isn't an error!
	}

	// GET the url
	resp, err := http.Get(page.URL.String())
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
