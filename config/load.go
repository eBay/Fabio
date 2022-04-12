package config

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	gs "github.com/hashicorp/go-sockaddr/template"
	"github.com/magiconair/properties"
)

func Load(args, environ []string) (cfg *Config, err error) {
	var props *properties.Properties

	cmdline, path, version, err := parse(args)
	switch {
	case err != nil:
		return nil, err
	case version:
		return nil, nil
	case path != "":
		switch {
		case strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"):
			props, err = properties.LoadURL(path)
		case path != "":
			props, err = properties.LoadFile(path, properties.UTF8)
		}
		if err != nil {
			return nil, err
		}
	}
	envprefix := []string{"FABIO_", ""}
	return load(cmdline, environ, envprefix, props)
}

var errInvalidConfig = errors.New("invalid or missing path to config file")

// parse extracts the version and config file flags from the command
// line arguments and returns the individual parts. Test flags are
// ignored.
func parse(args []string) (cmdline []string, path string, version bool, err error) {
	if len(args) < 1 {
		panic("missing exec name")
	}

	// always copy the name of the executable
	cmdline = args[:1]

	// parse rest of the arguments
	for i := 1; i < len(args); i++ {
		arg := args[i]

		switch {
		// version flag
		case arg == "-v" || arg == "-version" || arg == "--version":
			return nil, "", true, nil

		// config file without '='
		case arg == "-cfg" || arg == "--cfg":
			if i >= len(args)-1 {
				return nil, "", false, errInvalidConfig
			}
			path = args[i+1]
			i++

		// config file with '='. needs unquoting
		case strings.HasPrefix(arg, "-cfg=") || strings.HasPrefix(arg, "--cfg="):
			if strings.HasPrefix(arg, "-cfg=") {
				path = arg[len("-cfg="):]
			} else {
				path = arg[len("--cfg="):]
			}
			switch {
			case path == "":
				return nil, "", false, errInvalidConfig
			case path[0] == '\'':
				path = strings.Trim(path, "'")
			case path[0] == '"':
				path = strings.Trim(path, "\"")
			}
			if path == "" {
				return nil, "", false, errInvalidConfig
			}

		// ignore test flags
		case strings.HasPrefix(arg, "-test."):
			continue

		default:
			cmdline = append(cmdline, arg)
		}
	}
	return cmdline, path, false, nil
}

