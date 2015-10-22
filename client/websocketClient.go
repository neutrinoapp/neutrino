package client

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/go-neutrino/neutrino/log"
	"time"
)

type WebsocketClient struct {
	*Client
}

func NewWebsocketClient(addr string) *WebsocketClient {
	wsDialer := websocket.Dialer{}
	connect := func () (interface{}, error) {
			var conn *websocket.Conn
			var err error

			conn, _, err = wsDialer.Dial(addr, nil)

			return conn, err
		}

	c := NewClient(connect, addr)

	wsClient := &WebsocketClient{
		c,
	}

	wsClient.autoProcess()

	return wsClient
}

func (w *WebsocketClient) GetConnection() *websocket.Conn {
	if (w.connection == nil) {
		return nil
	}

	return w.connection.(*websocket.Conn)
}

func (w *WebsocketClient) Disconnected() {
	w.Client.Disconnected()
	conn := w.GetConnection()
	if (conn != nil) {
		conn.Close()
	}
}

func (w *WebsocketClient) autoProcess() {
	var conn *websocket.Conn

	establishConnection := func () *websocket.Conn {
		log.Info("Trying to connect to:", w.Addr)
		w.Connect()
		return w.GetConnection()
	}

	onError := func(err error) {
		conn = nil
		w.Disconnected()
		log.Error("Connection error:", w.Addr, err)
		conn = establishConnection()
	}


	go func () {
		for {
			select {
			case err := <- w.error:
				onError(err)
			}
		}
	}()

	go func () {
		for {
			if conn != nil {
				_, m, err := conn.ReadMessage()
				message := string(m)
				if err != nil {
					w.error <- err
				} else {
					log.Info("Client got message from service:", message)
					w.Message <- message
				}
			}
		}
	}()

	time.Sleep(time.Second * 2)
	onError(nil)
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