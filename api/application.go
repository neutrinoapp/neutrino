package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"realbase/core"
	"gopkg.in/mgo.v2/bson"
)

type ApplicationModel struct {
	Name string `json: "name"`
}

func CreateApplicationHandler(w rest.ResponseWriter, r *rest.Request) {
	body := &ApplicationModel{}

	if err := r.DecodeJsonPayload(body); err != nil {
		RestGeneralError(w, err)
		return
	}

	if body.Name == "" {
		RestErrorInvalidBody(w)
		return
	}

	db := realbase.NewApplicationsDbService()

	username := r.Env["user"].(string)
	doc := bson.M{
		"name": body.Name,
		"owner": username,
		"types": []string{"users"},
		"keys": []string{},//TODO:
	}

	if err := db.Insert(doc); err != nil {
		RestGeneralError(w, err)
		return
	}
}