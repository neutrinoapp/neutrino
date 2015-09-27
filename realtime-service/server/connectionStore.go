package server

import (
	"sync"
)

var c ConnectionStore
type ConnectionStore interface {
	Get(string) []RealtimeConnection
	Put(string, RealtimeConnection)
	Remove(string, RealtimeConnection)
}

type connectionStore struct {
	connections map[string][]RealtimeConnection
	l *sync.Mutex
}

func (c *connectionStore) Get(g string) []RealtimeConnection {
	c.l.Lock()
	defer c.l.Unlock()

	return c.connections[g]
}

func (c *connectionStore) Put(g string, conn RealtimeConnection) {
	c.l.Lock()
	defer c.l.Unlock()

	c.connections[g] = append(c.connections[g], conn)
}

func (c *connectionStore) Remove(g string, conn RealtimeConnection) {
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
		c  = &connectionStore{
			connections: make(map[string][]RealtimeConnection),
			l: new(sync.Mutex),
		}
	}

	return c
}
