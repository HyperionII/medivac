package auth

import (
	"net/http"

	"github.com/hyperionii/medivac/config"
)

type Handler interface {
	Login(http.ResponseWriter, *http.Request) error
}

// Handler structure for the auth handler.
type handler struct {
	cfg *config.Config
}

// NewHandler initializes an auth handler struct.
func NewHandler(cfg *config.Config) Handler {
	return &handler{
		cfg: cfg,
	}
}
