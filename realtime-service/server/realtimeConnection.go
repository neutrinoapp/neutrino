package server

import (
	"github.com/go-neutrino/neutrino/log"
	"github.com/gorilla/websocket"
)

type RealtimeConnection interface {
	GetProcessor() MessageProcessor
	GetConnection() *websocket.Conn
	Broadcast(string) error
	Listen() error
	Close() error
}

type realtimeConnection struct {
	conn      *websocket.Conn
	processor MessageProcessor
	group     string
}

func (r *realtimeConnection) GetConnection() *websocket.Conn {
	return r.conn
}

func (r *realtimeConnection) Broadcast(m string) error {
	log.Info("Broadcasting message to clients:", m)
	return r.conn.WriteMessage(MESSAGE_TYPE_STRING, []byte(m))
}

func (r *realtimeConnection) Listen() error {
	for {
		messageType, m, err := r.conn.ReadMessage()
		if err != nil {
			return err
		}

		processErr := r.processor.Process(messageType, m)
		if processErr != nil {
			return processErr
		}
	}

	return nil
}

func (r *realtimeConnection) Close() error {
	return r.conn.Close()
}

func (r *realtimeConnection) GetProcessor() MessageProcessor {
	return r.processor
}

func NewConnection(c *websocket.Conn, g string) RealtimeConnection {
	r := &realtimeConnection{
		conn:      c,
		group:     g,
		processor: NewMessageProcessor(),
	}

	go func() {
		err := r.Listen()

		if err != nil {
			defer r.Close()

			log.Error(err)
			GetConnectionStore().Remove(g, r)
		}
	}()

	return r
}
