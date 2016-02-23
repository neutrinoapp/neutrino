package notification

import (
	"github.com/go-neutrino/neutrino/src/common/client"
	"github.com/go-neutrino/neutrino/src/common/config"
	"github.com/go-neutrino/neutrino/src/common/log"
	"github.com/go-neutrino/neutrino/src/common/messaging"
)

var (
	natsClient   *client.NatsClient
	queueSubject string
)

func init() {
	natsClient = client.NewNatsClient(config.Get(config.KEY_QUEUE_ADDR))
	queueSubject = config.Get(config.CONST_REALTIME_JOBS_SUBJ)
}

func Notify(m messaging.Message) {
	conn := natsClient.GetConnection()
	if conn != nil {
		model, err := m.ToJson()
		if err != nil {
			log.Error(err)
			return
		}

		str, err := model.String()
		if err != nil {
			log.Error(err)
			return
		}

		log.Info("Publishing to queue subject: " + queueSubject + " data: " + str)
		conn.Publish(queueSubject, model)
	} else {
		log.Info("Queue service not available, realtime updates will not arrive.")
	}
}
