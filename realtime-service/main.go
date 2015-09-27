package main

import (
	"net/http"
	"fmt"
	"github.com/go-neutrino/neutrino-config"
	"github.com/go-neutrino/neutrino-core/realtime-service/server"
)

func main() {
	c := nconfig.Load()

	server.Initialize(c)

	fmt.Println(http.ListenAndServe(c.GetString(nconfig.KEY_REALTIME_PORT), nil))
}