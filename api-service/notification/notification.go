package notification

import (
	"github.com/go-neutrino/neutrino-config"
	"github.com/go-neutrino/neutrino-core/log"
	"github.com/go-neutrino/neutrino-core/models"
	"github.com/nats-io/nats"
	"github.com/spf13/viper"
)

var (
	qConn *nats.EncodedConn
	c     *viper.Viper
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
	c = nconfig.Load()
	qAddr := c.GetString(nconfig.KEY_QUEUE_ADDR)
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
	subj := c.GetString(nconfig.CONST_REALTIME_JOBS_SUBJ)
	log.Info("Publishing to queue subject: " + subj + " data: " + data.String())
	qConn.Publish(subj, data)
}

func Build(o op, og origin, pld interface{}, opts models.JSON) models.JSON {
	return models.JSON{
		"op":      o,
		"origin":  og,
		"options": opts,
		"payload": pld,
	}
}
