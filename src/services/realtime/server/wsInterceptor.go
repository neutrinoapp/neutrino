package server

import "gopkg.in/jcelliott/turnpike.v2"

type wsInterceptor struct {
	m chan turnpike.Message
}

func (i *wsInterceptor) Intercept(session turnpike.Session, msg *turnpike.Message) {
	m := *msg
	i.m <- m
}
