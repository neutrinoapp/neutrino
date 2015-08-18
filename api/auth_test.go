package api

import (
	"testing"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"bytes"
	"math/rand"
	"strconv"
	"encoding/json"
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

func RequestPut(path string, body map[string]interface{}, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Put(path, handler)
	r := request(echo.PUT, path, body, e)
	return r.Code, r.Body.String()
}

func RequestPost(path string, body map[string]interface{}, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Post(path, handler)
	r := request(echo.POST, path, body, e)
	return r.Code, r.Body.String()
}

func RequestGet(path string, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Get(path, handler)
	r := request(echo.GET, path, make(map[string]interface{}), e)
	return r.Code, r.Body.String()
}

func RequestDelete(path string, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Delete(path, handler)
	r := request(echo.DELETE, path, make(map[string]interface{}), e)
	return r.Code, r.Body.String()
}

func request(method, path string, body map[string]interface{}, e *echo.Echo) (*httptest.ResponseRecorder) {
	Initialize(e)

	bodyStr, _ := json.Marshal(body)
	bodyBuf := bytes.NewBuffer([]byte(bodyStr))

	req, _ := http.NewRequest(method, path, bodyBuf)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}