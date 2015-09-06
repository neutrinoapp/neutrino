package api

import (
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
)
func RestErrorInvalidBody(c *gin.Context) {
	RestError(c, "invalid request body")
}

func RestErrorNotFound(c *gin.Context) {
	RestError(c, "not found")
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

	if msg == "not found" {
		status = http.StatusNotFound
	} else if msg == "invalid request body" {
		status = http.StatusBadRequest
	} else if msg == "ns not found" {
		//mongo throws this error when a collection does not exist
		//but we call drop
		return
	}

	c.AbortWithError(status, errors.New(msg))
}