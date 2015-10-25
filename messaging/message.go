package messaging

import (
	"github.com/go-neutrino/neutrino/models"
	"github.com/gorilla/websocket"
)

type (
	Op     string
	Origin string
)

const (
	OP_UPDATE Op = "update"
	OP_CREATE Op = "create"
	OP_DELETE Op = "delete"

	ORIGIN_API    Origin = "api"
	ORIGIN_CLIENT Origin = "client"
)

var MESSAGE_TYPE_STRING int = 1

type Message struct {
	Operation Op
	Origin    Origin
	Options   models.JSON
	Payload   models.JSON
	Type      string
	App       string
}

func (m Message) Send(c *websocket.Conn) error {
	msg, err := m.Serialize().String()
	if err != nil {
		return err
	}

	return c.WriteMessage(MESSAGE_TYPE_STRING, []byte(msg))
}

func (m Message) Serialize() models.JSON {
	return models.JSON{
		"op":      m.Operation,
		"origin":  m.Origin,
		"options": m.Options,
		"pld":     m.Payload,
		"type":    m.Type,
		"app":     m.App,
	}
}
