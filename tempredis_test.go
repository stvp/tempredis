package main

import (
	"testing"
)

func TestRedisServerStart(t *testing.T) {
	server := NewServer(Config{"port": "6379"})
	err := server.Start()
	if err != nil {
		panic(err)
	}
	// Test successful command
	err = server.Stop()
	if err != nil {
		panic(err)
	}
}

// Test failed startup

// Test proper stop