func load(cmdline, environ, envprefix []string, props *properties.Properties) (cfg *Config, err error) {
	cfg = &Config{}
	f := NewFlagSet(cmdline[0], flag.ExitOnError)

	// dummy values which were parsed earlier
	f.String("cfg", "", "Path or URL to config file")
	f.Bool("v", false, "Show version")
	f.Bool("version", false, "Show version")

	// config values
	var listenerValue string
	var uiListenerValue string
	var certSourcesValue string
	var authSchemesValue string
	var readTimeout, writeTimeout time.Duration
	var gzipContentTypesValue string

	var obsoleteStr string

	f.BoolVar(&cfg.Insecure, "insecure", defaultConfig.Insecure, "allow fabio to run as root when set to true")
	f.IntVar(&cfg.Proxy.MaxConn, "proxy.maxconn", defaultConfig.Proxy.MaxConn, "maximum number of cached connections")
	f.StringVar(&cfg.Proxy.Strategy, "proxy.strategy", defaultConfig.Proxy.Strategy, "load balancing strategy")
	f.StringVar(&cfg.Proxy.Matcher, "proxy.matcher", defaultConfig.Proxy.Matcher, "path matching algorithm")
	f.IntVar(&cfg.Proxy.NoRouteStatus, "proxy.noroutestatus", defaultConfig.Proxy.NoRouteStatus, "status code for invalid route. Must be three digits")
	f.DurationVar(&cfg.Proxy.ShutdownWait, "proxy.shutdownwait", defaultConfig.Proxy.ShutdownWait, "time for graceful shutdown")
	f.DurationVar(&cfg.Proxy.DialTimeout, "proxy.dialtimeout", defaultConfig.Proxy.DialTimeout, "connection timeout for backend connections")
	f.DurationVar(&cfg.Proxy.ResponseHeaderTimeout, "proxy.responseheadertimeout", defaultConfig.Proxy.ResponseHeaderTimeout, "response header timeout")
	f.DurationVar(&cfg.Proxy.KeepAliveTimeout, "proxy.keepalivetimeout", defaultConfig.Proxy.KeepAliveTimeout, "keep-alive timeout")
	f.DurationVar(&cfg.Proxy.IdleConnTimeout, "proxy.idleconntimeout", defaultConfig.Proxy.IdleConnTimeout, "idle timeout, when to close (keep-alive) connections")
	f.StringVar(&cfg.Proxy.LocalIP, "proxy.localip", defaultConfig.Proxy.LocalIP, "fabio address in Forward headers")
	f.StringVar(&cfg.Proxy.ClientIPHeader, "proxy.header.clientip", defaultConfig.Proxy.ClientIPHeader, "header for the request ip")
	f.StringVar(&cfg.Proxy.TLSHeader, "proxy.header.tls", defaultConfig.Proxy.TLSHeader, "header for TLS connections")
	f.StringVar(&cfg.Proxy.TLSHeaderValue, "proxy.header.tls.value", defaultConfig.Proxy.TLSHeaderValue, "value for TLS connection header")
	f.StringVar(&cfg.Proxy.RequestID, "proxy.header.requestid", defaultConfig.Proxy.RequestID, "header for reqest id")
	f.IntVar(&cfg.Proxy.STSHeader.MaxAge, "proxy.header.sts.maxage", defaultConfig.Proxy.STSHeader.MaxAge, "enable and set the max-age value for HSTS")
	f.BoolVar(&cfg.Proxy.STSHeader.Subdomains, "proxy.header.sts.subdomains", defaultConfig.Proxy.STSHeader.Subdomains, "direct HSTS to include subdomains")
	f.BoolVar(&cfg.Proxy.STSHeader.Preload, "proxy.header.sts.preload", defaultConfig.Proxy.STSHeader.Preload, "direct HSTS to pass the preload directive")
	f.StringVar(&gzipContentTypesValue, "proxy.gzip.contenttype", defaultValues.GZIPContentTypesValue, "regexp of content types to compress")
	f.StringVar(&listenerValue, "proxy.addr", defaultValues.ListenerValue, "listener config")
	f.StringVar(&certSourcesValue, "proxy.cs", defaultValues.CertSourcesValue, "certificate sources")
	f.DurationVar(&readTimeout, "proxy.readtimeout", defaultValues.ReadTimeout, "read timeout for incoming requests")
	f.DurationVar(&writeTimeout, "proxy.writetimeout", defaultValues.WriteTimeout, "write timeout for outgoing responses")
	f.DurationVar(&cfg.Proxy.FlushInterval, "proxy.flushinterval", defaultConfig.Proxy.FlushInterval, "flush interval for streaming responses")
	f.DurationVar(&cfg.Proxy.GlobalFlushInterval, "proxy.globalflushinterval", defaultConfig.Proxy.GlobalFlushInterval, "flush interval for non-streaming responses")
	f.StringVar(&authSchemesValue, "proxy.auth", defaultValues.AuthSchemesValue, "auth schemes")
	f.StringVar(&cfg.Log.AccessFormat, "log.access.format", defaultConfig.Log.AccessFormat, "access log format")
	f.StringVar(&cfg.Log.AccessTarget, "log.access.target", defaultConfig.Log.AccessTarget, "access log target")
	f.StringVar(&cfg.Log.RoutesFormat, "log.routes.format", defaultConfig.Log.RoutesFormat, "log format of routing table updates")
	f.StringVar(&cfg.Log.Level, "log.level", defaultConfig.Log.Level, "log level: TRACE, DEBUG, INFO, WARN, ERROR, FATAL")
	f.StringVar(&cfg.Metrics.Target, "metrics.target", defaultConfig.Metrics.Target, "metrics backend")
	f.StringVar(&cfg.Metrics.Prefix, "metrics.prefix", defaultConfig.Metrics.Prefix, "prefix for reported metrics")
	f.StringVar(&cfg.Metrics.Names, "metrics.names", defaultConfig.Metrics.Names, "route metric name template")
	f.DurationVar(&cfg.Metrics.Interval, "metrics.interval", defaultConfig.Metrics.Interval, "metrics reporting interval")
	f.DurationVar(&cfg.Metrics.Timeout, "metrics.timeout", defaultConfig.Metrics.Timeout, "timeout for metrics to become available")
	f.DurationVar(&cfg.Metrics.Retry, "metrics.retry", defaultConfig.Metrics.Retry, "retry interval during startup")
	f.StringVar(&cfg.Metrics.GraphiteAddr, "metrics.graphite.addr", defaultConfig.Metrics.GraphiteAddr, "graphite server address")
	f.StringVar(&cfg.Metrics.StatsDAddr, "metrics.statsd.addr", defaultConfig.Metrics.StatsDAddr, "statsd server address")
	f.StringVar(&cfg.Metrics.DogstatsdAddr, "metrics.dogstatsd.addr", defaultConfig.Metrics.DogstatsdAddr, "dogstatsd server address")
	f.StringVar(&cfg.Metrics.Circonus.APIKey, "metrics.circonus.apikey", defaultConfig.Metrics.Circonus.APIKey, "Circonus API token key")
	f.StringVar(&cfg.Metrics.Circonus.APIApp, "metrics.circonus.apiapp", defaultConfig.Metrics.Circonus.APIApp, "Circonus API token app")
	f.StringVar(&cfg.Metrics.Circonus.APIURL, "metrics.circonus.apiurl", defaultConfig.Metrics.Circonus.APIURL, "Circonus API URL")
	f.StringVar(&cfg.Metrics.Circonus.BrokerID, "metrics.circonus.brokerid", defaultConfig.Metrics.Circonus.BrokerID, "Circonus Broker ID")
	f.StringVar(&cfg.Metrics.Circonus.CheckID, "metrics.circonus.checkid", defaultConfig.Metrics.Circonus.CheckID, "Circonus Check ID")
	f.StringVar(&cfg.Metrics.Circonus.SubmissionURL, "metrics.circonus.submissionurl", defaultConfig.Metrics.Circonus.SubmissionURL, "Circonus Check SubmissionURL")
	f.StringVar(&cfg.Metrics.Prometheus.Subsystem, "metrics.prometheus.subsystem", defaultConfig.Metrics.Prometheus.Subsystem, "Prometheus system")
	f.StringVar(&cfg.Metrics.Prometheus.Path, "metrics.prometheus.path", defaultConfig.Metrics.Prometheus.Path, "Prometheus http handler path")
	f.FloatSliceVar(&cfg.Metrics.Prometheus.Buckets, "metrics.prometheus.buckets", defaultConfig.Metrics.Prometheus.Buckets, "Prometheus histogram buckets")
	f.StringVar(&cfg.Registry.Backend, "registry.backend", defaultConfig.Registry.Backend, "registry backend")
	f.DurationVar(&cfg.Registry.Timeout, "registry.timeout", defaultConfig.Registry.Timeout, "timeout for registry to become available")
	f.DurationVar(&cfg.Registry.Retry, "registry.retry", defaultConfig.Registry.Retry, "retry interval during startup")
	f.StringVar(&cfg.Registry.File.RoutesPath, "registry.file.path", defaultConfig.Registry.File.RoutesPath, "path to file based routing table")
	f.StringVar(&cfg.Registry.File.NoRouteHTMLPath, "registry.file.noroutehtmlpath", defaultConfig.Registry.File.NoRouteHTMLPath, "path to file for HTML returned when no route is found")
	f.StringVar(&cfg.Registry.Static.Routes, "registry.static.routes", defaultConfig.Registry.Static.Routes, "static routes")
	f.StringVar(&cfg.Registry.Static.NoRouteHTML, "registry.static.noroutehtml", defaultConfig.Registry.Static.NoRouteHTML, "HTML which is returned when no route is found")
	f.StringVar(&cfg.Registry.Consul.Addr, "registry.consul.addr", defaultConfig.Registry.Consul.Addr, "address of the consul agent")
	f.StringVar(&cfg.Registry.Consul.Token, "registry.consul.token", defaultConfig.Registry.Consul.Token, "token for consul agent")
	f.StringVar(&cfg.Registry.Consul.KVPath, "registry.consul.kvpath", defaultConfig.Registry.Consul.KVPath, "consul KV path for manual overrides")
	f.StringVar(&cfg.Registry.Consul.NoRouteHTMLPath, "registry.consul.noroutehtmlpath", defaultConfig.Registry.Consul.NoRouteHTMLPath, "consul KV path for HTML returned when no route is found")
	f.StringVar(&cfg.Registry.Consul.TagPrefix, "registry.consul.tagprefix", defaultConfig.Registry.Consul.TagPrefix, "prefix for consul tags")
	f.StringVar(&cfg.Registry.Consul.TLS.KeyFile, "registry.consul.tls.keyfile", defaultConfig.Registry.Consul.TLS.KeyFile, "path to consul key file")
	f.StringVar(&cfg.Registry.Consul.TLS.CertFile, "registry.consul.tls.certfile", defaultConfig.Registry.Consul.TLS.CertFile, "path to consul cert file")
	f.StringVar(&cfg.Registry.Consul.TLS.CAFile, "registry.consul.tls.cafile", defaultConfig.Registry.Consul.TLS.CAFile, "path to consul CA file")
	f.StringVar(&cfg.Registry.Consul.TLS.CAPath, "registry.consul.tls.capath", defaultConfig.Registry.Consul.TLS.CAPath, "path to consul CA directory")
	f.BoolVar(&cfg.Registry.Consul.TLS.InsecureSkipVerify, "registry.consul.tls.insecureskipverify", defaultConfig.Registry.Consul.TLS.InsecureSkipVerify, "is tls check enabled")
	f.BoolVar(&cfg.Registry.Consul.Register, "registry.consul.register.enabled", defaultConfig.Registry.Consul.Register, "register fabio in consul")
	f.StringVar(&cfg.Registry.Consul.Namespace, "registry.consul.namespace", defaultConfig.Registry.Consul.Namespace, "consul namespace in which fabio is active")
	f.StringVar(&cfg.Registry.Consul.ServiceAddr, "registry.consul.register.addr", "<ui.addr>", "service registration address")
	f.StringVar(&cfg.Registry.Consul.ServiceName, "registry.consul.register.name", defaultConfig.Registry.Consul.ServiceName, "service registration name")
	f.StringSliceVar(&cfg.Registry.Consul.ServiceTags, "registry.consul.register.tags", defaultConfig.Registry.Consul.ServiceTags, "service registration tags")
	f.StringSliceVar(&cfg.Registry.Consul.ServiceStatus, "registry.consul.service.status", defaultConfig.Registry.Consul.ServiceStatus, "valid service status values")
	f.DurationVar(&cfg.Registry.Consul.CheckInterval, "registry.consul.register.checkInterval", defaultConfig.Registry.Consul.CheckInterval, "service check interval")
	f.DurationVar(&cfg.Registry.Consul.CheckTimeout, "registry.consul.register.checkTimeout", defaultConfig.Registry.Consul.CheckTimeout, "service check timeout")
	f.BoolVar(&cfg.Registry.Consul.CheckTLSSkipVerify, "registry.consul.register.checkTLSSkipVerify", defaultConfig.Registry.Consul.CheckTLSSkipVerify, "service check TLS verification")
	f.StringVar(&obsoleteStr, "registry.consul.register.checkDeregisterCriticalServiceAfter", "", "This option is deprecated and has no effect.")
	f.StringVar(&cfg.Registry.Consul.ChecksRequired, "registry.consul.checksRequired", defaultConfig.Registry.Consul.ChecksRequired, "number of checks which must pass: one or all")
	f.IntVar(&cfg.Registry.Consul.ServiceMonitors, "registry.consul.serviceMonitors", defaultConfig.Registry.Consul.ServiceMonitors, "concurrency for route updates")
	f.DurationVar(&cfg.Registry.Consul.PollInterval, "registry.consul.pollinterval", defaultConfig.Registry.Consul.PollInterval, "poll interval for route updates")
	f.IntVar(&cfg.Runtime.GOGC, "runtime.gogc", defaultConfig.Runtime.GOGC, "sets runtime.GOGC")
	f.IntVar(&cfg.Runtime.GOMAXPROCS, "runtime.gomaxprocs", defaultConfig.Runtime.GOMAXPROCS, "sets runtime.GOMAXPROCS")
	f.StringVar(&cfg.UI.Access, "ui.access", defaultConfig.UI.Access, "access mode, one of [ro, rw]")
	f.StringVar(&uiListenerValue, "ui.addr", defaultValues.UIListenerValue, "Address the UI/API is listening on")
	f.StringVar(&cfg.UI.Color, "ui.color", defaultConfig.UI.Color, "background color of the UI")
	f.StringVar(&cfg.UI.Title, "ui.title", defaultConfig.UI.Title, "optional title for the UI")
	f.StringVar(&cfg.ProfileMode, "profile.mode", defaultConfig.ProfileMode, "enable profiling mode, one of [cpu, mem, mutex, block, trace]")
	f.StringVar(&cfg.ProfilePath, "profile.path", defaultConfig.ProfilePath, "path to profile dump file")
	f.BoolVar(&cfg.Tracing.TracingEnabled, "tracing.TracingEnabled", defaultConfig.Tracing.TracingEnabled, "Enable/Disable OpenTrace, one of [true, false]")
	f.StringVar(&cfg.Tracing.CollectorType, "tracing.CollectorType", defaultConfig.Tracing.CollectorType, "OpenTrace Collector Type, one of [http, kafka]")
	f.StringVar(&cfg.Tracing.ConnectString, "tracing.ConnectString", defaultConfig.Tracing.ConnectString, "OpenTrace Collector host:port")
	f.StringVar(&cfg.Tracing.ServiceName, "tracing.ServiceName", defaultConfig.Tracing.ServiceName, "Service name to embed in OpenTrace span")
	f.StringVar(&cfg.Tracing.SpanName, "tracing.SpanName", defaultConfig.Tracing.SpanName, "Span name template used to embed in OpenTrace span")
	f.StringVar(&cfg.Tracing.Topic, "tracing.Topic", defaultConfig.Tracing.Topic, "OpenTrace Collector Kafka Topic")
	f.Float64Var(&cfg.Tracing.SamplerRate, "tracing.SamplerRate", defaultConfig.Tracing.SamplerRate, "OpenTrace sample rate percentage in decimal form")
	f.StringVar(&cfg.Tracing.SpanHost, "tracing.SpanHost", defaultConfig.Tracing.SpanHost, "Host:Port info to add to spans")
	f.BoolVar(&cfg.Tracing.TraceID128Bit, "tracing.TraceID128Bit", defaultConfig.Tracing.TraceID128Bit, "Generate 128 bit trace IDs")
	f.BoolVar(&cfg.GlobMatchingDisabled, "glob.matching.disabled", defaultConfig.GlobMatchingDisabled, "Disable Glob Matching on routes, one of [true, false]")
	f.IntVar(&cfg.GlobCacheSize, "glob.cache.size", defaultConfig.GlobCacheSize, "sets the size of the glob cache")

	f.StringVar(&cfg.Registry.Custom.Host, "registry.custom.host", defaultConfig.Registry.Custom.Host, "custom back end hostname/port")
	f.StringVar(&cfg.Registry.Custom.Scheme, "registry.custom.scheme", defaultConfig.Registry.Custom.Scheme, "custom back end scheme - http/https")
	f.StringVar(&cfg.Registry.Custom.NoRouteHTML, "registry.custom.noroutehtml", defaultConfig.Registry.Custom.NoRouteHTML, "path to file for HTML returned when no route is found")
	f.BoolVar(&cfg.Registry.Custom.CheckTLSSkipVerify, "registry.custom.checkTLSSkipVerify", defaultConfig.Registry.Custom.CheckTLSSkipVerify, "custom back end check TLS verification")
	f.DurationVar(&cfg.Registry.Custom.Timeout, "registry.custom.timeout", defaultConfig.Registry.Custom.Timeout, "timeout for API request to custom back end")
	f.DurationVar(&cfg.Registry.Custom.PollInterval, "registry.custom.pollinterval", defaultConfig.Registry.Custom.PollInterval, "poll interval for API request to custom back end")
	f.StringVar(&cfg.Registry.Custom.Path, "registry.custom.path", defaultConfig.Registry.Custom.Path, "custom back end path in the URL")
	f.StringVar(&cfg.Registry.Custom.QueryParams, "registry.custom.queryparams", defaultConfig.Registry.Custom.QueryParams, "custom back end query parameters in the URL")

	// deprecated flags
	var proxyLogRoutes string
	f.StringVar(&proxyLogRoutes, "proxy.log.routes", "", "deprecated. use log.routes.format instead")

	var awsApiGWCertCN string
	f.StringVar(&awsApiGWCertCN, "aws.apigw.cert.cn", "", "deprecated. use caupgcn=<CN> for cert source")

	// parse configuration
	if err := f.ParseFlags(cmdline[1:], environ, envprefix, props); err != nil {
		return nil, err
	}

	// post configuration
	if cfg.Runtime.GOMAXPROCS == -1 {
		cfg.Runtime.GOMAXPROCS = runtime.NumCPU()
	}

	cfg.Registry.Consul.Scheme, cfg.Registry.Consul.Addr = parseScheme(cfg.Registry.Consul.Addr)

	certSources, err := parseCertSources(certSourcesValue)
	if err != nil {
		return nil, err
	}

	authSchemes, err := parseAuthSchemes(authSchemesValue)

	if err != nil {
		return nil, err
	}

	cfg.Proxy.AuthSchemes = authSchemes

	if uiListenerValue != "" {
		kvs, err := parseKVSlice(uiListenerValue)
		if err != nil {
			return nil, err
		}
		if len(kvs) != 1 {
			return nil, fmt.Errorf("ui.addr must contain only one listener")
		}
		cfg.UI.Listen, err = parseListen(kvs[0], certSources, 0, 0)
		if err != nil {
			return nil, err
		}
	}

	cfg.Listen, err = parseListeners(listenerValue, certSources, readTimeout, writeTimeout)
	if err != nil {
		return nil, err
	}

	if cfg.Proxy.LocalIP != "" {
		if cfg.Proxy.LocalIP, err = gs.Parse(cfg.Proxy.LocalIP); err != nil {
			return nil, fmt.Errorf("failed to parse local ip: %s", err)
		}
	}

	// Unless registry.consul.register.addr has been set explicitly it should
	// be the same as ui.addr. See issue 657.
	if !f.IsSet("registry.consul.register.addr") {
		cfg.Registry.Consul.ServiceAddr = cfg.UI.Listen.Addr
	}
	if cfg.Registry.Consul.ServiceAddr != "" {
		if cfg.Registry.Consul.ServiceAddr, err = gs.Parse(cfg.Registry.Consul.ServiceAddr); err != nil {
			return nil, fmt.Errorf("failed to consul service address: %s", err)
		}
	}

	cfg.Registry.Consul.CheckScheme = defaultConfig.Registry.Consul.CheckScheme
	if cfg.UI.Listen.CertSource.Name != "" {
		cfg.Registry.Consul.CheckScheme = "https"
	}

	if cfg.Registry.Consul.ServiceMonitors <= 0 {
		cfg.Registry.Consul.ServiceMonitors = 1
	}

	if gzipContentTypesValue != "" {
		cfg.Proxy.GZIPContentTypes, err = regexp.Compile(gzipContentTypesValue)
		if err != nil {
			return nil, fmt.Errorf("invalid expression for content types: %s", err)
		}
	}

	if cfg.Proxy.Strategy != "rr" && cfg.Proxy.Strategy != "rnd" {
		return nil, fmt.Errorf("invalid proxy.strategy: %s", cfg.Proxy.Strategy)
	}

	if cfg.Proxy.Matcher != "prefix" && cfg.Proxy.Matcher != "glob" && cfg.Proxy.Matcher != "iprefix" {
		return nil, fmt.Errorf("invalid proxy.matcher: %s", cfg.Proxy.Matcher)
	}

	if cfg.UI.Access != "ro" && cfg.UI.Access != "rw" {
		return nil, fmt.Errorf("invalid ui.access: %s", cfg.UI.Access)
	}

	// go1.10 will not accept a non-three digit status code
	if cfg.Proxy.NoRouteStatus < 100 || cfg.Proxy.NoRouteStatus > 999 {
		return nil, fmt.Errorf("proxy.noroutestatus must be between 100 and 999")
	}

	// handle deprecations
	deprecate := func(name, msg string) {
		if f.IsSet(name) {
			log.Print("[WARN] ", msg)
		}
	}
	deprecate("proxy.log.routes", "proxy.log.routes has been deprecated. Please use 'log.routes.format' instead")

	if proxyLogRoutes != "" {
		cfg.Log.RoutesFormat = proxyLogRoutes
	}

	return cfg, nil
}

