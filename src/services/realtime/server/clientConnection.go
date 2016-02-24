package server

import (
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/gorilla/websocket"
)

type ClientConnection interface {
	GetProcessor() messaging.MessageProcessor
	GetConnection() *websocket.Conn
	GetClientId() string
	GetAppId() string
	Broadcast(string) error
	Listen() error
	Close() error
}

type clientConnection struct {
	conn      *websocket.Conn
	processor messaging.MessageProcessor
	appId     string
	clientId  string
}

func (r *clientConnection) GetClientId() string {
	return r.clientId
}

func (r *clientConnection) GetAppId() string {
	return r.appId
}

func (r *clientConnection) GetConnection() *websocket.Conn {
	return r.conn
}

func (r *clientConnection) Broadcast(m string) error {
	log.Info("Broadcasting message to client:", m, "with id:", r.clientId)
	return r.conn.WriteMessage(messaging.MESSAGE_TYPE_STRING, []byte(m))
}

func (r *clientConnection) Listen() error {
	for {
		messageType, m, err := r.conn.ReadMessage()
		if err != nil {
			return err
		}

		messageJson := models.JSON{}
		messageJson.FromString(m)

		message := messaging.GetMessageBuilder().BuildFromModel(messageJson)
		message.App = r.appId

		processErr := r.processor.Process(messageType, message)
		if processErr != nil {
			return processErr
		}
	}

	return nil
}

func (r *clientConnection) Close() error {
	return r.conn.Close()
}

func (r *clientConnection) GetProcessor() messaging.MessageProcessor {
	return r.processor
}

func NewConnection(c *websocket.Conn, appId, clientId string) ClientConnection {
	r := &clientConnection{
		conn:      c,
		appId:     appId,
		clientId:  clientId,
		processor: NewClientMessageProcessor(),
	}

	go func() {
		err := r.Listen()

		if err != nil {
			defer r.Close()

			log.Error(err)
			GetConnectionStore().Remove(appId, r)
		}
	}()

	return r
}
