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
	server, err := tempredis.Start(tempredis.Config{"databases": "8"})
	if err != nil {
		panic(err)
	}
	defer server.Term()

	server.WaitFor(tempredis.Ready)

	conn, err := redis.Dial("tcp", server.Config.Host())
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Do("SET", "foo", "bar")
}
```

If you don't care about normal shutdown behavior or want to simulate a crash,
you can send a KILL signal to the server with:

```go
server.Kill()
```

Should I use this outside of testing?
-------------------------------------

No.

