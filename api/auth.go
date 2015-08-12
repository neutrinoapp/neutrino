package api

import (
	"github.com/labstack/echo"
	"realbase/core"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserHandler(c *echo.Context) error {
	b := JsonBody(c.Request())
	username, email := b["username"], b["email"]
	val, _ := b["password"].(string)

	password := []byte(val)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		panic(err)
	}

	db := realbase.NewUsersDbService()
	doc := map[string]interface{}{
		"username": username,
		"email": email,
		"password": hashedPassword,
	}

	return db.Insert(doc)

	//err = bcrypt.CompareHashAndPassword(hashedPassword, password)
}