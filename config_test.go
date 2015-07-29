package tempredis

import "testing"

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
