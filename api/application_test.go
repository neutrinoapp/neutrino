package api

import (
	"testing"
)

func TestCreateApplication(t *testing.T) {
	rec := sendAuthenticatedRequest("POST", "application", ApplicationModel{randomString()}, t)
	rec.CodeIs(200)
}