package messaging

import (
	"fmt"
	"strings"
)

func GetTopic(m Message) string {
	return fmt.Sprintf("%s.%s.%s", m.App, m.Type, strings.ToLower(m.Operation))
}
