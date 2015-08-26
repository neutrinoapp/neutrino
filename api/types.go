package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"realbase/core"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func CreateTypeHandler(w rest.ResponseWriter, r *rest.Request) {
	var body map[string]string
	r.DecodeJsonPayload(&body)
	typeName := body["name"]

	app := r.Env["app"].(ApplicationModel)

	appsDb := realbase.NewApplicationsDbService()
	appsDb.Update(
		bson.M{
			"_id": app.Id,
		},
		bson.M{
			"$push": bson.M{
				"types": typeName,
			},
		},
	)

	w.WriteHeader(http.StatusOK)
}