package server

import (
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/jcelliott/turnpike.v2"
	"gopkg.in/redis.v3"
)

type JobsMessageProcessor struct {
	Client   *redis.Client
	WsClient *turnpike.Client
}

func (p JobsMessageProcessor) Process() {
	go func() {
		psub, err := p.Client.Subscribe(config.CONST_REALTIME_JOBS_SUBJ)
		if err != nil {
			log.Error(err)
			return
		}

		for {
			msg, err := psub.ReceiveMessage()
			if err != nil {
				log.Error(err)
				continue
			}

			log.Info("Processing redis message:", msg.Payload)

			var m messaging.Message
			parseError := m.FromString(msg.Payload)
			if parseError != nil {
				log.Error(parseError)
				return
			}

			m.Origin = messaging.ORIGIN_API
			publishErr := p.WsClient.Publish(m.Topic, []interface{}{m}, nil)

			if publishErr != nil {
				log.Error(publishErr)
				return
			}
		}
	}()
}
