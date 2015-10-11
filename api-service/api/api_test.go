package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-neutrino/neutrino/api-service/db"
	"github.com/go-neutrino/neutrino/models"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	"runtime"
)

var (
	apiHandler http.Handler
	user       map[string]interface{}
	token      string
)

type ResRecorder struct {
	*httptest.ResponseRecorder
	t *testing.T
}

func (r *ResRecorder) CodeIs(s int) {
	if r.Code != s {
		_, file, line, _ := runtime.Caller(2)
		r.t.Error(r.Code, "is different from", s, file+":"+strconv.Itoa(line))
	}
}

func (r *ResRecorder) BodyString() string {
	return r.Body.String()
}

func (r *ResRecorder) BodyHas(str string) {
	if !strings.Contains(r.BodyString(), str) {
		r.t.Error(r.BodyString(), "does not contain", str)
	}
}

func (r *ResRecorder) BodyJSON() models.JSON {
	b, _ := ioutil.ReadAll(r.Body)
	var res models.JSON
	json.Unmarshal(b, &res)
	return res
}

func (r *ResRecorder) Decode(o interface{}) {
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, o)
}

func sendAuthenticatedRequest(method, path string, body interface{}, t *testing.T) *ResRecorder {
	login(t)
	return sendRequest(method, path, body, t)
}

func sendRequest(method, path string, body interface{}, t *testing.T) *ResRecorder {
	if !IsInitialized() {
		e := gin.Default()
		apiHandler = e
		Initialize(e)
		db.Initialize()
		httptest.NewServer(e)

		e.Use(func() gin.HandlerFunc {
			return func(c *gin.Context) {
				fmt.Println("###")
				fmt.Println("URL: -> " + c.Request.URL.String())
				fmt.Println("Method: -> " + c.Request.Method)

				fmt.Println("Headers: ->")
				for k := range c.Request.Header {
					fmt.Println(c.Request.Header[k])
				}
				fmt.Println("<-")

				if user, exists := c.Get("user"); exists {
					fmt.Println("User: -> ", user)
				}

				b, _ := ioutil.ReadAll(c.Request.Body)
				fmt.Println("Body: -> ", string(b))

				fmt.Println("### -->>")

				c.Next()
			}
		}())
	}

	var b string
	if body != nil {
		bodyStr, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		b = fmt.Sprintf("%s", bodyStr)
	}

	req, err := http.NewRequest(method, "/v1"+path, strings.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	apiHandler.ServeHTTP(w, req)
	return &ResRecorder{w, t}
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())

	return "r" + strconv.Itoa(rand.Int())
}

func register(t *testing.T) map[string]interface{} {
	b := models.JSON{
		"email":    randomString() + "@gmail.com",
		"password": "pass",
	}

	rec := sendRequest("POST", "/register", b, t)

	rec.CodeIs(http.StatusOK)

	return b
}

func login(t *testing.T) (*UserModel, string) {
	if token == "" {
		user = register(t)
		rec := sendRequest("POST", "/login", models.JSON{
			"email":    user["email"],
			"password": user["password"],
		}, t)

		response := rec.BodyJSON()
		token = response["token"].(string)
	}

	return &UserModel{
		Email:    user["email"].(string),
		Password: user["password"].(string),
	}, token
}
