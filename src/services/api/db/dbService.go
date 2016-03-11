package db

import (
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"

	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/utils"
)

var session *r.Session

type DbService interface {
	GetSession() *r.Session
	GetDb() r.Term
	GetTable() r.Term
	//Insert(doc map[string]interface{}) error
	//Update(q, u map[string]interface{}) error
	//FindId(id interface{}) (map[string]interface{}, error)
	//Find(query interface{}) ([]map[string]interface{}, error)
	//FindOne(query interface{}) (map[string]interface{}, error)
	//RemoveId(id interface{}) error
	//UpdateId(id, u interface{}) error
	//ReplaceId(id, u interface{}) error
	Changes(filter, channel interface{}) error
	Query() r.Term
	//Path(path ...string) DbService
}

type dbService struct {
	address, dbName, tableName string
	query                      r.Term
}

func NewDbService(dbName, tableName string) DbService {
	address := config.Get(config.KEY_REDIS_ADDR)
	d := &dbService{address, dbName, tableName, r.Term{}}
	d.query = d.GetTable()

	return d
}

func (d *dbService) Path(path ...string) DbService {
	t := d.Query()

	for _, p := range path {
		t = t.Field(p)
	}

	d.query = t
	return d
}

func (d *dbService) GetSession() *r.Session {
	if session == nil {
		s, err := r.Connect(r.ConnectOpts{
			Address: d.address,
		})

		if err != nil {
			log.Error(err)
			panic(err)
		}

		session = s
	}

	return session
}

func (d *dbService) GetDb() r.Term {
	return r.DB(d.dbName)
}

func (d *dbService) GetTable() r.Term {
	return d.GetDb()
}

func (d *dbService) Query() r.Term {
	return d.query
}

func (d *dbService) Insert(doc map[string]interface{}) error {
	if doc["id"] == nil {
		doc["id"] = utils.GetCleanUUID()
	}

	_, err := d.Query().Insert(doc).RunWrite(d.GetSession())
	return err
}

func (d *dbService) Update(q, u map[string]interface{}) error {
	_, err := d.Query().Filter(q).Update(u).RunWrite(d.GetSession())

	return err
}

func (d *dbService) FindId(id interface{}) (map[string]interface{}, error) {
	c, err := d.Query().Get(id).Run(d.GetSession())
	defer c.Close()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = c.One(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *dbService) Find(query interface{}) ([]map[string]interface{}, error) {
	c, err := d.Query().Filter(query).Run(d.GetSession())
	defer c.Close()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = c.All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *dbService) FindOne(query interface{}) (map[string]interface{}, error) {
	c, err := d.Query().Filter(query).Limit(1).Run(d.GetSession())
	defer c.Close()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = c.All(&result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (d *dbService) RemoveId(id interface{}) error {
	_, err := d.Query().Get(id).Delete().RunWrite(d.GetSession())
	return err
}

func (d *dbService) UpdateId(id, u interface{}) error {
	_, err := d.Query().Get(id).Update(u).RunWrite(d.GetSession())
	return err
}

func (d *dbService) ReplaceId(id, u interface{}) error {
	_, err := d.Query().Get(id).Replace(u).RunWrite(d.GetSession())
	return err
}

func (d *dbService) Changes(filter, channel interface{}) error {
	c, err := d.Query().Filter(filter).Changes().Run(d.GetSession())
	if err != nil {
		return err
	}

	c.Listen(channel)
	return nil
}
