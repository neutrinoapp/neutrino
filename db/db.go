package db

import "github.com/go-neutrino/go-env-config"

var config envconfig.Config
func Initialize(c envconfig.Config) {
	if config != nil {
		panic("Initialize must be called once")
	}

	config = c
}