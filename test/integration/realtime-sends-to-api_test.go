package integrationtests

import (
	"github.com/go-neutrino/neutrino/src/common/log"
	"github.com/go-neutrino/neutrino/src/common/models"
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

		sleep()
	}
}

func sleep() {
	log.Info("Sleeping for a second...")
	time.Sleep(time.Second * 1)
}

func TestCreateItemFromClient(t *testing.T) {
	testType := randomString()
	d := RealtimeClient.Data(testType)

	d.Create(models.JSON{
		"test": "integration",
	})

	sleep()

	items, err := ApiClient.GetItems(testType)
	if err != nil {
		t.Error(err)
		return
	}

	if len(items) != 1 || items[0]["test"] != "integration" {
		t.Error("Invalid items created")
	}
}

func TestUpdateItemFromClient(t *testing.T) {
	testType := randomString()
	d := RealtimeClient.Data(testType)

	d.Create(models.JSON{
		"test": "",
	})
	sleep()

	items, err := ApiClient.GetItems(testType)
	if err != nil {
		t.Error(err)
		return
	}

	item := items[0]
	d.Update(item["_id"].(string), models.JSON{
		"test": "updated",
	})
	sleep()

	items, err = ApiClient.GetItems(testType)
	if err != nil {
		t.Error(err)
		return
	}

	item = items[0]
	if item["test"] != "updated" {
		t.Error("Item not updated properly")
	}
}

func TestDeleteItemFromClient(t *testing.T) {
	dType := randomString()
	d := RealtimeClient.Data(dType)

	d.Create(models.JSON{
		"test": "",
	})
	sleep()

	items, err := ApiClient.GetItems(dType)
	if err != nil {
		t.Error(err)
		return
	}

	item := items[0]
	d.Delete(item["_id"].(string))
	sleep()

	items, err = ApiClient.GetItems(dType)
	if err != nil {
		t.Error(err)
		return
	}

	if len(items) > 0 {
		t.Error("Item not deleted properly")
	}
}
