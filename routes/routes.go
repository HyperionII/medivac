package routes

import (
	"github.com/hyperionii/medivac/config"
	"github.com/hyperionii/medivac/handlers/auth"
)

// Routes contains all routes for the application.
type Routes struct {
	APIRoutes []Route
}

// NewRoutes creates a new Router instance and initializes all routes.
func NewRoutes(cfg *config.Config) *Routes {
	authHandler := auth.NewHandler(cfg)

	r := &Routes{
		APIRoutes: []Route{
			&route{
				pattern:       "/auth/login/",
				method:        "GET",
				handlerFunc:   authHandler.Login,
				requiresAuth:  false,
				requiredRoles: []string{},
			},
		},
	}

	return r
}
