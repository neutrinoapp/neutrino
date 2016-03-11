package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
	"github.com/neutrinoapp/neutrino/src/common/utils/webUtils"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
)

type ApplicationModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
	d := db.NewUserDbService(username, "")

	appId := utils.GetCleanUUID()
	app := models.JSON{
		"id":        appId,
		"name":      body.Name,
		"owner":     username,
		"createdAt": time.Now(),
		"masterKey": strings.ToUpper(utils.GetCleanUUID()),
	}

	if err := d.CreateApp(app); err != nil {
		log.Error(RestError(c, err))
		return
	}

	dataDb := db.NewDbService(db.DATABASE_NAME, db.DATA_TABLE)
	_, err := dataDb.Query().Insert(models.JSON{
		"id":    appId,
		"types": make([]interface{}, 0),
		"users": make([]interface{}, 0),
	}).RunWrite(dataDb.GetSession())
	if err != nil {
		//TODO: rollback
		log.Error(RestError(c, err))
		return
	}

	RespondId(appId, c)
}

func (a *ApplicationController) GetApplicationsHandler(c *gin.Context) {
	user := ApiUser(c).Name
	d := db.NewUserDbService(user, "")
	apps, err := d.GetApps()
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, apps)
}

func (a *ApplicationController) GetApplicationHandler(c *gin.Context) {
	user := ApiUser(c).Name
	appId := c.Param("appId")

	d := db.NewUserDbService(user, appId)
	app, err := d.App()
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, app)
}

func (a *ApplicationController) DeleteApplicationHandler(c *gin.Context) {
	user := ApiUser(c).Name
	appId := c.Param("appId")

	d := db.NewUserDbService(user, appId)
	err := d.DeleteApp()

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	dataDb := db.NewDataDbService(appId, "")
	err = dataDb.RemoveApp()
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.Status(http.StatusOK)
}

func (a *ApplicationController) UpdateApplicationHandler(c *gin.Context) {
	appId := c.Param("appId")
	user := ApiUser(c).Name

	d := db.NewUserDbService(user, appId)
	app := utils.WhitelistFields([]string{"name"}, webUtils.GetBody(c))

	err := d.UpdateApp(app)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.Status(http.StatusOK)
}
