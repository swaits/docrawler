package main

import (
	"io/ioutil"
	"net/http"
)

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
	resp.Body.Close()
	return string(body), nil
}
