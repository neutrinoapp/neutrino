package server

import (
	"sync"
)

var c ConnectionStore

type ConnectionStore interface {
	Get(string) []ClientConnection
	Put(string, ClientConnection)
	Remove(string, ClientConnection)
}

type connectionStore struct {
	connections map[string][]ClientConnection
	l           *sync.Mutex
}

func (c *connectionStore) Get(g string) []ClientConnection {
	c.l.Lock()
	defer c.l.Unlock()

	return c.connections[g]
}

func (c *connectionStore) Put(g string, conn ClientConnection) {
	c.l.Lock()
	defer c.l.Unlock()

	c.connections[g] = append(c.connections[g], conn)
}

func (c *connectionStore) Remove(g string, conn ClientConnection) {
	c.l.Lock()
	defer c.l.Unlock()

	connections := c.connections[g]

	connIndex := -1
	for i, v := range connections {
		if v == conn {
			connIndex = i
		}
	}

	if connIndex != -1 {
		c.connections[g] = append(connections[:connIndex], connections[connIndex+1:]...)
	}
}

func GetConnectionStore() ConnectionStore {
	if c == nil {
		c = &connectionStore{
			connections: make(map[string][]ClientConnection),
			l:           new(sync.Mutex),
		}
	}

	return c
}
