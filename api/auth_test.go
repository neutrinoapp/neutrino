package api_test

import (
	"testing"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
)

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	statusCode, _ := request(echo.POST, "/auth", e)

	if statusCode != http.StatusOK {
		t.Fatal("Wrong status code expected 200, got", statusCode);
	}
}

func request(method, path string, e *echo.Echo) (int, string) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}