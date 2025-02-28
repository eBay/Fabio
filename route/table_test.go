package route

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const (
	// helper constants for the Lookup function
	globEnabled  = false
	globDisabled = true
)

// Global GlobCache for Testing
var globCache = NewGlobCache(1000)

func TestTableParse(t *testing.T) {
	genRoutes := func(n int, format string) (a []string) {
		for i := 0; i < n; i++ {
			a = append(a, fmt.Sprintf(format, i))
		}
		return a
	}

	tests := []struct {
		desc    string
		in, out []string
	}{

		{"1 service, 1 prefix",
			[]string{
				`route add svc-a / http://aaa.com/`,
			},
			[]string{
				`route add svc-a / http://aaa.com/ weight 1.0000`,
			},
		},

		{"1 service, 1 prefix, 3 instances",
			[]string{
				`route add svc-a / http://aaa.com:1111/`,
				`route add svc-a / http://aaa.com:2222/`,
				`route add svc-a / http://aaa.com:3333/`,
			},
			[]string{
				`route add svc-a / http://aaa.com:1111/ weight 0.3333`,
				`route add svc-a / http://aaa.com:2222/ weight 0.3333`,
				`route add svc-a / http://aaa.com:3333/ weight 0.3333`,
			},
		},

		{"1 service, 1 prefix with option",
			[]string{
				`route add svc-a / http://aaa.com/ opts "strip=/foo"`,
				`route add svc-b / http://bbb.com/ opts "strip=/bar"`,
			},
			[]string{
				`route add svc-a / http://aaa.com/ weight 0.5000 opts "strip=/foo"`,
				`route add svc-b / http://bbb.com/ weight 0.5000 opts "strip=/bar"`,
			},
		},

		{"1 service, 1 prefix, 2 instances with different options",
			[]string{
				`route add svc-a / http://aaa.com/ opts "strip=/foo"`,
				`route add svc-b / http://bbb.com/ opts "strip=/bar"`,
			},
			[]string{
				`route add svc-a / http://aaa.com/ weight 0.5000 opts "strip=/foo"`,
				`route add svc-b / http://bbb.com/ weight 0.5000 opts "strip=/bar"`,
			},
		},

		{"2 service, 1 prefix",
			[]string{
				`route add svc-a / http://aaa.com/`,
				`route add svc-b / http://bbb.com/`,
			},
			[]string{
				`route add svc-a / http://aaa.com/ weight 0.5000`,
				`route add svc-b / http://bbb.com/ weight 0.5000`,
			},
		},

		{"1 service, 2 prefix",
			[]string{
				`route add svc-a /one http://aaa.com/`,
				`route add svc-a /two http://aaa.com/`,
			},
			[]string{
				`route add svc-a /two http://aaa.com/ weight 1.0000`,
				`route add svc-a /one http://aaa.com/ weight 1.0000`,
			},
		},

		{"2 service, 2 prefix",
			[]string{
				`route add svc-a /a http://aaa.com/`,
				`route add svc-b /b http://bbb.com/`,
			},
			[]string{
				`route add svc-b /b http://bbb.com/ weight 1.0000`,
				`route add svc-a /a http://aaa.com/ weight 1.0000`,
			},
		},

		{"sort by more specific prefix",
			[]string{
				`route add svc-a / http://aaa.com/`,
				`route add svc-b /b http://bbb.com/`,
			},
			[]string{
				`route add svc-b /b http://bbb.com/ weight 1.0000`,
				`route add svc-a / http://aaa.com/ weight 1.0000`,
			},
		},

		{"sort prefix with host before prefix without host",
			[]string{
				`route add svc-a / http://aaa.com/`,
				`route add svc-b b.com/ http://bbb.com/`,
			},
			[]string{
				`route add svc-b b.com/ http://bbb.com/ weight 1.0000`,
				`route add svc-a / http://aaa.com/ weight 1.0000`,
			},
		},

		{"add more specific prefix to existing host",
			[]string{
				`route add svc-a a.com/ http://aaa.com/`,
				`route add svc-a a.com/a http://aaa.com/`,
			},
			[]string{
				`route add svc-a a.com/a http://aaa.com/ weight 1.0000`,
				`route add svc-a a.com/ http://aaa.com/ weight 1.0000`,
			},
		},

		{"delete route by service, path and target",
			[]string{
				`route add svc-a / http://aaa.com/`,
				`route add svc-b / http://bbb.com/`,
				`route del svc-b / http://bbb.com/`,
			},
			[]string{
				`route add svc-a / http://aaa.com/ weight 1.0000`,
			},
		},

		{"delete route by service and path",
			[]string{
				`route add svc-a / http://aaa.com/`,
				`route add svc-a / http://aaa.com:2222/`,
				`route add svc-b / http://bbb.com/`,
				`route del svc-a /`,
			},
			[]string{
				`route add svc-b / http://bbb.com/ weight 1.0000`,
			},
		},

		{"delete route by service",
			[]string{
				`route add svc-a /a http://aaa.com/`,
				`route add svc-a / http://aaa.com/`,
				`route add svc-b / http://bbb.com/`,
				`route del svc-a`,
			},
			[]string{
				`route add svc-b / http://bbb.com/ weight 1.0000`,
			},
		},

		{"delete route by service and tags",
			[]string{
				`route add svc-a /a http://aaa.com/ tags "a,b"`,
				`route add svc-a /  http://aaa.com/ tags "b,c"`,
				`route add svc-b /  http://bbb.com/ tags "c,d"`,
				`route del svc-a tags "a,b"`,
			},
			[]string{
				`route add svc-a / http://aaa.com/ weight 0.5000 tags "b,c"`,
				`route add svc-b / http://bbb.com/ weight 0.5000 tags "c,d"`,
			},
		},

		{"delete route by tags",
			[]string{
				`route add svc-a /a http://aaa.com/ tags "a,b"`,
				`route add svc-a /  http://aaa.com/ tags "b,c"`,
				`route add svc-b /  http://bbb.com/ tags "c,d"`,
				`route del tags "b"`,
			},
			[]string{
				`route add svc-b / http://bbb.com/ weight 1.0000 tags "c,d"`,
			},
		},

		{"weigh fixed weight 0 -> auto distribution",
			[]string{
				`route add svc / http://bar:111/ weight 0`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 1.0000`,
			},
		},

		{"weigh only fixed weights and sum(fixedWeight) < 1 -> normalize to 100%",
			[]string{
				`route add svc / http://bar:111/ weight 0.2`,
				`route add svc / http://bar:222/ weight 0.3`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 0.4000`,
				`route add svc / http://bar:222/ weight 0.6000`,
			},
		},

		{"weigh only fixed weights and sum(fixedWeight) > 1 -> normalize to 100%",
			[]string{
				`route add svc / http://bar:111/ weight 2`,
				`route add svc / http://bar:222/ weight 3`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 0.4000`,
				`route add svc / http://bar:222/ weight 0.6000`,
			},
		},

		{"weigh multiple entries for same instance with no fixed weight -> de-duplication",
			[]string{
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:111/`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 1.0000`,
			},
		},

		{"weigh multiple entries with no fixed weight -> even distribution",
			[]string{
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:222/`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 0.5000`,
				`route add svc / http://bar:222/ weight 0.5000`,
			},
		},

		{"weigh multiple entries with de-dup and no fixed weight -> even distribution",
			[]string{
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:222/`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 0.5000`,
				`route add svc / http://bar:222/ weight 0.5000`,
			},
		},

		{"weigh mixed fixed and auto weights -> even distribution of remaining weight across non-fixed weighted targets",
			[]string{
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:222/`,
				`route add svc / http://bar:333/ weight 0.5`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 0.2500`,
				`route add svc / http://bar:222/ weight 0.2500`,
				`route add svc / http://bar:333/ weight 0.5000`,
			},
		},

		{"weigh fixed weight == 100% -> route only to fixed weighted targets",
			[]string{
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:222/ weight 0.2500`,
				`route add svc / http://bar:333/ weight 0.7500`,
			},
			[]string{
				`route add svc / http://bar:222/ weight 0.2500`,
				`route add svc / http://bar:333/ weight 0.7500`,
			},
		},

		{"weigh fixed weight > 100%  -> route only to fixed weighted targets and normalize weight",
			[]string{
				`route add svc / http://bar:111/`,
				`route add svc / http://bar:222/ weight 1`,
				`route add svc / http://bar:333/ weight 3`,
			},
			[]string{
				`route add svc / http://bar:222/ weight 0.2500`,
				`route add svc / http://bar:333/ weight 0.7500`,
			},
		},

		{"weigh dynamic weight matched on service name",
			[]string{
				`route add svca / http://bar:111/`,
				`route add svcb / http://bar:222/`,
				`route add svcb / http://bar:333/`,
				`route weight svcb / weight 0.1`,
			},
			[]string{
				`route add svca / http://bar:111/ weight 0.9000`,
				`route add svcb / http://bar:222/ weight 0.0500`,
				`route add svcb / http://bar:333/ weight 0.0500`,
			},
		},

		{"weigh dynamic weight matched on service name and tags",
			[]string{
				`route add svc / http://bar:111/ tags "a"`,
				`route add svc / http://bar:222/ tags "b"`,
				`route add svc / http://bar:333/ tags "b"`,
				`route weight svc / weight 0.1 tags "b"`,
			},
			[]string{
				`route add svc / http://bar:111/ weight 0.9000 tags "a"`,
				`route add svc / http://bar:222/ weight 0.0500 tags "b"`,
				`route add svc / http://bar:333/ weight 0.0500 tags "b"`,
			},
		},

		{"weigh dynamic weight matched on tags",
			[]string{
				`route add svca / http://bar:111/ tags "a"`,
				`route add svcb / http://bar:222/ tags "b"`,
				`route add svcb / http://bar:333/ tags "b"`,
				`route weight / weight 0.1 tags "b"`,
			},
			[]string{
				`route add svca / http://bar:111/ weight 0.9000 tags "a"`,
				`route add svcb / http://bar:222/ weight 0.0500 tags "b"`,
				`route add svcb / http://bar:333/ weight 0.0500 tags "b"`,
			},
		},

		{"weigh more than 1000 routes",
			genRoutes(1234, `route add svc / http://bar:%d/`),
			genRoutes(1234, `route add svc / http://bar:%d/ weight 0.0008`),
		},

		{"weigh more than 1000 routes with a fixed route target",
			func() (a []string) {
				a = genRoutes(1234, `route add svc / http://bar:%d/`)
				a = append(a, `route add svc / http://static:12345/ tags "a"`)
				a = append(a, `route weight svc / weight 0.2 tags "a"`)
				return a
			}(),
			func() (a []string) {
				a = genRoutes(1234, `route add svc / http://bar:%d/ weight 0.0006`)
				a = append(a, `route add svc / http://static:12345/ weight 0.2000 tags "a"`)
				return a
			}(),
		},
	}

	atof := func(s string) float64 {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return n
	}

	for _, tt := range tests {
		// perform a test which parses the tt.in routes into a table and
		// compares the weighted, generated routing table with tt.out. verify,
		// that the distribution of the target URLs for each prefix in the
		// generated routing table matches the weight This test assumes that
		// the table generates the correct routing table but does not test the
		// actual lookup which it probably should.
		t.Run(tt.desc, func(t *testing.T) {
			// parse the routes
			tbl, err := NewTable(bytes.NewBufferString(strings.Join(tt.in, "\n")))
			if err != nil {
				t.Fatalf("got %v want nil", err)
			}

			// compare the generated routes with the normalized weights
			if got, want := tbl.config(true), tt.out; !reflect.DeepEqual(got, want) {
				t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}

			// check that the weights returned in the generated config match
			// the distribution in the wTargets array of the corresponding route.
			checked := map[string]bool{}

			for _, s := range tt.out {
				// route add <svc> <path> ...
				path := strings.Fields(s)[3]

				// if we have already checked this path then skip this.
				// Otherwise, this test becomes O(n^2) and will time out for
				// the large number of routes.
				if checked[path] {
					continue
				}
				checked[path] = true

				// fetch the route
				r := tbl.route(hostpath(path))
				if r == nil {
					t.Fatalf("got nil want route %s", path)
				}

				// check that there are at least some slots
				if len(r.wTargets) == 0 {
					t.Fatalf("got 0 targets want some")
				}

				// pre-generate the target urls for comparison as this
				// will otherwise slow the test down significantly
				targetURLs := make([]string, len(r.wTargets))
				for i, tg := range r.wTargets {
					targetURLs[i] = tg.URL.Scheme + "://" + tg.URL.Host + tg.URL.Path
				}

				// count how often the 'url' from 'route add svc <path> <url>'
				// appears in the list of wTargets for all the URLs
				// from the routes to determine whether the actual
				// distribution of each target within the wTarget slice
				// matches what we expect
				for _, s := range tt.out {
					// route add <svc> <path> <url> weight <weight> ...`,
					p := strings.Fields(s)

					// skip if the path doesn't match
					if path != p[3] {
						continue
					}

					// count how often the target url appears in the list of wTargets
					count := 0
					for _, u := range targetURLs {
						if u == p[4] {
							count++
						}
					}

					// calc the weight as nSlots/totalSlots
					gotWeight := float64(count) / float64(len(r.wTargets))

					// round the weight down to the number of decimal points
					// supported by maxSlots
					gotWeight = float64(int(gotWeight*float64(maxSlots))) / float64(maxSlots)

					// compare to the weight from the generated config
					wantWeight := atof(p[6])

					// check that the actual weight is within 2% of the computed weight
					if math.Abs(gotWeight-wantWeight) > 0.02 {
						t.Errorf("got weight %f want %f", gotWeight, wantWeight)
					}

					// TODO(fs): verify distriibution of targets across the ring
					// TODO(fs): verify lookup with 'rr' works as expected. Current test is by proxy of generated config.
				}
			}
		})
	}
}

