package main

import (
	"net/http"

	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/services/realtime/server"
)

func main() {
	err := server.Initialize()
	if err != nil {
		panic(err)
	}

	log.Info("Realtime service listening")
	log.Error(http.ListenAndServe(config.Get(config.KEY_REALTIME_PORT), nil))
}
