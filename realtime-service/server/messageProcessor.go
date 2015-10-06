package server

import (
	"errors"
	"strconv"
	"github.com/go-neutrino/neutrino-core/log"
)

var (
	MESSAGE_TYPE_STRING int = 1
)

type MessageProcessor interface {
	Process(mType int, m []byte) error
}

type messageProcessor struct {
}

func (p *messageProcessor) Process(mType int, m []byte) error {
	if mType != MESSAGE_TYPE_STRING {
		return errors.New("Unsupported message type: " + strconv.Itoa(mType))
	}

	//TODO:
	log.Info(string(m))

	return nil
}

func NewMessageProcessor() MessageProcessor {
	return &messageProcessor{}
}
