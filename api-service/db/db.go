package db

import (
	"github.com/spf13/viper"
)

var config *viper.Viper
func Initialize(c *viper.Viper) {
	if config != nil {
		panic("Initialize must be called once")
	}

	config = c
}