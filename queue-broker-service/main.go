package main

import (
	"github.com/go-neutrino/neutrino/log"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
	"net/http"
	"strconv"
	"github.com/go-neutrino/neutrino/config"
	"github.com/go-neutrino/neutrino/client"
)

var (
	upgrader = client.NewWebsocketUpgrader()
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
	qAddr := config.Get(config.KEY_QUEUE_ADDR)
	n, e := nats.Connect(qAddr)

	if e != nil {
		panic(e)
	}

	conn, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	conn.Subscribe("realtime-jobs", jobsHandler)
	log.Info("Connected to NATS on " + qAddr)

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
