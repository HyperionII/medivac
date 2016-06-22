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
	server := &Server{}

	log.Println("Configuring server..")
	server.cfg = config.NewConfig()
	server.cfg.Print()

	if err := server.cfg.Validate(); err != nil {
		return err
	}

	return server
}

// ListenAndServe attaches the current server to the specified configuration port.
func (s *Server) ListenAndServe() error {
	port := ":" + s.cfg.App.Port

	return http.ListenAndServe(port, nil)
}
