package notification

import (
	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/redis.v3"
)

var (
	redisClient         *redis.Client
	realtimeJobsSubject string
)

func init() {
	redisClient = client.GetNewRedisClient()

	realtimeJobsSubject = config.Get(config.CONST_REALTIME_JOBS_SUBJ)
}

func Notify(m messaging.Message) {
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

	log.Info("Publishing to queue subject: " + realtimeJobsSubject + " data: " + str)
	pubErr := redisClient.Publish(realtimeJobsSubject, str).Err()
	if pubErr != nil {
		log.Error(pubErr)
		return
	}
}
