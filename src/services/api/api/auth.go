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

	r "github.com/dancannon/gorethink"
)

type UserModel struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthController struct {
}

//TODO: refactor
func registerUser(c *gin.Context, t r.Term, s *r.Session, isApp bool) {
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

	user := models.JSON{
		"id":        u["email"].(string),
		"password":  hashedPassword,
		"createdAt": time.Now(),
	}

	var query r.Term
	if isApp {
		query = t.Update(func(row r.Term) interface{} {
			return models.JSON{
				db.USERS_FIELD: row.Field(db.USERS_FIELD).Append(user),
			}
		})
	} else {
		user["apps"] = make([]models.JSON, 0)
		query = t.Insert(user)
	}
	_, err = query.RunWrite(s)

	if err != nil {
		//TODO: user exists
		log.Error(RestError(c, err))
		return
	}

	c.Status(http.StatusOK)
}

func loginUser(c *gin.Context, t r.Term, s *r.Session, isApp bool) {
	var u UserModel

	if err := c.Bind(&u); err != nil {
		log.Error(RestError(c, err))
		return
	}

	var cu *r.Cursor
	var err error
	if isApp {
		cu, err = t.Filter(models.JSON{
			db.ID_FIELD: u.Email,
		}).Run(s)
	} else {
		cu, err = t.Get(u.Email).Run(s)
	}

	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	var existingUser models.JSON
	err = cu.One(&existingUser)
	if err != nil {
		log.Error(RestError(c, err))
		return
	}

	err = bcrypt.CompareHashAndPassword(existingUser["password"].([]byte), []byte(u.Password))
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
	d := db.NewDbService(db.DATABASE_NAME, db.USERS_TABLE)
	registerUser(c, d.Query(), d.GetSession(), false)
}

func (a *AuthController) AppRegisterUserHandler(c *gin.Context) {
	appId := c.Param("appId")
	d := db.NewDbService(db.DATABASE_NAME, db.DATA_TABLE)

	registerUser(c, d.Query().Get(appId), d.GetSession(), true)
}

func (a *AuthController) LoginUserHandler(c *gin.Context) {
	d := db.NewDbService(db.DATABASE_NAME, db.USERS_TABLE)

	loginUser(c, d.Query(), d.GetSession(), false)
}

func (a *AuthController) AppLoginUserHandler(c *gin.Context) {
	appId := c.Param("appId")
	d := db.NewDbService(db.DATABASE_NAME, db.DATA_TABLE)

	loginUser(c, d.Query().Get(appId).Field(db.USERS_FIELD), d.GetSession(), true)
}
