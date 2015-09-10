package api

import (
	"testing"
	"net/http"
	"strings"
	"github.com/go-neutrino/neutrino-core/db"
)

func TestRegisterUser(t *testing.T) {
	body := register(t)

	res, err := db.NewUsersDbService().FindId(body["email"], nil)

	if res == nil || err != nil {
		t.Fatal("User not created correctly", res, err);
	}
}

func TestLoginUser(t *testing.T) {
	body := register(t)

	rec := sendRequest("POST", "/login", body, t)
	rec.CodeIs(http.StatusOK)

	bodyStr := rec.B()

	contains := strings.Contains(bodyStr, "token")

	if !contains {
		t.Fatal("Incorrect login response")
	}
}

func TestAppRegisterUser(t *testing.T) {
	app := createApp(t)
	b := JSON{
		"email": randomString() + "@gmail.com",
		"password": "pass",
	}

	r := sendAuthenticatedRequest("POST", "/app/" + app.Id + "/register", b, t)
	r.CodeIs(200)

	res, err := db.NewAppUsersDbService(app.Id).FindId(b["email"], nil)

	if res == nil || err != nil {
		t.Fatal("User not created correctly", res, err);
	}
}

func TestAppLoginUser(t *testing.T) {
	app := createApp(t)
	email := randomString() + "@gmail.com"
	password := "pass"

	sendAuthenticatedRequest("POST", "/app/" + app.Id + "/register", JSON{
		"email": email,
		"password": password,
	}, t)

	rec := sendAuthenticatedRequest("POST", "/app/" + app.Id + "/login", JSON{
		"email": email,
		"password": "pass",
	}, t)
	rec.CodeIs(200)

	res := rec.BObj()

	token := res["token"].(string)

	if len(token) <= 1 {
		t.Fatal("Incorrect token")
	}
}