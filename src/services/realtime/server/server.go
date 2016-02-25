package server

import (
	"net/http"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/redis.v3"
)

var (
	realtimeRedisSubject string
	redisClient          *redis.Client
	messageProcessor     messaging.MessageProcessor
)

func init() {
	realtimeRedisSubject = config.Get(config.CONST_REALTIME_JOBS_SUBJ)
	redisClient = client.GetNewRedisClient()
	messageProcessor = NewClientMessageProcessor()
}

type wsInterceptor struct {
	m chan turnpike.Message
}

func (i *wsInterceptor) Intercept(session turnpike.Session, msg *turnpike.Message) {
	i.m <- *msg
}

func Initialize() (*http.Server, error) {
	interceptor := &wsInterceptor{
		m: make(chan turnpike.Message),
	}

	r := turnpike.Realm{}
	r.Interceptor = interceptor

	realms := map[string]turnpike.Realm{}
	realms[config.CONST_DEFAULT_REALM] = r
	wsServer, err := turnpike.NewWebsocketServer(realms)
	if err != nil {
		return nil, err
	}

	wsServer.Upgrader.CheckOrigin = func(r *http.Request) bool {
		//allow connections from any origin
		return true
	}

	go func() {
		for {
			select {
			case m := <-interceptor.m:
				switch msg := m.(type) {
				case *turnpike.Subscribe:
					//redisClient.Set(msg.Topic, 1, 0)
				case *turnpike.Publish:
					if len(msg.Arguments) > 0 {
						m, ok := msg.Arguments[0].(string)
						if ok {
							apiError := messageProcessor.Process(messaging.MESSAGE_TYPE_STRING, m)
							if apiError != nil {
								log.Error(apiError)
							}
						}
					}
				}
			}
		}
	}()

	server := &http.Server{
		Handler: wsServer,
		Addr:    config.Get(config.KEY_REALTIME_PORT),
	}

	////TODO: do not fail
	//realtimeRedisSub, err := redisClient.Subscribe(realtimeRedisSubject)
	//if err != nil {
	//	return nil, err
	//}
	//
	//go func() {
	//	for {
	//		redisMsg, err := realtimeRedisSub.ReceiveMessage()
	//		if err != nil {
	//			log.Error(err)
	//			continue
	//		}
	//
	//		msg := redisMsg.Payload
	//		log.Info("Realtime service got message from redis, broadcasting:", msg)
	//		var m messaging.Message
	//		if err := m.FromString(msg); err != nil {
	//			log.Error(err)
	//			continue
	//		}
	//
	//		appId := m.App
	//		if appId == "" {
	//			log.Error("No appId provided with realtime notification.", m)
	//			continue
	//		}
	//
	//		apiError := messageProcessor.Process(messaging.MESSAGE_TYPE_STRING, m)
	//		if apiError != nil {
	//			log.Error(apiError)
	//		}
	//
	//		c, err := wsServer.GetLocalClient(config.CONST_DEFAULT_REALM, nil)
	//		if err != nil {
	//			log.Error(err)
	//			continue
	//		}
	//
	//		topic := appId + "_" + m.Type
	//		log.Info("Trying to publish message to topic:", topic, msg)
	//		publishErr := c.Publish(topic, []interface{}{msg}, nil)
	//		if publishErr != nil {
	//			log.Error("Failed publishing message to topic:", topic, msg)
	//		}
	//	}
	//}()

	return server, nil
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

//log.Info("Broadcasting:", msg, "to", appId)
//connsForApp := GetConnectionStore().Get(appId)
//for _, conn := range connsForApp {
//	connClientId := conn.GetClientId()
//	if m.Origin == messaging.ORIGIN_CLIENT &&
//		m.Options != nil &&
//		m.Options["clientId"] != nil &&
//		m.Options["clientId"] == connClientId {
//
//		log.Info("Skipping broadcast to client", connClientId, "has same id.")
//		continue
//	}
//
//	conn.Broadcast(msg)
//}
