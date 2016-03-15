package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/models"
)

type UserDbService interface {
	DbService

	AppTerm() r.Term
	App() (interface{}, error)
	CreateApp(app interface{}) error
	GetApps() ([]interface{}, error)
	DeleteApp() error
	UpdateApp(interface{}) error
	GetUser(interface{}) (interface{}, error)
}

type userDbService struct {
	DbService
	user, appId string
}

func (u *userDbService) AppTerm() r.Term {
	return u.Query().
		Get(u.user).
		Field(APPS_FIELD).
		Filter(models.JSON{USERS_TABLE
			ID_FIELD: u.appId,
		}).
		Nth(0)
}

func (u *userDbService) App() (app interface{}, err error) {
	c, err := u.AppTerm().Run(u.GetSession())
	if err != nil {
		return
	}

	err = c.One(&app)
	return
}

func (u *userDbService) CreateApp(app interface{}) (err error) {
	_, err = u.Query().Get(u.user).Update(func(user r.Term) interface{} {
		return models.JSON{
			APPS_FIELD: user.Field(APPS_FIELD).Append(app),
		}
	}).RunWrite(u.GetSession())

	return
}

func (u *userDbService) GetApps() (apps []interface{}, err error) {
	c, err := u.Query().Get(u.user).Field(APPS_FIELD).Run(u.GetSession())
	err = c.All(&apps)
	return
}

func (u *userDbService) DeleteApp() (err error) {
	_, err = u.AppTerm().Delete().RunWrite(u.GetSession())
	return
}

func (u *userDbService) UpdateApp(v interface{}) (err error) {
	_, err = u.AppTerm().Update(v).RunWrite(u.GetSession())
	return
}

func (u *userDbService) GetUser(id interface{}) (user interface{}, err error) {
	c, err := u.Query().Get(u.user).Run(u.GetSession())
	err = c.One(&user)
	return
}
