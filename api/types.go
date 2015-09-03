package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/go-neutrino/neutrino-core/core"
)

type TypesController struct {
}

func (t *TypesController) Path() string {
	return "/types"
}

func (t *TypesController) CreateTypeHandler(w rest.ResponseWriter, r *rest.Request) {
	var body map[string]string
	r.DecodeJsonPayload(&body)
	typeName := body["name"]

	app, err := GetAppFromRequest(r)

	if err != nil {
		RestError(w, err)
		return
	}

	appsDb := neutrino.NewAppsDbService(r.Env["user"].(string))
	appsDb.UpdateId(app["_id"],
		bson.M{
			"$push": bson.M{
				"types": typeName,
			},
		},
	)

	w.WriteHeader(http.StatusOK)
}

func (t * TypesController) DeleteType(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")

	app, err := GetAppFromRequest(r)

	if err != nil {
		RestError(w, err)
		return
	}

	appsDb := neutrino.NewAppsDbService(r.Env["user"].(string))
	appsDb.UpdateId(app["_id"],
		bson.M{
			"$pull": bson.M{
				"types": typeName,
			},
		},
	)

	db := neutrino.NewTypeDbService(appId, typeName)
	session, collection := db.GetCollection()
	defer session.Close()

	dropError := collection.DropCollection()

	if dropError != nil {
		RestError(w, dropError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *TypesController) InsertInTypeHandler(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")
	var body map[string]interface{}
	r.DecodeJsonPayload(&body)

	db := neutrino.NewTypeDbService(appId, typeName)
	err := db.Insert(body)

	if err != nil {
		RestError(w, err)
		return
	}

	RespondId(body["_id"], w)
}

func (t *TypesController) GetTypeDataHandler(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")

	app, appErr := GetAppFromRequest(r)
	if appErr != nil {
		RestError(w, appErr)
		return
	}

	types := app["types"].([]interface{})

	found := false

	for _, t := range types {
		if value, ok := t.(string); ok && value == typeName {
			found = true
			break
		}
	}

	if !found {
		RestErrorNotFound(w)
		return
	}

	db := neutrino.NewTypeDbService(appId, typeName)

	typeData, err := db.Find(nil, nil)

	if err != nil {
		RestError(w, err)
		return
	}

	w.WriteJson(typeData)
}

func (t *TypesController) GetTypeItemById(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")
	itemId := r.PathParam("itemId")

	db := neutrino.NewTypeDbService(appId, typeName)

	item, err := db.FindId(itemId, nil)

	if (err != nil) {
		RestError(w, err)
		return
	}

	w.WriteJson(item)
}

func (t *TypesController) UpdateTypeItemById(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")
	itemId := r.PathParam("itemId")

	db := neutrino.NewTypeDbService(appId, typeName)

	var body bson.M
	r.DecodeJsonPayload(&body)

	err := db.UpdateId(itemId, body)

	if err != nil {
		RestError(w, err)
		return
	}

	w.WriteHeader(200)
}

func (t *TypesController) DeleteTypeItemById(w rest.ResponseWriter, r *rest.Request) {
	appId := r.PathParam("appId")
	typeName := r.PathParam("typeName")
	itemId := r.PathParam("itemId")

	db := neutrino.NewTypeDbService(appId, typeName)

	err := db.RemoveId(itemId)

	if err != nil {
		RestError(w, err)
		return
	}

	w.WriteHeader(200)
}