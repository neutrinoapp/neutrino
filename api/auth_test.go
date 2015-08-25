package api

import (
	"testing"
	"net/http"
	"realbase/core"
	"strings"
)

func TestRegisterUser(t *testing.T) {
	body := register(t)

	res, err := realbase.NewUsersDbService().FindId(body.Username)

	if res == nil || err != nil {
		t.Fatal("User not created correctly", res, err);
	}
}

func TestLoginUser(t *testing.T) {
	body := register(t)

	rec := sendRequest("POST", "/auth", body, t)
	rec.CodeIs(http.StatusOK)

	bodyStr := rec.Recorder.Body.String()

	contains := strings.Contains(bodyStr, "token")

	if !contains {
		t.Fatal("Incorrect login response")
	}
}