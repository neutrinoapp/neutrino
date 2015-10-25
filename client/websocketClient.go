package client

import (
	"github.com/go-neutrino/neutrino/log"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WebsocketClient struct {
	*Client
}

func NewWebsocketClient(addr string) *WebsocketClient {
	wsDialer := websocket.Dialer{}
	connect := func() (interface{}, error) {
		log.Info("Connecting to websocket server:", addr)
		var conn *websocket.Conn
		var err error

		conn, _, err = wsDialer.Dial(addr, nil)

		return conn, err
	}

	c := NewClient(connect, addr)

	wsClient := &WebsocketClient{c}
	wsClient.handleConnection()
	return wsClient
}

func (w *WebsocketClient) GetConnection() *websocket.Conn {
	if w.connection == nil {
		return nil
	}

	return w.connection.(*websocket.Conn)
}

func (w *WebsocketClient) Disconnected() {
	w.Client.Disconnected()
	conn := w.GetConnection()
	if conn != nil {
		conn.Close()
	}
}

func (w *WebsocketClient) handleConnection() {
	var conn *websocket.Conn

	establishConnection := func() *websocket.Conn {
		log.Info("Trying to connect to:", w.Addr)
		w.Connect()
		return w.GetConnection()
	}

	onError := func(err error, initial bool) {
		conn = nil

		if !initial {
			w.Disconnected()
			log.Error("Connection error:", w.Addr, err, "dispatching error event")
			w.Error <- err
		}

		conn = establishConnection()
	}

	go func() {
		for {
			select {
			case err := <-w.error:
				onError(err, false)
			}
		}
	}()

	go func() {
		for {
			if conn != nil {
				_, m, err := conn.ReadMessage()
				message := string(m)
				if err != nil {
					w.error <- err
				} else {
					w.Message <- message
				}
			}
		}
	}()

	time.Sleep(time.Second * 2)
	go onError(nil, true)
}

func NewWebsocketUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//allow connections from any origin
			return true
		},
	}
}
