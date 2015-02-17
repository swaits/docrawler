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

// implement Location slice sorting (by URL)
type byURL []*Location

// implement Sort interface on byURL (which is a []*Location)
func (l byURL) Len() int           { return len(l) }
func (l byURL) Less(i, j int) bool { return l[i].URL < l[j].URL }
func (l byURL) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// sitemapToLocations converts an itemSlice to []*Location, which is appropriate for marshalling
func sitemapToLocations(pages itemSlice) []*Location {
	// build a map of our pages
	pageMap := make(map[string]*httpItem)
	for _, p := range pages {
		pageMap[p.url.String()] = p
	}

	// build a slice of locations (one per page)
	var locations []*Location
	for _, p := range pages {
		if p.linkType == tHTMLPage {
			// create a location for this page
			l := &Location{URL: p.url.String(), Title: p.title}

			// add its children
			for _, c := range p.children {
				// look up this child's media type from the root list of pages
				//mediaType := pageMap[c.url.String()].mediaType
				if c.linkType == tRemote {
					l.Remote = append(l.Remote, c.url.String())
				} else if c.linkType == tHTMLPage {
					l.Links = append(l.Links, c.url.String())
				} else if c.linkType == tBroken {
					l.Broken = append(l.Broken, c.url.String())
				} else if c.linkType == tAsset {
					l.Assets = append(l.Assets, c.url.String())
				} else {
					// unknown link here, which means it failed to crawl, let's call it "broken"
					l.Broken = append(l.Broken, c.url.String())
				}

			}

			// now uniq & sort the children slices
			l.Remote = uniqStrings(l.Remote)
			l.Links = uniqStrings(l.Links)
			l.Broken = uniqStrings(l.Broken)
			l.Assets = uniqStrings(l.Assets)
			sort.Strings(l.Remote)
			sort.Strings(l.Links)
			sort.Strings(l.Broken)
			sort.Strings(l.Assets)

			// and add this location to our slice
			locations = append(locations, l)
		}
	}

	// sort the Locations themselves and return
	sort.Sort(byURL(locations))
	return locations
}

// locationsToJSON takes a *Location slice and marshals it into a JSON string
func locationsToJSON(locations []*Location) (string, error) {
	b, err := json.MarshalIndent(locations, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
