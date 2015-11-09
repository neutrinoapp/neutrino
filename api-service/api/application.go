package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/api-service/db"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"github.com/go-neutrino/neutrino/utils"
	"github.com/go-neutrino/neutrino/utils/webUtils"
	"net/http"
	"time"
"strings"
)

type ApplicationModel struct {
	Id   string `json: _id`
	Name string `json: "name"`
}

type ApplicationController struct {
}

func (a *ApplicationController) CreateApplicationHandler(c *gin.Context) {
	body := &ApplicationModel{}

	if err := c.Bind(body); err != nil {
		log.Error(RestError(c, err))
		return
	}

	if body.Name == "" {
		log.Error(RestErrorInvalidBody(c))
		return
	}

	username := ApiUser(c).Name
	d := db.NewAppsDbService(username)

	doc := models.JSON{
		"name":      body.Name,
		"owner":     username,
		"types":     []string{"users"},
		"createdAt": time.Now(),
		"masterKey": strings.ToUpper(utils.GetCleanUUID()),
	}

	if err := d.Insert(doc); err != nil {
		log.Error(RestError(c, err))
		return
	}

	appId := doc["_id"]
	appsMapDb := db.NewAppsMapDbService()
	err := appsMapDb.Insert(models.JSON{
		"appId": appId,
		"masterKey": doc["masterKey"],
		"user":  username,
	})

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	RespondId(appId, c)
}

func (a *ApplicationController) GetApplicationsHandler(c *gin.Context) {
	user := ApiUser(c).Name
	d := db.NewAppsDbService(user)

	res, err := d.Find(
		models.JSON{
			"owner": user,
		},
		models.JSON{
			"name": 1,
		},
	)

	if err != nil {
		RestError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *ApplicationController) GetApplicationHandler(c *gin.Context) {
	app := Application(c, c.Param("appId"))
	if app != nil {
		c.JSON(http.StatusOK, app)
	}
}

func (a *ApplicationController) DeleteApplicationHandler(c *gin.Context) {
	appId := c.Param("appId")

	d := db.NewAppsDbService(ApiUser(c).Name)
	err := d.RemoveId(appId)

	if err != nil {
		RestError(c, err)
		return
	}
}

func (a *ApplicationController) UpdateApplicationHandler(c *gin.Context) {
	appId := c.Param("appId")
	d := db.NewAppsDbService(ApiUser(c).Name)
	doc := utils.WhitelistFields([]string{"name"}, webUtils.GetBody(c))

	err := d.Update(models.JSON{
		"_id": appId,
	}, models.JSON{
		"$set": doc,
	})

	if err != nil {
		RestError(c, err)
		return
	}
}
