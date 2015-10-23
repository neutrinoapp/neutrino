package integrationtests

import (
	"net/http"
	"encoding/json"
	"strings"
	"math/rand"
	"time"
	"strconv"
	"github.com/go-neutrino/neutrino/models"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/realtime-service/client"
)

var (
	ApiBaseUrl = "http://localhost" + config.Get(config.KEY_API_PORT) + "/v1/"
	AppId = ""
	Token = ""
	Email = ""
	Password = ""
	Client *neutrinoclient.NeutrinoClient
	Data *neutrinoclient.NeutrinoData
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

	Client = neutrinoclient.NewClient(AppId)
	Data = Client.Data("test")
}

func SendRequest(baseUrl, url, method string, body interface{}) models.JSON {
	var bodyStr = ""
	if body != nil {
		b, err := json.Marshal(body)
		if err !=nil {
			panic (err)
		}

		bodyStr = string(b)
	}

	req, err := http.NewRequest(method, baseUrl + url, strings.NewReader(bodyStr))
	if err != nil {
		panic(err)
	}

	if Token != "" {
		req.Header.Set("Authorization", "Bearer "+Token)
	}

	client := http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	json := models.JSON{}
	err = json.FromResponse(res)
	if err != nil {
		log.Error(err)
	}

	return json
}

func CreateApp(name string) string {
	res := SendRequest(ApiBaseUrl, "app", "POST", models.JSON{
		"name": name,
	})

	return res["_id"].(string)
}

func Register(email, password string) {
	SendRequest(ApiBaseUrl, "register", "POST", models.JSON{
		"email": email,
		"password": password,
	})
}

func Login(email, password string) string {
	res := SendRequest(ApiBaseUrl, "login", "POST", models.JSON{
		"email": email,
		"password": password,
	})

	return res["token"].(string)
}

func CreateItem(t string, m models.JSON) models.JSON {
	return SendRequest(ApiBaseUrl, "app/" + AppId + "/data/" + t, "POST", m)
}

func UpdateItem(t, id string, m models.JSON) models.JSON {
	return SendRequest(ApiBaseUrl, "app/" + AppId + "/data/" + t + "/" + id, "PUT", m)
}

func DeleteItem(t, id string) models.JSON {
	return SendRequest(ApiBaseUrl, "app/" + AppId + "/data/" + t + "/" + id, "DELETE", nil)
}