package server

import (
	"errors"
//	"github.com/go-neutrino/neutrino/log"
	"strconv"
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

//	message := string(m)

	//TODO: this is a message from the client?
	//how about common message handling logic?

	//TODO: send message only to the right clients, sort of filtering on the go
//	for _, conn := range GetConnectionStore().Get("") {
//		conn.Broadcast(message)
//	}
//
//	log.Info(message)

	return nil
}

func NewMessageProcessor() MessageProcessor {
	return &messageProcessor{}
}
