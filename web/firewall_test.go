package web

import (
	"fmt"
	"testing"

	"github.com/shoenig/nomad-holepunch/configuration"
	"github.com/shoenig/test/must"
)

func Test_allow(t *testing.T) {
	cases := []struct {
		path  string
		rules *configuration.Firewall
		exp   bool
	}{
		{
			path:  "/ui", // non-api is never allowable
			rules: &configuration.Firewall{All: true},
			exp:   false,
		},
		{
			path:  "/v1/any",
			rules: &configuration.Firewall{All: false},
			exp:   false,
		},
		{
			path:  "/v1/any",
			rules: &configuration.Firewall{All: true},
			exp:   true,
		},
		{
			path:  "/v1/metrics",
			rules: &configuration.Firewall{Metrics: false},
			exp:   false,
		},
		{
			path:  "/v1/metrics",
			rules: &configuration.Firewall{Metrics: true},
			exp:   true,
		},
		{
			path:  "/v1/nodes",
			rules: &configuration.Firewall{Nodes: false},
			exp:   false,
		},
		{
			path:  "/v1/nodes",
			rules: &configuration.Firewall{Nodes: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/health",
			rules: &configuration.Firewall{AgentHealth: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/health",
			rules: &configuration.Firewall{AgentHealth: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/members",
			rules: &configuration.Firewall{AgentMembers: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/members",
			rules: &configuration.Firewall{AgentMembers: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/servers",
			rules: &configuration.Firewall{AgentServers: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/servers",
			rules: &configuration.Firewall{AgentServers: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/self",
			rules: &configuration.Firewall{AgentSelf: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/self",
			rules: &configuration.Firewall{AgentSelf: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/host",
			rules: &configuration.Firewall{AgentHost: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/host",
			rules: &configuration.Firewall{AgentHost: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/schedulers",
			rules: &configuration.Firewall{AgentSchedulers: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/schedulers",
			rules: &configuration.Firewall{AgentSchedulers: true},
			exp:   true,
		},
		{
			path:  "/v1/agent/schedulers/config",
			rules: &configuration.Firewall{AgentSchedulers: false},
			exp:   false,
		},
		{
			path:  "/v1/agent/schedulers/config",
			rules: &configuration.Firewall{AgentSchedulers: true},
			exp:   true,
		},
		{
			path:  "/v1/plugins",
			rules: &configuration.Firewall{Plugins: false},
			exp:   false,
		},
		{
			path:  "/v1/plugins",
			rules: &configuration.Firewall{Plugins: true},
			exp:   true,
		},
		{
			path:  "/v1/services",
			rules: &configuration.Firewall{Services: false},
			exp:   false,
		},
		{
			path:  "/v1/services",
			rules: &configuration.Firewall{Services: true},
			exp:   true,
		},
		{
			path:  "/v1/status/leader",
			rules: &configuration.Firewall{Status: false},
			exp:   false,
		},
		{
			path:  "/v1/status/leader",
			rules: &configuration.Firewall{Status: true},
			exp:   true,
		},
		{
			path:  "/v1/status/peers",
			rules: &configuration.Firewall{Status: false},
			exp:   false,
		},
		{
			path:  "/v1/status/peers",
			rules: &configuration.Firewall{Status: true},
			exp:   true,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%s#%t", tc.path, tc.exp)
		t.Run(name, func(t *testing.T) {
			f := &firewall{ruleset: initialize(tc.rules)}
			result := f.allow(tc.path)
			must.Eq(t, tc.exp, result)
		})
	}
}
