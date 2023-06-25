package web

import (
	"net/http"
	"strings"

	"github.com/shoenig/loggy"
	"github.com/shoenig/nomad-holepunch/configuration"
)

type firewall struct {
	log   loggy.Logger
	rules *configuration.Firewall
	next  http.Handler
}

func newFirewall(rules *configuration.Firewall, next http.Handler) http.Handler {
	return &firewall{
		log:   loggy.New("firewall"),
		rules: rules,
		next:  next,
	}
}

func (f *firewall) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fast path where the firewall is disabled
	if f.rules.Disable {
		f.next.ServeHTTP(w, r)
		return
	}

	// get the original path
	path := r.URL.Path

	// fast path where endpoint is not the api
	if !strings.HasPrefix(path, "/v1/") {
		http.Error(w, "non-api access is diabled", http.StatusForbidden)
		return
	}

	// strip off /v1/ prefix, we can assume api endpoints from here
	path = strings.TrimPrefix(path, "/v1/")

	elem0 := strings.Split(path, "/")[0]
	f.log.Infof("path is %q, elem0 is %s", path, elem0)

	// check if the firewall rules allow this
	// this could do with a good bit of refactoring
	if elem0 == "metrics" && f.rules.AllowMetrics {
		f.next.ServeHTTP(w, r)
		return
	}

	http.Error(w, "forbidden", http.StatusForbidden)
}
