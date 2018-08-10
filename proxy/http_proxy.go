package proxy

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/logger"
	"github.com/fabiolb/fabio/metrics"
	"github.com/fabiolb/fabio/noroute"
	"github.com/fabiolb/fabio/proxy/fastcgi"
	"github.com/fabiolb/fabio/proxy/gzip"
	"github.com/fabiolb/fabio/route"
	"github.com/fabiolb/fabio/uuid"
)

// HTTPProxy is a dynamic reverse proxy for HTTP and HTTPS protocols.
type HTTPProxy struct {
	// Config is the proxy configuration as provided during startup.
	Config *config.Config

	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, time.Now is used.
	Time func() time.Time

	// Transport is the http connection pool configured with timeouts.
	// The proxy will panic if this value is nil.
	Transport http.RoundTripper

	// InsecureTransport is the http connection pool configured with
	// InsecureSkipVerify set. This is used for https proxies with
	// self-signed certs.
	InsecureTransport http.RoundTripper

	// Lookup returns a target host for the given request.
	// The proxy will panic if this value is nil.
	Lookup func(*http.Request) *route.Target

	// Requests is a timer metric which is updated for every request.
	Requests metrics.Timer

	// Noroute is a counter metric which is updated for every request
	// where Lookup() returns nil.
	Noroute metrics.Counter

	// Logger is the access logger for the requests.
	Logger logger.Logger

	// UUID returns a unique id in uuid format.
	// If UUID is nil, uuid.NewUUID() is used.
	UUID func() string
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Lookup == nil {
		panic("no lookup function")
	}

	if p.Config.Proxy.RequestID != "" {
		id := p.UUID
		if id == nil {
			id = uuid.NewUUID
		}
		r.Header.Set(p.Config.Proxy.RequestID, id())
	}

	t := p.Lookup(r)
	if t == nil {
		status := p.Config.Proxy.NoRouteStatus
		if status < 100 || status > 999 {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		html := noroute.GetHTML()
		if html != "" {
			io.WriteString(w, html)
		}
		return
	}

	if t.AccessDeniedHTTP(r) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// build the request url since r.URL will get modified
	// by the reverse proxy and contains only the RequestURI anyway
	requestURL := &url.URL{
		Scheme:   scheme(r),
		Host:     r.Host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	if t.RedirectCode != 0 && t.RedirectURL != nil {
		http.Redirect(w, r, t.RedirectURL.String(), t.RedirectCode)
		if t.Timer != nil {
			t.Timer.Update(0)
		}
		metrics.DefaultRegistry.GetTimer(key(t.RedirectCode)).Update(0)
		return
	}

	// build the real target url that is passed to the proxy
	targetURL := &url.URL{
		Scheme: t.URL.Scheme,
		Host:   t.URL.Host,
		Path:   r.URL.Path,
	}
	if t.URL.RawQuery == "" || r.URL.RawQuery == "" {
		targetURL.RawQuery = t.URL.RawQuery + r.URL.RawQuery
	} else {
		targetURL.RawQuery = t.URL.RawQuery + "&" + r.URL.RawQuery
	}

	if t.Host == "dst" {
		r.Host = targetURL.Host
	} else if t.Host != "" {
		r.Host = t.Host
	}

	// TODO(fs): The HasPrefix check seems redundant since the lookup function should
	// TODO(fs): have found the target based on the prefix but there may be other
	// TODO(fs): matchers which may have different rules. I'll keep this for
	// TODO(fs): a defensive approach.
	if t.StripPath != "" && strings.HasPrefix(r.URL.Path, t.StripPath) {
		targetURL.Path = targetURL.Path[len(t.StripPath):]
	}

	if err := addHeaders(r, p.Config.Proxy, t.StripPath); err != nil {
		http.Error(w, "cannot parse "+r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	if err := addResponseHeaders(w, r, p.Config.Proxy); err != nil {
		http.Error(w, "cannot add response headers", http.StatusInternalServerError)
		return
	}

	upgrade, accept := r.Header.Get("Upgrade"), r.Header.Get("Accept")

	tr := p.Transport
	if t.TLSSkipVerify {
		tr = p.InsecureTransport
	}

	isFCGI := false
	if v, ok := t.Opts["proto"]; ok && v == "fcgi" {
		isFCGI = true
	}

	var h http.Handler
	switch {
	case isFCGI:
		fcgiProxy := fastcgi.NewProxy(p.Config, targetURL.Host)
		if fcgiRoot, ok := t.Opts["root"]; ok {
			fcgiProxy.SetRoot(fcgiRoot)
		}
		if stripPrefix, ok := t.Opts["strip"]; ok {
			fcgiProxy.SetStripPathPrefix(stripPrefix)
		}
		if indexFile, ok := t.Opts["index"]; ok {
			fcgiProxy.SetIndex(indexFile)
		}

		h = fcgiProxy

	case upgrade == "websocket" || upgrade == "Websocket":
		r.URL = targetURL
		if targetURL.Scheme == "https" || targetURL.Scheme == "wss" {
			h = newWSHandler(targetURL.Host, func(network, address string) (net.Conn, error) {
				return tls.Dial(network, address, tr.(*http.Transport).TLSClientConfig)
			})
		} else {
			h = newWSHandler(targetURL.Host, net.Dial)
		}

	case accept == "text/event-stream":
		// use the flush interval for SSE (server-sent events)
		// must be > 0s to be effective
		h = newHTTPProxy(targetURL, tr, p.Config.Proxy.FlushInterval)

	default:
		h = newHTTPProxy(targetURL, tr, p.Config.Proxy.GlobalFlushInterval)
	}

	if p.Config.Proxy.GZIPContentTypes != nil {
		h = gzip.NewGzipHandler(h, p.Config.Proxy.GZIPContentTypes)
	}

	timeNow := p.Time
	if timeNow == nil {
		timeNow = time.Now
	}

	start := timeNow()
	rw := &responseWriter{w: w}
	h.ServeHTTP(rw, r)
	end := timeNow()
	dur := end.Sub(start)

	if p.Requests != nil {
		p.Requests.Update(dur)
	}
	if t.Timer != nil {
		t.Timer.Update(dur)
	}
	if rw.code <= 0 {
		return
	}

	metrics.DefaultRegistry.GetTimer(key(rw.code)).Update(dur)

	// write access log
	if p.Logger != nil {
		p.Logger.Log(&logger.Event{
			Start:   start,
			End:     end,
			Request: r,
			Response: &http.Response{
				StatusCode:    rw.code,
				ContentLength: int64(rw.size),
			},
			RequestURL:      requestURL,
			UpstreamAddr:    targetURL.Host,
			UpstreamService: t.Service,
			UpstreamURL:     targetURL,
		})
	}
}

func key(code int) string {
	b := []byte("http.status.")
	b = strconv.AppendInt(b, int64(code), 10)
	return string(b)
}

// responseWriter wraps an http.ResponseWriter to capture the status code and
// the size of the response. It also implements http.Hijacker to forward
// hijacking the connection to the wrapped writer if supported.
type responseWriter struct {
	w    http.ResponseWriter
	code int
	size int
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.w.Write(b)
	rw.size += n
	return n, err
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.w.WriteHeader(statusCode)
	rw.code = statusCode
}

var errNoHijacker = errors.New("not a hijacker")

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rw.w.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errNoHijacker
}
