# To Do

## Bugs

* deal with anchor issue (repeatedly crawling index.html#about, index.html#contact, etc.)
  * should only crawl these pages once
  * but should be included in the links part
* when sorting links, "uniq" them as well
* be smart about index.html?

## Cleanup

* helper function to do the whole crawl process
* document code better
  * test with godoc
* omit fields in json when nil? (Broken, etc.)

## Features

* command line options, flag parsing
* check RFCs on URL character sets, compliance, etc.
* packagize parts of app
