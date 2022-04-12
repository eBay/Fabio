package consul

import (
	"errors"
	"log"

	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/registry"

	"github.com/hashicorp/consul/api"
)

// be is an implementation of a registry backend for consul.
type be struct {
	c     *api.Client
	dc    string
	cfg   *config.Consul
	dereg map[string](chan bool)
}

func NewBackend(cfg *config.Consul) (registry.Backend, error) {

	consulCfg := &api.Config{Address: cfg.Addr, Scheme: cfg.Scheme, Token: cfg.Token, Namespace: cfg.Namespace}
	if cfg.Scheme == "https" {
		consulCfg.TLSConfig.KeyFile = cfg.TLS.KeyFile
		consulCfg.TLSConfig.CertFile = cfg.TLS.CertFile
		consulCfg.TLSConfig.CAFile = cfg.TLS.CAFile
		consulCfg.TLSConfig.CAPath = cfg.TLS.CAPath
		consulCfg.TLSConfig.InsecureSkipVerify = cfg.TLS.InsecureSkipVerify
	}

	// create a reusable client
	c, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	// ping the agent
	dc, err := datacenter(c)
	if err != nil {
		return nil, err
	}

	// we're good
	log.Printf("[INFO] consul: Connecting to %q in datacenter %q", cfg.Addr, dc)
	if cfg.Namespace != "" {
		log.Printf("[INFO] consul: Connecting to namespace %q", cfg.Namespace)
	} else {
		log.Printf("[INFO] consul: Connecting to default namespace")
	}
	return &be{c: c, dc: dc, cfg: cfg}, nil
}

func (b *be) Register(services []string) error {
	if b.dereg == nil {
		b.dereg = make(map[string](chan bool))
	}

	if b.cfg.Register {
		services = append(services, b.cfg.ServiceName)
	}

	// deregister unneeded services
	for service := range b.dereg {
		if stringInSlice(service, services) {
			continue
		}
		err := b.Deregister(service)
		if err != nil {
			return err
		}
	}

	// register new services
	for _, service := range services {
		if b.dereg[service] != nil {
			log.Printf("[DEBUG] %q already registered", service)
			continue
		}

		serviceReg, err := serviceRegistration(b.cfg, service)
		if err != nil {
			return err
		}

		b.dereg[service] = register(b.c, serviceReg)
	}

	return nil
}

func (b *be) Deregister(service string) error {
	dereg := b.dereg[service]
	if dereg == nil {
		log.Printf("[WARN]: Attempted to deregister unknown service %q", service)
		return nil
	}
	dereg <- true // trigger deregistration
	<-dereg       // wait for completion
	delete(b.dereg, service)

	return nil
}

func (b *be) DeregisterAll() error {
	log.Printf("[DEBUG]: consul: Deregistering all registered aliases.")
	for _, dereg := range b.dereg {
		if dereg == nil {
			continue
		}
		dereg <- true // trigger deregistration
		<-dereg       // wait for completion
	}
	return nil
}

func (b *be) ManualPaths() ([]string, error) {
	keys, _, err := listKeys(b.c, b.cfg.KVPath, 0)
	return keys, err
}

func (b *be) ReadManual(path string) (value string, version uint64, err error) {
	// we cannot rely on the value provided by WatchManual() since
	// someone has to call that method first to kick off the go routine.
	return getKV(b.c, b.cfg.KVPath+path, 0)
}

func (b *be) WriteManual(path string, value string, version uint64) (ok bool, err error) {
	// try to create the key first by using version 0
	if ok, err = putKV(b.c, b.cfg.KVPath+path, value, 0); ok {
		return
	}

	// then try the CAS update
	return putKV(b.c, b.cfg.KVPath+path, value, version)
}

func (b *be) WatchServices() chan string {
	log.Printf("[INFO] consul: Using dynamic routes")
	log.Printf("[INFO] consul: Using tag prefix %q", b.cfg.TagPrefix)

	m := NewServiceMonitor(b.c, b.cfg, b.dc)
	svc := make(chan string)
	go m.Watch(svc)
	return svc
}

func (b *be) WatchManual() chan string {
	log.Printf("[INFO] consul: Watching KV path %q", b.cfg.KVPath)

	kv := make(chan string)
	go watchKV(b.c, b.cfg.KVPath, kv, true)
	return kv
}

func (b *be) WatchNoRouteHTML() chan string {
	log.Printf("[INFO] consul: Watching KV path %q", b.cfg.NoRouteHTMLPath)

	html := make(chan string)
	go watchKV(b.c, b.cfg.NoRouteHTMLPath, html, false)
	return html
}

// datacenter returns the datacenter of the local agent
func datacenter(c *api.Client) (string, error) {
	self, err := c.Agent().Self()
	if err != nil {
		return "", err
	}

	cfg, ok := self["Config"]
	if !ok {
		return "", errors.New("consul: self.Config not found")
	}
	dc, ok := cfg["Datacenter"].(string)
	if !ok {
		return "", errors.New("consul: self.Datacenter not found")
	}
	return dc, nil
}

func stringInSlice(str string, strSlice []string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}
	return false
}
