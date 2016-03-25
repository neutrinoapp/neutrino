package server

import "github.com/gngeorgiev/gowamp"

type wsInterceptor struct {
	OnMessage chan interceptorMessage
}

type interceptorMessage struct {
	msg         gowamp.Message
	sess        gowamp.Session
	messageType gowamp.MessageType
}

func NewWsInterceptor() *wsInterceptor {
	return &wsInterceptor{
		OnMessage: make(chan interceptorMessage),
	}
}

func (i *wsInterceptor) Intercept(session gowamp.Session, msg *gowamp.Message) {
	innerMessage := *msg
	m := interceptorMessage{
		msg:         innerMessage,
		sess:        session,
		messageType: innerMessage.MessageType(),
	}
	i.OnMessage <- m
}
