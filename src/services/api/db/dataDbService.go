package db

import (
	"errors"

	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
)

type DataDbService interface {
	DbService

	InsertData(v map[string]interface{}) (string, error)
	GetData(filter interface{}) ([]interface{}, error)
	GetDataId(filter interface{}) (interface{}, error)
	UpdateId(data map[string]interface{}) error
	RemoveId(id interface{}) error
	RemoveType() error
	RemoveApp() error
	Changes(filter, channel interface{}) error
}

type dataDbService struct {
	DbService
	t     string
	appId string
}

func (t *dataDbService) InsertData(v map[string]interface{}) (key string, err error) {
	if v[ID_FIELD] == nil {
		v[ID_FIELD] = utils.GetCleanUUID()
	}

	key = v[ID_FIELD].(string)

	_, err = t.Query().Get(t.appId).Update(func(app r.Term) r.Term {
		return r.Branch(
			app.Field(TYPES_FIELD).Field(t.t).Default(false).Ne(false),
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

	return
}

func (t *dataDbService) GetData(filter interface{}) (data []interface{}, err error) {
	c, err := t.Query().Get(t.appId).Field(TYPES_FIELD).Field(t.t).Default([]interface{}{}).Filter(filter).Run(t.GetSession())
	if err != nil {
		return
	}

	err = c.All(&data)
	return
}

func (t *dataDbService) GetDataId(id interface{}) (data interface{}, err error) {
	c, err := t.Query().Get(t.appId).Field(TYPES_FIELD).Field(t.t).Filter(models.JSON{ID_FIELD: id}).Run(t.GetSession())
	if err != nil {
		return
	}

	allItems := make([]interface{}, 0)
	err = c.All(&allItems)
	if err != nil {
		return
	}

	if len(allItems) == 0 {
		//TODO: generalize
		err = errors.New("not found")
		return
	}

	data = allItems[0]
	return
}

func (t *dataDbService) UpdateId(data map[string]interface{}) (err error) {
	id := data[ID_FIELD]
	_, err = t.Query().Get(t.appId).Update(func(app r.Term) interface{} {
		return models.JSON{
			TYPES_FIELD: models.JSON{
				t.t: app.Field(TYPES_FIELD).Field(t.t).Map(func(row r.Term) interface{} {
					return r.Branch(
						row.Field(ID_FIELD).Eq(id),
						row.Merge(data),
						row,
					)
				}),
			},
		}
	}).RunWrite(t.GetSession())

	return
}

func (t *dataDbService) RemoveId(id interface{}) (err error) {
	_, err = t.Query().Get(t.appId).Update(func(app r.Term) interface{} {
		return models.JSON{
			TYPES_FIELD: models.JSON{
				t.t: app.Field(TYPES_FIELD).Field(t.t).Filter(func(row r.Term) interface{} {
					return row.Field(ID_FIELD).Ne(id)
				}),
			},
		}
	}).RunWrite(t.GetSession())
	return
}

func (t *dataDbService) RemoveType() (err error) {
	_, err = t.Query().Get(t.appId).Replace(func(app r.Term) interface{} {
		return app.Without(models.JSON{
			TYPES_FIELD: models.JSON{
				t.t: true,
			},
		})
	}).RunWrite(t.GetSession())
	return
}

func (t *dataDbService) RemoveApp() (err error) {
	_, err = t.Query().Get(t.appId).Delete().RunWrite(t.GetSession())
	return
}

func (d *dataDbService) Changes(filter, channel interface{}) error {
	c, err := d.Query().
		GetAll(d.appId).
		Map(func(app r.Term) interface{} {
			return app.Field(TYPES_FIELD).Field(d.t).Filter(filter)
		}).
		Changes().
		Map(func(update r.Term) interface{} {
			return models.JSON{
				"new_val": update.Field("new_val").Difference(update.Field("old_val")).Union([]models.JSON{nil}).Nth(0),
				"old_val": update.Field("old_val").Difference(update.Field("new_val")).Union([]models.JSON{nil}).Nth(0),
			}
		}).
		Run(d.GetSession())

	if err != nil {
		return err
	}

	c.Listen(channel)
	return nil
}
