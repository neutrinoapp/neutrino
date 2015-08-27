package api

import (
	"testing"
)

func createApp(t *testing.T) *ApplicationModel {
	app := &ApplicationModel{
		Name: randomString(),
	}

	rec := sendAuthenticatedRequest("POST", "/applications", app, t)
	rec.CodeIs(200)

	var res map[string]interface{}
	rec.DecodeJsonPayload(&res)

	app.Id = res["_id"].(string)

	return app
}

func TestCreateAndGetApplication(t *testing.T) {
	app := createApp(t)

	getRec := sendAuthenticatedRequest("GET", "/applications", nil, t)
	getRec.CodeIs(200)

	var result []*ApplicationModel
	getRec.DecodeJsonPayload(&result)


	if len(result) == 0 {
		t.Error("Application not created")
	}

	if result[0].Name != app.Name {
		t.Error("Application not created correctly")
	}
}
