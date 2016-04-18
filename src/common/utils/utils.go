package utils

import (
	"strings"

	"os"
	"os/signal"

	"syscall"

	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/twinj/uuid"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

func GetCleanUUID() string {
	return strings.Replace(GetUUID(), "-", "", -1)
}

func BlacklistFields(fields []string, data interface{}) map[string]interface{} {
	if obj, ok := data.(map[string]interface{}); ok {
		for _, k := range fields {
			delete(obj, k)
		}

		return obj
	} else if obj, ok := data.(models.JSON); ok {
		for _, k := range fields {
			delete(obj, k)
		}

		return obj
	}

	log.Info("Invalid object to blacklist", data)
	return make(map[string]interface{})
}

func Recover() {
	e := recover()
	if e != nil {
		log.Error(e)
	}
}

func ListenSignals() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for si := range s {
			log.Warn("OS signal received:", si)
			os.Exit(0)
		}
	}()
}
