package api

import (
	"github.com/labstack/echo"
	"realbase/core"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

func RegisterUserHandler(c *echo.Context) error {
	b, err := JsonBody(c.Request())

	if err != nil {
		return err
	}

	username, email := b["username"], b["email"]
	val, ok := b["password"].(string)

	if !ok {
		return errors.New("Invalid password type")
	}

	password := []byte(val)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		return err;
	}

	db := realbase.NewUsersDbService()
	doc := map[string]interface{}{
		"_id": username,
		"email": email,
		"password": hashedPassword,
	}

	return db.Insert(doc)

	//err = bcrypt.CompareHashAndPassword(hashedPassword, password)
}