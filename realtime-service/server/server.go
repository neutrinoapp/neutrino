package server

import (
	"github.com/go-neutrino/neutrino-config"
	"github.com/go-neutrino/neutrino-core/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
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
	config *viper.Viper
)

func Initialize(c *viper.Viper) {
	config = c

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

	wsDialer := websocket.Dialer{}
	conn, _, err := wsDialer.Dial(c.GetString(nconfig.KEY_BROKER_HOST)+c.GetString(nconfig.KEY_BROKER_PORT)+"/register", nil)
	if err != nil {
		panic(err)
	}

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Info("read: ", err)
				break
			}

			log.Info("recv: %s", message)
		}
	}()
}
