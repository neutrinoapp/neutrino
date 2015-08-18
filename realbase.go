package main

import (
	"github.com/labstack/echo"
	"realbase/core"
	"fmt"
	"realbase/api"
)

func main() {
	e := echo.New()

	realbase.Initialize("localhost:27017")
	api.Initialize(e)

	port := ":1234";

	fmt.Println("Listening on port", port)
	e.Run(port)
}
