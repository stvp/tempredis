package tempredis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const (
	// RDBLoaded is the string redis-server prints to stdout after the RDB file
	// has been loaded successfully.
	RDBLoaded = "DB loaded from disk"

	// AOFLoaded is the string redis-server prints to stdout after the AOF file
	// has been loaded successfully.
	AOFLoaded = "DB loaded from append only file"

	// Ready is the string redis-server prints to stdout after starting
	// successfully.
	Ready = "The server is now ready to accept connections"
)

// Server encapsulates the starting, configuration, and stopping of a single
// redis-server process.
type Server struct {
	Config    Config
	cmd       *exec.Cmd
	stdout    io.Reader
	stdoutBuf bytes.Buffer
	stderr    io.Reader
}

// Start starts a redis-server process with the given config and returns a
// Server object. If the server failed to start, Start will return an error.
func Start(config Config) (server *Server, err error) {
	port, ok := config["port"]
	if !ok {
		port = ephemeralPort()
		config["port"] = port
	}
	server = &Server{Config: config}
	err = server.start()
	return server, err
}

// Stdout blocks until redis-server returns and then returns the full stdout.
func (s *Server) Stdout() string {
	io.Copy(&s.stdoutBuf, s.stdout)
	return s.stdoutBuf.String()
}

// Stderr blocks until redis-server returns and then returns the full stdout.
func (s *Server) Stderr() string {
	bytes, _ := ioutil.ReadAll(s.stderr)
	return string(bytes)
}

// WaitFor blocks until redis-server prints the given string to stdout.
func (s *Server) WaitFor(search string) (err error) {
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
		err = writeConfig(s.Config, stdin)
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

func (s *Server) signal(sig syscall.Signal) error {
	s.cmd.Process.Signal(sig)
	_, err := s.cmd.Process.Wait()
	return err
}

// ephemeralPort returns a local ephemeral TCP port that we can bind to.
func ephemeralPort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	listener.Close()

	return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
}
