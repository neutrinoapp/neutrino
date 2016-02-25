package notification

import (
	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
)

var (
	natsClient *client.NatsClient
)

func init() {
	natsClient = client.NewNatsClient(config.Get(config.KEY_QUEUE_ADDR))
}

func Notify(m messaging.Message) {
	topic := messaging.BuildTopic(m)
	m.Topic = topic

	log.Info("Publishing to nats topic: "+config.CONST_REALTIME_JOBS_SUBJ+" data:", m)
	publishErr := natsClient.Publish(config.CONST_REALTIME_JOBS_SUBJ, m)
	if publishErr != nil {
		log.Error(publishErr)
		return
	}
}
