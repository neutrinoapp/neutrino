package messaging

import (
	"encoding/json"

	"github.com/neutrinoapp/neutrino/src/common/models"
)

const (
	OP_UPDATE string = "update"
	OP_CREATE string = "create"
	OP_DELETE string = "delete"

	ORIGIN_API    string = "api"
	ORIGIN_CLIENT string = "client"
)

type Message struct {
	Operation string         `json:"op"`
	Origin    string         `json:"origin"`
	Options   models.Options `json:"options"`
	Payload   models.JSON    `json:"pld"`
	Type      string         `json:"type"`
	App       string         `json:"app"`
	Token     string         `json:"token"`
	Topic     string         `json:"topic"`
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

func (m Message) String() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
