package server

import (
	"net/http"

	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"gopkg.in/jcelliott/turnpike.v2"
)

func NewWebSocketServer() (*turnpike.WebsocketServer, *http.Server, *turnpike.Client, *wsInterceptor, error) {
	interceptor := &wsInterceptor{
		m: make(chan turnpike.Message),
	}

	r := turnpike.Realm{}
	r.Interceptor = interceptor

	realms := map[string]turnpike.Realm{}
	realms[config.CONST_DEFAULT_REALM] = r
	wsServer, err := turnpike.NewWebsocketServer(realms)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	wsServer.Upgrader.CheckOrigin = func(r *http.Request) bool {
		//allow connections from any origin
		return true
	}

	c, err := wsServer.GetLocalClient(config.CONST_DEFAULT_REALM, nil)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil, err
	}

	server := &http.Server{
		Handler: wsServer,
		Addr:    config.Get(config.KEY_REALTIME_PORT),
	}

	return wsServer, server, c, interceptor, nil
}
