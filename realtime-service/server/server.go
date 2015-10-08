package server

import (
	"github.com/go-neutrino/neutrino-core/log"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"fmt"
	"github.com/go-neutrino/neutrino-core/config"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//allow connections from any origin
			return true
		},
	}
)

func connectToBroker() *websocket.Conn {
	wsDialer := websocket.Dialer{}
	var conn *websocket.Conn
	brokerHost := config.Get(config.KEY_BROKER_HOST)
	brokerPort := config.Get(config.KEY_BROKER_PORT)
	//retry the connection to the broker until established
	for {
		c, _, err := wsDialer.Dial(brokerHost+brokerPort+"/register", nil)
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second * 5)
		} else {
			log.Info("Connected to broker.")
			conn = c
			break;
		}
	}

	return conn
}

func Initialize() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
			return
		}

		//TODO: token authentication
		token := r.URL.Query().Get("token") //the hash is a unique, per user
		realtimeConn := NewConnection(conn, token)

		GetConnectionStore().Put(token, realtimeConn)
	})

	go func() {
		conn := connectToBroker()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				log.Error("Broker connection error: ", err, "reconnecting....")
				conn = connectToBroker()
				//just in case if something bad happens
				defer conn.Close()
			}

			log.Info(fmt.Sprintf("recv: %s", message))
		}
	}()
}
