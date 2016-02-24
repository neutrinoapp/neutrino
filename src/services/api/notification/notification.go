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
	wsClient    *client.WebSocketClient
)

func init() {
	redisClient = client.GetNewRedisClient()
	wsClient = client.NewWebsocketClient([]string{config.CONST_DEFAULT_REALM})
	go wsClient.Connect()
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

	topic := messaging.GetTopic(m)
	log.Info("Publishing topic: " + topic + " data: " + str)
	//TODO: throws error if the connection was lost
	publishErr := wsClient.GetConnection().Publish(topic, []interface{}{str}, nil)
	if publishErr != nil {
		log.Error(publishErr)
		return
	}
}