func TestNormalizeHost(t *testing.T) {
	tests := []struct {
		req  *http.Request
		host string
	}{
		{&http.Request{Host: "foo.com"}, "foo.com"},
		{&http.Request{Host: "foo.com:80"}, "foo.com"},
		{&http.Request{Host: "foo.com:81"}, "foo.com:81"},
		{&http.Request{Host: "foo.com", TLS: &tls.ConnectionState{}}, "foo.com"},
		{&http.Request{Host: "foo.com:443", TLS: &tls.ConnectionState{}}, "foo.com"},
		{&http.Request{Host: "foo.com:444", TLS: &tls.ConnectionState{}}, "foo.com:444"},
	}

	for i, tt := range tests {
		if got, want := normalizeHost(tt.req.Host, tt.req.TLS != nil), tt.host; got != want {
			t.Errorf("%d: got %v want %v", i, got, want)
		}
	}
}

// see https://github.com/fabiolb/fabio/issues/448
// for more information on the issue and purpose of this test
func TestTableLookupIssue448(t *testing.T) {
	s := `
	route add mock foo.com:80/ https://foo.com/ opts "redirect=301"
	route add mock aaa.com:80/ http://bbb.com/ opts "redirect=301"
	route add mock ccc.com:443/bar https://ccc.com/baz opts "redirect=301"
	route add mock / http://foo.com/
	`

	tbl, err := NewTable(bytes.NewBufferString(s))
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		req         *http.Request
		dst         string
		globEnabled bool
	}{
		{
			req: &http.Request{
				Host: "foo.com",
				URL:  mustParse("/"),
			},
			dst: "https://foo.com/",
			// empty upstream header should follow redirect - standard behavior
		},
		{
			req: &http.Request{
				Host:   "foo.com",
				URL:    mustParse("/"),
				Header: http.Header{"X-Forwarded-Proto": {"http"}},
			},
			dst: "https://foo.com/",
			// upstream http request to same host and path should follow redirect
		},
		{
			req: &http.Request{
				Host:   "foo.com",
				URL:    mustParse("/"),
				Header: http.Header{"X-Forwarded-Proto": {"https"}},
				TLS:    &tls.ConnectionState{},
			},
			dst: "http://foo.com/",
			// upstream https request to same host and path should NOT follow redirect"
		},
		{
			req: &http.Request{
				Host:   "aaa.com",
				URL:    mustParse("/"),
				Header: http.Header{"X-Forwarded-Proto": {"http"}},
			},
			dst: "http://bbb.com/",
			// upstream http request to different http host should follow redirect
		},
		{
			req: &http.Request{
				Host:   "ccc.com",
				URL:    mustParse("/bar"),
				Header: http.Header{"X-Forwarded-Proto": {"https"}},
				TLS:    &tls.ConnectionState{},
			},
			dst: "https://ccc.com/baz",
			// upstream https request to same https host with different path should follow redirect"
		},
	}

	for i, tt := range tests {
		if got, want := tbl.Lookup(tt.req, rndPicker, prefixMatcher, globCache, globEnabled).URL.String(), tt.dst; got != want {
			t.Errorf("%d: got %v want %v", i, got, want)
		}
	}
}

