package main

import (
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/services/realtime/server"
)

func main() {
	wsServer, err := server.Initialize()
	if err != nil {
		panic(err)
	}

	log.Info("Realtime service listening:", wsServer.Addr)
	log.Error(wsServer.ListenAndServe())
}
