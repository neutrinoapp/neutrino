package client

import (
	"encoding/json"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
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

func (c *ApiClient) SendRequest(url, method string, body interface{}) (models.JSON, error) {
	log.Info("Sending request", url, method, body)
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

	var json models.JSON
	err = json.FromResponse(res)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return json, nil
}

func (c *ApiClient) CreateApp(name string) (string, error) {
	res, err := c.SendRequest("app", "POST", models.JSON{
		"name": name,
	})

	if res == nil {
		return "", err
	}

	return res["_id"].(string), nil
}

func (c *ApiClient) Register(email, password string) error {
	_, err := c.SendRequest("register", "POST", models.JSON{
		"email":    email,
		"password": password,
	})

	return err
}

func (c *ApiClient) Login(email, password string) (string, error) {
	res, err := c.SendRequest("login", "POST", models.JSON{
		"email":    email,
		"password": password,
	})

	if res == nil {
		return "", err
	}

	c.Token = res["token"].(string)

	return c.Token, nil
}

func (c *ApiClient) CreateItem(t string, m models.JSON) (models.JSON, error) {
	return c.SendRequest("app/"+c.AppId+"/data/"+t, "POST", m)
}

func (c *ApiClient) UpdateItem(t, id string, m models.JSON) (models.JSON, error) {
	return c.SendRequest("app/"+c.AppId+"/data/"+t+"/"+id, "PUT", m)
}

func (c *ApiClient) DeleteItem(t, id string) (models.JSON, error) {
	return c.SendRequest("app/"+c.AppId+"/data/"+t+"/"+id, "DELETE", nil)
}
