package api

import (
	"net/http"
	"encoding/json"
)

func JsonBody(r *http.Request) map[string]interface{} {
	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	return t
}