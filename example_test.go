package tempredis

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
	if err := server.Start(); err != nil {
		panic(err)
	}
	defer server.Stop()

	conn, err := redis.Dial("tcp", ":11001")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Do("SET", "foo", "bar")
}
