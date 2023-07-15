package configuration

import (
	"path/filepath"

	"github.com/shoenig/extractors/env"
	"github.com/shoenig/go-conceal"
	"github.com/shoenig/loggy"
)

type Config struct {
	Bind       string
	Port       string
	NomadToken *conceal.Text
	SocketPath string
	Rules      *Firewall
}

type Firewall struct {
	All             bool // /v1/*
	Metrics         bool // /v1/metrics              (allow by default)
	Nodes           bool // /v1/nodes
	AgentHealth     bool // /v1/agent/health         (allow by default)
	AgentMembers    bool // /v1/agent/members
	AgentServers    bool // /v1/agent/servers
	AgentSelf       bool // /v1/agent/self
	AgentHost       bool // /v1/agent/host
	AgentSchedulers bool // /v1/agent/schedulers[/config]
	Plugins         bool // /v1/plugins
	Services        bool // /v1/service[s/*]
	Status          bool // /v1/status/[leader,peers]
}

func (c *Config) Log(log loggy.Logger) {
	log.Tracef("HOLEPUNCH_BIND = %s", c.Bind)
	log.Tracef("HOLEPUNCH_PORT = %s", c.Port)
	log.Tracef("HOLEPUNCH_TOKEN = %s", c.NomadToken)
	log.Tracef("HOLEPUNCH_SOCKET_PATH = %s", c.SocketPath)
	log.Tracef("HOLEPUNCH_ALLOW_ALL = %t", c.Rules.All)
	log.Tracef("HOLEPUNCH_ALLOW_METRICS = %t", c.Rules.Metrics)
	log.Tracef("HOLEPUNCH_ALLOW_NODES = %t", c.Rules.Nodes)
	log.Tracef("HOLEPUNCH_ALLOW_AGENT_HEALTH = %t", c.Rules.AgentHealth)
	log.Tracef("HOLEPUNCH_ALLOW_AGENT_MEMBERS = %t", c.Rules.AgentMembers)
	log.Tracef("HOLEPUNCH_ALLOW_AGENT_SERVERS = %t", c.Rules.AgentServers)
	log.Tracef("HOLEPUNCH_ALLOW_AGENT_SELF = %t", c.Rules.AgentSelf)
	log.Tracef("HOLEPUNCH_ALLOW_AGENT_HOST = %t", c.Rules.AgentHost)
	log.Tracef("HOLEPUNCH_ALLOW_AGENT_SCHEDULERS = %t", c.Rules.AgentSchedulers)
	log.Tracef("HOLEPUNCH_ALLOW_PLUGINS = %t", c.Rules.Plugins)
	log.Tracef("HOLEPUNCH_ALLOW_SERVICES= %t", c.Rules.Services)
	log.Tracef("HOLEPUNCH_ALLOW_STATUS = %t", c.Rules.Status)
}

func Load(log loggy.Logger) *Config {
	c := new(Config)
	c.Rules = new(Firewall)

	err := env.ParseOS(env.Schema{
		// setup
		"HOLEPUNCH_BIND":        env.StringOr(&c.Bind, "0.0.0.0"),
		"HOLEPUNCH_PORT":        env.StringOr(&c.Port, "6120"),
		"NOMAD_TOKEN":           env.Secret(&c.NomadToken, true),
		"HOLEPUNCH_SOCKET_PATH": env.StringOr(&c.SocketPath, defaultSocketPath()),

		// firewall
		"HOLEPUNCH_ALLOW_ALL":              env.BoolOr(&c.Rules.All, false),
		"HOLEPUNCH_ALLOW_METRICS":          env.BoolOr(&c.Rules.Metrics, true),
		"HOLEPUNCH_ALLOW_NODES":            env.BoolOr(&c.Rules.Nodes, false),
		"HOLEPUNCH_ALLOW_AGENT_HEALTH":     env.BoolOr(&c.Rules.AgentHealth, true),
		"HOLEPUNCH_ALLOW_AGENT_MEMBERS":    env.BoolOr(&c.Rules.AgentMembers, false),
		"HOLEPUNCH_ALLOW_AGENT_SERVERS":    env.BoolOr(&c.Rules.AgentServers, false),
		"HOLEPUNCH_ALLOW_AGENT_SELF":       env.BoolOr(&c.Rules.AgentSelf, false),
		"HOLEPUNCH_ALLOW_AGENT_HOST":       env.BoolOr(&c.Rules.AgentHost, false),
		"HOLEPUNCH_ALLOW_AGENT_SCHEDULERS": env.BoolOr(&c.Rules.AgentSchedulers, false),
		"HOLEPUNCH_ALLOW_PLUGINS":          env.BoolOr(&c.Rules.Plugins, false),
		"HOLEPUNCH_ALLOW_SERVICES":         env.BoolOr(&c.Rules.Services, false),
		"HOLEPUNCH_ALLOW_STATUS":           env.BoolOr(&c.Rules.Status, false),
	})
	if err != nil {
		log.Errorf("error extracting config from environment: %v", err)
	}

	return c
}

func defaultSocketPath() string {
	var dir string
	_ = env.ParseOS(env.Schema{
		"NOMAD_SECRETS_DIR": env.StringOr(&dir, "/secrets"),
	})
	socket := filepath.Join(dir, "api.sock")
	return socket
}
