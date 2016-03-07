package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
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
		})

		if err != nil || appMapDoc["user"] == nil {
			log.Error(RestError(c, err))
			return nil
		}

		u.Name = appMapDoc["user"].(string)
	}

	d := db.NewAppsDbService(u.Name)
	app, err := d.FindId(appId)
	if err != nil {
		log.Error(RestErrorAppNotFound(c))
		return nil
	}

	return models.JSON{}.FromMap(app)
}
