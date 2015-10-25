package integrationtests

import (
	"github.com/go-neutrino/neutrino/log"
	"github.com/go-neutrino/neutrino/models"
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func init() {
	var conn *websocket.Conn
	for {
		if conn == nil {
			log.Info("Connection to websocket service not yet established.")
			conn = RealtimeClient.WebsocketClient.GetConnection()
		} else {
			break
		}

		time.Sleep(time.Second * 5)
	}
}

func TestCreateItemFromClient(t *testing.T) {
	RealtimeData.Create(models.JSON{
		"test": "integration",
	})

}
