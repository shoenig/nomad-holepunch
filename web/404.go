package web

import (
	"net/http"

	"github.com/shoenig/loggy"
)

type fallback struct {
	log loggy.Logger
}

func newFallback() http.Handler {
	return &fallback{
		log: loggy.New("fallback"),
	}
}

func (f *fallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.log.Warnf("serving 404 to %s", r.RemoteAddr)
	http.Error(w, "try another endpoint", http.StatusNotFound)
}
