package api

import (
	"net/http/httptest"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
)

func RequestPut(path string, body map[string]interface{}, handler (func(*echo.Context) error)) (int, string) {
	r := SendRequest(echo.PUT, path, body)
	return r.Code, r.Body.String()
}

func RequestPost(path string, body map[string]interface{}, handler (func(*echo.Context) error)) (int, string) {
	r := SendRequest(echo.POST, path, body)
	return r.Code, r.Body.String()
}

func RequestGet(path string, handler (func(*echo.Context) error)) (int, string) {
	r := SendRequest(echo.GET, path, make(map[string]interface{}))
	return r.Code, r.Body.String()
}

func RequestDelete(path string, handler (func(*echo.Context) error)) (int, string) {
	r := SendRequest(echo.DELETE, path, make(map[string]interface{}))
	return r.Code, r.Body.String()
}

func SendRequest(method, path string, body map[string]interface{}) (*httptest.ResponseRecorder) {
	e := echo.New()

	Initialize(e)

	bodyStr, _ := json.Marshal(body)
	bodyBuf := bytes.NewBuffer([]byte(bodyStr))

	req, _ := http.NewRequest(method, path, bodyBuf)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}