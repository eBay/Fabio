package proxy

import (
	"bufio"
	"crypto/tls"
	"errors"
	gkm "github.com/go-kit/kit/metrics"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fabiolb/fabio/auth"
	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/logger"
	"github.com/fabiolb/fabio/noroute"
	"github.com/fabiolb/fabio/proxy/gzip"
	"github.com/fabiolb/fabio/route"
	"github.com/fabiolb/fabio/trace"
	"github.com/fabiolb/fabio/uuid"
)

type HttpStatsHandler struct {
	// Requests is a histogram metric which is updated for every request.
	Requests gkm.Histogram

	// Noroute is a counter metric which is updated for every request
	// where Lookup() returns nil.
	Noroute gkm.Counter

	// WSConn counts the number of open web socket connections.
	WSConn gkm.Gauge

	// StatusTimer is a histogram for given status codes
	StatusTimer gkm.Histogram

	// RedirectCounter - counts redirects
	RedirectCounter gkm.Counter
}

// HTTPProxy is a dynamic reverse proxy for HTTP and HTTPS protocols.
type HTTPProxy struct {
	// Config is the proxy configuration as provided during startup.
	Config config.Proxy

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

	// Logger is the access logger for the requests.
	Logger logger.Logger

	// TracerCfg is the Open Tracing configuration as provided during startup
	TracerCfg config.Tracing

	// UUID returns a unique id in uuid format.
	// If UUID is nil, uuid.NewUUID() is used.
	UUID func() string

	// Auth schemes registered with the server
	AuthSchemes map[string]auth.AuthScheme

	// stats contains all of the stats bits
	Stats HttpStatsHandler
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Lookup == nil {
		panic("no lookup function")
	}

	if p.Config.RequestID != "" {
		id := p.UUID
		if id == nil {
			id = uuid.NewUUID
		}
		r.Header.Set(p.Config.RequestID, id())
	}

	//Create Span
	span := trace.CreateSpan(r, &p.TracerCfg)
	defer span.Finish()

	t := p.Lookup(r)

	if t == nil {
		status := p.Config.NoRouteStatus
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

	authDecision := t.Authorized(r, w, p.AuthSchemes)
	if authDecision.Done {
		return
	}
	if !authDecision.Authorized {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
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
		if p.Stats.RedirectCounter != nil {
			p.Stats.RedirectCounter.With("code", strconv.Itoa(t.RedirectCode)).Add(1)
		}
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
		// ensure absolute path after stripping to maintain compliance with
		// section 5.3 of RFC7230 (https://tools.ietf.org/html/rfc7230#section-5.3)
		if !strings.HasPrefix(targetURL.Path, "/") {
			targetURL.Path = "/" + targetURL.Path
		}
	}

	if t.PrependPath != "" {
		targetURL.Path = t.PrependPath + targetURL.Path
		// ensure absolute path after stripping to maintain compliance with
		// section 5.3 of RFC7230 (https://tools.ietf.org/html/rfc7230#section-5.3)
		if !strings.HasPrefix(targetURL.Path, "/") {
			targetURL.Path = "/" + targetURL.Path
		}
	}

	if err := addHeaders(r, p.Config, t.StripPath); err != nil {
		http.Error(w, "cannot parse "+r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	if err := addResponseHeaders(w, r, p.Config); err != nil {
		http.Error(w, "cannot add response headers", http.StatusInternalServerError)
		return
	}

	//Add OpenTrace Headers to response
	trace.InjectHeaders(span, r)

	upgrade, accept := r.Header.Get("Upgrade"), r.Header.Get("Accept")

	tr := p.Transport
	if t.Transport != nil {
		tr = t.Transport
	} else if t.TLSSkipVerify {
		tr = p.InsecureTransport
	}

	var h http.Handler
	switch {
	case upgrade == "websocket" || upgrade == "Websocket":
		r.URL = targetURL
		if targetURL.Scheme == "https" || targetURL.Scheme == "wss" {
			h = newWSHandler(targetURL.Host, func(network, address string) (net.Conn, error) {
				return tls.Dial(network, address, tr.(*http.Transport).TLSClientConfig)
			}, p.Stats.WSConn)
		} else {
			h = newWSHandler(targetURL.Host, net.Dial, p.Stats.WSConn)
		}

	case accept == "text/event-stream":
		// use the flush interval for SSE (server-sent events)
		// must be > 0s to be effective
		h = newHTTPProxy(targetURL, tr, p.Config.FlushInterval)

	default:
		h = newHTTPProxy(targetURL, tr, p.Config.GlobalFlushInterval)
	}

	if p.Config.GZIPContentTypes != nil {
		h = gzip.NewGzipHandler(h, p.Config.GZIPContentTypes)
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

	if p.Stats.Requests != nil {
		p.Stats.Requests.Observe(dur.Seconds())
	}
	if t.Timer != nil {
		t.Timer.Observe(dur.Seconds())
	}
	if rw.code <= 0 {
		return
	}

	if p.Stats.StatusTimer != nil {
		p.Stats.StatusTimer.With("code", strconv.Itoa(rw.code)).Observe(dur.Seconds())
	}

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

func (rw *responseWriter) Flush() {
	if fl, ok := rw.w.(http.Flusher); ok {
		fl.Flush()
	}
}

var errNoHijacker = errors.New("not a hijacker")

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rw.w.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errNoHijacker
}
