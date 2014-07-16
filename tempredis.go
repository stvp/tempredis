package tempredis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	// The presence of this string in redis-server's stdout stream indicates that
	// the server has successfully stood up.
	RedisStartupSuccess = "The server is now ready to accept connections"

	// Duration before returning a timeout error while waiting for redis-server
	// to start. You'll want to lengthen this significantly if you are loading a
	// large backup from disk.
	RedisStartupTimeout = time.Second
)

// Server encapsulates the starting, configuration, and stopping a single
// redis-server process.
type Server struct {
	Config Config
	Stdout io.Reader
	Stderr io.Reader
	cmd    *exec.Cmd
}

// Config is a simple map of config keys to config values. These config values
// will be fed to redis-server on startup.
type Config map[string]string

// Address returns the dial-able address for a Redis server configured with
// this Config struct.
func (c Config) Address() string {
	return c.Bind() + ":" + c.Port()
}

// Bind returns the local bind interface for a Redis server configured with
// this Config struct.
func (c Config) Bind() string {
	bind, ok := c["bind"]
	if !ok {
		bind = "0.0.0.0"
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

// Start a new redis-server with the given config and return both the Server
// and the Start() error, if any.
func Start(config Config) (server *Server, err error) {
	server = &Server{Config: config}
	err = server.Start()
	return server, err
}

// Start will attempt to start and configure redis-server. If the startup fails
// for any reason, an error will be returned and the redis-server process will
// be stopped.
func (s *Server) Start() (err error) {
	if s.cmd != nil {
		return fmt.Errorf("redis-server has already been started")
	}

	s.cmd = exec.Command("redis-server", "-")

	stdin, err := s.cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Hook up Stdout and Stderr readers
	stdoutCopy := new(bytes.Buffer)
	s.Stdout = stdoutCopy
	stdout = TeeReadCloser(stdout, stdoutCopy)
	s.Stderr = stderr

	// Start server and write config
	if err = s.cmd.Start(); err != nil {
		return err
	}
	if err = writeConfig(s.Config, stdin); err != nil {
		s.Term()
		return err
	}

	// Wait until Redis can accept connections
	if err = waitForSuccessfulStartup(stdout); err != nil {
		s.Term()
		return err
	}

	return nil
}

// Term sends a TERM signal to redis-server, if running. It returns an error if
// redis-server isn't running or if redis-server fails to exit.
func (s *Server) Term() (err error) {
	if s.cmd == nil {
		return fmt.Errorf("redis-server is not running")
	}

	s.cmd.Process.Signal(syscall.SIGTERM)
	_, err = s.cmd.Process.Wait()
	if err != nil {
		return err
	}

	s.cmd = nil
	return nil
}

// Kill sends a KILL signal to redis-server, if running. It returns an error if
// redis-server isn't running or if redis-server fails to exit.
func (s *Server) Kill() (err error) {
	if s.cmd == nil {
		return fmt.Errorf("redis-server is not running")
	}

	s.cmd.Process.Signal(syscall.SIGKILL)
	_, err = s.cmd.Process.Wait()
	if err != nil {
		return err
	}

	s.cmd = nil
	return nil
}

func writeConfig(config Config, w io.WriteCloser) (err error) {
	for key, value := range config {
		_, err = fmt.Fprintf(w, "%s %s\n", key, value)
		if err != nil {
			return err
		}
	}
	return w.Close()
}

func waitForSuccessfulStartup(r io.Reader) (err error) {
	scanner := bufio.NewScanner(r)
	line := ""

	success := make(chan bool, 1)
	failure := make(chan bool, 1)
	stopWaiting := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-stopWaiting: // Timeout
				return
			default:
				if scanner.Scan() {
					line = scanner.Text()
					if strings.Contains(line, RedisStartupSuccess) {
						success <- true
						return
					}
				} else {
					failure <- true
					return
				}
			}
		}
	}()

	select {
	case <-success:
		return nil
	case <-failure:
		if err = scanner.Err(); err != nil {
			return fmt.Errorf("Couldn't read redis-server's stdout: %s", err.Error())
		} else {
			return fmt.Errorf("redis-server failed to start up: %s", line)
		}
	case <-time.After(RedisStartupTimeout):
		stopWaiting <- true
		return fmt.Errorf("Timed-out waiting for redis-server to start up successfully. Last line received was: %s", line)
	}
}

// Temp runs a server with the given Config for the duration of the
// given function. If there is an error starting up the server, the
// error will be passed to the given function.
func Temp(config Config, fn func(err error)) {
	server, err := Start(config)
	defer server.Term()
	fn(err)
}
