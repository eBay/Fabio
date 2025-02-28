package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fabiolb/fabio/bgp"
	"github.com/fabiolb/fabio/transport"
	gkm "github.com/go-kit/kit/metrics"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fabiolb/fabio/admin"
	"github.com/fabiolb/fabio/auth"
	"github.com/fabiolb/fabio/cert"
	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/exit"
	"github.com/fabiolb/fabio/logger"
	"github.com/fabiolb/fabio/metrics"
	"github.com/fabiolb/fabio/noroute"
	"github.com/fabiolb/fabio/proxy"
	"github.com/fabiolb/fabio/proxy/tcp"
	"github.com/fabiolb/fabio/registry"
	"github.com/fabiolb/fabio/registry/consul"
	"github.com/fabiolb/fabio/registry/custom"
	"github.com/fabiolb/fabio/registry/file"
	"github.com/fabiolb/fabio/registry/static"
	"github.com/fabiolb/fabio/route"

	grpc_proxy "github.com/mwitkow/grpc-proxy/proxy"
	"github.com/pkg/profile"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"google.golang.org/grpc"
)

// version contains the version number
//
// It is set by build/release.sh for tagged releases
// so that 'go get' just works.
//
// It is also set by the linker when fabio
// is built via the Makefile or the build/docker.sh
// script to ensure the correct version number
var version = "1.6.5"

var shuttingDown int32

func main() {
	logOutput := logger.NewLevelWriter(os.Stderr, "INFO", "2017/01/01 00:00:00 ")
	log.SetOutput(logOutput)

	cfg, err := config.Load(os.Args, os.Environ())
	if err != nil {
		exit.Fatalf("[FATAL] %s. %s", version, err)
	}
	if cfg == nil {
		fmt.Printf("%s %s\n", version, runtime.Version())
		return
	}

	transport.SetConfig(cfg)

	log.Printf("[INFO] Setting log level to %s", logOutput.Level())
	if !logOutput.SetLevel(cfg.Log.Level) {
		log.Printf("[INFO] Cannot set log level to %s", cfg.Log.Level)
	}

	log.Printf("%s", "[INFO] Runtime config\n" + toJSON(cfg))
	log.Printf("[INFO] Version %s starting", version)
	log.Printf("[INFO] Go runtime is %s", runtime.Version())

	// warn once so that it is at the beginning of the log
	// this will also start the reminder go routine if necessary.
	WarnIfRunAsRoot(cfg.Insecure)

	// setup profiling if enabled
	var prof interface {
		Stop()
	}
	if cfg.ProfileMode != "" {
		var mode func(*profile.Profile)
		switch cfg.ProfileMode {
		case "":
			// do nothing
		case "cpu":
			mode = profile.CPUProfile
		case "mem":
			mode = profile.MemProfile
		case "mutex":
			mode = profile.MutexProfile
		case "block":
			mode = profile.BlockProfile
		case "trace":
			mode = profile.TraceProfile
		default:
			log.Fatalf("[FATAL] Invalid profile mode %q", cfg.ProfileMode)
		}

		prof = profile.Start(mode, profile.ProfilePath(cfg.ProfilePath), profile.NoShutdownHook)
		log.Printf("[INFO] Profile mode %q", cfg.ProfileMode)
		log.Printf("[INFO] Profile path %q", cfg.ProfilePath)
	}

	if cfg.BGP.BGPEnabled {
		err = bgp.ValidateConfig(&cfg.BGP)
		if err != nil {
			log.Fatalf("[FATAL] BGP configuration invalid: %s", err)
		}
	}

	exit.Listen(func(s os.Signal) {
		atomic.StoreInt32(&shuttingDown, 1)
		if registry.Default != nil {
			registry.Default.DeregisterAll()
		}
		time.Sleep(cfg.Proxy.DeregisterGracePeriod)
		proxy.Shutdown(cfg.Proxy.ShutdownWait)
		if prof != nil {
			prof.Stop()
		}
	})

	metrics, err := metrics.Initialize(&cfg.Metrics)
	if err != nil {
		exit.Fatal("[FATAL] ", err)
	}
	route.SetMetricsProvider(metrics)
	initRuntime(cfg)
	initBackend(cfg)

	startAdmin(cfg)

	go watchNoRouteHTML(cfg)

	first := make(chan bool)
	go watchBackend(cfg, metrics, first)
	log.Print("[INFO] Waiting for first routing table")
	<-first

	// create proxies after metrics since they use the metrics registry.
	startServers(cfg, metrics)

	// warn again so that it is visible in the terminal
	WarnIfRunAsRoot(cfg.Insecure)
	if cfg.BGP.BGPEnabled {
		startBGP(&cfg.BGP)
	}
	exit.Wait()
	log.Print("[INFO] Down")
}

