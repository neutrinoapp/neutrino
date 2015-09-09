package api

import (
	"testing"
)

func createApp(t *testing.T) *ApplicationModel {
	app := &ApplicationModel{
		Name: randomString(),
	}

	rec := sendAuthenticatedRequest("POST", "/app/", app, t)
	rec.CodeIs(200)

	res := rec.BObj()
	app.Id = res["_id"].(string)

	return app
}

func TestCreateAndGetApplication(t *testing.T) {
	app := createApp(t)

	getRec := sendAuthenticatedRequest("GET", "/app/", nil, t)
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
	app := createApp(t);

	delRec := sendAuthenticatedRequest("DELETE", "/app/" + app.Id, nil, t)
	delRec.CodeIs(200)

	getRec := sendAuthenticatedRequest("GET", "/app/" + app.Id, nil, t)
	getRec.CodeIs(404)

	result := getRec.BObj()

	err := result["error"]

	if err != "not found" {
		t.Fatal("App not deleted")
	}
}

func TestUpdateApplication(t *testing.T) {
	app := createApp(t);

	randomName := randomString() + "updated!"
	putRec := sendAuthenticatedRequest("PUT", "/app/" + app.Id, map[string]interface{}{
		"name": randomName,
	}, t)
	putRec.CodeIs(200)

	getRec := sendAuthenticatedRequest("GET", "/app/" + app.Id, nil, t)
	getRec.CodeIs(200)

	result := getRec.BObj()

	res := result["name"]

	if res != randomName {
		t.Fatal("App not updated")
	}
}