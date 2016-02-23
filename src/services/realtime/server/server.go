package server

import (
	"errors"
	"net/http"

	"github.com/go-neutrino/neutrino/src/common/client"
	"github.com/go-neutrino/neutrino/src/common/config"
	"github.com/go-neutrino/neutrino/src/common/log"
	"github.com/go-neutrino/neutrino/src/common/messaging"
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
		clientId := r.URL.Query().Get("id")
		realtimeConn := NewConnection(conn, appId, clientId)

		log.Info("New connection for app:", appId)

		GetConnectionStore().Put(appId, realtimeConn)
	})

	go func() {
		for {
			select {
			case msg := <-brokerClient.Message:
				{
					log.Info("Realtime service got message from broker, broadcasting:", msg)
					var m messaging.Message
					if err := m.FromString(msg); err != nil {
						log.Error(err)
						continue
					}

					appId := m.App
					if appId == "" {
						log.Error(errors.New("No appId provided with realtime notification."), m)
						continue
					}

					log.Info("Broadcasting:", msg, "to", appId)
					connsForApp := GetConnectionStore().Get(appId)
					for _, conn := range connsForApp {
						connClientId := conn.GetClientId()
						if m.Origin == messaging.ORIGIN_CLIENT &&
							m.Options != nil &&
							m.Options["clientId"] != nil &&
							m.Options["clientId"] == connClientId {

							log.Info("Skipping broadcast to client", connClientId, "has same id.")
							continue
						}

						conn.Broadcast(msg)
					}
				}
			}
		}
	}()
}
