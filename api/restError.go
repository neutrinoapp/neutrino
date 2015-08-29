package api

import (
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"runtime/debug"
	"log"
)

type RestError struct {
	error
}

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

	if _, ok := falsePositiveStatusCodes[statusCode]; !ok {
		log.Println(message)
		debug.PrintStack()
	}
}

func RestErrorInvalidBody(w rest.ResponseWriter) {
	restError(w, http.StatusBadRequest, "Invalid request body.")
}

func RestGeneralError(w rest.ResponseWriter, e error) {
	status := http.StatusInternalServerError
	if e.Error() == "not found" {
		status = http.StatusNotFound
	}

	restError(w, status, e.Error())
}