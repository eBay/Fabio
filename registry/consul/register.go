package consul

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fabiolb/fabio/config"
	"github.com/hashicorp/consul/api"
)

const (
	TTLInterval                       = time.Second * 15
	TTLRefreshInterval                = time.Second * 10
	TTLDeregisterCriticalServiceAfter = time.Minute
)

// register keeps a service registered in consul.
//
// When a value is sent in the dereg channel the service is deregistered from
// consul. To wait for completion the caller should read the next value from
// the dereg channel.
//
//    dereg <- true // trigger deregistration
//    <-dereg       // wait for completion
//
func register(c *api.Client, service *api.AgentServiceRegistration) chan bool {
	registered := func(serviceID string) bool {
		if serviceID == "" {
			return false
		}
		services, err := c.Agent().Services()
		if err != nil {
			log.Printf("[ERROR] consul: Cannot get service list. %s", err)
			return false
		}
		return services[serviceID] != nil
	}

	register := func() string {
		if err := c.Agent().ServiceRegister(service); err != nil {
			log.Printf("[ERROR] consul: Cannot register fabio [name:%q] in Consul. %s", service.Name, err)
			return ""
		}

		log.Printf("[INFO] consul: Registered fabio as %q", service.Name)
		log.Printf("[INFO] consul: Registered fabio with id %q", service.ID)
		log.Printf("[INFO] consul: Registered fabio with address %q", service.Address)
		log.Printf("[INFO] consul: Registered fabio with tags %q", strings.Join(service.Tags, ","))
		for _, check := range service.Checks {
			log.Printf("[INFO] consul: Registered fabio with check %+v", check)
		}

		return service.ID
	}

	deregister := func(serviceID string) {
		log.Printf("[INFO] consul: Deregistering %q", service.Name)
		c.Agent().ServiceDeregister(serviceID)
	}

	passTTL := func(serviceTTLID string) {
		c.Agent().UpdateTTL(serviceTTLID, "", api.HealthPassing)
	}

	dereg := make(chan bool)
	go func() {
		var serviceID string
		var serviceTTLCheckId string

		for {
			if !registered(serviceID) {
				serviceID = register()
				serviceTTLCheckId = computeServiceTTLCheckId(serviceID)
				// Pass the TTL check right now so traffic can be served immediately.
				passTTL(serviceTTLCheckId)
			}

			select {
			case <-dereg:
				deregister(serviceID)
				dereg <- true
				return
			case <-time.After(TTLRefreshInterval):
				// Reset the TTL check clock.
				passTTL(serviceTTLCheckId)
			}
		}
	}()
	return dereg
}

func serviceRegistration(cfg *config.Consul, serviceName string) (*api.AgentServiceRegistration, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	ipstr, portstr, err := net.SplitHostPort(cfg.ServiceAddr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(ipstr)
	if ip == nil {
		ip, err = config.LocalIP()
		if err != nil {
			return nil, err
		}
		if ip == nil {
			return nil, errors.New("no local ip")
		}
	}

	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, hostname, port)

	checkURL := fmt.Sprintf("%s://%s:%d/health", cfg.CheckScheme, ip, port)
	if ip.To16() != nil {
		checkURL = fmt.Sprintf("%s://[%s]:%d/health", cfg.CheckScheme, ip, port)
	}

	service := &api.AgentServiceRegistration{
		ID:        serviceID,
		Name:      serviceName,
		Address:   ip.String(),
		Port:      port,
		Tags:      cfg.ServiceTags,
		Namespace: cfg.Namespace,
		// Set the checks for the service.
		//
		// Both checks must pass for Consul to consider the service healthy and therefore serve the fabio instance to clients.
		Checks: []*api.AgentServiceCheck{
			// If fabio doesn't exit cleanly, it doesn't auto-deregister itself
			// from Consul. In order to address this, we introduce a TTL check
			// to confirm that the fabio instance is alive and able to route
			// this service.
			//
			// The TTL check must be refreshed before its timeout is crossed,
			// otherwise the check fails. If the check fails, Consul considers
			// this service to have become unhealthy. If the check is failing
			// (critical) for the DeregisterCriticalServiceAfter duration, the
			// Consul reaper will remove it from Consul.
			//
			// For more info, read https://www.consul.io/api/agent/check.html#deregistercriticalserviceafter.
			{
				CheckID:                        computeServiceTTLCheckId(serviceID),
				TTL:                            TTLInterval.String(),
				DeregisterCriticalServiceAfter: TTLDeregisterCriticalServiceAfter.String(),
			},
			// HTTP check is meant to confirm fabio health endpoint is
			// reachable from the Consul agent. If the check fails, Consul
			// considers this service to have become unhealthy.
			{
				HTTP:          checkURL,
				Interval:      cfg.CheckInterval.String(),
				Timeout:       cfg.CheckTimeout.String(),
				TLSSkipVerify: cfg.CheckTLSSkipVerify,
			},
		},
	}

	return service, nil
}

func computeServiceTTLCheckId(serviceID string) string {
	return serviceID + "-ttl"
}
