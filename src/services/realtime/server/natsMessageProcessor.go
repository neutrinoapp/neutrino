package server

import (
	"sync"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/jcelliott/turnpike.v2"
)

type NatsMessageProcessor struct {
	Client   *client.NatsClient
	WsClient *turnpike.Client
}

func (p NatsMessageProcessor) Process() {
	var mu sync.Mutex
	go func() {
		err := p.Client.Subscribe(config.CONST_REALTIME_JOBS_SUBJ, func(mStr string) {
			log.Info("Processing nats message:", mStr)

			var m messaging.Message
			err := m.FromString(mStr)
			if err != nil {
				log.Error(err)
				return
			}

			mu.Lock()
			publishErr := p.WsClient.Publish(m.Topic, []interface{}{mStr}, nil)
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
