package neutrinoclient

import (
	"github.com/go-neutrino/neutrino/client"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"sync"
)

type (
	NeutrinoClient struct {
		RealtimeAddr    string
		ApiAddr         string
		AppId           string
		Token           string
		WebsocketClient *client.WebsocketClient
		ApiClient       *client.ApiClient
		listeners       []*NeutrinoData
	}
)

func NewClient(appId string) *NeutrinoClient {
	wsAddr := "ws://localhost" + config.Get(config.KEY_REALTIME_PORT) + "/data?app=" + appId
	apiAddr := "http://localhost" + config.Get(config.KEY_API_PORT) + "/v1/"

	c := &NeutrinoClient{
		AppId:           appId,
		RealtimeAddr:    wsAddr,
		ApiAddr:         apiAddr,
		WebsocketClient: client.NewWebsocketClient(wsAddr),
		ApiClient:       client.NewApiClient(apiAddr, appId),
	}

	go func() {
		processMessage := func(msg string) {
			log.Info("Neutrino client sending message to data callback:", msg)
			var m models.JSON
			err := m.FromString([]byte(msg))
			if err != nil {
				log.Error(err)
				return
			}

			typeName := m["type"]
			for _, listener := range c.listeners {
				if listener.DataName == typeName {
					listener.onDataMessage(m)
				}
			}
		}

		for {
			select {
			case msg := <-c.WebsocketClient.Message:
				log.Info("Neutrino client got message:", msg)
				if msg != "" {
					processMessage(msg)
				}
			}
		}
	}()

	return c
}

func (c *NeutrinoClient) Register(u, p string) error {
	return c.ApiClient.Register(u, p)
}

func (c *NeutrinoClient) Login(u, p string) (string, error) {
	token, err := c.ApiClient.Login(u, p)
	c.Token = token

	return token, err
}

func (c *NeutrinoClient) registerDataListener(d *NeutrinoData) {
	m := &sync.Mutex{}
	m.Lock()
	c.listeners = append(c.listeners, d)
	m.Unlock()
}
