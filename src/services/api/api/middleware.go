package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"gopkg.in/dgrijalva/jwt-go.v2"
)

const (
	CONTEXT_HEADER_OPTIONS = "NeutrinoOptions"
	CONTEXT_EXPRESSION     = "Expression"
)

type apiUser struct {
	Email, Key    string
	Master, InApp bool
}

func authWithToken(c *gin.Context, userToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			//TODO: "Invalid signing token algorithm."
			return nil, nil
		}

		//TODO: cache this
		//tokenSecretRecord, accountSecretError := db.NewSystemDbService().FindId("accountSecret")
		//
		//if accountSecretError != nil {
		//	//we probably do not have such collection. Use a default secret and warn.
		//	log.Info("Account secret error: ", accountSecretError)
		//	tokenSecretRecord = models.JSON{
		//		"value": "",
		//	}
		//}

		//tokenSecret := tokenSecretRecord["value"].(string)

		return []byte(""), nil
	})

	return token, err
}

func authWithMaster(c *gin.Context, key string) (string, error) {
	//d := db.NewDbService(db.DATABASE_NAME, db.USERS_TABLE)
	//res, err := d.Query().Filter(models.JSON{
	//	"masterKey": key,
	//}).Nth(0).Run(d.GetSession())
	//if err != nil {
	//	return "", err
	//}
	//
	//if res == nil || res["user"] == nil {
	//	return "", BuildError("Invalid master key")
	//}
	//
	//return res["user"].(string), nil
	//TODO:
	return "", nil
}

func authorizeMiddleware(stop bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader != "" {
			authHeaderParts := strings.SplitN(authHeader, " ", 2)
			if len(authHeaderParts) != 2 {
				log.Error(RestErrorUnauthorized(c))
				return
			}

			authType := strings.ToLower(authHeaderParts[0])
			authValue := authHeaderParts[1]

			//TODO: authorization for master token, master key, normal token, app id only
			var token *jwt.Token
			var err error
			user := &apiUser{}
			if authType == "bearer" {
				token, err = authWithToken(c, authValue)
				if err == nil {
					user.Email = token.Claims["user"].(string)
					user.InApp = token.Claims["inApp"].(bool)
					user.Master = !user.InApp //we can use the token instead of a master key
					user.Key = authValue
				}
			} else if authType == "masterkey" {
				email, err := authWithMaster(c, authValue)
				if err == nil {
					user.Email = email
					user.InApp = false
					user.Master = true
					user.Key = authValue
				}
			} else {
				c.Next()
				return
			}

			c.Set("user", user)

			if err != nil {
				log.Error(RestError(c, err))
				return
			}

			c.Next()
		} else {
			if !stop {
				c.Next()
			} else {
				log.Error(RestErrorUnauthorized(c))
			}
		}
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
		appId := c.Param("appId")
		log.Info(appId)
		if appId == "" {
			log.Error(RestError(c, "Invalid app id."))
		} else {
			c.Next()
		}
	}
}

func validateAppPermissionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := ApiUser(c)
		if !user.Master {
			log.Error(RestErrorUnauthorized(c), user)
			c.Next()
			return
		}
	}
}

func validateAppOperationsAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := ApiUser(c)
		if user.InApp && !user.Master {
			log.Error(RestErrorUnauthorized(c), user)
			return
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, NeutrinoOptions")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func processHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		optionsHeader := c.Request.Header.Get(CONTEXT_HEADER_OPTIONS)
		if optionsHeader == "" {
			optionsHeader = "{}"
		}

		var options models.Options
		options.FromString(optionsHeader)

		if options.Notify == nil {
			notify := true
			options.Notify = &notify
		}

		if options.Filter == nil {
			options.Filter = models.JSON{}
		}

		if options.Origin == "" {
			options.Origin = messaging.ORIGIN_API
		}

		log.Info("Request options:", optionsHeader)

		c.Set(CONTEXT_HEADER_OPTIONS, options)
		c.Next()
	}
}

func parseExpressionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//query := c.Request.URL.Query()
		//g, err := expression.ParseExpressionGroup(query)
		//if err != nil {
		//	log.Error(err)
		//	c.Next()
		//	return
		//}

		//log.Info("Expression: ", g)
		//
		//c.Set(CONTEXT_EXPRESSION, g)
	}
}
