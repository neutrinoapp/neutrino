package messaging

import (
	"github.com/go-neutrino/neutrino/models"
	"github.com/gorilla/websocket"
)

var MESSAGE_TYPE_STRING int = 1

type Message struct {
	Operation op
	Origin    origin
	Options   models.JSON
	Payload   models.JSON
	Type      string
}

func (m Message) Send(c *websocket.Conn) error {
	return c.WriteMessage(MESSAGE_TYPE_STRING, m)
}

func (m Message) Serialize() models.JSON {
	return models.JSON{
		"op":      m.Operation,
		"origin":  m.Origin,
		"options": m.Options,
		"pld":     m.Payload,
		"type":    m.Type,
	}
}
