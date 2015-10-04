package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino-core/api-service/db"
	"github.com/go-neutrino/neutrino-core/api-service/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"net/http"
	"time"
)

type UserModel struct {
	Id       string `json: "_id"`
	Email    string `json: "email"`
	Password string `json: "password"`
}

type AuthController struct {
}

func registerUser(c *gin.Context, d db.DbService) {
	var u JSON

	if err := c.Bind(&u); err != nil {
		RestErrorInvalidBody(c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u["password"].(string)), 10)
	if err != nil {
		RestError(c, err)
		return
	}

	doc := JSON{
		"_id":       u["email"].(string),
		"password":  hashedPassword,
		"createdAt": time.Now(),
	}

	if err := d.Insert(doc); err != nil {
		RestError(c, err)
		return
	}

	utils.OK(c)
}

func loginUser(c *gin.Context, d db.DbService) {
	var u JSON

	if err := c.Bind(&u); err != nil {
		RestError(c, err)
		return
	}

	existingUser, err := d.FindId(u["email"].(string), nil)

	if err != nil {
		RestError(c, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(existingUser["password"].([]byte), []byte(u["password"].(string)))

	if err != nil {
		RestError(c, err)
		return
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = u["email"].(string)
	token.Claims["expiration"] = time.Now().Add(time.Minute + 60).Unix()

	tokenStr, err := token.SignedString([]byte(""))

	if err != nil {
		RestError(c, err)
		return
	}

	c.JSON(http.StatusOK, JSON{
		"token": tokenStr,
	})
}

func (a *AuthController) RegisterUserHandler(c *gin.Context) {
	registerUser(c, db.NewUsersDbService())
}

func (a *AuthController) AppRegisterUserHandler(c *gin.Context) {
	registerUser(c, db.NewAppUsersDbService(c.Param("appId")))
}

func (a *AuthController) LoginUserHandler(c *gin.Context) {
	loginUser(c, db.NewUsersDbService())
}

func (a *AuthController) AppLoginUserHandler(c *gin.Context) {
	loginUser(c, db.NewAppUsersDbService(c.Param("appId")))
}
