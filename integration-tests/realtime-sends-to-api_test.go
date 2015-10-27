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

func sleep() {
	log.Info("Sleeping for a second...")
	time.Sleep(time.Second * 1)
}

func TestCreateItemFromClient(t *testing.T) {
	RealtimeData.Create(models.JSON{
		"test": "integration",
	})

	sleep()

	items, err := ApiClient.GetItems(DataType)
	if err != nil {
		t.Error(err)
		return
	}

	if len(items) != 1 || items[0]["test"] != "integration" {
		t.Error("Invalid items created")
	}
}

func TestUpdateItemFromClient(t *testing.T) {
	RealtimeData.Create(models.JSON{
		"test":"",
	})
	sleep()

	items, err := ApiClient.GetItems(DataType)
	if err != nil {
		t.Error(err)
		return
	}

	item := items[0]
	RealtimeData.Update(item["_id"].(string), models.JSON{
		"test": "updated",
	})
	sleep()

	items, err = ApiClient.GetItems(DataType)
	if err != nil {
		t.Error(err)
		return
	}

	item = items[0]
	if item["test"] != "updated" {
		t.Error("Item not updated properly")
	}
}