// parseScheme splits a url into scheme and address and defaults
// to "http" if no scheme was given.
func parseScheme(s string) (scheme, addr string) {
	s = strings.ToLower(s)
	switch {
	case strings.HasPrefix(s, "https://"):
		scheme, addr = "https", s[len("https://"):]
	case strings.HasPrefix(s, "http://"):
		scheme, addr = "http", s[len("http://"):]
	default:
		scheme, addr = "http", s
	}

	// strip off anything after a final slash
	if n := strings.Index(addr, "/"); n >= 0 {
		addr = addr[:n]
	}
	return
}

func parseListeners(cfgs string, cs map[string]CertSource, readTimeout, writeTimeout time.Duration) (listen []Listen, err error) {
	kvs, err := parseKVSlice(cfgs)
	for _, cfg := range kvs {
		l, err := parseListen(cfg, cs, readTimeout, writeTimeout)
		if err != nil {
			return nil, err
		}
		listen = append(listen, l)
	}
	return
}

func parseListen(cfg map[string]string, cs map[string]CertSource, readTimeout, writeTimeout time.Duration) (l Listen, err error) {
	l = Listen{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	var csName string
	for k, v := range cfg {
		switch k {
		case "", "addr":
			if l.Addr, err = gs.Parse(v); err != nil {
				return Listen{}, err
			}
		case "proto":
			l.Proto = v
			switch l.Proto {
			case "tcp", "tcp+sni", "tcp-dynamic", "http", "https", "grpc", "grpcs", "https+tcp+sni", "prometheus":
				// ok
			default:
				return Listen{}, fmt.Errorf("unknown protocol %q", v)
			}
		case "rt": // read timeout
			d, err := time.ParseDuration(v)
			if err != nil {
				return Listen{}, err
			}
			l.ReadTimeout = d
		case "wt": // write timeout
			d, err := time.ParseDuration(v)
			if err != nil {
				return Listen{}, err
			}
			l.WriteTimeout = d
		case "it": // idle timeout
			d, err := time.ParseDuration(v)
			if err != nil {
				return Listen{}, err
			}
			l.IdleTimeout = d
		case "cs": // cert source
			csName = v
			c, ok := cs[v]
			if !ok {
				return Listen{}, fmt.Errorf("unknown certificate source %q", v)
			}
			l.CertSource = c
			if l.Proto == "" {
				l.Proto = "https"
			}
		case "strictmatch":
			l.StrictMatch = (v == "true")
		case "tlsmin":
			n, err := parseTLSVersion(v)
			if err != nil {
				return Listen{}, err
			}
			l.TLSMinVersion = n
		case "tlsmax":
			n, err := parseTLSVersion(v)
			if err != nil {
				return Listen{}, err
			}
			l.TLSMaxVersion = n
		case "tlsciphers":
			c, err := parseTLSCiphers(v)
			if err != nil {
				return Listen{}, err
			}
			l.TLSCiphers = c
		case "pxyproto":
			l.ProxyProto = (v == "true")
		case "pxytimeout":
			d, err := time.ParseDuration(v)
			if err != nil {
				return Listen{}, err
			}
			l.ProxyHeaderTimeout = d
		case "refresh":
			d, err := time.ParseDuration(v)
			if err != nil {
				return Listen{}, err
			}
			l.Refresh = d
		}
	}

	if l.Proto == "" {
		l.Proto = "http"
	}
	if l.Addr == "" {
		return Listen{}, fmt.Errorf("need listening host:port")
	}
	if csName != "" && l.Proto != "https" && l.Proto != "tcp" && l.Proto != "tcp-dynamic" && l.Proto != "grpcs" && l.Proto != "prometheus" && l.Proto != "https+tcp+sni" {
		return Listen{}, fmt.Errorf("cert source requires proto 'https', 'tcp', 'tcp-dynamic', 'https+tcp+sni', 'prometheus', or 'grpcs'")
	}
	if csName == "" && l.Proto == "https" {
		return Listen{}, fmt.Errorf("proto 'https' requires cert source")
	}
	if csName == "" && l.Proto == "grpcs" {
		return Listen{}, fmt.Errorf("proto 'grpcs' requires cert source")
	}
	if cs[csName].Type == "vault-pki" && !l.StrictMatch {
		// Without StrictMatch the first issued certificate is used for all
		// subsequent requests, even if the common name doesn't match.
		log.Printf("[INFO] vault-pki requires strictmatch; enabling strictmatch for listener %s", l.Addr)
		l.StrictMatch = true
	}
	if l.ProxyProto && l.ProxyHeaderTimeout == 0 {
		// We should define a safe default if proxy-protocol was enabled but no header timeout was set.
		// See https://github.com/fabiolb/fabio/issues/524 for more information.
		l.ProxyHeaderTimeout = 250 * time.Millisecond
	}
	return
}

var tlsver = map[string]uint16{
	"ssl30": tls.VersionSSL30,
	"tls10": tls.VersionTLS10,
	"tls11": tls.VersionTLS11,
	"tls12": tls.VersionTLS12,
}

var tlsciphers = map[string]uint16{
	"TLS_RSA_WITH_RC4_128_SHA":                0x0005,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           0x000a,
	"TLS_RSA_WITH_AES_128_CBC_SHA":            0x002f,
	"TLS_RSA_WITH_AES_256_CBC_SHA":            0x0035,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":         0x003c,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         0x009c,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         0x009d,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        0xc007,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    0xc009,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    0xc00a,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          0xc011,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     0xc012,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      0xc013,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      0xc014,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": 0xc023,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   0xc027,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   0xc02f,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": 0xc02b,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   0xc030,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": 0xc02c,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    0xcca8,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  0xcca9,
}

func parseTLSVersion(s string) (uint16, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if n, ok := tlsver[s]; ok {
		return n, nil
	}
	return parseUint16(s)
}

func parseTLSCiphers(s string) ([]uint16, error) {
	var c []uint16
	for _, v := range strings.Split(s, ",") {
		v = strings.ToUpper(strings.TrimSpace(v))
		if n, ok := tlsciphers[v]; ok {
			c = append(c, n)
			continue
		}
		n, err := parseUint16(v)
		if err != nil {
			return nil, err
		}
		c = append(c, n)
	}
	return c, nil
}

func parseUint16(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}
	if n > 1<<16 {
		return 0, fmt.Errorf("%d out of range: [0..65535]", n)
	}
	return uint16(n), nil
}

