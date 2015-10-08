package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino-core/api-service/api"
	"github.com/go-neutrino/neutrino-core/api-service/db"
	"github.com/go-neutrino/neutrino-core/config"
)

func main() {
	engine := gin.Default()

	db.Initialize()
	api.Initialize(engine)

	engine.Run(config.Get(config.KEY_API_PORT))
}
