package main

import (
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/utils"
	"github.com/neutrinoapp/neutrino/src/services/realtime/server"
)

func main() {
	r.SetVerbose(true)

	defer utils.Recover()
	utils.ListenSignals()
	utils.Liveness()

	err := server.Initialize()
	if err != nil {
		panic(err)
	}

	log.Info("Realtime service listening")
	log.Error(http.ListenAndServe(config.Get(config.KEY_REALTIME_PORT), nil))
}
