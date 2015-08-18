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
)

func TestRegisterUser(t *testing.T) {
	body := map[string]interface{}{
		"username": "u" + strconv.Itoa(rand.Int()),
		"password": "pass",
		"email": "e" + strconv.Itoa(rand.Int()) + "@gmail.com",
	}

	statusCode, _ := requestPut("/auth", body, RegisterUserHandler)
	if statusCode != http.StatusOK {
		t.Fatal("Wrong status code expected 200, got", statusCode);
	}
}

func requestPut(path string, body map[string]interface{}, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Put(path, handler)
	r := request(echo.PUT, path, body, e)
	return r.Code, r.Body.String()
}

func requestPost(path string, body map[string]interface{}, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Post(path, handler)
	r := request(echo.POST, path, body, e)
	return r.Code, r.Body.String()
}

func requestGet(path string, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Get(path, handler)
	r := request(echo.GET, path, make(map[string]interface{}), e)
	return r.Code, r.Body.String()
}

func requestDelete(path string, handler (func(*echo.Context) error)) (int, string) {
	e := echo.New()
	e.Delete(path, handler)
	r := request(echo.DELETE, path, make(map[string]interface{}), e)
	return r.Code, r.Body.String()
}

func request(method, path string, body map[string]interface{}, e *echo.Echo) (*httptest.ResponseRecorder) {
	bodyStr, _ := json.Marshal(body)
	bodyBuf := bytes.NewBuffer([]byte(bodyStr))

	req, _ := http.NewRequest(method, path, bodyBuf)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}