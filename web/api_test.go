package web

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shoenig/go-conceal"
	"github.com/shoenig/nomad-holepunch/configuration"
	"github.com/shoenig/test/must"
)

func socket() string {
	return filepath.Join(os.TempDir(), "hp.sock")
}

func upstream(t *testing.T) *http.Server {
	s := &http.Server{
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "ok", 200)
			},
		),
	}
	t.Cleanup(func() { _ = s.Close() })
	t.Cleanup(func() { _ = os.Remove(socket()) })

	listener, err := net.Listen("unix", socket())
	must.NoError(t, err)

	go func() { _ = s.Serve(listener) }()
	return s
}

func server(t *testing.T, c *configuration.Config) *http.Server {
	s := &http.Server{
		Addr:    "127.0.0.1:5431",
		Handler: newFirewall(c.Authorization, newProxy(c)),
	}
	t.Cleanup(func() { _ = s.Close() })

	go func() { _ = s.ListenAndServe() }()
	return s
}

func config(f *configuration.Firewall) *configuration.Config {
	return &configuration.Config{
		NomadToken:    conceal.New("abc123"),
		Bind:          "127.0.0.1",
		Port:          "5431",
		SocketPath:    socket(),
		Authorization: f,
	}
}

func Test_ServeHTTP(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		firewall *configuration.Firewall
		exp      int // status code
	}{
		{
			name: "defaults - metrics",
			path: "/v1/metrics",
			firewall: &configuration.Firewall{
				AllowMetrics: true,
			},
			exp: 200,
		},
		{
			name: "defaults - nodes",
			path: "/v1/nodes",
			firewall: &configuration.Firewall{
				AllowMetrics: true,
			},
			exp: 403,
		},
		{
			name: "defaults - not api",
			path: "/any",
			firewall: &configuration.Firewall{
				AllowMetrics: true,
			},
			exp: 403,
		},
		{
			name: "allow all - metrics",
			path: "/v1/metrics",
			firewall: &configuration.Firewall{
				AllowAll: true,
			},
			exp: 200,
		},
		{
			name: "allow all - nodes",
			path: "/v1/nodes",
			firewall: &configuration.Firewall{
				AllowAll: true,
			},
			exp: 200,
		},
		{
			name: "allow all - not api",
			path: "/any",
			firewall: &configuration.Firewall{
				AllowAll: true,
			},
			exp: 403, // only allows api access
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := config(tc.firewall)
			_ = upstream(t)
			_ = server(t, c)

			// wait for http servers to startup
			time.Sleep(150 * time.Millisecond)

			// do request and see what happens
			r, err := http.Get("http://127.0.0.1:5431" + tc.path)
			must.NoError(t, err)
			must.Eq(t, tc.exp, r.StatusCode)
		})
	}
}
