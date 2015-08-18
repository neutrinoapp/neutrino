package api
import (
	"net/http"
	"github.com/labstack/echo"
)

type RestError struct {
	error
}

func restError(statusCode int, message string) error {
	return echo.NewHTTPError(statusCode, message)
}

func RestErrorInvalidBody() error {
	return restError(http.StatusBadRequest, "Invalid body.")
}