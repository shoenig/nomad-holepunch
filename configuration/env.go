package configuration

import (
	"os"
	"path/filepath"

	"github.com/shoenig/loggy"
)

type Config struct {
	Bind          string
	Port          string
	NomadToken    string
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
	return &Config{
		Bind:       get("HOLEPUNCH_BIND", "0.0.0.0"),
		Port:       get("HOLEPUNCH_PORT", "6120"),
		NomadToken: get("NOMAD_TOKEN", "unset"),
		SocketPath: get("HOLEPUNCH_SOCKET_PATH", defaultSocketPath()),
		Authorization: &Firewall{
			AllowAll:     getBool("HOLEPUNCH_ALLOW_ALL", false),
			AllowMetrics: getBool("HOLEPUNCH_ALLOW_METRICS", true),
		},
	}
}

func defaultSocketPath() string {
	dir := get("NOMAD_SECRETS_DIR", "/secrets")
	socket := filepath.Join(dir, "api.sock")
	return socket
}

func getBool(key string, fallback bool) bool {
	switch get(key, "") {
	case "":
		return fallback
	case "1", "true":
		return true
	default:
		return false
	}
}

func get(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
