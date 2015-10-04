package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino-config"
	"github.com/go-neutrino/neutrino-core/api-service/api"
	"github.com/go-neutrino/neutrino-core/api-service/db"
)

func main() {
	engine := gin.Default()

	c := nconfig.Load()

	db.Initialize(c)
	api.Initialize(engine, c)

	engine.Run(c.GetString(nconfig.KEY_CORE_PORT))
}
