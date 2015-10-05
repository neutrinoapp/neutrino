package main

import (
	"github.com/go-neutrino/neutrino-config"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	c        *viper.Viper
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//allow connections from any origin
			return true
		},
	}
	connections []*websocket.Conn
)

func jobsHandler(m *nats.Msg) {
	log.Println("Got message " + string(m.Data))

	for i, c := range connections {
		log.Println("Sending message to connection: " + strconv.Itoa(i+1))
		c.WriteMessage(websocket.TextMessage, m.Data)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c = nconfig.Load()
	qAddr := c.GetString(nconfig.KEY_QUEUE_ADDR)
	n, e := nats.Connect(qAddr)

	if e != nil {
		panic(e)
	}

	conn, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	conn.Subscribe("realtime-jobs", jobsHandler)
	log.Println("Connected to NATS on " + qAddr)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
			return
		}

		connections = append(connections, conn)
	})

	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Println("Sending message!")
			conn.Publish("realtime-jobs", "Hello World")
		}
	}()

	port := c.GetString(nconfig.KEY_BROKER_PORT)
	log.Println("Starting WS service on port " + port)
	http.ListenAndServe(port, nil)
}
