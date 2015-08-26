package api

import (
	"testing"
)

func TestCreateType(t *testing.T) {
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

	types := createdApp["types"].([]string)

	if types[1] != typeName {
		t.Error("Type not created correctly")
	}
}