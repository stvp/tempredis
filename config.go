package tempredis

import "fmt"

// Config is a key-value map of Redis config settings.
type Config map[string]string

// Host returns the host for a Redis server configured with this Config as
// "host:port".
func (c Config) Host() string {
	bind, ok := c["bind"]
	if !ok {
		bind = "127.0.0.1"
	}

	port, ok := c["port"]
	if !ok {
		port = "6379"
	}

	return fmt.Sprintf("%s:%s", bind, port)
}

// URL returns a Redis URL for a Redis server configured with this Config.
func (c Config) URL() string {
	password := c.Password()
	if len(password) == 0 {
		return fmt.Sprintf("redis://%s/", c.Host())
	} else {
		return fmt.Sprintf("redis://:%s@%s/", password, c.Host())
	}
}

// Password returns the password for a Redis server configured with this
// Config. If the server doesn't require authentication, an empty string will
// be returned.
func (c Config) Password() string {
	return c["requirepass"]
}
