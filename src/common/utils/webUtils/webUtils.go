package webUtils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBody(c *gin.Context) map[string]interface{} {
	var res map[string]interface{}
	c.Bind(&res)
	return res
}

func OK(c *gin.Context) {
	c.String(http.StatusOK, "")
}
