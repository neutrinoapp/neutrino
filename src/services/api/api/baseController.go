package api

import "github.com/neutrinoapp/neutrino/src/services/api/db"

type BaseController struct {
	DbService db.DbService
}

func NewBaseController() *BaseController {
	return &BaseController{db.NewDbService()}
}
