package integrationtests

import (
	"net/http"
	"encoding/json"
	"strings"
	"math/rand"
	"time"
	"strconv"
	"github.com/go-neutrino/neutrino-core/models"
	"github.com/go-neutrino/neutrino-core/config"
)

var (
	ApiBaseUrl = "http://localhost" + config.Get(config.KEY_API_PORT) + "/v1/"
	AppId = ""
	Token = ""
	Email = ""
	Password = ""
)

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	return "r" + strconv.Itoa(rand.Int())
}

func init() {
	//initialize the tests
	Email = randomString() + "@gmail.com"
	Password = randomString()

	Register(Email, Password)
	Token = Login(Email, Password)
	AppId = CreateApp(randomString())
}

func SendRequest(baseUrl, url, method string, body interface{}) models.JSON {
	var bodyStr = ""
	if body != nil {
		b, err := json.Marshal(body)
		if err !=nil {
			panic (err)
		}

		bodyStr = b
	}

	req, err := http.NewRequest(method, baseUrl + url, strings.NewReader(bodyStr))
	if err != nil {
		panic(err)
	}

	if Token != "" {
		req.Header.Set("Authorization", "Bearer "+Token)
	}

	res, err := http.Client.Do(req)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	json := models.JSON{}
	err = json.FromResponse(res)
	if err != nil {
		panic(err)
	}

	return json
}

func CreateApp(name string) string {
	res := SendRequest(ApiBaseUrl, "app", "post", models.JSON{
		"name": name,
	})

	return res["_id"]
}

func Register(email, password string) {
	SendRequest(ApiBaseUrl, "register", "post", models.JSON{
		"email": email,
		"password": password,
	})
}

func Login(email, password string) string {
	res := SendRequest(ApiBaseUrl, "login", "post", models.JSON{
		"email": email,
		"password": password,
	})

	return res["token"]
}