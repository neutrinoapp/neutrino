package api

import (
	"github.com/gin-gonic/gin"
	"strings"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"github.com/go-neutrino/neutrino-core/db"
)

func authWithToken(c *gin.Context, userToken string) error {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if(jwt.GetSigningMethod("HS256") != token.Method){
			//TODO: "Invalid signing token algorithm."
			return nil, nil
		}

		return []byte(""), nil
	})

	c.Set("token", token)
	c.Set("user", token.Claims["user"])

	return err
}

func authWithMaster(c *gin.Context, key string) error {
	return nil
}

func authorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			//TODO: not authorized
			return
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 {
			//TODO: not authorized
			return
		}

		authType := authHeaderParts[0]
		authValue := authHeaderParts[1]

		var err error
		if authType == "Bearer" {
			err = authWithToken(c, authValue)
		} else if authType == "MasterKey" {
			err = authWithMaster(c, authValue)
		}

		if err != nil {
			//TODO: err
			return
		}

		c.Next()
	}
}

func defaultContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentTypeHeader := c.Request.Header.Get("Content-Type")
		if contentTypeHeader == "" {
			c.Request.Header.Set("Content-Type", "application/json")
		}

		c.Next()
	}
}

func injectAppMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appId := c.Param("appId")

		if appId != "" {
			//TODO: cache this
			d := db.NewAppsDbService(c.MustGet("user").(string))
			app, err := d.FindId(appId, nil)
			if err != nil {
				RestError(c, err)
				return
			}

			c.Set("app", JSON{}.FromMap(app))

		} else {
			RestError(c, "Invalid app id.")
		}
	}
}