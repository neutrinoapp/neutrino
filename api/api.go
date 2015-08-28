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
		&rest.PoweredByMiddleware{"realbase"},
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
	router, err := rest.MakeRouter(
		rest.Put("/auth", RegisterUserHandler),
		rest.Post("/auth", LoginUserHandler),

		rest.Post("/applications", CreateApplicationHandler),
		rest.Get("/applications", GetApplicationsHandler),
		rest.Get("/applications/:appId", GetApplicationHandler),
		//TODO: moar apps

		rest.Post("/:appId/types", CreateTypeHandler),
		rest.Post("/:appId/types/:typeName", InsertInTypeHandler),
		rest.Get("/:appId/types/:typeName", GetTypeDataHandler),
		rest.Get("/:appId/types/:typeName/:itemId", GetTypeItemById),
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