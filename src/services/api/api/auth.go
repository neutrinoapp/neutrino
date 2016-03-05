package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils/webUtils"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/dgrijalva/jwt-go.v2"
)

type UserModel struct {
	Id       string `json: "_id"`
	Email    string `json: "email"`
	Password string `json: "password"`
}

type AuthController struct {
}

func registerUser(c *gin.Context, d db.DbService) {
	var u models.JSON

	if err := c.Bind(&u); err != nil {
		log.Error(RestErrorInvalidBody(c))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u["password"].(string)), 10)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	doc := models.JSON{
		"_id":       u["email"].(string),
		"password":  hashedPassword,
		"createdAt": time.Now(),
	}

	if err := d.Insert(doc); err != nil {
		//TODO: user exists
		log.Error(RestError(c, err))
		return
	}

	webUtils.OK(c)
}

func loginUser(c *gin.Context, d db.DbService, isApp bool) {
	var u models.JSON

	if err := c.Bind(&u); err != nil {
		log.Error(RestError(c, err))
		return
	}

	existingUser, err := d.FindId(u["email"].(string))

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	err = bcrypt.CompareHashAndPassword(existingUser["password"].([]byte), []byte(u["password"].(string)))

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = u["email"].(string)
	token.Claims["expiration"] = time.Now().Add(time.Minute + 60).Unix()
	token.Claims["inApp"] = isApp

	tokenStr, err := token.SignedString([]byte(""))

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	c.JSON(http.StatusOK, models.JSON{
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
	loginUser(c, db.NewUsersDbService(), false)
}

func (a *AuthController) AppLoginUserHandler(c *gin.Context) {
	loginUser(c, db.NewAppUsersDbService(c.Param("appId")), true)
}
