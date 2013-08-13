package tempredis

import (
	"github.com/garyburd/redigo/redis"
)

func ExampleUsage() {
	server, err := Start(
		Config{
			"port":      "11001",
			"databases": "8",
		},
	)
	if err := server.Start(); err != nil {
		panic(err)
	}
	defer server.Term()

	conn, err := redis.Dial("tcp", ":11001")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Do("SET", "foo", "bar")
}
