package tempredis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const (
	// ready is the string redis-server prints to stdout after starting
	// successfully.
	ready = "The server is now ready to accept connections"
)

// Server encapsulates the configuration, starting, and stopping of a single
// redis-server process.
type Server struct {
	config    Config
	cmd       *exec.Cmd
	stdout    io.Reader
	stdoutBuf bytes.Buffer
	stderr    io.Reader
}

// Start initiates a new redis-server process configured with the given
// configuration. If "port" is not specified in the config, redis-server will
// bind to a free port. If config is nil, redis-server will use the default
// config and bind to a free port. An error is returned if redis-server is
// unable to successfully start for any reason.
func Start(config Config) (server *Server, err error) {
	if config == nil {
		config = Config{}
	}

	port, ok := config["port"]
	if !ok {
		port, err = ephemeralPort()
		if err != nil {
			return nil, err
		}
		config["port"] = port
	}

	server = &Server{config: config}
	err = server.start()
	if err != nil {
		return server, err
	}

	// Block until Redis is ready to accept connections.
	err = server.waitFor(ready)

	return server, err
}

func (s *Server) start() (err error) {
	if s.cmd != nil {
		return fmt.Errorf("redis-server has already been started")
	}

	s.cmd = exec.Command("redis-server", "-")

	stdin, _ := s.cmd.StdinPipe()
	s.stdout, _ = s.cmd.StdoutPipe()
	s.stderr, _ = s.cmd.StderrPipe()

	err = s.cmd.Start()
	if err == nil {
		err = writeConfig(s.config, stdin)
	}

	return err
}

func writeConfig(config Config, w io.WriteCloser) (err error) {
	for key, value := range config {
		if value == "" {
			value = "\"\""
		}
		_, err = fmt.Fprintf(w, "%s %s\n", key, value)
		if err != nil {
			return err
		}
	}
	return w.Close()
}

// waitFor blocks until redis-server prints the given string to stdout.
func (s *Server) waitFor(search string) (err error) {
	var line string

	scanner := bufio.NewScanner(s.stdout)
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Fprintf(&s.stdoutBuf, "%s\n", line)
		if strings.Contains(line, search) {
			return nil
		}
	}
	err = scanner.Err()
	if err == nil {
		err = io.EOF
	}
	return err
}

// URL returns a dial-able URL for this Redis server process.
func (s *Server) URL() *url.URL {
	return s.config.URL()
}

// Stdout blocks until redis-server returns and then returns the full stdout
// output.
func (s *Server) Stdout() string {
	io.Copy(&s.stdoutBuf, s.stdout)
	return s.stdoutBuf.String()
}

// Stderr blocks until redis-server returns and then returns the full stdout
// output.
func (s *Server) Stderr() string {
	bytes, _ := ioutil.ReadAll(s.stderr)
	return string(bytes)
}

// Term gracefully shuts down redis-server. It returns an error if redis-server
// fails to terminate.
func (s *Server) Term() (err error) {
	return s.signal(syscall.SIGTERM)
}

// Kill forcefully shuts down redis-server. It returns an error if redis-server
// fails to die.
func (s *Server) Kill() error {
	return s.signal(syscall.SIGKILL)
}

func (s *Server) signal(sig syscall.Signal) error {
	s.cmd.Process.Signal(sig)
	_, err := s.cmd.Process.Wait()
	return err
}

// ephemeralPort returns a local ephemeral TCP port that we can bind to.
func ephemeralPort() (port string, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}

	listener.Close()
	return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), nil
}
