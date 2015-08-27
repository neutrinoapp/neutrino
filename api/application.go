package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"realbase/core"
	"gopkg.in/mgo.v2/bson"
	"realbase/utils"
	"errors"
	"time"
)

type ApplicationModel struct {
	Id string `json: _id`
	Name string `json: "name"`
}

func GetAppFromRequest(r *rest.Request) (map[string]interface{}, error) {
	appId := r.PathParam("appId")

	if appId != "" {
		//TODO: cache this
		appDb := realbase.NewApplicationsDbService()
		return appDb.FindId(appId, nil)
	} else {
		return nil, errors.New("Invalid app id.")
	}
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
		"createdAt": time.Now(),
		"keys": bson.M{ //TODO:
			"Master Key": bson.M{
				"key": utils.GetCleanUUID(),
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

	RespondId(doc["_id"], w)
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
	app, err := GetAppFromRequest(r)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	w.WriteJson(app)
}