package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var initialized bool
var config map[string]interface{}

func initMiddleware(e *gin.Engine) {
	e.Use(defaultContentTypeMiddleware())
}

func initRoutes(e *gin.Engine) {
	authController := &AuthController{}
	appController := &ApplicationController{}
	typesController := &TypesController{}

	v1 := e.Group("/v1")
	{
		v1.POST("/login", authController.LoginUserHandler)
		v1.POST("/register", authController.RegisterUserHandler)

		appGroup := v1.Group("/app", authorizeMiddleware())
		{
			appGroup.POST("", appController.CreateApplicationHandler)
			appGroup.GET("", appController.GetApplicationsHandler)

			appIdGroup := appGroup.Group("/:appId", injectAppMiddleware())
			{
				appIdGroup.GET("", appController.GetApplicationHandler)
				appIdGroup.DELETE("", appController.DeleteApplicationHandler)
				appIdGroup.PUT("", appController.UpdateApplicationHandler)

				appIdGroup.POST("/register", authController.AppRegisterUserHandler)
				appIdGroup.POST("/login", authController.AppLoginUserHandler)

				dataGroup := appIdGroup.Group("/data")
				{
					dataGroup.POST("", typesController.CreateTypeHandler)
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

func Initialize(e *gin.Engine, c map[string]interface{}) {
	if IsInitialized() {
		return
	}

	initialized = true
	config = c

	initMiddleware(e)
	initRoutes(e)
}

func IsInitialized() bool {
	return initialized
}

func RespondId(id interface{}, c *gin.Context) {
	i := JSON{}

	switch t := id.(type) {
	case JSON:
		i["_id"] =  t["_id"]
	default:
		i["_id"] =  t
	}

	c.JSON(http.StatusOK, i)
}