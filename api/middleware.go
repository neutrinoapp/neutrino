package api
import (
"github.com/ant0ine/go-json-rest/rest"
"strings"
	"gopkg.in/dgrijalva/jwt-go.v2"
)

type authMiddleware struct {}
type environmentMiddleware struct {}
type defaultContentTypeMiddleware struct {
	contentType string
}

func (a *authMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			rest.Error(w, "Not authorized.", 401)
			return
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			rest.Error(w, "Not authorized.", 401)
			return
		}

		token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
			if(jwt.GetSigningMethod("HS256") != token.Method){
				rest.Error(w, "Invalid signing token algorithm.", 500)
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
		handler(w, r)
	}
}

func (e *defaultContentTypeMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		contentTypeHeader := r.Header.Get("Content-Type")
		if contentTypeHeader == "" {
			r.Header.Set("Content-Type", e.contentType)
		}

		handler(w, r)
	}
}