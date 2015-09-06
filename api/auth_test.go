package api

import (
	"testing"
	"net/http"
	"strings"
	"github.com/go-neutrino/neutrino-core/db"
)

func TestRegisterUser(t *testing.T) {
	body := register(t)

	res, err := neutrino.NewUsersDbService().FindId(body["email"], nil)

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

func TestAppRegisterUser(t *testing.T) {
	app := createApp(t)
	b := map[string]interface{}{
		"email": randomString() + "@gmail.com",
		"password": "pass",
	}

	r := sendAuthenticatedRequest("PUT", "/" + app.Id + "/auth", b, t)
	r.CodeIs(200)

	res, err := neutrino.NewAppUsersDbService(app.Id).FindId(b["email"], nil)

	if res == nil || err != nil {
		t.Fatal("User not created correctly", res, err);
	}
}

func TestAppLoginUser(t *testing.T) {
	app := createApp(t)
	email := randomString() + "@gmail.com"
	password := "pass"

	sendAuthenticatedRequest("PUT", "/" + app.Id + "/auth", map[string]interface{}{
		"email": email,
		"password": password,
	}, t)

	rec := sendAuthenticatedRequest("POST", "/" + app.Id + "/auth", map[string]interface{}{
		"email": email,
		"password": "pass",
	}, t)
	rec.CodeIs(200)

	var res map[string]interface{}
	rec.DecodeJsonPayload(&res)

	token := res["token"].(string)

	if len(token) <= 1 {
		t.Fatal("Incorrect token")
	}
}