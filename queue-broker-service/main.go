package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-neutrino/neutrino-config"
	"github.com/nats-io/nats"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	services []string
	c        *viper.Viper
)

func refreshServices() {
	qHost := c.GetString(nconfig.KEY_QUEUE_STATS_HOST) + "/connz"

	res, err := http.Get(qHost)

	if err != nil {
		log.Println(err)
		return
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var o map[string]interface{}

	json.Unmarshal(b, &o)

	connz := o["connections"].([]interface{})

	services = make([]string, 0)
	realtimePort := c.GetString(nconfig.KEY_REALTIME_PORT)

	for _, e := range connz {
		con := e.(map[string]interface{})

		ip := con["ip"].(string)
		host := ip + realtimePort

		services = append(services, host)
	}
}

func jobsHandler(m *nats.Msg) {
	go func() {
		b := bytes.NewBuffer(m.Data)

		sCopy := services[:]

		for _, s := range sCopy {
			http.Post(s+"/message", "application/json", b)

			log.Println("Send message to " + s + ", " + string(m.Data))
		}
	}()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var wg sync.WaitGroup

	wg.Add(1)

	c = nconfig.Load()

	qHost := c.GetString(nconfig.KEY_QUEUE_HOST)

	n, e := nats.Connect(qHost)

	if e != nil {
		panic(e)
	}

	conn, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)

	if err != nil {
		panic(err)
	}

	refreshServices()

	conn.Subscribe("realtime-jobs", jobsHandler)

	wg.Wait()
}
