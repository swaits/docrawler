package main

import (
	"encoding/json"
	"sort"
)

// Location is a struct which defines a single URL, which URLs (links and assets) it contains, etc.
type Location struct {
	URL    string
	Title  string
	Links  []string
	Assets []string
	Broken []string
	Remote []string
}

func sitemapToLocations(pages []*httpItem) []*Location {
	// build a map of our pages
	pageMap := make(map[string]*httpItem)
	for _, p := range pages {
		pageMap[p.url.String()] = p
	}

	// build a slice of locations (one per page)
	var locations []*Location
	for _, p := range pages {
		if p.mediaType == "text/html" {
			// create a location for this page
			l := &Location{URL: p.url.String(), Title: p.title}

			// add its children
			for _, c := range p.children {
				// look up this child's media type from the root list of pages
				//mediaType := pageMap[c.url.String()].mediaType
				if c.skipped {
					l.Remote = append(l.Remote, c.url.String())
				} else if c.mediaType == "text/html" {
					l.Links = append(l.Links, c.url.String())
				} else if c.broken {
					l.Broken = append(l.Broken, c.url.String())
				} else {
					l.Assets = append(l.Assets, c.url.String())
				}
			}

			// now sort the children slices
			sort.Strings(l.Remote)
			sort.Strings(l.Links)
			sort.Strings(l.Broken)
			sort.Strings(l.Assets)

			// and add this location to our slice
			locations = append(locations, l)
		}
	}

	return locations
}

func locationsToJSON(locations []*Location) (string, error) {
	b, err := json.Marshal(locations)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
