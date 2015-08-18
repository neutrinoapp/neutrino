package api

import (
	"net/http"
	"encoding/json"
)

func JsonBody(r *http.Request) (map[string]interface{}, error) {
	var t map[string]interface{}

	if r.Body == nil {
		return t, RestErrorInvalidBody()
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		return t, err
	}

	return t, nil
}