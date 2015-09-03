package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
)

var initialized bool

func initMiddleware(restApi *rest.Api) {
	restApi.Use(
		&defaultContentTypeMiddleware{"application/json"},
		&rest.AccessLogJsonMiddleware{},
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{"neutrino"},
		&rest.RecoverMiddleware{EnableResponseStackTrace: false},
		&rest.JsonIndentMiddleware{},
		&rest.IfMiddleware{
			Condition: func(request *rest.Request) bool {
				return request.URL.Path != "/auth"
			},
			IfTrue: &authMiddleware{},
		},
		&rest.IfMiddleware{
			Condition: func(request *rest.Request) bool {
				return request.URL.Path != "/auth" && request.URL.Path != "/applications"
			},
			IfTrue: &environmentMiddleware{},
		},
	)
}

func initRoutes(restApi *rest.Api) {
	authController := &AuthController{}
	appController := &ApplicationController{}
	typesController := &TypesController{}

	router, err := rest.MakeRouter(
		rest.Put(authController.Path(), authController.RegisterUserHandler),
		rest.Post(authController.Path(), authController.LoginUserHandler),
		rest.Put("/:appId" + authController.Path(), authController.AppRegisterUserHandler),
		rest.Post("/:appId" + authController.Path(), authController.AppLoginUserHandler),

		rest.Post(appController.Path(), appController.CreateApplicationHandler),
		rest.Get(appController.Path(), appController.GetApplicationsHandler),
		rest.Get(appController.Path() + "/:appId", appController.GetApplicationHandler),
		rest.Delete(appController.Path() + "/:appId", appController.DeleteApplicationHandler),
		rest.Put(appController.Path() + "/:appId", appController.UpdateApplicationHandler),

		rest.Post("/:appId" + typesController.Path(), typesController.CreateTypeHandler),
		rest.Delete("/:appId" + typesController.Path() + "/:typeName", typesController.DeleteType),
		rest.Post("/:appId" + typesController.Path() + "/:typeName", typesController.InsertInTypeHandler),
		rest.Get("/:appId" + typesController.Path() + "/:typeName", typesController.GetTypeDataHandler),
		rest.Get("/:appId" + typesController.Path() + "/:typeName/:itemId", typesController.GetTypeItemById),
		rest.Put("/:appId" + typesController.Path() + "/:typeName/:itemId", typesController.UpdateTypeItemById),
		rest.Delete("/:appId" + typesController.Path() + "/:typeName/:itemId", typesController.DeleteTypeItemById),
	)

	if err != nil {
		log.Fatal(err)
	}

	restApi.SetApp(router)
}

func Initialize(restApi *rest.Api) {
	if IsInitialized() {
		return
	}

	initialized = true
	initMiddleware(restApi)
	initRoutes(restApi)
}

func IsInitialized() bool {
	return initialized
}

func RespondId(id interface{}, w rest.ResponseWriter) {
	w.WriteJson(map[string]interface{}{
		"_id": id,
	})
}