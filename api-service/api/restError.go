package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/go-neutrino/neutrino/models"
	"errors"
)

type restError struct {
	error
	Message string
	Code int
}

func (e restError) Error() string {
	return e.Message
}

func RestErrorInvalidBody(c *gin.Context) error {
	return RestError(c, "invalid request body")
}

func RestErrorNotFound(c *gin.Context) error {
	return RestError(c, "not found")
}

func RestErrorAppNotFound(c *gin.Context) error {
	return RestError(c, "app not found")
}

func RestErrorUnauthorized(c *gin.Context) error {
	return RestError(c, "not authorized")
}

func BuildError(err interface{}) restError {
	status := http.StatusInternalServerError

	var msg string
	switch t := err.(type) {
	case error:
		msg = t.Error()
	case string:
		msg = t
		if msg == "not found" || msg == "app not found" {
			status = http.StatusNotFound
		} else if msg == "invalid request body" {
			status = http.StatusBadRequest
		} else if msg == "not authorized" {
			status = http.StatusUnauthorized
		}
	case int:
		status = t
		if status == http.StatusNotFound {
			msg = "not found"
		} else if status == http.StatusBadRequest {
			msg = "invalid request body"
		} else if status == http.StatusUnauthorized {
			msg = "not authorized"
		}
	}

	return restError{
		Message: msg,
		Code: status,
	}
}

func RestError(c *gin.Context, err interface{}) error {
	restError := BuildError(err)

	c.JSON(restError.Code, models.JSON{"error": restError.Message});
	c.Abort()

	return errors.New(restError.Message)
}
