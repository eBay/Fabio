package auth

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fabiolb/fabio/config"
	"github.com/tg123/go-htpasswd"
)

// basic is an implementation of AuthScheme
type basic struct {
	realm   string
	secrets *htpasswd.File
}

func newBasicAuth(cfg config.BasicAuth) (AuthScheme, error) {
	bad := func(err error) {
		log.Println("[WARN] Error processing a line in an htpasswd file:", err)
	}

	secrets, err := htpasswd.New(cfg.File, htpasswd.DefaultSystems, bad)
	if err != nil {
		return nil, err
	}

	if cfg.Refresh > 0 {
		stat, err := os.Stat(cfg.File)
		if err != nil {
			return nil, err
		}
		cfg.ModTime = stat.ModTime()

		go func() {
			cleared := false
			ticker := time.NewTicker(cfg.Refresh).C
			for range ticker {
				stat, err := os.Stat(cfg.File)
				if err != nil {
					log.Println("[WARN] Error accessing htpasswd file:", err)
					if !cleared {
						err = secrets.ReloadFromReader(&bytes.Buffer{}, bad)
						if err != nil {
							log.Println("[WARN] Error clearing the htpasswd credentials:", err)
						} else {
							log.Println("[INFO] The htpasswd credentials have been cleared")
							cleared = true
						}
					}
					continue
				}

				// refresh the htpasswd file only if its modification time has changed
				// even if the new htpasswd file is older than previously loaded
				if cfg.ModTime != stat.ModTime() {
					if err := secrets.Reload(bad); err == nil {
						log.Println("[INFO] The htpasswd file has been successfully reloaded")
						cfg.ModTime = stat.ModTime()
						cleared = false
					} else {
						log.Println("[WARN] Error reloading htpasswd file:", err)
					}
				}
			}
		}()
	}

	return &basic{
		secrets: secrets,
		realm:   cfg.Realm,
	}, nil
}

func (b *basic) Authorized(request *http.Request, response http.ResponseWriter) AuthDecision {
	user, password, ok := request.BasicAuth()

	if !ok {
		response.Header().Set("WWW-Authenticate", "Basic realm=\""+b.realm+"\"")
		return unauthorized()
	}

	if b.secrets.Match(user, password) {
		return authorized()
	} else {
		return unauthorized()
	}
}
