package config

import (
	"github.com/spf13/viper"
)

const (
	KEY_MONGO_ADDR = "mongo-host"
	KEY_QUEUE_ADDR = "queue-host"

	KEY_REDIS_ADDR = "redis-addr"

	KEY_API_PORT      = "core-port"
	KEY_REALTIME_PORT = "realtime-port"

	KEY_BROKER_PORT = "broker-port"
	KEY_BROKER_HOST = "broker-host"

	CONST_REALTIME_JOBS_SUBJ = "realtime-jobs"
	CoNST_DEFAULT_REALM      = "default"
)

var c *viper.Viper

func setDefaults(v *viper.Viper) {
	v.SetDefault(KEY_MONGO_ADDR, "localhost:27017")
	v.SetDefault(KEY_REDIS_ADDR, "localhost:6379")
	v.SetDefault(KEY_QUEUE_ADDR, "nats://localhost:4222")

	v.SetDefault(KEY_API_PORT, ":5000")
	v.SetDefault(KEY_REALTIME_PORT, ":6000")
	v.SetDefault(KEY_BROKER_PORT, ":4000")
	v.SetDefault(KEY_BROKER_HOST, "ws://localhost")

	v.SetDefault(CONST_REALTIME_JOBS_SUBJ, "realtime-jobs")
}

func load() *viper.Viper {
	v := viper.New()

	//TODO: load from env variables

	setDefaults(v)

	return v
}

func init() {
	c = load()
}

func Raw() *viper.Viper {
	return c
}

func Get(k string) string {
	return c.GetString(k)
}
