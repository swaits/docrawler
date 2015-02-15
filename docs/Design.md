# Crawler Requirements

## Minimum Features

This is the assignment (directly) from DO.

* We do ask that you put a heavy emphasis on testing in the assignment.
* written in a modern language
* limited to one domain - so when crawling digitalocean.com it would crawl all pages within the digitalocean.com domain, but not follow the links to our Facebook or Instagram accounts or subdomains like cloud.digitalocean.com. 
* It should be Given a URL
* it should output a site map
* it should show which static assets each page depends on, and the links between pages.
* Choose the most appropriate data structure to store & display this site map.
* Build this as you would build something for production - it's fine if you don't finish everything, but focus on code quality and write tests.
* We're interested in how you code and how you test your code.

## Desired

These are features I'd like to include.

* efficient, optimal
* handles base URL properly (for root & relative links)
* documented
* uses standard Go practices
* uses Go stdlib (only!)
* thoroughly tested
* configurable on command line (defaults are for company specs)
* variable output formats: json, dot
* adheres to URL RFC (as far as case sensitivity, acceptable character sets, etc.)
* go gettable
* go fmt'ed, vet'ed, lint'ed
* polite - robots.txt support
* robust - detect infinite loops
* throttling
* use Go context pattern
* supports http & https

## Wishlist

These are features that sound cool but are probably out of scope.

* distributed
* sitemap support https://en.wikipedia.org/wiki/Sitemaps
* noindex tag
* noidext http header
