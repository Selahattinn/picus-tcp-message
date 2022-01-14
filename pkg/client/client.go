package client

import (
	"net"
)

// Client defines Client
type Client struct {
	cfg *Config

	conn net.Conn
}

// Config for client
type Config struct {
	ServerAddr string
}

// New creates a Client instance
func New(cfg *Config) (*Client, error) {

	conn, err := net.Dial("tcp", cfg.ServerAddr)
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:  cfg,
		conn: conn,
	}, nil
}
