package neutrinoclient

import (
	"github.com/go-neutrino/neutrino/client"
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/messaging"
	"github.com/go-neutrino/neutrino/models"
	"sync"
)

type (
	NeutrinoClient struct {
		Addr            string
		AppId           string
		WebsocketClient *client.WebsocketClient
		listeners       []*NeutrinoData
	}
)

func NewClient(appId string) *NeutrinoClient {
	addr := "ws://localhost:6000/data?app=" + appId

	c := &NeutrinoClient{
		AppId:           appId,
		Addr:            addr,
		WebsocketClient: client.NewWebsocketClient(addr),
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
				if listener.Name == typeName {
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

func (c *NeutrinoClient) Login(u, p string) {
	payload := models.JSON{
		"email", u,
		"password": p,
	}

	messaging.GetMessageBuilder().Build(
		messaging.OP_LOGIN,
		messaging.ORIGIN_CLIENT,
		payload,
		nil,
		"users",
	).Send(c.WebsocketClient.GetConnection())
}

func (c *NeutrinoClient) Register(u, p string) {
	payload := models.JSON{
		"email", u,
		"password": p,
	}

	messaging.GetMessageBuilder().Build(
		messaging.OP_REGISTER,
		messaging.ORIGIN_CLIENT,
		payload,
		nil,
		"users",
	).Send(c.WebsocketClient.GetConnection())
}

func (c *NeutrinoClient) getAppUrl() string {
	return c.Addr + "/app/" + c.AppId
}

func (c *NeutrinoClient) registerDataListener(d *NeutrinoData) {
	m := &sync.Mutex{}
	m.Lock()
	c.listeners = append(c.listeners, d)
	m.Unlock()
}
