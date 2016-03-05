package config

import (
	"github.com/spf13/viper"
)

const (
	KEY_RETHINK_ADDR = "rethink-host"
	KEY_QUEUE_ADDR   = "queue-host"

	KEY_REDIS_ADDR = "redis-addr"

	KEY_API_PORT = "api-port"
	KEY_API_ADDR = "api-addr"

	KEY_REALTIME_PORT = "realtime-port"
	KEY_REALTIME_ADDR = "realtime-addr"

	CONST_REALTIME_JOBS_SUBJ = "realtime-jobs"
	CONST_DEFAULT_REALM      = "default"
)

var c *viper.Viper

func setDefaults(v *viper.Viper) {
	v.SetDefault(KEY_RETHINK_ADDR, "localhost:28015")
	v.SetDefault(KEY_REDIS_ADDR, "localhost:6379")
	v.SetDefault(KEY_QUEUE_ADDR, "nats://localhost:4222")

	v.SetDefault(KEY_API_PORT, ":5000")
	v.SetDefault(KEY_API_ADDR, "http://localhost"+v.GetString(KEY_API_PORT)+"/v1/")

	v.SetDefault(KEY_REALTIME_PORT, ":6000")
	v.SetDefault(KEY_REALTIME_ADDR, "ws://localhost"+v.GetString(KEY_REALTIME_PORT))

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
