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

func TestDeleteApplication(t *testing.T) {
	app := createApp(t);

	delRec := sendAuthenticatedRequest("DELETE", "/applications/" + app.Id, nil, t)
	delRec.CodeIs(200)

	getRec := sendAuthenticatedRequest("GET", "/applications/" + app.Id, nil, t)
	getRec.CodeIs(404)

	var result map[string]interface{}
	getRec.DecodeJsonPayload(&result)

	err := result["error"]

	if err != "not found" {
		t.Fatal("App not deleted")
	}
}

func TestUpdateApplication(t *testing.T) {
	app := createApp(t);

	randomName := randomString() + "updated!"
	putRec := sendAuthenticatedRequest("PUT", "/applications/" + app.Id, map[string]interface{}{
		"name": randomName,
	}, t)
	putRec.CodeIs(200)

	getRec := sendAuthenticatedRequest("GET", "/applications/" + app.Id, nil, t)
	getRec.CodeIs(200)

	var result map[string]interface{}
	getRec.DecodeJsonPayload(&result)

	res := result["name"]

	if res != randomName {
		t.Fatal("App not updated")
	}
}