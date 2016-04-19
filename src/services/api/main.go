package main

import (
	"os"

	r "github.com/dancannon/gorethink"
	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/utils"
	"github.com/neutrinoapp/neutrino/src/services/api/api"
)

func main() {
	r.SetVerbose(true)

	defer utils.Recover()
	utils.ListenSignals()
	utils.Liveness()

	if os.Getenv("DEBUG_N") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			defer c.Next()
			log.Info(c.Request.Method, c.Request.URL.Path, c.Writer.Status())
		}
	}())

	api.Initialize(engine)

	engine.Run(config.Get(config.KEY_API_PORT))
}
