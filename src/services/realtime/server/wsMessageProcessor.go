package server

import (
	"fmt"
	"strings"

	"github.com/neutrinoapp/neutrino/src/common"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/redis.v3"
)

type WsMessageProcessor struct {
	Interceptor     *wsInterceptor
	RedisClient     *redis.Client
	ClientProcessor messaging.MessageProcessor
	WsClient        *turnpike.Client
	broadcaster     *common.Broadcaster
}

func NewWsMessageProcessor(
	interceptor *wsInterceptor,
	redisClient *redis.Client,
	clientMessageProcessor messaging.MessageProcessor,
	wsClient *turnpike.Client,
) WsMessageProcessor {
	broadcaster := common.NewBroadcaster()
	p := WsMessageProcessor{interceptor, redisClient, clientMessageProcessor, wsClient, broadcaster}
	return p
}

func (p WsMessageProcessor) Process() {
	go func() {
		for {
			select {
			case m := <-p.Interceptor.OnMessage:
				msgType := m.messageType
				if msgType == turnpike.SUBSCRIBE {
					p.HandleSubscribe(m, m.msg.(*turnpike.Subscribe))
				} else if msgType == turnpike.PUBLISH {
					p.HandlePublish(m, m.msg.(*turnpike.Publish))
				}
			}
		}
	}()
}

func (p WsMessageProcessor) HandlePublish(im interceptorMessage, msg *turnpike.Publish) (interface{}, error) {
	if string(msg.Topic) == "wamp.session.on_leave" {
		args := msg.Arguments
		if len(args) == 0 {
			return nil, nil
		}

		if leavingSessionId, ok := args[0].(turnpike.ID); ok {
			log.Info("Broadcasting session leave:", leavingSessionId)
			p.broadcaster.Broadcast(leavingSessionId)
		}

		return nil, nil
	}

	if string(msg.Topic) == "wamp.session.on_join" {
		return nil, nil
	}

	if len(msg.Arguments) == 0 {
		return nil, nil
	}

	m, ok := msg.Arguments[0].(string)
	if !ok {
		m = models.String(msg.Arguments[0])
	}

	data, apiError := p.ClientProcessor.Process(m)
	if apiError != nil {
		log.Error(apiError)
	}

	return data, apiError
}

func (p WsMessageProcessor) HandleSubscribe(im interceptorMessage, msg *turnpike.Subscribe) {
	opts := models.SubscribeOptions{}
	err := models.Convert(msg.Options, &opts)
	if err != nil {
		log.Error(err)
		return
	}

	//remove the last part from 8139ed1ec39a467b96b0250dcf520749.todos.create.2882717310567
	topic := fmt.Sprintf("%v", msg.Topic)
	topicArguments := strings.Split(topic, ".")
	//uniqueTopicId := topicArguments[len(topicArguments)-1]

	baseTopic := messaging.BuildTopicArbitrary(topicArguments[:len(topicArguments)-1]...)
	//opts.BaseTopic = baseTopic
	opts.Topic = topic
	opts.BaseTopic = baseTopic
	//opts.ClientId = msg.Request
	//opts.TopicId = uniqueTopicId

	d := db.NewDbService()

	newValuesChan := make(chan map[string]interface{})
	err = d.Changes(opts.AppId, opts.Type, opts.Filter, newValuesChan)
	if err != nil {
		log.Error(err)
		return
	}

	messageBuilder := messaging.GetMessageBuilder()
	go func() {
		leaveChan := make(chan interface{})
		p.broadcaster.Subscribe(leaveChan)

		for {
			select {
			case val := <-newValuesChan:
				p.processDatabaseUpdate(val, messageBuilder, opts)
			case leaveVal := <-leaveChan:
				sessionId := leaveVal.(turnpike.ID)
				if im.sess.Id == sessionId {
					//TODO: newValuesChan seems to be automatically closed by the rethinkdb driver
					//investigate whether we need to do something else

					_, leaveChanOpened := <-leaveChan
					if leaveChanOpened {
						close(leaveChan)
					}

					p.broadcaster.Remove(leaveChan)
					return
				}
			}
		}
	}()
}

func (p WsMessageProcessor) processDatabaseUpdate(
	val map[string]interface{},
	messageBuilder messaging.MessageBuilder,
	opts models.SubscribeOptions,
) {

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

	//only emit messages with the same operation as the subscriber
	if dbOp != opts.Operation {
		return
	}

	var dbVal map[string]interface{}
	if dbOp == messaging.OP_CREATE {
		dbVal = newVal.(map[string]interface{})
	} else if dbOp == messaging.OP_DELETE {
		dbVal = oldVal.(map[string]interface{})
	} else {
		dbVal = newVal.(map[string]interface{})
	}

	log.Info("new_val:", newVal, "old_val:", oldVal, "db_val:", dbVal, "dbOp:", dbOp, "op:", opts.Operation)

	pld = pld.FromMap(dbVal)
	msg := messageBuilder.Build(
		dbOp,
		messaging.ORIGIN_API,
		pld,
		models.Options{},
		opts.Type,
		opts.AppId,
		"",
	)
	msg.Topic = opts.Topic

	log.Info("Publishing realtime data: ", val, opts.Topic, opts.BaseTopic, msg)

	publishArgs := []interface{}{msg}
	if opts.Topic != opts.BaseTopic {
		p.WsClient.Publish(opts.BaseTopic, publishArgs, nil)
	}

	p.WsClient.Publish(opts.Topic, publishArgs, nil)

}
