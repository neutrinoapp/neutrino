package server

import "github.com/neutrinoapp/neutrino/src/common/client"

func Initialize() error {
	redisClient := client.GetNewRedisClient()

	_, wsClient, interceptor, err := NewWebSocketServer()
	if err != nil {
		return err
	}

	wsReceiver := NewWsMessageReceiver(interceptor, redisClient, wsClient)
	wsReceiver.Receive()

	rpcReceiver := NewRpcMessageReceiver(wsClient, wsReceiver)
	rpcReceiver.Receive()

	return nil
}
