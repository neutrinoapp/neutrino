package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/src/common/config"
	"github.com/go-neutrino/neutrino/src/services/api/api"
	"github.com/go-neutrino/neutrino/src/services/api/db"
)

func main() {
	engine := gin.Default()

	db.Initialize()
	api.Initialize(engine)

	engine.Run(config.Get(config.KEY_API_PORT))
}
