package main

import (
	"fmt"
	"net/http"

	"github.com/shoenig/loggy"
	"github.com/shoenig/nomad-holepunch/configuration"
	"github.com/shoenig/nomad-holepunch/web"
)

func main() {
	log := loggy.New("main")
	log.Infof("^^ startup nomad-holepunch ^^")

	config := configuration.Load()
	config.Log(log)

	mux := web.New(config)
	address := fmt.Sprintf("%s:%s", config.Bind, config.Port)
	http.ListenAndServe(address, mux)

	select {
	// empty
	}
}
