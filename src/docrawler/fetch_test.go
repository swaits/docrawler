package main

import (
	"testing"
)

// TestFailedHeaderFetching tests fetchFiletype() to see if we are getting expected failures
func TestFailedHeaderFetching(t *testing.T) {
	page, err := newHTTPItem(nil, "http://doesntexist23492387492837492374982734.com/")
	if err != nil {
		t.Error("problem creating New Page struct")
	}
	if err := fetchFiletype(page); err == nil {
		t.Error("tired fetching bogus page but didn't get nil back from fetchFiletype")
	}
}

// TestFailedFetching tests fetchFiletype() to see if we are getting expected failures
func TestFailedFetching(t *testing.T) {
	page, err := newHTTPItem(nil, "http://doesntexist23492387492837492374982734.com/")
	if err != nil {
		t.Error("problem creating New Page struct")
	}
	if _, err := fetchItem(page); err == nil {
		t.Error("tired fetching bogus page but didn't get nil back from fetchPage")
	}
}
