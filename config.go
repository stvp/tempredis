package tempredis

import "fmt"

// Config is a simple map of config keys to config values. These config values
// will be fed to redis-server on startup.
type Config map[string]string

// Address returns the dial-able address for a Redis server configured with
// this Config struct.
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
		return fmt.Sprintf("redis://%s", c.Host())
	} else {
		return fmt.Sprintf("redis://:%s@%s", password, c.Host())
	}
}

// Password returns the required password for a Redis server configured with
// this Config struct.  If the server doesn't require authentication, an empty
// string will be returned.
func (c Config) Password() string {
	return c["requirepass"]
}
