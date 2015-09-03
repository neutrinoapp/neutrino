package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"testing"
	"net/http"
	"math/rand"
	"time"
	"strconv"
)

var (
	apiHandler http.Handler
	user map[string]interface{}
	token string
)

func sendAuthenticatedRequest(method, path string, body interface{}, t *testing.T) *test.Recorded {
	login(t)
	return sendRequest(method, path, body, t)
}

func sendRequest(method, path string, body interface{}, t *testing.T) *test.Recorded {
	if !IsInitialized() {
		restApi := rest.NewApi()
		Initialize(restApi)
		apiHandler = restApi.MakeHandler()
	}

	req := test.MakeSimpleRequest(method, "http://localhost" + path, body)

	if token != "" {
		req.Header.Add("Authorization", "Bearer " + token)
	}

	return test.RunRequest(t, apiHandler, req)
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())

	return "r" + strconv.Itoa(rand.Int())
}

func register(t *testing.T) map[string]interface{} {
	b := map[string]interface{}{
		"email": randomString() + "@gmail.com",
		"password": "pass",
	}

	rec := sendRequest("PUT", "/auth", b, t)

	rec.CodeIs(http.StatusOK)

	return b
}

func login(t *testing.T) (*UserModel, string) {
	if token == "" {
		user = register(t)
		rec := sendRequest("POST", "/auth", map[string]interface{}{
			"email": user["email"],
			"password": user["password"],
		}, t)

		var response map[string]interface{}
		rec.DecodeJsonPayload(&response)
		token = response["token"].(string)
	}

	return &UserModel{
		Email: user["email"].(string),
		Password: user["password"].(string),
	}, token
}