package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils/webUtils"
)

type TypesController struct {
	*BaseController
}

func NewTypesController() *TypesController {
	return &TypesController{NewBaseController()}
}

func (t *TypesController) GetTypesHandler(c *gin.Context) {
	appId := c.Param("appId")

	types, err := t.DbService.GetTypes(appId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, types)
}

func (t *TypesController) DeleteType(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")

	err := t.DbService.DeleteAllItems(appId, typeName)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}
}

func (t *TypesController) InsertInTypeHandler(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	body := webUtils.GetBody(c)

	id, err := t.DbService.CreateItem(appId, typeName, body)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	RespondId(id, c)
}

func (t *TypesController) GetTypeDataHandler(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	opts := GetHeaderOptions(c)

	items, err := t.DbService.GetItems(appId, typeName, opts.Filter)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	if items == nil {
		items = make([]models.JSON, 0)
	}

	c.JSON(http.StatusOK, items)
}

func (t *TypesController) GetTypeItemById(c *gin.Context) {
	//appId := c.Param("appId")
	//typeName := c.Param("typeName")
	itemId := c.Param("itemId")
	//TODO: this should work out of the box since all ids are unique in the data table
	//no specific app and type needed
	//TODO: permissions

	item, err := t.DbService.GetItemById(itemId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, item)
}

func (t *TypesController) UpdateTypeItemById(c *gin.Context) {
	//appId := c.Param("appId")
	//typeName := c.Param("typeName")
	itemId := c.Param("itemId")
	//TODO: same as get
	body := webUtils.GetBody(c)
	err := t.DbService.UpdateItemById(itemId, body)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}
}

func (t *TypesController) DeleteTypeItemById(c *gin.Context) {
	//appId := c.Param("appId")
	//typeName := c.Param("typeName")
	itemId := c.Param("itemId")

	//TODO: same as above
	err := t.DbService.DeleteItemById(itemId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}
}
