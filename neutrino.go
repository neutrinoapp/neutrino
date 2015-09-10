package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino-core/db"
	"github.com/go-neutrino/neutrino-core/api"
	"github.com/go-neutrino/go-env-config"
)

func main() {
	engine := gin.Default()

	c, err := envconfig.LoadSimple("neutrino")

	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	db.Initialize(c)
	api.Initialize(engine, c)

	engine.Run(c["port"].(string))
}
