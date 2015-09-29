package tempredis

// Config is a simple map of config keys to config values. These config values
// will be fed to redis-server on startup.
type Config map[string]string

// Address returns the dial-able address for a Redis server configured with
// this Config struct.
func (c Config) Host() string {
	return c.Bind() + ":" + c.Port()
}

// Bind returns the local bind interface for a Redis server configured with
// this Config struct.
func (c Config) Bind() string {
	bind, ok := c["bind"]
	if !ok {
		bind = "127.0.0.1"
	}
	return bind
}

// Port returns the local port for a Redis server configured with this Config
// struct.
func (c Config) Port() string {
	port, ok := c["port"]
	if !ok {
		port = "6379"
	}
	return port
}

// Password returns the required password for a Redis server configured with
// this Config struct.  If the server doesn't require authentication, an empty
// string will be returned.
func (c Config) Password() string {
	return c["requirepass"]
}
