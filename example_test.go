package tempredis

import (
	"github.com/garyburd/redigo/redis"
)

func ExampleUsage() {
	server, err := Start(Config{"databases": "8"})
	if err != nil {
		panic(err)
	}
	defer server.Term()

	server.WaitFor(Ready)

	conn, err := redis.Dial("tcp", server.Config.Host())
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Do("SET", "foo", "bar")
}
