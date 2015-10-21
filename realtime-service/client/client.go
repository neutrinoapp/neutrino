package neutrinoclient

import (
"github.com/go-neutrino/neutrino/client"
"github.com/go-neutrino/neutrino/log"
)

type (
	NeutrinoClient struct {
		Addr string
		AppId string
		WebsocketClient *client.WebsocketClient
	}
)


func NewClient(appId string) *NeutrinoClient {
	addr := "ws://localhost:6000/data"

	c := &NeutrinoClient{
		AppId: appId,
		Addr: addr,
		WebsocketClient: client.NewWebsocketClient(addr),
	}

	go func () {
		processMessage := func (msg string) {
			log.Info(msg)
		}

		for {
			select {
			case msg := <- c.WebsocketClient.Message:
				if (msg != "") {
					processMessage(msg)
				}
			}
		}
	}()

	return c
}

func (c *NeutrinoClient) getAppUrl() string {
	return c.Addr + "/app/" + c.AppId
}