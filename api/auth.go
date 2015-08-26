package api

import (
	"realbase/core"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"time"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"gopkg.in/mgo.v2/bson"
)

type UserModel struct {
	Password string `json: "password"`
	Email string `json: "email`
}

func RegisterUserHandler (w rest.ResponseWriter, r *rest.Request) {
	u := UserModel{}

	if err := r.DecodeJsonPayload(&u); err != nil {
		RestErrorInvalidBody(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		RestGeneralError(w, err)
		return
	}

	db := realbase.NewUsersDbService()
	doc := bson.M{
		"_id": u.Email,
		"password": hashedPassword,
	}

	if err := db.Insert(doc); err != nil {
		RestGeneralError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoginUserHandler(w rest.ResponseWriter, r *rest.Request) {
	u := UserModel{}
	if err := r.DecodeJsonPayload(&u); err != nil {
		RestGeneralError(w, err)
		return
	}

	db := realbase.NewUsersDbService()
	existingUser, err := db.FindId(u.Email, nil)

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(existingUser["password"].([]byte), []byte(u.Password))

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = u.Email
	token.Claims["expiration"] = time.Now().Add(time.Minute + 60).Unix()

	tokenStr, err := token.SignedString([]byte(""))

	if err != nil {
		RestGeneralError(w, err)
		return
	}

	w.WriteJson(map[string]string{"token": tokenStr})
}