func parseCertSources(cfgs string) (cs map[string]CertSource, err error) {
	kvs, err := parseKVSlice(cfgs)
	if err != nil {
		return nil, err
	}
	cs = map[string]CertSource{}
	for _, cfg := range kvs {
		src, err := parseCertSource(cfg)
		if err != nil {
			return nil, err
		}
		cs[src.Name] = src
	}
	return
}

func parseCertSource(cfg map[string]string) (c CertSource, err error) {
	if cfg == nil {
		return CertSource{}, nil
	}

	c.Refresh = 3 * time.Second

	for k, v := range cfg {
		switch k {
		case "cs":
			c.Name = v
		case "type":
			c.Type = v
		case "cert":
			c.CertPath = v
		case "key":
			c.KeyPath = v
		case "clientca":
			c.ClientCAPath = v
		case "caupgcn":
			c.CAUpgradeCN = v
		case "refresh":
			d, err := time.ParseDuration(v)
			if err != nil {
				return CertSource{}, err
			}
			c.Refresh = d
		case "vaultfetchtoken":
			c.VaultFetchToken = v
		case "hdr":
			p := strings.SplitN(v, ": ", 2)
			if len(p) != 2 {
				return CertSource{}, fmt.Errorf("invalid header %s", v)
			}
			if c.Header == nil {
				c.Header = http.Header{}
			}
			c.Header.Set(p[0], p[1])
		}
	}
	if c.Name == "" {
		return CertSource{}, fmt.Errorf("missing 'cs' in %s", cfg)
	}
	if c.CertPath == "" {
		return CertSource{}, fmt.Errorf("missing 'cert' in %s", cfg)
	}
	switch c.Type {
	case "":
		return CertSource{}, fmt.Errorf("missing 'type' in %s", cfg)
	case "file", "consul":
		c.Refresh = 0
	case "path", "http", "vault", "vault-pki":
		// no-op
	default:
		return CertSource{}, fmt.Errorf("unknown cert source type %s", c.Type)
	}

	return
}

