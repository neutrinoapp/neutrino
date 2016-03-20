package server

import (
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
)

type clientMessageProcessor struct {
	DbService db.DbService
}

func (p *clientMessageProcessor) Process(m string) (interface{}, error) {
	var msg messaging.Message
	err := msg.FromString(m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if msg.Origin == messaging.ORIGIN_API {
		return nil, nil
	}

	log.Info("Got message from client....processing:", m)

	//TODO: this  should not be repeated
	var resp interface{}
	var opErr error
	if msg.Type == messaging.OP_CREATE {
		resp, opErr = p.DbService.CreateItem(msg.App, msg.Type, msg.Payload)
	} else if msg.Type == messaging.OP_DELETE {
		//TODO: handle delete all
		opErr = p.DbService.DeleteItemById(msg.Payload[db.ID_FIELD].(string))
	} else if msg.Type == messaging.OP_UPDATE {
		opErr = p.DbService.UpdateItemById(msg.Payload[db.ID_FIELD].(string), msg.Payload)
	}

	return resp, opErr
}

func NewClientMessageProcessor() messaging.MessageProcessor {
	return &clientMessageProcessor{db.NewDbService()}
}
