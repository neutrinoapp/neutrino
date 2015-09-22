package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino-core/db"
	"github.com/go-neutrino/neutrino-core/api"
	"github.com/go-neutrino/neutrino-config"
)

func main() {
	engine := gin.Default()

	c := nconfig.Load()

	db.Initialize(c)
	api.Initialize(engine, c)


	engine.Run(c.GetString(nconfig.KEY_CORE_PORT))
}
