package notification

import (
	"github.com/go-neutrino/neutrino/client"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/messaging"
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
	str, err := m.Serialize().String()
	if err != nil {
		log.Error(err)
		return
	}

	conn := natsClient.GetConnection()
	if conn != nil {
		log.Info("Publishing to queue subject: " + queueSubject + " data: " + str)
		conn.Publish(queueSubject, m)
	} else {
		log.Info("Queue service not available, realtime updates will not arrive.")
	}
}
