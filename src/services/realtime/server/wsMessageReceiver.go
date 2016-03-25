package server

import (
	"fmt"
	"strings"

	"github.com/gngeorgiev/gowamp"
	"github.com/neutrinoapp/neutrino/src/common"
	"github.com/neutrinoapp/neutrino/src/common/db"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
	"gopkg.in/redis.v3"
)

type WsMessageReceiver struct {
	Interceptor     *wsInterceptor
	RedisClient     *redis.Client
	ClientProcessor messaging.MessageProcessor
	WsClient        *gowamp.Client
	broadcaster     *common.Broadcaster
}

func NewWsMessageReceiver(
	interceptor *wsInterceptor,
	redisClient *redis.Client,
	wsClient *gowamp.Client,
) WsMessageReceiver {
	broadcaster := common.NewBroadcaster()
	clientMessageProcessor := messaging.NewMessageProcessor()
	p := WsMessageReceiver{interceptor, redisClient, clientMessageProcessor, wsClient, broadcaster}
	return p
}

func (p WsMessageReceiver) Receive() {
	go func() {
		for {
			select {
			case m := <-p.Interceptor.OnMessage:
				msgType := m.messageType
				if msgType == gowamp.SUBSCRIBE {
					p.handleSubscribe(m, m.msg.(*gowamp.Subscribe))
				} else if msgType == gowamp.PUBLISH {
					p.handlePublish(m, m.msg.(*gowamp.Publish))
				}
			}
		}
	}()
}

func (p WsMessageReceiver) handlePublish(im interceptorMessage, msg *gowamp.Publish) (interface{}, error) {

	if string(msg.Topic) == "wamp.session.on_leave" {
		args := msg.Arguments
		if len(args) == 0 {
			return nil, nil
		}

		if leavingSessionId, ok := args[0].(gowamp.ID); ok {
			log.Info("Broadcasting session leave:", leavingSessionId)
			p.broadcaster.Broadcast(leavingSessionId)
			return nil, nil
		}

		log.Info("Leave:", args)
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

	var clientMessage messaging.Message
	converError := clientMessage.FromString(m)
	if converError != nil {
		log.Error(converError)
		return nil, converError
	}

	if clientMessage.Origin == messaging.ORIGIN_API {
		return nil, nil
	}

	log.Info("Websocket message received:", clientMessage)
	data, apiError := p.ClientProcessor.Process(clientMessage)
	if apiError != nil {
		log.Error(apiError)
	}

	return data, apiError
}

func (p WsMessageReceiver) handleSubscribe(im interceptorMessage, msg *gowamp.Subscribe) {
	opts := models.SubscribeOptions{}
	err := models.Convert(msg.Options, &opts)
	if err != nil {
		log.Error(err)
		return
	}

	topic := fmt.Sprintf("%v", msg.Topic)
	topicArguments := strings.Split(topic, ".")

	opts.Topic = topic
	if opts.IsSpecial() {
		baseTopic := messaging.BuildTopicArbitrary(topicArguments[:len(topicArguments)-1]...)
		opts.BaseTopic = baseTopic
	} else {
		opts.BaseTopic = topic
	}

	log.Info("Listening for changefeed:", opts)
	d := db.NewDbService()
	newValuesChan := make(chan map[string]interface{})

	//for create and delete we want to listen for changes for the whole collection, for update only for the specific item
	if opts.Operation == messaging.OP_CREATE || opts.Operation == messaging.OP_DELETE {
		err = d.Changes(opts.AppId, opts.Type, opts.Filter, newValuesChan)
	} else if opts.Operation == messaging.OP_UPDATE {
		opts.ItemId = topicArguments[len(topicArguments)-1]
		err = d.ChangesId(opts.ItemId, newValuesChan)
	}

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
				p.processDatabaseUpdate(val, &messageBuilder, opts)
			case leaveVal := <-leaveChan:
				sessionId := leaveVal.(gowamp.ID)
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

func (p WsMessageReceiver) processDatabaseUpdate(
	val map[string]interface{},
	builder *messaging.MessageBuilder,
	opts models.SubscribeOptions,
) {
	messageBuilder := *builder

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
		if newVal == nil {
			log.Error("Invalid message:", val)
			return
		}

		dbVal = newVal.(map[string]interface{})
	}

	pld = pld.FromMap(dbVal)
	pld = utils.BlacklistFields([]string{db.TYPE_FIELD, db.APP_ID_FIELD}, pld)

	msg := messageBuilder.Build(
		dbOp,
		messaging.ORIGIN_API,
		pld,
		models.Options{},
		opts.Type,
		opts.AppId,
		"", //TODO: token?
	)
	msg.Topic = opts.Topic

	log.Info("Publishing database feed: ", strings.ToUpper(msg.Operation), opts, val, msg)

	publishArgs := []interface{}{msg}
	p.WsClient.Publish(opts.Topic, publishArgs, nil)
}
