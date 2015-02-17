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
		panic("unable to open port for http server")
	}
	http.Handle("/", http.FileServer(http.Dir("./testsite")))
	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			panic("unable to start http server")
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
	page, err := newHTTPItem(nil, baseURL+"fetch_test.html")
	if err != nil {
		t.Error("problem creating New httpItem struct")
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

	// because our crawl is non-deterministic, we have to do a complete
	// cycle through every page, counting stuff, finding specific pages
	// so that we can then verify certain aspects of each page
	skipped, children := 0, 0
	var homepage, aboutpage, skippage *httpItem
	for _, p := range pages {
		if p.title == "Home" {
			homepage = p
		} else if p.title == "About Test" {
			aboutpage = p
		} else if p.url.String() == "http://doesntexist23492387492837492374982734.com/" {
			skippage = p
		}
		if p.linkType == tRemote {
			skipped++
		}
		children += len(p.children)
	}

	// check page counts
	if len(pages) != 6 {
		t.Logf("got %v, wanted %v\n", len(pages), 6)
		t.Fatal("got wrong number of pages")
	}
	if children != 9 {
		t.Logf("got %v, wanted %v\n", children, 8)
		t.Fatal("got wrong number of total children")
	}
	if skipped != 1 {
		t.Logf("got %v, wanted %v\n", skipped, 1)
		t.Fatal("got wrong number of total skipped")
	}

	// check first page URL
	if homepage.url.String() != baseURL {
		t.Error("page URL is invalid")
	}

	// specific page number of children
	if len(homepage.children) != 5 {
		t.Logf("got %v, wanted %v\n", len(homepage.children), 4)
		t.Fatal("got wrong number of links")
	}
	if len(aboutpage.children) != 4 {
		t.Logf("got %v, wanted %v\n", len(aboutpage.children), 4)
		t.Fatal("got wrong number of links")
	}

	// make sure remote pages were skipped
	if homepage.linkType == tRemote {
		t.Logf("found a page that should NOT have been skipped %q", homepage.children[0].url.String())
		t.Error("child was skipped")
	}
	if skippage.linkType != tRemote {
		t.Logf("found a page that should have been skipped %q", pages[3].children[2].url.String())
		t.Error("child wasn't skipped")
	}

	// verify page content
	// TODO find a better way to do this
	//if homepage.children[0].url.String() != baseURL+"about.html" {
	//t.Logf("   Got: %q\n", homepage.children[0].url.String())
	//t.Logf("Wanted: %q\n", baseURL+"about.html")
	//t.Error("link name is incorrect")
	//}
	//if homepage.children[1].url.String() != baseURL+"assets/image.png" {
	//t.Logf("   Got: %q\n", homepage.children[1].url.String())
	//t.Logf("Wanted: %q\n", baseURL+"assets/image.png")
	//t.Error("link name is incorrect")
	//}
	//if homepage.children[3].url.String() != baseURL+"scripts/blah.js" {
	//t.Logf("   Got: %q\n", homepage.children[3].url.String())
	//t.Logf("Wanted: %q\n", baseURL+"scripts/blah.js")
	//t.Error("link name is incorrect")
	//}
}

// TestHeaderFetching tests fetchFiletype() to see if we are getting expected results
func TestHeaderFetching(t *testing.T) {
	basepage, err := newHTTPItem(nil, baseURL)
	if err != nil {
		t.Error("problem creating New httpItem struct")
	}
	if err := fetchFiletype(basepage); err != nil || basepage.linkType != tHTMLPage {
		t.Error("problem fetching filetype")
	}

	page, err := newHTTPItem(basepage, "about.html")
	if err != nil {
		t.Error("problem creating New httpItem struct")
	}
	if err := fetchFiletype(page); err != nil || page.linkType != tHTMLPage {
		t.Error("problem fetching filetype")
	}

	page, err = newHTTPItem(basepage, "assets/image.png")
	if err != nil {
		t.Error("problem creating New httpItem struct")
	}
	if err := fetchFiletype(page); err != nil || page.linkType != tAsset {
		t.Logf("got %q, wanted %q", page.linkType, tAsset)
		t.Error("problem fetching filetype")
	}
}

// TestJsonOutput gets a sitemap and then converts it to json
func TestJsonOutput(t *testing.T) {
	pages := docrawl(baseURL)
	l := sitemapToLocations(pages)
	if len(l) != 2 {
		t.Error("sitemapToLocations has the wrong number of locations")
	}
	if l[0].URL != baseURL {
		t.Error("got an incorrect location")
	}
	if len(l[0].Links) != 1 {
		t.Error("location 0 has wrong number of links")
	}
	if len(l[0].Assets) != 2 {
		t.Logf("got %v", len(l[0].Assets))
		t.Error("location 0 has wrong number of assets")
	}
	if len(l[0].Broken) != 1 {
		t.Logf("got %v", len(l[0].Broken))
		t.Error("location 0 has wrong number of broken urls")
	}
	if len(l[0].Remote) != 0 {
		t.Logf("got %v", len(l[1].Remote))
		t.Error("location 0 has wrong number of remote urls")
	}
	if l[1].URL != baseURL+"about.html" {
		t.Error("got an incorrect location")
	}
	if len(l[1].Links) != 1 {
		t.Error("location 1 has wrong number of links")
	}
	if len(l[1].Assets) != 2 {
		t.Error("location 1 has wrong number of assets")
	}
	if len(l[1].Broken) != 0 {
		t.Logf("got %v", len(l[1].Broken))
		t.Error("location 1 has wrong number of broken urls")
	}
	if len(l[1].Remote) != 1 {
		t.Logf("got %v", len(l[1].Remote))
		t.Error("location 1 has wrong number of remote urls")
	}

	// this is a cheater test because the output is from a run of the code being
	// test itself. but, it's been examined and I think it's right. and, there's not
	// many other ways to test this without some significant pain
	cheaterTest := `[{"URL":"http://localhost:8765/","Title":"Home","Links":["http://localhost:8765/about.html"],"Assets":["http://localhost:8765/assets/image.png","http://localhost:8765/scripts/blah.js"],"Broken":["http://localhost:8765/zzzbroken.html"],"Remote":null},{"URL":"http://localhost:8765/about.html","Title":"About Test","Links":["http://localhost:8765/"],"Assets":["http://localhost:8765/assets/image.png","http://localhost:8765/scripts/blah.js"],"Broken":null,"Remote":["http://doesntexist23492387492837492374982734.com/"]}]`
	j, err := locationsToJSON(l)
	if err != nil {
		t.Error("locationsToJSON failed")
	}
	if j != cheaterTest {
		t.Error("locationsToJSON text doesn't match")
	}
	t.Log(j)
}
