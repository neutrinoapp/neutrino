package api

import (
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"runtime/debug"
	"log"
	"os"
	"errors"
)

var falsePositiveStatusCodes map[int]bool

func init() {
	falsePositiveStatusCodes = map[int]bool{
		http.StatusNotFound: true,
	}
}

func restError(w rest.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	err := w.WriteJson(map[string]interface{}{
		"error": message,
	})

	if err != nil {
		panic(err)
	}

	if _, ok := falsePositiveStatusCodes[statusCode]; !ok || os.Getenv("PROD") != "TRUE" {
		log.Println(message)
		debug.PrintStack()
	}
}

func RestErrorInvalidBody(w rest.ResponseWriter) {
	RestError(w, errors.New("invalid request body"))
}

func RestErrorNotFound(w rest.ResponseWriter) {
	RestError(w, errors.New("not found"))
}

func RestError(w rest.ResponseWriter, e error) {
	status := http.StatusInternalServerError
	errMsg := e.Error()
	if errMsg == "not found" {
		status = http.StatusNotFound
	} else if errMsg == "invalid request body" {
		status = http.StatusBadRequest
	} else if errMsg == "ns not found" {
		//mongo throws this error when a collection does not exist
		//but we call drop
		return
	}

	restError(w, status, e.Error())
}