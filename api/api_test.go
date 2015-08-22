package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"testing"
	"net/http"
)

var apiHandler http.Handler

func SendRequest(method, path string, body interface{}, t *testing.T) *test.Recorded {
	if !IsInitialized() {
		restApi := rest.NewApi()
		Initialize(restApi)
		apiHandler = restApi.MakeHandler()
	}

	req := test.MakeSimpleRequest(method, "http://localhost" + path, body)
	return test.RunRequest(t, apiHandler, req)
}