func parseAuthSchemes(cfgs string) (as map[string]AuthScheme, err error) {
	kvs, err := parseKVSlice(cfgs)
	if err != nil {
		return nil, err
	}
	as = map[string]AuthScheme{}
	for _, cfg := range kvs {
		src, err := parseAuthScheme(cfg)
		if err != nil {
			return nil, err
		}
		as[src.Name] = src
	}
	return
}

func parseAuthScheme(cfg map[string]string) (a AuthScheme, err error) {
	if cfg == nil {
		return
	}

	for k, v := range cfg {
		switch k {
		case "name":
			a.Name = v
		case "type":
			a.Type = v
		}
	}

	if a.Name == "" {
		return AuthScheme{}, errors.New("missing 'name' in auth")
	}

	switch a.Type {
	case "":
		return AuthScheme{}, fmt.Errorf("missing 'type' in auth '%s'", a.Name)
	case "basic":
		a.Basic = BasicAuth{
			File:    cfg["file"],
			Realm:   cfg["realm"],
			Refresh: 0, // the htpasswd file refresh is disabled by default
		}

		if a.Basic.File == "" {
			return AuthScheme{}, fmt.Errorf("missing 'file' in auth '%s'", a.Name)
		}
		if a.Basic.Realm == "" {
			a.Basic.Realm = a.Name
		}

		if cfg["refresh"] != "" {
			d, err := time.ParseDuration(cfg["refresh"])
			if err != nil {
				return AuthScheme{}, err
			}
			if d < time.Second {
				d = time.Second
			}
			a.Basic.Refresh = d
		}

	default:
		return AuthScheme{}, fmt.Errorf("unknown auth type '%s'", a.Type)
	}

	return
}
