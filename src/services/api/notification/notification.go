package notification

import (
	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"gopkg.in/redis.v3"
)

var (
	redisClient *redis.Client
)

func init() {
	redisClient = client.GetNewRedisClient()
}

func Notify(m messaging.Message) {
	if m.Topic == "" {
		topic := messaging.BuildTopic(m)
		m.Topic = topic
	}

	log.Info("Publishing to redis topic: "+config.CONST_REALTIME_JOBS_SUBJ+" data:", m)
	messageString, err := m.String()
	if err != nil {
		log.Error(err)
	}

	pubCmd := redisClient.Publish(config.CONST_REALTIME_JOBS_SUBJ, messageString)
	publishErr := pubCmd.Err()
	if publishErr != nil {
		log.Error(publishErr)
		return
	}
}
