package integrationtests

import (
	"testing"
	"github.com/go-neutrino/neutrino/realtime-service/client"
	"time"
	"github.com/go-neutrino/neutrino/models"
	"sync"
)

var c = neutrinoclient.NewClient("15f7fee2b1dc41a1b9b3398321358277")

func readMessage(t *testing.T, expectedEv string, times int) sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(times)

	go func () {
		for i := 0; i < times; i++ {
			select {
			case ev := <- c.Data("test").Event:
				t.Log(ev.Code, ev.Data)

				if (ev.Code != expectedEv) {
					t.Error("Unexptected event.", ev.Data)
				}
			}

			wg.Done()
		}
	}()

	return wg
}

func TestInsertIntoType(t *testing.T) {
	wg := readMessage(t, neutrinoclient.EVENT_CREATE, 1)

	time.Sleep(time.Second * 1)

	CreateItem("test", models.JSON{
		"name": "test",
	})

	wg.Wait()
}
