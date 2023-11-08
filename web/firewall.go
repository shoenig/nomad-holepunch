package web

import (
	"net/http"
	"strings"

	"github.com/hashicorp/go-set/v2"
	"github.com/shoenig/loggy"
	"github.com/shoenig/nomad-holepunch/configuration"
)

type firewall struct {
	log     loggy.Logger
	ruleset *set.Set[string]
	next    http.Handler
}

func newFirewall(rules *configuration.Firewall, next http.Handler) http.Handler {
	return &firewall{
		log:     loggy.New("firewall"),
		ruleset: initialize(rules),
		next:    next,
	}
}

func (f *firewall) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the original path
	path := r.URL.Path

	f.log.Tracef("checking path %q", path)

	// compare path with allowable rules
	if f.allow(path) {
		f.next.ServeHTTP(w, r)
		return
	}

	// otherwise the request is forbidden
	http.Error(w, "forbidden", http.StatusForbidden)
}

func (f *firewall) allow(path string) bool {
	// fast path where endpoint is not the api
	if !strings.HasPrefix(path, "/v1/") {
		return false
	}

	// fast path for allow all
	if f.ruleset.Contains("*") {
		return true
	}

	// check for service lookup
	if strings.HasPrefix(path, "/v1/service/") && f.ruleset.Contains("/v1/services") {
		return true
	}

	// check if ruleset allows access to path
	if f.ruleset.Contains(path) {
		return true
	}

	return false
}

func initialize(rules *configuration.Firewall) *set.Set[string] {
	s := set.New[string](10)
	if rules.All {
		s.Insert("*")
	}
	if rules.Metrics {
		s.Insert("/v1/metrics")
	}
	if rules.Nodes {
		s.Insert("/v1/nodes")
	}
	if rules.AgentHealth {
		s.Insert("/v1/agent/health")
	}
	if rules.AgentMembers {
		s.Insert("/v1/agent/members")
	}
	if rules.AgentServers {
		s.Insert("/v1/agent/servers")
	}
	if rules.AgentSelf {
		s.Insert("/v1/agent/self")
	}
	if rules.AgentHost {
		s.Insert("/v1/agent/host")
	}
	if rules.AgentSchedulers {
		s.Insert("/v1/agent/schedulers")
		s.Insert("/v1/agent/schedulers/config")
	}
	if rules.Plugins {
		s.Insert("/v1/plugins")
	}
	if rules.Services {
		s.Insert("/v1/services")
	}
	if rules.Status {
		s.Insert("/v1/status/leader")
		s.Insert("/v1/status/peers")
	}
	return s
}
