package api

import (
	"github.com/gin-gonic/gin"
	"testing"
	"net/http"
	"math/rand"
	"time"
	"strconv"
	"net/http/httptest"
	"strings"
)

var (
	apiHandler http.Handler
	user map[string]interface{}
	token string
)

type ResRecorder struct {
	rec *httptest.ResponseRecorder
	t *testing.T
}

func (r *ResRecorder) CodeIs(s int) {
	if r.rec.Code != s {
		r.t.Error(r.rec.Code, "is different from", s)
	}
}

func (r *ResRecorder) B() string {
	return r.rec.Body.String()
}

func (r *ResRecorder) BHas(str string) {
	if !strings.Contains(r.B(), str) {
		r.t.Error(r.B(), "does not contain", str)
	}
}

func sendAuthenticatedRequest(method, path string, body interface{}, t *testing.T) *ResRecorder {
	login(t)
	return sendRequest(method, path, body, t)
}

func sendRequest(method, path string, body interface{}, t *testing.T) *ResRecorder {
	if !IsInitialized() {
		e := gin.Default()
		apiHandler = e
		Initialize(e)
		httptest.NewServer(e)
	}

	req, err := http.NewRequest(method, "http://localhost" + path, nil)
	if err {
		panic(err)
	}

	w := httptest.NewRecorder()
	apiHandler.ServeHTTP(req, w)
	return &ResRecorder{w, t}
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