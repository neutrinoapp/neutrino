package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"realbase/core"
	"fmt"
	"realbase/api"
	"net/http"
	"log"
)

func main() {
	restApi := rest.NewApi()

	realbase.Initialize("localhost:27017")
	api.Initialize(restApi)

	port := ":1234";

	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, restApi.MakeHandler()))
}
