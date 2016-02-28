package server

import (
	"net/http"
	"sync"

	"encoding/json"

	"fmt"

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
	wsServer, server, err := handlerWebSocketServer()
	if err != nil {
		return nil, err
	}

	c, err := wsServer.GetLocalClient(config.CONST_DEFAULT_REALM, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	handleNatsConnection(c)
	handleRpc(c)

	return server, nil
}

func handleRpc(c *turnpike.Client) {
	getArgs := func(args []interface{}) (messaging.Message, *client.ApiClient, error) {
		var m messaging.Message

		b, err := json.Marshal(args[0])
		if err != nil {
			return m, nil, err
		}

		err = json.Unmarshal(b, &m)
		if err != nil {
			return m, nil, err
		}

		c := client.NewApiClientCached(m.App)
		c.Token = m.Token
		return m, c, nil
	}

	dataRead := func(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
		m, c, err := getArgs(args)
		if err != nil {
			log.Error(err)
			return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
		}

		var clientResult interface{}
		if id, ok := m.Payload["_id"].(string); ok {
			clientResult, err = c.GetItem(m.Type, id)
		} else {
			clientResult, err = c.GetItems(m.Type)
		}

		if err != nil {
			log.Error(err)
			return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
		}

		return &turnpike.CallResult{Args: []interface{}{clientResult}}
	}

	dataCreate := func(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
		m, c, err := getArgs(args)
		if err != nil {
			log.Error(err)
			return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
		}

		resp, err := c.CreateItem(m.Type, m.Payload)
		if err != nil {
			log.Error(err)
			return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
		}

		return &turnpike.CallResult{Args: []interface{}{resp}}
	}

	dataRemove := func(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
		m, c, err := getArgs(args)
		if err != nil {
			log.Error(err)
			return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
		}

		id, ok := m.Payload["_id"].(string)
		if !ok {
			return &turnpike.CallResult{Err: turnpike.URI(fmt.Sprintf("Incorrect payload, %v", m.Payload))}
		}

		_, err = c.DeleteItem(m.Type, id)
		if err != nil {
			log.Error(err)
			return &turnpike.CallResult{Err: turnpike.URI(err.Error())}
		}

		return &turnpike.CallResult{Args: []interface{}{id}}
	}

	c.BasicRegister("data.read", dataRead)
	c.BasicRegister("data.create", dataCreate)
	c.BasicRegister("data.remove", dataRemove)
}

func handleNatsConnection(c *turnpike.Client) {
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
}

func handlerWebSocketServer() (*turnpike.WebsocketServer, *http.Server, error) {
	interceptor := &wsInterceptor{
		m: make(chan turnpike.Message),
	}

	r := turnpike.Realm{}
	r.Interceptor = interceptor

	realms := map[string]turnpike.Realm{}
	realms[config.CONST_DEFAULT_REALM] = r
	wsServer, err := turnpike.NewWebsocketServer(realms)
	if err != nil {
		return nil, nil, err
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

	return wsServer, server, nil
}
