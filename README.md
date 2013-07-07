tempredis
=========

Tempredis makes it easy to start and stop temporary `redis-server`
processes with custom configs for testing.

[API documentation](http://godoc.org/github.com/stvp/tempredis)

Example
-------

```go
package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/stvp/tempredis"
)

func main() {
	server := tempredis.Server{
		Config: tempredis.Config{
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
```

Should I use this outside of testing?
-------------------------------------

No.

