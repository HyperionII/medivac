package main

import (
	"log"
	"net/http"

	"github.com/hyperionii/medivac/config"
)

// Server contains instance configuration variables for the server.
type Server struct {
	cfg *config.Config
}

// NewServer returns a new instance of the server.
func NewServer() (*Server, error) {
	log.Println("Configuring server..")
	cfg, err := config.NewConfig()
	cfg.Print()
	if err != nil {
		return nil, err
	}

	server := &Server{
		cfg: cfg,
	}
	return server, nil
}

// ListenAndServe attaches the current server to the specified configuration port.
func (s *Server) ListenAndServe() error {
	port := ":" + s.cfg.App.Port

	return http.ListenAndServe(port, nil)
}
