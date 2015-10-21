package server

import (
	"github.com/go-neutrino/neutrino/log"
	"net/http"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/client"
)

var (
	upgrader = client.NewWebsocketUpgrader()
	brokerClient *client.WebsocketClient
)

func init() {
	brokerHost := config.Get(config.KEY_BROKER_HOST)
	brokerPort := config.Get(config.KEY_BROKER_PORT)
	brokerClient = client.NewWebsocketClient(brokerHost+brokerPort+"/register")
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
		for {
			select {
			case msg := <- brokerClient.Message:
				log.Info("Got message:", msg)
			}
		}
	}()
}
