package messaging

import "github.com/go-neutrino/neutrino/models"

type MessageBuilder interface {
	Build(op Op, origin Origin, payload models.JSON, options models.JSON, t, app string) Message
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

func (b *messageBuilder) Build(op Op, og Origin, pld models.JSON, opts models.JSON, t, app string) Message {
	return Message{
		Operation: op,
		Origin:    og,
		Payload:   pld,
		Options:   opts,
		Type:      t,
		App:       app,
	}
}

func (b *messageBuilder) BuildFromModel(m models.JSON) Message {
	return b.Build(
		m["op"].(Op),
		m["origin"].(Origin),
		m["pld"].(models.JSON),
		m["options"].(models.JSON),
		m["type"].(string),
		m["app"].(string),
	)
}
