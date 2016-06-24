package auth

import (
	"net/http"

	"github.com/hyperionii/medivac/handlers/httputils"
)

// Login does basic account/password login.
func (h *handler) Login(w http.ResponseWriter, r *http.Request) error {
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
