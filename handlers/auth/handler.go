package auth

import (
	"net/http"

	"github.com/hyperionii/config"
	"github.com/hyperionii/handlers/httputils"
)

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

// Login does basic account/password login.
func (h *handler) Login(w http.ResponseWriter, r *httpRequest) error {
	var payload struct {
		User     string
		Password string
	}

	err := httputils.DecodeJSON(r.Body, &payload)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "")
		return nil
	}

	return nil
}
