package server

import (
	"net/http"
	"sync"

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
	natsClient           *client.NatsClient
	messageProcessor     messaging.MessageProcessor
)

func init() {
	realtimeRedisSubject = config.Get(config.CONST_REALTIME_JOBS_SUBJ)
	redisClient = client.GetNewRedisClient()
	messageProcessor = NewClientMessageProcessor()
	natsClient = client.NewNatsClient(config.Get(config.KEY_QUEUE_ADDR))
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
					//TODO: put special subscriptions into redis, e.g. filtering
				case *turnpike.Publish:
					if len(msg.Arguments) > 0 {
						m, ok := msg.Arguments[0].(string)
						if ok {
							apiError := messageProcessor.Process(m)
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

	c, err := wsServer.GetLocalClient(config.CONST_DEFAULT_REALM, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var mu sync.Mutex
	go func() {
		err := natsClient.Subscribe(config.CONST_REALTIME_JOBS_SUBJ, func(mStr string) {
			log.Info("Processing nats message:", mStr)

			var m messaging.Message
			err := m.FromString(mStr)
			if err != nil {
				log.Error(err)
				return
			}

			mu.Lock()
			publishErr := c.Publish(m.Topic, []interface{}{mStr}, nil)
			mu.Unlock()

			if publishErr != nil {
				log.Error(publishErr)
				return
			}
		})

		if err != nil {
			log.Error(err)
			return
		}
	}()

	return server, nil
}
