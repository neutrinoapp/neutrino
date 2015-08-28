package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/realbas3/realbas3/core"
)

func CreateTypeHandler(w rest.ResponseWriter, r *rest.Request) {
	var body map[string]string
	r.DecodeJsonPayload(&body)
	typeName := body["name"]

	app, err := GetAppFromRequest(r)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	appsDb := realbase.NewApplicationsDbService()
	appsDb.Update(
		bson.M{
			"_id": app["_id"],
		},
		bson.M{
			"$push": bson.M{
				"types": typeName,
			},
		},
	)

	w.WriteHeader(http.StatusOK)
}

func InsertInTypeHandler(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")
	var body map[string]interface{}
	r.DecodeJsonPayload(&body)

	db := realbase.NewTypeDbService(appId, typeName)
	err := db.Insert(body)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	RespondId(body["_id"], w)
}

func GetTypeDataHandler(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")

	db := realbase.NewTypeDbService(appId, typeName)

	typeData, err := db.Find(nil, nil)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	w.WriteJson(typeData)
}

func GetTypeItemById(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")
	itemId := r.PathParam("itemId")

	db := realbase.NewTypeDbService(appId, typeName)

	item, err := db.FindId(itemId, nil)

	if (err != nil) {
		RestGeneralError(w, err)
		return
	}

	w.WriteJson(item)
}