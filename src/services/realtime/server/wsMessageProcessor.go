package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/redis.v3"
)

type WsMessageProcessor struct {
	Interceptor     wsInterceptor
	RedisClient     *redis.Client
	ClientProcessor clientMessageProcessor
	WsClient        *turnpike.Client
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
	if len(msg.Arguments) == 0 {
		return
	}

	m, ok := msg.Arguments[0].(string)
	if !ok {
		return
	}

	apiError := p.ClientProcessor.Process(m)
	if apiError != nil {
		log.Error(apiError)
	}

	topic := string(msg.Topic)
	log.Info("Sending out special messages:", topic)
	msgRaw := models.JSON{}
	err := msgRaw.FromString([]byte(m))
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(msgRaw)
	if payload, ok := msgRaw["pld"].(map[string]interface{}); ok {
		clientIds := p.RedisClient.SMembers(topic).Val()
		log.Info(clientIds)
		for _, clientId := range clientIds {
			filterString := p.RedisClient.HGet(clientId, "filter").Val()
			filter := models.JSON{}
			filter.FromString([]byte(filterString))
			log.Info(filter)

			passes := true
			for k, v := range filter {
				log.Info(payload, k, v)
				if payload[k] != v {
					passes = false
					break
				}
			}

			if passes {
				topic := p.RedisClient.HGet(clientId, "topic").Val()
				log.Info("Publishing to special topic: ", topic, m)
				err := p.WsClient.Publish(topic, []interface{}{msgRaw}, nil)
				if err != nil {
					log.Error(err)
					continue
				}
			}

			log.Info(filter)
		}
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
		clientId := strconv.FormatUint(uint64(msg.Request), 10)

		baseTopic := messaging.BuildTopicArbitrary(topicArguments[:len(topicArguments)-1]...)
		opts.BaseTopic = baseTopic
		opts.Topic = topic
		opts.ClientId = msg.Request
		opts.TopicId = uniqueTopicId

		p.RedisClient.SAdd(baseTopic, clientId)

		p.RedisClient.HSet(clientId, "baseTopic", opts.BaseTopic)
		p.RedisClient.HSet(clientId, "topic", opts.Topic)
		p.RedisClient.HSet(clientId, "clientId", clientId)
		p.RedisClient.HSet(clientId, "topicId", opts.TopicId)
		p.RedisClient.HSet(clientId, "filter", models.String(opts.Filter))
	}
}
