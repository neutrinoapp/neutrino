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
	Exec(r.Term) error
	CreateApp(userEmail string, app models.JSON) (appId string, err error)
	GetApps(userEmail string) (apps []models.JSON, err error)
	GetApp(appId string) (app models.JSON, err error)
}

type dbService struct {
}

func (d *dbService) setId(item models.JSON) string {
	if item[ID_FIELD] == nil {
		item[ID_FIELD] = utils.GetCleanUUID()
	}

	return item[ID_FIELD]
}

func (d *dbService) getUser(email string) r.Term {
	return d.db().Table(USERS_TABLE).GetAllByIndex(EMAIL_INDEX, email).Nth(0)
}

func (d *dbService) db() r.Term {
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
		d.db().Table(APPS_TABLE).Insert(app),
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
			return d.db().Table(APPS_TABLE).Get(appId)
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
		d.db().Table(APPS_TABLE).Get(appId),
	)
	if err != nil {
		return
	}

	err = c.One(&app)
	return
}
