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

	r "github.com/dancannon/gorethink"
)

type TypesController struct {
}

func (t *TypesController) ensureType(typeName string, c *gin.Context) {
	appId := c.Param("appId")
	user := ApiUser(c).Name

	//we do not need to wait for this op
	d := db.NewAppsDbService(user)
	err := d.GetTable().Get(appId).Update(func(row r.Term) interface{} {
		return r.Branch(
			row.Field("types").Contains(typeName),
			nil,
			models.JSON{
				"types": row.Field("types").Append(typeName),
			},
		)
	}).Exec(d.GetSession(), r.ExecOpts{NoReply: true})
	if err != nil {
		log.Error(err)
	}
}

func (t *TypesController) GetTypesHandler(c *gin.Context) {
	app := Application(c, c.Param("appId"))
	if app != nil {
		c.JSON(http.StatusOK, app["types"])
	}
}

func (t *TypesController) DeleteType(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")

	d := db.NewAppsDbService(ApiUser(c).Name)
	_, err := d.GetTable().Get(appId).Update(func(row r.Term) interface{} {
		return models.JSON{
			"types": row.Field("types").Filter(func(item r.Term) interface{} {
				return item.Ne(typeName)
			}),
		}
	}).RunWrite(d.GetSession())

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	database := db.NewTypeDbService(appId, typeName)
	_, dropError := database.GetTable().TableDrop().RunWrite(database.GetSession())

	if dropError != nil {
		log.Error(RestError(c, dropError))
		return
	}
}

func (t *TypesController) InsertInTypeHandler(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	body := webUtils.GetBody(c)

	t.ensureType(typeName, c)

	d := db.NewTypeDbService(appId, typeName)
	if body == nil {
		body = make(map[string]interface{})
	}

	err := d.Insert(body)
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

	RespondId(body["id"], c)
}

func (t *TypesController) GetTypeDataHandler(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")

	t.ensureType(typeName, c)
	d := db.NewTypeDbService(appId, typeName)

	opts := GetHeaderOptions(c)
	log.Info("Filter: ", opts.Filter)

	typeData, err := d.Find(opts.Filter)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	if typeData == nil {
		typeData = make([]map[string]interface{}, 0)
	}

	c.JSON(http.StatusOK, typeData)
}

func (t *TypesController) GetTypeItemById(c *gin.Context) {
	appId := c.Param("appId")
	typeName := c.Param("typeName")
	itemId := c.Param("itemId")

	t.ensureType(typeName, c)

	d := db.NewTypeDbService(appId, typeName)

	item, err := d.FindId(itemId)

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

	t.ensureType(typeName, c)

	d := db.NewTypeDbService(appId, typeName)
	body := webUtils.GetBody(c)
	body["id"] = itemId

	err := d.ReplaceId(itemId, body)

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

	t.ensureType(typeName, c)

	d := db.NewTypeDbService(appId, typeName)

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
