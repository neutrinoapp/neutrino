package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils/webUtils"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
	"github.com/neutrinoapp/neutrino/src/services/api/notification"
)

type TypesController struct {
}

func (t *TypesController) GetTypesHandler(c *gin.Context) {
	appId := c.Param("appId")

	d := db.NewDbService(db.DATABASE_NAME, db.DATA_TABLE)
	cu, err := d.Query().Get(appId).Field(db.TYPES_FIELD).Keys().Run(d.GetSession())
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	var types []interface{}
	err = cu.All(&types)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, types)
}

func (t *TypesController) DeleteType(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")

	dataDb := db.NewDataDbService(appId, typeName)
	err := dataDb.RemoveType()
	if err != nil {
		log.Error(RestError(c, err))
		return
	}
}

func (t *TypesController) InsertInTypeHandler(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	body := webUtils.GetBody(c)

	d := db.NewDataDbService(appId, typeName)
	if body == nil {
		body = make(map[string]interface{})
	}

	id, err := d.InsertData(body)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	opts := GetHeaderOptions(c)
	if *opts.Notify {
		messageBuilder := messaging.GetMessageBuilder()
		token := ApiUser(c).Key
		notification.Notify(messageBuilder.Build(
			messaging.OP_CREATE,
			messaging.ORIGIN_API,
			body,
			opts,
			typeName,
			appId,
			token,
		))
	}

	RespondId(id, c)
}

func (t *TypesController) GetTypeDataHandler(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")

	d := db.NewDataDbService(appId, typeName)

	opts := GetHeaderOptions(c)
	log.Info("Filter: ", opts.Filter)

	typeData, err := d.GetData(opts.Filter)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	if typeData == nil {
		typeData = make([]interface{}, 0)
	}

	c.JSON(http.StatusOK, typeData)
}

func (t *TypesController) GetTypeItemById(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	itemId := c.Param("itemId")

	d := db.NewDataDbService(appId, typeName)

	item, err := d.GetDataId(itemId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, item)
}

func (t *TypesController) UpdateTypeItemById(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	itemId := c.Param("itemId")

	d := db.NewDataDbService(appId, typeName)
	body := webUtils.GetBody(c)
	body[db.ID_FIELD] = itemId

	err := d.ReplaceId(body)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	opts := GetHeaderOptions(c)
	if *opts.Notify {
		messageBuilder := messaging.GetMessageBuilder()
		token := ApiUser(c).Key
		notification.Notify(messageBuilder.Build(
			messaging.OP_UPDATE,
			messaging.ORIGIN_API,
			body,
			opts,
			typeName,
			appId,
			token,
		))
	}
}

func (t *TypesController) DeleteTypeItemById(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	itemId := c.Param("itemId")

	d := db.NewDataDbService(appId, typeName)

	err := d.RemoveId(itemId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	opts := GetHeaderOptions(c)
	if *opts.Notify {
		messageBuilder := messaging.GetMessageBuilder()
		token := ApiUser(c).Key
		notification.Notify(messageBuilder.Build(
			messaging.OP_DELETE,
			messaging.ORIGIN_API,
			models.JSON{"id": itemId},
			opts,
			typeName,
			appId,
			token,
		))
	}
}
