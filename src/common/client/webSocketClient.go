package client

import (
	"github.com/gngeorgiev/gowamp"
	"github.com/neutrinoapp/neutrino/src/common/config"
)

type WebSocketClient struct {
	*Client
}

func NewWebsocketClient(realms []string) *WebSocketClient {
	connect := func() (interface{}, error) {
		c, err := gowamp.NewWebsocketClient(gowamp.JSON, config.Get(config.KEY_REALTIME_ADDR))
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

func (w *WebSocketClient) GetConnection() *gowamp.Client {
	if w.connection != nil {
		return w.connection.(*gowamp.Client)
	}

	return nil
}
