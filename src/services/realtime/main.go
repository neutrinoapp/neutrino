package main

import (
	"net/http"

	"github.com/go-neutrino/neutrino/src/common/config"
	"github.com/go-neutrino/neutrino/src/common/log"
	"github.com/go-neutrino/neutrino/src/services/realtime/server"
)

func main() {
	server.Initialize()

	port := config.Get(config.KEY_REALTIME_PORT)
	log.Info("Listening on port: " + port)
	log.Info(http.ListenAndServe(port, nil))
}
