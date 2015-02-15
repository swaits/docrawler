package main

import (
	"io/ioutil"
	"net/http"
)

// fetch takes a URL, GETs it, and returns the body as a string
func fetch(url string) (string, error) {
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
