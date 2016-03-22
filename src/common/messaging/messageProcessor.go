package messaging

import (
	"errors"
	"fmt"

	"github.com/neutrinoapp/neutrino/src/common/db"
)

type MessageProcessor struct {
	dbService db.DbService
}

func NewMessageProcessor() MessageProcessor {
	return MessageProcessor{db.NewDbService()}
}

func (p MessageProcessor) Process(m Message) (res interface{}, err error) {
	if m.Operation == OP_CREATE {
		return p.handleCreateMessage(m)
	} else if m.Operation == OP_READ {
		return p.handleReadMessage(m)
	} else if m.Operation == OP_UPDATE {
		err = p.handleUpdateMessage(m)
		return
	} else if m.Operation == OP_DELETE {
		err = p.handleDeleteMessage(m)
		return
	}

	return nil, errors.New(fmt.Sprintf("Invalid message! %v", m))
}

func (p MessageProcessor) handleCreateMessage(m Message) (id string, err error) {
	return p.dbService.CreateItem(m.App, m.Type, m.Payload)
}

func (p MessageProcessor) handleReadMessage(m Message) (data interface{}, err error) {
	if id, ok := m.Payload[db.ID_FIELD].(string); ok {
		return p.dbService.GetItemById(id)
	}

	return p.dbService.GetItems(m.App, m.Type, m.Options.Filter)
}

func (p MessageProcessor) handleUpdateMessage(m Message) (err error) {
	if id, ok := m.Payload[db.ID_FIELD].(string); ok {
		return p.dbService.UpdateItemById(id, m.Payload)
	}

	return errors.New(fmt.Sprintf("Cannot update item without id %v", m))
}

func (p MessageProcessor) handleDeleteMessage(m Message) (err error) {
	if id, ok := m.Payload[db.ID_FIELD].(string); ok {
		return p.dbService.DeleteItemById(id)
	}

	return p.dbService.DeleteAllItems(m.App, m.Type)
}
