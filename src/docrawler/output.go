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

func sitemapToLocations(pages []*Page) []*Location {
	// build a map of our pages
	pageMap := make(map[string]*Page)
	for _, p := range pages {
		pageMap[p.URL.String()] = p
	}

	// build a slice of locations (one per page)
	var locations []*Location
	for _, p := range pages {
		if p.MediaType == "text/html" {
			// create a location for this page
			l := &Location{URL: p.URL.String(), Title: p.Title}

			// add its children
			for _, c := range p.Children {
				// look up this child's media type from the root list of pages
				//mediaType := pageMap[c.URL.String()].MediaType
				if c.Skipped == true {
					l.Remote = append(l.Remote, c.URL.String())
				} else if c.MediaType == "text/html" {
					l.Links = append(l.Links, c.URL.String())
				} else if c.MediaType == "" {
					l.Broken = append(l.Broken, c.URL.String())
				} else {
					l.Assets = append(l.Assets, c.URL.String())
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
