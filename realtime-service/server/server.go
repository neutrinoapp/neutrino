package server

import (
	"fmt"
	"github.com/go-neutrino/neutrino-config"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//allow connections from any origin
			return true
		},
	}
	config *viper.Viper
)

func Initialize(c *viper.Viper) {
	config = c

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
			return
		}

		//TODO: token authentication
		token := r.URL.Query().Get("token") //the hash is a unique, per user
		realtimeConn := NewConnection(conn, token)

		GetConnectionStore().Put(token, realtimeConn)
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got message!")

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	})

	nc, err := nats.Connect(config.GetString(nconfig.KEY_QUEUE_HOST))
	if err != nil {
		panic(err)
	}

	conn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
		return
	}

	conn.Subscribe("realtime-jobs", func(s string) {
		fmt.Printf("Received a message: %s\n", s)
	})

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println("Sending message!")
			conn.Publish("realtime-jobs", "Hello World")
		}
	}()
}
