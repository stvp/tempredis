package tempredis

import "testing"

func TestConfigAddress(t *testing.T) {
	config := Config{}
	if config.Host() != "127.0.0.1:6379" {
		t.Errorf("Expected: %#v, got: %#v", "127.0.0.1:6379", config.Host())
	}
	if config.Password() != "" {
		t.Errorf("Expected: %#v, got: %#v", "", config.Host())
	}
	config["bind"] = "127.0.0.10"
	config["port"] = "1234"
	config["requirepass"] = "pw"
	if config.Host() != "127.0.0.10:1234" {
		t.Errorf("Expected: %#v, got: %#v", "127.0.0.10:1234", config.Host())
	}
	if config.Password() != "pw" {
		t.Errorf("Expected: %#v, got: %#v", "pw", config.Password())
	}
}
