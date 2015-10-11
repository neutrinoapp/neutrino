package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/api-service/api"
	"github.com/go-neutrino/neutrino/api-service/db"
	"github.com/go-neutrino/neutrino/config"
)

func main() {
	engine := gin.Default()

	db.Initialize()
	api.Initialize(engine)

	engine.Run(config.Get(config.KEY_API_PORT))
}
