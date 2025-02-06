package auth

import (
	"fmt"
	"net/http"

	"github.com/fabiolb/fabio/config"
)

type AuthDecision struct {
	Authorized bool
	Done       bool
}

func authorized() AuthDecision {
	return AuthDecision{Authorized: true, Done: false}
}

func unauthorized() AuthDecision {
	return AuthDecision{Authorized: false, Done: false}
}

type AuthScheme interface {
	Authorized(request *http.Request, response http.ResponseWriter) AuthDecision
}

func LoadAuthSchemes(cfg map[string]config.AuthScheme) (map[string]AuthScheme, error) {
	auths := map[string]AuthScheme{}
	for _, a := range cfg {
		switch a.Type {
		case "basic":
			b, err := newBasicAuth(a.Basic)
			if err != nil {
				return nil, err
			}
			auths[a.Name] = b
		case "external":
			d, err := newExternalAuth(a.External)
			if err != nil {
				return nil, err
			}
			auths[a.Name] = d
		default:
			return nil, fmt.Errorf("unknown auth type '%s'", a.Type)
		}
	}

	return auths, nil
}
