package webUtils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/models"
)

func GetBody(c *gin.Context) models.JSON {
	var res models.JSON
	c.Bind(&res)
	return res
}

func OK(c *gin.Context) {
	c.String(http.StatusOK, "")
}
