package messaging

import (
	"encoding/json"

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
	Operation string      `json:"op"`
	Origin    string      `json:"origin"`
	Options   models.JSON `json:"options"`
	Payload   models.JSON `json:"pld"`
	Type      string      `json:"type"`
	App       string      `json:"app"`
	Token     string      `json:"token"`
}

func (m *Message) FromString(s string) error {
	if err := json.Unmarshal([]byte(s), m); err != nil {
		return err
	}

	return nil
}

func (m Message) Send(c *websocket.Conn) error {
	model, err := m.Serialize()
	if err != nil {
		return err
	}

	msg, err := model.String()
	if err != nil {
		return err
	}

	return c.WriteMessage(MESSAGE_TYPE_STRING, []byte(msg))
}

func (m Message) Serialize() (models.JSON, error) {
	//return models.JSON{
	//	"op":      m.Operation,
	//	"origin":  m.Origin,
	//	"options": m.Options,
	//	"pld":     m.Payload,
	//	"type":    m.Type,
	//	"app":     m.App,
	//	"token":   m.Token,
	//}

	var model models.JSON

	if err := model.FromObject(m); err != nil {
		return model, err
	}

	return model, nil
}