func TestTableLookup(t *testing.T) {
	s := `
	route add svc / http://foo.com:800
	route add svc /foo http://foo.com:900
	route add svc abc.com/ http://foo.com:1000
	route add svc abc.com/foo http://foo.com:1500
	route add svc abc.com/foo/ http://foo.com:2000
	route add svc abc.com/foo/bar http://foo.com:2500
	route add svc abc.com/foo/bar/ http://foo.com:3000
	route add svc z.abc.com/foo/ http://foo.com:3100
	route add svc *.abc.com/ http://foo.com:4000
	route add svc *.abc.com/foo/ http://foo.com:5000
	route add svc *.aaa.abc.com/ http://foo.com:6000
	route add svc *.bbb.abc.com/ http://foo.com:6100
	route add svc xyz.com:80/ https://xyz.com
	`

	tbl, err := NewTable(bytes.NewBufferString(s))
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		req         *http.Request
		dst         string
		globEnabled bool
	}{
		// match on host and path with and without trailing slash
		{&http.Request{Host: "abc.com", URL: mustParse("/")}, "http://foo.com:1000", globEnabled},
		{&http.Request{Host: "abc.com", URL: mustParse("/bar")}, "http://foo.com:1000", globEnabled},
		{&http.Request{Host: "abc.com", URL: mustParse("/foo")}, "http://foo.com:1500", globEnabled},
		{&http.Request{Host: "abc.com", URL: mustParse("/foo/")}, "http://foo.com:2000", globEnabled},
		{&http.Request{Host: "abc.com", URL: mustParse("/foo/bar")}, "http://foo.com:2500", globEnabled},
		{&http.Request{Host: "abc.com", URL: mustParse("/foo/bar/")}, "http://foo.com:3000", globEnabled},

		// do not match on host but maybe on path
		{&http.Request{Host: "def.com", URL: mustParse("/")}, "http://foo.com:800", globEnabled},
		{&http.Request{Host: "def.com", URL: mustParse("/bar")}, "http://foo.com:800", globEnabled},
		{&http.Request{Host: "def.com", URL: mustParse("/foo")}, "http://foo.com:900", globEnabled},

		// strip default port
		{&http.Request{Host: "abc.com:80", URL: mustParse("/")}, "http://foo.com:1000", globEnabled},
		{&http.Request{Host: "abc.com:443", URL: mustParse("/"), TLS: &tls.ConnectionState{}}, "http://foo.com:1000", globEnabled},

		// not using default port
		{&http.Request{Host: "abc.com:443", URL: mustParse("/")}, "http://foo.com:800", globEnabled},
		{&http.Request{Host: "abc.com:80", URL: mustParse("/"), TLS: &tls.ConnectionState{}}, "http://foo.com:800", globEnabled},

		// glob match the host
		{&http.Request{Host: "x.abc.com", URL: mustParse("/")}, "http://foo.com:4000", globEnabled},
		{&http.Request{Host: "y.abc.com", URL: mustParse("/abc")}, "http://foo.com:4000", globEnabled},
		{&http.Request{Host: "x.abc.com", URL: mustParse("/foo/")}, "http://foo.com:5000", globEnabled},
		{&http.Request{Host: "y.abc.com", URL: mustParse("/foo/")}, "http://foo.com:5000", globEnabled},
		{&http.Request{Host: ".abc.com", URL: mustParse("/foo/")}, "http://foo.com:5000", globEnabled},
		{&http.Request{Host: "x.y.abc.com", URL: mustParse("/foo/")}, "http://foo.com:5000", globEnabled},
		{&http.Request{Host: "y.abc.com:80", URL: mustParse("/foo/")}, "http://foo.com:5000", globEnabled},
		{&http.Request{Host: "x.aaa.abc.com", URL: mustParse("/")}, "http://foo.com:6000", globEnabled},
		{&http.Request{Host: "x.aaa.abc.com", URL: mustParse("/foo")}, "http://foo.com:6000", globEnabled},
		{&http.Request{Host: "x.bbb.abc.com", URL: mustParse("/")}, "http://foo.com:6100", globEnabled},
		{&http.Request{Host: "x.bbb.abc.com", URL: mustParse("/foo")}, "http://foo.com:6100", globEnabled},
		{&http.Request{Host: "y.abc.com:443", URL: mustParse("/foo/"), TLS: &tls.ConnectionState{}}, "http://foo.com:5000", globEnabled},

		// exact match has precedence over glob match
		{&http.Request{Host: "z.abc.com", URL: mustParse("/foo/")}, "http://foo.com:3100", globEnabled},

		// explicit port on route
		{&http.Request{Host: "xyz.com", URL: mustParse("/")}, "https://xyz.com", globEnabled},
	}

	for i, tt := range tests {
		if got, want := tbl.Lookup(tt.req, rndPicker, prefixMatcher, globCache, globEnabled).URL.String(), tt.dst; got != want {
			t.Errorf("%d: got %v want %v", i, got, want)
		}
	}
}

