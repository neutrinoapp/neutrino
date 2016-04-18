package server

import (
	"github.com/gngeorgiev/gowamp"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
	"github.com/neutrinoapp/neutrino/src/common/utils"
)

type RpcMessageReceiver struct {
	WsClient         *gowamp.Client
	WsReceiver       WsMessageReceiver
	MessageProcessor messaging.MessageProcessor
}

func NewRpcMessageReceiver(wsClient *gowamp.Client, wsProcessor WsMessageReceiver) RpcMessageReceiver {
	return RpcMessageReceiver{wsClient, wsProcessor, messaging.NewMessageProcessor()}
}

func (p RpcMessageReceiver) Receive() {
	p.WsClient.BasicRegister("data.read", p.handleRpc)
	p.WsClient.BasicRegister("data.create", p.handleRpc)
	p.WsClient.BasicRegister("data.remove", p.handleRpc)
	p.WsClient.BasicRegister("data.update", p.handleRpc)
}

func (p RpcMessageReceiver) makeResult(data interface{}) *gowamp.CallResult {
	return &gowamp.CallResult{Args: []interface{}{data}}
}

func (p RpcMessageReceiver) makeErrorResult(err error) *gowamp.CallResult {
	log.Error(err)
	return &gowamp.CallResult{
		Err:  gowamp.URI(err.Error()),
		Args: []interface{}{err},
	}
}

func (p RpcMessageReceiver) handleRpc(args []interface{}, kwargs map[string]interface{}) *gowamp.CallResult {
	defer utils.Recover()

	var msg messaging.Message

	if msgStr, ok := args[0].(string); ok {
		err := msg.FromString(msgStr)
		if err != nil {
			return p.makeErrorResult(err)
		}
	} else {
		err := models.Convert(args[0], &msg)
		if err != nil {
			return p.makeErrorResult(err)
		}
	}

	log.Info("RPC message received:", msg)
	res, err := p.MessageProcessor.Process(msg)
	if err != nil {
		return p.makeErrorResult(err)
	}

	return p.makeResult(res)
}
