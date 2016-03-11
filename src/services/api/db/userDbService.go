package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/models"
)

type UserDbService interface {
	DbService

	AppTerm() r.Term
	App() (interface{}, error)
}

type userDbService struct {
	*dbService
	user, appId string
}

func (u *userDbService) AppTerm() r.Term {
	return u.Query().
		Get(u.user).
		Field(APPS_FIELD).
		Filter(models.JSON{
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
