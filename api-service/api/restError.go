package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/go-neutrino/neutrino-core/models"
)

func RestErrorInvalidBody(c *gin.Context) {
	RestError(c, "invalid request body")
}

func RestErrorNotFound(c *gin.Context) {
	RestError(c, "not found")
}

func RestErrorAppNotFound(c *gin.Context) {
	RestError(c, "app not found")
}

func RestError(c *gin.Context, err interface{}) {
	status := http.StatusInternalServerError

	var msg string
	switch t := err.(type) {
	case error:
		msg = t.Error()
	case string:
		msg = t
	}

	if msg == "not found" || msg == "app not found" {
		status = http.StatusNotFound
	} else if msg == "invalid request body" {
		status = http.StatusBadRequest
	} else if msg == "ns not found" {
		//mongo throws this error when a collection does not exist
		//but we call drop
		return
	}

	c.Error(errors.New(msg))

	c.JSON(status, models.JSON{
		"error": msg,
	})
}
