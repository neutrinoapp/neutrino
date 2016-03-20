package db

import (
	"github.com/neutrinoapp/neutrino/src/common/log"

	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
)

var session *r.Session

type DbService interface {
	GetSession() *r.Session
	Run(r.Term) (*r.Cursor, error)
	Exec(terms ...r.Term) error
	Db() r.Term

	CreateApp(userEmail string, app models.JSON) (appId string, err error)
	GetApps(userEmail string) (apps []models.JSON, err error)
	GetApp(appId string) (app models.JSON, err error)

	CreateItem(appId, t string, data models.JSON) (id string, err error)
	GetItems(appId, t string, filter interface{}) (data []models.JSON, err error)
	GetItemById(id string) (item models.JSON, err error)
	UpdateItemById(id string, data interface{}) (err error)
	DeleteItemById(id string) (err error)
	DeleteAllItems(appId, t string) (err error)
	GetTypes(appId string) (types []models.JSON, err error)

	GetUser(email string, isApp bool, appId string) (user models.JSON, err error)
	CreateUser(user models.JSON, isApp bool) (err error)

	Changes(appId, t string, filter, channel interface{}) (err error)
}

type dbService struct {
}

func (d *dbService) setId(item models.JSON) string {
	if item[ID_FIELD] == nil {
		item[ID_FIELD] = utils.GetCleanUUID()
	}

	return item[ID_FIELD].(string)
}

func (d *dbService) getUser(email string) r.Term {
	return d.Db().Table(USERS_TABLE).GetAllByIndex(USERS_TABLE_EMAIL_INDEX, email).Nth(0)
}

func (d *dbService) Db() r.Term {
	return r.DB(DATABASE_NAME)
}

func (d *dbService) GetSession() *r.Session {
	if session == nil {
		addr := config.Get(config.KEY_RETHINK_ADDR)

		s, err := r.Connect(r.ConnectOpts{
			Address: addr,
		})

		if err != nil {
			log.Error(err)
		}

		//TODO: retry until connected
		session = s
	}

	return session
}

func (d *dbService) Run(t r.Term) (*r.Cursor, error) {
	return t.Run(d.GetSession())
}

func (d *dbService) Exec(terms ...r.Term) (err error) {
	s := d.GetSession()
	for _, t := range terms {
		err = t.Exec(s)
		if err != nil {
			return
		}
	}

	return
}

func (d *dbService) CreateApp(userEmail string, app models.JSON) (appId string, err error) {
	appId = d.setId(app)
	err = d.Exec(
		d.Db().Table(APPS_TABLE).Insert(app),
		d.getUser(userEmail).Update(func(user r.Term) interface{} {
			return models.JSON{
				APPS_FIELD: user.Field(APPS_FIELD).Append(appId),
			}
		}),
	)
	return
}

func (d *dbService) GetApps(userEmail string) (apps []models.JSON, err error) {
	c, err := d.Run(
		d.getUser(userEmail).Field(APPS_FIELD).Map(func(appId r.Term) interface{} {
			return d.Db().Table(APPS_TABLE).Get(appId)
		}),
	)
	if err != nil {
		return
	}

	err = c.All(&apps)
	return
}

func (d *dbService) GetApp(appId string) (app models.JSON, err error) {
	c, err := d.Run(
		d.Db().Table(APPS_TABLE).Get(appId),
	)
	if err != nil {
		return
	}

	err = c.One(&app)
	return
}

func (d *dbService) CreateItem(appId, t string, data models.JSON) (id string, err error) {
	id = d.setId(data)
	data[APP_ID_FIELD] = appId
	data[TYPE_FIELD] = t

	log.Info(appId)
	log.Info(t)
	log.Info(data)

	//dataMap := data.ToMap()

	err = d.Exec(
		//d.Db().Table(DATA_TABLE).Insert(models.JSON{
		//	"text":     "teST",
		//	"id":       utils.GetCleanUUID(),
		//	"complete": true,
		//	"appId":    appId,
		//	"type":     t,
		//}),
		d.Db().Table(DATA_TABLE).Insert(data),
	)

	return
}

func (d *dbService) GetItems(appId, t string, filter interface{}) (data []models.JSON, err error) {
	c, err := d.Run(
		d.Db().Table(DATA_TABLE).GetAllByIndex(DATA_TABLE_APPIDTYPE_INDEX, []interface{}{appId, t}).Filter(filter),
	)
	if err != nil {
		return
	}

	err = c.All(&data)
	if err == nil {
		for i, v := range data {
			data[i] = utils.BlacklistFields(DB_FIELDS, v)
		}
	}

	return
}

func (d *dbService) GetItemById(id string) (item models.JSON, err error) {
	c, err := d.Run(
		d.Db().Table(DATA_TABLE).Get(id),
	)
	if err != nil {
		return
	}

	err = c.One(&item)
	if err == nil {
		item = utils.BlacklistFields(DB_FIELDS, item)
	}

	return
}

func (d *dbService) UpdateItemById(id string, data interface{}) (err error) {
	err = d.Exec(
		d.Db().Table(DATA_TABLE).Get(id).Update(data),
	)
	return
}

func (d *dbService) DeleteItemById(id string) (err error) {
	err = d.Exec(
		d.Db().Table(DATA_TABLE).Get(id).Delete(),
	)
	return
}

func (d *dbService) DeleteAllItems(appId, t string) (err error) {
	err = d.Exec(
		d.Db().Table(DATA_TABLE).GetAllByIndex(DATA_TABLE_APPIDTYPE_INDEX, []interface{}{appId, t}).Delete(),
	)
	return
}

func (d *dbService) GetTypes(appId string) (types []models.JSON, err error) {
	c, err := d.Run(
		d.Db().Table(APPS_TABLE).Get(appId).Field(TYPES_FIELD),
	)
	if err != nil {
		return
	}

	err = c.All(types)
	return
}

func (d *dbService) GetUser(email string, isApp bool, appId string) (user models.JSON, err error) {
	var c *r.Cursor
	if isApp {
		c, err = d.Run(
			d.Db().Table(APPS_USERS_TABLE).GetAllByIndex(APPS_USERS_TABLE_EMAILAPPID_INDEX, []interface{}{email, appId}).Nth(0),
		)
	} else {
		c, err = d.Run(
			d.Db().Table(USERS_TABLE).GetAllByIndex(USERS_TABLE_EMAIL_INDEX, email).Nth(0),
		)
	}

	if err != nil {
		return
	}

	err = c.One(&user)
	return
}

func (d *dbService) CreateUser(user models.JSON, isApp bool) (err error) {
	var table string
	if isApp {
		table = APPS_USERS_TABLE
	} else {
		table = USERS_TABLE
	}

	err = d.Exec(
		d.Db().Table(table).Insert(user),
	)
	return
}

func (d *dbService) Changes(appId, t string, filter, channel interface{}) (err error) {
	c, err := d.Run(
		d.Db().Table(DATA_TABLE).GetAllByIndex(DATA_TABLE_APPIDTYPE_INDEX, []interface{}{appId, t}).
			Filter(filter).Changes(),
	)
	if err != nil {
		return nil
	}

	c.Listen(channel)
	return
}
