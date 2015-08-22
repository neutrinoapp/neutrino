package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"strings"
	"gopkg.in/dgrijalva/jwt-go.v2"
)

var initialized bool

type authMiddleware struct {
}

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
	}
}

func initMiddleware(restApi *rest.Api) {
	restApi.Use(rest.DefaultDevStack...)
	restApi.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/auth"
		},
		IfTrue: &authMiddleware{},
	})
}

func initRoutes(restApi *rest.Api) {
	router, err := rest.MakeRouter(
		rest.Put("/auth", RegisterUserHandler),
		rest.Post("/auth", LoginUserHandler),

		rest.Post("/application", CreateApplicationHandler),
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