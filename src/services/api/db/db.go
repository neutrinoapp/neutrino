package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/config"
)

const (
	DATABASE_NAME = "neutrino"
	USERS_TABLE   = "users"
	DATA_TABLE    = "data"
	TYPES_FIELD   = "types"
	APPS_FIELD    = "apps"
	USERS_FIELD   = "users"
	ID_FIELD      = "id"
)

//TODO: indexes
func NewDbService(dbName, tableName string) DbService {
	address := config.Get(config.KEY_RETHINK_ADDR)
	d := &dbService{address, dbName, tableName, r.Term{}}
	d.query = d.GetTable()

	return d
}

func NewUserDbService(u, appId string) UserDbService {
	d := NewDbService(DATABASE_NAME, USERS_TABLE)
	return &userDbService{d, u, appId}
}

func NewTypeDbService(t, appId string) TypeDbService {
	d := NewDbService(DATABASE_NAME, DATA_TABLE)
	return &typeDbService{d, appId, t}
}
