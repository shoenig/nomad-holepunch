package web

import (
	"net/http"

	"github.com/shoenig/loggy"
)

type health struct {
	log loggy.Logger
}

func newHealth() http.Handler {
	return &health{
		log: loggy.New("health"),
	}
}

func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("serving health endpoint to %s", r.RemoteAddr)
	http.Error(w, "ok", http.StatusOK)
}
