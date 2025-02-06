package route

import (
	"log"
	"net/http"

	"github.com/fabiolb/fabio/auth"
)

func (t *Target) Authorized(r *http.Request, w http.ResponseWriter, authSchemes map[string]auth.AuthScheme) auth.AuthDecision {
	if t.AuthScheme == "" {
		return auth.AuthDecision{Authorized: true, Done: false}
	}

	scheme := authSchemes[t.AuthScheme]

	if scheme == nil {
		log.Printf("[ERROR] unknown auth scheme '%s'\n", t.AuthScheme)
		return auth.AuthDecision{Authorized: false, Done: false}
	}

	return scheme.Authorized(r, w)
}
