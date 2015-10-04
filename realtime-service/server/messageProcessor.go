package server

import (
	"errors"
	"fmt"
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

	//TODO:
	fmt.Println(string(m))

	return nil
}

func NewMessageProcessor() MessageProcessor {
	return &messageProcessor{}
}
