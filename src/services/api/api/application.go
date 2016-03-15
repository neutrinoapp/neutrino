package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
)

type ApplicationModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ApplicationController struct {
	BaseController
}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{NewBaseController()}
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

	email := ApiUser(c).Email
	app := models.JSON{
		db.NAME_FIELD:       body.Name,
		db.OWNER_FIELD:      email,
		db.MASTER_KEY_FIELD: strings.ToUpper(utils.GetCleanUUID()),
	}

	appId, err := a.DbService.CreateApp(email, app)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	RespondId(appId, c)
}

func (a *ApplicationController) GetApplicationsHandler(c *gin.Context) {
	email := ApiUser(c).Email
	apps, err := a.DbService.GetApps(email)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, apps)
}

func (a *ApplicationController) GetApplicationHandler(c *gin.Context) {
	//user := ApiUser(c).Email
	appId := c.Param("appId")

	//TODO: permissions
	app, err := a.DbService.GetApp(appId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, app)
}

func (a *ApplicationController) DeleteApplicationHandler(c *gin.Context) {
	//user := ApiUser(c).Email
	//appId := c.Param("appId")
	//TODO:
	//d := db.NewUserDbService(user, appId)
	//err := d.DeleteApp()
	//
	//if err != nil {
	//	log.Error(RestError(c, err))
	//	return
	//}
	//
	//dataDb := db.NewDataDbService(appId, "")
	//err = dataDb.RemoveApp()
	//if err != nil {
	//	log.Error(RestError(c, err))
	//	return
	//}

	c.Status(http.StatusOK)
}

func (a *ApplicationController) UpdateApplicationHandler(c *gin.Context) {
	//TODO:
	//appId := c.Param("appId")
	//user := ApiUser(c).Email
	//
	//d := db.NewUserDbService(user, appId)
	//app := utils.WhitelistFields([]string{"name"}, webUtils.GetBody(c))
	//
	//err := d.UpdateApp(app)
	//if err != nil {
	//	log.Error(RestError(c, err))
	//	return
	//}

	c.Status(http.StatusOK)
}
