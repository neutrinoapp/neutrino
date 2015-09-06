package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino-core/db"
	"github.com/go-neutrino/neutrino-core/api"
)

func main() {
	engine := gin.Default()

	db.Initialize("localhost:27017")
	api.Initialize(engine)

	port := ":1234";
	engine.Run(port)
}
