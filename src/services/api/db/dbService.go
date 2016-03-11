package db

import (
	"github.com/neutrinoapp/neutrino/src/common/log"

	r "github.com/dancannon/gorethink"
)

var session *r.Session

type DbService interface {
	GetSession() *r.Session
	GetDb() r.Term
	GetTable() r.Term
	Query() r.Term
}

type dbService struct {
	address, dbName, tableName string
	query                      r.Term
}

func (d *dbService) GetSession() *r.Session {
	if session == nil {
		s, err := r.Connect(r.ConnectOpts{
			Address: d.address,
		})

		log.Info("Connected to rethinkdb:", d.address)

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
	return d.GetDb().Table(d.tableName)
}

func (d *dbService) Query() r.Term {
	return d.query
}
