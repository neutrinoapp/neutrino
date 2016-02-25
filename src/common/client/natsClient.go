package client

import (
	"errors"
	"time"

	"github.com/nats-io/nats"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/log"
)

type NatsClient struct {
	*Client
}

func NewNatsClient(addr string) *NatsClient {
	connect := func() (interface{}, error) {
		log.Info("Connecting to nats:", addr)
		n, err := nats.Connect(addr)

		if err != nil {
			return nil, err
		}

		conn, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	c := NewClient(connect, addr)

	natsClient := &NatsClient{c}
	natsClient.handleConnection()
	return natsClient
}

func (n *NatsClient) handleConnection() {
	go func() {
		for {
			conn := n.GetConnection()
			if conn != nil {
				//TODO: this flush does not seem to work as a proper ping
				err := conn.Flush()
				if err != nil {
					n.Error <- err
					n.Connect()
					continue
				}
			}

			time.Sleep(time.Second * 5)
		}
	}()

	go n.Connect()
}

func (n *NatsClient) Subscribe(c string, cb nats.Handler) error {
	for {
		if n.GetConnection() != nil {
			_, err := n.GetConnection().Subscribe(config.Get(config.CONST_REALTIME_JOBS_SUBJ), cb)
			return err
		} else {
			time.Sleep(time.Second * 1)
		}
	}
}

func (n *NatsClient) Publish(t string, msg interface{}) error {
	conn := n.GetConnection()
	if conn == nil {
		return errors.New("No available connection.")
	}

	return conn.Publish(t, msg)
}

func (n *NatsClient) GetConnection() *nats.EncodedConn {
	if n.connection == nil {
		return nil
	}

	return n.connection.(*nats.EncodedConn)
}
