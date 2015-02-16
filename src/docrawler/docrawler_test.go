package main

import (
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
		t.Logf("got %v, err = %v\n", resp.StatusCode, err)
		t.Fatal(err)
	}
}

// TestBadRequest verifies we get a 404 for a request to a missing file
func TestBadRequest(t *testing.T) {
	resp, err := http.Get(baseURL + "doesnotexist.html")
	if resp.StatusCode != 404 || err != nil {
		t.Logf("got %v, err = %v\n", resp.StatusCode, err)
		t.Fatal("request for non existent file didn't fail properly")
	}
}

// TestSimpleFetch makes sure we can fetch a file and we get exactly what we expect
func TestSimpleFetch(t *testing.T) {
	page, err := NewPage(nil, baseURL+"fetch_test.html")
	if err != nil {
		t.Error("problem creating New Page struct")
	}
	body, err := fetchPage(page)
	if err != nil {
		t.Fatal(err)
	}
	desired := "<html><head></head><body></body></html>\n"
	if body != desired {
		t.Logf("   Got: '%v'\n", body)
		t.Logf("Wanted: '%v'\n", desired)
		t.Error("fetch succeeded, but content mismatched")
	}
}

// TestSimpleMap figures out the site map for the site in baseURL
func TestSimpleMap(t *testing.T) {
	pages := docrawl(baseURL)
	if len(pages) != 4 {
		t.Logf("got %v, wanted %v\n", len(pages), 1)
		t.Fatal("got wrong number of pages")
	}
	if pages[0].URL.String() != baseURL {
		t.Error("page URL is invalid")
	}
	if len(pages[0].Children) != 3 {
		t.Logf("got %v, wanted %v\n", len(pages[0].Children), 1)
		t.Fatal("got wrong number of links")
	}
	if pages[0].Children[0].URL.String() != baseURL+"about.html" {
		t.Logf("   Got: %q\n", pages[0].Children[0].URL.String())
		t.Logf("Wanted: %q\n", baseURL+"about.html")
		t.Error("link name is incorrect")
	}
	if pages[0].Children[1].URL.String() != baseURL+"assets/image.png" {
		t.Logf("   Got: %q\n", pages[0].Children[1].URL.String())
		t.Logf("Wanted: %q\n", baseURL+"assets/image.png")
		t.Error("link name is incorrect")
	}
	if pages[0].Children[2].URL.String() != baseURL+"scripts/blah.js" {
		t.Logf("   Got: %q\n", pages[0].Children[2].URL.String())
		t.Logf("Wanted: %q\n", baseURL+"scripts/blah.js")
		t.Error("link name is incorrect")
	}
}

// TestHeaderFetching tests fetchFiletype() to see if we are getting expected results
func TestHeaderFetching(t *testing.T) {
	basepage, err := NewPage(nil, baseURL)
	if err != nil {
		t.Error("problem creating New Page struct")
	}
	if err := fetchFiletype(basepage); err != nil || basepage.MediaType != "text/html" {
		t.Error("problem fetching filetype")
	}

	page, err := NewPage(basepage, "about.html")
	if err != nil {
		t.Error("problem creating New Page struct")
	}
	if err := fetchFiletype(page); err != nil || page.MediaType != "text/html" {
		t.Error("problem fetching filetype")
	}

	page, err = NewPage(basepage, "assets/image.png")
	if err != nil {
		t.Error("problem creating New Page struct")
	}
	if err := fetchFiletype(page); err != nil || page.MediaType != "text/plain" {
		t.Logf("got %q, wanted %q", page.MediaType, "text/plain")
		t.Error("problem fetching filetype")
	}
}
