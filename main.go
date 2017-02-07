package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"github.com/eBay/fabio/mdllog"
	"github.com/eBay/fabio/admin"
	"github.com/eBay/fabio/config"
	"github.com/eBay/fabio/exit"
	"github.com/eBay/fabio/metrics"
	"github.com/eBay/fabio/proxy"
	"github.com/eBay/fabio/registry"
	"github.com/eBay/fabio/registry/consul"
	"github.com/eBay/fabio/registry/file"
	"github.com/eBay/fabio/registry/static"
	"github.com/eBay/fabio/route"
)

// version contains the version number
//
// It is set by build/release.sh for tagged releases
// so that 'go get' just works.
//
// It is also set by the linker when fabio
// is built via the Makefile or the build/docker.sh
// script to ensure the correct version nubmer
var version = "1.3.7"

func main() {
	cfg, err := config.Load(os.Args, os.Environ())
	if err != nil {
		exit.Fatalf("[FATAL] %s. %s", version, err)
	}
	if cfg == nil {
		fmt.Println(version)
		return
	}
	mdllog.Info.Printf("[INFO] Runtime config\n" + toJSON(cfg))
	mdllog.Info.Printf("[INFO] Version %s starting", version)
	mdllog.Info.Printf("[INFO] Go runtime is %s", runtime.Version())

	exit.Listen(func(s os.Signal) {
		if registry.Default == nil {
			return
		}
		registry.Default.Deregister()
	})

	// init metrics early since that create the global metric registries
	// that are used by other parts of the code.
	initMetrics(cfg)

	initRuntime(cfg)
	initBackend(cfg)
	go watchBackend()
	startAdmin(cfg)

	// create proxies after metrics since they use the metrics registry.
	httpProxy := newHTTPProxy(cfg)
	tcpProxy := proxy.NewTCPSNIProxy(cfg.Proxy)
	startListeners(cfg.Listen, cfg.Proxy.ShutdownWait, httpProxy, tcpProxy)
	exit.Wait()
}

func newHTTPProxy(cfg *config.Config) http.Handler {
	if err := route.SetPickerStrategy(cfg.Proxy.Strategy); err != nil {
		exit.Fatal("[FATAL] ", err)
	}
	mdllog.Info.Printf("[INFO] Using routing strategy %q", cfg.Proxy.Strategy)

	if err := route.SetMatcher(cfg.Proxy.Matcher); err != nil {
		exit.Fatal("[FATAL] ", err)
	}
	mdllog.Info.Printf("[INFO] Using routing matching %q", cfg.Proxy.Matcher)

	tr := &http.Transport{
		ResponseHeaderTimeout: cfg.Proxy.ResponseHeaderTimeout,
		MaxIdleConnsPerHost:   cfg.Proxy.MaxConn,
		Dial: (&net.Dialer{
			Timeout:   cfg.Proxy.DialTimeout,
			KeepAlive: cfg.Proxy.KeepAliveTimeout,
		}).Dial,
	}

	return proxy.NewHTTPProxy(tr, cfg.Proxy)
}

func startAdmin(cfg *config.Config) {
	mdllog.Info.Printf("[INFO] Admin server listening on %q", cfg.UI.Addr)
	go func() {
		srv := &admin.Server{
			Color:    cfg.UI.Color,
			Title:    cfg.UI.Title,
			Version:  version,
			Commands: route.Commands,
			Cfg:      cfg,
		}
		if err := srv.ListenAndServe(cfg.UI.Addr); err != nil {
			exit.Fatal("[FATAL] ui: ", err)
		}
	}()
}

func initMetrics(cfg *config.Config) {
	if cfg.Metrics.Target == "" {
		mdllog.Info.Printf("[INFO] Metrics disabled")
		return
	}

	var err error
	if metrics.DefaultRegistry, err = metrics.NewRegistry(cfg.Metrics); err != nil {
		exit.Fatal("[FATAL] ", err)
	}
	if route.ServiceRegistry, err = metrics.NewRegistry(cfg.Metrics); err != nil {
		exit.Fatal("[FATAL] ", err)
	}
}

func initRuntime(cfg *config.Config) {
	if os.Getenv("GOGC") == "" {
		mdllog.Info.Print("[INFO] Setting GOGC=", cfg.Runtime.GOGC)
		debug.SetGCPercent(cfg.Runtime.GOGC)
	} else {
		mdllog.Info.Print("[INFO] Using GOGC=", os.Getenv("GOGC"), " from env")
	}

	if os.Getenv("GOMAXPROCS") == "" {
		mdllog.Info.Print("[INFO] Setting GOMAXPROCS=", cfg.Runtime.GOMAXPROCS)
		runtime.GOMAXPROCS(cfg.Runtime.GOMAXPROCS)
	} else {
		mdllog.Info.Print("[INFO] Using GOMAXPROCS=", os.Getenv("GOMAXPROCS"), " from env")
	}
}

func initBackend(cfg *config.Config) {
	var err error

	switch cfg.Registry.Backend {
	case "file":
		registry.Default, err = file.NewBackend(cfg.Registry.File.Path)
	case "static":
		registry.Default, err = static.NewBackend(cfg.Registry.Static.Routes)
	case "consul":
		registry.Default, err = consul.NewBackend(&cfg.Registry.Consul)
	default:
		exit.Fatal("[FATAL] Unknown registry backend ", cfg.Registry.Backend)
	}

	if err != nil {
		exit.Fatal("[FATAL] Error initializing backend. ", err)
	}
	if err := registry.Default.Register(); err != nil {
		exit.Fatal("[FATAL] Error registering backend. ", err)
	}
}

func watchBackend() {
	var (
		last   string
		svccfg string
		mancfg string
	)

	svc := registry.Default.WatchServices()
	man := registry.Default.WatchManual()

	for {
		select {
		case svccfg = <-svc:
		case mancfg = <-man:
		}

		// manual config overrides service config
		// order matters
		next := svccfg + "\n" + mancfg
		if next == last {
			continue
		}

		t, err := route.NewTable(next)
		if err != nil {
			mdllog.Warning.Printf("[WARN] %s", err)
			continue
		}
		route.SetTable(t)
		mdllog.Info.Printf("[INFO] Updated config to\n%s", t)

		last = next
	}
}

func toJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic("json: " + err.Error())
	}
	return string(data)
}
