package server

import (
	"errors"
	"github.com/go-neutrino/neutrino/client"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"net/http"
)

var (
	upgrader     = client.NewWebsocketUpgrader()
	brokerClient *client.WebsocketClient
)

func init() {
	brokerHost := config.Get(config.KEY_BROKER_HOST)
	brokerPort := config.Get(config.KEY_BROKER_PORT)
	brokerClient = client.NewWebsocketClient(brokerHost + brokerPort + "/register")
}

func Initialize() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Error(err)
			return
		}

		//TODO: token authentication
		appId := r.URL.Query().Get("app")
		realtimeConn := NewConnection(conn, appId)

		log.Info("New connection for app:", appId)

		GetConnectionStore().Put(appId, realtimeConn)
	})

	go func() {
		for {
			select {
			case msg := <-brokerClient.Message:
				{
					log.Info("Realtime service got message from broker, broadcasting:", msg)
					var m models.JSON
					m.FromString([]byte(msg))
					opts := m["options"]
					noAppIdErr := errors.New("No appId provided with realtime notification.")
					if opts == nil {
						log.Error(noAppIdErr)
						return
					}

					appId := opts.(map[string]interface{})["appId"]
					if appId == nil {
						log.Error(noAppIdErr)
						return
					}

					idStr := appId.(string)
					log.Info("Broadcasting:", msg, "to", idStr)
					for _, conn := range GetConnectionStore().Get(idStr) {
						conn.Broadcast(msg)
					}
				}
			}
		}
	}()
}
