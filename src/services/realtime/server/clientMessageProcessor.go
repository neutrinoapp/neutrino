package server

import (
	"errors"
	"strconv"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
)

type clientMessageProcessor struct {
	OpProcessors map[string]func(messaging.Message, *client.ApiClient) error
}

var clientsCache map[string]*client.ApiClient

func init() {
	clientsCache = make(map[string]*client.ApiClient)
}

func (p *clientMessageProcessor) Process(mType int, m string) error {
	if mType != messaging.MESSAGE_TYPE_STRING {
		return errors.New("Unsupported message type: " + strconv.Itoa(mType))
	}

	var msg messaging.Message
	err := msg.FromString(m)
	if err != nil {
		log.Error(err)
		return err
	}

	if msg.Origin == messaging.ORIGIN_API {
		log.Info("Skipping processing message from API", m)
		return nil
	}

	log.Info("Got message from client....processing:", m)

	opProcessor := p.OpProcessors[msg.Operation]

	var c *client.ApiClient
	if clientsCache[msg.App] != nil {
		c = clientsCache[msg.App]
	} else {
		apiAddr := config.Get(config.KEY_API_ADDR)
		c = client.NewApiClient(apiAddr, msg.App)
	}

	c.Token = msg.Token

	opts := msg.Options
	if opts != nil && opts["clientId"] != nil {
		c.ClientId = opts["clientId"].(string)
	}

	return opProcessor(msg, c)
}

func opCreate(m messaging.Message, c *client.ApiClient) error {
	_, err := c.CreateItem(m.Type, m.Payload)
	return err
}

func opUpdate(m messaging.Message, c *client.ApiClient) error {
	_, err := c.UpdateItem(m.Type, m.Payload["_id"].(string), m.Payload)
	return err
}

func opDelete(m messaging.Message, c *client.ApiClient) error {
	_, err := c.DeleteItem(m.Type, m.Payload["_id"].(string))
	return err
}

func NewClientMessageProcessor() messaging.MessageProcessor {
	opProcessors := make(map[string]func(messaging.Message, *client.ApiClient) error)
	opProcessors[messaging.OP_CREATE] = opCreate
	opProcessors[messaging.OP_UPDATE] = opUpdate
	opProcessors[messaging.OP_DELETE] = opDelete

	return &clientMessageProcessor{opProcessors}
}
