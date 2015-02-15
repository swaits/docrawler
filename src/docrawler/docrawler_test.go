package main

import (
	"net"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// start our simple web server
	l, err := net.Listen("tcp", ":8765")
	if err != nil {
		panic("unable to listen on port 8765") // TODO remove panic?
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

func TestServerRunning(t *testing.T) {
	resp, err := http.Get("http://localhost:8765/")
	if resp.StatusCode != 200 || err != nil {
		t.Error(err)
	}
}
