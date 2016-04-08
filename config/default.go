package config

import (
	"runtime"
	"time"
)

var Default = &Config{
	Proxy: Proxy{
		MaxConn:      10000,
		Strategy:     "rnd",
		DialTimeout:  30 * time.Second,
		LocalIP:      LocalIPString(),
		ReadTimeout:  time.Duration(0),
		WriteTimeout: time.Duration(0),
	},
	Registry: Registry{
		Backend: "consul",
		Consul: Consul{
			Addr:          "localhost:8500",
			KVPath:        "/fabio/config",
			TagPrefix:     "urlprefix-",
			ServiceAddr:   ":9998",
			ServiceName:   "fabio",
			CheckInterval: time.Second,
			CheckTimeout:  3 * time.Second,
		},
	},
	Listen: []Listen{
		{
			Addr: ":9999",
		},
	},
	Runtime: Runtime{
		GOGC:       800,
		GOMAXPROCS: runtime.NumCPU(),
	},
	UI: UI{
		Addr:  ":9998",
		Color: "light-green",
	},
	Metrics: []Metrics{
		{
			Target:   "",
			Prefix:   "default",
			Addr:     "",
			Interval: 30 * time.Second,
		},
	},
}
