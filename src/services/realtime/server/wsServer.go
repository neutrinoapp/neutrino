package server

import (
	"net/http"

	"github.com/gngeorgiev/gowamp"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
)

func NewWebSocketServer() (*gowamp.WebsocketServer, *gowamp.Client, *wsInterceptor, error) {
	interceptor := NewWsInterceptor()

	r := gowamp.Realm{}
	r.Interceptor = interceptor

	realms := map[string]gowamp.Realm{}
	realms[config.CONST_DEFAULT_REALM] = r
	wsServer, err := gowamp.NewWebsocketServer(realms)
	if err != nil {
		return nil, nil, nil, err
	}

	wsServer.Upgrader.CheckOrigin = func(r *http.Request) bool {
		//allow connections from any origin
		return true
	}

	c, err := wsServer.GetLocalClient(config.CONST_DEFAULT_REALM, nil)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, err
	}

	http.Handle("/", wsServer)
	http.HandleFunc("/_status", func(w http.ResponseWriter, r *http.Request) {
		//we are fine
	})

	return wsServer, c, interceptor, nil
}
