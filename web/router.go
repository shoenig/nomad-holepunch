package web

import (
	"net/http"

	"github.com/shoenig/nomad-holepunch/configuration"
)

func New(config *configuration.Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/health", newHealth())
	mux.Handle("/v1/", newFirewall(config.Rules, newProxy(config)))
	mux.Handle("/", newFallback())
	return mux
}
