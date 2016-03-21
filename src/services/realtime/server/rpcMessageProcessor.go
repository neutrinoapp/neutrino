package server

import (
	"encoding/json"
	"fmt"

	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
	"gopkg.in/jcelliott/turnpike.v2"
)

type RpcMessageProcessor struct {
	WsClient    *turnpike.Client
	WsProcessor WsMessageProcessor
	DbService   db.DbService
}

func NewRpcMessageProcessor(wsClient *turnpike.Client, wsProcessor WsMessageProcessor) RpcMessageProcessor {
	return RpcMessageProcessor{wsClient, wsProcessor, db.NewDbService()}
}

func (p RpcMessageProcessor) Process() {
	p.WsClient.BasicRegister("data.read", p.handleDataRead)
	p.WsClient.BasicRegister("data.create", p.handleDataCreate)
	p.WsClient.BasicRegister("data.remove", p.handleDataRemove)
	p.WsClient.BasicRegister("data.update", p.handleDataUpdate)
}

func (p RpcMessageProcessor) getArgs(args []interface{}) (messaging.Message, error) {
	var m messaging.Message

	incomingMsg := args[0]
	log.Info("RPC message:", incomingMsg)

	b, err := json.Marshal(incomingMsg)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (p RpcMessageProcessor) handleDataRead(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	var item interface{}
	if id, ok := m.Payload["id"].(string); ok {
		item, err = p.DbService.GetItemById(id)
	} else {
		item, err = p.DbService.GetItems(m.App, m.Type, m.Options.Filter)
	}

	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{item}}
}

func (p RpcMessageProcessor) handleDataCreate(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	id, err := p.DbService.CreateItem(m.App, m.Type, m.Payload)
	if err != nil {
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{id}}
}

func (p RpcMessageProcessor) handleDataRemove(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	id, ok := m.Payload["id"].(string)
	if !ok {
		return &turnpike.CallResult{Err: turnpike.URI(fmt.Sprintf("Incorrect payload, %v", m.Payload))}
	}

	err = p.DbService.DeleteItemById(id)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{id}}
}

func (p RpcMessageProcessor) handleDataUpdate(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	id, ok := m.Payload["id"].(string)
	if !ok {
		return &turnpike.CallResult{Err: turnpike.URI(fmt.Sprintf("Incorrect payload, %v", m.Payload))}
	}

	err = p.DbService.UpdateItemById(id, m.Payload)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{id}}
}
