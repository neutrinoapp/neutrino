package server

import (
	"errors"
	"github.com/go-neutrino/neutrino/client"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/messaging"
	"strconv"
)

type clientMessageProcessor struct {
	OpProcessors map[string]func(messaging.Message, *client.ApiClient) error
}

func (p *clientMessageProcessor) Process(mType int, m messaging.Message) error {
	if mType != messaging.MESSAGE_TYPE_STRING {
		return errors.New("Unsupported message type: " + strconv.Itoa(mType))
	}

	opProcessor := p.OpProcessors[m.Operation]

	apiPort := config.Get(config.KEY_API_PORT)
	//TODO: guess not
	c := client.NewApiClient("http://localhost"+apiPort+"/v1/", m.App)
	return opProcessor(m, c)
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
