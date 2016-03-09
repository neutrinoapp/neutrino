package server

import (
	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
)

type clientMessageProcessor struct {
	OpProcessors map[string]func(messaging.Message, *client.ApiClient) (interface{}, error)
}

func (p *clientMessageProcessor) Process(m string) (interface{}, error) {
	var msg messaging.Message
	err := msg.FromString(m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if msg.Origin == messaging.ORIGIN_API {
		log.Info("Skipping processing message from API", m)
		return nil, nil
	}

	log.Info("Got message from client....processing:", m)

	opProcessor := p.OpProcessors[msg.Operation]

	c := client.NewApiClientCached(msg.App)
	c.Token = msg.Token
	if msg.Options.Notify != nil {
		c.NotifyRealTime = *msg.Options.Notify
	} else {
		c.NotifyRealTime = true
	}

	opts := msg.Options
	if *opts.ClientId != "" {
		c.ClientId = *opts.ClientId
	}

	resp, err := opProcessor(msg, c)
	log.Info("Api response: ", resp, err)
	return resp, err
}

func opCreate(m messaging.Message, c *client.ApiClient) (interface{}, error) {
	return c.CreateItem(m.Type, m.Payload)
}

func opUpdate(m messaging.Message, c *client.ApiClient) (interface{}, error) {
	return c.UpdateItem(m.Type, m.Payload["id"].(string), m.Payload)
}

func opDelete(m messaging.Message, c *client.ApiClient) (interface{}, error) {
	return c.DeleteItem(m.Type, m.Payload["id"].(string))
}

func NewClientMessageProcessor() messaging.MessageProcessor {
	opProcessors := make(map[string]func(messaging.Message, *client.ApiClient) (interface{}, error))
	opProcessors[messaging.OP_CREATE] = opCreate
	opProcessors[messaging.OP_UPDATE] = opUpdate
	opProcessors[messaging.OP_DELETE] = opDelete

	return &clientMessageProcessor{opProcessors}
}
