package server

import "gopkg.in/jcelliott/turnpike.v2"

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
	i.OnMessage <- m
}
