package api

import (
	"github.com/gin-gonic/gin"
	"testing"
	"net/http"
	"math/rand"
	"time"
	"strconv"
	"net/http/httptest"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/go-neutrino/neutrino-core/db"
	"github.com/go-neutrino/go-env-config"
)

var (
	apiHandler http.Handler
	user map[string]interface{}
	token string
)

type ResRecorder struct {
	*httptest.ResponseRecorder
	t *testing.T
}

func (r *ResRecorder) CodeIs(s int) {
	if r.Code != s {
		r.t.Error(r.Code, "is different from", s)
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

func (r *ResRecorder) BodyJSON() JSON {
	b, _ := ioutil.ReadAll(r.Body)
	var res JSON
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
		testConfig := envconfig.NewConfig()
		testConfig.M = map[string]interface{}{
			"mongoHost": "localhost:27017",
		}
		Initialize(e, testConfig)
		db.Initialize(testConfig)
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

	req, err := http.NewRequest(method, "/v1" + path, strings.NewReader(b))
	req.Header.Set("Authorization", "Bearer " + token)

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
	b := JSON{
		"email": randomString() + "@gmail.com",
		"password": "pass",
	}

	rec := sendRequest("POST", "/register", b, t)

	rec.CodeIs(http.StatusOK)

	return b
}

func login(t *testing.T) (*UserModel, string) {
	if token == "" {
		user = register(t)
		rec := sendRequest("POST", "/login", JSON{
			"email": user["email"],
			"password": user["password"],
		}, t)

		response := rec.BodyJSON()
		token = response["token"].(string)
	}

	return &UserModel{
		Email: user["email"].(string),
		Password: user["password"].(string),
	}, token
}