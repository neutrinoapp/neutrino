package api

import (
	"testing"
	"net/http"
	"math/rand"
	"strconv"
	"realbase/core"
	"time"
)

func TestRegisterUser(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	body := map[string]interface{}{
		"username": "u" + strconv.Itoa(rand.Int()),
		"password": "pass",
		"email": "e" + strconv.Itoa(rand.Int()) + "@gmail.com",
	}

	statusCode, _ := RequestPut("/auth", body, RegisterUserHandler)
	if statusCode != http.StatusOK {
		t.Fatal("Wrong status code expected 200, got", statusCode);
	}

	res, err := realbase.NewUsersDbService().FindId(body["username"])

	if res == nil || err != nil {
		t.Fatal("User not created correctly", res, err);
	}
}