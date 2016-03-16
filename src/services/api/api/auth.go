package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/dgrijalva/jwt-go.v2"
)

type UserModel struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthController struct {
	*BaseController
}

func NewAuthController() *AuthController {
	return &AuthController{NewBaseController()}
}

func (a *AuthController) registerUser(c *gin.Context, isApp bool) {
	var u models.JSON

	if err := c.Bind(&u); err != nil {
		log.Error(RestErrorInvalidBody(c))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u[db.PASSWORD_FIELD].(string)), 10)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	user := models.JSON{
		db.EMAIL_FIELD:         u[db.EMAIL_FIELD],
		db.PASSWORD_FIELD:      hashedPassword,
		db.REGISTERED_AT_FIELD: time.Now(),
	}

	if isApp {
		user[db.APP_ID_FIELD] = c.Param(db.APP_ID_FIELD)
	}

	err = a.DbService.CreateUser(user, isApp)
	if err != nil {
		log.Error(RestError(c, err))
	}
}

func (a *AuthController) loginUser(c *gin.Context, isApp bool) {
	var u UserModel

	if err := c.Bind(&u); err != nil {
		log.Error(RestError(c, err))
		return
	}

	appId := ""
	if isApp {
		appId = c.Param(db.APP_ID_FIELD)
	}

	user, err := a.DbService.GetUser(u.Email, isApp, appId)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	err = bcrypt.CompareHashAndPassword(user[db.PASSWORD_FIELD].([]byte), []byte(u.Password))
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = u.Email
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
	a.registerUser(c, false)
}

func (a *AuthController) AppRegisterUserHandler(c *gin.Context) {
	a.registerUser(c, true)
}

func (a *AuthController) LoginUserHandler(c *gin.Context) {
	a.loginUser(c, false)
}

func (a *AuthController) AppLoginUserHandler(c *gin.Context) {
	a.loginUser(c, true)
}