func newGrpcProxy(cfg *config.Config, tlscfg *tls.Config, statsHandler *proxy.GrpcStatsHandler) []grpc.ServerOption {

	//Init Glob Cache
	globCache := route.NewGlobCache(cfg.GlobCacheSize)

	proxyInterceptor := proxy.GrpcProxyInterceptor{
		Config:       cfg,
		StatsHandler: statsHandler,
		GlobCache:    globCache,
	}

	handler := grpc_proxy.TransparentHandler(proxy.GetGRPCDirector(tlscfg, cfg))

	return []grpc.ServerOption{
		grpc.CustomCodec(grpc_proxy.Codec()),
		grpc.UnknownServiceHandler(handler),
		grpc.StreamInterceptor(proxyInterceptor.Stream),
		grpc.StatsHandler(statsHandler),
		grpc.MaxRecvMsgSize(cfg.Proxy.GRPCMaxRxMsgSize),
		grpc.MaxSendMsgSize(cfg.Proxy.GRPCMaxTxMsgSize),
	}
}

func newHTTPProxy(cfg *config.Config, statsHandler *proxy.HttpStatsHandler) *proxy.HTTPProxy {
	var w io.Writer

	//Init Glob Cache
	globCache := route.NewGlobCache(cfg.GlobCacheSize)

	switch cfg.Log.AccessTarget {
	case "":
		log.Printf("[INFO] Access logging disabled")
	case "stdout":
		log.Printf("[INFO] Writing access log to stdout")
		w = os.Stdout
	default:
		exit.Fatal("[FATAL] Invalid access log target ", cfg.Log.AccessTarget)
	}

	format := cfg.Log.AccessFormat
	switch format {
	case "common":
		format = logger.CommonFormat
	case "combined":
		format = logger.CombinedFormat
	}

	l, err := logger.New(w, format)
	if err != nil {
		exit.Fatal("[FATAL] Invalid log format: ", err)
	}

	pick := route.Picker[cfg.Proxy.Strategy]
	match := route.Matcher[cfg.Proxy.Matcher]
	log.Printf("[INFO] Using routing strategy %q", cfg.Proxy.Strategy)
	log.Printf("[INFO] Using route matching %q", cfg.Proxy.Matcher)

	authSchemes, err := auth.LoadAuthSchemes(cfg.Proxy.AuthSchemes)

	if err != nil {
		exit.Fatal("[FATAL] ", err)
	}

	return &proxy.HTTPProxy{
		Config:            cfg.Proxy,
		Transport:         transport.NewTransport(nil),
		InsecureTransport: transport.NewTransport(&tls.Config{InsecureSkipVerify: true}),
		Lookup: func(r *http.Request) *route.Target {
			t := route.GetTable().Lookup(r, pick, match, globCache, cfg.GlobMatchingDisabled)
			if t == nil {
				statsHandler.Noroute.Add(1)
				log.Print("[WARN] No route for ", r.Host, r.URL)
			}
			return t
		},
		Logger:      l,
		AuthSchemes: authSchemes,
		Stats:       *statsHandler,
	}
}

