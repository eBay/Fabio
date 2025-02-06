package auth

import (
	"bytes"
	"log"
	"net/http"

	"github.com/fabiolb/fabio/config"
)

type external struct {
	endpoint          string
	appendAuthHeaders []string
	setAuthHeaders    []string
}

func newExternalAuth(cfg config.ExternalAuth) (AuthScheme, error) {
	return &external{
		endpoint:          cfg.Endpoint,
		appendAuthHeaders: cfg.AppendAuthHeaders,
		setAuthHeaders:    cfg.SetAuthHeaders,
	}, nil
}

func (b *external) Authorized(request *http.Request, response http.ResponseWriter) AuthDecision {
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}

	authRequest, err := http.NewRequest(request.Method, b.endpoint+request.URL.RequestURI(), bytes.NewReader(nil))
	if err != nil {
		log.Println("[ERROR] Can't make auth external request value:", err.Error())
		return unauthorized()
	}
	authRequest.Header = request.Header
	authRequest.Host = request.Host

	resp, err := client.Do(authRequest)
	if err != nil {
		log.Println("[ERROR] External request error:", err.Error())
		return unauthorized()
	}

	if resp.StatusCode == 200 {
		for _, header := range b.setAuthHeaders {
			response.Header().Set(header, resp.Header.Get(header))
		}
		for _, header := range b.appendAuthHeaders {
			response.Header().Add(header, resp.Header.Get(header))
		}
		return authorized()
	}

	if resp.StatusCode == 302 {
		http.Redirect(response, request, resp.Header.Get("location"), 302)
		return AuthDecision{Authorized: false, Done: true}
	}

	if resp.StatusCode != 401 {
		log.Println("[WARN] Unexpected status code", resp.StatusCode, "treated as unauthorized")
		return unauthorized()
	}

	return unauthorized()
}
