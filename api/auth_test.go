package api

import (
	"testing"
	"net/http"
	"strings"
	"github.com/realbas3/realbas3/core"
)

func TestRegisterUser(t *testing.T) {
	body := register(t)

	res, err := realbase.NewUsersDbService().FindId(body.Email, nil)

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