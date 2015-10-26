package messaging

import (
	"github.com/go-neutrino/neutrino/models"
	"github.com/gorilla/websocket"
)

const (
	OP_UPDATE string = "update"
	OP_CREATE string = "create"
	OP_DELETE string = "delete"

	ORIGIN_API    string = "api"
	ORIGIN_CLIENT string = "client"
)

var MESSAGE_TYPE_STRING int = 1

type Message struct {
	Operation string
	Origin    string
	Options   models.JSON
	Payload   models.JSON
	Type      string
	App       string
	Token     string
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
		"token":   m.Token,
	}
}
