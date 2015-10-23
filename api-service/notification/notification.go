package notification

import (
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"github.com/nats-io/nats"
	"github.com/go-neutrino/neutrino/config"
)

var (
	qConn *nats.EncodedConn
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
	qAddr := config.Get(config.KEY_QUEUE_ADDR)
	n, e := nats.Connect(qAddr)

	if e != nil {
		panic(e)
	}

	conn, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	qConn = conn
}

func Notify(data models.JSON) {
	subj := config.Get(config.CONST_REALTIME_JOBS_SUBJ)
	str, err := data.String()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Publishing to queue subject: " + subj + " data: " + str)
	qConn.Publish(subj, data)
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
