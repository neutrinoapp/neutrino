package main

import (
	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/services/api/api"
)

func main() {
	engine := gin.Default()

	api.Initialize(engine)

	engine.Run(config.Get(config.KEY_API_PORT))
}
