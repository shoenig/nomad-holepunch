package configuration

import (
	"path/filepath"
	"testing"

	"github.com/shoenig/loggy"
	"github.com/shoenig/test/must"
)

func Test_defaultSocketPath(t *testing.T) {
	cases := []struct {
		name                   string
		envNomadSecretsDir     string
		envHolepunchSocketPath string
		exp                    string
	}{
		{
			name: "no env set",
			exp:  "/secrets/api.sock",
		},
		{
			name:               "only secrets dir set",
			envNomadSecretsDir: "/alloc/secrets",
			exp:                "/alloc/secrets/api.sock",
		},
		{
			name:                   "only socket path set",
			envHolepunchSocketPath: "/var/my.sock",
			exp:                    "/var/my.sock",
		},
		{
			name:                   "both env set",
			envNomadSecretsDir:     "/alloc/secrets",
			envHolepunchSocketPath: "/var/my.sock",
			exp:                    "/var/my.sock",
		},
	}

	log := loggy.New("test")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("NOMAD_TOKEN", "abc123")
			t.Setenv("NOMAD_SECRETS_DIR", tc.envNomadSecretsDir)
			t.Setenv("HOLEPUNCH_SOCKET_PATH", tc.envHolepunchSocketPath)
			c := Load(log)
			unixPath := filepath.ToSlash(c.SocketPath)
			must.Eq(t, tc.exp, unixPath)
		})
	}
}
