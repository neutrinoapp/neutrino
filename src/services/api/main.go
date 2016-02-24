package main

import (
	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/services/api/api"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
)

func main() {
	engine := gin.Default()

	db.Initialize()
	api.Initialize(engine)

	engine.Run(config.Get(config.KEY_API_PORT))
}
