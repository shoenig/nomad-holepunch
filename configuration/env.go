package configuration

import (
	"path/filepath"

	"github.com/shoenig/extractors/env"
	"github.com/shoenig/go-conceal"
	"github.com/shoenig/loggy"
)

type Config struct {
	Bind          string
	Port          string
	NomadToken    *conceal.Text
	SocketPath    string
	Authorization *Firewall
}

type Firewall struct {
	AllowAll     bool
	AllowMetrics bool
}

func (c *Config) Log(log loggy.Logger) {
	log.Tracef("HOLEPUNCH_BIND = %s", c.Bind)
	log.Tracef("HOLEPUNCH_PORT = %s", c.Port)
	log.Tracef("HOLEPUNCH_TOKEN = %s", "<redacted>")
	log.Tracef("HOLEPUNCH_SOCKET_PATH = %s", c.SocketPath)
	log.Tracef("HOLEPUNCH_ALLOW_ALL = %t", c.Authorization.AllowAll)
	log.Tracef("HOLEPUNCH_ALLOW_METRICS = %t", c.Authorization.AllowMetrics)
}

func Load() *Config {
	c := new(Config)
	c.Authorization = new(Firewall)

	_ = env.ParseOS(env.Schema{
		// setup
		"HOLEPUNCH_BIND":        env.StringOr(&c.Bind, "0.0.0.0"),
		"HOLEPUNCH_PORT":        env.StringOr(&c.Port, "6120"),
		"NOMAD_TOKEN":           env.Secret(&c.NomadToken, true),
		"HOLEPUNCH_SOCKET_PATH": env.StringOr(&c.SocketPath, defaultSocketPath()),

		// firewall
		"HOLEPUNCH_ALLOW_ALL":     env.BoolOr(&c.Authorization.AllowAll, false),
		"HOLEPUNCH_ALLOW_METRICS": env.BoolOr(&c.Authorization.AllowMetrics, true),
	})

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
