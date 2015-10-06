package main

import (
	"github.com/go-neutrino/neutrino-config"
	"net/http"
	"github.com/go-neutrino/neutrino-core/log"
	"github.com/go-neutrino/neutrino-core/realtime-service/server"
)

func main() {
	c := nconfig.Load()

	server.Initialize(c)

	port := c.GetString(nconfig.KEY_REALTIME_PORT)
	log.Info("Listening on port: " + port)
	log.Info(http.ListenAndServe(port, nil))
}