func lookupHostFn(cfg *config.Config, notFound gkm.Counter) func(string) *route.Target {
	pick := route.Picker[cfg.Proxy.Strategy]
	return func(host string) *route.Target {
		t := route.GetTable().LookupHost(host, pick)
		if t == nil {
			notFound.Add(1)
			log.Print("[WARN] No route for ", host)
		}
		return t
	}
}

// Returns a matcher function compatible with tcpproxy Matcher from github.com/inetaf/tcpproxy
func lookupHostMatcher(cfg *config.Config) func(context.Context, string) bool {
	pick := route.Picker[cfg.Proxy.Strategy]
	return func(ctx context.Context, host string) bool {
		t := route.GetTable().LookupHost(host, pick)
		if t == nil {
			return false
		}

		// Make sure this is supposed to be a tcp proxy.
		// opts proto= overrides scheme if present.
		var (
			ok    bool
			proto string
		)
		if proto, ok = t.Opts["proto"]; !ok && t.URL != nil {
			proto = t.URL.Scheme
		}
		return "tcp" == proto
	}
}

func makeTLSConfig(l config.Listen) (*tls.Config, error) {
	if l.CertSource.Name == "" {
		return nil, nil
	}
	src, err := cert.NewSource(l.CertSource)
	if err != nil {
		return nil, fmt.Errorf("Failed to create cert source %s. %s", l.CertSource.Name, err)
	}
	tlscfg, err := cert.TLSConfig(src, l.StrictMatch, l.TLSMinVersion, l.TLSMaxVersion, l.TLSCiphers)
	if err != nil {
		return nil, fmt.Errorf("[FATAL] Failed to create TLS config for cert source %s. %s", l.CertSource.Name, err)
	}
	return tlscfg, nil
}

func startAdmin(cfg *config.Config) {
	log.Printf("[INFO] Admin server access mode %q", cfg.UI.Access)
	log.Printf("[INFO] Admin server listening on %q", cfg.UI.Listen.Addr)
	go func() {
		l := cfg.UI.Listen
		tlscfg, err := makeTLSConfig(l)
		if err != nil {
			exit.Fatal("[FATAL] ", err)
		}
		srv := &admin.Server{
			Access:   cfg.UI.Access,
			Color:    cfg.UI.Color,
			Title:    cfg.UI.Title,
			Version:  version,
			Commands: route.Commands,
			Cfg:      cfg,
		}
		if err := srv.ListenAndServe(l, tlscfg); err != nil {
			exit.Fatal("[FATAL] ui: ", err)
		}
	}()
}

