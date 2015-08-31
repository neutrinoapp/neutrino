package api

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"time"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/go-realbase/realbase/core"
)

type UserModel struct {
	Id string `json: "_id,email"`
	Email string `json: "_id,email"`
	Password string `json: "password"`
}

type AuthController struct {
}

func (a *AuthController) Path() string {
	return "/auth"
}

func registerUser(w rest.ResponseWriter, r *rest.Request, db realbase.DbService) {
	var u bson.M

	if err := r.DecodeJsonPayload(&u); err != nil {
		RestErrorInvalidBody(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u["password"].(string)), 10)
	if err != nil {
		RestError(w, err)
		return
	}

	doc := bson.M{
		"_id": u["email"].(string),
		"password": hashedPassword,
		"createdAt": time.Now(),
	}

	if err := db.Insert(doc); err != nil {
		RestError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func loginUser(w rest.ResponseWriter, r *rest.Request, db realbase.DbService) {
	u := UserModel{}

	if err := r.DecodeJsonPayload(&u); err != nil {
		RestError(w, err)
		return
	}

	existingUser, err := db.FindId(u.Id, nil)

	if err != nil {
		RestError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(existingUser["password"].([]byte), []byte(u.Password))

	if err != nil {
		RestError(w, err)
		return
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = u.Id
	token.Claims["expiration"] = time.Now().Add(time.Minute + 60).Unix()

	tokenStr, err := token.SignedString([]byte(""))

	if err != nil {
		RestError(w, err)
		return
	}

	w.WriteJson(map[string]string{"token": tokenStr})
}

func (a *AuthController) RegisterUserHandler (w rest.ResponseWriter, r *rest.Request) {
	registerUser(w, r, realbase.NewUsersDbService())
}

func (a *AuthController) AppRegisterUserHandler (w rest.ResponseWriter, r *rest.Request) {
	registerUser(w, r, realbase.NewAppUsersDbService(r.PathParam("appId")))
}

func (a *AuthController) LoginUserHandler(w rest.ResponseWriter, r *rest.Request) {
	loginUser(w, r, realbase.NewUsersDbService())
}

func (a *AuthController) AppLoginUserHandler(w rest.ResponseWriter, r *rest.Request) {
	loginUser(w, r, realbase.NewAppUsersDbService(r.PathParam("appId")))
}
