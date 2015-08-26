package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"realbase/core"
	"gopkg.in/mgo.v2/bson"
)

type ApplicationModel struct {
	Id bson.ObjectId `json: _id`
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
		"keys": bson.M{ //TODO:
			"Master Key": bson.M{
				"key": GetCleanUUID(),
				"name": "Master Key",
				"permissions": bson.M{
					"types": bson.M{
						"read": true,
						"write": true,
					},
				},
			},
		},
	}

	if err := db.Insert(doc); err != nil {
		RestGeneralError(w, err)
		return
	}
}

func GetApplicationsHandler(w rest.ResponseWriter, r *rest.Request) {
	db := realbase.NewApplicationsDbService()

	res, err := db.Find(
		bson.M{
			"owner": r.Env["user"],
		},
		bson.M{
			"name": 1,
		},
	)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	w.WriteJson(res)
}

func GetApplicationHandler(w rest.ResponseWriter, r *rest.Request) {
	res, err := realbase.NewApplicationsDbService().FindId(r.PathParam("appId"), nil)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	w.WriteJson(res)
}