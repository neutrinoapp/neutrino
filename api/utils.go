package api

import (
	"github.com/labstack/echo"
)

func GetBody(c *echo.Context) (map[string]interface{}, error) {
	b := make(map[string]interface{})
	err := c.Bind(&b)

	return b, err
}

