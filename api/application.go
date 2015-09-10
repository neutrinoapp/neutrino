package api

import (
	"time"
	"net/http"
	"github.com/go-neutrino/neutrino-core/db"
	"github.com/go-neutrino/neutrino-core/utils"
	"github.com/gin-gonic/gin"
)

type ApplicationModel struct {
	Id string `json: _id`
	Name string `json: "name"`
}

type ApplicationController struct {
}

func (a *ApplicationController) CreateApplicationHandler(c *gin.Context) {
	body := &ApplicationModel{}

	if err := c.Bind(body); err != nil {
		RestError(c, err)
		return
	}

	if body.Name == "" {
		RestErrorInvalidBody(c)
		return
	}

	d := db.NewAppsDbService(c.MustGet("user").(string))

	username := c.MustGet("user").(string)
	doc := JSON{
		"name": body.Name,
		"owner": username,
		"types": []string{"users"},
		"createdAt": time.Now(),
		"masterKey": utils.GetCleanUUID(),
	}

	if err := d.Insert(doc); err != nil {
		RestError(c, err)
		return
	}

	RespondId(doc["_id"], c)
}

func (a *ApplicationController) GetApplicationsHandler(c *gin.Context) {
	user := c.MustGet("user").(string)
	d := db.NewAppsDbService(user)

	res, err := d.Find(
		JSON{
			"owner": user,
		},
		JSON{
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
	app := c.MustGet("app")
	c.JSON(http.StatusOK, app)
}

func (a *ApplicationController) DeleteApplicationHandler(c *gin.Context) {
	appId := c.MustGet("app").(JSON)["_id"]

	d := db.NewAppsDbService(c.MustGet("user").(string))
	err := d.RemoveId(appId)

	if err != nil {
		RestError(c, err)
		return
	}
}

func (a *ApplicationController) UpdateApplicationHandler(c *gin.Context) {
	appId := c.MustGet("app").(JSON)["_id"]
	d := db.NewAppsDbService(c.MustGet("user").(string))
	doc := utils.WhitelistFields([]string{"name"}, utils.GetBody(c))

	err := d.Update(JSON{
		"_id": appId,
	}, JSON{
		"$set": doc,
	})

	if err != nil {
		RestError(c, err)
		return
	}
}