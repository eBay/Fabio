package consul

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
)

// routecmd builds a route command.
type routecmd struct {
	// svc is the consul service instance.
	svc *api.CatalogService

	// prefix is the prefix of urlprefix tags. e.g. 'urlprefix-'.
	prefix string

	env map[string]string
}

func (r routecmd) build() []string {
	var svctags, routetags []string
	for _, t := range r.svc.ServiceTags {

		t = strings.TrimSpace(t)

		if strings.HasPrefix(t, r.prefix) {
			routetags = append(routetags, t)
		} else {
			svctags = append(svctags, t)
		}
	}

	// generate route commands
	var config []string
	for _, tag := range routetags {
		if route, opts, ok := parseURLPrefixTag(tag, r.prefix, r.env); ok {
			name, addr, port := r.svc.ServiceName, r.svc.ServiceAddress, r.svc.ServicePort

			// use consul node address if service address is not set
			if addr == "" {
				addr = r.svc.Address
			}

			// add .local suffix on OSX for simple host names w/o domain
			if runtime.GOOS == "darwin" && !strings.Contains(addr, ".") && !strings.HasSuffix(addr, ".local") {
				addr += ".local"
			}

			addr = net.JoinHostPort(addr, strconv.Itoa(port))
			//tags := strings.Join(r.tags, ",")
			dst := "http://" + addr + "/"

			var weight string
			var ropts []string
			for _, o := range strings.Fields(opts) {
				switch {
				case o == "proto=tcp":
					dst = "tcp://" + addr

				case o == "proto=https":
					dst = "https://" + addr

				case o == "proto=grpcs":
					dst = "grpcs://" + addr

				case o == "proto=grpc":
					dst = "grpc://" + addr

				case strings.HasPrefix(o, "weight="):
					weight = o[len("weight="):]

				case strings.HasPrefix(o, "redirect="):
					redir := strings.Split(o[len("redirect="):], ",")
					if len(redir) == 2 {
						dst = redir[1]
						ropts = append(ropts, fmt.Sprintf("redirect=%s", redir[0]))
					} else {
						log.Printf("[ERROR] Invalid syntax for redirect: %s. should be redirect=<code>,<url>", o)
						continue
					}
				default:
					ropts = append(ropts, o)
				}
			}

			cfg := "route add " + name + " " + route + " " + dst
			if weight != "" {
				cfg += " weight " + weight
			}
			if len(svctags) > 0 {
				cfg += " tags " + strconv.Quote(strings.Join(svctags, ","))
			}
			if len(ropts) > 0 {
				cfg += " opts " + strconv.Quote(strings.Join(ropts, " "))
			}

			config = append(config, cfg)
		}
	}
	return config
}

// parseURLPrefixTag expects an input in the form of 'tag-host/path[ opts]'
// and returns the lower cased host and the unaltered path if the
// prefix matches the tag.
func parseURLPrefixTag(s, prefix string, env map[string]string) (route, opts string, ok bool) {
	// expand $x or ${x} to env[x] or ""
	expand := func(s string) string {
		return os.Expand(s, func(x string) string {
			if env == nil {
				return "$" + x
			}
			if _, ok := env[x]; !ok {
				return "$" + x
			}
			return env[x]
		})
	}

	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, prefix) {
		return "", "", false
	}
	s = strings.TrimSpace(s[len(prefix):])

	p := strings.SplitN(s, " ", 2)
	if len(p) == 2 {
		opts = p[1]
	}
	s = p[0]

	// prefix is ":port"
	if strings.HasPrefix(s, ":") {
		return s, opts, true
	}

	if !strings.Contains(s, "/") {
		return s, opts, true
	}
	// prefix is "host/path"
	p = strings.SplitN(s, "/", 2)
	if len(p) == 1 {
		log.Printf("[WARN] consul: Invalid %s tag %q - You need to have a trailing slash!", prefix, s)
		return "", "", false
	}
	host, path := p[0], p[1]

	return strings.ToLower(expand(host)) + "/" + expand(path), opts, true
}