func startServers(cfg *config.Config, stats metrics.Provider) {
	notFound := stats.NewCounter("notfound")

	var (
		tcpConn          gkm.Counter
		tcpConnFail      gkm.Counter
		tcpNoRoute       gkm.Counter
		tcpSniConn       gkm.Counter
		tcpSniConnFail   gkm.Counter
		tcpSniNoRoute    gkm.Counter
		grpStatsHandler  *proxy.GrpcStatsHandler
		httpStatsHandler *proxy.HttpStatsHandler
	)

	grpcCounters := func() {
		grpStatsHandler = &proxy.GrpcStatsHandler{
			Connect: stats.NewCounter("grpc.conn"),
			Request: stats.NewHistogram("grpc.requests"),
			NoRoute: stats.NewCounter("grpc.noroute"),
			Status:  stats.NewHistogram("grep.status", "code"),
		}
	}

	var grpcOnce sync.Once

	httpCounters := func() {
		httpStatsHandler = &proxy.HttpStatsHandler{
			Requests:        stats.NewHistogram("requests"),
			Noroute:         notFound,
			WSConn:          stats.NewGauge("ws.conn"),
			StatusTimer:     stats.NewHistogram("http.status", "code"),
			RedirectCounter: stats.NewCounter("http.redirect.count", "code"),
		}
	}

	var httpOnce sync.Once

	tcpCounters := func() {
		tcpConn = stats.NewCounter("tcp.conn")
		tcpConnFail = stats.NewCounter("tcp.connfail")
		tcpNoRoute = stats.NewCounter("tcp.noroute")
	}
	var tcpOnce sync.Once

	tcpSniCounters := func() {
		tcpSniConn = stats.NewCounter("tcp_sni.conn")
		tcpSniConnFail = stats.NewCounter("tcp_sni.connfail")
		tcpSniNoRoute = stats.NewCounter("tcp_sni.noroute")
	}

	var tcpSniOnce sync.Once

	for _, l := range cfg.Listen {
		l := l // capture loop var for go routines below
		tlscfg, err := makeTLSConfig(l)
		if err != nil {
			exit.Fatal("[FATAL] ", err)
		}

		log.Printf("[INFO] %s proxy listening on %s", strings.ToUpper(l.Proto), l.Addr)
		if tlscfg != nil && tlscfg.ClientAuth == tls.RequireAndVerifyClientCert {
			log.Printf("[INFO] Client certificate authentication enabled on %s", l.Addr)
		}
		switch l.Proto {
		case "http", "https":
			httpOnce.Do(httpCounters)
			go func() {
				h := newHTTPProxy(cfg, httpStatsHandler)
				// reset the ws.conn gauge
				h.Stats.WSConn.Set(0)
				if err := proxy.ListenAndServeHTTP(l, h, tlscfg); err != nil {
					exit.Fatal("[FATAL] ", err)
				}
			}()
		case "grpc", "grpcs":
			grpcOnce.Do(grpcCounters)
			go func() {
				h := newGrpcProxy(cfg, tlscfg, grpStatsHandler)
				if err := proxy.ListenAndServeGRPC(l, h, tlscfg); err != nil {
					exit.Fatal("[FATAL] ", err)
				}
			}()
		case "tcp":
			tcpOnce.Do(tcpCounters)
			go func() {
				h := &tcp.Proxy{
					DialTimeout: cfg.Proxy.DialTimeout,
					Lookup:      lookupHostFn(cfg, notFound),
					Conn:        tcpConn,
					ConnFail:    tcpConnFail,
					Noroute:     tcpNoRoute,
				}
				if err := proxy.ListenAndServeTCP(l, h, tlscfg); err != nil {
					exit.Fatal("[FATAL] ", err)
				}
			}()
		case "tcp+sni":
			tcpSniOnce.Do(tcpSniCounters)
			go func() {
				h := &tcp.SNIProxy{
					DialTimeout: cfg.Proxy.DialTimeout,
					Lookup:      lookupHostFn(cfg, notFound),
					Conn:        tcpSniConn,
					ConnFail:    tcpSniConnFail,
					Noroute:     tcpSniNoRoute,
				}
				if err := proxy.ListenAndServeTCP(l, h, tlscfg); err != nil {
					exit.Fatal("[FATAL] ", err)
				}
			}()
		case "tcp-dynamic":
			tcpOnce.Do(tcpCounters)
			go func() {
				var buffer strings.Builder
				lastPorts := []string{}
				for {
					time.Sleep(l.Refresh)
					table := route.GetTable()
					ports := []string{}
					for target, rts := range table {
						if strings.Contains(target, ":") {
							buffer.WriteString(":")
							buffer.WriteString(strings.Split(target, ":")[1])

							schemes := tableSchemes(rts)
							if len(schemes) == 1 && schemes[0] == "tcp" {
								ports = append(ports, buffer.String())
							}
							buffer.Reset()
						}
						ports = unique(ports)
					}
					for _, port := range difference(lastPorts, ports) {
						log.Printf("[DEBUG] Dynamic TCP listener on %s eligable for termination", port)
						proxy.CloseProxy(port)
					}
					for _, port := range ports {
						l := l
						port := port
						conn, err := net.Listen("tcp", port)
						if err != nil {
							log.Printf("[DEBUG] Dynamic TCP port %s in use", port)
							continue
						}
						conn.Close()
						log.Printf("[INFO] Starting dynamic TCP listener on port %s ", port)
						go func() {
							h := &tcp.DynamicProxy{
								DialTimeout: cfg.Proxy.DialTimeout,
								Lookup:      lookupHostFn(cfg, notFound),
								Conn:        tcpConn,
								ConnFail:    tcpConnFail,
								Noroute:     tcpNoRoute,
							}
							l.Addr = port
							if err := proxy.ListenAndServeTCP(l, h, tlscfg); err != nil {
								exit.Fatal("[FATAL] ", err)
							}
						}()
					}
					lastPorts = ports
				}
			}()
		case "https+tcp+sni":
			tcpSniOnce.Do(tcpSniCounters)
			httpOnce.Do(httpCounters)
			go func() {
				hp := newHTTPProxy(cfg, httpStatsHandler)
				tp := &tcp.SNIProxy{
					DialTimeout: cfg.Proxy.DialTimeout,
					Lookup:      lookupHostFn(cfg, notFound),
					Conn:        tcpSniConn,
					ConnFail:    tcpSniConnFail,
					Noroute:     tcpSniNoRoute}
				if err := proxy.ListenAndServeHTTPSTCPSNI(l, hp, tp, tlscfg, lookupHostMatcher(cfg)); err != nil {
					exit.Fatal("[FATAL] ", err)
				}
			}()
		case "prometheus":
			go func() {
				if err := proxy.ListenAndServePrometheus(l, cfg.Metrics.Prometheus, tlscfg); err != nil {
					exit.Fatal("[FATAL] ", err)
				}
			}()
		default:
			exit.Fatal("[FATAL] Invalid protocol ", l.Proto)
		}
	}
}

