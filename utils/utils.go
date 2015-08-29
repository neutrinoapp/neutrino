package utils

import (
	"github.com/twinj/uuid"
	"strings"
	"github.com/ant0ine/go-json-rest/rest"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

func GetCleanUUID() string {
	return strings.Replace(GetUUID(), "-", "", -1)
}

func WhitelistFields(fields []string, obj map[string]interface{}) map[string]interface{}{
	result := make(map[string]interface{})

	for _, k := range fields {
		result[k] = obj[k]
	}

	return result
}

func GetBody(r *rest.Request) map[string]interface{} {
	var res map[string]interface{}
	r.DecodeJsonPayload(&res)
	return res
}