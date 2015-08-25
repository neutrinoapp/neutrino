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
	user *UserModel
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
	return "r" + strconv.Itoa(rand.Int())
}

func register(t *testing.T) *UserModel {
	rand.Seed(time.Now().UnixNano())

	body := &UserModel{
		Username: randomString(),
		Password: "pass",
		Email: randomString() + "@gmail.com",
	}

	rec := sendRequest("PUT", "/auth", body, t)
	rec.CodeIs(http.StatusOK)

	return body
}

func login(t *testing.T) (*UserModel, string) {
	if token == "" {
		user = register(t)
	}

	rec := sendRequest("POST", "/auth", user, t)

	response := make(map[string]interface)
	rec.DecodeJsonPayload(&response)

	return user, response["token"]
}