func startBGP(cfg *config.BGP) {
	h, err := bgp.NewBGPHandler(cfg)
	if err != nil {
		exit.Fatal("[FATAL] ", err)
	}
	go func() {
		if err := h.Start(); err != nil {
			exit.Fatal("[FATAL] ", err)
		}
	}()
}

func initRuntime(cfg *config.Config) {
	if os.Getenv("GOGC") == "" {
		log.Print("[INFO] Setting GOGC=", cfg.Runtime.GOGC)
		debug.SetGCPercent(cfg.Runtime.GOGC)
	} else {
		log.Print("[INFO] Using GOGC=", os.Getenv("GOGC"), " from env")
	}

	if os.Getenv("GOMAXPROCS") == "" {
		log.Print("[INFO] Setting GOMAXPROCS=", cfg.Runtime.GOMAXPROCS)
		runtime.GOMAXPROCS(cfg.Runtime.GOMAXPROCS)
	} else {
		log.Print("[INFO] Using GOMAXPROCS=", os.Getenv("GOMAXPROCS"), " from env")
	}
}

func initBackend(cfg *config.Config) {
	var deadline = time.Now().Add(cfg.Registry.Timeout)
	var err error
	for {
		switch cfg.Registry.Backend {
		case "file":
			registry.Default, err = file.NewBackend(&cfg.Registry.File)
		case "static":
			registry.Default, err = static.NewBackend(&cfg.Registry.Static)
		case "consul":
			registry.Default, err = consul.NewBackend(&cfg.Registry.Consul)
		case "custom":
			registry.Default, err = custom.NewBackend(&cfg.Registry.Custom)
		default:
			exit.Fatal("[FATAL] Unknown registry backend ", cfg.Registry.Backend)
		}

		if err == nil {
			if err = registry.Default.Register(nil); err == nil {
				return
			}
		}
		log.Print("[WARN] Error initializing backend. ", err)

		if time.Now().After(deadline) {
			exit.Fatal("[FATAL] Timeout registering backend.")
		}

		time.Sleep(cfg.Registry.Retry)
		if atomic.LoadInt32(&shuttingDown) > 0 {
			exit.Exit(1)
		}
	}
}

