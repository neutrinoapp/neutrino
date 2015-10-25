package client

import (
	"github.com/go-neutrino/neutrino/log"
	"time"
)

type Client struct {
	connect     func() (interface{}, error)
	isConnected bool
	connection  interface{}
	error       chan error

	Error   chan error
	Message chan string
	Addr    string
}

func NewClient(connect func() (interface{}, error), addr string) *Client {
	return &Client{
		connect:     connect,
		isConnected: false,
		Addr:        addr,
		Message:     make(chan string),
		Error:       make(chan error),
		error:       make(chan error),
	}
}

func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) Disconnected() {
	c.isConnected = false
}

func (c *Client) Connect() {
	for {
		conn, err := c.connect()
		if err == nil {
			c.isConnected = true
			c.connection = conn
			log.Info("Connected to", c.Addr)
			break
		} else {
			log.Error("Error connecting to", c.Addr, err)
			time.Sleep(time.Second * 5)
		}
	}
}
