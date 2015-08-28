package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"fmt"
	"net/http"
	"log"
	"github.com/realbas3/realbas3/core"
	"github.com/realbas3/realbas3/api"
)

func main() {
	restApi := rest.NewApi()

	realbase.Initialize("localhost:27017")
	api.Initialize(restApi)

	port := ":1234";

	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, restApi.MakeHandler()))
}
