package tempredis

import (
	"fmt"
	"net/url"
)

// Config is a key-value map of Redis config settings.
type Config map[string]string

// URL returns the dial-able URL for a Redis server configured with this
// Config.
func (c Config) URL() (redisURL *url.URL) {
	bind, ok := c["bind"]
	if !ok {
		// Redis' default bind
		bind = "0.0.0.0"
	}

	port, ok := c["port"]
	if !ok {
		// Redis' default port
		port = "6379"
	}

	redisURL = &url.URL{
		Scheme: "redis",
		Host:   fmt.Sprintf("%s:%s", bind, port),
	}

	if len(c["requirepass"]) > 0 {
		redisURL.User = url.UserPassword("", c["requirepass"])
	}

	return redisURL
}