func TestTableLookup_656(t *testing.T) {
	// A typical HTTPS redirect
	s := `
	route add my-service example.com:80/ https://example.com$path opts "redirect=301"
	route add my-service example.com/ http://127.0.0.1:3000/
	`

	tbl, err := NewTable(bytes.NewBufferString(s))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	target := tbl.Lookup(req, rrPicker, prefixMatcher, globCache, globDisabled)

	if target == nil {
		t.Fatal("No route match")
	}
	if got, want := target.RedirectCode, 301; got != want {
		t.Errorf("target.RedirectCode = %d, want %d", got, want)
	}
	if got, want := fmt.Sprint(target.RedirectURL), "https://example.com/foo"; got != want {
		t.Errorf("target.RedirectURL = %s, want %s", got, want)
	}
}

func TestNewTableCustom(t *testing.T) {

	var routes []RouteDef
	var tags = []string{"tag1", "tag2"}
	var opts = make(map[string]string)
	opts["tlsskipverify"] = "true"
	opts["proto"] = "http"

	var route1 = RouteDef{
		Cmd:     "route add",
		Service: "service1",
		Src:     "app.com",
		Dst:     "http://10.1.1.1:8080",
		Weight:  0.50,
		Tags:    tags,
		Opts:    opts,
	}
	var route2 = RouteDef{
		Cmd:     "route add",
		Service: "service1",
		Src:     "app.com",
		Dst:     "http://10.1.1.2:8080",
		Weight:  0.50,
		Tags:    tags,
		Opts:    opts,
	}
	var route3 = RouteDef{
		Cmd:     "route add",
		Service: "service2",
		Src:     "app.com",
		Dst:     "http://10.1.1.3:8080",
		Weight:  0.25,
		Tags:    tags,
		Opts:    opts,
	}

	routes = append(routes, route1)
	routes = append(routes, route2)
	routes = append(routes, route3)

	table, err := NewTableCustom(&routes)

	if err != nil {
		fmt.Printf("Got error from NewTableCustom - %s", err.Error())
		t.FailNow()
	}

	tableString := table.String()
	if !strings.Contains(tableString, route1.Dst) {
		fmt.Printf("Table Missing Destination %s -- Table -- %s", route1.Dst, tableString)
		t.FailNow()
	}

	if !strings.Contains(tableString, route2.Dst) {
		fmt.Printf("Table Missing Destination %s -- Table -- %s", route1.Dst, tableString)
		t.FailNow()
	}

	if !strings.Contains(tableString, route3.Dst) {
		fmt.Printf("Table Missing Destination %s -- Table -- %s", route1.Dst, tableString)
		t.FailNow()
	}
}

func TestTable_Dump(t *testing.T) {
	s := `
	route add svc / http://foo.com:800
	route add svc /foo http://foo.com:900
	route add svc abc.com/ http://foo.com:1000
	`

	tbl, err := NewTable(bytes.NewBufferString(s))
	if err != nil {
		t.Fatal(err)
	}

	want := `+-- host=
|   |-- path=/foo
|   |    +-- addr=foo.com:900 weight 1.00 slots 1/1
|   +-- path=/
|       +-- addr=foo.com:800 weight 1.00 slots 1/1
+-- host=abc.com
    +-- path=/
        +-- addr=foo.com:1000 weight 1.00 slots 1/1
`

	got := tbl.Dump()

	if want != got {
		t.Errorf("Unexpected Dump() output:\nwant:\n%s\ngot:\n%s\n", want, got)
	}
}
