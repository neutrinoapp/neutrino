package client

import (
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"gopkg.in/redis.v3"
)

func GetNewRedisClient() *redis.Client {
	redisAddr := config.Get(config.KEY_REDIS_ADDR)
	//TODO: check if will reconnect automatically on lost connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	log.Info("Connected to redis client on:", redisAddr)
	return redisClient
}
