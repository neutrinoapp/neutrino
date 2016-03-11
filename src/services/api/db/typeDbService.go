package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/models"
)

type TypeDbService interface {
	DbService

	InsertData(v interface{}) (string, error)
	GetData(filter interface{}) ([]interface{}, error)
	GetDataId(filter interface{}) (interface{}, error)
	ReplaceId(id, data interface{}) error
	RemoveId(id interface{}) error
}

type typeDbService struct {
	DbService
	t     string
	appId string
}

func (t *typeDbService) InsertData(v interface{}) (key string, err error) {
	res, err := t.Query().Get(t.appId).Update(func(app r.Term) r.Term {
		return r.Branch(
			app.Field(TYPES_FIELD).Default(false).Ne(false),
			models.JSON{
				TYPES_FIELD: models.JSON{
					t.t: app.Field(TYPES_FIELD).Field(t.t).Append(v),
				},
			},
			models.JSON{
				TYPES_FIELD: models.JSON{
					t.t: []interface{}{v},
				},
			},
		)
	}).RunWrite(t.GetSession())

	if len(res.GeneratedKeys) > 0 {
		key = res.GeneratedKeys[0]
	}

	return
}

func (t *typeDbService) GetData(filter interface{}) (data []interface{}, err error) {
	c, err := t.Query().Get(t.appId).Field(TYPES_FIELD).Field(t.t).Filter(filter).Run(t.GetSession())
	if err != nil {
		return
	}

	err = c.All(&data)
	return
}

func (t *typeDbService) GetDataId(id interface{}) (data interface{}, err error) {
	c, err := t.Query().Get(t.appId).Field(TYPES_FIELD).Field(t.t).Get(id).Nth(0).Run(t.GetSession())
	if err != nil {
		return
	}

	err = c.One(data)
	return
}

func (t *typeDbService) ReplaceId(id, data interface{}) (err error) {
	_, err = t.Query().Get(t.appId).Field(TYPES_FIELD).Field(t.t).Get(id).Replace(data).RunWrite(t.GetSession())
	return
}

func (t *typeDbService) RemoveId(id interface{}) (err error) {
	_, err = t.Query().Get(t.appId).Delete().RunWrite(t.GetSession())
	return
}
