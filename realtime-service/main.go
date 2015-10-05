package main

import (
	"fmt"
	"github.com/go-neutrino/neutrino-config"
	"github.com/go-neutrino/neutrino-core/realtime-service/server"
	"net/http"
)

func main() {
	c := nconfig.Load()

	server.Initialize(c)

	port := c.GetString(nconfig.KEY_REALTIME_PORT)
	fmt.Println("Listening on port: " + port)
	fmt.Println(http.ListenAndServe(port, nil))
}
