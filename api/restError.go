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

func restError(w rest.ResponseWriter, statusCode int, message string) {
	rest.Error(w, message, statusCode)
	log.Println(message)
	debug.PrintStack()
}

func RestErrorInvalidBody(w rest.ResponseWriter) {
	restError(w, http.StatusBadRequest, "Invalid request body.")
}

func RestGeneralError(w rest.ResponseWriter, e error) {
	restError(w, http.StatusInternalServerError, e.Error())
}