package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"strings"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"realbase/core"
	"gopkg.in/mgo.v2/bson"
)

var initialized bool

type authMiddleware struct {}
type environmentMiddleware struct {}

func (a *authMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			rest.Error(w, "Not authorized", 401)
			return
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			rest.Error(w, "Not authorized", 401)
			return
		}

		token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
			if(jwt.GetSigningMethod("HS256") != token.Method){
				rest.Error(w, "Invalid signing token algorithm", 500)
				return nil, nil
			}

			return []byte(""), nil
		})

		r.Env["token"] = token
		r.Env["user"] = token.Claims["user"]

		if err != nil {
			rest.Error(w, err.Error(), 500)
			return
		}

		handler(w, r)
	}
}

func (e *environmentMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		appId := r.PathParam("appId")

		if appId != "" {
			//TODO: cache this
			appDb := realbase.NewApplicationsDbService()
			id := bson.ObjectId(appId)
			app, err := appDb.FindId(id, nil)

			if err != nil {
				RestGeneralError(w, err)
				handler(w, r)
				return
			}

			r.Env["app"] = ApplicationModel{
				Id: app["_id"].(bson.ObjectId),
				Name: app["Name"].(string),
			}
		}

		handler(w, r)
	}
}

func initMiddleware(restApi *rest.Api) {
	restApi.Use(
		&rest.AccessLogJsonMiddleware{},
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{"Realbase"},
		&rest.RecoverMiddleware{EnableResponseStackTrace: true},
		&rest.JsonIndentMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
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
		rest.Get("/applications/#appId", GetApplicationHandler),
		//TODO: moar apps

		rest.Post("/#appId/types", CreateTypeHandler),
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