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
	"io/ioutil"
	"encoding/json"
	"fmt"
)

var (
	apiHandler http.Handler
	user map[string]interface{}
	token string
)

type ResRecorder struct {
	*httptest.ResponseRecorder
	t *testing.T
}

func (r *ResRecorder) CodeIs(s int) {
	if r.Code != s {
		r.t.Error(r.Code, "is different from", s)
	}
}

func (r *ResRecorder) B() string {
	return r.Body.String()
}

func (r *ResRecorder) BHas(str string) {
	if !strings.Contains(r.B(), str) {
		r.t.Error(r.B(), "does not contain", str)
	}
}

func (r *ResRecorder) BObj() JSON {
	b, _ := ioutil.ReadAll(r.Body)
	var res JSON
	json.Unmarshal(b, &res)
	return res
}

func (r *ResRecorder) Decode(o interface{}) {
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, o)
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

	var b string
	if body != nil {
		bodyStr, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		b = fmt.Sprintf("%s", bodyStr)
	}

	req, err := http.NewRequest(method, "/v1" + path, strings.NewReader(b))
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	apiHandler.ServeHTTP(w, req)
	return &ResRecorder{w, t}
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())

	return "r" + strconv.Itoa(rand.Int())
}

func register(t *testing.T) map[string]interface{} {
	b := JSON{
		"email": randomString() + "@gmail.com",
		"password": "pass",
	}

	rec := sendRequest("POST", "/register", b, t)

	rec.CodeIs(http.StatusOK)

	return b
}

func login(t *testing.T) (*UserModel, string) {
	if token == "" {
		user = register(t)
		rec := sendRequest("POST", "/login", JSON{
			"email": user["email"],
			"password": user["password"],
		}, t)

		response := rec.BObj()
		token = response["token"].(string)
	}

	return &UserModel{
		Email: user["email"].(string),
		Password: user["password"].(string),
	}, token
}