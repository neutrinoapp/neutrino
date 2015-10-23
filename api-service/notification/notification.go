package notification

import (
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/client"
)

var (
	natsClient *client.NatsClient
	queueSubject string
)

type op string
type origin string

const (
	OP_UPDATE op = "update"
	OP_CREATE op = "create"
	OP_DELETE op = "delete"

	ORIGIN_API    origin = "api"
	ORIGIN_CLIENT origin = "client"
)

func init() {
	natsClient = client.NewNatsClient(config.Get(config.KEY_QUEUE_ADDR))
	queueSubject = config.Get(config.CONST_REALTIME_JOBS_SUBJ)
}

func Notify(data models.JSON) {
	str, err := data.String()
	if err != nil {
		log.Error(err)
		return
	}

	conn := natsClient.GetConnection()
	if conn != nil {
		log.Info("Publishing to queue subject: " + queueSubject + " data: " + str)
		conn.Publish(queueSubject, data)
	} else {
		log.Info("Queue service not available, realtime updates will not arrive.")
	}
}

func Build(o op, og origin, pld interface{}, opts models.JSON, t string) models.JSON {
	return models.JSON{
		"op":      o,
		"origin":  og,
		"options": opts,
		"type": t,
		"payload": pld,
	}
}
