package main

import (
	"github.com/go-neutrino/neutrino/client"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/log"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
	"net/http"
	"strconv"
)

var (
	upgrader    = client.NewWebsocketUpgrader()
	connections []*websocket.Conn
)

func jobsHandler(m *nats.Msg) {
	log.Info("Got message " + string(m.Data))

	for i, c := range connections {
		log.Info("Sending message to connection:", strconv.Itoa(i+1), string(m.Data))
		c.WriteMessage(websocket.TextMessage, m.Data)
	}
}

func main() {
	c := client.NewNatsClient(config.Get(config.KEY_QUEUE_ADDR))
	//TODO: handle subscription after the connection to nats is lost and restored
	c.Subscribe(config.Get(config.CONST_REALTIME_JOBS_SUBJ), jobsHandler)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
			return
		}

		connections = append(connections, conn)
	})

	port := config.Get(config.KEY_BROKER_PORT)
	log.Info("Starting WS service on port " + port)
	http.ListenAndServe(port, nil)
}
