package server

import (
	"fmt"
	"strings"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
	"github.com/neutrinoapp/neutrino/src/services/api/notification"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/redis.v3"
)

type WsMessageProcessor struct {
	Interceptor     *wsInterceptor
	RedisClient     *redis.Client
	ClientProcessor messaging.MessageProcessor
	WsClient        *turnpike.Client
	NatsClient      *client.NatsClient
}

func (p WsMessageProcessor) Process() {
	go func() {
		for {
			select {
			case m := <-p.Interceptor.m:
				switch msg := m.(type) {
				case *turnpike.Subscribe:
					p.HandleSubscribe(msg)
				case *turnpike.Publish:
					p.HandlePublish(msg)
				case *turnpike.Goodbye:
					p.HandleGoodbye(msg)
				}
			}
		}
	}()
}

func (p WsMessageProcessor) HandleGoodbye(msg *turnpike.Goodbye) {
	clientId := string(msg.Request)
	p.RedisClient.Del(clientId)
}

func (p WsMessageProcessor) HandlePublish(msg *turnpike.Publish) {
	if string(msg.Topic) == "wamp.session.on_join" {
		return
	}

	if len(msg.Arguments) == 0 {
		return
	}

	m, ok := msg.Arguments[0].(string)
	if !ok {
		m = models.String(msg.Arguments[0])
	}

	apiError := p.ClientProcessor.Process(m)
	if apiError != nil {
		log.Error(apiError)
	}
}

func (p WsMessageProcessor) HandleSubscribe(msg *turnpike.Subscribe) {
	opts := models.SubscribeOptions{}
	err := models.Convert(msg.Options, &opts)
	if err != nil {
		log.Error(err)
		return
	}

	if opts.IsSpecial() {
		//remove the last part from 8139ed1ec39a467b96b0250dcf520749.todos.create.2882717310567
		topic := fmt.Sprintf("%v", msg.Topic)
		topicArguments := strings.Split(topic, ".")
		uniqueTopicId := topicArguments[len(topicArguments)-1]
		//clientId := strconv.FormatUint(uint64(msg.Request), 10)

		baseTopic := messaging.BuildTopicArbitrary(topicArguments[:len(topicArguments)-1]...)
		opts.BaseTopic = baseTopic
		opts.Topic = topic
		opts.ClientId = msg.Request
		opts.TopicId = uniqueTopicId

		d := db.NewTypeDbService(opts.AppId, opts.Type)

		newValuesChan := make(chan map[string]interface{})
		err := d.Changes(opts.Filter, newValuesChan)
		if err != nil {
			log.Error(err)
			return
		}

		messageBuilder := messaging.GetMessageBuilder()
		go func() {
			for {
				select {
				case val := <-newValuesChan:
					pld := models.JSON{}
					var dbOp string
					newVal := val["new_val"]
					oldVal := val["old_val"]
					if newVal != nil && oldVal == nil {
						dbOp = messaging.OP_CREATE
					} else if newVal == nil && oldVal != nil {
						dbOp = messaging.OP_DELETE
					} else {
						dbOp = messaging.OP_UPDATE
					}

					if dbOp != opts.Operation {
						//only emit messages with the same operation as the subscriber
						continue
					}

					var dbVal map[string]interface{}
					if dbOp == messaging.OP_CREATE {
						dbVal = newVal.(map[string]interface{})
					} else if dbOp == messaging.OP_DELETE {
						dbVal = oldVal.(map[string]interface{})
					} else {
						dbVal = newVal.(map[string]interface{})
					}

					pld = pld.FromMap(dbVal)
					notify := false
					msg := messageBuilder.Build(
						dbOp,
						messaging.ORIGIN_API,
						pld,
						models.Options{
							Notify: &notify,
						},
						opts.Type,
						opts.AppId,
						"",
					)
					msg.Topic = topic

					log.Info("Publishing filtered data: ", val, opts.Topic, msg)
					notification.Notify(msg)
				}
			}
		}()
	}
}
