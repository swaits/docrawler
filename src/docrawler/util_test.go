package main

import (
	"sort"
	"testing"
)

// TestUniqStrings passes a slice of strings to uniqStrings() to see if it removes
// all duplicates
func TestUniqStrings(t *testing.T) {
	strs := []string{"a", "a", "b", "a", "b", "c", "d", "a", "b", "e", "d"}
	uniq := uniqStrings(strs)
	if len(uniq) != 5 {
		t.Error("wrong number of strings after uniq'ing")
	}
	sort.Strings(uniq)
	if uniq[0] != "a" || uniq[1] != "b" || uniq[2] != "c" || uniq[3] != "d" || uniq[4] != "e" {
		t.Error("uniq values are in correct!")
	}
}

// TestUniqStringsNilParm passes a nil slice to uniqStrings to be sure it returns nil
func TestUniqStringsNilParm(t *testing.T) {
	uniq := uniqStrings(nil)
	if uniq != nil {
		t.Error("uniqStrings should return a nil slice when passed a nil slice")
	}
}
