package messaging

import (
	"strings"

	"github.com/go-neutrino/neutrino/messaging"
)

func BuildTopic(m Message) string {
	tokens := []string{m.App, m.Type, strings.ToLower(m.Operation)}
	if m.Operation == messaging.OP_UPDATE {
		tokens = append(tokens, m.Payload["_id"].(string))
	}

	return strings.Join(tokens, ".")
}
