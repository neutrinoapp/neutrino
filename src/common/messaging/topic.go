package messaging

import (
	"fmt"
	"strings"
)

func BuildTopic(m Message) string {
	topic := fmt.Sprintf("%s.%s.%s", m.App, m.Type, strings.ToLower(m.Operation))
	if m.Operation == OP_UPDATE {
		topic += ("." + m.Payload["id"].(string))
	}

	return topic
}

func BuildTopicArbitrary(s ...string) string {
	return strings.Join(s, ".")
}
