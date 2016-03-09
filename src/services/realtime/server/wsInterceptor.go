package server

import (
	"github.com/neutrinoapp/neutrino/src/common/log"
	"gopkg.in/jcelliott/turnpike.v2"
)

type wsInterceptor struct {
	OnMessage chan interceptorMessage
}

type interceptorMessage struct {
	msg         turnpike.Message
	sess        turnpike.Session
	messageType turnpike.MessageType
}

func NewWsInterceptor() *wsInterceptor {
	return &wsInterceptor{
		OnMessage: make(chan interceptorMessage),
	}
}

func (i *wsInterceptor) Intercept(session turnpike.Session, msg *turnpike.Message) {
	innerMessage := *msg
	m := interceptorMessage{
		msg:         innerMessage,
		sess:        session,
		messageType: innerMessage.MessageType(),
	}
	log.Info("Intercepting message:", innerMessage, session, innerMessage.MessageType())
	i.OnMessage <- m
}
