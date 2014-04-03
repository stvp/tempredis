package tempredis

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

const (
	PORT = "7725"
)

func startServer(config Config) (*Server, error) {
	server := &Server{Config: config}
	err := server.Start()
	return server, err
}

// -- Tests

func TestConfigAddress(t *testing.T) {
	config := Config{}
	if config.Address() != "0.0.0.0:6379" {
		t.Errorf("Expected: %#v, got: %#v", "0.0.0.0:6379", config.Address())
	}
	if config.Password() != "" {
		t.Errorf("Expected: %#v, got: %#v", "", config.Address())
	}
	config["bind"] = "127.0.0.1"
	config["port"] = "1234"
	config["requirepass"] = "pw"
	if config.Address() != "127.0.0.1:1234" {
		t.Errorf("Expected: %#v, got: %#v", "127.0.0.1:1234", config.Address())
	}
	if config.Password() != "pw" {
		t.Errorf("Expected: %#v, got: %#v", "pw", config.Address())
	}
}

func TestRedisServerStartAndStop(t *testing.T) {
	server, err := startServer(Config{"port": PORT, "databases": "3"})
	if err != nil {
		t.Fatalf("Creating a server failed: %s", err.Error())
	}

	r, err := redis.Dial("tcp", ":"+PORT)
	defer r.Close()
	if err != nil {
		t.Fatalf("Couldn't connect to running server", err.Error())
	}
	databases, err := redis.Strings(r.Do("CONFIG", "GET", "databases"))
	if err != nil {
		t.Fatalf("CONFIG GET failed on running server: %s", err.Error())
	}
	if databases[1] != "3" {
		t.Fatalf("Config wasn't properly set. Should have got 3, but got %s", databases)
	}

	if err := server.Term(); err != nil {
		t.Fatalf("Stopping a running server failed: %s", err.Error())
	}
	if err := server.Term(); err == nil {
		t.Fatal("Stopping an already stopped server should fail")
	}
}

func TestRedisServerTerm(t *testing.T) {
	server := Server{Config: Config{"port": PORT}}
	if err := server.Term(); err == nil {
		t.Fatal("Term() on a server that isn't running should fail")
	}

	err := server.Start()
	if err != nil {
		t.Fatalf("Starting a server failed: %s", err.Error())
	}
	if err := server.Term(); err != nil {
		t.Fatalf("Failed to TERM server: %s", err.Error())
	}

	_, err = redis.Dial("tcp", ":"+PORT)
	if err == nil {
		t.Fatal("Server is running, but it shouldn't be")
	}
}

func TestRedisServerKill(t *testing.T) {
	server := Server{Config: Config{"port": PORT}}
	if err := server.Kill(); err == nil {
		t.Fatal("Kill() on a server that isn't running should fail")
	}

	err := server.Start()
	if err != nil {
		t.Fatalf("Starting a server failed: %s", err.Error())
	}

	// Block server with sleep
	r, err := redis.Dial("tcp", ":"+PORT)
	defer r.Close()
	if err != nil {
		t.Fatalf("Couldn't connect to running server", err.Error())
	}
	go r.Do("DEBUG", "SLEEP", "10")

	if err := server.Kill(); err != nil {
		t.Fatalf("Failed to KILL server: %s", err.Error())
	}

	_, err = redis.Dial("tcp", ":"+PORT)
	if err == nil {
		t.Fatal("Server is running, but it shouldn't be")
	}
}

func TestRedisServerStartFailure(t *testing.T) {
	s, err := startServer(Config{"port": PORT})
	defer s.Term()
	if err != nil {
		t.Fatalf("Creating a server failed: %s", err.Error())
	}

	server, err := startServer(Config{"port": PORT})
	defer server.Term()
	if err == nil {
		t.Fatal("Exptected server to fail starting (port in use), but it didn't")
	}

	if err := server.Term(); err == nil {
		t.Fatal(err)
	}
}
