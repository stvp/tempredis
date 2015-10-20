package tempredis

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestServer(t *testing.T) {
	server, err := Start(Config{"databases": "3"})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	err = server.WaitFor(Ready)
	if err != nil {
		t.Fatal(err)
	}

	r, err := redis.DialURL(server.URL().String())
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

	err = server.WaitFor(Ready)
	if err != nil {
		t.Fatal(err)
	}

	r, err := redis.DialURL(server.URL().String())
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	_, err = r.Do("PING")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWaitForStdoutFail(t *testing.T) {
	server, err := Start(Config{"oops": "borked"})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	err = server.WaitFor(Ready)
	if err == nil {
		t.Fatal(err)
	}
}

func TestRDBLoaded(t *testing.T) {
	server, err := Start(Config{"dbfilename": "_dump.rdb"})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	err = server.WaitFor(RDBLoaded)
	if err != nil {
		fmt.Println(err)
	}
}

func TestAOFLoaded(t *testing.T) {
	server, err := Start(Config{"appendonly": "yes", "appendfilename": "_appendonly.aof"})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	err = server.WaitFor(AOFLoaded)
	if err != nil {
		fmt.Println(err)
	}
}
