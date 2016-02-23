package api

import (
	"errors"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/src/common/models"
)

type restError struct {
	error
	Message string
	Code    int
}

func (e restError) Error() string {
	return e.Message
}

func (e restError) String() string {
	return fmt.Sprintf("Message: %s, Code: %v", e.Message, e.Code)
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

	setStatus := func(msg string) {
		if msg == "not found" || msg == "app not found" {
			status = http.StatusNotFound
		} else if msg == "invalid request body" {
			status = http.StatusBadRequest
		} else if msg == "not authorized" {
			status = http.StatusUnauthorized
		}
	}

	setMessage := func(status int) {
		if status == http.StatusNotFound {
			msg = "not found"
		} else if status == http.StatusBadRequest {
			msg = "invalid request body"
		} else if status == http.StatusUnauthorized {
			msg = "not authorized"
		}
	}

	switch t := err.(type) {
	case error:
		msg = t.Error()
		setStatus(msg)
	case string:
		msg = t
		setStatus(msg)
	case int:
		status = t
		setMessage(status)
	}

	return restError{
		Message: msg,
		Code:    status,
	}
}

func RestError(c *gin.Context, err interface{}) error {
	restError := BuildError(err)

	c.JSON(restError.Code, models.JSON{"error": restError.Message})
	c.Abort()

	return errors.New(restError.String())
}
