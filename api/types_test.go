package api

import (
	"testing"
)

func setupTypeTests(t *testing.T) (map[string]interface{}, *ApplicationModel, string) {
	typeName := randomString()
	app := createApp(t)

	createTypeRec := sendAuthenticatedRequest("POST", "/" + app.Id + "/types",
		map[string]interface{}{
			"name": typeName,
		}, t)

	createTypeRec.CodeIs(200)

	getRec := sendAuthenticatedRequest("GET", "/applications/" + app.Id, nil, t)
	getRec.CodeIs(200)

	var createdApp map[string]interface{}
	getRec.DecodeJsonPayload(&createdApp)

	return createdApp, app, typeName
}

func TestCreateType(t *testing.T) {
	createdApp, _, typeName := setupTypeTests(t)

	types := createdApp["types"].([]interface{})

	if types[1].(string) != typeName {
		t.Error("Type not created correctly")
	}
}

func TestDeleteType(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	deleteReq := sendAuthenticatedRequest("DELETE", "/" + app.Id + "/types/" + typeName, nil, t)
	deleteReq.CodeIs(200)

	getReq := sendAuthenticatedRequest("GET", "/" + app.Id + "/types/" + typeName, nil, t)
	getReq.CodeIs(404)

	appReq := sendAuthenticatedRequest("GET", "/applications/" + app.Id, nil, t)
	var updatedApp map[string]interface{}
	appReq.DecodeJsonPayload(&updatedApp)

	types := updatedApp["types"].([]interface{})

	if len(types) > 1 {
		t.Error("Type not deleted correctly")
	}
}

func TestGetAndInsertTypeData(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	sendAuthenticatedRequest("POST", "/" + app.Id + "/types/" + typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)

	getRec := sendAuthenticatedRequest("GET", "/" + app.Id + "/types/" + typeName, nil, t)
	getRec.CodeIs(200)

	var data []map[string]interface{}
	getRec.DecodeJsonPayload(&data)

	record := data[0]

	if record["field1"] != "test" || record["field2"] != "test" {
		t.Error("Item not written correctly")
	}
}

func TestGetByIdTypeData(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	rec := sendAuthenticatedRequest("POST", "/" + app.Id + "/types/" + typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)

	var res map[string]interface{}
	rec.DecodeJsonPayload(&res)
	id := res["_id"].(string)

	rec1 := sendAuthenticatedRequest("GET", "/" + app.Id + "/types/" + typeName + "/" + id, nil, t)
	var item map[string]interface{}
	rec1.DecodeJsonPayload(&item)

	if item["field1"] != "test" || item["field2"] != "test" {
		t.Error("Item not written correctly")
	}
}

func TestUpdateTypeItemById(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	rec := sendAuthenticatedRequest("POST", "/" + app.Id + "/types/" + typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)
	rec.CodeIs(200)

	var res map[string]interface{}
	rec.DecodeJsonPayload(&res)
	id := res["_id"].(string)


	sendAuthenticatedRequest("PUT", "/" + app.Id + "/types/" + typeName + "/" + id, map[string]interface{}{
		"field1": "testupdated",
		"field2": "testupdated",
	}, t)

	rec1 := sendAuthenticatedRequest("GET", "/" + app.Id + "/types/" + typeName + "/" + id, nil, t)
	var item map[string]interface{}
	rec1.DecodeJsonPayload(&item)

	if item["field1"] != "testupdated" || item["field2"] != "testupdated" {
		t.Fatal("Item not updated correctly")
	}
}

func TestDeleteTypeItemById(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	rec := sendAuthenticatedRequest("POST", "/" + app.Id + "/types/" + typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)
	rec.CodeIs(200)

	var res map[string]interface{}
	rec.DecodeJsonPayload(&res)
	id := res["_id"].(string)

	sendAuthenticatedRequest("DELETE", "/" + app.Id + "/types/" + typeName + "/" + id, nil, t)

	rec1 := sendAuthenticatedRequest("GET", "/" + app.Id + "/types/" + typeName + "/" + id, nil, t)
	rec1.CodeIs(404)
}