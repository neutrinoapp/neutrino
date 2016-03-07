package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/expression"
	"github.com/neutrinoapp/neutrino/src/common/models"
)

var initialized bool

func initMiddleware(e *gin.Engine) {
	e.Use(defaultContentTypeMiddleware(), CORSMiddleware())
}

func initRoutes(e *gin.Engine) {
	authController := &AuthController{}
	appController := &ApplicationController{}
	typesController := &TypesController{}

	v1 := e.Group("/v1", processHeadersMiddleware())
	{
		v1.POST("/login", authController.LoginUserHandler)
		v1.POST("/register", authController.RegisterUserHandler)

		appGroup := v1.Group("/app", authorizeMiddleware(false))
		{
			appGroup.POST("", appController.CreateApplicationHandler, validateAppOperationsAuthorizationMiddleware())
			appGroup.GET("", appController.GetApplicationsHandler, validateAppOperationsAuthorizationMiddleware())

			appIdGroup := appGroup.Group("/:appId", validateAppMiddleware())
			{
				appIdGroup.GET("", appController.GetApplicationHandler, validateAppPermissionsMiddleware())
				appIdGroup.DELETE("", appController.DeleteApplicationHandler, validateAppPermissionsMiddleware())
				appIdGroup.PUT("", appController.UpdateApplicationHandler, validateAppPermissionsMiddleware())

				appIdGroup.POST("/register", authController.AppRegisterUserHandler)
				appIdGroup.POST("/login", authController.AppLoginUserHandler)

				dataGroup := appIdGroup.Group("/data", authorizeMiddleware(true), parseExpressionsMiddleware())
				{
					dataGroup.GET("", typesController.GetTypesHandler)
					dataGroup.DELETE("/:typeName", typesController.DeleteType)
					dataGroup.POST("/:typeName", typesController.InsertInTypeHandler)
					dataGroup.GET("/:typeName", typesController.GetTypeDataHandler)
					dataGroup.GET("/:typeName/:itemId", typesController.GetTypeItemById)
					dataGroup.PUT("/:typeName/:itemId", typesController.UpdateTypeItemById)
					dataGroup.DELETE("/:typeName/:itemId", typesController.DeleteTypeItemById)
				}
			}
		}
	}
}

func Initialize(e *gin.Engine) {
	if IsInitialized() {
		return
	}

	initialized = true

	initMiddleware(e)
	initRoutes(e)
}

func IsInitialized() bool {
	return initialized
}

func RespondId(id interface{}, c *gin.Context) {
	i := models.JSON{}

	switch t := id.(type) {
	case models.JSON:
		i["id"] = t["id"]
	default:
		i["id"] = t
	}

	c.JSON(http.StatusOK, i)
}

func GetHeaderOptions(c *gin.Context) models.Options {
	v, exists := c.Get(CONTEXT_HEADER_OPTIONS)
	if !exists {
		return models.Options{}
	}

	return v.(models.Options)
}

func GetExpression(c *gin.Context) (expression.ExpressionGroup, models.JSON, error) {
	v, exists := c.Get(CONTEXT_EXPRESSION)
	if !exists {
		return expression.ExpressionGroup{}, models.JSON{}, nil
	}

	exp := v.(expression.ExpressionGroup)
	query, err := exp.ToMongo()
	return exp, query, err
}
