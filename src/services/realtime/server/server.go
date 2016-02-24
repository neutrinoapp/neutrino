package server

import (
	"errors"
	"net/http"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/redis.v3"
)

var (
	upgrader             = client.NewWebsocketUpgrader()
	realtimeRedisSubject string
	redisClient          *redis.Client
)

func init() {
	realtimeRedisSubject = config.Get(config.CONST_REALTIME_JOBS_SUBJ)
	redisClient = client.GetNewRedisClient()
}

func Initialize() *http.Server {
	turnpike.Debug()

	s := turnpike.NewBasicWebsocketServer(config.CONST_DEFAULT_REALM)
	s.Upgrader.CheckOrigin = func(r *http.Request) bool {
		//allow connections from any origin
		return true
	}

	server := &http.Server{
		Handler: s,
		Addr:    config.Get(config.KEY_REALTIME_PORT),
	}

	//http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
	//	conn, err := upgrader.Upgrade(w, r, nil)
	//
	//	if err != nil {
	//		log.Error(err)
	//		return
	//	}
	//
	//	//TODO: token authentication
	//	appId := r.URL.Query().Get("app")
	//	clientId := r.URL.Query().Get("id")
	//	realtimeConn := NewConnection(conn, appId, clientId)
	//
	//	log.Info("New connection for app:", appId)
	//
	//	GetConnectionStore().Put(appId, realtimeConn)
	//})

	//TODO: do not fail
	realtimeSub, err := redisClient.Subscribe(realtimeRedisSubject)
	if err != nil {
		log.Error(err)
		return nil
	}

	go func() {
		for {
			redisMsg, err := realtimeSub.ReceiveMessage()
			if err != nil {
				log.Error(err)
				continue
			}

			msg := redisMsg.Payload
			log.Info("Realtime service got message from redis, broadcasting:", msg)
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
	}()

	return server
}
