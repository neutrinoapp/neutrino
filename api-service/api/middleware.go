package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/api-service/db"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"strings"
)

type GetAppFunc func() (models.JSON, error)

func authWithToken(c *gin.Context, userToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			//TODO: "Invalid signing token algorithm."
			return nil, nil
		}

		//TODO: cache this
		tokenSecretRecord, err := db.NewSystemDbService().FindId("accountSecret", nil)

		if err != nil {
			//we probably do not have such collection. Use a default secret and warn.
			log.Error("Account secret error: ", err)
			tokenSecretRecord = models.JSON{
				"value": "",
			}
		}

		tokenSecret := tokenSecretRecord["value"].(string)

		return []byte(tokenSecret), nil
	})

	return token, err
}

func authWithMaster(c *gin.Context, key string) (*jwt.Token, error) {
	return nil, nil
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

		authType := strings.ToLower(authHeaderParts[0])
		authValue := authHeaderParts[1]

		//TODO: authorization for master token, master key, normal token, app id only
		var token *jwt.Token
		var err error
		if authType == "bearer" {
			token, err = authWithToken(c, authValue)
		} else if authType == "masterkey" {
			token, err = authWithMaster(c, authValue)
		} else {
			c.Next()
			return
		}

		c.Set("user", token.Claims["user"])
		c.Set("inApp", token.Claims["inApp"])
		c.Set("token", authValue)

		if err != nil {
			log.Error(RestError(c, err))
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

func validateAppMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: do we always need the app injected?
		appId := c.Param("appId")

		if appId == "" {
			log.Error(RestError(c, "Invalid app id."))
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
