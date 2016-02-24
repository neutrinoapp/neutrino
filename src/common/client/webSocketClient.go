package client

import (
	"github.com/neutrinoapp/neutrino/src/common/config"
	"gopkg.in/jcelliott/turnpike.v2"
)

type WebSocketClient struct {
	*Client
}

func NewWebsocketClient(realms []string) *WebSocketClient {
	connect := func() (interface{}, error) {
		c, err := turnpike.NewWebsocketClient(turnpike.JSON, config.Get(config.KEY_REALTIME_ADDR))
		if err != nil {
			return nil, err
		}

		for _, r := range realms {
			c.JoinRealm(r, nil)
		}

		return c, err
	}

	baseClient := NewClient(connect, config.Get(config.KEY_REALTIME_ADDR))

	return &WebSocketClient{baseClient}
}

func (w *WebSocketClient) GetConnection() *turnpike.Client {
	if w.connection != nil {
		return w.connection.(*turnpike.Client)
	}

	return nil
}
