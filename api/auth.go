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
	Username string `json: "username"`
	Password string `json: "password"`
	Email string `json: "email`
}

func RegisterUserHandler (w rest.ResponseWriter, r *rest.Request) {
	u := UserModel{}

	if err := r.DecodeJsonPayload(&u); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := realbase.NewUsersDbService()
	doc := bson.M{
		"_id": u.Username,
		"email": u.Email,
		"password": hashedPassword,
		"createdAt": time.Now(),
	}

	if err := db.Insert(doc); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoginUserHandler(w rest.ResponseWriter, r *rest.Request) {
	u := UserModel{}
	if err := r.DecodeJsonPayload(&u); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	db := realbase.NewUsersDbService()
	existingUser, err := db.FindId(u.Username)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword(existingUser["password"].([]byte), []byte(u.Password))

	if err != nil {
		return
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = existingUser["username"]
	token.Claims["expiration"] = time.Now().Add(time.Minute + 60).Unix()

	tokenStr, err := token.SignedString([]byte(""))

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteJson(map[string]string{"token": tokenStr})
}
