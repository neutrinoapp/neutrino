package api

import "github.com/neutrinoapp/neutrino/src/common/db"

type BaseController struct {
	DbService db.DbService
}

func NewBaseController() *BaseController {
	return &BaseController{db.NewDbService()}
}
