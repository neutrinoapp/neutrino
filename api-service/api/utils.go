package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/api-service/db"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"github.com/go-neutrino/neutrino/utils"
)

func ApiUser(c *gin.Context) *apiUser {
	val, exists := c.Get("user")
	if exists {
		return val.(*apiUser)
	}

	return &apiUser{}
}

func Application(c *gin.Context, appId string) models.JSON {
	//TODO: cache all this
	u := ApiUser(c)
	userExists := u != nil
	p := utils.PathOfUrl(c.Request.URL.String())
	if !userExists && p != "login" && p != "register" {
		//TODO: handle non authorized data access - anonymous
		log.Error(RestErrorUnauthorized(c))
		return nil
	}

	if userExists {
		//check if the user is inApp (not the owner of the app)
		//if it is, we need to find the app by id
		isInAppUser := u.InApp
		if isInAppUser {
			userExists = false
		}
	}

	if !userExists {
		u = &apiUser{}
		d := db.NewAppsMapDbService()
		appMapDoc, err := d.FindOne(models.JSON{
			"appId": appId,
		}, nil)

		if err != nil || appMapDoc["user"] == nil {
			log.Error(RestError(c, err))
			return nil
		}

		u.Name = appMapDoc["user"].(string)
	}

	d := db.NewAppsDbService(u.Name)
	app, err := d.FindId(appId, nil)
	if err != nil {
		log.Error(RestErrorAppNotFound(c))
		return nil
	}

	return models.JSON{}.FromMap(app)
}
