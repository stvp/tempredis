package tempredis

import (
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestServer(t *testing.T) {
	server, err := Start(Config{"databases": "3"})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	r, err := redis.Dial("unix", server.Socket())
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	databases, err := redis.Strings(r.Do("CONFIG", "GET", "databases"))
	if err != nil {
		t.Fatal(err)
	}
	if databases[1] != "3" {
		t.Fatalf("databases config should be 3, but got %s", databases)
	}

	if err := server.Term(); err != nil {
		t.Fatal(err)
	}
	if err := server.Term(); err == nil {
		t.Fatal("stopping an already stopped server should fail")
	}
}

func TestStartWithDefaultConfig(t *testing.T) {
	server, err := Start(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	r, err := redis.Dial("unix", server.Socket())
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	_, err = r.Do("PING")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartFail(t *testing.T) {
	server, err := Start(Config{"oops": "borked"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	defer server.Kill()
}
