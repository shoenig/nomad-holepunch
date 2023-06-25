package configuration

import (
	"os"

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
	Disable      bool
	AllowMetrics bool
}

func (c *Config) Log(log loggy.Logger) {
	log.Tracef("HOLEPUNCH_BIND = %s", c.Bind)
	log.Tracef("HOLEPUNCH_PORT = %s", c.Port)
	log.Tracef("HOLEPUNCH_TOKEN = %s", "<redacted>")
	log.Tracef("socket path = %s", c.SocketPath)
	log.Tracef("HOLEPUNCH_ALLOW_ALL = %t", c.Authorization.Disable)
	log.Tracef("HOLEPUNCH_ALLOW_METRICS = %t", c.Authorization.AllowMetrics)
}

func Load() *Config {
	return &Config{
		Bind:       get("HOLEPUNCH_BIND", "0.0.0.0"),
		Port:       get("HOLEPUNCH_PORT", "3333"),
		NomadToken: get("NOMAD_TOKEN", "unset"),
		SocketPath: "secrets/api.sock",
		Authorization: &Firewall{
			Disable:      false,
			AllowMetrics: getBool("HOLEPUNCH_ALLOW_METRICS", true),
		},
	}
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
