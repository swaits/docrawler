package main

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
)

var (
	ErrContentTypeNotFound = errors.New("no Content-Type header found")
)

func fetch_filetype(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}

	// pull the content type out of the http header
	contentType, ok := resp.Header["Content-Type"]
	if !ok {
		// no Content-Type header found
		return "", ErrContentTypeNotFound
	}

	// parse mime type
	mediatype, _, err := mime.ParseMediaType(contentType[0])
	if err != nil {
		return "", err
	}

	return mediatype, nil
}

// fetch takes a URL, GETs it, and returns the body as a string
func fetch(url string) (string, error) {
	// figure out the file type
	//filetype, err := fetch_filetype(url)

	// GET the url
	resp, err := http.Get(url)
	if err != nil {
		return "", err
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
