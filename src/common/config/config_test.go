package config

import (
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	if Get(KEY_QUEUE_ADDR) != "nats://localhost:4222" {
		t.Fatal("Incorrect key", KEY_QUEUE_ADDR)
	}

	if Get(KEY_RETHINK_KEY) != "localhost:27017" {
		t.Fatal("Incorrect key", KEY_RETHINK_KEY)
	}

	if Get(KEY_API_PORT) != ":5000" {
		t.Fatal("Incorrect key", KEY_API_PORT)
	}

	if Get(KEY_REALTIME_PORT) != ":6000" {
		t.Fatal("Incorrect key", KEY_REALTIME_PORT)
	}
}
