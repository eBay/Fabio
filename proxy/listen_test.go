package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/route"
)

func TestGracefulShutdown(t *testing.T) {

	// start a server which responds after the shutdown has been triggered.
	trigger := make(chan bool)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-trigger
		return
	}))
	defer srv.Close()

	// start proxy
	addr := "127.0.0.1:57777"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		h := &HTTPProxy{
			Transport: http.DefaultTransport,
			Lookup: func(r *http.Request) *route.Target {
				tbl, _ := route.NewTable(bytes.NewBufferString("route add svc / " + srv.URL))
				return tbl.Lookup(r, route.Picker["rr"], route.Matcher["prefix"], globCache, globEnabled)
			},
		}
		l := config.Listen{Addr: addr}
		if err := ListenAndServeHTTP(l, h, nil); err != nil {
			t.Log("ListenAndServeHTTP: ", err)
		}
	}()

	// trigger shutdown after some time
	delay := 100 * time.Millisecond
	go func() {
		time.Sleep(delay)
		close(trigger)
		Shutdown(delay)
	}()

	// give server some time to start up
	time.Sleep(delay / 2)

	makeReq := func() (int, error) {
		resp, err := http.Get("http://" + addr + "/")
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		return resp.StatusCode, nil
	}

	// make 200 OK request
	// start before and complete after shutdown was triggered
	code, err := makeReq()
	if err != nil {
		t.Fatalf("request 1: got error %q want nil", err)
	}
	if got, want := code, 200; got != want {
		t.Fatalf("request 1: got %v want %v", got, want)
	}

	// make request to closed server
	_, err = makeReq()
	if got, want := err.Error(), "connection refused"; !strings.Contains(got, want) {
		t.Fatalf("request 2: got error %q want %q", got, want)
	}

	// wait for listen() to return
	// note that the actual listeners have not returned yet
	wg.Wait()
}
