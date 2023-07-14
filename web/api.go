package web

import (
	"io"
	"net"
	"net/http"

	"github.com/shoenig/go-conceal"
	"github.com/shoenig/loggy"
	"github.com/shoenig/nomad-holepunch/configuration"
)

type proxy struct {
	log        loggy.Logger
	httpClient *http.Client
	token      *conceal.Text
}

func newProxy(config *configuration.Config) http.Handler {
	return &proxy{
		log:   loggy.New("proxy"),
		token: config.NomadToken,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Dial: func(string, string) (net.Conn, error) {
					return net.Dial("unix", config.SocketPath)
				},
			},
		},
	}
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.log.Tracef("serving proxy endpoint to %s", r.RemoteAddr)

	request, err := p.toProxy(r)
	if err != nil {
		p.log.Errorf("failed to create request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := p.httpClient.Do(request)
	if err != nil {
		p.log.Errorf("failed to do request: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(response.StatusCode)
	_, _ = io.Copy(w, response.Body)
}

// toProxy creates a new http.Request to send over the Nomad identity
// unix domain socketj
func (p *proxy) toProxy(original *http.Request) (*http.Request, error) {
	method := original.Method
	url := "http://nomad" + original.URL.Path + "?" + original.URL.RawQuery
	request, err := http.NewRequest(method, url, nil)
	request.Header.Set("X-Nomad-Token", p.token.Unveil())
	return request, err
}
