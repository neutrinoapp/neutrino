package common

import (
	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/utils"
)

var (
	c *client.ApiClient
)

func GetClient() *client.ApiClient {
	return c
}

func init() {
	c = client.NewApiClientClean()

	email := utils.GetCleanUUID()
	password := utils.GetCleanUUID()

	err := c.Register(email, password)
	if err != nil {
		panic(err)
		return
	}

	token, err := c.Login(email, password)
	if err != nil {
		panic(err)
		return
	}

	c.Token = token

	appId, err := c.CreateApp(utils.GetCleanUUID())
	if err != nil {
		panic(err)
		return
	}

	c.AppId = appId
	err = c.AppRegister(email, password)
	if err != nil {
		panic(err)
		return
	}

	appToken, err := c.AppLogin(email, password)
	if err != nil {
		panic(err)
		return
	}

	c.Token = appToken
}
