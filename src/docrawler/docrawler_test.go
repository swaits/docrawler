package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
)

const (
	port    = 8765
	baseURL = "http://localhost:8765/"
)

// TestMain is used so that we can setup an http server, run tests against it, and tear it down
func TestMain(m *testing.M) {
	// start our simple web server
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic("unable to open port for http server") // TODO remove panic?
	}
	http.Handle("/", http.FileServer(http.Dir("./testsite")))
	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			panic("unable to start http server") // TODO remove panic?
		}
	}()

	// run tests
	exitcode := m.Run()

	// tear down web server
	l.Close()

	// tell OS our result
	os.Exit(exitcode)
}

// TestServerRunning verifies we can get the baseURL from our test harness http server
func TestServerRunning(t *testing.T) {
	resp, err := http.Get(baseURL)
	if resp.StatusCode != 200 || err != nil {
		fmt.Printf("got %v, err = %v\n", resp.StatusCode, err)
		t.Fatal(err)
	}
}

// TestBadRequest verifies we get a 404 for a request to a missing file
func TestBadRequest(t *testing.T) {
	resp, err := http.Get(baseURL + "doesnotexist.html")
	if resp.StatusCode != 404 || err != nil {
		fmt.Printf("got %v, err = %v\n", resp.StatusCode, err)
		t.Fatal("request for non existent file didn't fail properly")
	}
}

// TestSimpleFetch makes sure we can fetch a file and we get exactly what we expect
func TestSimpleFetch(t *testing.T) {
	body, err := fetch(baseURL + "fetch_test.html")
	if err != nil {
		t.Fatal(err)
	}
	desired := "<html><head></head><body></body></html>\n"
	if body != desired {
		fmt.Printf("   Got: '%v'\n", body)
		fmt.Printf("Wanted: '%v'\n", desired)
		t.Error("fetch succeeded, but content mismatched")
	}
}

// TestSimpleParse verifies that we can extract URLs (etc.) from a known document
func TestSimpleParse(t *testing.T) {
	doc := `<html>
	<head>
	</head>
	<body>
		<img src="/assets/image.png"/>
		<a href="/about.html">
	</body>
</html>`
	matches, err := parse(doc)
	if err != nil {
		t.Error(err)
	}
	if len(matches) != 2 {
		t.Fatal("invalid number of matches in parse")
	}
	if matches[0] != "/about.html" {
		t.Error("match text is invalid")
	}
	if matches[1] != "/assets/image.png" {
		t.Error("match text is invalid")
	}
}

// TestSimpleMap figures out the site map for the site in baseURL
func TestSimpleMap(t *testing.T) {
	pages := docrawl(baseURL)
	if len(pages) != 1 {
		t.Fatal("got wrong number of pages")
	}
	if pages[0].URL != baseURL {
		t.Error("page URL is invalid")
	}
	if len(pages[0].Assets) != 1 {
		t.Fatal("got wrong number of assets")
	}
	if pages[0].Assets[0] != baseURL+"/assets/image.png" {
		t.Error("asset name is incorrect")
	}
	if len(pages[0].Links) != 1 {
		t.Fatal("got wrong number of links")
	}
	if pages[0].Links[0] != baseURL+"/about.html" {
		t.Error("link name is incorrect")
	}
}
