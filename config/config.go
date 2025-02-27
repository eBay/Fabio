package config

import (
	"net/http"
	"regexp"
	"time"
)

type Config struct {
	Proxy                Proxy
	Registry             Registry
	Listen               []Listen
	Log                  Log
	Metrics              Metrics
	UI                   UI
	Runtime              Runtime
	ProfileMode          string
	ProfilePath          string
	Insecure             bool
	GlobMatchingDisabled bool
	GlobCacheSize        int
	BGP                  BGP
}

type CertSource struct {
	Name            string
	Type            string
	CertPath        string
	KeyPath         string
	ClientCAPath    string
	CAUpgradeCN     string
	Refresh         time.Duration
	Header          http.Header
	VaultFetchToken string
}

type Listen struct {
	Addr               string
	Proto              string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	CertSource         CertSource
	StrictMatch        bool
	TLSMinVersion      uint16
	TLSMaxVersion      uint16
	TLSCiphers         []uint16
	ProxyProto         bool
	ProxyHeaderTimeout time.Duration
	Refresh            time.Duration
}

type Source struct {
	LinkEnabled bool
	NewTab      bool
	Scheme      string
	Host        string
	Port        string
}

type RoutingTable struct {
	Source Source
}

type UI struct {
	Listen       Listen
	Color        string
	Title        string
	Access       string
	RoutingTable RoutingTable
}

type Proxy struct {
	Strategy              string
	Matcher               string
	NoRouteStatus         int
	MaxConn               int
	ShutdownWait          time.Duration
	DeregisterGracePeriod time.Duration
	DialTimeout           time.Duration
	ResponseHeaderTimeout time.Duration
	KeepAliveTimeout      time.Duration
	IdleConnTimeout       time.Duration
	FlushInterval         time.Duration
	GlobalFlushInterval   time.Duration
	LocalIP               string
	ClientIPHeader        string
	TLSHeader             string
	TLSHeaderValue        string
	GZIPContentTypes      *regexp.Regexp
	RequestID             string
	STSHeader             STSHeader
	AuthSchemes           map[string]AuthScheme
	GRPCMaxRxMsgSize      int
	GRPCMaxTxMsgSize      int
	GRPCGShutdownTimeout  time.Duration
}

type STSHeader struct {
	MaxAge     int
	Subdomains bool
	Preload    bool
}

type Runtime struct {
	GOGC       int
	GOMAXPROCS int
}

type Circonus struct {
	APIKey        string
	APIApp        string
	APIURL        string
	CheckID       string
	BrokerID      string
	SubmissionURL string
}

type Prometheus struct {
	Subsystem string
	Path      string
	Buckets   []float64
}

type Log struct {
	AccessFormat string
	AccessTarget string
	RoutesFormat string
	Level        string
}

type Metrics struct {
	Target        string
	Prefix        string
	Names         string
	Interval      time.Duration
	Timeout       time.Duration
	Retry         time.Duration
	GraphiteAddr  string
	StatsDAddr    string
	DogstatsdAddr string
	Circonus      Circonus
	Prometheus    Prometheus
}

type Registry struct {
	Backend string
	Static  Static
	File    File
	Consul  Consul
	Custom  Custom
	Timeout time.Duration
	Retry   time.Duration
}

type Static struct {
	NoRouteHTML string
	Routes      string
}

type File struct {
	NoRouteHTMLPath string
	RoutesPath      string
}

type Consul struct {
	Addr               string
	Scheme             string
	Token              string
	KVPath             string
	NoRouteHTMLPath    string
	TagPrefix          string
	Register           bool
	ServiceAddr        string
	ServiceName        string
	ServiceTags        []string
	ServiceStatus      []string
	CheckInterval      time.Duration
	CheckTimeout       time.Duration
	CheckScheme        string
	CheckTLSSkipVerify bool
	ChecksRequired     string
	ServiceMonitors    int
	TLS                ConsulTlS
	PollInterval       time.Duration
	RequireConsistent  bool
	AllowStale         bool
}

type Custom struct {
	Host               string
	Path               string
	QueryParams        string
	Scheme             string
	CheckTLSSkipVerify bool
	PollInterval       time.Duration
	NoRouteHTML        string
	Timeout            time.Duration
}

type AuthScheme struct {
	Name  string
	Type  string
	Basic BasicAuth
}

type BasicAuth struct {
	Realm   string
	File    string
	Refresh time.Duration
	ModTime time.Time // the htpasswd file last modification time
}

type ConsulTlS struct {
	KeyFile            string
	CertFile           string
	CAFile             string
	CAPath             string
	InsecureSkipVerify bool
}

type BGP struct {
	BGPEnabled        bool
	Asn               uint
	AnycastAddresses  []string
	RouterID          string
	ListenPort        int
	ListenAddresses   []string
	Peers             []BGPPeer
	EnableGRPC        bool
	GRPCListenAddress string
	GRPCTLS           bool
	CertFile          string
	KeyFile           string
	GOBGPDCfgFile     string
	NextHop           string
}

type BGPPeer struct {
	NeighborAddress string
	NeighborPort    uint
	Asn             uint
	MultiHop        bool
	MultiHopLength  uint
	Password        string
}
