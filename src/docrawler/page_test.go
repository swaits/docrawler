package main

import (
	"testing"
)

// TestBadNewPage tests NewPage to see if we are getting expected failures
func TestBadNewPage(t *testing.T) {
	_, err := NewPage(nil, "h%20t\tp://doesntexist23492387492837492374982734.com/")
	if err == nil {
		t.Error("tried creating a NewPage with an invalid URL but didn't get an error")
	}
}
