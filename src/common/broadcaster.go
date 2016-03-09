package common

import "github.com/neutrinoapp/neutrino/src/common/log"

type Broadcaster struct {
	listeners []chan interface{}
}

func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		listeners: make([]chan interface{}, 0),
	}
	return b
}

func (b *Broadcaster) Broadcast(v interface{}) {
	log.Info(b.listeners)
	for _, ch := range b.listeners {
		ch <- v
	}
}

func (b *Broadcaster) Subscribe(ch chan interface{}) {
	b.listeners = append(b.listeners, ch)
}

func (b *Broadcaster) Remove(ch chan interface{}) {
	for i, listener := range b.listeners {
		if listener == ch {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			return
		}
	}
}
