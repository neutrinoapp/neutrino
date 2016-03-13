package server

import "github.com/neutrinoapp/neutrino/src/common/client"

func Initialize() error {
	redisClient := client.GetNewRedisClient()
	clientMessageProcessor := NewClientMessageProcessor()

	_, wsClient, interceptor, err := NewWebSocketServer()
	if err != nil {
		return err
	}

	natsProcessor := JobsMessageProcessor{redisClient, wsClient}
	natsProcessor.Process()

	wsProcessor := NewWsMessageProcessor(interceptor, redisClient, clientMessageProcessor, wsClient)
	wsProcessor.Process()

	rpcProcessor := RpcMessageProcessor{wsClient, wsProcessor}
	rpcProcessor.Process()

	return nil
}
