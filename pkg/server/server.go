package server

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// Server defines server
type Server struct {
	cfg      *Config
	listener net.Listener
}

// Config for server
type Config struct {
	Addr string
}

// New creates a Server instance
func New(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// Start starts server
func (s *Server) Start() error {

	listener, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err
	}

	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.WithError(err).Warn("connection error")
			continue
		}

		go s.handleConnection(conn)
	}

}

// Shutdown closes the tcp socket
func (s *Server) Shutdown() error {
	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}

func (s *Server) handleConnection(conn net.Conn) {
	fmt.Println("asdasd")
}
