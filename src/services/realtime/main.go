package main

import (
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/services/realtime/server"
)

func main() {
	wsServer := server.Initialize()

	log.Info("Realtime service listening:", wsServer.Addr)
	log.Error(wsServer.ListenAndServe())
}
