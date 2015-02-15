package main

import ()

// Page is a struct which defines a single page, which URLs (links and assets) it contains
type Page struct {
	URL    string
	Assets []string
	Links  []string
}

// docrawl begins crawling the site at "url"
func docrawl(url string) []Page {
	var result []Page
	return result
}

func main() {
}
