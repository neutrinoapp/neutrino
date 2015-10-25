package messaging

import "github.com/go-neutrino/neutrino/models"

type (
	op     string
	origin string
)

const (
	OP_UPDATE   op = "update"
	OP_CREATE   op = "create"
	OP_DELETE   op = "delete"
	OP_LOGIN    op = "login"
	OP_REGISTER op = "register"

	ORIGIN_API    origin = "api"
	ORIGIN_CLIENT origin = "client"
)

type MessageBuilder interface {
	Build(op, origin, interface{}, models.JSON, string) Message
	BuildFromModel(m models.JSON) Message
}

type messageBuilder struct {
}

var b MessageBuilder

func GetMessageBuilder() MessageBuilder {
	if b == nil {
		b = &messageBuilder{}
	}

	return b
}

func (m *messageBuilder) Build(op op, og origin, pld interface{}, opts models.JSON, t string) Message {
	return Message{
		Operation: op,
		Origin:    og,
		Payload:   pld,
		Options:   opts,
		Type:      t,
	}
}

func (m *messageBuilder) BuildFromModel(m models.JSON) Message {
	return m.Build(
		m["op"].(op),
		m["origin"].(origin),
		m["pld"],
		m["options"].(models.JSON),
		m["type"].(string),
	)
}
