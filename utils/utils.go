package utils

import (
	"github.com/twinj/uuid"
	"strings"
	"github.com/gin-gonic/gin"
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

func GetBody(c *gin.Context) map[string]interface{} {
	var res map[string]interface{}
	c.Bind(&res)
	return res
}