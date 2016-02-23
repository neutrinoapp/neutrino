package messaging

import (
	"encoding/json"

	"github.com/go-neutrino/neutrino/src/common/models"
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

func (m Message) Send(c *websocket.Conn) error {
	model, err := m.ToJson()
	if err != nil {
		return err
	}

	msg, err := model.String()
	if err != nil {
		return err
	}

	return c.WriteMessage(MESSAGE_TYPE_STRING, []byte(msg))
}

func (m *Message) FromString(s string) error {
	if err := json.Unmarshal([]byte(s), m); err != nil {
		return err
	}

	return nil
}

func (m Message) ToJson() (models.JSON, error) {
	var model models.JSON

	if err := model.FromObject(m); err != nil {
		return model, err
	}

	return model, nil
}
