package api

import (
	"testing"
)

func createApp(t *testing.T) *ApplicationModel {
	app := &ApplicationModel{
		Name: randomString(),
	}

	rec := sendAuthenticatedRequest("POST", "/app", app, t)
	rec.CodeIs(200)

	res := rec.BodyJSON()
	app.Id = res["_id"].(string)

	return app
}

func TestCreateAndGetApplication(t *testing.T) {
	app := createApp(t)

	getRec := sendAuthenticatedRequest("GET", "/app", nil, t)
	getRec.CodeIs(200)

	var result []*ApplicationModel
	getRec.Decode(&result)

	if len(result) == 0 {
		t.Error("Application not created")
	}

	if result[0].Name != app.Name {
		t.Error("Application not created correctly")
	}
}

func TestDeleteApplication(t *testing.T) {
	app := createApp(t)

	sendAuthenticatedRequest("DELETE", "/app/"+app.Id, nil, t)

	getRec := sendAuthenticatedRequest("GET", "/app/"+app.Id, nil, t)
	result := getRec.BodyJSON()

	err := result["error"]

	if err != "app not found" {
		t.Fatal("App not deleted")
	}
}

func TestUpdateApplication(t *testing.T) {
	app := createApp(t)

	randomName := randomString() + "updated!"
	sendAuthenticatedRequest("PUT", "/app/"+app.Id, map[string]interface{}{
		"name": randomName,
	}, t)

	getRec := sendAuthenticatedRequest("GET", "/app/"+app.Id, nil, t)

	result := getRec.BodyJSON()

	res := result["name"]

	if res != randomName {
		t.Fatal("App not updated")
	}
}
