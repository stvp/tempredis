package main

import (
	"github.com/garyburd/redigo/redis"
)

func ExampleUsage() {
	server := Server{
		Config: Config{
			"port":      "11001",
			"databases": "8",
		},
	}
	defer server.Stop()
	if err := server.Start(); err != nil {
		panic(err)
	}

	conn, err := redis.Dial("tcp", ":11001")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Do("SET", "foo", "bar")
}
