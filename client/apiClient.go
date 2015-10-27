package client

import (
	"encoding/json"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiClient struct {
	BaseUrl, Token, AppId string
}

func NewApiClient(url, appId string) *ApiClient {
	return &ApiClient{
		BaseUrl: url,
		Token:   "",
		AppId:   appId,
	}
}

func (c *ApiClient) SendRequest(url, method string, body interface{}, isArray bool) (interface{}, error) {
	log.Info(
		"Sending request",
		"BaseUrl:", c.BaseUrl,
		"Url:", url,
		"Method:", method,
		"Body:", body,
		"Token:", c.Token,
		"AppId:", c.AppId,
	)
	var bodyStr = ""
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		bodyStr = string(b)
	}

	req, err := http.NewRequest(method, c.BaseUrl+url, strings.NewReader(bodyStr))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	bodyRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if string(bodyRes) == "" {
		return nil, nil
	}

	var result interface{}
	if isArray {
		jsonArray := make([]models.JSON, 0)
		err = json.Unmarshal(bodyRes, &jsonArray)
		result = jsonArray
	} else {
		m := models.JSON{}
		err = json.Unmarshal(bodyRes, &m)
		result = m
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return result, nil
}

func (c *ApiClient) CreateApp(name string) (string, error) {
	res, err := c.SendRequest("app", "POST", models.JSON{
		"name": name,
	}, false)

	if res == nil {
		return "", err
	}

	return res.(models.JSON)["_id"].(string), nil
}

func (c *ApiClient) Register(email, password string) error {
	_, err := c.SendRequest("register", "POST", models.JSON{
		"email":    email,
		"password": password,
	}, false)

	return err
}

func (c *ApiClient) Login(email, password string) (string, error) {
	res, err := c.SendRequest("login", "POST", models.JSON{
		"email":    email,
		"password": password,
	}, false)

	if res == nil {
		return "", err
	}

	c.Token = res.(models.JSON)["token"].(string)

	return c.Token, nil
}

func (c *ApiClient) CreateItem(t string, m models.JSON) (models.JSON, error) {
	res, err := c.SendRequest("app/"+c.AppId+"/data/"+t, "POST", m, false)
	if res == nil {
		return nil, err
	}

	return res.(models.JSON), err
}

func (c *ApiClient) UpdateItem(t, id string, m models.JSON) (models.JSON, error) {
	res, err := c.SendRequest("app/"+c.AppId+"/data/"+t+"/"+id, "PUT", m, false)
	if res == nil {
		return nil, err
	}

	return res.(models.JSON), err
}

func (c *ApiClient) DeleteItem(t, id string) (models.JSON, error) {
	res, err := c.SendRequest("app/"+c.AppId+"/data/"+t+"/"+id, "DELETE", nil, false)
	if res == nil {
		return nil, err
	}

	return res.(models.JSON), err
}

func (c *ApiClient) GetItem(t, id string) (models.JSON, error) {
	res, err := c.SendRequest("app/"+c.AppId+"/data/"+t+"/"+id, "GET", nil, false)
	if res == nil {
		return nil, err
	}

	return res.(models.JSON), err
}

func (c *ApiClient) GetItems(t string) ([]models.JSON, error) {
	res, err := c.SendRequest("app/"+c.AppId+"/data/"+t, "GET", nil, true)
	if res == nil {
		return nil, err
	}

	return res.([]models.JSON), err
}
