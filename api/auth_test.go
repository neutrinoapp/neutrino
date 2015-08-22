package api

import (
	"testing"
	"net/http"
	"math/rand"
	"strconv"
	"realbase/core"
	"time"
	"strings"
)

func registerUser(t *testing.T) *UserModel {
	rand.Seed(time.Now().UnixNano())

	body := &UserModel{
		Username: "u" + strconv.Itoa(rand.Int()),
		Password: "pass",
		Email: "e" + strconv.Itoa(rand.Int()) + "@gmail.com",
	}

	rec := SendRequest("PUT", "/auth", body, t)
	rec.CodeIs(http.StatusOK)

	return body
}

func TestRegisterUser(t *testing.T) {
	body := registerUser(t)

	res, err := realbase.NewUsersDbService().FindId(body.Username)

	if res == nil || err != nil {
		t.Fatal("User not created correctly", res, err);
	}
}

func TestLoginUser(t *testing.T) {
	body := registerUser(t)

	rec := SendRequest("POST", "/auth", body, t)
	rec.CodeIs(http.StatusOK)

	bodyStr := rec.Recorder.Body.String()

	contains := strings.Contains(bodyStr, "token")

	if !contains {
		t.Fatal("Incorrect login response")
	}
}