package main

import (
	"testing"
)

// TestBadnewHTTPItem tests newHTTPItem to see if we are getting expected failures
func TestBadnewHTTPItem(t *testing.T) {
	_, err := newHTTPItem(nil, "h%20t\tp://doesntexist23492387492837492374982734.com/")
	if err == nil {
		t.Error("tried creating a newHTTPItem with an invalid URL but didn't get an error")
	}
}
