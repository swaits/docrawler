# DoCrawler #

This is a toy web crawler written in Go.

It crawls a single host (i.e. anything.com) and outputs a site map in JSON format.

For each page crawled, it distinguishes between links to other pages, links to assets, broken links, and remote links (i.e. someotherhost.com).

### An Example

    ❯ bin/docrawler https://goregex.com/
    
    D.O. Crawler 1.0  Copyright (c) 2015 Stephen Waits <steve@waits.net>  2015-02-17
    
    [
      {
        "URL": "https://goregex.com/",
        "Title": "GoRegEx.com | Go Regular Expression Tester",
        "Links": [
          "https://goregex.com/"
        ],
        "Assets": [
          "https://goregex.com/assets/css/bootstrap-responsive.min.css",
          "https://goregex.com/assets/css/bootstrap.min.css",
          "https://goregex.com/assets/css/goregex.css",
          "https://goregex.com/assets/img/apple-touch-icon-114x114.png",
          "https://goregex.com/assets/img/apple-touch-icon-120x120.png",
          "https://goregex.com/assets/img/apple-touch-icon-144x144.png",
          "https://goregex.com/assets/img/apple-touch-icon-152x152.png",
          "https://goregex.com/assets/img/apple-touch-icon-180x180.png",
          "https://goregex.com/assets/img/apple-touch-icon-57x57.png",
          "https://goregex.com/assets/img/apple-touch-icon-60x60.png",
          "https://goregex.com/assets/img/apple-touch-icon-72x72.png",
          "https://goregex.com/assets/img/apple-touch-icon-76x76.png",
          "https://goregex.com/assets/img/favicon-160x160.png",
          "https://goregex.com/assets/img/favicon-16x16.png",
          "https://goregex.com/assets/img/favicon-192x192.png",
          "https://goregex.com/assets/img/favicon-32x32.png",
          "https://goregex.com/assets/img/favicon-96x96.png",
          "https://goregex.com/assets/img/favicon.ico",
          "https://goregex.com/assets/js/goregex.js"
        ],
        "Broken": null,
        "Remote": [
          "http://code.google.com/p/re2/wiki/Syntax",
          "http://golang.org/",
          "http://golang.org/pkg/regexp/",
          "https://ajax.googleapis.com/ajax/libs/jquery/1.8.2/jquery.min.js",
          "https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"
        ]
      }
    ]

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

This will `fmt`, `lint`, `vet`, and `build` the source into an executable at `bin/docrawler`.

In order to run the tests:

    make test

If you'd like to automatically run tests any time a file is changed:

    make autotest

To see stats about the code, install cloc (`brew install cloc` on OS X) and then:

    make stats

To look at the test coverage reports in your browser:

    make cover

To deploy, just copy the executable `bin/docrawler` to a destination of your choosing.

*Developed on OS X. Also tested on Windows 7 under MinGW+bash.*

### Why a Makefile? Aren't those old???

* It's clean.
* It's portable.
* It lets me package the app in a way that's trivially easy for anyone to build.
* It makes it so I can easily vendor third-party packages outside of the main `./src` tree.

## Who do I talk to? ##

Written by Stephen Waits. Please contact me at <mailto:steve@waits.net> with any questions.
