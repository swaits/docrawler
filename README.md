# DoCrawler #

This is a toy web crawler written in Go.

It crawls a single host (i.e. anything.com) and outputs a site map in JSON format.

For each page crawled, it distinguishes between links to other pages, links to assets, broken links, and remote links (i.e. someotherhost.com).

## Design ##

* The main crawler loop is in `docrawler.go:docrawl()`.
* From there I maintain a hash of links we've already crawled.
* Each link gets an `httpItem{}` struct instanced, which holds its crawl state.
* A number of crawler goroutines are fired off in the beginning so that we can control precisely how many http fetches happen at a single time. This number is configurable from via command line parameter.
* These goroutines listen for `*httpItem{}` on one channel, crawl it, fill out its results, and send it back to `docrawl()` on another channel.
* `docrawl()` is predominantly a for-select loop, selecting on a ticker (which is used to update status to the console and check for crawl completion), and on the "work finished" channel that the crawlers send data back on.

## How do I get set up? ##

After cloning this repository, run:

    make vendor && make

In order to run the tests:

    make test

If you'd like to automatically run tests any time a file is changed:

    make autotest

To deploy, just copy the executable `bin/docrawler` to a destination of your choosing.

*Developed on OS X. Also tested on Windows 7 under MinGW+bash.*

## Who do I talk to? ##

Written by Stephen Waits. Please contact me at <mailto:steve@waits.net> with any questions.
