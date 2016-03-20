package utils

import (
	"strings"

	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/twinj/uuid"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

func GetCleanUUID() string {
	return strings.Replace(GetUUID(), "-", "", -1)
}

func BlacklistFields(fields []string, obj models.JSON) models.JSON {
	for _, k := range fields {
		delete(obj, k)
	}

	return obj
}
