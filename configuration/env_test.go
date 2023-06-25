package configuration

import (
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_defaultSocketPath(t *testing.T) {
	cases := []struct {
		name                   string
		envNomadSecretsDir     string
		envHolepunchSocketPath string
		exp                    string
	}{{
		name: "no env set",
		exp:  "/secrets/api.sock",
	}, {
		name:               "only secrets dir set",
		envNomadSecretsDir: "/alloc/secrets",
		exp:                "/alloc/secrets/api.sock",
	}, {
		name:                   "only socket path set",
		envHolepunchSocketPath: "/var/my.sock",
		exp:                    "/var/my.sock",
	}, {
		name:                   "both env set",
		envNomadSecretsDir:     "/alloc/secrets",
		envHolepunchSocketPath: "/var/my.sock",
		exp:                    "/var/my.sock",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("NOMAD_SECRETS_DIR", tc.envNomadSecretsDir)
			t.Setenv("HOLEPUNCH_SOCKET_PATH", tc.envHolepunchSocketPath)
			c := Load()
			unixPath := filepath.ToSlash(c.SocketPath)
			must.Eq(t, tc.exp, unixPath)
		})
	}
}

func Test_getBool(t *testing.T) {
	cases := []struct {
		name     string
		env      string
		fallback bool
		exp      bool
	}{
		{
			name:     "unset and true",
			env:      "",
			fallback: true,
			exp:      true,
		},
		{
			name:     "1 and false",
			env:      "1",
			fallback: false,
			exp:      true,
		},
		{
			name:     "true and false",
			env:      "true",
			fallback: false,
			exp:      true,
		},
		{
			name:     "bleh and true",
			env:      "bleh",
			fallback: true,
			exp:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TESTVAR", tc.env)
			result := getBool("TESTVAR", tc.fallback)
			must.Eq(t, tc.exp, result)
		})
	}
}

func Test_get(t *testing.T) {
	cases := []struct {
		name     string
		env      string
		fallback string
		exp      string
	}{
		{
			name:     "unset",
			env:      "",
			fallback: "default",
			exp:      "default",
		},
		{
			name:     "set",
			env:      "blah",
			fallback: "default",
			exp:      "blah",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TESTVAR", tc.env)
			result := get("TESTVAR", tc.fallback)
			must.Eq(t, tc.exp, result)
		})
	}
}
