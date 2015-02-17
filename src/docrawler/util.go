package main

// uniqStrings takes a slice of strings and removes any duplicates
// note: does not guarantee any order (or stability of order)
func uniqStrings(strs []string) []string {
	// handle special case of a nil parameter
	if strs == nil {
		return nil
	}

	// add each string to the slice as keys (used like a "set")
	set := make(map[string]struct{})
	for _, s := range strs {
		set[s] = struct{}{}
	}

	// now copy all the keys over to a new slice and return it
	result := []string{}
	for key := range set {
		result = append(result, key)
	}
	return result
}