func watchBackend(cfg *config.Config, p metrics.Provider, first chan bool) {
	var (
		nextTable   string
		lastTable   string
		svccfg      string
		mancfg      string
		customBE    string
		once        sync.Once
		tableBuffer = new(bytes.Buffer) // fix crash on reset before used (#650)
	)

	switch cfg.Registry.Backend {
	// custom back end receives JSON from a remote source that contains a slice of route.RouteDef
	// the route table is created directly from that input
	case "custom":
		svc := registry.Default.WatchServices()
		for {
			customBE = <-svc
			if customBE != "OK" {
				log.Printf("[ERROR] error during update from custom back end - %s", customBE)
			}
			once.Do(func() { close(first) })
		}
	// all other backend types
	default:
		svc := registry.Default.WatchServices()
		man := registry.Default.WatchManual()

		for {
			select {
			case svccfg = <-svc:
			case mancfg = <-man:
			}
			// manual config overrides service config - order matters
			tableBuffer.Reset()
			tableBuffer.WriteString(svccfg)
			tableBuffer.WriteString("\n")
			tableBuffer.WriteString(mancfg)
			// set nextTable here to preserve the state.  The buffer is altered
			// when calling route.NewTable and we lose change logging (#737)
			if nextTable = tableBuffer.String(); nextTable == lastTable {
				continue
			}
			aliases, err := route.ParseAliases(nextTable)
			if err != nil {
				log.Printf("[WARN]: %s", err)
			}
			registry.Default.Register(aliases)
			t, err := route.NewTable(tableBuffer)
			if err != nil {
				log.Printf("[WARN] %s", err)
				continue
			}
			route.SetTable(t)
			logRoutes(t, lastTable, nextTable, cfg.Log.RoutesFormat)
			lastTable = nextTable
			once.Do(func() { close(first) })
		}
	}
}

func watchNoRouteHTML(cfg *config.Config) {
	html := registry.Default.WatchNoRouteHTML()
	for {
		next := <-html
		if next == noroute.GetHTML() {
			continue
		}
		noroute.SetHTML(next)
		if next == "" {
			log.Print("[INFO] Unset noroute HTML")
		} else {
			log.Printf("[INFO] Set noroute HTML (%d bytes)", len(next))
		}
	}
}

func logRoutes(t route.Table, last, next, format string) {
	fmtDiff := func(diffs []dmp.Diff) string {
		var b bytes.Buffer
		for _, d := range diffs {
			t := strings.TrimSpace(d.Text)
			if t == "" {
				continue
			}
			switch d.Type {
			case dmp.DiffDelete:
				b.WriteString("- ")
				b.WriteString(strings.Replace(t, "\n", "\n- ", -1))
			case dmp.DiffInsert:
				b.WriteString("+ ")
				b.WriteString(strings.Replace(t, "\n", "\n+ ", -1))
			}
		}
		return b.String()
	}

	const defFormat = "delta"
	switch format {
	case "detail":
		log.Printf("[INFO] Updated config to\n%s", t.Dump())

	case "delta":
		if delta := fmtDiff(dmp.New().DiffMain(last, next, true)); delta != "" {
			log.Printf("[INFO] Config updates\n%s", delta)
		}

	case "all":
		log.Printf("[INFO] Updated config to\n%s", next)

	default:
		log.Printf("[WARN] Invalid route format %q. Defaulting to %q", format, defFormat)
		logRoutes(t, last, next, defFormat)
	}
}

func toJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic("json: " + err.Error())
	}
	return string(data)
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// difference returns elements in `a` that aren't in `b`
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func tableSchemes(r route.Routes) []string {
	schemes := []string{}
	for _, rt := range r {
		for _, target := range rt.Targets {
			schemes = append(schemes, target.URL.Scheme)
		}
	}
	return unique(schemes)
}
