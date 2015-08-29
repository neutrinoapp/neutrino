package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"fmt"
	"net/http"
	"log"
	"github.com/go-realbase/realbase/core"
	"github.com/go-realbase/realbase/api"
)

func main() {
	restApi := rest.NewApi()

	realbase.Initialize("localhost:27017")
	api.Initialize(restApi)

	port := ":1234";

	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, restApi.MakeHandler()))
}
