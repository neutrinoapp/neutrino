package server

import (
	"encoding/json"
	"fmt"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/jcelliott/turnpike.v2"
)

type RpcMessageProcessor struct {
	WsClient    *turnpike.Client
	WsProcessor WsMessageProcessor
}

func (p RpcMessageProcessor) Process() {
	p.WsClient.BasicRegister("data.read", p.handleDataRead)
	p.WsClient.BasicRegister("data.create", p.handleDataCreate)
	p.WsClient.BasicRegister("data.remove", p.handleDataRemove)
	p.WsClient.BasicRegister("data.update", p.handleDataUpdate)
}

func (p RpcMessageProcessor) getArgs(args []interface{}) (messaging.Message, *client.ApiClient, error) {
	var m messaging.Message

	incomingMsg := args[0]
	log.Info("RPC message:", incomingMsg)

	b, err := json.Marshal(incomingMsg)
	if err != nil {
		return m, nil, err
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		return m, nil, err
	}

	c := client.NewApiClientCached(m.App)
	c.Token = m.Token
	if m.Options.Notify != nil {
		c.NotifyRealTime = *m.Options.Notify
	} else {
		c.NotifyRealTime = false
	}

	return m, c, nil
}

func (p RpcMessageProcessor) handleDataRead(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, c, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	var clientResult interface{}
	if id, ok := m.Payload["_id"].(string); ok {
		clientResult, err = c.GetItem(m.Type, id)
	} else {
		clientResult, err = c.GetItems(m.Type)
	}

	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{clientResult}}
}

func (p RpcMessageProcessor) handleDataCreate(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, c, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	resp, err := c.CreateItem(m.Type, m.Payload)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	p.WsProcessor.HandlePublish(&turnpike.Publish{
		Topic:     turnpike.URI(m.Topic),
		Arguments: args,
	})

	log.Info(resp)
	return &turnpike.CallResult{Args: []interface{}{resp["_id"]}}
}

func (p RpcMessageProcessor) handleDataRemove(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, c, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	id, ok := m.Payload["_id"].(string)
	if !ok {
		return &turnpike.CallResult{Err: turnpike.URI(fmt.Sprintf("Incorrect payload, %v", m.Payload))}
	}

	_, err = c.DeleteItem(m.Type, id)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{id}}
}

func (p RpcMessageProcessor) handleDataUpdate(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
	m, c, err := p.getArgs(args)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	id, ok := m.Payload["_id"].(string)
	if !ok {
		return &turnpike.CallResult{Err: turnpike.URI(fmt.Sprintf("Incorrect payload, %v", m.Payload))}
	}

	_, err = c.UpdateItem(m.Type, id, m.Payload)
	if err != nil {
		log.Error(err)
		return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
	}

	return &turnpike.CallResult{Args: []interface{}{id}}
}
